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
	pb "github.com/ennoo/fabric-client/grpc/proto/raft"
	"github.com/ennoo/fabric-client/raft"
	"github.com/ennoo/rivet"
	"github.com/ennoo/rivet/utils/log"
	str "github.com/ennoo/rivet/utils/string"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"time"
)

type RaftServer struct {
}

func (r *RaftServer) HeartBeat(ctx context.Context, in *pb.Beat) (*pb.Beat, error) {
	brokerID := string(in.Beat)
	// 如果心跳发送发为当前所处 Leader 节点，则刷新计时
	if brokerID == raft.Leader.BrokerID {
		raft.RefreshTimeOut()
		if !raft.Sync {
			leader := raft.Nodes[raft.Leader.BrokerID]
			uri := strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, "")
			log.Self.Debug("syncConfig", log.String("uri", uri))
			_, err := rivet.Request().RestJSON(http.MethodGet, uri, "config/sync/ask", nil)
			if nil != err {
				log.Self.Debug("sync ask", log.Error(err))
			} else {
				raft.Sync = true
			}
		}
		return &pb.Beat{}, nil
	}
	return nil, nil
	//return nil, errors.New("you are not leader for me")
}

func (r *RaftServer) RequestVote(ctx context.Context, in *pb.ReqElection) (*pb.Resp, error) {
	if in.Term <= raft.Term {
		log.Self.Info("raft", log.Int32("Term", raft.Term), log.String("RequestVote", "拒绝投票给对方"))
		if raft.Leader.BrokerID == raft.ID { // 如果自身为Leader，则要求对方Follow
			raft.FollowMe(strings.Join([]string{in.Node.Addr, ":", in.Node.Rpc}, ""), &pb.ReqFollow{
				Node: raft.Nodes[raft.ID],
				Term: raft.Term,
			})
		}
		return &pb.Resp{
			Type: pb.Type_TERM,
			Result: &pb.Resp_Election{
				Election: &pb.ReqElection{
					Term: raft.Term,
				},
			},
		}, nil
	}
	// 更新自身区间
	raft.Term = in.Term
	// 如果索要区间大于当前所处区间，则同投票并重新计时
	raft.RefreshTimeOut()
	log.Self.Info("raft", log.Int32("Term", raft.Term), log.String("RequestVote", "同意投票给对方"))
	// 临时将状态切回FOLLOW，等待可能存在的Leader节点介入信息
	raft.Nodes[raft.ID].Status = pb.Status_FOLLOWER
	raft.RefreshTimeOut()
	return &pb.Resp{Type: pb.Type_OK}, nil
}

func (r *RaftServer) FollowMe(ctx context.Context, in *pb.ReqFollow) (*pb.Resp, error) {
	raft.RefreshTimeOut()
	// 检查该节点是否在集群中，如果不在，则新增
	if nil == raft.Nodes[in.Node.Id] {
		raft.Nodes[in.Node.Id] = in.Node
	}
	if in.Term < raft.Term {
		// 告知对方当前最新Term
		return &pb.Resp{
			Type: pb.Type_TERM,
			Result: &pb.Resp_Election{
				Election: &pb.ReqElection{
					Term: raft.Term,
				},
			},
		}, nil
	}
	log.Self.Info("raft", log.Int32("Term", raft.Term), log.String("FollowMe", "find leader, now status FOLLOWER"))
	// 更新自身区间
	raft.Term = in.Term
	// 刷新计时，并设置对方节点成为了 Leader
	raft.RefreshTimeOut()
	raft.Leader = &pb.Leader{BrokerID: in.Node.Id, Term: in.Term}
	raft.Nodes[in.Node.Id] = in.Node
	raft.Nodes[raft.ID].Status = pb.Status_FOLLOWER
	return &pb.Resp{Type: pb.Type_OK}, nil
}

func (r *RaftServer) LeaderMe(ctx context.Context, in *pb.Node) (*pb.Resp, error) {
	log.Self.Info("raft", log.Int32("Term", raft.Term), log.String("LeaderMe", "请求被跟随"))
	// 检查该节点是否在集群中，如果不在，则新增
	if nil == raft.Nodes[in.Id] {
		raft.Nodes[in.Id] = in
	}
	return termLittle(in)
}

func (r *RaftServer) SyncNode(ctx context.Context, in *pb.NodeMap) (*pb.NodeMap, error) {
	log.Self.Debug("raft", log.Int32("Term", raft.Term), log.String("SyncNode", "接收同步节点信息"))
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
	// 同步节点信息当做一次心跳处理
	raft.RefreshTimeOut()
	return &pb.NodeMap{Nodes: raft.Nodes}, nil
}

func (r *RaftServer) NodeList(ctx context.Context, in *pb.Beat) (*pb.Nodes, error) {
	log.Self.Debug("raft", log.Int32("Term", raft.Term), log.String("NodeList", "获取节点集合"))
	var nodes []*pb.Node
	if str.IsEmpty(raft.ID) || (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		for _, node := range raft.Nodes {
			if time.Now().UnixNano()/1e6-node.LastActive < 1500 {
				nodes = append(nodes, node)
			}
		}
		nodes = append(nodes, raft.Nodes[raft.ID])
		return &pb.Nodes{
			Nodes: nodes,
		}, nil
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		return &pb.Nodes{
			Node: &pb.Node{
				Addr: leader.Addr,
				Rpc:  leader.Rpc,
			},
		}, nil
	}
}

func termLittle(req *pb.Node) (*pb.Resp, error) {
	// 如果索要区间小于当前所处区间
	if raft.Nodes[raft.ID].Status == pb.Status_LEADER {
		// 如果自身为 Leader 状态，则告知跟随
		log.Self.Debug("raft", log.Int32("Term", raft.Term), log.String("termLittle", "自身为 Leader 状态，告知跟随"))
		return &pb.Resp{Type: pb.Type_FOLLOW_ME,
			Result: &pb.Resp_Election{Election: &pb.ReqElection{Node: raft.Nodes[raft.ID], Term: raft.Term}}}, nil
	} else if raft.Nodes[raft.ID].Status == pb.Status_CANDIDATE {
		// 如果自身为 CANDIDATE 状态，则告知投票
		log.Self.Debug("raft", log.Int32("Term", raft.Term), log.String("termLittle", "自身为 CANDIDATE 状态，告知投票"))
		raft.RequestVote(&raft.RV{
			URL: strings.Join([]string{req.Addr, ":", req.Rpc}, ""),
			Req: &pb.ReqElection{
				Node: raft.Nodes[raft.ID],
				Term: raft.Term,
			},
		})
		return &pb.Resp{Type: pb.Type_VOTE_ME,
			Result: &pb.Resp_Election{Election: &pb.ReqElection{Node: raft.Nodes[raft.ID], Term: raft.Term}}}, nil
	}
	// 如果自身为 FOLLOW 状态，则告知当前Leader节点
	log.Self.Debug("raft", log.Int32("Term", raft.Term), log.String("termLittle", "自身为 FOLLOW 状态，告知当前Leader节点"))
	return &pb.Resp{Type: pb.Type_LEADER_NODE,
		Result: &pb.Resp_Leader{Leader: &pb.LeaderNode{
			Leader: raft.Nodes[raft.Leader.BrokerID],
			Term:   raft.Term,
		}}}, nil
}
