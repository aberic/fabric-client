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

// Package scheduled 定时服务包
package scheduled

import (
	"github.com/ennoo/rivet/discovery"
	"github.com/ennoo/rivet/server"
	"github.com/ennoo/rivet/utils/log"
	"github.com/ennoo/rivet/utils/slip"
	"github.com/ennoo/rivet/utils/string"
	"go.uber.org/zap"
	"time"
)

var (
	selfServiceID          string
	selfServiceName        string
	selfDiscoveryComponent string
)

// CheckService 检查可用负载服务列表
func CheckService(serviceID, serviceName, component string) {
	log.Scheduled.Info("CheckService 检查可用负载服务列表")
	selfServiceID = serviceID
	selfServiceName = serviceName
	selfDiscoveryComponent = component
	if str.IsEmpty(selfDiscoveryComponent) {
		go checkSelfService()
	} else {
		go checkDiscoveryService()
	}
}

func checkDiscoveryService() {
	timeout := time.After(30 * time.Second) // timeout 是一个计时信道, 如果达到时间了，就会发一个信号出来
	// 发现服务通道
	abortDiscovery := make(chan int, 1)
	switch selfDiscoveryComponent {
	case discovery.ComponentConsul:
		go startCheckServicesByConsul(abortDiscovery)
		for isTimeout := false; !isTimeout; {
			select {
			case discoveryChanStatus := <-abortDiscovery:
				timeout = time.After(30 * time.Second) // 超时重置
				switch discoveryChanStatus {
				// 启动定时任务出错
				// 请求调用有误
				// 请求对方网络有误
				// 请求返回数据转JSON失败
				case slip.StartError, slip.RestRequestError, slip.RestResponseError, slip.JSONUnmarshalError:
					log.Scheduled.Debug("自身有误，则自处理完成后关闭当前协程，进入自检查模式")
					// 自身有误，则自处理完成后关闭当前协程
					isTimeout = true
					// 进入自检查模式
					checkSelfService()
				}
			case <-timeout:
				log.Scheduled.Debug("发现服务检查超时，则自处理完成后关闭当前协程，进入自检查模式")
				checkSelfService()
				isTimeout = true // 超时
			}
		}
	default:
		log.Scheduled.Warn("没有发现组件，进入自检查模式", zap.String("component", selfDiscoveryComponent))
		checkSelfService()
	}
}

func checkSelfService() {
	timeout := time.After(30 * time.Second) // timeout 是一个计时信道, 如果达到时间了，就会发一个信号出来
	// 自检查服务通道
	abortServices := make(chan int, 1)
	go startCheckServices(abortServices)
	for isTimeout := false; !isTimeout; {
		select {
		case servicesChanStatus := <-abortServices:
			timeout = time.After(30 * time.Second) // 超时重置
			switch servicesChanStatus {
			// 启动定时任务出错
			case slip.StartError:
				log.Scheduled.Debug("启动定时任务出错，关闭当前协程，开启新的自检查模式")
				checkSelfService()
				isTimeout = true
			// 启动发现服务
			case slip.DiscoveryStart:
				log.Scheduled.Debug("发现服务可用，退出自检查模式，重新启用注册发现服务")
				isTimeout = true
			}
		case <-timeout:
			log.Scheduled.Debug("自检查服务超时，关闭当前协程，开启新的自检查模式")
			checkSelfService()
			isTimeout = true // 超时
		}
	}
}

// compareAndResetServices 通过比较对象移除原本对象中多余项
func compareAndResetServices(services, servicesCompare *server.Services) {
	// todo 发现服务中没有的 service 应该交由自检查服务进行管理
	servicesArr := services.Services
	size := len(servicesArr)
	for i := 0; i < size; i++ {
		have := false
		for position := range servicesCompare.Services {
			if servicesCompare.Services[position].EqualService(services.Services[i]) {
				have = true
			}
		}
		if !have {
			services.Remove(i)
			i--
			size--
		}
	}
}
