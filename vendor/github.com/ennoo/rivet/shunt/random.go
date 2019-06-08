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
	"math/rand"
)

// RunRandom 负载均衡 random 策略实现
func RunRandom(serviceName string) (add *server.Service, err error) {
	services := server.ServiceGroup()[serviceName].Services
	if len(services) == 0 {
		err = errors.New("no instance")
		return
	}

	lens := len(services)

	index := rand.Intn(lens)
	add = services[index]
	return
}
