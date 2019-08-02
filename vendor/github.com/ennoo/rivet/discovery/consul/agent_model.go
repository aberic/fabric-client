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

package consul

// AgentServiceCheck consul 中根据服务名称获取到服务列表中单个对象
type AgentServiceCheck struct {
	AggregatedStatus string       `json:"AggregatedStatus"`
	Service          AgentService `json:"Service"`
	Checks           []AgentCheck `json:"Checks"`
}

// AgentService AgentServiceCheck 中所属服务对象
type AgentService struct {
	ID      string `json:"ID"`
	Service string `json:"Service"`
	Port    int    `json:"Port"`
	Address string `json:"Address"`
}

// AgentCheck AgentServiceCheck 中所属服务健康检查 URL
type AgentCheck struct {
	Output string `json:"Output"`
}
