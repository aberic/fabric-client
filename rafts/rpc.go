/*
 * Copyright (c) 2019.. Aberic - All Rights Reserved.
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
	"google.golang.org/grpc"
)

// rpc 通过rpc进行通信
func rpc(url string, business func(conn *grpc.ClientConn) (interface{}, error)) (interface{}, error) {
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
