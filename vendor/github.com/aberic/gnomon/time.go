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

import "time"

// TimeCommon 时间工具
type TimeCommon struct{}

// String2Timestamp 字符串转时间戳
//
// date 待转换时间字符串 如：2019/09/17 10:16:56
//
// format 时间字符串格式化类型 如：2006/01/02 15:04:05
//
// zone 时区 如：time.Local / time.UTC
func (t *TimeCommon) String2Timestamp(date, format string, zone *time.Location) (int64, error) {
	var (
		theTime time.Time
		err     error
	)
	if theTime, err = time.ParseInLocation(format, date, zone); nil != err {
		return 0, err
	}
	return theTime.Unix(), nil
}

// Timestamp2String 时间戳转字符串
//
// timestampSec 时间戳秒值
//
// timestampNSec 时间戳纳秒值
//
// format 时间字符串格式化类型 如：2006/01/02 15:04:05
func (t *TimeCommon) Timestamp2String(timestampSec, timestampNSec int64, format string, zone *time.Location) string {
	switch zone {
	default:
		return time.Unix(timestampSec, timestampNSec).Local().Format(format) //设置时间戳 使用模板格式化为日期字符串
	case time.UTC:
		return time.Unix(timestampSec, timestampNSec).UTC().Format(format) //设置时间戳 使用模板格式化为日期字符串
	}
}
