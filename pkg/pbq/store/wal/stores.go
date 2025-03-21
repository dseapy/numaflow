/*
Copyright 2022 The Numaproj Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package wal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	dfv1 "github.com/numaproj/numaflow/pkg/apis/numaflow/v1alpha1"
	"github.com/numaproj/numaflow/pkg/pbq/partition"
	"github.com/numaproj/numaflow/pkg/pbq/store"
)

type walStores struct {
	storePath string
	// maxBufferSize max size of batch before it's flushed to store
	maxBatchSize int64
	// syncDuration timeout to sync to store
	syncDuration time.Duration
}

func NewWALStores(opts ...Option) store.StoreProvider {
	s := &walStores{
		storePath:    dfv1.DefaultStorePath,
		maxBatchSize: dfv1.DefaultStoreMaxBufferSize,
		syncDuration: dfv1.DefaultStoreSyncDuration,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

func (ws *walStores) CreateStore(ctx context.Context, partitionID partition.ID) (store.Store, error) {
	// Create wal dir if not exist
	var err error
	if _, err = os.Stat(ws.storePath); os.IsNotExist(err) {
		err = os.Mkdir(ws.storePath, 0644)
		if err != nil {
			return nil, err
		}
	}

	// let's open or create, initialize and return a new WAL
	wal, err := ws.openOrCreateWAL(&partitionID)

	return wal, err
}

// openOrCreateWAL open or creates a new WAL segment.
func (ws *walStores) openOrCreateWAL(id *partition.ID) (*WAL, error) {
	var err error

	filePath := getSegmentFilePath(id, ws.storePath)
	stat, err := os.Stat(filePath)

	var fp *os.File
	var wal *WAL
	if os.IsNotExist(err) {
		// here we are explicitly giving O_WRONLY because we will not be using this to read. Our read is only during
		// boot up.
		fp, err = os.OpenFile(getSegmentFilePath(id, ws.storePath), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		filesCount.With(map[string]string{}).Inc()
		activeFilesCount.With(map[string]string{}).Inc()
		wal = &WAL{
			fp:                fp,
			openMode:          os.O_WRONLY,
			createTime:        time.Now(),
			wOffset:           0,
			rOffset:           0,
			readUpTo:          0,
			partitionID:       id,
			prevSyncedWOffset: 0,
			prevSyncedTime:    time.Time{},
			walStores:         ws,
			numOfUnsyncedMsgs: 0,
		}

		err = wal.writeHeader()
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// here we are explicitly giving O_RDWR because we will be using this to read too. Our read is only during
		// bootstrap.
		fp, err = os.OpenFile(filePath, os.O_RDWR, stat.Mode())
		if err != nil {
			return nil, err
		}
		wal = &WAL{
			fp:                fp,
			openMode:          os.O_RDWR,
			createTime:        time.Now(),
			wOffset:           0,
			rOffset:           0,
			readUpTo:          stat.Size(),
			partitionID:       id,
			prevSyncedWOffset: 0,
			prevSyncedTime:    time.Time{},
			walStores:         ws,
			numOfUnsyncedMsgs: 0,
		}
		readPartition, err := wal.readHeader()
		if err != nil {
			return nil, err
		}
		if id.Key != readPartition.Key {
			return nil, fmt.Errorf("expected partition key %s, but got %s", id.Key, readPartition.Key)
		}
	}

	return wal, err
}

func (ws *walStores) DiscoverPartitions(ctx context.Context) ([]partition.ID, error) {
	files, err := os.ReadDir(ws.storePath)
	if os.IsNotExist(err) {
		return []partition.ID{}, nil
	} else if err != nil {
		return nil, err
	}
	partitions := make([]partition.ID, 0)

	for _, f := range files {
		if strings.HasPrefix(f.Name(), SegmentPrefix) && !f.IsDir() {
			filePath := filepath.Join(ws.storePath, f.Name())
			wal, err := OpenWAL(filePath)
			if err != nil {
				return nil, err
			}
			partitions = append(partitions, *wal.partitionID)
		}
	}

	return partitions, nil
}
