/*
 * Copyright (c) 2019. Aberic - All Rights Reserved.
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
	"context"
	"github.com/ennoo/fabric-client/grpc/proto/utils"
	"github.com/ennoo/fabric-client/service"
	"github.com/ennoo/rivet/utils/log"
	"github.com/panjf2000/ants"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

// leader 负责接收客户端的请求，将日志复制到其他节点并告知其他节点何时应用这些日志是安全的
type leader struct {
	// raft服务
	raft *Raft
	// 发送心跳协程池
	heartBeatPool *ants.PoolWithFunc
}

func (l *leader) become(raft *Raft) {
	log.Self.Info("raft", log.String("become", "Leader"))
	l.raft = raft
	l.raft.persistence.version = 0
	l.raft.persistence.currentTerm = l.raft.term
	l.raft.persistence.leaderID = l.raft.self.Id
	l.raft.persistence.votedFor.id = ""
	l.raft.persistence.votedFor.term = 0
	l.heartBeatPool, _ = ants.NewPoolWithFunc(len(l.raft.nodes), func(i interface{}) {
		l.heartbeat(i)
	})
	l.raft.scheduled.tickerStart()
}

func (l *leader) leader() {}

func (l *leader) candidate() {
	l.release()
	l.raft.role = &candidate{}
	l.raft.role.become(l.raft)
}

func (l *leader) follower() {
	l.release()
	l.raft.role = &follower{}
	l.raft.role.become(l.raft)
}

func (l *leader) release() {
	_ = l.heartBeatPool.Release()
	if nil != l.raft.scheduled.tickerEnd {
		l.raft.scheduled.tickerEnd <- 1
	}
}

func (l *leader) role() int {
	return RoleLeader
}

func (l *leader) work() {
	l.sendHeartbeats()
}

// HB 组合心跳发送参数
type HB struct {
	// Node 节点信息
	node *Node
	// AppendEntries 用于Leader节点复制日志给其他节点，也作为心跳
	hBeat *HBeat
}

// heartBeat 发送心跳
func (l *leader) heartbeat(i interface{}) {
	var (
		hbr    interface{}
		result *HBeatReturn
		err    error
	)
	hb := i.(*HB)
	hbr, err = utils.RPC(hb.node.Url, func(conn *grpc.ClientConn) (interface{}, error) {
		// 创建grpc客户端
		cli := NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if result, err = cli.Heartbeat(context.Background(), hb.hBeat); nil != err {
			return nil, err
		}
		return result, nil
	})
	if nil != err {
		log.Self.Warn("raft", log.Error(err))
		return
	}
	if heartbeatReturn := hbr.(*HBeatReturn); !heartbeatReturn.Success {
		l.raft.term = heartbeatReturn.Term
		l.candidate()
	}
}

// sendHeartBeats 遍历发送心跳
func (l *leader) sendHeartbeats() {
	var cs = service.Configs
	configStr, err := yaml.Marshal(cs)
	if nil != err {
		return
	}
	hBeat := &HBeat{
		Term:     l.raft.term,
		LeaderId: l.raft.self.Id,
		Version:  l.raft.persistence.version,
		Config:   configStr,
	}
	l.heartBeatPool.Tune(len(l.raft.nodes))
	log.Self.Debug("raft", log.Reflect("send heartbeat", hBeat), log.Reflect("nodes", l.raft.nodes))
	// 遍历发送心跳
	for _, node := range l.raft.nodes {
		if node.Id == l.raft.self.Id {
			continue
		}
		if err := l.heartBeatPool.Invoke(&HB{
			node:  node,
			hBeat: hBeat,
		}); nil != err {
			return
		}
	}
}
