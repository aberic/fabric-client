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

// follower 负责响应来自Leader或者Candidate的请求
type follower struct {
	// raft服务
	raft *Raft
	// persistence 所有角色都拥有的持久化的状态（在响应RPC请求之前变更且持久化的状态）
	persistence *persistence
	// nonPersistence 所有角色都拥有的非持久化的状态
	nonPersistence *nonPersistence
}

func (f *follower) become() {
	f.raft.role = roleFollower
	f.raft.leader.release()
	f.raft.candidate.release()
}

func (f *follower) release() {

}
