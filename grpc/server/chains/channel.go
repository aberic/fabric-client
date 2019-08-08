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

package chains

import (
	"errors"
	"github.com/ennoo/fabric-client/config"
	"github.com/ennoo/fabric-client/core"
	"github.com/ennoo/fabric-client/geneses"
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	genesis "github.com/ennoo/fabric-client/grpc/proto/geneses"
	"github.com/ennoo/fabric-client/service"
	"github.com/ennoo/rivet/trans/response"
	"golang.org/x/net/context"
)

type ChannelServer struct {
}

func (c *ChannelServer) Create(ctx context.Context, in *pb.ChannelCreate) (*pb.Result, error) {
	var (
		res         *response.Result
		conf        *config.Config
		orderOrgURL string
		orgName     string
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return nil, errors.New("config client is not exist")
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
			LedgerName:  in.LeagueName,
			ChannelName: in.ChannelID,
			Force:       true,
		}); nil != err {
		return nil, err
	}
	channelTXFilePath := geneses.ChannelTXFilePath(in.LeagueName, in.ChannelID)
	if res = sdk.Create(geneses.OrdererOrgName, "Admin", orderOrgURL, orgName, "Admin",
		in.ChannelID, channelTXFilePath, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return &pb.Result{Code: pb.Code_Success, Data: res.Data.(string)}, nil
	}
	return &pb.Result{Code: pb.Code_Fail, ErrMsg: res.Msg}, errors.New(res.Msg)
}

func (c *ChannelServer) Join(ctx context.Context, in *pb.ChannelJoin) (*pb.Result, error) {
	var (
		res         *response.Result
		conf        *config.Config
		orderOrgURL string
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return nil, errors.New("config client is not exist")
	}
	for _, order := range conf.Orderers {
		orderOrgURL = order.URL
	}
	if res = sdk.Join(orderOrgURL, in.OrgName, in.OrgUser, in.ChannelID, in.PeerName, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return &pb.Result{Code: pb.Code_Success, Data: res.Data.(string)}, nil
	}
	return &pb.Result{Code: pb.Code_Fail, ErrMsg: res.Msg}, errors.New(res.Msg)
}

func (c *ChannelServer) List(ctx context.Context, in *pb.ChannelList) (*pb.ResultArr, error) {
	var (
		res        *response.Result
		conf       *config.Config
		channelArr *sdk.ChannelArr
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return nil, errors.New("config client is not exist")
	}
	if res = sdk.Channels(in.OrgName, in.OrgUser, in.PeerName, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		channelArr = res.Data.(*sdk.ChannelArr)
		channels := channelArr.Channels
		data := make([]string, len(channels))
		for index := range channels {
			data[index] = channels[index].ChannelId
		}
		return &pb.ResultArr{Code: pb.Code_Success, Data: data}, nil
	}
	return &pb.ResultArr{Code: pb.Code_Fail, ErrMsg: res.Msg}, errors.New(res.Msg)
}
