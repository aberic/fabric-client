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
	"github.com/ennoo/rivet/discovery/consul"
	"github.com/robfig/cron"
)

// ConsulEnroll 定时调用此方法检查注册 consul 是否失效
//
// consulUrl：consul 注册地址，包括端口号（优先通过环境变量 CONSUL_URL 获取）
//
// serviceID：注册到 consul 的服务 ID
//
// serviceName：注册到 consul 的服务名称（优先通过环境变量 SERVICE_NAME 获取）
//
// hostname：注册到 consul 的服务地址（如果为空，则尝试通过 /etc/hostname 获取）
//
// port：注册到 consul 的服务端口（优先通过环境变量 PORT 获取）
func ConsulEnroll(consulURL, serviceID, serviceName, hostname string, port int) {
	enrollOnce := false
	c := cron.New()
	// 每隔5秒执行一次：*/5 * * * * ?
	// 每隔1分钟执行一次：0 */1 * * * ?
	// 每天23点执行一次：0 0 23 * * ?
	// 每天凌晨1点执行一次：0 0 1 * * ?
	// 每月1号凌晨1点执行一次：0 0 1 1 * ?
	// 在26分、29分、33分执行一次：0 26,29,33 * * * ?
	// 每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
	_ = c.AddFunc("*/10 * * * * ?", func() {
		if enrollOnce {
			consul.ReEnroll()
		} else {
			enrollOnce = true
			consul.Enroll(consulURL, serviceID, serviceName, hostname, port)
		}
	})
	c.Start()
}
