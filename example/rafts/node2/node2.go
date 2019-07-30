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

package main

import (
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	pbGeneses "github.com/ennoo/fabric-client/grpc/proto/geneses"
	"github.com/ennoo/fabric-client/grpc/server/chains"
	"github.com/ennoo/fabric-client/grpc/server/geneses"
	pbRafts "github.com/ennoo/fabric-client/rafts"
	"github.com/ennoo/rivet"
	"github.com/ennoo/rivet/utils/env"
	"github.com/ennoo/rivet/utils/log"
	"google.golang.org/grpc"
	"net"
)

func init() {
	rivet.Initialize(false, false, false, false)
	rivet.Log().Init(env.GetEnvDefault(env.LogPath, "./logs/node1"), "node1", &log.Config{
		Level:      log.InfoLevel,
		MaxSize:    128,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
	}, false)
}

func main() {
	go pbRafts.NewRaft()
	grpcListener()
}

func grpcListener() {
	var (
		listener net.Listener
		err      error
	)
	//  创建server端监听端口
	if listener, err = net.Listen("tcp", ":19878"); nil != err {
		panic(err)
	}
	//  创建grpc的server
	rpcServer := grpc.NewServer()

	//  注册我们自定义的helloworld服务
	pb.RegisterLedgerConfigServer(rpcServer, &chains.ConfigServer{})
	pb.RegisterLedgerChannelServer(rpcServer, &chains.ChannelServer{})
	pb.RegisterLedgerChainCodeServer(rpcServer, &chains.ChainCodeServer{})
	pb.RegisterLedgerServer(rpcServer, &chains.LedgerServer{})
	pbGeneses.RegisterGenesisServer(rpcServer, &geneses.GenesisServer{})
	pbRafts.RegisterRaftServer(rpcServer, &pbRafts.Server{})

	//  启动grpc服务
	if err = rpcServer.Serve(listener); nil != err {
		panic(err)
	}
}
