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

package main

import (
	"github.com/aberic/gnomon"
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	pbGeneses "github.com/ennoo/fabric-client/grpc/proto/geneses"
	"github.com/ennoo/fabric-client/grpc/server/chains"
	"github.com/ennoo/fabric-client/grpc/server/geneses"
	"github.com/ennoo/fabric-client/rafts"
	"google.golang.org/grpc"
	"net"
)

const (
	// LogPath 日志文件输出路径
	LogPath = "LOG_PATH"
)

func main() {
	if id := gnomon.Env().Get(rafts.BrokerID); gnomon.String().IsNotEmpty(id) {
		gnomon.Log().Info("raft self", gnomon.Log().Field("BrokerID", id))
		rafts.NewRaft()
	} else if k8s := gnomon.Env().GetBool(rafts.K8S); k8s {
		gnomon.Log().Info("raft k8s")
		rafts.NewRaft()
	}
	grpcListener()
}

func init() {
	logSet()
}

// logSet 日志设置
func logSet() {
	if err := gnomon.Log().Init(gnomon.Env().GetD(LogPath, "./logs"), 50, 7, false); nil != err {
		panic(err)
	}
	gnomon.Log().Set(gnomon.Log().WarnLevel(), true)
}

func grpcListener() {
	var (
		listener net.Listener
		err      error
	)
	//  创建server端监听端口
	if listener, err = net.Listen("tcp", ":19877"); nil != err {
		panic(err)
	}
	//  创建grpc的server
	rpcServer := grpc.NewServer()

	//  注册我们自定义的helloworld服务
	pb.RegisterLedgerConfigServer(rpcServer, &chains.ConfigServer{})
	pb.RegisterLedgerCAServer(rpcServer, &chains.CAServer{})
	pb.RegisterLedgerChannelServer(rpcServer, &chains.ChannelServer{})
	pb.RegisterLedgerChainCodeServer(rpcServer, &chains.ChainCodeServer{})
	pb.RegisterLedgerPeerServer(rpcServer, &chains.PeerServer{})
	pb.RegisterLedgerServer(rpcServer, &chains.LedgerServer{})
	pbGeneses.RegisterGenesisServer(rpcServer, &geneses.GenesisServer{})
	rafts.RegisterRaftServer(rpcServer, &rafts.Server{})

	//  启动grpc服务
	if err = rpcServer.Serve(listener); nil != err {
		panic(err)
	}
}
