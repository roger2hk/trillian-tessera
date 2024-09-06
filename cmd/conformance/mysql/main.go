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

// mysql is a simple personality allowing to run conformance/compliance/performance tests and showing how to use the Tessera MySQL storage implmentation.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tessera "github.com/transparency-dev/trillian-tessera"
	"github.com/transparency-dev/trillian-tessera/api/layout"
	"github.com/transparency-dev/trillian-tessera/storage/mysql"
	"golang.org/x/mod/sumdb/note"
	"k8s.io/klog/v2"
)

var (
	mysqlURI          = flag.String("mysql_uri", "user:password@tcp(db:3306)/tessera", "Connection string for a MySQL database")
	dbConnMaxLifetime = flag.Duration("db_conn_max_lifetime", 3*time.Minute, "")
	dbMaxOpenConns    = flag.Int("db_max_open_conns", 64, "")
	dbMaxIdleConns    = flag.Int("db_max_idle_conns", 64, "")
	initSchemaPath    = flag.String("init_schema_path", "", "Location of the schema file if database initialization is needed")
	listen            = flag.String("listen", ":2024", "Address:port to listen on")
	privateKeyPath    = flag.String("private_key_path", "", "Location of private key file")
	publicKeyPath     = flag.String("public_key_path", "", "Location of public key file")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	ctx := context.Background()

	db := createDatabaseOrDie(ctx)
	noteSigner := createSignerOrDie()
	noteVerifier := createVerifierOrDie()

	// Initialise the Tessera MySQL storage
	storage, err := mysql.New(ctx, db, tessera.WithCheckpointSignerVerifier(noteSigner, noteVerifier))
	if err != nil {
		klog.Exitf("Failed to create new MySQL storage: %v", err)
	}

	// Set up the handlers for the tlog-tiles GET methods, and a custom handler for HTTP POSTs to /add
	configureTilesReadAPI(http.DefaultServeMux, storage)
	http.HandleFunc("POST /add", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		idx, err := storage.Add(r.Context(), tessera.NewEntry(b))()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		if _, err = w.Write([]byte(fmt.Sprintf("%d", idx))); err != nil {
			klog.Errorf("/add: %v", err)
			return
		}
	})

	// Serve HTTP requests until the process is terminated
	if err := http.ListenAndServe(*listen, http.DefaultServeMux); err != nil {
		klog.Exitf("ListenAndServe: %v", err)
	}
}

func createDatabaseOrDie(ctx context.Context) *sql.DB {
	db, err := sql.Open("mysql", *mysqlURI)
	if err != nil {
		klog.Exitf("Failed to connect to DB: %v", err)
	}
	db.SetConnMaxLifetime(*dbConnMaxLifetime)
	db.SetMaxOpenConns(*dbMaxOpenConns)
	db.SetMaxIdleConns(*dbMaxIdleConns)

	initDatabaseSchema(ctx)
	return db
}

func createSignerOrDie() note.Signer {
	rawPrivateKey, err := os.ReadFile(*privateKeyPath)
	if err != nil {
		klog.Exitf("Failed to read private key file %q: %v", *privateKeyPath, err)
	}
	noteSigner, err := note.NewSigner(string(rawPrivateKey))
	if err != nil {
		klog.Exitf("Failed to create new signer: %v", err)
	}
	return noteSigner
}

func createVerifierOrDie() note.Verifier {
	rawPublicKey, err := os.ReadFile(*publicKeyPath)
	if err != nil {
		klog.Exitf("Failed to read public key file %q: %v", *publicKeyPath, err)
	}
	noteVerifier, err := note.NewVerifier(string(rawPublicKey))
	if err != nil {
		klog.Exitf("Failed to create new verifier: %v", err)
	}
	return noteVerifier
}

// configureTilesReadAPI adds the API methods from https://c2sp.org/tlog-tiles to the mux,
// routing the requests to the mysql storage.
// This method could be moved into the storage API as it's likely this will be
// the same for any implementation of a personality based on MySQL.
func configureTilesReadAPI(mux *http.ServeMux, storage *mysql.Storage) {
	mux.HandleFunc("GET /checkpoint", func(w http.ResponseWriter, r *http.Request) {
		checkpoint, err := storage.ReadCheckpoint(r.Context())
		if err != nil {
			klog.Errorf("/checkpoint: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if checkpoint == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if _, err := w.Write(checkpoint); err != nil {
			klog.Errorf("/checkpoint: %v", err)
			return
		}
	})

	mux.HandleFunc("GET /tile/{level}/{index...}", func(w http.ResponseWriter, r *http.Request) {
		level, index, width, err := layout.ParseTileLevelIndexWidth(r.PathValue("level"), r.PathValue("index"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if _, werr := w.Write([]byte(fmt.Sprintf("Malformed URL: %s", err.Error()))); werr != nil {
				klog.Errorf("/tile/{level}/{index...}: %v", werr)
			}
			return
		}

		tile, err := storage.ReadTile(r.Context(), level, index, width)
		if err != nil {
			klog.Errorf("/tile/{level}/{index...}: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if tile == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

		if _, err := w.Write(tile); err != nil {
			klog.Errorf("/tile/{level}/{index...}: %v", err)
			return
		}
	})

	mux.HandleFunc("GET /tile/entries/{index...}", func(w http.ResponseWriter, r *http.Request) {
		index, _, err := layout.ParseTileIndexWidth(r.PathValue("index"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if _, werr := w.Write([]byte(fmt.Sprintf("Malformed URL: %s", err.Error()))); werr != nil {
				klog.Errorf("/tile/entries/{index...}: %v", werr)
			}
			return
		}

		entryBundle, err := storage.ReadEntryBundle(r.Context(), index)
		if err != nil {
			klog.Errorf("/tile/entries/{index...}: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if entryBundle == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// TODO: Add immutable Cache-Control header.

		if _, err := w.Write(entryBundle); err != nil {
			klog.Errorf("/tile/entries/{index...}: %v", err)
			return
		}
	})
}

func initDatabaseSchema(ctx context.Context) {
	if *initSchemaPath != "" {
		klog.Infof("Initializing database schema")

		db, err := sql.Open("mysql", *mysqlURI+"?multiStatements=true")
		if err != nil {
			klog.Exitf("Failed to connect to DB: %v", err)
		}
		defer func() {
			if err := db.Close(); err != nil {
				klog.Warningf("Failed to close db: %v", err)
			}
		}()

		rawSchema, err := os.ReadFile(*initSchemaPath)
		if err != nil {
			klog.Exitf("Failed to read init schema file %q: %v", *initSchemaPath, err)
		}
		if _, err := db.ExecContext(ctx, string(rawSchema)); err != nil {
			klog.Exitf("Failed to execute init database schema: %v", err)
		}

		klog.Infof("Database schema initialized")
	}
}