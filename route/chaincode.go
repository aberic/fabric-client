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

func ChainCode(router *response.Router) {
	// 仓库相关路由设置
	router.Group = router.Engine.Group("/code")
	router.POST("/install", install)
	router.POST("/installed", installed)
	router.POST("/instantiate", instantiate)
	router.POST("/instantiated", instantiated)
	router.POST("/upgrade", upgrade)
	router.POST("/invoke", invoke)
	router.POST("/query", query)
}

func install(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var install = new(service.Install)
		if err := router.Context.ShouldBindJSON(install); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if nil == service.Get(install.ConfigID) {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		sdk.Install(install.OrderOrgName, install.OrgAdmin, install.Name, install.Source, install.Path, install.Version,
			service.GetBytes(install.ConfigID)).Say(router.Context)
	})
}

func installed(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var installed = new(service.Installed)
		if err := router.Context.ShouldBindJSON(installed); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if nil == service.Get(installed.ConfigID) {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		sdk.Installed(installed.OrgName, installed.OrgAdmin, installed.PeerName,
			service.GetBytes(installed.ConfigID)).Say(router.Context)
	})
}

func instantiate(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var instantiate = new(service.Instantiate)
		if err := router.Context.ShouldBindJSON(instantiate); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if nil == service.Get(instantiate.ConfigID) {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		sdk.Instantiate(instantiate.OrderOrgName, instantiate.OrgAdmin, instantiate.ChannelID, instantiate.Name,
			instantiate.Path, instantiate.Version, instantiate.OrgPolicies, instantiate.Args,
			service.GetBytes(instantiate.ConfigID)).Say(router.Context)
	})
}

func instantiated(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var instantiated = new(service.Instantiated)
		if err := router.Context.ShouldBindJSON(instantiated); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if nil == service.Get(instantiated.ConfigID) {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		sdk.Instantiated(instantiated.OrgName, instantiated.OrgAdmin, instantiated.ChannelID, instantiated.PeerName,
			service.GetBytes(instantiated.ConfigID)).Say(router.Context)
	})
}

func upgrade(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var upgrade = new(service.Upgrade)
		if err := router.Context.ShouldBindJSON(upgrade); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if nil == service.Get(upgrade.ConfigID) {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		sdk.Instantiate(upgrade.OrderOrgName, upgrade.OrgAdmin, upgrade.ChannelID, upgrade.Name,
			upgrade.Path, upgrade.Version, upgrade.OrgPolicies, upgrade.Args,
			service.GetBytes(upgrade.ConfigID)).Say(router.Context)
	})
}

func invoke(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var invoke = new(service.Invoke)
		if err := router.Context.ShouldBindJSON(invoke); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if nil == service.Get(invoke.ConfigID) {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		sdk.Invoke(invoke.ChainCodeID, invoke.OrgName, invoke.OrgUser, invoke.ChannelID, invoke.Fcn, invoke.Args,
			service.GetBytes(invoke.ConfigID)).Say(router.Context)
	})
}

func query(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var query = new(service.Query)
		if err := router.Context.ShouldBindJSON(query); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		if nil == service.Get(query.ConfigID) {
			result.SayFail(router.Context, "config client is not exist")
			return
		}
		if nil == query.TargetEndpoints {
			query.TargetEndpoints = []string{}
		}
		sdk.Query(query.ChainCodeID, query.OrgName, query.OrgUser, query.ChannelID, query.Fcn, query.Args, query.TargetEndpoints,
			service.GetBytes(query.ConfigID)).Say(router.Context)
	})
}
