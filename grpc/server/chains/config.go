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
	configer "github.com/ennoo/fabric-client/config"
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

func (c *ConfigServer) InitClient(ctx context.Context, in *pb.ReqClient) (*pb.Result, error) {
	if i, err := c.proxy(
		true,
		func() (i interface{}, e error) {
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.InitClient(&service.Client{
				ConfigID:     in.ConfigID,
				TlS:          in.Tls,
				Organization: in.Organization,
				Level:        in.Level,
				CryptoConfig: in.CryptoConfig,
				KeyPath:      in.KeyPath,
				CertPath:     in.CertPath,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.InitClientSelf(&service.ClientSelf{
				ConfigID:     in.ConfigID,
				TlS:          in.Tls,
				LeagueName:   in.LeagueName,
				UserName:     in.UserName,
				Organization: in.Organization,
				Level:        in.Level,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.InitClientCustom(&service.ClientCustom{
				ConfigID: in.ConfigID,
				Client: &service.Client{
					ConfigID:     in.ConfigID,
					TlS:          in.Client.Tls,
					Organization: in.Client.Organization,
					Level:        in.Client.Level,
					CryptoConfig: in.Client.CryptoConfig,
					KeyPath:      in.Client.KeyPath,
					CertPath:     in.Client.CertPath,
				},
				Peer: &configer.ClientPeer{
					Timeout: &configer.ClientPeerTimeout{
						Connection: in.Peer.Timeout.Connection,
						Response:   in.Peer.Timeout.Response,
						Discovery: &configer.ClientPeerTimeoutDiscovery{
							GreyListExpiry: in.Peer.Timeout.Discovery.GreyListExpiry,
						},
					},
				},
				EventService: &configer.ClientEventService{
					Timeout: &configer.ClientEventServiceTimeout{
						RegistrationResponse: in.EventService.Timeout.RegistrationResponse,
					},
				},
				Order: &configer.ClientOrder{
					Timeout: &configer.ClientOrderTimeout{
						Connection: in.Order.Timeout.Connection,
						Response:   in.Order.Timeout.Response,
					},
				},
				Global: &configer.ClientGlobal{
					Timeout: &configer.ClientGlobalTimeout{
						Query:   in.Global.Timeout.Query,
						Execute: in.Global.Timeout.Execute,
						Resmgmt: in.Global.Timeout.Resmgmt,
					},
					Cache: &configer.ClientGlobalCache{
						ConnectionIdle:    in.Global.Cache.ConnectionIdle,
						EventServiceIdle:  in.Global.Cache.EventServiceIdle,
						ChannelMembership: in.Global.Cache.ChannelMembership,
						ChannelConfig:     in.Global.Cache.ChannelConfig,
						Discovery:         in.Global.Cache.Discovery,
						Selection:         in.Global.Cache.Selection,
					},
				},
				BCCSP: &configer.ClientBCCSP{
					Security: &configer.ClientBCCSPSecurity{
						Enabled: in.BCCSP.Security.Enabled,
						Default: &configer.ClientBCCSPSecurityDefault{
							Provider: in.BCCSP.Security.Default.Provider,
						},
						HashAlgorithm: in.BCCSP.Security.HashAlgorithm,
						SoftVerify:    in.BCCSP.Security.SoftVerify,
						Level:         in.BCCSP.Security.Level,
					},
				},
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetPeerForChannel(&service.ChannelPeer{
				ConfigID:       in.ConfigID,
				ChannelName:    in.ChannelName,
				PeerName:       in.PeerName,
				EndorsingPeer:  in.EndorsingPeer,
				ChainCodeQuery: in.ChainCodeQuery,
				LedgerQuery:    in.LedgerQuery,
				EventSource:    in.EventSource,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetQueryChannelPolicyForChannel(&service.ChannelPolicyQuery{
				ConfigID:       in.ConfigID,
				ChannelName:    in.ChannelName,
				InitialBackOff: in.InitialBackOff,
				MaxBackOff:     in.MaxBackOff,
				MaxTargets:     in.MaxTargets,
				MinResponses:   in.MinResponses,
				Attempts:       in.Attempts,
				BackOffFactor:  in.BackOffFactor,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetDiscoveryPolicyForChannel(&service.ChannelPolicyDiscovery{
				ConfigID:       in.ConfigID,
				ChannelName:    in.ChannelName,
				InitialBackOff: in.InitialBackOff,
				MaxBackOff:     in.MaxBackOff,
				MaxTargets:     in.MaxTargets,
				Attempts:       in.Attempts,
				BackOffFactor:  in.BackOffFactor,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetEventServicePolicyForChannel(&service.ChannelPolicyEvent{
				ConfigID:                         in.ConfigID,
				ChannelName:                      in.ChannelName,
				ReconnectBlockHeightLagThreshold: in.ReconnectBlockHeightLagThreshold,
				ResolverStrategy:                 in.ResolverStrategy,
				BlockHeightLagThreshold:          in.BlockHeightLagThreshold,
				Balance:                          in.Balance,
				PeerMonitorPeriod:                in.PeerMonitorPeriod,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetOrdererForOrganizations(&service.OrganizationsOrder{
				ConfigID:   in.ConfigID,
				MspID:      in.MspID,
				CryptoPath: in.CryptoPath,
				Users:      in.Users,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetOrdererForOrganizationsSelf(&service.OrganizationsOrderSelf{
				ConfigID:   in.ConfigID,
				LeagueName: in.LeagueName,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetOrgForOrganizations(&service.OrganizationsOrg{
				ConfigID:               in.ConfigID,
				MspID:                  in.MspID,
				CryptoPath:             in.CryptoPath,
				OrgName:                in.OrgName,
				Users:                  in.Users,
				Peers:                  in.Peers,
				CertificateAuthorities: in.CertificateAuthorities,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetOrgForOrganizationsSelf(&service.OrganizationsOrgSelf{
				ConfigID:               in.ConfigID,
				LeagueName:             in.LeagueName,
				Peers:                  in.Peers,
				CertificateAuthorities: in.CertificateAuthorities,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetOrderer(&service.Order{
				ConfigID:              in.ConfigID,
				OrderName:             in.OrderName,
				URL:                   in.Url,
				TLSCACerts:            in.TlsCACerts,
				SSLTargetNameOverride: in.SslTargetNameOverride,
				KeepAliveTime:         in.KeepAliveTime,
				KeepAliveTimeout:      in.KeepAliveTimeout,
				KeepAlivePermit:       in.KeepAlivePermit,
				FailFast:              in.FailFast,
				AllowInsecure:         in.AllowInsecure,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetOrdererSelf(&service.OrderSelf{
				ConfigID:         in.ConfigID,
				OrderName:        in.OrderName,
				URL:              in.Url,
				LeagueName:       in.LeagueName,
				KeepAliveTime:    in.KeepAliveTime,
				KeepAliveTimeout: in.KeepAliveTimeout,
				KeepAlivePermit:  in.KeepAlivePermit,
				FailFast:         in.FailFast,
				AllowInsecure:    in.AllowInsecure,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetPeer(&service.Peer{
				ConfigID:              in.ConfigID,
				PeerName:              in.PeerName,
				URL:                   in.Url,
				EventUrl:              in.EventUrl,
				TLSCACerts:            in.TlsCACerts,
				SSLTargetNameOverride: in.SslTargetNameOverride,
				KeepAliveTime:         in.KeepAliveTime,
				KeepAliveTimeout:      in.KeepAliveTimeout,
				KeepAlivePermit:       in.KeepAlivePermit,
				FailFast:              in.FailFast,
				AllowInsecure:         in.AllowInsecure,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetPeerSelf(&service.PeerSelf{
				ConfigID:         in.ConfigID,
				PeerName:         in.PeerName,
				URL:              in.Url,
				EventUrl:         in.EventUrl,
				LeagueName:       in.LeagueName,
				KeepAliveTime:    in.KeepAliveTime,
				KeepAliveTimeout: in.KeepAliveTimeout,
				KeepAlivePermit:  in.KeepAlivePermit,
				FailFast:         in.FailFast,
				AllowInsecure:    in.AllowInsecure,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetCertificateAuthority(&service.CertificateAuthority{
				ConfigID:                in.ConfigID,
				CertName:                in.CertName,
				URL:                     in.Url,
				TLSCACertPath:           in.TlsCACertPath,
				TLSCACertClientKeyPath:  in.TlsCACertClientKeyPath,
				TLSCACertClientCertPath: in.TlsCACertClientCertPath,
				CAName:                  in.CaName,
				EnrollId:                in.EnrollId,
				EnrollSecret:            in.EnrollSecret,
			})
			return pbStr, err
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
			pbStr, err := &pb.Result{Code: pb.Code_Success, Data: "success"}, service.AddOrSetCertificateAuthoritySelf(&service.CertificateAuthoritySelf{
				ConfigID:     in.ConfigID,
				CertName:     in.CertName,
				URL:          in.Url,
				LeagueName:   in.LeagueName,
				CAName:       in.CaName,
				EnrollId:     in.EnrollId,
				EnrollSecret: in.EnrollSecret,
			})
			return pbStr, err
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
