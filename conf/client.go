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
 */

package conf

import (
	"errors"
	"github.com/aberic/fabric-client/geneses"
	"github.com/aberic/fabric-client/grpc/proto/chain"
	"github.com/aberic/gnomon"
	"path"
	"path/filepath"
)

// Client go sdk 使用的客户端
type Client struct {
	// Organization 这个应用程序实例属于哪个组织?值必须是在“组织”下定义的组织的名称，如：Org1或league-org1
	Organization string `yaml:"organization"`
	// Logging 日志级别，debug、info、warn、error等
	Logging *ClientLogging `yaml:"logging"`
	// 节点超时的全局配置，如果省略此部分，则将使用缺省值
	Peer *ClientPeer `yaml:"peer"`
	// 事件服务超时的全局配置，如果省略此部分，则将使用缺省值
	EventService *ClientEventService `yaml:"eventService"`
	// orderer超时的全局配置，如果省略此部分，则将使用缺省值
	Order *ClientOrder `yaml:"orderer"`
	// 超时的全局配置，如果省略此部分，则将使用缺省值
	Global *ClientGlobal `yaml:"global"`
	// CryptoConfig 客户端
	CryptoConfig    *ClientCryptoConfig    `yaml:"cryptoconfig"`
	CredentialStore *ClientCredentialStore `yaml:"credentialStore"`
	// BCCSP 客户端的BCCSP配置
	BCCSP    *ClientBCCSP    `yaml:"BCCSP"`
	TLSCerts *ClientTLSCerts `yaml:"tlsCerts"`
}

// ClientLogging 客户端日志设置对象
type ClientLogging struct {
	Level string `yaml:"level"` // info
}

// ClientCryptoConfig 客户端
type ClientCryptoConfig struct {
	// Path 带有密钥和证书的MSP目录的根目录
	Path string `yaml:"path"` // /Users/Documents/fabric/crypto-config
}

type ClientCredentialStore struct {
	Path        string                            `yaml:"path"` // /tmp/state-store"
	CryptoStore *ClientCredentialStoreCryptoStore `yaml:"cryptoStore"`
}

type ClientCredentialStoreCryptoStore struct {
	Path string `yaml:"path"` // /tmp/msp
}

type ClientBCCSP struct {
	Security *ClientBCCSPSecurity `yaml:"security"`
}

type ClientBCCSPSecurity struct {
	Enabled       bool                        `yaml:"enabled"`
	Default       *ClientBCCSPSecurityDefault `yaml:"default"`
	HashAlgorithm string                      `yaml:"hashAlgorithm"`
	SoftVerify    bool                        `yaml:"softVerify"`
	Level         int32                       `yaml:"level"`
}

type ClientBCCSPSecurityDefault struct {
	Provider string `yaml:"provider"`
}

type ClientTLSCerts struct {
	// SystemCertPool 是否开启TLS，默认false
	SystemCertPool bool `yaml:"systemCertPool"`
	// Client 客户端密钥和证书，用于TLS与节点和排序服务的握手
	Client *ClientTLSCertsClient `yaml:"client"`
}

type ClientTLSCertsClient struct {
	Key  *ClientTLSCertsClientKey  `yaml:"key"`
	Cert *ClientTLSCertsClientCert `yaml:"cert"`
}

type ClientTLSCertsClientKey struct {
	Path string `yaml:"path"` // /fabric/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.key
}

type ClientTLSCertsClientCert struct {
	Path string `yaml:"path"` // /fabric/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.crt
}

type ClientPeer struct {
	Timeout *ClientPeerTimeout `yaml:"timeout"`
}

type ClientPeerTimeout struct {
	Connection string                      `yaml:"connection"`
	Response   string                      `yaml:"response"`
	Discovery  *ClientPeerTimeoutDiscovery `yaml:"discovery"`
}

type ClientPeerTimeoutDiscovery struct {
	// GreyListExpiry 发现服务失效列表筛选器的有效期。
	//
	// 通道客户端将列出脱机的失效节点名单，防止在后续重试中重新选择它们。
	//
	// 这个间隔将定义一个节点被灰列出的时间
	GreyListExpiry string `yaml:"greylistExpiry"`
}

type ClientEventService struct {
	Timeout *ClientEventServiceTimeout `yaml:"timeout"`
}

type ClientEventServiceTimeout struct {
	RegistrationResponse string `yaml:"registrationResponse"`
}

type ClientOrder struct {
	Timeout *ClientOrderTimeout `yaml:"timeout"`
}

type ClientOrderTimeout struct {
	Connection string `yaml:"connection"`
	Response   string `yaml:"response"`
}

type ClientGlobal struct {
	Timeout *ClientGlobalTimeout `yaml:"timeout"`
	Cache   *ClientGlobalCache   `yaml:"cache"`
}

type ClientGlobalTimeout struct {
	Query   string `yaml:"query"`
	Execute string `yaml:"execute"`
	Resmgmt string `yaml:"resmgmt"`
}

type ClientGlobalCache struct {
	ConnectionIdle    string `yaml:"connectionIdle"`
	EventServiceIdle  string `yaml:"eventServiceIdle"`
	ChannelConfig     string `yaml:"channelConfig"`
	ChannelMembership string `yaml:"channelMembership"`
	Discovery         string `yaml:"discovery"`
	Selection         string `yaml:"selection"`
}

func NewConfigClient(org *chain.OrgInfo) *Client {
	cryptoConfigPath := geneses.CryptoConfigPath(org.LeagueDomain)
	_, orgUserPath := geneses.CryptoOrgAndNodePath(org.LeagueDomain, org.Domain, org.Name, org.Username, true, geneses.CcnAdmin)
	return &Client{
		Logging: &ClientLogging{
			Level: "info",
		},
		Peer: &ClientPeer{
			Timeout: &ClientPeerTimeout{
				Connection: "10s",
				Response:   "180s",
				Discovery: &ClientPeerTimeoutDiscovery{
					GreyListExpiry: "10s",
				},
			},
		},
		EventService: &ClientEventService{
			Timeout: &ClientEventServiceTimeout{
				RegistrationResponse: "15s",
			},
		},
		Order: &ClientOrder{
			Timeout: &ClientOrderTimeout{
				Connection: "15s",
				Response:   "15s",
			},
		},
		Global: &ClientGlobal{
			Timeout: &ClientGlobalTimeout{
				Query:   "180s",
				Execute: "180s",
				Resmgmt: "180s",
			},
			Cache: &ClientGlobalCache{
				ConnectionIdle:    "30s",
				EventServiceIdle:  "2m",
				ChannelMembership: "30m",
				ChannelConfig:     "30s",
				Discovery:         "10s",
				Selection:         "10m",
			},
		},
		CryptoConfig: &ClientCryptoConfig{
			Path: cryptoConfigPath,
		},
		CredentialStore: &ClientCredentialStore{
			Path:        path.Join(orgUserPath, "msp", "signcerts"),
			CryptoStore: &ClientCredentialStoreCryptoStore{Path: path.Join(orgUserPath, "msp")},
		},
		BCCSP: &ClientBCCSP{
			Security: &ClientBCCSPSecurity{
				Enabled: true,
				Default: &ClientBCCSPSecurityDefault{
					Provider: "SW",
				},
				HashAlgorithm: "SHA2",
				SoftVerify:    true,
				Level:         256,
			},
		},
		TLSCerts: &ClientTLSCerts{
			SystemCertPool: true,
			Client: &ClientTLSCertsClient{
				Key: &ClientTLSCertsClientKey{
					Path: filepath.Join(orgUserPath, "tls", "client.key"),
				},
				Cert: &ClientTLSCertsClientCert{
					Path: filepath.Join(orgUserPath, "tls", "client.crt"),
				},
			},
		},
	}
}

func (c *Client) set(in *chain.Client) error {
	if gnomon.String().IsEmpty(in.Organization) {
		return errors.New("client org can't be empty")
	}
	c.Organization = in.Organization
	c.setLogging(in)
	c.setPeer(in)
	c.setEventService(in)
	c.setOrder(in)
	c.setGlobal(in)
	c.setCryptoConfig(in)
	c.setCredentialStore(in)
	c.setBCCSP(in)
	c.setTLSCerts(in)
	return nil
}

func (c *Client) setLogging(in *chain.Client) {
	if nil != in.Logging && gnomon.String().IsNotEmpty(in.Logging.Level) {
		c.Logging.Level = in.Logging.Level
	} else {
		c.Logging.Level = "info"
	}
}

func (c *Client) setPeer(in *chain.Client) {
	if nil != in.Peer && nil != in.Peer.Timeout {
		if gnomon.String().IsNotEmpty(in.Peer.Timeout.Connection) {
			c.Peer.Timeout.Connection = in.Peer.Timeout.Connection
			c.Peer.Timeout.Response = in.Peer.Timeout.Response
			c.Peer.Timeout.Discovery.GreyListExpiry = in.Peer.Timeout.Discovery.GreyListExpiry
		}
		if gnomon.String().IsNotEmpty(in.Peer.Timeout.Response) {
			c.Peer.Timeout.Response = in.Peer.Timeout.Response
		}
		if nil != in.Peer.Timeout.Discovery && gnomon.String().IsNotEmpty(in.Peer.Timeout.Discovery.GreyListExpiry) {
			c.Peer.Timeout.Discovery.GreyListExpiry = in.Peer.Timeout.Discovery.GreyListExpiry
		}
	}
}

func (c *Client) setEventService(in *chain.Client) {
	if nil != in.EventService && nil != in.EventService.Timeout && gnomon.String().IsNotEmpty(in.EventService.Timeout.RegistrationResponse) {
		c.EventService.Timeout.RegistrationResponse = in.EventService.Timeout.RegistrationResponse
	}
}

func (c *Client) setOrder(in *chain.Client) {
	if nil != in.Order && nil != in.Order.Timeout {
		if gnomon.String().IsNotEmpty(in.Order.Timeout.Connection) {
			c.Order.Timeout.Connection = in.Order.Timeout.Connection
		}
		if gnomon.String().IsNotEmpty(in.Order.Timeout.Response) {
			c.Order.Timeout.Response = in.Order.Timeout.Response
		}
	}
}

func (c *Client) setGlobal(in *chain.Client) {
	if nil != in.Global {
		if nil != in.Global.Timeout {
			if gnomon.String().IsNotEmpty(in.Global.Timeout.Query) {
				c.Global.Timeout.Query = in.Global.Timeout.Query
			}
			if gnomon.String().IsNotEmpty(in.Global.Timeout.Execute) {
				c.Global.Timeout.Execute = in.Global.Timeout.Execute
			}
			if gnomon.String().IsNotEmpty(in.Global.Timeout.Resmgmt) {
				c.Global.Timeout.Resmgmt = in.Global.Timeout.Resmgmt
			}
		}
		if nil != in.Global.Cache {
			if gnomon.String().IsNotEmpty(in.Global.Cache.ConnectionIdle) {
				c.Global.Cache.ConnectionIdle = in.Global.Cache.ConnectionIdle
			}
			if gnomon.String().IsNotEmpty(in.Global.Cache.EventServiceIdle) {
				c.Global.Cache.EventServiceIdle = in.Global.Cache.EventServiceIdle
			}
			if gnomon.String().IsNotEmpty(in.Global.Cache.ChannelMembership) {
				c.Global.Cache.ChannelMembership = in.Global.Cache.ChannelMembership
			}
			if gnomon.String().IsNotEmpty(in.Global.Cache.ChannelConfig) {
				c.Global.Cache.ChannelConfig = in.Global.Cache.ChannelConfig
			}
			if gnomon.String().IsNotEmpty(in.Global.Cache.Discovery) {
				c.Global.Cache.Discovery = in.Global.Cache.Discovery
			}
			if gnomon.String().IsNotEmpty(in.Global.Cache.Selection) {
				c.Global.Cache.Selection = in.Global.Cache.Selection
			}
		}
	}
}

func (c *Client) setCryptoConfig(in *chain.Client) {
	if nil != in.CryptoConfig && gnomon.String().IsNotEmpty(in.CryptoConfig.Path) {
		c.CryptoConfig.Path = in.CryptoConfig.Path
	}
}

func (c *Client) setCredentialStore(in *chain.Client) {
	if nil != in.CredentialStore {
		if gnomon.String().IsNotEmpty(in.CredentialStore.Path) {
			c.CredentialStore.Path = in.CredentialStore.Path
		}
		if nil != in.CredentialStore.CryptoStore && gnomon.String().IsNotEmpty(in.CredentialStore.CryptoStore.Path) {
			c.CredentialStore.CryptoStore.Path = in.CredentialStore.CryptoStore.Path
		}
	}
}

func (c *Client) setBCCSP(in *chain.Client) {
	if nil != in.BCCSP && nil != in.BCCSP.Security {
		c.BCCSP.Security.Enabled = in.BCCSP.Security.Enabled
		c.BCCSP.Security.SoftVerify = in.BCCSP.Security.SoftVerify
		if nil != in.BCCSP.Security.Default && gnomon.String().IsNotEmpty(in.BCCSP.Security.Default.Provider) {
			c.BCCSP.Security.Default.Provider = in.BCCSP.Security.Default.Provider
		}
		if gnomon.String().IsNotEmpty(in.BCCSP.Security.HashAlgorithm) {
			c.BCCSP.Security.HashAlgorithm = in.BCCSP.Security.HashAlgorithm
		}
		if in.BCCSP.Security.Level > 0 {
			c.BCCSP.Security.Level = in.BCCSP.Security.Level
		}
	}
}

func (c *Client) setTLSCerts(in *chain.Client) {
	if nil != in.TlsCerts {
		c.TLSCerts.SystemCertPool = in.TlsCerts.SystemCertPool
		if nil != in.TlsCerts.Client {
			if nil != in.TlsCerts.Client.Key && gnomon.String().IsNotEmpty(in.TlsCerts.Client.Key.Path) {
				c.TLSCerts.Client.Key.Path = in.TlsCerts.Client.Key.Path
			}
			if nil != in.TlsCerts.Client.Key && gnomon.String().IsNotEmpty(in.TlsCerts.Client.Cert.Path) {
				c.TLSCerts.Client.Cert.Path = in.TlsCerts.Client.Cert.Path
			}
		}
	}
}
