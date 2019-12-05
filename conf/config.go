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
	"fmt"
	"github.com/aberic/fabric-client/grpc/proto/chain"
	"github.com/aberic/gnomon"
)

// Config 网络连接配置为客户端应用程序提供有关目标区块链网络的信息
type Config struct {
	Version       string                   `yaml:"version"`       // Version 内容的版本。用于SDK应用相应的解析规则
	Client        *Client                  `yaml:"client"`        // Client GO SDK使用的客户端
	Channels      map[string]*Channel      `yaml:"channels"`      // Channels 可选，如果有通道操作则需要补充完整
	Organizations map[string]*Organization `yaml:"organizations"` // Organizations 此网络的参与机构名单
	// Orderers
	//
	// 发送事务和通道创建/更新请求的Order列表。如果定义了多个，那么SDK将根据文档定义来使用特定的Order
	Orderers map[string]*Orderer `yaml:"orderers"`
	// Peers
	//
	// 发送各种请求的节点列表，包括背书、查询和事件侦听器注册。
	Peers map[string]*Peer `yaml:"peers"`
	// CertificateAuthorities
	//
	// Fabric- ca是由Hyperledger Fabric提供的一种特殊的证书颁发机构，它允许通过REST api进行证书管理。
	//
	// 应用程序可以选择使用标准的证书颁发机构，而不是Fabric-CA，在这种情况下，不会指定此部分。
	CertificateAuthorities map[string]*CertificateAuthority `yaml:"certificateAuthorities"`
}

func (c *Config) Set(in *chain.ReqConfigSet) error {
	var err error
	if gnomon.String().IsEmpty(in.Version) {
		err = errors.New("version can't be nil")
		goto ERR
	}
	c.Version = in.Version
	if err = c.setClient(in); nil != err {
		goto ERR
	}
	c.setChannels(in)
	c.setOrganizations(in)
	if err = c.setOrderers(in); nil != err {
		goto ERR
	}
	if err = c.setPeers(in); nil != err {
		goto ERR
	}
	if err = c.setCertificateAuthorities(in); nil != err {
		goto ERR
	}
	return nil
ERR:
	return fmt.Errorf("config set error: %w", err)
}

func (c *Config) setClient(in *chain.ReqConfigSet) error {
	client := NewConfigClient(in.OrgInfo)
	if err := client.set(in.Client); nil != err {
		return err
	}
	c.Client = client
	return nil
}

func (c *Config) setChannels(in *chain.ReqConfigSet) {
	c.Channels = make(map[string]*Channel)
	c.Channels["_default"] = NewConfigChannel()
	for channelName, channel := range in.Channels {
		ch := NewConfigChannel()
		ch.set(channel)
		c.Channels[channelName] = ch
	}
}

func (c *Config) setOrganizations(in *chain.ReqConfigSet) {
	c.Organizations = make(map[string]*Organization)
	for orgName, organization := range in.Organizations {
		org := &Organization{}
		org.set(in.OrgInfo.LeagueDomain, orgName, organization)
		c.Organizations[orgName] = org
	}
}

func (c *Config) setOrderers(in *chain.ReqConfigSet) error {
	c.Orderers = make(map[string]*Orderer)
	for ordererName, orderer := range in.Orderers {
		order := &Orderer{GRPCOptions: &OrdererGRPCOptions{}, TLSCACerts: &OrdererTLSCACerts{}}
		if err := order.set(in.OrgInfo.LeagueDomain, orderer); nil != err {
			return err
		}
		c.Orderers[ordererName] = order
	}
	return nil
}

func (c *Config) setPeers(in *chain.ReqConfigSet) error {
	c.Peers = make(map[string]*Peer)
	for peerName, peer := range in.Peers {
		p := &Peer{GRPCOptions: &PeerGRPCOptions{}, TLSCACerts: &PeerTLSCACerts{}}
		if err := p.set(in.OrgInfo.LeagueDomain, peer); nil != err {
			return err
		}
		c.Peers[peerName] = p
	}
	return nil
}

func (c *Config) setCertificateAuthorities(in *chain.ReqConfigSet) error {
	c.CertificateAuthorities = make(map[string]*CertificateAuthority)
	for caName, certificateAuthority := range in.CertificateAuthorities {
		ca := &CertificateAuthority{
			TLSCACerts: &CertificateAuthorityTLSCACerts{Client: &CertificateAuthorityTLSCACertsClient{
				Key:  &CertificateAuthorityTLSCACertsClientKey{},
				Cert: &CertificateAuthorityTLSCACertsClientCert{},
			}},
			Registrar: &CertificateAuthorityRegistrar{},
		}
		if err := ca.set(in.OrgInfo.LeagueDomain, caName, certificateAuthority); nil != err {
			return err
		}
		c.CertificateAuthorities[caName] = ca
	}
	return nil
}
