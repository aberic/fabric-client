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

package server

import (
	"errors"
	"github.com/ennoo/fabric-client/core"
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	"github.com/ennoo/fabric-client/service"
	"github.com/ennoo/rivet/trans/response"
	"golang.org/x/net/context"
)

type LedgerServer struct {
}

func (l *LedgerServer) QueryLedgerInfo(ctx context.Context, in *pb.ReqInfo) (*pb.ChannelInfo, error) {
	var (
		res *response.Result
	)
	if res = sdk.QueryLedgerInfo(in.ConfigID, in.ChannelID, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return res.Data.(*pb.ChannelInfo), nil
	}
	return nil, errors.New(res.Msg)
}

func (l *LedgerServer) QueryLedgerBlockByHeight(ctx context.Context, in *pb.ReqBlockByHeight) (*pb.Block, error) {
	var (
		res *response.Result
	)
	if res = sdk.QueryLedgerBlockByHeight(in.ConfigID, in.ChannelID, in.Height, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return res.Data.(*pb.Block), nil
	}
	return nil, errors.New(res.Msg)
}

func (l *LedgerServer) QueryLedgerBlockByHash(ctx context.Context, in *pb.ReqBlockByHash) (*pb.Block, error) {
	var (
		res *response.Result
	)
	if res = sdk.QueryLedgerBlockByHash(in.ConfigID, in.ChannelID, in.Hash, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return res.Data.(*pb.Block), nil
	}
	return nil, errors.New(res.Msg)
}

func (l *LedgerServer) QueryLedgerBlockByTxID(ctx context.Context, in *pb.ReqBlockByTxID) (*pb.Block, error) {
	var (
		res *response.Result
	)
	if res = sdk.QueryLedgerBlockByTxID(in.ConfigID, in.ChannelID, in.TxID, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return res.Data.(*pb.Block), nil
	}
	return nil, errors.New(res.Msg)
}

func (l *LedgerServer) QueryLedgerInfoSpec(ctx context.Context, in *pb.ReqInfoSpec) (*pb.ChannelInfo, error) {
	var (
		res *response.Result
	)
	if res = sdk.QueryLedgerInfoSpec(in.ChannelID, in.OrgName, in.OrgUser, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return res.Data.(*pb.ChannelInfo), nil
	}
	return nil, errors.New(res.Msg)
}

func (l *LedgerServer) QueryLedgerBlockByHeightSpec(ctx context.Context, in *pb.ReqBlockByHeightSpec) (*pb.Block, error) {
	var (
		res *response.Result
	)
	if res = sdk.QueryLedgerBlockByHeightSpec(in.ChannelID, in.OrgName, in.OrgUser, in.Height, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return res.Data.(*pb.Block), nil
	}
	return nil, errors.New(res.Msg)
}

func (l *LedgerServer) QueryLedgerBlockByHashSpec(ctx context.Context, in *pb.ReqBlockByHashSpec) (*pb.Block, error) {
	var (
		res *response.Result
	)
	if res = sdk.QueryLedgerBlockByHashSpec(in.ChannelID, in.OrgName, in.OrgUser, in.Hash, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return res.Data.(*pb.Block), nil
	}
	return nil, errors.New(res.Msg)
}

func (l *LedgerServer) QueryLedgerBlockByTxIDSpec(ctx context.Context, in *pb.ReqBlockByTxIDSpec) (*pb.Block, error) {
	var (
		res *response.Result
	)
	if res = sdk.QueryLedgerBlockByTxIDSpec(in.ChannelID, in.OrgName, in.OrgUser, in.TxID, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return res.Data.(*pb.Block), nil
	}
	return nil, errors.New(res.Msg)
}
