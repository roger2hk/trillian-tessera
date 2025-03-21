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

package tessera

import (
	"bytes"
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/transparency-dev/trillian-tessera/api/layout"
	"github.com/transparency-dev/trillian-tessera/client"
	"golang.org/x/sync/errgroup"
	"k8s.io/klog/v2"
)

// copier controls the migration work.
type copier struct {
	// storage is the target we're migrating to.
	storage    MigrationStorage
	getEntries client.EntryBundleFetcherFunc

	// sourceSize is the size of the source log.
	sourceSize uint64
	// sourceRoot is the root hash of the source log at sourceSize.
	sourceRoot []byte

	// todo contains work items to be completed.
	todo chan bundle

	// bundlesMigrated is the number of entry bundles migrated so far.
	bundlesMigrated atomic.Uint64
}

// bundle represents the address of an individual entry bundle.
type bundle struct {
	Index   uint64
	Partial uint8
}

// MigrationStorage describes the required functionality from the target storage driver.
//
// It's expected that the implementation of this interface will attempt to integrate the entry bundles
// being set as soon as is reasonably possible. It's up to the implementation whether that's done as a
// background task, or is in-line with the call to AwaitIntegration.
type MigrationStorage interface {
	// SetEntryBundle is called to store the provided entry bundle bytes at the given coordinates.
	//
	// Implementations SHOULD treat calls to this function as being idempotent - i.e. attempts to set a previously
	// set entry bundle should succeed if the bundle data are identical.
	//
	// This will be called as many times as necessary to set all entry bundles being migrated, quite likely in parallel.
	SetEntryBundle(ctx context.Context, index uint64, partial uint8, bundle []byte) error
	// AwaitIntegration should block until the storage driver has received and integrated all outstanding entry bundles implied by sourceSize,
	// and return the locally calculated root hash.
	// An error should be returned if there is a problem integrating.
	AwaitIntegration(ctx context.Context, sourceSize uint64) ([]byte, error)
	// IntegratedSize returns the current integrated size of the local tree.
	IntegratedSize(ctx context.Context) (uint64, error)
}

// Migrate starts the work of copying sourceSize entries from the source to the target log.
//
// Only the entry bundles are copied as the target storage is expected to integrate them and recalculate the root.
// This is done to ensure the correctness of both the source log as well as the copy process itself.
//
// A call to this function will block until either the copying is done, or an error has occurred.
// It is an error if the resource copying completes ok but the resulting root hash does not match the provided sourceRoot.
func Migrate(ctx context.Context, numWorkers int, sourceSize uint64, sourceRoot []byte, getEntries client.EntryBundleFetcherFunc, storage MigrationStorage) error {
	klog.Infof("Starting migration; source size %d root %x", sourceSize, sourceRoot)

	m := &copier{
		storage:    storage,
		sourceSize: sourceSize,
		sourceRoot: sourceRoot,
		getEntries: getEntries,
		todo:       make(chan bundle, numWorkers),
	}

	// init
	targetSize, err := m.storage.IntegratedSize(ctx)
	if err != nil {
		return fmt.Errorf("size: %v", err)
	}
	if targetSize > sourceSize {
		return fmt.Errorf("target size %d > source size %d", targetSize, sourceSize)
	}

	go m.populateWork(targetSize, sourceSize)

	// Print stats
	go func() {
		bundlesToMigrate := (sourceSize / layout.EntryBundleWidth) - (targetSize / layout.EntryBundleWidth) + 1
		if bundlesToMigrate == 0 {
			return
		}
		for {
			time.Sleep(time.Second)
			bn := m.bundlesMigrated.Load()
			bnp := float64(bn*100) / float64(bundlesToMigrate)
			s, err := m.storage.IntegratedSize(ctx)
			if err != nil {
				klog.Warningf("Size: %v", err)
			}
			intp := float64(s*100) / float64(sourceSize)
			klog.Infof("integration: %d (%.2f%%)  bundles: %d (%.2f%%)", s, intp, bn, bnp)
		}
	}()

	// Do the copying
	eg := errgroup.Group{}
	for range numWorkers {
		eg.Go(func() error {
			return m.migrateWorker(ctx)
		})
	}
	var root []byte
	eg.Go(func() error {
		r, err := m.storage.AwaitIntegration(ctx, sourceSize)
		if err != nil {
			return fmt.Errorf("migration failed: %v", err)
		}
		root = r
		return nil
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}

	if !bytes.Equal(root, sourceRoot) {
		return fmt.Errorf("migration completed, but local root hash %x != source root hash %x", root, sourceRoot)
	}

	klog.Infof("Migration successful.")
	return nil
}

// populateWork sends entries to the `todo` work channel.
// Each entry corresponds to an individual entryBundle which needs to be copied.
func (m *copier) populateWork(from, treeSize uint64) {
	klog.Infof("Spans for entry range [%d, %d)", from, treeSize)
	defer close(m.todo)

	for ri := range layout.Range(from, treeSize-from, treeSize) {
		m.todo <- bundle{Index: ri.Index, Partial: ri.Partial}
	}
}

// migrateWorker undertakes work items from the `todo` channel.
//
// It will attempt to retry failed operations several times before giving up, this should help
// deal with any transient errors which may occur.
func (m *copier) migrateWorker(ctx context.Context) error {
	for b := range m.todo {
		err := retry.Do(func() error {
			d, err := m.getEntries(ctx, b.Index, uint8(b.Partial))
			if err != nil {
				wErr := fmt.Errorf("failed to fetch entrybundle %d (p=%d): %v", b.Index, b.Partial, err)
				klog.Infof("%v", wErr)
				return wErr
			}
			if err := m.storage.SetEntryBundle(ctx, b.Index, b.Partial, d); err != nil {
				wErr := fmt.Errorf("failed to store entrybundle %d (p=%d): %v", b.Index, b.Partial, err)
				klog.Infof("%v", wErr)
				return wErr
			}
			m.bundlesMigrated.Add(1)
			return nil
		},
			retry.Attempts(10),
			retry.DelayType(retry.BackOffDelay))
		if err != nil {
			klog.Infof("retry: %v", err)
			return err
		}
	}
	return nil
}
