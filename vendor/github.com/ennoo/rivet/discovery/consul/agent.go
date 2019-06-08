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

import (
	"encoding/json"
	"github.com/ennoo/rivet/trans/request"
	"github.com/ennoo/rivet/utils/env"
	"github.com/ennoo/rivet/utils/file"
	"github.com/ennoo/rivet/utils/log"
	"github.com/ennoo/rivet/utils/slip"
	"github.com/ennoo/rivet/utils/string"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// agentRegister 注册服务到 consul
//
// consulUrl：consul 注册地址，包括端口号（优先通过环境变量 DISCOVERY_URL 获取）
//
// serviceID：注册到 consul 的服务 ID
//
// serviceName：注册到 consul 的服务名称（优先通过环境变量 SERVICE_NAME 获取）
//
// hostname：注册到 consul 的服务地址（如果为空，则尝试通过 /etc/hostname 获取）
//
// port：注册到 consul 的服务端口（优先通过环境变量 PORT 获取）
func agentRegister(consulURL, serviceID, serviceName, hostname string, port int) {
	agentServiceChecks, slips := ServiceCheck(selfServiceName)
	if nil == slips {
		have := false
		for index := range agentServiceChecks {
			if agentServiceChecks[index].Service.ID == selfServiceID {
				have = true
			}
		}
		if have {
			return
		}
	} else {
		time.Sleep(1 * time.Second)
		ReEnroll()
		return
	}
	if containerID, err := file.ReadFileFirstLine("/etc/hostname"); nil == err && str.IsEmpty(hostname) {
		hostname = containerID
	} else {
		log.Discovery.Info("open /etc/hostname: no such file or directory")
	}
	method := http.MethodPut
	remote := strings.Join([]string{"http://", env.GetEnvDefault(env.DiscoveryURL, consulURL)}, "")
	uri := "v1/agent/service/register"
	if envPort := env.GetEnv(env.Port); str.IsNotEmpty(envPort) {
		envPortInt, err := strconv.Atoi(envPort)
		if err != nil {
			port = envPortInt
		}
	}
	healthCheckURL := strings.Join([]string{"http://", hostname, ":", strconv.Itoa(port), "/health/check"}, "")

	log.Discovery.Info("consul register info",
		zap.String("serviceID", serviceID),
		zap.String("hostname", hostname),
		zap.String("method", http.MethodPut),
		zap.String("remote", remote),
		zap.String("health", healthCheckURL))
	param := Register{
		ID:                serviceID,
		Name:              env.GetEnvDefault(env.ServiceName, serviceName),
		Address:           hostname,
		Port:              port,
		EnableTagOverride: false,
		Check: Check{
			DeregisterCriticalServiceAfter: "1m",
			HTTP:                           healthCheckURL,
			Interval:                       "10s"},
	}

	if body, err := request.SyncPoolGetRequest().RestJSON(method, remote, uri, param); nil != err {
		log.Discovery.Warn(err.Error(),
			zap.String("url", strings.Join([]string{remote, "/", uri}, "")),
			zap.String("method", method))
	} else {
		log.Discovery.Info(string(body),
			zap.String("url", strings.Join([]string{remote, "/", uri}, "")),
			zap.String("method", method),
		)
	}
}

// agentCheck 检查 consul 中各服务状态
//
// consulUrl：consul 注册地址，包括端口号（优先通过环境变量 CONSUL_URL 获取）
func agentCheck(consulURL string) {
	method := http.MethodGet
	remote := strings.Join([]string{"http://", env.GetEnvDefault(env.DiscoveryURL, consulURL)}, "")
	uri := "v1/agent/checks"

	if body, err := request.SyncPoolGetRequest().RestJSON(method, remote, uri, nil); nil != err {
		log.Discovery.Warn(err.Error(),
			zap.String("url", strings.Join([]string{remote, "/", uri}, "")),
			zap.String("method", method))
	} else {
		log.Discovery.Info(string(body),
			zap.String("url", strings.Join([]string{remote, "/", uri}, "")),
			zap.String("method", method),
		)
	}
}

// serviceCheck 检查 consul 中各服务状态
//
// consulUrl：consul 注册地址，包括端口号（优先通过环境变量 CONSUL_URL 获取）
//
// serviceName：想要检出的服务名称
func serviceCheck(consulURL, serviceName string) ([]*AgentServiceCheck, *slip.Slip) {
	method := http.MethodGet
	remote := strings.Join([]string{"http://", env.GetEnvDefault(env.DiscoveryURL, consulURL)}, "")
	uri := strings.Join([]string{"v1/agent/health/service/name/", serviceName}, "")
	slips := slip.Slip{}
	if body, err := request.SyncPoolGetRequest().RestJSON(method, remote, uri, nil); nil != err {
		log.Discovery.Info(err.Error(),
			zap.String("url", strings.Join([]string{remote, "/", uri}, "")),
			zap.String("method", method))
		slips.FormatError(slip.RestResponseError, err)
	} else {
		var agentServiceChecks []*AgentServiceCheck
		if err = json.Unmarshal(body, &agentServiceChecks); nil == err {
			return agentServiceChecks, nil
		}
		slips.FormatError(slip.JSONUnmarshalError, err)
		log.Discovery.Warn(err.Error(),
			zap.String("url", strings.Join([]string{remote, "/", uri}, "")),
			zap.String("method", method),
			zap.String("body", string(body)))
	}
	return nil, &slips
}
