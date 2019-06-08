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

package route

import (
	"github.com/ennoo/fabric-go-client/service"
	"github.com/ennoo/rivet"
	"github.com/ennoo/rivet/trans/response"
	"strings"
)

func Config(router *response.Router) {
	// 仓库相关路由设置
	router.Group = router.Engine.Group("/config")
	router.GET("/get/:configID", get)
	router.POST("/client", initClient)
	router.POST("/client/custom", initClientCustom)
	router.POST("/channel/peer", addOrSetPeerForChannel)
	router.POST("/channel/policy/query", addOrSetQueryChannelPolicyForChannel)
	router.POST("/channel/policy/discovery", addOrSetDiscoveryPolicyForChannel)
	router.POST("/channel/policy/event", addOrSetEventServicePolicyForChannel)
	router.POST("/organizations/order", addOrSetOrdererForOrganizations)
	router.POST("/organizations/org", addOrSetOrgForOrganizations)
	router.POST("/order", addOrSetOrderer)
	router.POST("/peer", addOrSetPeer)
	router.POST("/ca", addOrSetCertificateAuthority)
}

func get(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		configID := router.Context.Param("configID")
		config := service.Get(configID)
		if nil == config {
			result.SayFail(router.Context, strings.Join([]string{"config ", configID, " is not exist"}, ""))
		} else {
			result.SaySuccess(router.Context, service.Get(configID))
		}
	})
}

func initClient(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var client = new(service.Client)
		if err := router.Context.ShouldBindJSON(client); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if err := service.InitClient(client); nil != err {
			result.SayFail(router.Context, err.Error())
			return
		}
		result.Say(router.Context)
	})
}

func initClientCustom(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var clientCustom = new(service.ClientCustom)
		if err := router.Context.ShouldBindJSON(clientCustom); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if err := service.InitClientCustom(clientCustom); nil != err {
			result.SayFail(router.Context, err.Error())
			return
		}
		result.Say(router.Context)
	})
}

func addOrSetPeerForChannel(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var channelPeer = new(service.ChannelPeer)
		if err := router.Context.ShouldBindJSON(channelPeer); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if err := service.AddOrSetPeerForChannel(channelPeer); nil != err {
			result.SayFail(router.Context, err.Error())
			return
		}
		result.Say(router.Context)
	})
}

func addOrSetQueryChannelPolicyForChannel(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var query = new(service.ChannelPolicyQuery)
		if err := router.Context.ShouldBindJSON(query); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if err := service.AddOrSetQueryChannelPolicyForChannel(query); nil != err {
			result.SayFail(router.Context, err.Error())
			return
		}
		result.Say(router.Context)
	})
}

func addOrSetDiscoveryPolicyForChannel(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var discovery = new(service.ChannelPolicyDiscovery)
		if err := router.Context.ShouldBindJSON(discovery); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if err := service.AddOrSetDiscoveryPolicyForChannel(discovery); nil != err {
			result.SayFail(router.Context, err.Error())
			return
		}
		result.Say(router.Context)
	})
}

func addOrSetEventServicePolicyForChannel(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var event = new(service.ChannelPolicyEvent)
		if err := router.Context.ShouldBindJSON(event); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if err := service.AddOrSetEventServicePolicyForChannel(event); nil != err {
			result.SayFail(router.Context, err.Error())
			return
		}
		result.Say(router.Context)
	})
}

func addOrSetOrdererForOrganizations(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var order = new(service.OrganizationsOrder)
		if err := router.Context.ShouldBindJSON(order); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if err := service.AddOrSetOrdererForOrganizations(order); nil != err {
			result.SayFail(router.Context, err.Error())
			return
		}
		result.Say(router.Context)
	})
}

func addOrSetOrgForOrganizations(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var org = new(service.OrganizationsOrg)
		if err := router.Context.ShouldBindJSON(org); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if err := service.AddOrSetOrgForOrganizations(org); nil != err {
			result.SayFail(router.Context, err.Error())
			return
		}
		result.Say(router.Context)
	})
}

func addOrSetOrderer(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var orderer = new(service.Orderer)
		if err := router.Context.ShouldBindJSON(orderer); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if err := service.AddOrSetOrderer(orderer); nil != err {
			result.SayFail(router.Context, err.Error())
			return
		}
		result.Say(router.Context)
	})
}

func addOrSetPeer(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var peer = new(service.Peer)
		if err := router.Context.ShouldBindJSON(peer); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if err := service.AddOrSetPeer(peer); nil != err {
			result.SayFail(router.Context, err.Error())
			return
		}
		result.Say(router.Context)
	})
}

func addOrSetCertificateAuthority(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var certificateAuthority = new(service.CertificateAuthority)
		if err := router.Context.ShouldBindJSON(certificateAuthority); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if err := service.AddOrSetCertificateAuthority(certificateAuthority); nil != err {
			result.SayFail(router.Context, err.Error())
			return
		}
		result.Say(router.Context)
	})
}
