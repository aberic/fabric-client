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

func (c *ConfigServer) InitConfig(ctx context.Context, in *pb.ReqInit) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.InitConfig(in)
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

func (c *ConfigServer) InitClient(ctx context.Context, in *pb.ReqClient) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.InitClient(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := InitClient(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) InitClientSelf(ctx context.Context, in *pb.ReqClientSelf) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.InitClientSelf(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := InitClientSelf(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) InitClientCustom(ctx context.Context, in *pb.ReqClientCustom) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.InitClientCustom(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := InitClientCustom(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetPeerForChannel(ctx context.Context, in *pb.ReqChannelPeer) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetPeerForChannel(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetPeerForChannel(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetQueryChannelPolicyForChannel(ctx context.Context, in *pb.ReqChannelPolicyQuery) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetQueryChannelPolicyForChannel(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetQueryChannelPolicyForChannel(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetDiscoveryPolicyForChannel(ctx context.Context, in *pb.ReqChannelPolicyDiscovery) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetDiscoveryPolicyForChannel(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetDiscoveryPolicyForChannel(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetEventServicePolicyForChannel(ctx context.Context, in *pb.ReqChannelPolicyEvent) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetEventServicePolicyForChannel(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetEventServicePolicyForChannel(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetOrdererForOrganizations(ctx context.Context, in *pb.ReqOrganizationsOrder) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetOrdererForOrganizations(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetOrdererForOrganizations(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetOrdererForOrganizationsSelf(ctx context.Context, in *pb.ReqOrganizationsOrderSelf) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetOrdererForOrganizationsSelf(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetOrdererForOrganizationsSelf(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetOrgForOrganizations(ctx context.Context, in *pb.ReqOrganizationsOrg) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetOrgForOrganizations(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetOrgForOrganizations(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetOrgForOrganizationsSelf(ctx context.Context, in *pb.ReqOrganizationsOrgSelf) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetOrgForOrganizationsSelf(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetOrgForOrganizationsSelf(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetOrderer(ctx context.Context, in *pb.ReqOrder) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetOrderer(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetOrderer(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetOrdererSelf(ctx context.Context, in *pb.ReqOrderSelf) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetOrdererSelf(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetOrdererSelf(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetPeer(ctx context.Context, in *pb.ReqPeer) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetPeer(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetPeer(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetPeerSelf(ctx context.Context, in *pb.ReqPeerSelf) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetPeerSelf(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetPeerSelf(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetCertificateAuthority(ctx context.Context, in *pb.ReqCertificateAuthority) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetCertificateAuthority(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetCertificateAuthority(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
}

func (c *ConfigServer) AddOrSetCertificateAuthoritySelf(ctx context.Context, in *pb.ReqCertificateAuthoritySelf) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			return &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetCertificateAuthoritySelf(in)
		},
		func() (i interface{}, e error) {
			pbStr, err := AddOrSetCertificateAuthoritySelf(rafts.LeaderURL(), in)
			return pbStr.(*pb.Result), err
		},
	); nil != err {
		return nil, err
	} else {
		return i.(*pb.Result), nil
	}
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

// InitClient 初始化区块链配置信息
func InitClient(url string, req *pb.ReqClient) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.InitClient(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// InitClientSelf 初始化区块链配置信息
func InitClientSelf(url string, req *pb.ReqClientSelf) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.InitClientSelf(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// InitClientCustom 初始化区块链更自定义配置信息
func InitClientCustom(url string, req *pb.ReqClientCustom) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.InitClientCustom(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetPeerForChannel 为通道新增或覆盖节点信息
func AddOrSetPeerForChannel(url string, req *pb.ReqChannelPeer) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetPeerForChannel(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetQueryChannelPolicyForChannel 为通道新增或覆盖检索节点策略
func AddOrSetQueryChannelPolicyForChannel(url string, req *pb.ReqChannelPolicyQuery) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetQueryChannelPolicyForChannel(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetDiscoveryPolicyForChannel 为通道新增或覆盖发现节点策略
func AddOrSetDiscoveryPolicyForChannel(url string, req *pb.ReqChannelPolicyDiscovery) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetDiscoveryPolicyForChannel(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetEventServicePolicyForChannel 为通道新增或覆盖事件策略
func AddOrSetEventServicePolicyForChannel(url string, req *pb.ReqChannelPolicyEvent) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetEventServicePolicyForChannel(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetOrdererForOrganizations 新增或覆盖组织内排序机构信息
func AddOrSetOrdererForOrganizations(url string, req *pb.ReqOrganizationsOrder) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetOrdererForOrganizations(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetOrdererForOrganizationsSelf 新增或覆盖组织内排序机构信息
func AddOrSetOrdererForOrganizationsSelf(url string, req *pb.ReqOrganizationsOrderSelf) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetOrdererForOrganizationsSelf(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetOrgForOrganizations 新增或覆盖组织内组织机构信息
func AddOrSetOrgForOrganizations(url string, req *pb.ReqOrganizationsOrg) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetOrgForOrganizations(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetOrgForOrganizationsSelf 新增或覆盖组织内组织机构信息
func AddOrSetOrgForOrganizationsSelf(url string, req *pb.ReqOrganizationsOrgSelf) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetOrgForOrganizationsSelf(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetOrderer 新增或覆盖排序服务信息
func AddOrSetOrderer(url string, req *pb.ReqOrder) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetOrderer(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetOrderer 新增或覆盖排序服务信息
func AddOrSetOrdererSelf(url string, req *pb.ReqOrderSelf) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetOrdererSelf(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetPeer 新增或覆盖节点服务信息
func AddOrSetPeer(url string, req *pb.ReqPeer) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetPeer(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetPeerSelf 新增或覆盖节点服务信息
func AddOrSetPeerSelf(url string, req *pb.ReqPeerSelf) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetPeerSelf(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetCertificateAuthority 新增或覆盖CA服务信息
func AddOrSetCertificateAuthority(url string, req *pb.ReqCertificateAuthority) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetCertificateAuthority(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// AddOrSetCertificateAuthority 新增或覆盖CA服务信息
func AddOrSetCertificateAuthoritySelf(url string, req *pb.ReqCertificateAuthoritySelf) (interface{}, error) {
	return utils.RPC(url, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *pb.Result
			err    error
		)
		// 创建grpc客户端
		c := pb.NewLedgerConfigClient(conn)
		// 客户端向grpc服务端发起请求
		if result, err = c.AddOrSetCertificateAuthoritySelf(context.Background(), req); nil != err {
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
