package raft

import (
	"context"
	pb "github.com/ennoo/fabric-client/grpc/proto/raft"
	"github.com/ennoo/fabric-client/grpc/proto/utils"
	"google.golang.org/grpc"
)

type HB struct {
	URL string
	Req *pb.Beat
}

type RV struct {
	URL string
	Req *pb.ReqElection
}

type FM struct {
	URL string
	Req *pb.ReqFollow
}

type LM struct {
	URL string
	Req *pb.ReqLeader
}

type SN struct {
	URL string
	Req *pb.NodeMap
}

func hb(hb *HB) {

}

func rv(rv *RV) {

}

func fm(fm *FM) {

}

func lm(lm *LM) {

}

func sn(sn *SN) {

}

// heartBeat 发送心跳
func heartBeat(hb *HB) (interface{}, error) {
	return utils.RPC(hb.URL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Beat
			err    error
		)
		// 创建grpc客户端
		c := pb.NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if result, err = c.HeartBeat(context.Background(), hb.Req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// requestVote 发起选举，索要选票
func requestVote(rv *RV) (interface{}, error) {
	return utils.RPC(rv.URL, func(conn *grpc.ClientConn) (interface{}, error) {
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
}

// followMe 成为Leader并要求被跟随
func followMe(fm *FM) (interface{}, error) {
	return utils.RPC(fm.URL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Resp
			err    error
		)
		// 创建grpc客户端
		c := pb.NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if result, err = c.FollowMe(context.Background(), fm.Req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// leaderMe 请求Leader将自身加入follows
func leaderMe(lm *LM) (interface{}, error) {
	return utils.RPC(lm.URL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Resp
			err    error
		)
		// 创建grpc客户端
		c := pb.NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if result, err = c.LeaderMe(context.Background(), lm.Req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// syncNode 同步节点信息
func syncNode(sn *SN) (interface{}, error) {
	return utils.RPC(sn.URL, func(conn *grpc.ClientConn) (interface{}, error) {
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
}
