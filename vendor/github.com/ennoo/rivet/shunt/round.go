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
 *
 */

package shunt

import (
	"errors"
	"github.com/ennoo/rivet/server"
)

var roundRobinBalances map[string]*RoundRobinBalance

// RoundRobinBalance 负载均衡 round 策略实体
type RoundRobinBalance struct {
	serviceName string
	rrbCh       chan int
}

// generaCount 自增生成一个0~65535的数，到达65535则重0开始计数
func generaCount() chan int {
	var ch = make(chan int)
	go func() {
		for i := 0; ; i++ {
			ch <- i // 等待索要数据
			if i == 65535 {
				i = 0
			}
		}
	}()
	return ch
}

// RunRound 负载均衡 round 策略实现
func RunRound(serviceName string) (service *server.Service, err error) {
	services := server.ServiceGroup()[serviceName].Services
	var lens int
	if lens = len(services); nil == services || lens == 0 {
		err = errors.New("no instance")
		return
	}
	roundRobinBalance := roundRobinBalances[serviceName]
	var position int
	if position = <-roundRobinBalance.rrbCh; position >= lens {
		position = position % lens
	}
	service = services[position]
	return
}
