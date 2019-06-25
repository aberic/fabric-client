package raft

import (
	"context"
	pb "github.com/ennoo/fabric-client/grpc/proto/raft"
	"github.com/ennoo/fabric-client/grpc/proto/utils"
	"google.golang.org/grpc"
)

// HeartBeat 发送心跳
func HeartBeat(url string, req *pb.Beat) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Beat
			err    error
		)
		// 创建grpc客户端
		c := pb.NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if result, err = c.HeartBeat(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// RequestVote 发起选举，索要选票
func RequestVote(url string, req *pb.ReqElection) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Resp
			err    error
		)
		// 创建grpc客户端
		c := pb.NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if result, err = c.RequestVote(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// FollowMe 成为Leader并要求被跟随
func FollowMe(url string, req *pb.ReqFollow) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
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
}

// LeaderMe 请求Leader将自身加入follows
func LeaderMe(url string, req *pb.ReqLeader) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
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
}

// SyncNode 同步节点信息
func SyncNode(url string, req *pb.NodeMap) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.NodeMap
			err    error
		)
		// 创建grpc客户端
		c := pb.NewRaftClient(conn)
		//客户端向grpc服务端发起请求
		if result, err = c.SyncNode(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}
