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

type LedgerServer struct {
}

func (c *LedgerServer) QueryLedgerInfo(ctx context.Context, in *pb.ReqInfo) (*pb.ChannelInfo, error) {
	var (
		res      *response.Result
		orgName  string
		userName string
	)

	if res = sdk.QueryLedgerInfo(in.ChannelID, orgName, userName, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return &pb.String{Data: res.Data.(string)}, nil
	}
	return nil, errors.New(res.Msg)
}

func (c *LedgerServer) InstalledCC(ctx context.Context, in *pb.Installed) (*pb.CCList, error) {
	var (
		res              *response.Result
		chainCodeInfoArr *sdk.ChainCodeInfoArr
	)
	if res = sdk.Installed(in.OrgName, in.OrgUser, in.PeerName, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		chainCodeInfoArr = res.Data.(*sdk.ChainCodeInfoArr)
		chainCodes := chainCodeInfoArr.ChainCodes
		data := make([]*pb.ChainCodeInfo, len(chainCodes))
		for index := range chainCodes {
			data[index].Name = chainCodes[index].Name
			data[index].Version = chainCodes[index].Name
			data[index].Path = chainCodes[index].Name
			data[index].Input = chainCodes[index].Name
			data[index].Escc = chainCodes[index].Name
			data[index].Vscc = chainCodes[index].Name
			data[index].Id = chainCodes[index].Id
		}
		return &pb.CCList{Data: data}, nil
	}
	return nil, errors.New(res.Msg)
}
