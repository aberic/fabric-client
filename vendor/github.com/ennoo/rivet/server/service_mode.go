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

// Package server 路由、负载、监听外部服务包
package server

import "sync"

// ServiceReq 服务器新增请求对象
type ServiceReq struct {
	Name    string  `json:"name"`
	Service Service `json:"service"`
}

// Services 服务器对象集合
type Services struct {
	Services []*Service `json:"services"`
}

// Service 服务器信息
type Service struct {
	ID     string `json:"id"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Health string `json:"health"`
}

// Add Service 服务器对象集合内新增
func (services *Services) Add(service Service) {
	add := true
	for index := range services.Services {
		if services.Services[index].ID == service.ID {
			services.Services[index] = &service
			add = false
		}
	}
	if add {
		services.Services = append(services.Services, &service)
	}
}

// Remove Service 服务器对象集合内移除
func (services *Services) Remove(position int) {
	services.Services = append(services.Services[:position], services.Services[position+1:]...)
}

// NewService 服务器对象新建
func NewService(host string, port int) *Service {
	return &Service{
		Host: host,
		Port: port,
	}
}

// GetHost 获取服务器信息地址
func (a *Service) GetHost() string {
	return a.Host
}

// GetPort 获取服务器信息端口号
func (a *Service) GetPort() int {
	return a.Port
}

// Equal 比较服务器是否相同
func (a *Service) Equal(host string, port int) bool {
	return a.Host == host && a.Port == port
}

// EqualService 比较服务器是否相同
func (a *Service) EqualService(service *Service) bool {
	return a.Host == service.Host && a.Port == service.Port
}

var (
	serviceGroup map[string]*Services
	once         sync.Once
)

// ServiceGroup 全局服务器群组
func ServiceGroup() map[string]*Services {
	once.Do(func() {
		serviceGroup = make(map[string]*Services)
	})
	return serviceGroup
}

// GetServices 根据服务名称获取本地服务器列表，如果没有，则新建
func GetServices(serviceName string) *Services {
	services := ServiceGroup()[serviceName]
	if nil == services {
		services = &Services{}
		ServiceGroup()[serviceName] = services
	}
	return ServiceGroup()[serviceName]
}
