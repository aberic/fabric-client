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
	"github.com/ennoo/fabric-client/core"
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	"github.com/ennoo/fabric-client/service"
	"github.com/ennoo/rivet"
	"github.com/ennoo/rivet/trans/response"
)

func Ledger(router *response.Router) {
	// 仓库相关路由设置
	router.Group = router.Engine.Group("/ledger")
	router.POST("/info", queryLedgerInfo)
	router.POST("/height", queryLedgerBlockByHeight)
	router.POST("/hash", queryLedgerBlockByHash)
	router.POST("/tx", queryLedgerBlockByTxID)
	router.POST("/info/spec", queryLedgerInfoSpec)
	router.POST("/height/spec", queryLedgerBlockByHeightSpec)
	router.POST("/hash/spec", queryLedgerBlockByHashSpec)
	router.POST("/tx/spec", queryLedgerBlockByTxIDSpec)
}

func queryLedgerInfo(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(pb.ReqInfo)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		sdk.QueryLedgerInfo(in.ConfigID, in.ChannelID, service.GetBytes(in.ConfigID)).Say(router.Context)
	})
}

func queryLedgerBlockByHeight(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(pb.ReqBlockByHeight)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		sdk.QueryLedgerBlockByHeight(in.ConfigID, in.ChannelID, in.Height, service.GetBytes(in.ConfigID)).Say(router.Context)
	})
}

func queryLedgerBlockByHash(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(pb.ReqBlockByHash)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		sdk.QueryLedgerBlockByHash(in.ConfigID, in.ChannelID, in.Hash, service.GetBytes(in.ConfigID)).Say(router.Context)
	})
}

func queryLedgerBlockByTxID(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(pb.ReqBlockByTxID)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		sdk.QueryLedgerBlockByTxID(in.ConfigID, in.ChannelID, in.TxID, service.GetBytes(in.ConfigID)).Say(router.Context)
	})
}

func queryLedgerInfoSpec(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(pb.ReqInfoSpec)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		sdk.QueryLedgerInfoSpec(in.ChannelID, in.OrgName, in.OrgUser, service.GetBytes(in.ConfigID)).Say(router.Context)
	})
}

func queryLedgerBlockByHeightSpec(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(pb.ReqBlockByHeightSpec)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		sdk.QueryLedgerBlockByHeightSpec(in.ChannelID, in.OrgName, in.OrgUser, in.Height, service.GetBytes(in.ConfigID)).Say(router.Context)
	})
}

func queryLedgerBlockByHashSpec(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(pb.ReqBlockByHashSpec)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		sdk.QueryLedgerBlockByHashSpec(in.ChannelID, in.OrgName, in.OrgUser, in.Hash, service.GetBytes(in.ConfigID)).Say(router.Context)
	})
}

func queryLedgerBlockByTxIDSpec(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(pb.ReqBlockByTxIDSpec)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		sdk.QueryLedgerBlockByTxIDSpec(in.ChannelID, in.OrgName, in.OrgUser, in.TxID, service.GetBytes(in.ConfigID)).Say(router.Context)
	})
}
