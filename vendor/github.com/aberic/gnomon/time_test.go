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
	"testing"
	"time"
)

func TestTimeCommon_String2Timestamp(t *testing.T) {
	i64, err := Time().String2Timestamp("2019/09/17 10:16:56", "2006/01/02 15:04:05", time.Local)
	t.Log("i64", i64, err)
	i64, err = Time().String2Timestamp("2019/09/17 10:16:56", "2006/01/02 15:04:05", time.UTC)
	t.Log("i64", i64, err)
}

func TestTimeCommon_String2Timestamp_Fail(t *testing.T) {
	_, err := Time().String2Timestamp("hello world", "2006/01/02 15:04:05", time.Local)
	t.Skip(err)
}

func TestTimeCommon_Timestamp2String(t *testing.T) {
	t.Log(Time().Timestamp2String(1568686616, 28889, "2006/01/02 15:04:05", time.Local))
	t.Log(Time().Timestamp2String(1568686626, 98882, "2006/01/02 15:04:05", time.UTC))
}
