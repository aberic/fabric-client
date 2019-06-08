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

package scheduled

import (
	"github.com/ennoo/rivet/discovery"
	"github.com/ennoo/rivet/discovery/consul"
	"github.com/ennoo/rivet/server"
	"github.com/ennoo/rivet/shunt"
	"github.com/ennoo/rivet/trans/request"
	"github.com/ennoo/rivet/utils/log"
	"github.com/ennoo/rivet/utils/slip"
	"github.com/ennoo/rivet/utils/string"
	"github.com/robfig/cron"
	"go.uber.org/zap"
	"net/http"
)

// startCheckServices 定时检查已存在的服务列表
func startCheckServices(abortServices chan int) {
	c := cron.New()
	// 每隔5秒执行一次：*/5 * * * * ?
	// 每隔1分钟执行一次：0 */1 * * * ?
	// 每天23点执行一次：0 0 23 * * ?
	// 每天凌晨1点执行一次：0 0 1 * * ?
	// 每月1号凌晨1点执行一次：0 0 1 1 * ?
	// 在26分、29分、33分执行一次：0 26,29,33 * * * ?
	// 每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
	err := c.AddFunc("*/10 * * * * ?", func() {
		checkServices(abortServices)
	})
	if nil != err {
		abortServices <- slip.StartError
	} else {
		c.Start()
	}
}

// checkServicesByConsul
//
// 获取本地可负载服务名称列表
//
// 获取本地本地列表
//
// 遍历本地列表
//
// 如不可用，则移除
//
// 如可用，则继续下一轮循环
func checkServices(abortServices chan int) {
	switch selfDiscoveryComponent {
	case discovery.ComponentConsul:
		// 检查发现服务状态
		_, slips := consul.ServiceCheck(selfServiceName)
		if nil == slips {
			abortServices <- slip.DiscoveryStart
			return
		}
	}
	// 获取本地可负载服务列表
	allWay := shunt.GetShuntInstance().AllWay
	// 根据本地可负载服务列表遍历发现服务(线上)中是否存在
	for serviceName := range allWay {
		// 获取本地服务列表
		services := server.GetServices(serviceName)
		log.Scheduled.Debug("获取本地服务列表", zap.Any("servicesArr", services.Services))
		servicesArr := services.Services
		size := len(servicesArr)
		for i := 0; i < size; i++ {
			log.Scheduled.Debug("自检查 servicesArr", zap.Int("i", i))
			if str.IsNotEmpty(servicesArr[i].Health) {
				method := http.MethodGet
				body, err := request.SyncPoolGetRequest().RestJSONByURL(method, servicesArr[i].Health, nil)
				log.Scheduled.Debug("自检查 body 结果", zap.String("body", string(body)))
				if nil != err {
					log.Scheduled.Debug("自检查 err 结果", zap.String("err", err.Error()))
					services.Remove(i)
					i--
					size--
				}
			} else {
				services.Remove(i)
				i--
				size--
			}
		}
		abortServices <- slip.StartSuccess
	}
}
