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

func AddOrSetOrderer(order *Orderer) error {
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
