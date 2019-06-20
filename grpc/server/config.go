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

package server

import (
	"errors"
	configer "github.com/ennoo/fabric-go-client/config"
	pb "github.com/ennoo/fabric-go-client/grpc/proto"
	"github.com/ennoo/fabric-go-client/service"
	"golang.org/x/net/context"
)

type ConfigServer struct {
}

func (c *ConfigServer) GetConfig(ctx context.Context, in *pb.String) (*pb.Config, error) {
	config := service.Get(in.Data)
	if nil == config {
		return nil, errors.New("config is nil")
	} else {
		return config.GetPBConfig(), nil
	}
}

func (c *ConfigServer) InitClient(ctx context.Context, in *pb.ReqClient) (*pb.String, error) {
	config := service.Get(in.ConfigID)
	if nil == config {
		return nil, errors.New("config is nil")
	} else {
		return &pb.String{Data: "success"}, service.InitClient(&service.Client{
			ConfigID:     in.ConfigID,
			TlS:          in.Tls,
			Organization: in.Organization,
			Level:        in.Level,
			CryptoConfig: in.CryptoConfig,
			KeyPath:      in.KeyPath,
			CertPath:     in.CertPath,
		})
	}
}

func (c *ConfigServer) InitClientCustom(ctx context.Context, in *pb.ReqClientCustom) (*pb.String, error) {
	config := service.Get(in.ConfigID)
	if nil == config {
		return nil, errors.New("config is nil")
	} else {
		return &pb.String{Data: "success"}, service.InitClientCustom(&service.ClientCustom{
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
	}
}

func (c *ConfigServer) AddOrSetPeerForChannel(ctx context.Context, in *pb.ReqChannelPeer) (*pb.String, error) {
	config := service.Get(in.ConfigID)
	if nil == config {
		return nil, errors.New("config is nil")
	} else {
		return &pb.String{Data: "success"}, service.AddOrSetPeerForChannel(&service.ChannelPeer{
			ConfigID:       in.ConfigID,
			ChannelName:    in.ChannelName,
			PeerName:       in.PeerName,
			EndorsingPeer:  in.EndorsingPeer,
			ChainCodeQuery: in.ChainCodeQuery,
			LedgerQuery:    in.LedgerQuery,
			EventSource:    in.EventSource,
		})
	}
}

func (c *ConfigServer) AddOrSetQueryChannelPolicyForChannel(ctx context.Context, in *pb.ReqChannelPolicyQuery) (*pb.String, error) {
	config := service.Get(in.ConfigID)
	if nil == config {
		return nil, errors.New("config is nil")
	} else {
		return &pb.String{Data: "success"}, service.AddOrSetQueryChannelPolicyForChannel(&service.ChannelPolicyQuery{
			ConfigID:       in.ConfigID,
			ChannelName:    in.ChannelName,
			InitialBackOff: in.InitialBackOff,
			MaxBackOff:     in.MaxBackOff,
			MaxTargets:     in.MaxTargets,
			MinResponses:   in.MinResponses,
			Attempts:       in.Attempts,
			BackOffFactor:  in.BackOffFactor,
		})
	}
}

func (c *ConfigServer) AddOrSetDiscoveryPolicyForChannel(ctx context.Context, in *pb.ReqChannelPolicyDiscovery) (*pb.String, error) {
	config := service.Get(in.ConfigID)
	if nil == config {
		return nil, errors.New("config is nil")
	} else {
		return &pb.String{Data: "success"}, service.AddOrSetDiscoveryPolicyForChannel(&service.ChannelPolicyDiscovery{
			ConfigID:       in.ConfigID,
			ChannelName:    in.ChannelName,
			InitialBackOff: in.InitialBackOff,
			MaxBackOff:     in.MaxBackOff,
			MaxTargets:     in.MaxTargets,
			Attempts:       in.Attempts,
			BackOffFactor:  in.BackOffFactor,
		})
	}
}

func (c *ConfigServer) AddOrSetEventServicePolicyForChannel(ctx context.Context, in *pb.ReqChannelPolicyEvent) (*pb.String, error) {
	config := service.Get(in.ConfigID)
	if nil == config {
		return nil, errors.New("config is nil")
	} else {
		return &pb.String{Data: "success"}, service.AddOrSetEventServicePolicyForChannel(&service.ChannelPolicyEvent{
			ConfigID:                         in.ConfigID,
			ChannelName:                      in.ChannelName,
			ReconnectBlockHeightLagThreshold: in.ReconnectBlockHeightLagThreshold,
			ResolverStrategy:                 in.ResolverStrategy,
			BlockHeightLagThreshold:          in.BlockHeightLagThreshold,
			Balance:                          in.Balance,
			PeerMonitorPeriod:                in.PeerMonitorPeriod,
		})
	}
}

func (c *ConfigServer) AddOrSetOrdererForOrganizations(ctx context.Context, in *pb.ReqOrganizationsOrder) (*pb.String, error) {
	config := service.Get(in.ConfigID)
	if nil == config {
		return nil, errors.New("config is nil")
	} else {
		return &pb.String{Data: "success"}, service.AddOrSetOrdererForOrganizations(&service.OrganizationsOrder{
			ConfigID:   in.ConfigID,
			MspID:      in.MspID,
			CryptoPath: in.CryptoPath,
		})
	}
}

func (c *ConfigServer) AddOrSetOrgForOrganizations(ctx context.Context, in *pb.ReqOrganizationsOrg) (*pb.String, error) {
	config := service.Get(in.ConfigID)
	if nil == config {
		return nil, errors.New("config is nil")
	} else {
		return &pb.String{Data: "success"}, service.AddOrSetOrgForOrganizations(&service.OrganizationsOrg{
			ConfigID:               in.ConfigID,
			MspID:                  in.MspID,
			CryptoPath:             in.CryptoPath,
			OrgName:                in.OrgName,
			Peers:                  in.Peers,
			CertificateAuthorities: in.CertificateAuthorities,
		})
	}
}

func (c *ConfigServer) AddOrSetOrderer(ctx context.Context, in *pb.ReqOrder) (*pb.String, error) {
	config := service.Get(in.ConfigID)
	if nil == config {
		return nil, errors.New("config is nil")
	} else {
		return &pb.String{Data: "success"}, service.AddOrSetOrderer(&service.Order{
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
	}
}

func (c *ConfigServer) AddOrSetPeer(ctx context.Context, in *pb.ReqPeer) (*pb.String, error) {
	config := service.Get(in.ConfigID)
	if nil == config {
		return nil, errors.New("config is nil")
	} else {
		return &pb.String{Data: "success"}, service.AddOrSetPeer(&service.Peer{
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
	}
}

func (c *ConfigServer) AddOrSetCertificateAuthority(ctx context.Context, in *pb.ReqCertificateAuthority) (*pb.String, error) {
	config := service.Get(in.ConfigID)
	if nil == config {
		return nil, errors.New("config is nil")
	} else {
		return &pb.String{Data: "success"}, service.AddOrSetCertificateAuthority(&service.CertificateAuthority{
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
	}
}
