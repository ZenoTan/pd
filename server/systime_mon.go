// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"time"

	"github.com/pingcap/log"
	errs "github.com/pingcap/pd/v4/pkg/errors"
	"go.uber.org/zap"
)

// StartMonitor calls systimeErrHandler if system time jump backward.
func StartMonitor(ctx context.Context, now func() time.Time, systimeErrHandler func()) {
	log.Info("start system time monitor")
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	for {
		last := now().UnixNano()
		select {
		case <-tick.C:
			if now().UnixNano() < last {
				log.Error("system time jump backward", zap.Int64("last", last), zap.Error(errs.ErrOtherSystemTime.FastGenByArgs()))
				systimeErrHandler()
			}
		case <-ctx.Done():
			return
		}
	}
}
