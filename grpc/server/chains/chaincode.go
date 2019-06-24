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
	"github.com/ennoo/fabric-client/core"
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	"github.com/ennoo/fabric-client/service"
	"github.com/ennoo/rivet/trans/response"
	"golang.org/x/net/context"
)

type ChainCodeServer struct {
}

func (c *ChainCodeServer) InstallCC(ctx context.Context, in *pb.Install) (*pb.String, error) {
	var (
		res *response.Result
	)
	if res = sdk.Install(in.OrgName, in.OrgUser, in.Name, in.Source, in.Path, in.Version, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return &pb.String{Data: res.Data.(string)}, nil
	}
	return nil, errors.New(res.Msg)
}

func (c *ChainCodeServer) InstalledCC(ctx context.Context, in *pb.Installed) (*pb.CCList, error) {
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

func (c *ChainCodeServer) InstantiateCC(ctx context.Context, in *pb.Instantiate) (*pb.String, error) {
	var (
		res *response.Result
	)
	if res = sdk.Instantiate(in.OrgName, in.OrgUser, in.ChannelID, in.Name, in.Path, in.Version, in.OrgPolicies,
		in.Args, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return &pb.String{Data: res.Data.(string)}, nil
	}
	return nil, errors.New(res.Msg)
}

func (c *ChainCodeServer) InstantiatedCC(ctx context.Context, in *pb.Instantiated) (*pb.CCList, error) {
	var (
		res              *response.Result
		chainCodeInfoArr *sdk.ChainCodeInfoArr
	)
	if res = sdk.Instantiated(in.OrgName, in.OrgUser, in.ChannelID, in.PeerName, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
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

func (c *ChainCodeServer) UpgradeCC(ctx context.Context, in *pb.Upgrade) (*pb.String, error) {
	var (
		res *response.Result
	)
	if res = sdk.Upgrade(in.OrgName, in.OrgUser, in.ChannelID, in.Name, in.Path, in.Version, in.OrgPolicies,
		in.Args, service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return &pb.String{Data: res.Data.(string)}, nil
	}
	return nil, errors.New(res.Msg)
}

func (c *ChainCodeServer) InvokeCC(ctx context.Context, in *pb.Invoke) (*pb.String, error) {
	var (
		res *response.Result
	)
	if res = sdk.Invoke(in.ChainCodeID, in.OrgName, in.OrgUser, in.ChannelID, in.Fcn, in.Args,
		service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return &pb.String{Data: res.Data.(string)}, nil
	}
	return nil, errors.New(res.Msg)
}

func (c *ChainCodeServer) QueryCC(ctx context.Context, in *pb.Query) (*pb.String, error) {
	var (
		res *response.Result
	)
	if res = sdk.Query(in.ChainCodeID, in.OrgName, in.OrgUser, in.ChannelID, in.Fcn, in.Args, in.TargetEndpoints,
		service.GetBytes(in.ConfigID)); res.ResultCode == response.Success {
		return &pb.String{Data: res.Data.(string)}, nil
	}
	return nil, errors.New(res.Msg)
}
