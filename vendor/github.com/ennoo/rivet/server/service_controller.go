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

package server

import (
	"github.com/ennoo/rivet/trans/response"
	"github.com/gin-gonic/gin"
	"strings"
)

// Server 服务器管理服务路由
func Server(engine *gin.Engine) {
	// 仓库相关路由设置
	vRepo := engine.Group("/service")
	vRepo.GET("/group/list", listGroup)
	vRepo.DELETE("/group/rm/:serviceName", rmGroup)
	vRepo.GET("/list/:serviceName", listService)
	vRepo.POST("/add", addService)
	vRepo.DELETE("/rm/:serviceName/:serviceId", rmService)
}

func addService(context *gin.Context) {
	resp := response.Response{}
	resp.Do(context, func(result *response.Result) {
		serviceReq := new(ServiceReq)
		if err := context.ShouldBindJSON(serviceReq); err != nil {
			result.SayFail(context, err.Error())
		}
		name := serviceReq.Name
		services := ServiceGroup()[name]
		if nil == services {
			services = &Services{}
			ServiceGroup()[name] = services
		}
		services.Add(serviceReq.Service)
		result.SaySuccess(context, "add service success")
	})
}

func rmService(context *gin.Context) {
	resp := response.Response{}
	resp.Do(context, func(result *response.Result) {
		serviceName := context.Param("serviceName")
		serviceID := context.Param("serviceId")
		if nil == ServiceGroup()[serviceName] {
			panic(response.ExpNotExist.Fit(strings.Join([]string{"service", serviceName}, " ")))
		}
		have := false
		service := ServiceGroup()[serviceName]
		services := service.Services
		for i := 0; i < len(services); i++ {
			if serviceID == services[i].ID {
				have = true
				service.Remove(i)
			}
		}
		if have {
			result.SaySuccess(context, strings.Join([]string{"remove service", serviceName, "id =", serviceID, "success"}, " "))
		} else {
			panic(response.ExpNotExist.Fit(strings.Join([]string{"service", serviceName, "id =", serviceID}, " ")))
		}
	})
}

func listService(context *gin.Context) {
	resp := response.Response{}
	resp.Do(context, func(result *response.Result) {
		serviceName := context.Param("serviceName")
		if nil != ServiceGroup()[serviceName] {
			result.SaySuccess(context, ServiceGroup()[serviceName].Services)
		} else {
			result.SaySuccess(context, []Service{})
		}
	})
}

func rmGroup(context *gin.Context) {
	resp := response.Response{}
	resp.Do(context, func(result *response.Result) {
		serviceName := context.Param("serviceName")
		if nil == ServiceGroup()[serviceName] {
			panic(response.ExpNotExist.Fit(strings.Join([]string{"service", serviceName}, " ")))
		}
		delete(ServiceGroup(), serviceName)
		result.SaySuccess(context, strings.Join([]string{"remove ", serviceName, " balance success"}, ""))
	})
}

func listGroup(context *gin.Context) {
	resp := response.Response{}
	resp.Do(context, func(result *response.Result) {
		shunts := make([]string, len(ServiceGroup()))
		index := 0
		for k := range ServiceGroup() {
			shunts[index] = k
			index++
		}
		result.SaySuccess(context, shunts)
	})
}
