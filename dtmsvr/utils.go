/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package dtmsvr

import (
	"fmt"
	"time"

	"github.com/lithammer/shortuuid/v3"
	"github.com/sllt/dtm/client/dtmcli/dtmimp"
	"github.com/sllt/dtm/dtmsvr/config"
	"github.com/sllt/dtm/dtmsvr/storage"
	"github.com/sllt/dtm/dtmsvr/storage/registry"
)

type branchStatus struct {
	gid        string
	branchID   string
	op         string
	status     string
	finishTime *time.Time
}

var e2p = dtmimp.E2P

var conf = &config.Config

// GetStore returns storage.Store
func GetStore() storage.Store {
	return registry.GetStore()
}

// TransProcessedTestChan only for test usage. when transaction processed once, write gid to this chan
var TransProcessedTestChan chan string

// GenGid generate gid, use uuid
func GenGid() string {
	return shortuuid.New()
}

// GetTransGlobal construct trans from db
func GetTransGlobal(gid string) *TransGlobal {
	trans := GetStore().FindTransGlobalStore(gid)
	//nolint:staticcheck
	dtmimp.PanicIf(trans == nil, fmt.Errorf("no TransGlobal with gid: %s found", gid))
	//nolint:staticcheck
	return &TransGlobal{TransGlobalStore: *trans}
}
