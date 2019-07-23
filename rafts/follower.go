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

import "github.com/ennoo/rivet/utils/log"

// follower 负责响应来自Leader或者Candidate的请求
type follower struct {
	// raft服务
	raft *Raft
}

func (f *follower) become(raft *Raft) {
	log.Self.Info("raft", log.String("become", "Follower"))
	f.raft = raft
	f.raft.persistence.votedFor.id = ""
	f.raft.persistence.votedFor.term = 0
	f.raft.scheduled.refreshLastHeartBeatTime()
}

func (f *follower) leader() {
	f.release()
	f.raft.role = &leader{}
	f.raft.role.become(f.raft)
}

func (f *follower) candidate() {
	f.release()
	f.raft.role = &candidate{}
	f.raft.role.become(f.raft)
}

func (f *follower) follower() {}

func (f *follower) release() {}

func (f *follower) role() int {
	return roleFollower
}

func (f *follower) work() {}
