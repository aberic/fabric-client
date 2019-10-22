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

package rafts

import (
	"github.com/aberic/gnomon"
	"github.com/ennoo/fabric-client/config"
	"github.com/ennoo/fabric-client/service"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v3"
)

type Server struct{}

// HeartBeat 发送心跳
func (s *Server) Heartbeat(_ context.Context, hBeat *HBeat) (hbr *HBeatReturn, err error) {
	gnomon.Log().Debug("raft", gnomon.Log().Field("receive heartbeat", hBeat))
	hbr = &HBeatReturn{}
	if hBeat.Term < obtainRaft().term {
		hbr.Success = false
	} else if hBeat.Term == obtainRaft().term {
		switch obtainRaft().role.role() {
		case RoleLeader:
			obtainRaft().role.candidate()
			hbr.Success = false
		case RoleCandidate:
			obtainRaft().role.follower()
			s.syncConfig(hBeat)
			hbr.Success = true
		case RoleFollower:
			s.syncConfig(hBeat)
			hbr.Success = true
		}
	} else if hBeat.Term > obtainRaft().term {
		switch obtainRaft().role.role() {
		case RoleLeader, RoleCandidate:
			obtainRaft().role.follower()
		}
		s.syncConfig(hBeat)
		hbr.Success = true
	}
	hbr.Term = obtainRaft().term
	return
}

// RequestVote 发起选举，索要选票
func (s *Server) RequestVote(_ context.Context, rv *ReqVote) (rvr *ReqVoteReturn, err error) {
	gnomon.Log().Info("raft", gnomon.Log().Field("receive RequestVote", rv))
	rvr = &ReqVoteReturn{}
	rvr.Term = obtainRaft().term
	if rv.Term < obtainRaft().term {
		gnomon.Log().Info("raft", gnomon.Log().Field("refuse", rv),
			gnomon.Log().Field("termLocal", obtainRaft().term),
			gnomon.Log().Field("termReceive", rv.Term))
		rvr.VoteGranted = false
	} else if rv.Term >= obtainRaft().term {
		rvr.VoteGranted = s.voteFor(rv)
	}
	s.syncNodes(&Node{
		Id:  rv.CandidateId,
		Url: rv.Url,
	})
	return
}

func (s *Server) syncNodes(node *Node) {
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
func (s *Server) syncConfig(hBeat *HBeat) {
	if obtainRaft().persistence.leaderID == hBeat.LeaderId && obtainRaft().persistence.version == hBeat.Version && obtainRaft().term == hBeat.Term {
		obtainRaft().scheduled.refreshLastHeartBeatTime()
		return
	}
	var cs map[string]*config.Config
	if err := yaml.Unmarshal(hBeat.Config, &cs); nil != err {
		gnomon.Log().Error("raft", gnomon.Log().Err(err))
	} else {
		service.RecoverConfig(cs)
		gnomon.Log().Debug("raft", gnomon.Log().Field("syncConfig", "refresh time"))
		obtainRaft().term = hBeat.Term
		obtainRaft().persistence.leaderID = hBeat.LeaderId
		obtainRaft().persistence.version = hBeat.Version
		obtainRaft().persistence.currentTerm = hBeat.Term
		obtainRaft().scheduled.refreshLastHeartBeatTime()
	}
}

// voteFor 要求投票节点任期大于当前任期返回方案
func (s *Server) voteFor(rv *ReqVote) bool {
	if rv.Term == obtainRaft().term {
		if rv.Timestamp < obtainRaft().persistence.votedFor.timestamp {
			s.vote(rv)
			return true
		}
	}
	if rv.Term > obtainRaft().persistence.votedFor.term {
		s.vote(rv)
		return true
	}
	if rv.Term == obtainRaft().persistence.votedFor.term && gnomon.String().IsEmpty(obtainRaft().persistence.votedFor.id) {
		s.vote(rv)
		return true
	}
	gnomon.Log().Info("raft", gnomon.Log().Field("refuse", rv),
		gnomon.Log().Field("termLocal", obtainRaft().term),
		gnomon.Log().Field("termReceive", rv.Term))
	return false
}

func (s *Server) vote(rv *ReqVote) {
	obtainRaft().persistence.votedFor.id = rv.CandidateId
	obtainRaft().persistence.votedFor.term = rv.Term
	obtainRaft().scheduled.refreshLastHeartBeatTime()
	gnomon.Log().Info("raft", gnomon.Log().Field("accept", rv),
		gnomon.Log().Field("termLocal", obtainRaft().term),
		gnomon.Log().Field("termReceive", rv.Term))
}
