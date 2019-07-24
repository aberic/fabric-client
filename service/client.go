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
 *
 */

package service

import (
	"github.com/ennoo/fabric-client/config"
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
)

type Client struct {
	ConfigID     string `json:"configID"`     // ConfigID 配置唯一ID
	TlS          bool   `json:"tls"`          // TlS 是否开启TLS，默认false
	Organization string `json:"organization"` // Organization 这个应用程序实例属于哪个组织?值必须是在“组织”下定义的组织的名称，如：Org1
	Level        string `json:"level"`        // Level 日志级别，debug、info、warn、error等
	CryptoConfig string `json:"cryptoConfig"` // CryptoConfig 带有密钥和证书的MSP目录的根目录
	KeyPath      string `json:"keyPath"`      // KeyPath 客户端密钥，用于TLS与节点和排序服务的握手
	CertPath     string `json:"certPath"`     // CertPath 客户端证书，用于TLS与节点和排序服务的握手
}

func (c *Client) Trans2pb() *pb.ReqClient {
	return &pb.ReqClient{
		ConfigID:     c.ConfigID,
		Tls:          c.TlS,
		Organization: c.Organization,
		Level:        c.Level,
		CryptoConfig: c.CryptoConfig,
		KeyPath:      c.KeyPath,
		CertPath:     c.CertPath,
	}
}

type ClientSelf struct {
	ConfigID     string `json:"configID"` // ConfigID 配置唯一ID
	TlS          bool   `json:"tls"`      // TlS 是否开启TLS，默认false
	LeagueName   string `json:"leagueName"`
	UserName     string `json:"userName"`
	Organization string `json:"organization"` // Organization 这个应用程序实例属于哪个组织?值必须是在“组织”下定义的组织的名称，如：Org1
	Level        string `json:"level"`        // Level 日志级别，debug、info、warn、error等
}

func (c *ClientSelf) Trans2pb() *pb.ReqClientSelf {
	return &pb.ReqClientSelf{
		ConfigID:     c.ConfigID,
		Tls:          c.TlS,
		LeagueName:   c.LeagueName,
		UserName:     c.UserName,
		Organization: c.Organization,
		Level:        c.Level,
	}
}

type ClientCustom struct {
	ConfigID     string                     `json:"configID"` // ConfigID 配置唯一ID
	Client       *Client                    `json:"client"`
	Peer         *config.ClientPeer         `json:"peer"`
	EventService *config.ClientEventService `json:"eventService"`
	Order        *config.ClientOrder        `json:"order"`
	Global       *config.ClientGlobal       `json:"global"`
	BCCSP        *config.ClientBCCSP        `json:"bccsp"`
}

func (c *ClientCustom) Trans2pb() *pb.ReqClientCustom {
	return &pb.ReqClientCustom{
		ConfigID: c.ConfigID,
		Client:   c.Client.Trans2pb(),
		Peer: &pb.ClientPeer{
			Timeout: &pb.ClientPeerTimeout{
				Connection: c.Peer.Timeout.Connection,
				Response:   c.Peer.Timeout.Response,
				Discovery: &pb.ClientPeerTimeoutDiscovery{
					GreyListExpiry: c.Peer.Timeout.Discovery.GreyListExpiry,
				},
			},
		},
		EventService: &pb.ClientEventService{
			Timeout: &pb.ClientEventServiceTimeout{
				RegistrationResponse: c.EventService.Timeout.RegistrationResponse,
			},
		},
		Order: &pb.ClientOrder{
			Timeout: &pb.ClientOrderTimeout{
				Connection: c.Order.Timeout.Connection,
				Response:   c.Order.Timeout.Response,
			},
		},
		Global: &pb.ClientGlobal{
			Timeout: &pb.ClientGlobalTimeout{
				Query:   c.Global.Timeout.Query,
				Execute: c.Global.Timeout.Execute,
				Resmgmt: c.Global.Timeout.Resmgmt,
			},
			Cache: &pb.ClientGlobalCache{
				ConnectionIdle:    c.Global.Cache.ConnectionIdle,
				EventServiceIdle:  c.Global.Cache.EventServiceIdle,
				ChannelMembership: c.Global.Cache.ChannelMembership,
				ChannelConfig:     c.Global.Cache.ChannelConfig,
				Discovery:         c.Global.Cache.Discovery,
				Selection:         c.Global.Cache.Selection,
			},
		},
		BCCSP: &pb.ClientBCCSP{
			Security: &pb.ClientBCCSPSecurity{
				Enabled: c.BCCSP.Security.Enabled,
				Default: &pb.ClientBCCSPSecurityDefault{
					Provider: c.BCCSP.Security.Default.Provider,
				},
				HashAlgorithm: c.BCCSP.Security.HashAlgorithm,
				SoftVerify:    c.BCCSP.Security.SoftVerify,
				Level:         c.BCCSP.Security.Level,
			},
		},
	}
}
