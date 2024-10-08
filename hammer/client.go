// Copyright 2024 The Tessera authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/transparency-dev/trillian-tessera/client"
	"k8s.io/klog/v2"
)

var ErrRetry = errors.New("retry")

// newLogClientsFromFlags returns a fetcher and a writer that will read
// and write leaves to all logs in the `log_url` flag set.
func newLogClientsFromFlags() (*roundRobinFetcher, *roundRobinLeafWriter) {
	if len(logURL) == 0 {
		klog.Exitf("--log_url must be provided")
	}

	if len(writeLogURL) == 0 {
		// If no write_log_url is provided, then default it to log_url
		writeLogURL = logURL
	}

	rootUrlOrDie := func(s string) *url.URL {
		// url must reference a directory, by definition
		if !strings.HasSuffix(s, "/") {
			s += "/"
		}
		rootURL, err := url.Parse(s)
		if err != nil {
			klog.Exitf("Invalid log URL: %v", err)
		}
		return rootURL
	}

	fetchers := []client.Fetcher{}
	for _, s := range logURL {
		fetchers = append(fetchers, newFetcher(rootUrlOrDie(s)))
	}
	writers := []httpLeafWriter{}
	for _, s := range writeLogURL {
		addURL, err := rootUrlOrDie(s).Parse("add")
		if err != nil {
			klog.Exitf("Failed to create add URL: %v", err)
		}
		writers = append(writers, httpLeafWriter{u: addURL})
	}
	return &roundRobinFetcher{f: fetchers}, &roundRobinLeafWriter{ws: writers}
}

// newFetcher creates a Fetcher for the log at the given root location.
func newFetcher(root *url.URL) client.Fetcher {
	get := getByScheme[root.Scheme]
	if get == nil {
		panic(fmt.Errorf("unsupported URL scheme %s", root.Scheme))
	}

	return func(ctx context.Context, p string) ([]byte, error) {
		u, err := root.Parse(p)
		if err != nil {
			return nil, err
		}
		return get(ctx, u)
	}
}

var getByScheme = map[string]func(context.Context, *url.URL) ([]byte, error){
	"http":  readHTTP,
	"https": readHTTP,
	"file": func(_ context.Context, u *url.URL) ([]byte, error) {
		return os.ReadFile(u.Path)
	},
}

func readHTTP(ctx context.Context, u *url.URL) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if *bearerToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *bearerToken))
	}

	resp, err := hc.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			klog.Errorf("resp.Body.Close(): %v", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		klog.Infof("Not found: %q", u.String())
		return nil, os.ErrNotExist
	case http.StatusOK:
		break
	default:
		return nil, fmt.Errorf("unexpected http status %q", resp.Status)
	}
	return body, nil
}

// roundRobinFetcher ensures that read requests are sent to all configured fetchers
// using a round-robin strategy.
type roundRobinFetcher struct {
	sync.Mutex
	idx int
	f   []client.Fetcher
}

func (rr *roundRobinFetcher) Fetch(ctx context.Context, path string) ([]byte, error) {
	f := rr.next()
	return f(ctx, path)
}

func (rr *roundRobinFetcher) next() client.Fetcher {
	rr.Lock()
	defer rr.Unlock()

	f := rr.f[rr.idx]
	rr.idx = (rr.idx + 1) % len(rr.f)

	return f
}

type httpLeafWriter struct {
	u *url.URL
}

func (w httpLeafWriter) Write(ctx context.Context, newLeaf []byte) (uint64, error) {
	req, err := http.NewRequest(http.MethodPost, w.u.String(), bytes.NewReader(newLeaf))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}
	if *bearerTokenWrite != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *bearerTokenWrite))
	}
	resp, err := hc.Do(req.WithContext(ctx))
	if err != nil {
		return 0, fmt.Errorf("failed to write leaf: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return 0, fmt.Errorf("failed to read body: %v", err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		if resp.Request.Method != http.MethodPost {
			return 0, fmt.Errorf("write leaf was redirected to %s", resp.Request.URL)
		}
		// Continue below
	case http.StatusServiceUnavailable, http.StatusBadGateway, http.StatusGatewayTimeout:
		// These status codes may indicate a delay before retrying, so handle that here:
		time.Sleep(retryDelay(resp.Header.Get("RetryAfter"), time.Second))

		return 0, fmt.Errorf("log not available. Status code: %d. Body: %q %w", resp.StatusCode, body, ErrRetry)
	default:
		return 0, fmt.Errorf("write leaf was not OK. Status code: %d. Body: %q", resp.StatusCode, body)
	}
	parts := bytes.Split(body, []byte("\n"))
	index, err := strconv.ParseUint(string(parts[0]), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("write leaf failed to parse response: %v", body)
	}
	return index, nil
}

func retryDelay(retryAfter string, defaultDur time.Duration) time.Duration {
	if retryAfter == "" {
		return defaultDur
	}
	d, err := time.Parse(http.TimeFormat, retryAfter)
	if err == nil {
		return time.Until(d)
	}
	s, err := strconv.Atoi(retryAfter)
	if err == nil {
		return time.Duration(s) * time.Second
	}
	return defaultDur
}

// roundRobinLeafWriter ensures that write requests are sent to all configured
// LeafWriters using a round-robin strategy.
type roundRobinLeafWriter struct {
	sync.Mutex
	idx int
	ws  []httpLeafWriter
}

func (rr *roundRobinLeafWriter) Write(ctx context.Context, newLeaf []byte) (uint64, error) {
	w := rr.next()
	return w(ctx, newLeaf)
}

func (rr *roundRobinLeafWriter) next() LeafWriter {
	rr.Lock()
	defer rr.Unlock()

	f := rr.ws[rr.idx]
	rr.idx = (rr.idx + 1) % len(rr.ws)

	return f.Write
}
