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
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	pbGeneses "github.com/ennoo/fabric-client/grpc/proto/geneses"
	pbRaft "github.com/ennoo/fabric-client/grpc/proto/raft"
	"github.com/ennoo/fabric-client/grpc/server/chains"
	"github.com/ennoo/fabric-client/grpc/server/geneses"
	"github.com/ennoo/fabric-client/grpc/server/raft"
	scheduled "github.com/ennoo/fabric-client/raft"
	"github.com/ennoo/fabric-client/route"
	"github.com/ennoo/rivet"
	"github.com/ennoo/rivet/utils/env"
	"github.com/ennoo/rivet/utils/log"
	str "github.com/ennoo/rivet/utils/string"
	"google.golang.org/grpc"
	"net"
	"strings"
)

func main() {
	if id := env.GetEnv(scheduled.BrokerID); str.IsNotEmpty(id) {
		scheduled.Start()
	} else if k8s := env.GetEnvBool(scheduled.K8S); k8s {
		scheduled.Start()
	}
	go httpListener()
	grpcListener()
}

func init() {
	var (
		level log.Level
	)
	rivet.Initialize(false, false, false, false)

	logLevel := strings.ToLower(env.GetEnvDefault("LOG_LEVEL", "warn"))
	switch logLevel {
	case "debug":
		level = log.DebugLevel
	case "info":
		level = log.InfoLevel
	case "warn":
		level = log.WarnLevel
	case "error":
		level = log.ErrorLevel
	}
	rivet.Log().Init(env.GetEnvDefault(env.LogPath, "./logs"), "fabric-client", &log.Config{
		Level:      level,
		MaxSize:    128,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
	}, false)
}

func httpListener() {
	rivet.ListenAndServe(&rivet.ListenServe{
		Engine: rivet.SetupRouter(
			route.Config,
			route.Channel,
			route.Order,
			route.ChainCode,
			route.Ledger,
		),
		DefaultPort: "19865",
	})
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
	pb.RegisterLedgerChannelServer(rpcServer, &chains.ChannelServer{})
	pb.RegisterLedgerChainCodeServer(rpcServer, &chains.ChainCodeServer{})
	pb.RegisterLedgerServer(rpcServer, &chains.LedgerServer{})
	pbGeneses.RegisterGenesisServer(rpcServer, &geneses.GenesisServer{})
	pbRaft.RegisterRaftServer(rpcServer, &rafts.RaftServer{})

	//  启动grpc服务
	if err = rpcServer.Serve(listener); nil != err {
		panic(err)
	}
}
