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

package scheduled

import (
	"github.com/ennoo/rivet/discovery/consul"
	"github.com/ennoo/rivet/server"
	"github.com/ennoo/rivet/shunt"
	"github.com/ennoo/rivet/utils/slip"
	"github.com/robfig/cron"
	"strings"
)

// startCheckServicesByConsul 使用 consul 定时检出Service列表
func startCheckServicesByConsul(abortDiscovery chan int) {
	c := cron.New()
	// 每隔5秒执行一次：*/5 * * * * ?
	// 每隔1分钟执行一次：0 */1 * * * ?
	// 每天23点执行一次：0 0 23 * * ?
	// 每天凌晨1点执行一次：0 0 1 * * ?
	// 每月1号凌晨1点执行一次：0 0 1 1 * ?
	// 在26分、29分、33分执行一次：0 26,29,33 * * * ?
	// 每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
	err := c.AddFunc("*/10 * * * * ?", func() {
		checkServicesByConsul(abortDiscovery)
	})
	if nil != err {
		abortDiscovery <- slip.StartError
	} else {
		c.Start()
	}
}

// checkServicesByConsul
//
// 获取本地可负载服务名称列表
//
// 根据本地可负载服务列表遍历发现服务(线上)中是否存在
//
// 如不存在，则继续下一轮遍历
//
// 如存在且列表大于0，遍历线上服务列表并检查线上服务状态是否为可用
//
// 获取本地本地列表 x
//
// 新建空服务列表 y
//
// 如不可用，且 x 中包含此服务，则移除 x 中的服务
//
// 如可用，且 x 中不包含此服务，则新增服务到 x,y 中
//
// 移除 x 中不包含 y 的服务
func checkServicesByConsul(abortDiscovery chan int) {
	// 检查发现服务状态
	agentServiceChecks, slips := consul.ServiceCheck(selfServiceName)
	if nil == slips {
		have := false
		for index := range agentServiceChecks {
			if agentServiceChecks[index].Service.ID == selfServiceID {
				have = true
			}
		}
		if !have {
			consul.ReEnroll()
		}
	} else {
		abortDiscovery <- slips.Code
		return
	}
	// 获取本地可负载服务列表
	allWay := shunt.GetShuntInstance().AllWay
	// 根据本地可负载服务列表遍历发现服务(线上)中是否存在
	for serviceName := range allWay {
		agentServiceChecks, slips = consul.ServiceCheck(serviceName)
		if nil != slips {
			abortDiscovery <- slips.Code
		}
		// 如不存在，则继续下一轮遍历
		if nil == agentServiceChecks || len(agentServiceChecks) <= 0 {
			continue
		}
		// 获取本地本地列表
		services := server.GetServices(serviceName)
		// 新建空服务列表
		servicesCompare := server.Services{}
		// 如存在且列表大于0，遍历线上服务列表并检查线上服务状态是否为可用
		checkUpAndLocalByConsul(agentServiceChecks, services, &servicesCompare)
		// 移除 x 中不包含 y 的服务
		compareAndResetServices(services, &servicesCompare)
	}
	abortDiscovery <- slip.StartSuccess
}

// checkUpAndLocalByConsul 如存在且列表大于0，遍历线上服务列表并检查线上服务状态是否为可用
func checkUpAndLocalByConsul(agentServiceChecks []*consul.AgentServiceCheck, services, servicesCompare *server.Services) {
	for index := range agentServiceChecks {
		agentServiceCheck := agentServiceChecks[index]
		// 如不可用，且本地列表中包含此服务，则移除本地列表中的服务
		if agentServiceCheck.AggregatedStatus != "passing" {
			checkRemoveServiceByConsul(services, agentServiceCheck)
		} else { // 如可用，且本地列表中不包含此服务，则新增服务到本地列表中
			checkAddServiceByConsul(services, servicesCompare, agentServiceCheck)
		}
	}
}

// checkRemoveServiceByConsul 移除本地列表中的服务
func checkRemoveServiceByConsul(services *server.Services, agentServiceCheck *consul.AgentServiceCheck) {
	servicesArr := services.Services
	size := len(servicesArr)
	for i := 0; i < size; i++ {
		if servicesArr[i].Equal(agentServiceCheck.Service.Address, agentServiceCheck.Service.Port) {
			services.Remove(i)
			i--
			size--
		}
	}
}

// checkAddServiceByConsul 新增服务到本地列表中
func checkAddServiceByConsul(services, servicesCompare *server.Services, agentServiceCheck *consul.AgentServiceCheck) {
	var health string // 服务健康检查地址
	for offset := range agentServiceCheck.Checks {
		if health = strings.Split(agentServiceCheck.Checks[offset].Output, " ")[2]; !strings.HasPrefix(health, "http") {
			continue
		}
		health = health[0 : len(health)-1]
		break
	}
	service := server.Service{
		ID:     agentServiceCheck.Service.ID,
		Host:   agentServiceCheck.Service.Address,
		Port:   agentServiceCheck.Service.Port,
		Health: health,
	}
	have := false
	for position := range services.Services {
		if nil != services.Services && services.Services[position].Equal(service.Host, service.Port) {
			have = true
			break
		}
	}
	if !have {
		services.Add(service)
	}
	servicesCompare.Add(service)
}
