/*
 * Copyright (c) 2019. ENNOO - All Rights Reserved.
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
	"github.com/ennoo/rivet/utils/log"
	str "github.com/ennoo/rivet/utils/string"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v3"
)

type RaftsServer struct{}

// HeartBeat 发送心跳
func (r *RaftsServer) Heartbeat(_ context.Context, hBeat *HBeat) (hbr *HBeatReturn, err error) {
	log.Self.Debug("raft", log.Reflect("receive heartbeat", hBeat))
	hbr = &HBeatReturn{}
	if hBeat.Term < obtainRaft().term {
		hbr.Success = false
	} else if hBeat.Term == obtainRaft().term {
		switch obtainRaft().role.role() {
		case roleLeader:
			obtainRaft().role.candidate()
			hbr.Success = false
		case roleCandidate:
			obtainRaft().role.follower()
			r.syncConfig(hBeat)
			hbr.Success = true
		case roleFollower:
			r.syncConfig(hBeat)
			hbr.Success = true
		}
	} else if hBeat.Term > obtainRaft().term {
		switch obtainRaft().role.role() {
		case roleLeader, roleCandidate:
			obtainRaft().role.follower()
		}
		r.syncConfig(hBeat)
		hbr.Success = true
	}
	hbr.Term = obtainRaft().term
	return
}

// RequestVote 发起选举，索要选票
func (r *RaftsServer) RequestVote(_ context.Context, rv *ReqVote) (rvr *ReqVoteReturn, err error) {
	log.Self.Info("raft", log.Reflect("receive RequestVote", rv))
	rvr = &ReqVoteReturn{}
	rvr.Term = obtainRaft().term
	if rv.Term < obtainRaft().term {
		log.Self.Info("raft", log.Reflect("refuse", rv),
			log.Int32("termLocal", obtainRaft().term),
			log.Int32("termReceive", rv.Term))
		rvr.VoteGranted = false
	} else if rv.Term >= obtainRaft().term {
		rvr.VoteGranted = r.voteFor(rv)
	}
	r.syncNodes(&Node{
		Id:  rv.CandidateId,
		Url: rv.Url,
	})
	return
}

func (r *RaftsServer) syncNodes(node *Node) {
	have := false
	for _, n := range obtainRaft().nodes {
		if n.Id == node.Id {
			have = true
		}
	}
	if !have {
		obtainRaft().appendNode(node)
	}
}

// syncConfig 同步配置信息方案
func (r *RaftsServer) syncConfig(hBeat *HBeat) {
	if obtainRaft().persistence.leaderID == hBeat.LeaderId && obtainRaft().persistence.version == hBeat.Version && obtainRaft().term == hBeat.Term {
		obtainRaft().scheduled.refreshLastHeartBeatTime()
		return
	}
	if err := yaml.Unmarshal(hBeat.Config, obtainRaft().persistence.configs); nil != err {
		log.Self.Error("raft", log.Error(err))
	} else {
		log.Self.Debug("raft", log.String("syncConfig", "refresh time"))
		obtainRaft().term = hBeat.Term
		obtainRaft().persistence.leaderID = hBeat.LeaderId
		obtainRaft().persistence.version = hBeat.Version
		obtainRaft().persistence.currentTerm = hBeat.Term
		// todo config storage
		obtainRaft().scheduled.refreshLastHeartBeatTime()
	}
}

// voteFor 要求投票节点任期大于当前任期返回方案
func (r *RaftsServer) voteFor(rv *ReqVote) bool {
	if rv.Term == obtainRaft().term {
		if rv.Timestamp < obtainRaft().persistence.votedFor.timestamp {
			r.vote(rv)
			return true
		}
	}
	if rv.Term > obtainRaft().persistence.votedFor.term {
		r.vote(rv)
		return true
	}
	if rv.Term == obtainRaft().persistence.votedFor.term && str.IsEmpty(obtainRaft().persistence.votedFor.id) {
		r.vote(rv)
		return true
	}
	log.Self.Info("raft", log.Reflect("refuse", rv),
		log.Int32("termLocal", obtainRaft().term),
		log.Int32("termReceive", rv.Term))
	return false
}

func (r *RaftsServer) vote(rv *ReqVote) {
	obtainRaft().persistence.votedFor.id = rv.CandidateId
	obtainRaft().persistence.votedFor.term = rv.Term
	obtainRaft().scheduled.refreshLastHeartBeatTime()
	log.Self.Info("raft", log.Reflect("accept", rv),
		log.Int32("termLocal", obtainRaft().term),
		log.Int32("termReceive", rv.Term))
}
