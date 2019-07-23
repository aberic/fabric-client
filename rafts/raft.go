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
	"github.com/ennoo/fabric-client/config"
	"github.com/ennoo/rivet"
	"github.com/ennoo/rivet/utils/log"
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
	roleLeader = iota
	roleCandidate
	roleFollower
)

var (
	// instance Raft 实例
	instance *Raft
	// once 确保Raft的启动方法只会被调用一次
	once          sync.Once
	logLevelLocal string
	logPathLocal  string
	selfLocal     *Node
	nodesLocal    []*Node
)

func NewRaft(logLevel, logPath string, self *Node, nodes []*Node) {
	log.Self.Info("raft", log.String("new", "新建Raft"))
	logLevelLocal = logLevel
	logPathLocal = logPath
	selfLocal = self
	nodesLocal = nodes
	_ = obtainRaft()
}

func obtainRaft() *Raft {
	once.Do(func() {
		instance = &Raft{}
		instance.start(logLevelLocal, logPathLocal, selfLocal, nodesLocal)
	})
	return instance
}

// Start Raft启用方法
func (r *Raft) start(logLevel, logPath string, self *Node, nodes []*Node) {
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
	log.Self.Info("raft", log.String("init", "初始化Raft"))
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
		configs: make(map[string]*config.Config),
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
