/*
 * Copyright (c) 2019. aberic - All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package gnomon

import (
	"errors"
	"sync"
	"testing"
	"time"
)

var logDir = "./tmp/log"

func logDo() {
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Debug("test", nil)
	Log().Info("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Warn("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Error("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true), Log().Err(errors.New("yes")))
}

func TestLogCommon_Fail(t *testing.T) {
	t.Log(Log().Init("/etc/log", 1, 1, false))
}

func TestLogCommon_Debug(t *testing.T) {
	_ = Log().Init("", 1, 1, false)
	Log().Set(Log().DebugLevel(), false)
	logDo()
}

func TestLogCommon_Info(t *testing.T) {
	Log().once = sync.Once{}
	_ = Log().Init(logDir, 1, 1, true)
	Log().Set(Log().InfoLevel(), false)
	logDo()
}

func TestLogCommon_Warn(t *testing.T) {
	_ = Log().Init(logDir, 1, 1, true)
	Log().Set(Log().WarnLevel(), false)
	logDo()
}

func TestLogCommon_Error(t *testing.T) {
	_ = Log().Init(logDir, 1, 1, false)
	Log().Set(Log().ErrorLevel(), false)
	logDo()
}

func TestLogCommon_Panic(t *testing.T) {
	_ = Log().Init(logDir, 1, 1, false)
	Log().Set(Log().PanicLevel(), false)
	logDo()
}

func TestLogCommon_Fatal(t *testing.T) {
	_ = Log().Init(logDir, 1, 1, false)
	Log().Set(Log().FatalLevel(), false)
	logDo()
}

func TestLogCommon_Fatal_BigStorage(t *testing.T) {
	_ = Log().Init(logDir, 1, 1, false)
	Log().Set(debugLevel, true)
	for i := 0; i < 100000; i++ {
		go Log().Debug("test", Log().Field("i", i), Log().Field("str", "str"), Log().Field("3", true))
	}
	time.Sleep(2 * time.Second)
}
