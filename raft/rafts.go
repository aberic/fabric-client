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

package raft

import (
	pb "github.com/ennoo/fabric-client/grpc/proto/raft"
	"github.com/ennoo/rivet/utils/env"
	"github.com/ennoo/rivet/utils/log"
	"github.com/ennoo/rivet/utils/string"
	"os"
	"strings"
)

const (
	brokerID = "BROKER_ID"    // BROKER_ID=1
	nodeAddr = "NODE_ADDRESS" // NODE_ADDRESS=example.com NODE_ADDRESS=127.0.0.1
	// CLUSTER=1=127.0.0.1:19865:19877,2=127.0.0.2:19865:19877,3=127.0.0.3:19865:19877
	cluster = "CLUSTER"
	timeOut = 1500 // raft心跳超时ms
)

var (
	Leader    *pb.Leader
	Nodes     map[string]*pb.Node
	ID        string // ID ID 为空则表示不启用集群模式
	Addr      string // Addr Addr 为空则表示不启用集群模式
	Term      int32  // 当前所处区间
	voteCount int32  // 获取到的票数
	Sync      bool   // 是否同步过
)

func init() {
	// 仅测试用
	_ = os.Setenv(brokerID, "8")
	_ = os.Setenv(nodeAddr, "127.0.0.1")
	_ = os.Setenv(cluster, "1=127.0.0.1:19865:19877,2=127.0.0.1:19866:19878,3=127.0.0.1:19867:19879,"+
		"4=127.0.0.1:19868:19880,5=127.0.0.1:19869:19881,6=127.0.0.1:19870:19882,7=127.0.0.1:19871:19883")

	if Addr = env.GetEnv(nodeAddr); str.IsEmpty(Addr) {
		ID = ""
		return
	}
	if id := env.GetEnv(brokerID); str.IsEmpty(id) {
		log.Self.Error("raft", log.String("note", "broker id is not appoint"))
		ID = ""
		return
	} else {
		ID = id
	}
	Term = 0
	voteCount = 0
	Sync = false
	Leader = &pb.Leader{}
	nodesStr := env.GetEnv(cluster)
	if str.IsEmpty(nodesStr) {
		Nodes = map[string]*pb.Node{}
	} else {
		clusterArr := strings.Split(nodesStr, ",")
		Nodes = make(map[string]*pb.Node, len(clusterArr))
		haveOne := false
		for _, cluster := range clusterArr {
			clusterSplit := strings.Split(cluster, "=")
			id := clusterSplit[0]
			if str.IsEmpty(id) {
				log.Self.Error("raft", log.String("cluster", "broker id is nil"))
				continue
			}
			nodeStr := clusterSplit[1]
			nodeArr := strings.Split(nodeStr, ":")
			Nodes[id] = &pb.Node{
				Id:   id,
				Addr: nodeArr[0],
				Http: nodeArr[1],
				Rpc:  nodeArr[2],
			}
			if id == ID {
				haveOne = true
				// 默认初始将自身定义为FOLLOWER
				Nodes[id].Status = pb.Status_FOLLOWER
			}
		}
		if !haveOne {
			Nodes[ID] = &pb.Node{
				Id:     ID,
				Addr:   Addr,
				Http:   "19872",
				Rpc:    "19884",
				Status: pb.Status_FOLLOWER,
			}
		}
	}
}
