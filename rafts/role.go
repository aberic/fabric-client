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

// role
//
// 所有节点初始状态都是Follower角色
//
// 超时时间内没有收到Leader的请求则转换为Candidate进行选举
//
// Candidate收到大多数节点的选票则转换为Leader；发现Leader或者收到更高任期的请求则转换为Follower
//
// Leader在收到更高任期的请求后转换为Follower
//
// Raft把时间切割为任意长度的任期（term），每个任期都有一个任期号，采用连续的整数

package raft

// Node 节点信息
type Node struct {
	ID  string // 节点ID
	URL string // 节点地址
}

// persistence 所有角色都拥有的持久化的状态（在响应RPC请求之前变更且持久化的状态）
type persistence struct {
	currentTerm int32   // 服务器的任期，初始为0，递增
	votedFor    string  // 在当前获得选票的候选人的 Id
	entries     []Entry // 日志条目内容 // 日志条目集；每一个条目包含一个用户状态机执行的指令，和收到时的任期号
}

// nonPersistence 所有角色都拥有的非持久化的状态
type nonPersistence struct {
	commitIndex int64 // 最大的已经被commit的日志的index
	lastApplied int64 // 最大的已经被应用到状态机的index
}

// leaderNonPersistence Leader节点上非持久化的状态（选举后重新初始化）
type leaderNonPersistence struct {
	nextIndex  map[string]nodeIndex // 每个节点下一次应该接收的日志的index（初始化为Leader节点最后一个日志的Index + 1）
	matchIndex map[string]nodeIndex // 每个节点已经复制的日志的最大的索引（初始化为0，之后递增）
}

// nodeIndex 节点与索引
type nodeIndex struct {
	node  Node
	index int64
}

// roleChange 角色转换接口
type roleChange interface {
	// become 角色转换方法
	become()
	// release 角色释放
	release()
}
