package utils

import (
	"google.golang.org/grpc"
)

// RPC 通过rpc进行通信 protoc --go_out=plugins=grpc:. grpc/proto/*.proto
func RPC(url string, business func(conn *grpc.ClientConn) (interface{}, error)) (interface{}, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	// 创建一个grpc连接器
	if conn, err = grpc.Dial(url, grpc.WithInsecure()); nil != err {
		return nil, err
	}
	// 请求完毕后关闭连接
	defer func() { _ = conn.Close() }()
	return business(conn)
}
