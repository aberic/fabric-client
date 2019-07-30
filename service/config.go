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
	configer "github.com/ennoo/fabric-client/config"
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
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

func InitConfig(in *pb.ReqInit) error {
	if nil == Configs[in.Client.ConfigID] {
		Configs[in.Client.ConfigID] = &config.Config{}
	}
	_ = InitClientSelf(in.Client)
	for _, peer := range in.ChannelPeer {
		_ = AddOrSetPeerForChannel(peer)
	}
	for _, query := range in.ChannelPolicyQuery {
		_ = AddOrSetQueryChannelPolicyForChannel(query)
	}
	for _, discovery := range in.ChannelPolicyDiscovery {
		_ = AddOrSetDiscoveryPolicyForChannel(discovery)
	}
	for _, event := range in.ChannelPolicyEvent {
		_ = AddOrSetEventServicePolicyForChannel(event)
	}
	_ = AddOrSetOrdererForOrganizationsSelf(in.OrganizationsOrder)
	_ = AddOrSetOrgForOrganizationsSelf(in.OrganizationsOrg)
	for _, order := range in.Order {
		_ = AddOrSetOrdererSelf(order)
	}
	for _, peer := range in.Peer {
		_ = AddOrSetPeerSelf(peer)
	}
	for _, cert := range in.CertificateAuthority {
		_ = AddOrSetCertificateAuthoritySelf(cert)
	}
	return nil
}

func InitClient(in *pb.ReqClient) error {
	client := &Client{
		ConfigID:     in.ConfigID,
		TlS:          in.Tls,
		Organization: in.Organization,
		Level:        in.Level,
		CryptoConfig: in.CryptoConfig,
		KeyPath:      in.KeyPath,
		CertPath:     in.CertPath,
	}
	if nil == Configs[client.ConfigID] {
		Configs[client.ConfigID] = &config.Config{}
	}
	return Configs[client.ConfigID].InitClient(client.TlS, client.Organization, client.Level, client.CryptoConfig,
		client.KeyPath, client.CertPath)
}

func InitClientSelf(in *pb.ReqClientSelf) error {
	client := &ClientSelf{
		ConfigID:     in.ConfigID,
		TlS:          in.Tls,
		LeagueName:   in.LeagueName,
		UserName:     in.UserName,
		Organization: in.Organization,
		Level:        in.Level,
	}
	if nil == Configs[client.ConfigID] {
		Configs[client.ConfigID] = &config.Config{}
	}
	return Configs[client.ConfigID].InitSelfClient(client.TlS, client.LeagueName, client.Organization, client.UserName,
		client.Level)
}

func InitClientCustom(in *pb.ReqClientCustom) error {
	clientCustom := &ClientCustom{
		ConfigID: in.ConfigID,
		Client: &Client{
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
	}
	if nil == Configs[clientCustom.ConfigID] {
		Configs[clientCustom.ConfigID] = &config.Config{}
	}
	return Configs[clientCustom.ConfigID].InitCustomClient(clientCustom.Client.TlS, clientCustom.Client.Organization,
		clientCustom.Client.Level, clientCustom.Client.CryptoConfig, clientCustom.Client.KeyPath,
		clientCustom.Client.CertPath, clientCustom.Peer, clientCustom.EventService, clientCustom.Order,
		clientCustom.Global, clientCustom.BCCSP)
}

func AddOrSetPeerForChannel(in *pb.ReqChannelPeer) error {
	channelPeer := &ChannelPeer{
		ConfigID:       in.ConfigID,
		ChannelName:    in.ChannelName,
		PeerName:       in.PeerName,
		EndorsingPeer:  in.EndorsingPeer,
		ChainCodeQuery: in.ChainCodeQuery,
		LedgerQuery:    in.LedgerQuery,
		EventSource:    in.EventSource,
	}
	if nil == Configs[channelPeer.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[channelPeer.ConfigID].AddOrSetPeerForChannel(channelPeer.ChannelName, channelPeer.PeerName,
		channelPeer.EndorsingPeer, channelPeer.ChainCodeQuery, channelPeer.LedgerQuery, channelPeer.EventSource)
	return nil
}

func AddOrSetQueryChannelPolicyForChannel(in *pb.ReqChannelPolicyQuery) error {
	query := &ChannelPolicyQuery{
		ConfigID:       in.ConfigID,
		ChannelName:    in.ChannelName,
		InitialBackOff: in.InitialBackOff,
		MaxBackOff:     in.MaxBackOff,
		MaxTargets:     in.MaxTargets,
		MinResponses:   in.MinResponses,
		Attempts:       in.Attempts,
		BackOffFactor:  in.BackOffFactor,
	}
	if nil == Configs[query.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[query.ConfigID].AddOrSetQueryChannelPolicyForChannel(query.ChannelName, query.InitialBackOff,
		query.MaxBackOff, query.MinResponses, query.MaxTargets, query.Attempts, query.BackOffFactor)
	return nil
}

func AddOrSetDiscoveryPolicyForChannel(in *pb.ReqChannelPolicyDiscovery) error {
	discovery := &ChannelPolicyDiscovery{
		ConfigID:       in.ConfigID,
		ChannelName:    in.ChannelName,
		InitialBackOff: in.InitialBackOff,
		MaxBackOff:     in.MaxBackOff,
		MaxTargets:     in.MaxTargets,
		Attempts:       in.Attempts,
		BackOffFactor:  in.BackOffFactor,
	}
	if nil == Configs[discovery.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[discovery.ConfigID].AddOrSetDiscoveryPolicyForChannel(discovery.ChannelName, discovery.InitialBackOff,
		discovery.MaxBackOff, discovery.MaxTargets, discovery.Attempts, discovery.BackOffFactor)
	return nil
}

func AddOrSetEventServicePolicyForChannel(in *pb.ReqChannelPolicyEvent) error {
	event := &ChannelPolicyEvent{
		ConfigID:                         in.ConfigID,
		ChannelName:                      in.ChannelName,
		ReconnectBlockHeightLagThreshold: in.ReconnectBlockHeightLagThreshold,
		ResolverStrategy:                 in.ResolverStrategy,
		BlockHeightLagThreshold:          in.BlockHeightLagThreshold,
		Balance:                          in.Balance,
		PeerMonitorPeriod:                in.PeerMonitorPeriod,
	}
	if nil == Configs[event.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[event.ConfigID].AddOrSetEventServicePolicyForChannel(event.ChannelName, event.ResolverStrategy,
		event.Balance, event.PeerMonitorPeriod, event.BlockHeightLagThreshold, event.ReconnectBlockHeightLagThreshold)
	return nil
}

func AddOrSetOrdererForOrganizations(in *pb.ReqOrganizationsOrder) error {
	order := &OrganizationsOrder{
		ConfigID:   in.ConfigID,
		MspID:      in.MspID,
		CryptoPath: in.CryptoPath,
		Users:      in.Users,
	}
	if nil == Configs[order.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[order.ConfigID].AddOrSetOrdererForOrganizations(order.MspID, order.CryptoPath, order.Users)
	return nil
}

func AddOrSetOrdererForOrganizationsSelf(in *pb.ReqOrganizationsOrderSelf) error {
	order := &OrganizationsOrderSelf{
		ConfigID:   in.ConfigID,
		LeagueName: in.LeagueName,
	}
	if nil == Configs[order.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[order.ConfigID].AddOrSetSelfOrdererForOrganizations(order.LeagueName)
	return nil
}

func AddOrSetOrgForOrganizations(in *pb.ReqOrganizationsOrg) error {
	org := &OrganizationsOrg{
		ConfigID:               in.ConfigID,
		MspID:                  in.MspID,
		CryptoPath:             in.CryptoPath,
		OrgName:                in.OrgName,
		Users:                  in.Users,
		Peers:                  in.Peers,
		CertificateAuthorities: in.CertificateAuthorities,
	}
	if nil == Configs[org.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[org.ConfigID].AddOrSetOrgForOrganizations(org.OrgName, org.MspID, org.CryptoPath, org.Users, org.Peers,
		org.CertificateAuthorities)
	return nil
}

func AddOrSetOrgForOrganizationsSelf(in *pb.ReqOrganizationsOrgSelf) error {
	org := &OrganizationsOrgSelf{
		ConfigID:               in.ConfigID,
		LeagueName:             in.LeagueName,
		Peers:                  in.Peers,
		CertificateAuthorities: in.CertificateAuthorities,
	}
	if nil == Configs[org.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[org.ConfigID].AddOrSetSelfOrgForOrganizations(org.LeagueName, org.Peers, org.CertificateAuthorities)
	return nil
}

func AddOrSetOrderer(in *pb.ReqOrder) error {
	order := &Order{
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
	}
	if nil == Configs[order.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[order.ConfigID].AddOrSetOrderer(order.OrderName, order.URL, order.SSLTargetNameOverride, order.KeepAliveTime,
		order.KeepAliveTimeout, order.TLSCACerts, order.KeepAlivePermit, order.FailFast, order.AllowInsecure)
	return nil
}

func AddOrSetOrdererSelf(in *pb.ReqOrderSelf) error {
	order := &OrderSelf{
		ConfigID:         in.ConfigID,
		OrderName:        in.OrderName,
		URL:              in.Url,
		LeagueName:       in.LeagueName,
		KeepAliveTime:    in.KeepAliveTime,
		KeepAliveTimeout: in.KeepAliveTimeout,
		KeepAlivePermit:  in.KeepAlivePermit,
		FailFast:         in.FailFast,
		AllowInsecure:    in.AllowInsecure,
	}
	if nil == Configs[order.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[order.ConfigID].AddOrSetSelfOrderer(order.LeagueName, order.OrderName, order.URL, order.KeepAliveTime,
		order.KeepAliveTimeout, order.KeepAlivePermit, order.FailFast, order.AllowInsecure)
	return nil
}

func AddOrSetPeer(in *pb.ReqPeer) error {
	peerUU := &Peer{
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
	}
	if nil == Configs[peerUU.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[peerUU.ConfigID].AddOrSetPeer(peerUU.PeerName, peerUU.URL, peerUU.EventUrl, peerUU.SSLTargetNameOverride,
		peerUU.KeepAliveTime, peerUU.KeepAliveTimeout, peerUU.TLSCACerts, peerUU.KeepAlivePermit, peerUU.FailFast,
		peerUU.AllowInsecure)
	return nil
}

func AddOrSetPeerSelf(in *pb.ReqPeerSelf) error {
	peerSelf := &PeerSelf{
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
	}
	if nil == Configs[peerSelf.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[peerSelf.ConfigID].AddOrSetSelfPeer(peerSelf.LeagueName, peerSelf.PeerName, peerSelf.URL, peerSelf.EventUrl, peerSelf.KeepAliveTime,
		peerSelf.KeepAliveTimeout, peerSelf.KeepAlivePermit, peerSelf.FailFast, peerSelf.AllowInsecure)
	return nil
}

func AddOrSetCertificateAuthority(in *pb.ReqCertificateAuthority) error {
	certificateAuthority := &CertificateAuthority{
		ConfigID:                in.ConfigID,
		CertName:                in.CertName,
		URL:                     in.Url,
		TLSCACertPath:           in.TlsCACertPath,
		TLSCACertClientKeyPath:  in.TlsCACertClientKeyPath,
		TLSCACertClientCertPath: in.TlsCACertClientCertPath,
		CAName:                  in.CaName,
		EnrollId:                in.EnrollId,
		EnrollSecret:            in.EnrollSecret,
	}
	if nil == Configs[certificateAuthority.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[certificateAuthority.ConfigID].AddOrSetCertificateAuthority(certificateAuthority.CertName,
		certificateAuthority.URL, certificateAuthority.TLSCACertPath,
		certificateAuthority.TLSCACertClientKeyPath, certificateAuthority.TLSCACertClientCertPath,
		certificateAuthority.CAName, certificateAuthority.EnrollId, certificateAuthority.EnrollSecret)
	return nil
}

func AddOrSetCertificateAuthoritySelf(in *pb.ReqCertificateAuthoritySelf) error {
	certificateAuthority := &CertificateAuthoritySelf{
		ConfigID:     in.ConfigID,
		CertName:     in.CertName,
		URL:          in.Url,
		LeagueName:   in.LeagueName,
		CAName:       in.CaName,
		EnrollId:     in.EnrollId,
		EnrollSecret: in.EnrollSecret,
	}
	if nil == Configs[certificateAuthority.ConfigID] {
		return errors.New("config client is not init")
	}
	Configs[certificateAuthority.ConfigID].AddOrSetSelfCertificateAuthority(certificateAuthority.LeagueName,
		certificateAuthority.CertName, certificateAuthority.URL, certificateAuthority.CAName,
		certificateAuthority.EnrollId, certificateAuthority.EnrollSecret)
	return nil
}
