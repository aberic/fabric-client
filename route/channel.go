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

package route

import (
	"github.com/ennoo/fabric-go-client/core"
	"github.com/ennoo/fabric-go-client/service"
	"github.com/ennoo/rivet"
	"github.com/ennoo/rivet/trans/response"
)

func Channel(router *response.Router) {
	// 仓库相关路由设置
	router.Group = router.Engine.Group("/channel")
	router.POST("/create", create)
	router.POST("/join", join)
	router.POST("/list", list)
}

func create(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var channelCreate = new(service.ChannelCreate)
		if err := router.Context.ShouldBindJSON(channelCreate); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if nil == service.Get(channelCreate.ConfigID) {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		sdk.Create(channelCreate.OrderOrgName, channelCreate.OrgName, channelCreate.OrgUser, channelCreate.ChannelID,
			channelCreate.ChannelConfigPath, service.GetBytes(channelCreate.ConfigID)).Say(router.Context)
	})
}

func join(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var channelJoin = new(service.ChannelJoin)
		if err := router.Context.ShouldBindJSON(channelJoin); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if nil == service.Get(channelJoin.ConfigID) {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		sdk.Join(channelJoin.OrderName, channelJoin.OrgName, channelJoin.OrgUser, channelJoin.ChannelID, channelJoin.PeerUrl,
			service.GetBytes(channelJoin.ConfigID)).Say(router.Context)
	})
}

func list(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var channelList = new(service.ChannelList)
		if err := router.Context.ShouldBindJSON(installed); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if nil == service.Get(channelList.ConfigID) {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		sdk.Channels(channelList.OrgName, channelList.OrgUser, channelList.PeerName,
			service.GetBytes(channelList.ConfigID)).Say(router.Context)
	})
}
