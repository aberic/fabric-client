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
	cfg "github.com/ennoo/fabric-client/config"
	"github.com/ennoo/fabric-client/core"
	"github.com/ennoo/fabric-client/geneses"
	genesis "github.com/ennoo/fabric-client/grpc/proto/geneses"
	"github.com/ennoo/fabric-client/service"
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
		var (
			conf        *cfg.Config
			orderOrgURL string
			orgName     string
		)
		if conf = service.Configs[channelCreate.ConfigID]; nil == conf {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		for _, order := range conf.Orderers {
			orderOrgURL = order.URL
		}
		for name, org := range conf.Organizations {
			if len(org.Peers) <= 0 {
				continue
			}
			orgName = name
		}
		if _, _, err := geneses.GenerateChannelTX(
			&genesis.ChannelTX{
				LedgerName:  channelCreate.LeagueName,
				ChannelName: channelCreate.ChannelID,
				Force:       true,
			}); nil != err {
			result.SayFail(router.Context, err.Error())
			return
		}
		channelTXFilePath := geneses.ChannelTXFilePath(channelCreate.LeagueName, channelCreate.ChannelID)

		sdk.Create(geneses.OrdererOrgName, "Admin", orderOrgURL, orgName, "Admin",
			channelCreate.ChannelID, channelTXFilePath, service.GetBytes(channelCreate.ConfigID)).Say(router.Context)
	})
}

func join(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var channelJoin = new(service.ChannelJoin)
		if err := router.Context.ShouldBindJSON(channelJoin); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		var (
			conf        *cfg.Config
			orderOrgURL string
		)
		if conf = service.Configs[channelJoin.ConfigID]; nil == conf {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		for _, order := range conf.Orderers {
			orderOrgURL = order.URL
		}
		sdk.Join(orderOrgURL, channelJoin.OrgName, channelJoin.OrgUser, channelJoin.ChannelID,
			service.GetBytes(channelJoin.ConfigID)).Say(router.Context)
	})
}

func list(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var channelList = new(service.ChannelList)
		if err := router.Context.ShouldBindJSON(channelList); err != nil {
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
