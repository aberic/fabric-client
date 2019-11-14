/*
 * Copyright (c) 2019.. Aberic - All Rights Reserved.
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
	"github.com/aberic/fabric-client/config"
	"github.com/aberic/fabric-client/core"
	"github.com/aberic/fabric-client/geneses"
	pb "github.com/aberic/fabric-client/grpc/proto/chain"
	"github.com/aberic/fabric-client/service"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"golang.org/x/net/context"
)

type ChannelServer struct {
}

func (c *ChannelServer) Create(ctx context.Context, in *pb.ChannelCreate) (*pb.Result, error) {
	var (
		conf        *config.Config
		orderOrgURL string
		orgName     string
		txID        string
		err         error
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return &pb.Result{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, nil
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
	//if _, _, err = geneses.GenerateChannelTX(
	//	&genesis.ChannelTX{
	//		LedgerName:  in.LeagueName,
	//		ChannelName: in.ChannelID,
	//		Force:       true,
	//	}); nil != err {
	//	return &pb.Result{Code: pb.Code_Fail, ErrMsg: err.Error()}, nil
	//}
	channelTXFilePath := geneses.ChannelTXFilePath(in.LeagueName, in.ChannelID)
	if txID, err = sdk.Create(geneses.OrdererOrgName, "Admin", orderOrgURL, orgName, "Admin",
		in.ChannelID, channelTXFilePath, service.GetBytes(in.ConfigID)); nil != err {
		return &pb.Result{Code: pb.Code_Fail, ErrMsg: err.Error()}, nil
	}
	return &pb.Result{Code: pb.Code_Success, Data: txID}, nil
}

func (c *ChannelServer) Join(ctx context.Context, in *pb.ChannelJoin) (*pb.Result, error) {
	var (
		conf        *config.Config
		orderOrgURL string
		err         error
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return &pb.Result{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, nil
	}
	for _, order := range conf.Orderers {
		orderOrgURL = order.URL
	}
	if err = sdk.Join(orderOrgURL, in.OrgName, in.OrgUser, in.ChannelID, in.PeerName, service.GetBytes(in.ConfigID)); nil != err {
		return &pb.Result{Code: pb.Code_Fail, ErrMsg: err.Error()}, nil
	}
	return &pb.Result{Code: pb.Code_Success, Data: "success"}, nil
}

func (c *ChannelServer) List(ctx context.Context, in *pb.ChannelList) (*pb.ResultArr, error) {
	var (
		conf *config.Config
		chs  []*peer.ChannelInfo
		err  error
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return &pb.ResultArr{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, nil
	}
	if chs, err = sdk.Channels(in.OrgName, in.OrgUser, in.PeerName, service.GetBytes(in.ConfigID)); nil != err {
		return &pb.ResultArr{Code: pb.Code_Fail, ErrMsg: err.Error()}, nil
	}
	data := make([]string, len(chs))
	for index := range chs {
		data[index] = chs[index].ChannelId
	}
	return &pb.ResultArr{Code: pb.Code_Success, Data: data}, nil
}
