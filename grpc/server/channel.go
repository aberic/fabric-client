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
	"github.com/ennoo/fabric-go-client/core"
	pb "github.com/ennoo/fabric-go-client/grpc/proto"
	"github.com/ennoo/fabric-go-client/service"
	"github.com/ennoo/rivet/trans/response"
	"golang.org/x/net/context"
)

type ChannelServer struct {
}

func (c *ChannelServer) Create(ctx context.Context, in *pb.ChannelCreate) (*pb.String, error) {
	var (
		res *response.Result
	)
	if res = sdk.Create(in.OrderOrgName, in.OrgName, in.OrgUser, in.ChannelID, in.ChannelConfigPath, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return &pb.String{Data: res.Data.(string)}, nil
	}
	return nil, errors.New(res.Msg)
}

func (c *ChannelServer) Join(ctx context.Context, in *pb.ChannelJoin) (*pb.String, error) {
	var (
		res *response.Result
	)
	if res = sdk.Join(in.OrderName, in.OrgName, in.OrgUser, in.ChannelID, in.PeerUrl, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return &pb.String{Data: res.Data.(string)}, nil
	}
	return nil, errors.New(res.Msg)
}

func (c *ChannelServer) List(ctx context.Context, in *pb.ChannelList) (*pb.StringArr, error) {
	var (
		res        *response.Result
		channelArr *sdk.ChannelArr
	)
	if res = sdk.Channels(in.OrgName, in.OrgUser, in.PeerName, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		channelArr = res.Data.(*sdk.ChannelArr)
		channels := channelArr.Channels
		data := make([]string, len(channels))
		for index := range channels {
			data[index] = channels[index].ChannelId
		}
		return &pb.StringArr{Data: data}, nil
	}
	return nil, errors.New(res.Msg)
}
