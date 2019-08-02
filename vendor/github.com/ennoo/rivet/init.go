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

package rivet

import (
	"github.com/ennoo/rivet/bow"
	"github.com/ennoo/rivet/discovery"
	"github.com/ennoo/rivet/discovery/consul"
	"github.com/ennoo/rivet/scheduled"
	"github.com/ennoo/rivet/server"
	"github.com/ennoo/rivet/shunt"
	"github.com/ennoo/rivet/trans/request"
	"github.com/ennoo/rivet/trans/response"
	"github.com/ennoo/rivet/utils/env"
	"github.com/ennoo/rivet/utils/log"
	"github.com/ennoo/rivet/utils/string"
	"github.com/gin-gonic/gin"
	"runtime"
	"time"
)

var (
	hc = false // 是否开启健康检查。开启后为 Get 请求，路径为 /health/check
	sm = false // 是否开启外界服务管理功能
	ud = false // 是否启用发现服务
	cp string  // 启用的发现服务组件类型
	sn string  // 注册到发现服务的服务名称（优先通过环境变量 SERVICE_NAME 获取）
	f  func(result *response.Result) bool
)

// ListenServe 启动监听端口服务对象
type ListenServe struct {
	Engine *gin.Engine
	// defaultPort 默认启用端口号，优先通过环境变量 PORT 获取
	DefaultPort string
	// connectTimeout 拨号等待连接完成的最长时间，TCP超时的时间一般在3s，默认3s
	ConnectTimeout time.Duration
	// keepAlive 指定保持活动网络连接的时间，如果为0，则不启用keep-alive，默认30s
	KeepAlive time.Duration

	// TLS 服务端私钥
	CertFile string
	// TLS 服务端的数字证书
	KeyFile string
}

// Initialize rivet 初始化方法，必须最先调用
//
// bow：是否开启网关路由
//
// healthCheck：是否开启健康检查。开启后为 Get 请求，路径为 /health/check
//
// serverManager：是否开启外界服务管理功能
//
// loadBalance：是否开启负载均衡
func Initialize(healthCheck, serverManager, loadBalance, dev bool) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	hc = healthCheck
	sm = serverManager
	request.LB = loadBalance
	if dev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

// UseBow 开启网关路由
//
// filter 自定义过滤方案
func UseBow(filter func(result *response.Result) bool) {
	request.Route = true
	f = filter
}

// UseDiscovery 启用指定的发现服务
//
// component：启用的发现服务组件类型
//
// url：consul 等发现服务注册地址，包括端口号（优先通过环境变量 CONSUL_URL 获取）
//
// serviceName：注册到发现服务的服务名称（优先通过环境变量 SERVICE_NAME 获取）
//
// hostname：注册到发现服务的服务地址（如果为空，则尝试通过 /etc/hostname 获取）
//
// port：注册到 consul 的服务端口（优先通过环境变量 PORT 获取）
func UseDiscovery(component, url, serviceName, hostname string, port int) {
	cp = env.GetEnvDefault(env.DiscoveryComponent, component)
	sn = serviceName
	discoveryURL := env.GetEnvDefault(env.DiscoveryURL, url)
	discoveryReceiveHost := env.GetEnvDefault(env.DiscoveryReceiveHost, hostname)
	switch cp {
	case discovery.ComponentConsul:
		if !ud {
			log.Rivet.Info("use discovery service {}" + discovery.ComponentConsul)
			ud = true
			if request.LB {
				go consul.Enroll(discoveryURL, serviceID, ServiceName(), discoveryReceiveHost, port)
			} else {
				go scheduled.ConsulEnroll(discoveryURL, serviceID, ServiceName(), discoveryReceiveHost, port)
			}
		}
	}
}

// SetupRouter 设置路由器相关选项
func SetupRouter(routes ...func(router *response.Router)) *gin.Engine {
	router := &response.Router{
		Engine: gin.Default(),
	}

	if request.Route {
		bow.Route(router.Engine, f)
	}
	if hc {
		Health(router.Engine)
	}
	if sm {
		server.Server(router.Engine)
	}
	if request.LB {
		if ud {
			scheduled.CheckService(serviceID, sn, cp)
		} else {
			scheduled.CheckService(serviceID, sn, "")
		}
		if !request.Route {
			shunt.Route(router.Engine)
		}
	}
	for _, route := range routes {
		route(router)
	}
	return router.Engine
}

// ListenAndServe 开始启用 rivet
//
// listenServe 启动监听端口服务对象
//
// caCertPaths 作为客户端发起 HTTPS 请求时所需客户端证书路径数组
func ListenAndServe(listenServe *ListenServe, caCertPaths ...string) {
	listenAndServe(listenServe, false, caCertPaths...)
}

// ListenAndServes 开始启用 rivet
//
// listenServe 启动监听端口服务对象
//
// caCertPaths 作为客户端发起 HTTPS 请求时所需客户端证书路径数组
func ListenAndServes(listenServe *ListenServe, caCertPaths []string) {
	listenAndServe(listenServe, false, caCertPaths...)
}

// ListenAndServeTLS 开始启用 rivet
//
// listenServe 启动监听端口服务对象
//
// caCertPaths 作为客户端发起 HTTPS 请求时所需客户端证书路径数组
func ListenAndServeTLS(listenServe *ListenServe, caCertPaths ...string) {
	listenAndServe(listenServe, true, caCertPaths...)
}

// ListenAndServesTLS 开始启用 rivet
//
// listenServe 启动监听端口服务对象
//
// caCertPaths 作为客户端发起 HTTPS 请求时所需客户端证书路径数组
func ListenAndServesTLS(listenServe *ListenServe, caCertPaths []string) {
	listenAndServe(listenServe, true, caCertPaths...)
}

func listenAndServe(listenServe *ListenServe, isTLS bool, caCertPaths ...string) {
	if nil == listenServe.Engine {
		log.Rivet.Fatal("HTTP Engine is nil")
	}
	if str.IsEmpty(listenServe.DefaultPort) {
		log.Rivet.Fatal("HTTP Listening Port is nil")
	}
	if listenServe.ConnectTimeout < 0 {
		listenServe.ConnectTimeout = 3 * time.Second
	}
	if listenServe.KeepAlive < 0 {
		listenServe.KeepAlive = 30 * time.Second
	}
	request.GetTPInstance().Timeout(listenServe.ConnectTimeout, listenServe.KeepAlive).RootCACerts(caCertPaths).Instantiate()
	log.Rivet.Info("listening http port bind")
	var err error
	if isTLS {
		err = listenServe.Engine.RunTLS(":"+env.GetEnvDefault(env.Port, listenServe.DefaultPort), listenServe.CertFile, listenServe.KeyFile)
	} else {
		err = listenServe.Engine.Run(":" + env.GetEnvDefault(env.Port, listenServe.DefaultPort))
	}
	if nil != err {
		log.Rivet.Info("exit because {}" + err.Error())
	}
}
