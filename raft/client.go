package raft

import (
	"context"
	pb "github.com/ennoo/fabric-client/grpc/proto/raft"
	"github.com/ennoo/fabric-client/grpc/proto/utils"
	"github.com/ennoo/rivet/utils/log"
	"google.golang.org/grpc"
	"strings"
	"time"
)

type HB struct {
	URL string
	Req *pb.ReqElection
}

type RV struct {
	URL    string
	Req    *pb.ReqElection
	Target *pb.ReqElection
}

type SN struct {
	URL string
	Req *pb.NodeMap
}

// heartBeat 发送心跳
func heartBeat(i interface{}) {
	hb := i.(*HB)
	_, _ = utils.RPC(hb.URL, func(conn *grpc.ClientConn) (interface{}, error) {
		// 创建grpc客户端
		c := pb.NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if _, err := c.HeartBeat(context.Background(), &pb.Beat{}); nil != err {
			return nil, err
		}
		return nil, nil
	})
	//if nil != err {
	//	if time.Now().UnixNano()/1e6-Nodes[hb.Req.Node.Id].LastActive > 1000 {
	//		delete(Nodes, hb.Req.Node.Id)
	//	}
	//} else {
	//	Nodes[hb.Req.Node.Id].LastActive = time.Now().UnixNano() / 1e6
	//}
}

// RequestVote 发起选举，索要选票
func RequestVote(i interface{}) {
	log.Self.Info("raft", log.Int32("Term", Term), log.String("RequestVote", "发起选举，索要选票"))
	rv := i.(*RV)
	pbI, err := utils.RPC(rv.URL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Resp
			err    error
		)
		// 创建grpc客户端
		c := pb.NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if result, err = c.RequestVote(context.Background(), rv.Req); nil != err {
			return nil, err
		}
		return result, nil
	})
	if nil != err {
		return
	}
	//if nil != err {
	//	if Nodes[rv.Target.Node.Id].LastActive == 0 { // 首次请求该节点最后活跃时间初始0，默认已连接成功过，并赋值一次
	//		Nodes[rv.Target.Node.Id].LastActive = time.Now().UnixNano() / 1e6
	//	}
	//	if time.Now().UnixNano()/1e6-Nodes[rv.Target.Node.Id].LastActive > 1000 {
	//		log.Self.Info("raft",
	//			log.Int32("Term", Term),
	//			log.String("BrokerID", rv.Target.Node.Id),
	//			log.String("RequestVote", "节点三次未连接成功，移除该节点"),
	//		)
	//		delete(Nodes, rv.Target.Node.Id)
	//	}
	//	return
	//}
	//Nodes[rv.Target.Node.Id].LastActive = time.Now().UnixNano() / 1e6
	RefreshTimeOut()
	resp := pbI.(*pb.Resp)
	switch resp.Type {
	case pb.Type_OK:
		log.Self.Info("raft", log.Int32("Term", Term), log.String("RequestVote", "获得选票"))
		voteCount += 1
		var aliveNodeCount int32
		aliveNodeCount = 0
		for _, node := range Nodes {
			if time.Now().UnixNano()/1e6-node.LastActive < 1000 {
				aliveNodeCount++
			}
		}
		if voteCount+1 >= aliveNodeCount/2 {
			voteCount = 0
			log.Self.Info("raft", log.Int32("Term", Term), log.String("RequestVote", "became leader, now status LEADER"))
			Leader = &pb.Leader{BrokerID: ID, Term: Term}
			Nodes[ID].Status = pb.Status_LEADER
			// 总得票数超过总节点数一半
			for _, node := range Nodes {
				go followMe(strings.Join([]string{node.Addr, ":", node.Rpc}, ""), &pb.ReqFollow{
					Node: Nodes[ID],
					Term: Term,
				})
			}
		}
	case pb.Type_TERM:
		Nodes[ID].Status = pb.Status_FOLLOWER
		Term = resp.GetElection().Term
	}
}

// followMe 成为Leader并要求被跟随
func followMe(url string, req *pb.ReqFollow) {
	log.Self.Info("raft", log.Int32("Term", Term), log.String("followMe", "成为Leader并要求被跟随"))
	pbI, err := utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Resp
			err    error
		)
		// 创建grpc客户端
		c := pb.NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if result, err = c.FollowMe(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
	if nil != err {
		return
	}
	resp := pbI.(*pb.Resp)
	switch resp.Type {
	case pb.Type_TERM:
		RefreshTimeOut()
		if Term >= resp.GetElection().Term {
			return
		}
		Term = resp.GetElection().Term
		// 临时将状态切回FOLLOW，等待可能存在的Leader节点介入信息
		Nodes[ID].Status = pb.Status_FOLLOWER
	}
}

// leaderMe 请求Leader将自身加入follows
func leaderMe(url string, req *pb.Node) {
	log.Self.Info("raft", log.Int32("Term", Term), log.String("leaderMe", "请求Leader将自身加入follows"))
	pbI, err := utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Resp
			err    error
		)
		// 创建grpc客户端
		c := pb.NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if result, err = c.LeaderMe(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
	if nil != err {
		return
	}
	resp := pbI.(*pb.Resp)
	switch resp.Type {
	case pb.Type_FOLLOW_ME:
		RefreshTimeOut()
		if Term >= resp.GetElection().Term {
			return
		}
		Term = resp.GetElection().Term
		log.Self.Info("raft", log.Int32("Term", Term), log.String("leaderMe", "find leader, now status FOLLOWER"))
		// 刷新计时，并设置对方节点成为了 Leader
		RefreshTimeOut()
		Leader = &pb.Leader{BrokerID: resp.GetElection().Node.Id, Term: resp.GetElection().Term}
		Nodes[ID].Status = pb.Status_FOLLOWER
	default:
		// 刷新计时，并等待可能存在的Leader节点介入信息
		RefreshTimeOut()
		//leader := resp.GetLeader().Leader
		//// 先判断本地是否有该节点，如果没有则新增
		//if Nodes[leader.Id] != nil {
		//	Nodes[leader.Id] = leader
		//}
		//leaderMe(strings.Join([]string{leader.Addr, ":", leader.Rpc}, ""), req)
	}
}

// syncNode 同步节点信息
func syncNode(i interface{}) {
	sn := i.(*SN)
	pbI, err := utils.RPC(sn.URL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.NodeMap
			err    error
		)
		// 创建grpc客户端
		c := pb.NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if result, err = c.SyncNode(context.Background(), sn.Req); nil != err {
			return nil, err
		}
		return result, nil
	})
	if nil != err {
		return
	}
	nodeMap := pbI.(*pb.NodeMap)
	for _, node := range nodeMap.Nodes {
		haveOne := false
		for _, localNode := range Nodes {
			if node.Id == localNode.Id {
				haveOne = true
				break
			}
		}
		if !haveOne {
			Nodes[node.Id] = node
		}
	}
}
