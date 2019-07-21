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

package raft

import (
	"context"
	"github.com/panjf2000/ants"
	"google.golang.org/grpc"
)

// leader 负责接收客户端的请求，将日志复制到其他节点并告知其他节点何时应用这些日志是安全的
type leader struct {
	// raft服务
	raft *Raft
	// persistence 所有角色都拥有的持久化的状态（在响应RPC请求之前变更且持久化的状态）
	persistence *persistence
	// nonPersistence 所有角色都拥有的非持久化的状态
	nonPersistence *nonPersistence
	// leaderNonPersistence Leader节点上非持久化的状态（选举后重新初始化）
	leaderNonPersistence *leaderNonPersistence
	// 发送心跳协程池
	heartBeatPool *ants.PoolWithFunc
}

func (l *leader) become() {
	l.raft.role = roleLeader
	l.raft.candidate.release()
	l.raft.follower.release()
	l.heartBeatPool, _ = ants.NewPoolWithFunc(len(l.raft.nodes), func(i interface{}) {
		l.heartBeat(i)
	})
	l.raft.scheduled.tickerStart()
}

func (l *leader) release() {
	_ = l.heartBeatPool.Release()
	if nil != l.raft.scheduled.tickerEnd {
		l.raft.scheduled.tickerEnd <- 1
	}
}

// HB 组合心跳发送参数
type HB struct {
	// Node 节点信息
	node *Node
	// AppendEntries 用于Leader节点复制日志给其他节点，也作为心跳
	appendEntries *AppendEntries
}

// heartBeat 发送心跳
func (l *leader) heartBeat(i interface{}) {
	hb := i.(*HB)
	_, _ = rpc(hb.node.URL, func(conn *grpc.ClientConn) (interface{}, error) {
		// 创建grpc客户端
		cli := NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if _, err := cli.HeartBeat(context.Background(), hb.appendEntries); nil != err {
			return nil, err
		}
		return nil, nil
	})
}

// sendHeartBeats 遍历发送心跳
func (l *leader) sendHeartBeats() {
	entryLen := len(l.persistence.entries)
	appendEntries := &AppendEntries{
		Term:         l.raft.term,
		LeaderId:     l.raft.self.ID,
		PrevLogIndex: l.persistence.entries[entryLen-1].Index,
		PrevLogTerm:  l.persistence.entries[entryLen-1].Term,
		Entries:      []*Entry{},
		// todo
	}
	l.heartBeatPool.Tune(len(l.raft.nodes))
	// 遍历发送心跳
	for _, node := range l.raft.nodes {
		if node.ID == l.raft.self.ID {
			continue
		}
		if err := l.heartBeatPool.Invoke(&HB{
			node:          node,
			appendEntries: appendEntries,
		}); nil != err {
			return
		}
	}
}
