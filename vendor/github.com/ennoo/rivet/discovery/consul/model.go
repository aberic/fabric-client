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

package consul

// Register consul 注册实体
type Register struct {
	ID                string `json:"ID"`
	Name              string `json:"Name"`
	Address           string `json:"Address"`
	Port              int    `json:"Port"`
	EnableTagOverride bool   `json:"EnableTagOverride"`
	Check             Check  `json:"Check"`
}

// Check consul 注册实体中健康检查实体
type Check struct {
	DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter"`
	HTTP                           string `json:"HTTP"`
	Interval                       string `json:"Interval"`
}
