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

// integration_test ensures that all storage implementations have a consistent contract.
package integration_test

import (
	"context"

	tessera "github.com/transparency-dev/trillian-tessera"
	"github.com/transparency-dev/trillian-tessera/storage/aws"
	"github.com/transparency-dev/trillian-tessera/storage/gcp"
	"github.com/transparency-dev/trillian-tessera/storage/mysql"
	"github.com/transparency-dev/trillian-tessera/storage/posix"
)

type StorageContract interface {
	Add(ctx context.Context, entry *tessera.Entry) tessera.IndexFuture
	ReadCheckpoint(ctx context.Context) ([]byte, error)
	ReadTile(ctx context.Context, level, index uint64, p uint8) ([]byte, error)
	ReadEntryBundle(ctx context.Context, index uint64, p uint8) ([]byte, error)
}

var (
	_ StorageContract = &mysql.Storage{}
	_ StorageContract = &gcp.Storage{}
	_ StorageContract = &posix.Storage{}
	_ StorageContract = &aws.Storage{}
)
