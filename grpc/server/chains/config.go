/*
 * Copyright (c) 2019.. ENNOO - All Rights Reserved.
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

package chains

import (
	"errors"
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	"github.com/ennoo/fabric-client/grpc/proto/utils"
	"github.com/ennoo/fabric-client/rafts"
	"github.com/ennoo/fabric-client/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

type ConfigServer struct {
}

func (c *ConfigServer) ListConfig(ctx context.Context, in *pb.ReqConfigList) (*pb.ResultConfigList, error) {
	configList := &pb.ResultConfigList{}
	configList.Code = pb.Code_Success
	configList.ConfigIDs = []string{}
	for id := range service.Configs {
		configList.ConfigIDs = append(configList.ConfigIDs, id)
	}
	return configList, nil
}

func (c *ConfigServer) GetConfig(ctx context.Context, in *pb.ReqConfig) (*pb.ResultConfig, error) {
	config := service.Get(in.ConfigID)
	if nil == config {
		errStr := "config is nil"
		return &pb.ResultConfig{Code: pb.Code_Fail, ErrMsg: errStr}, errors.New(errStr)
	} else {
		//confData, err := yaml.Marshal(&config)
		//if err != nil {
		//	log.Self.Debug("client", log.Error(err))
		//}
		//fmt.Printf("--- dump:\n%s\n\n", string(confData))
		return &pb.ResultConfig{Code: pb.Code_Success, Config: config.GetPBConfig()}, nil
	}
}

func (c *ConfigServer) RecoverConfig(ctx context.Context, in *pb.ReqConfigRecover) (*pb.Result, error) {
	service.Recover(in.ConfigIDs)
	return &pb.Result{Code: pb.Code_Success}, nil
}

func (c *ConfigServer) InitConfig(ctx context.Context, in *pb.ReqInit) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			service.InitConfig(in)
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, nil
		},
		func() (i interface{}, e error) {
			if pbStr, err := InitConfig(rafts.LeaderURL(), in); nil != err {
				return &pb.Result{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
			} else {
				return pbStr.(*pb.Result), nil
			}
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

// ListConfig 获取区块链配置信息ID集合
func ListConfig(url string, req *pb.ReqConfigList) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.ResultConfigList
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.ListConfig(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// GetConfig 获取区块链配置信息
func GetConfig(url string, req *pb.ReqConfig) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.ResultConfig
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.GetConfig(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// InitConfig 初始化区块链配置信息
func InitConfig(url string, req *pb.ReqInit) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.InitConfig(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

func (c *ConfigServer) proxy(sleep bool, exec func() (interface{}, error), trans func() (interface{}, error)) (interface{}, error) {
	switch rafts.Character() {
	case rafts.RoleLeader: // 自身即为 Leader 节点
		i, err := exec()
		if nil == err {
			rafts.VersionAdd()
		}
		return i, err
	case rafts.RoleCandidate: // 等待选举结果，如果超时则返回
		if sleep {
			time.Sleep(1000 * time.Millisecond)
			return c.proxy(false, exec, trans)
		} else {
			return nil, errors.New("leader is nil")
		}
	case rafts.RoleFollower: // 将该请求转发给Leader节点处理
		return trans()
	}
	return nil, errors.New("unknown err")
}
