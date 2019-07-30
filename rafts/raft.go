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

package rafts

import (
	"github.com/ennoo/rivet/utils/env"
	"github.com/ennoo/rivet/utils/log"
	str "github.com/ennoo/rivet/utils/string"
	"os"
	"strings"
	"sync"
	"time"
)

// Raft 接收客户端提交的同步内容，被封装在自定义的方法中
//
// 也返回客户端期望的同步结果及从其他节点同步过来的信息
type Raft struct {
	// 服务器的任期，初始为0，递增
	term int32
	// 自身节点信息
	self *Node
	// 当前Raft可见节点集合
	nodes []*Node
	// persistence 所有角色都拥有的持久化的状态（在响应RPC请求之前变更且持久化的状态）
	persistence *persistence
	// 自身节点角色状态
	role roleChange
	// Raft任务调度服务
	scheduled *scheduled
	// 确保Raft的启动方法只会被调用一次
	once sync.Once
}

// roleChange 角色转换接口
type roleChange interface {
	// become 角色转换方法
	become(raft *Raft)
	// leader 角色切换成为leader
	leader()
	// candidate 角色切换成为candidate
	candidate()
	// follower 角色切换成为follower
	follower()
	// release 角色释放
	release()
	// role 获取角色信息
	role() int
	// work 开始本职工作
	work()
}

const (
	RoleLeader = iota
	RoleCandidate
	RoleFollower
)

const (
	K8S      = "K8S"          // K8S=true
	BrokerID = "BROKER_ID"    // BROKER_ID=1
	nodeAddr = "NODE_ADDRESS" // NODE_ADDRESS=example.com NODE_ADDRESS=127.0.0.1
	// CLUSTER=1=127.0.0.1:19865:19877,2=127.0.0.2:19865:19877,3=127.0.0.3:19865:19877
	cluster = "CLUSTER"
)

var (
	// instance Raft 实例
	instance *Raft
	// once 确保Raft的启动方法只会被调用一次
	once  sync.Once
	self  *Node
	nodes []*Node
)

func init() {
	// 仅测试用
	_ = os.Setenv(BrokerID, "1")
	_ = os.Setenv(nodeAddr, "127.0.0.1:19880")
	_ = os.Setenv(cluster, "1=127.0.0.1:19877,2=127.0.0.1:19878,3=127.0.0.1:19879")

	self = &Node{}
	if k8s := env.GetEnvBool(K8S); k8s {
		if self.Url = env.GetEnv("HOSTNAME"); str.IsEmpty(self.Url) {
			log.Self.Info("raft k8s fail", log.String("addr", self.Url))
			return
		}
		log.Self.Info("raft k8s", log.String("addr", self.Url))
		self.Id = strings.Split(self.Url, "-")[1]
		log.Self.Info("raft k8s", log.String("id", self.Id))
	} else {
		if self.Url = env.GetEnv(nodeAddr); str.IsEmpty(self.Url) {
			return
		}
		if self.Id = env.GetEnv(BrokerID); str.IsEmpty(self.Id) {
			log.Self.Error("raft", log.String("note", "broker id is not appoint"))
			return
		}
	}
	nodesStr := env.GetEnv(cluster)
	if str.IsEmpty(nodesStr) {
		nodes = make([]*Node, 0)
	} else {
		clusterArr := strings.Split(nodesStr, ",")
		nodes = make([]*Node, 0)
		for _, cluster := range clusterArr {
			clusterSplit := strings.Split(cluster, "=")
			id := clusterSplit[0]
			if str.IsEmpty(id) {
				log.Self.Error("raft", log.String("cluster", "broker id is nil"))
				continue
			}
			if id == self.Id {
				continue
			}
			nodeUrl := clusterSplit[1]
			nodes = append(nodes, &Node{
				Id:  id,
				Url: nodeUrl,
			})
		}
	}
}

func NewRaft() {
	log.Self.Info("raft", log.String("new", "新建Raft"))
	_ = obtainRaft()
}

func obtainRaft() *Raft {
	once.Do(func() {
		instance = &Raft{}
		instance.start()
	})
	return instance
}

// Start Raft启用方法
func (r *Raft) start() {
	r.once.Do(func() {
		r.initRaft(self, nodes)
	})
}

// initRaft Raft初始化
func (r *Raft) initRaft(self *Node, nodes []*Node) {
	log.Self.Info("raft", log.String("initRaft", "初始化Raft"))
	if nil == self || nil == nodes {
		log.Self.Info("raft", log.String("initRaft", "未组网或参数配置有误，raft集群无法启动"))
		return
	}
	r.term = 0
	r.self = self
	r.nodes = nodes
	r.persistence = &persistence{
		currentTerm: r.term,
		votedFor: &votedFor{
			id:        "",
			term:      0,
			timestamp: time.Now().UnixNano(),
		},
		version: 0,
	}
	r.role = &follower{raft: r}
	r.scheduled = &scheduled{
		raft: r,
	}
	r.role.become(r)
	r.scheduled.start()
}

func (r *Raft) appendNode(node *Node) {
	r.nodes = append(r.nodes, node)
}

func (r *Raft) release() {
	r.scheduled.tickerEnd <- 1
	r.scheduled.checkRelease <- 1
}

func Character() int {
	return obtainRaft().role.role()
}

func LeaderURL() string {
	for _, node := range obtainRaft().nodes {
		if node.Id == obtainRaft().persistence.leaderID {
			return node.Url
		}
	}
	return ""
}

func VersionAdd() {
	obtainRaft().persistence.version++
}
