/*
 * Copyright (c) 2019.. ENNOO - All Rights Reserved.
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

package route

import (
	"github.com/ennoo/fabric-go-client/core"
	"github.com/ennoo/fabric-go-client/service"
	"github.com/ennoo/rivet"
	"github.com/ennoo/rivet/trans/response"
)

func Order(router *response.Router) {
	// 仓库相关路由设置
	router.Group = router.Engine.Group("/order")
	router.POST("/config", config)
}

func config(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var orderConfig = new(service.OrderConfig)
		if err := router.Context.ShouldBindJSON(orderConfig); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if nil == service.Get(orderConfig.ConfigID) {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		sdk.OrderConfig(orderConfig.OrgName, orderConfig.OrgUser, orderConfig.ChannelID, orderConfig.OrderURL,
			service.GetBytes(orderConfig.ConfigID)).Say(router.Context)
	})
}
