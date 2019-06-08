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

package config

import (
	"strings"
)

// Client go sdk 使用的客户端
type Client struct {
	// Organization 这个应用程序实例属于哪个组织?值必须是在“组织”下定义的组织的名称，如：Org1
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
	Global          *ClientGlobal          `yaml:"global"`
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
	Level         int                         `yaml:"level"`
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

func (c *Client) initClient(tls bool, organization, level, cryptoConfig, keyPath, certPath string) error {
	return c.initCustomClient(
		tls, organization, level, cryptoConfig, keyPath, certPath,
		&ClientPeer{
			Timeout: &ClientPeerTimeout{
				Connection: "10s",
				Response:   "180s",
				Discovery: &ClientPeerTimeoutDiscovery{
					GreyListExpiry: "10s",
				},
			},
		},
		&ClientEventService{
			Timeout: &ClientEventServiceTimeout{
				RegistrationResponse: "15s",
			},
		},
		&ClientOrder{
			Timeout: &ClientOrderTimeout{
				Connection: "15s",
				Response:   "15s",
			},
		},
		&ClientGlobal{
			Timeout: &ClientGlobalTimeout{
				Query:   "180s",
				Execute: "180s",
				Resmgmt: "180s",
			},
			Cache: &ClientGlobalCache{
				ConnectionIdle:    "30s",
				EventServiceIdle:  "2m",
				ChannelConfig:     "30m",
				ChannelMembership: "30s",
				Discovery:         "10s",
				Selection:         "10m",
			},
		},
		&ClientBCCSP{
			Security: &ClientBCCSPSecurity{
				Enabled:       true,
				HashAlgorithm: "SHA2",
				SoftVerify:    true,
				Level:         256,
				Default:       &ClientBCCSPSecurityDefault{Provider: "SW"},
			},
		})
}

func (c *Client) initCustomClient(tls bool, organization, level, cryptoConfig, keyPath, certPath string,
	peer *ClientPeer, eventService *ClientEventService, order *ClientOrder, global *ClientGlobal, bccsp *ClientBCCSP) error {
	c.Organization = organization
	c.Logging = &ClientLogging{Level: level}
	c.Peer = peer
	c.EventService = eventService
	c.Order = order
	c.Global = global
	c.CryptoConfig = &ClientCryptoConfig{Path: cryptoConfig}
	c.CredentialStore = &ClientCredentialStore{
		Path:        strings.Join([]string{"/tmp", organization, "state-store"}, "/"),
		CryptoStore: &ClientCredentialStoreCryptoStore{Path: strings.Join([]string{"/tmp", organization, "msp"}, "/")},
	}
	c.BCCSP = bccsp
	c.TLSCerts = &ClientTLSCerts{
		SystemCertPool: tls,
		Client: &ClientTLSCertsClient{
			Key:  &ClientTLSCertsClientKey{Path: keyPath},
			Cert: &ClientTLSCertsClientCert{Path: certPath},
		},
	}
	return nil
}
