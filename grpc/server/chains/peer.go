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

package chains

import (
	"github.com/ennoo/fabric-client/core"
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	"github.com/ennoo/fabric-client/grpc/proto/utils"
	"github.com/ennoo/fabric-client/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type PeerServer struct {
}

func (p *PeerServer) LocalPeers(ctx context.Context, in *pb.ReqLocalPeers) (*pb.ResultPeers, error) {
	if fabPeers, err := sdk.DiscoveryLocalPeers(in.OrgName, in.OrgUser, service.GetBytes(in.ConfigID)); nil != err {
		return &pb.ResultPeers{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	} else {
		var peers []*pb.DiscoveryPeer
		for _, peer := range fabPeers {
			peers = append(peers, &pb.DiscoveryPeer{
				MspID: peer.MSPID(),
				Url:   peer.URL(),
			})
		}
		return &pb.ResultPeers{Code: pb.Code_Success, Peer: peers}, nil
	}
}

func (p *PeerServer) ChannelPeers(ctx context.Context, in *pb.ReqChannelPeers) (*pb.ResultPeers, error) {
	if fabPeers, err := sdk.DiscoveryChannelPeers(in.ChannelID, in.OrgName, in.OrgUser, service.GetBytes(in.ConfigID)); nil != err {
		return &pb.ResultPeers{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	} else {
		var peers []*pb.DiscoveryPeer
		for _, peer := range fabPeers {
			peers = append(peers, &pb.DiscoveryPeer{
				MspID: peer.MSPID(),
				Url:   peer.URL(),
			})
		}
		return &pb.ResultPeers{Code: pb.Code_Success, Peer: peers}, nil
	}
}

// LocalPeers LocalPeers
func LocalPeers(url string, req *pb.ReqLocalPeers) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.ResultPeers
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerPeerClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.LocalPeers(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// ChannelPeers ChannelPeers
func ChannelPeers(url string, req *pb.ReqChannelPeers) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.ResultPeers
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerPeerClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.ChannelPeers(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}
