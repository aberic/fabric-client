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
	str "github.com/ennoo/rivet/utils/string"
	"strings"
	"time"
)

const (
	brokerID = "BROKER_ID"    // BROKER_ID=1
	NodeAddr = "NODE_ADDRESS" // NODE_ADDRESS=example.com:19865:19877 NODE_ADDRESS=127.0.0.1:19865:19877
	// CLUSTER=1=127.0.0.1:19865:19877,2=127.0.0.2:19865:19877,3=127.0.0.3:19865:19877
	cluster = "CLUSTER"
)

var (
	Leader *pb.Leader
	Nodes  map[string]*pb.Node
	ID     string // ID ID 为空则表示不启用集群模式
	Addr   string // Addr Addr 为空则表示不启用集群模式
	Time   int64  // 最后一次心跳时间戳ms
	Term   int32  // 当前所处区间
)

func init() {
	if Addr = env.GetEnv(NodeAddr); str.IsEmpty(Addr) {
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
	nodesStr := env.GetEnv(cluster)
	if str.IsEmpty(nodesStr) {
		Nodes = map[string]*pb.Node{}
	} else {
		clusterArr := strings.Split(nodesStr, ",")
		Nodes = make(map[string]*pb.Node, len(clusterArr))
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
				Nodes[id].Status = pb.Status_FOLLOWER
			}
		}
	}
}

func RefreshTimeOut() {
	Time = time.Now().UnixNano() / 1e6
}
