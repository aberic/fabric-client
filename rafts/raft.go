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

package raft

import (
	"github.com/ennoo/rivet"
	"github.com/ennoo/rivet/utils/log"
	"sync"
)

// Raft 接收客户端提交的同步内容，被封装在自定义的方法中
//
// 也返回客户端期望的同步结果及从其他节点同步过来的信息
type Raft struct {
	// 服务器的任期，初始为0，递增
	term int32
	// 自身节点角色状态
	role int8
	// 自身节点信息
	self *Node
	// 当前Raft可见节点集合
	nodes []*Node
	// Raft中Leader角色
	leader *leader
	// 用于选举Leader的一种角色
	candidate *candidate
	// 负责响应来自Leader或者Candidate的请求角色
	follower *follower
	// Raft任务调度服务
	scheduled *scheduled
	// 确保Raft的启动方法只会被调用一次
	once sync.Once
}

const (
	roleLeader = iota
	roleCandidate
	roleFollower
)

// Start Raft启用方法
func (r *Raft) Start(logLevel, logPath string, self *Node, nodes []*Node) {
	r.once.Do(func() {
		var (
			level log.Level
		)
		rivet.Initialize(false, false, false, false)

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
		rivet.Log().Init(logPath, "raft", &log.Config{
			Level:      level,
			MaxSize:    128,
			MaxBackups: 30,
			MaxAge:     30,
			Compress:   true,
		}, false)
		r.initRaft(self, nodes)
	})
}

// initRaft Raft初始化
func (r *Raft) initRaft(self *Node, nodes []*Node) {
	r.term = 0
	r.self = self
	r.nodes = nodes
	persistence := &persistence{
		currentTerm: r.term,
		votedFor:    "",
		entries:     []Entry{},
	}
	nonPersistence := &nonPersistence{
		commitIndex: 0,
		lastApplied: 0,
	}
	r.leader = &leader{
		raft:           r,
		persistence:    persistence,
		nonPersistence: nonPersistence,
		leaderNonPersistence: &leaderNonPersistence{
			nextIndex:  make(map[string]nodeIndex),
			matchIndex: make(map[string]nodeIndex),
		},
	}
	r.candidate = &candidate{
		persistence:    persistence,
		nonPersistence: nonPersistence,
	}
	r.follower = &follower{
		persistence:    persistence,
		nonPersistence: nonPersistence,
	}
	r.follower.become()
	r.scheduled = &scheduled{
		raft: r,
	}
	r.scheduled.start()
}
