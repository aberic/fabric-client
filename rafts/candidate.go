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
	"errors"
	"github.com/aberic/fabric-client/grpc/proto/utils"
	"github.com/aberic/gnomon"
	"github.com/panjf2000/ants"
	"google.golang.org/grpc"
	"sync"
	"time"
)

// candidate 用于选举Leader的一种角色
type candidate struct {
	// raft服务
	raft *Raft
	// 发送发起投票协程池
	requestVotePool *ants.PoolWithFunc
	// vote 发起投票结果
	vote *vote
}

// vote 发起投票结果
type vote struct {
	// lock 发起投票结果变更锁
	lock sync.Mutex
	// term 当前请求投票的任期
	term int32
	// voteIDsFrom 当前请求投票任期获得投票的节点id
	voteIDsFrom []string
	// voteChan 投票结果
	voteChan chan voteChan
	// voteDead 死亡节点
	voteDead chan int8
	// voteEnd 结束投票
	voteEnd chan error
}

type voteChan struct {
	// 投票节点id
	id string
	// 投票结果
	grant bool
}

func (c *candidate) become(raft *Raft) {
	gnomon.Log().Info("raft", gnomon.Log().Field("become", "Candidate"))
	c.raft = raft
	c.raft.term++
	c.raft.persistence.votedFor.id = c.raft.self.Id
	c.raft.persistence.votedFor.term = c.raft.term
	c.raft.persistence.votedFor.timestamp = time.Now().UnixNano()
	c.requestVotePool, _ = ants.NewPoolWithFunc(len(c.raft.nodes), func(i interface{}) {
		c.requestVote(i)
	})
	c.work()
}

func (c *candidate) leader() {
	c.release()
	c.raft.role = &leader{}
	c.raft.role.become(c.raft)
}

func (c *candidate) candidate() {}

func (c *candidate) follower() {
	c.release()
	c.raft.role = &follower{}
	c.raft.role.become(c.raft)
}

func (c *candidate) release() {
	if nil != c.vote {
		c.vote.voteEnd <- errors.New("candidate release")
	}
	c.requestVotePool.Release()
}

func (c *candidate) role() int {
	return RoleCandidate
}

func (c *candidate) work() {
	c.vote = &vote{
		term:        c.raft.term,
		voteIDsFrom: []string{c.raft.self.Id},
		voteChan:    make(chan voteChan, len(c.raft.nodes)),
		voteDead:    make(chan int8, len(c.raft.nodes)),
		voteEnd:     make(chan error),
	}
	c.sendRequestVotes()
	defer c.vote.lock.Unlock()
	c.vote.lock.Lock()
	position := 0 // 遍历数
	dead := 0     // 死亡数
	for position < len(c.raft.nodes) {
		select {
		case voteResult := <-c.vote.voteChan:
			if voteResult.grant {
				c.vote.voteIDsFrom = append(c.vote.voteIDsFrom, voteResult.id)
			}
			position++
		case <-c.vote.voteDead:
			dead++
			position++
		case err := <-c.vote.voteEnd:
			gnomon.Log().Error("raft", gnomon.Log().Err(err))
			position = len(c.raft.nodes)
			c.vote = nil
			c.follower()
			return
		}
	}
	if position-dead == 0 || 100*len(c.vote.voteIDsFrom)/(position-dead) > 50 {
		c.vote = nil
		c.leader()
	} else {
		c.vote = nil
		c.follower()
	}
}

// requestVote 发起选举，索要选票
func (c *candidate) requestVote(i interface{}) {
	var (
		rvr    interface{}
		result *ReqVoteReturn
		err    error
	)
	node := i.(*Node)
	rvr, err = utils.RPC(node.Url, func(conn *grpc.ClientConn) (interface{}, error) {
		// 创建grpc客户端
		cli := NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if result, err = cli.RequestVote(context.Background(), &ReqVote{
			Term:         c.raft.term,
			CandidateId:  c.raft.self.Id,
			Url:          c.raft.self.Url,
			LastLeaderId: c.raft.persistence.leaderID,
			LastTerm:     c.raft.persistence.currentTerm,
			LastVersion:  c.raft.persistence.version,
			Timestamp:    c.raft.persistence.votedFor.timestamp,
		}); nil != err {
			return nil, err
		}
		return result, nil
	})
	if nil != err && nil != c.vote.voteDead {
		c.vote.voteDead <- 0
		return
	}
	reqVoteReturn := rvr.(*ReqVoteReturn)
	if reqVoteReturn.VoteGranted && nil != c.vote.voteChan { // 如果投票
		c.vote.voteChan <- voteChan{
			id:    node.Id,
			grant: true,
		}
	} else if nil != c.vote.voteChan {
		c.vote.voteChan <- voteChan{
			grant: false,
		}
	}
}

// sendRequestVotes 批量发起选举，索要选票
func (c *candidate) sendRequestVotes() {
	c.requestVotePool.Tune(len(c.raft.nodes))
	gnomon.Log().Info("raft", gnomon.Log().Field("send requestVotes", c.raft.nodes))
	// 遍历发送心跳
	for _, node := range c.raft.nodes {
		if node.Id == c.raft.self.Id {
			continue
		}
		if err := c.requestVotePool.Invoke(node); nil != err {
			return
		}
	}
}
