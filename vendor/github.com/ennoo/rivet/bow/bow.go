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

// Package bow 网关服务包
package bow

import (
	"fmt"
	"github.com/ennoo/rivet/shunt"
	"github.com/ennoo/rivet/trans/request"
	"github.com/ennoo/rivet/trans/response"
	"github.com/ennoo/rivet/utils/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"net/http"
	"sync"
)

var (
	instance      *Bow
	once          sync.Once
	routeServices = make(map[string]*RouteService)
)

// GetBowInstance 获取路由管理对象 Bow 单例
func GetBowInstance() *Bow {
	once.Do(func() {
		instance = &Bow{AllWay: make(map[string]*RouteService)}
	})
	return instance
}

// Bow 路由入口对象
type Bow struct {
	AllWay map[string]*RouteService
}

// RouteServices 路由对象数组
type RouteServices struct {
	Routes []*RouteService `yaml:"routes"`
}

// RouteService 路由对象
type RouteService struct {
	Name      string `yaml:"Name"`      // 服务名称
	InURI     string `yaml:"InURI"`     // 路由入口 URI
	OutRemote string `yaml:"OutRemote"` // 路由出口地址
	Limit     *Limit `yaml:"Limit"`     // 服务限流策略
}

// YamlServices YML转路由对象数组
func YamlServices(data []byte) {
	routeServices := RouteServices{}
	err := yaml.Unmarshal([]byte(data), &routeServices)
	if err != nil {
		log.Bow.Panic("cannot unmarshal data: " + err.Error())
	}
	GetBowInstance().AddServices(routeServices.Routes)
}

// Add 新增路由服务数组
func (s *Bow) Add(routeServiceArr ...*RouteService) {
	for index := range routeServiceArr {
		routeService := routeServiceArr[index]
		routeServices[routeService.Name] = routeService
		GetBowInstance().register(routeService)
		if nil != routeService.Limit {
			routeService.Limit.LimitChan = make(chan int, routeService.Limit.LimitCount)
			go routeService.Limit.limit()
		}
	}
}

// AddServices 新增路由服务数组
func (s *Bow) AddServices(routeServiceArr []*RouteService) {
	for index := range routeServiceArr {
		routeService := routeServiceArr[index]
		routeServices[routeService.Name] = routeService
		GetBowInstance().register(routeService)
		if nil != routeService.Limit {
			routeService.Limit.LimitChan = make(chan int, routeService.Limit.LimitCount)
			go routeService.Limit.limit()
		}
	}
}

// AddService 新增路由服务
func (s *Bow) AddService(serviceName, inURI, outRemote string) {
	routeServices[serviceName] = &RouteService{
		Name:      serviceName,
		InURI:     inURI,
		OutRemote: outRemote,
	}
	GetBowInstance().register(&RouteService{
		Name:      serviceName,
		InURI:     inURI,
		OutRemote: outRemote,
	})
}

// Register 注册新的路由方式
func (s *Bow) register(routeService *RouteService) {
	instance.AllWay[routeService.Name] = routeService
}

// RunBow 开启路由
func RunBow(context *gin.Context, serviceName string, filter func(result *response.Result) bool) {
	RunBowCallback(context, serviceName, filter, nil)
}

// RunBowCallback 开启路由并处理降级
func RunBowCallback(context *gin.Context, serviceName string, filter func(result *response.Result) bool, f func() *response.Result) {
	routeService, ok := instance.AllWay[serviceName]
	result := response.Result{}
	if !ok {
		err := fmt.Errorf("routeService not fount")
		log.Shunt.Error(err.Error(), zap.String("serviceName", serviceName))
		result.Fail(err.Error())
		context.JSON(http.StatusOK, result)
		return
	}
	if !filter(&result) {
		context.JSON(http.StatusOK, result)
		return
	}
	// 限流
	if nil != routeService.Limit {
		routeService.Limit.LimitChan <- 1
	}
	outURI := context.Request.RequestURI[len(routeService.InURI)+2:]
	var OutRemote string
	if request.LB {
		service, err := shunt.GetService(routeService.OutRemote)
		if nil == err {
			OutRemote = request.FormatURL(context, service)
		} else {
			result.Fail(err.Error())
			context.JSON(http.StatusOK, result)
			return
		}
	}
	if nil == f {
		request.SyncPoolGetRequest().Call(context, context.Request.Method, OutRemote, outURI)
	} else {
		request.SyncPoolGetRequest().Callback(context, context.Request.Method, OutRemote, outURI, f)
	}
}
