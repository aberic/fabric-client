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

// Package shunt 负载服务包
package shunt

import (
	"fmt"
	"github.com/ennoo/rivet/server"
	"github.com/ennoo/rivet/trans/request"
	"github.com/ennoo/rivet/trans/response"
	"github.com/ennoo/rivet/utils/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"net/http"
	"sync"
)

const (
	// Random 负载均衡 random 策略
	Random = iota
	// Round 负载均衡 round 策略
	Round
	// Hash 负载均衡 hash 策略
	Hash
)

var (
	instance *Shunt
	once     sync.Once
	lbMap    = make(map[string]*LB)
)

// GetShuntInstance 获取负载管理对象 Shunt 单例
func GetShuntInstance() *Shunt {
	once.Do(func() {
		instance = &Shunt{AllWay: make(map[string]int)}
	})
	return instance
}

// Shunt 负载入口对象
type Shunt struct {
	AllWay map[string]int
}

// LBs 负载对象数组
type LBs struct {
	LBs []*LB `yaml:"shunt"`
}

// LB 负载均衡配置对象
type LB struct {
	Name      string `yaml:"Name"`      // 服务名称
	InURI     string `yaml:"InURI"`     // 负载入口 URI
	OutRemote string `yaml:"OutRemote"` // 负载出口地址
	Register  int    `yaml:"Register"`  // 负载均衡算法，1：随机；2：轮询；3：Hash一致性
}

// YamlLBs YML转负载对象数组
func YamlLBs(data []byte) {
	lbs := LBs{}
	err := yaml.Unmarshal([]byte(data), &lbs)
	if err != nil {
		log.Bow.Panic("cannot unmarshal data: " + err.Error())
	}
	if len(lbs.LBs) > 0 {
		for index := range lbs.LBs {
			lb := lbs.LBs[index]
			GetShuntInstance().Register(lb.Name, lb.Register)
			lbMap[lb.Name] = lb
		}
	}
}

// AddService 新增负载服务
func (s *Shunt) AddService(serviceName, inURI, outRemote string, way int) {
	lbMap[serviceName] = &LB{
		Name:      serviceName,
		InURI:     inURI,
		OutRemote: outRemote,
		Register:  way,
	}
	GetShuntInstance().Register(serviceName, way)
}

// Register 注册新的负载方式
func (s *Shunt) Register(serviceName string, way int) {
	switch way {
	case Round:
		if nil == roundRobinBalances {
			roundRobinBalances = make(map[string]*RoundRobinBalance)
		}
		roundRobinBalances[serviceName] = &RoundRobinBalance{
			serviceName: serviceName,
			rrbCh:       generaCount(),
		}
	}
	instance.AllWay[serviceName] = way
}

// GetService 根据服务名称获取待转发服务
func GetService(serviceName string) (service *server.Service, err error) {
	way, ok := instance.AllWay[serviceName]
	if !ok {
		err := fmt.Errorf("service not fount")
		fmt.Println("not found ", serviceName)
		log.Shunt.Error(err.Error(), zap.String("serviceName", serviceName))
		return nil, err
	}
	switch way {
	case Random:
		service, err = RunRandom(serviceName)
	case Round:
		service, err = RunRound(serviceName)
	case Hash:
		service, err = RunHash(serviceName)
	}
	if err != nil {
		err = fmt.Errorf(" %s erros", serviceName)
		return
	}
	return
}

// RunShunt 开启路由
func RunShunt(context *gin.Context, serviceName string) {
	RunShuntCallback(context, serviceName, nil)
}

// RunShuntCallback 开启路由并处理降级
func RunShuntCallback(context *gin.Context, serviceName string, f func() *response.Result) {
	lb := lbMap[serviceName]
	service, err := GetService(lb.OutRemote)
	outURI := context.Request.RequestURI[len(lb.InURI)+2:]
	var OutRemote string
	if nil == err {
		OutRemote = request.FormatURL(context, service)
	} else {
		result := response.Result{}
		result.Fail(err.Error())
		context.JSON(http.StatusOK, result)
		return
	}
	if nil == f {
		request.SyncPoolGetRequest().Call(context, context.Request.Method, OutRemote, outURI)
	} else {
		request.SyncPoolGetRequest().Callback(context, context.Request.Method, OutRemote, outURI, f)
	}
}
