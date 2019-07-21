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

// candidate 用于选举Leader的一种角色
type candidate struct {
	// raft服务
	raft *Raft
	// persistence 所有角色都拥有的持久化的状态（在响应RPC请求之前变更且持久化的状态）
	persistence *persistence
	// nonPersistence 所有角色都拥有的非持久化的状态
	nonPersistence *nonPersistence
	// 发送发起投票协程池
	requestVotePool *ants.PoolWithFunc
}

func (c *candidate) become() {
	c.raft.role = roleCandidate
	c.raft.leader.release()
	c.raft.follower.release()
	c.requestVotePool, _ = ants.NewPoolWithFunc(len(c.raft.nodes), func(i interface{}) {
		c.requestVote(i)
	})
}

func (c *candidate) release() {
	_ = c.requestVotePool.Release()
}

// requestVote 发起选举，索要选票
func (c *candidate) requestVote(i interface{}) {
	node := i.(*Node)
	_, _ = rpc(node.URL, func(conn *grpc.ClientConn) (interface{}, error) {
		// 创建grpc客户端
		cli := NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if _, err := cli.RequestVote(context.Background(), &ReqVote{
			Term:        c.raft.term,
			CandidateId: c.raft.self.ID,
			// todo
		}); nil != err {
			return nil, err
		}
		return nil, nil
	})
}

// sendRequestVotes 批量发起选举，索要选票
func (c *candidate) sendRequestVotes() {
	c.requestVotePool.Tune(len(c.raft.nodes))
	// 遍历发送心跳
	for _, node := range c.raft.nodes {
		if node.ID == c.raft.self.ID {
			continue
		}
		if err := c.requestVotePool.Invoke(&node); nil != err {
			return
		}
	}
}
