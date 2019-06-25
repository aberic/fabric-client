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
	"errors"
	pb "github.com/ennoo/fabric-client/grpc/proto/raft"
	"github.com/ennoo/fabric-client/raft"
	"golang.org/x/net/context"
)

type RaftServer struct {
}

func (r *RaftServer) HeartBeat(ctx context.Context, in *pb.Beat) (*pb.Beat, error) {
	brokerID := string(in.Beat)
	// 如果心跳发送发为当前所处 Leader 节点，则刷新计时
	if brokerID == raft.Leader.BrokerID {
		raft.RefreshTimeOut()
		return &pb.Beat{}, nil
	}
	return nil, errors.New("you are not leader for me")
}

func (r *RaftServer) RequestVote(ctx context.Context, in *pb.ReqElection) (*pb.Resp, error) {
	if in.Term <= raft.Term {
		return termLittle()
	}
	// 如果索要区间大于当前所处区间，则同投票并重新计时
	raft.RefreshTimeOut()
	return &pb.Resp{Type: pb.Type_OK}, nil
}

func (r *RaftServer) FollowMe(ctx context.Context, in *pb.ReqFollow) (*pb.Resp, error) {
	if in.Term <= raft.Term {
		return termLittle()
	}
	// 设置对方节点成为了 Leader，并刷新计时
	raft.Leader = &pb.Leader{BrokerID: in.Node.Id, Term: in.Term}
	raft.Nodes[in.Node.Id] = in.Node
	raft.RefreshTimeOut()
	return &pb.Resp{Type: pb.Type_OK}, nil
}

func (r *RaftServer) LeaderMe(ctx context.Context, in *pb.ReqLeader) (*pb.Resp, error) {
	return termLittle()
}

func (r *RaftServer) SyncNode(ctx context.Context, in *pb.NodeMap) (*pb.NodeMap, error) {
	go syncNode(in)
	// 同步节点信息当做一次心跳处理
	raft.RefreshTimeOut()
	return &pb.NodeMap{Nodes: raft.Nodes}, nil
}

func syncNode(in *pb.NodeMap) {
	for _, node := range in.Nodes {
		haveOne := false
		for _, localNode := range raft.Nodes {
			if node.Id == localNode.Id {
				haveOne = true
				break
			}
		}
		if !haveOne {
			raft.Nodes[node.Id] = node
		}
	}
}

func termLittle() (*pb.Resp, error) {
	// 如果索要区间小于当前所处区间
	if raft.Nodes[raft.ID].Status == pb.Status_LEADER {
		// 如果自身为 Leader 状态，则告知跟随
		return &pb.Resp{Type: pb.Type_FOLLOW_ME,
			Result: &pb.Resp_Election{Election: &pb.ReqElection{BrokerID: raft.ID, Term: raft.Term}}}, nil
	} else if raft.Nodes[raft.ID].Status == pb.Status_CANDIDATE {
		// 如果自身为 CANDIDATE 状态，则告知投票
		return &pb.Resp{Type: pb.Type_VOTE_ME,
			Result: &pb.Resp_Election{Election: &pb.ReqElection{BrokerID: raft.ID, Term: raft.Term}}}, nil
	}
	// 如果自身为 FOLLOW 状态，则告知当前Leader节点
	return &pb.Resp{Type: pb.Type_LEADER_NODE,
		Result: &pb.Resp_Leader{Leader: &pb.LeaderNode{
			Leader: raft.Nodes[raft.Leader.BrokerID],
			Term:   raft.Term,
		}}}, nil
}
