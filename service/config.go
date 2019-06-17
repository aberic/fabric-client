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
	"github.com/ennoo/fabric-go-client/config"
	"github.com/ennoo/rivet/utils/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	Configs map[string]*config.Config
)

func init() {
	Configs = map[string]*config.Config{}
}

func Get(configID string) *config.Config {
	return Configs[configID]
}

func GetBytes(configID string) []byte {
	confData, err := yaml.Marshal(Configs[configID])
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	return confData
}

func InitClient(client *Client) error {
	if nil == Configs[client.ConfigID] {
		Configs[client.ConfigID] = &config.Config{}
	}
	return Configs[client.ConfigID].InitClient(client.TlS, client.Organization, client.Level, client.CryptoConfig,
		client.KeyPath, client.CertPath)
}

func InitClientCustom(clientCustom *ClientCustom) error {
	if nil == Configs[clientCustom.ConfigID] {
		Configs[clientCustom.ConfigID] = &config.Config{}
	}
	return Configs[clientCustom.ConfigID].InitCustomClient(clientCustom.Client.TlS, clientCustom.Client.Organization,
		clientCustom.Client.Level, clientCustom.Client.CryptoConfig, clientCustom.Client.KeyPath,
		clientCustom.Client.CertPath, clientCustom.Peer, clientCustom.EventService, clientCustom.Order,
		clientCustom.Global, clientCustom.BCCSP)
}

func AddOrSetPeerForChannel(channelPeer *ChannelPeer) error {
	if nil == Configs[channelPeer.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[channelPeer.ConfigID].AddOrSetPeerForChannel(channelPeer.ChannelName, channelPeer.PeerName,
		channelPeer.EndorsingPeer, channelPeer.ChainCodeQuery, channelPeer.LedgerQuery, channelPeer.EventSource)
	return nil
}

func AddOrSetQueryChannelPolicyForChannel(query *ChannelPolicyQuery) error {
	if nil == Configs[query.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[query.ConfigID].AddOrSetQueryChannelPolicyForChannel(query.ChannelName, query.InitialBackOff,
		query.MaxBackOff, query.MinResponses, query.MaxTargets, query.Attempts, query.BackOffFactor)
	return nil
}

func AddOrSetDiscoveryPolicyForChannel(discovery *ChannelPolicyDiscovery) error {
	if nil == Configs[discovery.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[discovery.ConfigID].AddOrSetDiscoveryPolicyForChannel(discovery.ChannelName, discovery.InitialBackOff,
		discovery.MaxBackOff, discovery.MaxTargets, discovery.Attempts, discovery.BackOffFactor)
	return nil
}

func AddOrSetEventServicePolicyForChannel(event *ChannelPolicyEvent) error {
	if nil == Configs[event.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[event.ConfigID].AddOrSetEventServicePolicyForChannel(event.ChannelName, event.ResolverStrategy,
		event.Balance, event.PeerMonitorPeriod, event.BlockHeightLagThreshold, event.ReconnectBlockHeightLagThreshold)
	return nil
}

func AddOrSetOrdererForOrganizations(order *OrganizationsOrder) error {
	if nil == Configs[order.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[order.ConfigID].AddOrSetOrdererForOrganizations(order.MspID, order.CryptoPath)
	return nil
}

func AddOrSetOrgForOrganizations(org *OrganizationsOrg) error {
	if nil == Configs[org.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[org.ConfigID].AddOrSetOrgForOrganizations(org.OrgName, org.MspID, org.CryptoPath, org.Peers,
		org.CertificateAuthorities)
	return nil
}

func AddOrSetOrderer(order *Order) error {
	if nil == Configs[order.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[order.ConfigID].AddOrSetOrderer(order.OrderName, order.URL, order.SSLTargetNameOverride, order.KeepAliveTime,
		order.KeepAliveTimeout, order.TLSCACerts, order.KeepAlivePermit, order.FailFast, order.AllowInsecure)
	return nil
}

func AddOrSetPeer(peer *Peer) error {
	if nil == Configs[peer.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[peer.ConfigID].AddOrSetPeer(peer.PeerName, peer.URL, peer.EventUrl, peer.SSLTargetNameOverride,
		peer.KeepAliveTime, peer.KeepAliveTimeout, peer.TLSCACerts, peer.KeepAlivePermit, peer.FailFast,
		peer.AllowInsecure)
	return nil
}

func AddOrSetCertificateAuthority(certificateAuthority *CertificateAuthority) error {
	if nil == Configs[certificateAuthority.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[certificateAuthority.ConfigID].AddOrSetCertificateAuthority(certificateAuthority.CertName,
		certificateAuthority.URL, certificateAuthority.TLSCACertPath,
		certificateAuthority.TLSCACertClientKeyPath, certificateAuthority.TLSCACertClientCertPath,
		certificateAuthority.CAName, certificateAuthority.EnrollId, certificateAuthority.EnrollSecret)
	return nil
}

func TestConfig() *config.Config {
	//rootPath := "/Users/aberic/Documents/path/go/src/github.com/ennoo/fabric-go-client/example"
	rootPath := "/Users/admin/Documents/code/git/go/src/github.com/ennoo/fabric-go-client/example"
	conf := config.Config{}
	_ = conf.InitClient(false, "Org1", "debug",
		rootPath+"/config/crypto-config",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.key",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.crt")
	conf.AddOrSetPeerForChannel("mychannel", "peer0.league01-org1-vh-cn",
		true, true, true, true)
	conf.AddOrSetQueryChannelPolicyForChannel("mychannel", "500ms", "5s",
		1, 1, 5, 2.0)
	conf.AddOrSetDiscoveryPolicyForChannel("mychannel", "500ms", "5s",
		2, 4, 2.0)
	conf.AddOrSetEventServicePolicyForChannel("mychannel", "PreferOrg", "Random",
		"6s", 5, 8)
	conf.AddOrSetOrdererForOrganizations("OrdererMSP",
		rootPath+"/config/crypto-config/ordererOrganizations/league01-vh-cn/users/Admin@league01-vh-cn/msp")
	conf.AddOrSetOrgForOrganizations("Org1", "Org1MSP",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/msp",
		[]string{"peer0.league01-org1-vh-cn", "peer1.league01-org1-vh-cn"},
		[]string{"ca0.league01-org1-vh-cn"},
	)
	conf.AddOrSetOrgForOrganizations("Org2", "Org2MSP",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org2-vh-cn/users/Admin@league01-org2-vh-cn/msp",
		[]string{"peer1.league01-org2-vh-cn"},
		[]string{"ca0.league01-org2-vh-cn"},
	)
	conf.AddOrSetOrderer("order0.league01-vh-cn:7050", "grpc://10.10.203.51:30054",
		"order0.league01-vh-cn", "0s", "20s",
		rootPath+"/config/crypto-config/ordererOrganizations/league01-vh-cn/tlsca/tlsca.league01-vh-cn-cert.pem",
		false, false, false)
	conf.AddOrSetPeer("peer0.league01-org1-vh-cn", "grpc://10.10.203.51:30056",
		"grpc://10.10.203.51:30058", "peer0.league01-org1-vh-cn",
		"0s", "20s",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/tlsca/tlsca.league01-org1-vh-cn-cert.pem",
		false, false, false)
	conf.AddOrSetPeer("peer1.league01-org1-vh-cn", "grpc://10.10.203.51:30061",
		"grpc://10.10.203.51:30063", "peer1.league01-org1-vh-cn",
		"0s", "20s",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/tlsca/tlsca.league01-org1-vh-cn-cert.pem",
		false, false, false)
	conf.AddOrSetPeer("peer1.league01-org2-vh-cn", "grpc://10.10.203.51:30065",
		"grpc://10.10.203.51:30067", "peer1.league01-org2-vh-cn",
		"0s", "20s",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org2-vh-cn/tlsca/tlsca.league01-org2-vh-cn-cert.pem",
		false, false, false)
	conf.AddOrSetCertificateAuthority("ca.league01-vh-cn", "https://10.10.203.51:30059",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/tlsca/tlsca.league01-org1-vh-cn-cert.pem",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.key",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.crt",
		"admin", "adminpw", "ca.league01-vh-cn")
	return &conf
}
