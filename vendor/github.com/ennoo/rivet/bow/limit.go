/*
 * Copyright (c) 2019. ENNOO - All Rights Reserved.
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
 */

package bow

import (
	"github.com/ennoo/rivet/utils/log"
	"strconv"
	"time"
)

// Limit 限流对象
type Limit struct {
	LimitMillisecond         int64    `yaml:"LimitMillisecond"`         // 请求限定的时间段（毫秒）
	LimitCount               int      `yaml:"LimitCount"`               // 请求限定的时间段内允许的请求次数
	LimitIntervalMillisecond int64    `yaml:"LimitIntervalMillisecond"` // 请求允许的最小间隔时间（毫秒），0表示不限
	LimitChan                chan int // 限流通道
	Times                    []int64  // 请求时间数组
}

func (l *Limit) timeInit() {
	l.Times = make([]int64, 0)
	for i := 0; i < l.LimitCount; i++ {
		time.Sleep(10 * time.Millisecond)
		l.add(time.Now().UnixNano() / 1e6)
	}
}

func (l *Limit) limit() {
	for {
		timeNow := time.Now().UnixNano() / 1e6
		if len(l.Times) < l.LimitCount {
			l.add(time.Now().UnixNano() / 1e6)
		} else if timeNow-l.Times[0] > l.LimitMillisecond && timeNow-l.Times[len(l.Times)-1] > l.LimitIntervalMillisecond {
			limitChanResult := strconv.Itoa(<-l.LimitChan)
			log.Bow.Debug("取出一个元素，放行 <-c = " + limitChanResult)
			if len(l.Times) == l.LimitCount {
				l.remove()
			}
			l.add(time.Now().UnixNano() / 1e6)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// add 新增一个元素
func (l *Limit) add(time int64) {
	if len(l.Times) < l.LimitCount {
		l.Times = append(l.Times, time)
	}
}

// remove 移除第一个元素
func (l *Limit) remove() {
	l.Times = l.Times[1:]
}
