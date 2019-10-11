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
	"github.com/ennoo/rivet/utils/log"
	"gopkg.in/yaml.v2"
	"sync"
)

var (
	Configs map[string]*config.Config
	lock    sync.Mutex
)

func init() {
	Configs = map[string]*config.Config{}
}

func Get(configID string) *config.Config {
	return Configs[configID]
}

func GetBytes(configID string) []byte {
	if nil == Configs[configID] {
		return nil
	}
	confData, err := yaml.Marshal(Configs[configID])
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	return confData
}

func Recover(configIDs []string) {
	defer lock.Unlock()
	lock.Lock()
	for configID := range Configs {
		haveNo := true
		for _, recoverID := range configIDs {
			if configID == recoverID {
				haveNo = false
			}
		}
		if haveNo {
			delete(Configs, configID)
		}
	}
}

func RecoverConfig(configs map[string]*config.Config) {
	defer lock.Unlock()
	lock.Lock()
	Configs = configs
}

func InitConfig(in *pb.ReqInit) {
	conf := &config.Config{}
	conf.InitSelfClient(in.Client.Tls, in.Client.LeagueName, in.Client.Organization, in.Client.UserName, in.Client.Level)
	for _, peer := range in.ChannelPeer {
		conf.AddOrSetPeerForChannel(peer.ChannelName, peer.PeerName, peer.EndorsingPeer, peer.ChainCodeQuery,
			peer.LedgerQuery, peer.EventSource)
	}
	for _, query := range in.ChannelPolicyQuery {
		conf.AddOrSetQueryChannelPolicyForChannel(query.ChannelName, query.InitialBackOff,
			query.MaxBackOff, query.MinResponses, query.MaxTargets, query.Attempts, query.BackOffFactor)
	}
	for _, discovery := range in.ChannelPolicyDiscovery {
		conf.AddOrSetDiscoveryPolicyForChannel(discovery.ChannelName, discovery.InitialBackOff,
			discovery.MaxBackOff, discovery.MaxTargets, discovery.Attempts, discovery.BackOffFactor)
	}
	for _, event := range in.ChannelPolicyEvent {
		conf.AddOrSetEventServicePolicyForChannel(event.ChannelName, event.ResolverStrategy,
			event.Balance, event.PeerMonitorPeriod, event.BlockHeightLagThreshold, event.ReconnectBlockHeightLagThreshold)
	}
	conf.AddOrSetSelfOrdererForOrganizations(in.OrganizationsOrder.LeagueName)
	conf.AddOrSetSelfOrgForOrganizations(in.OrganizationsOrg.LeagueName, in.OrganizationsOrg.Peers, in.OrganizationsOrg.CertificateAuthorities)
	for _, order := range in.Order {
		conf.AddOrSetSelfOrderer(order.LeagueName, order.OrderName, order.Url, order.KeepAliveTime,
			order.KeepAliveTimeout, order.KeepAlivePermit, order.FailFast, order.AllowInsecure)
	}
	for _, peer := range in.Peer {
		conf.AddOrSetSelfPeer(peer.LeagueName, peer.PeerName, peer.Url, peer.EventUrl, peer.KeepAliveTime,
			peer.KeepAliveTimeout, peer.KeepAlivePermit, peer.FailFast, peer.AllowInsecure)
	}
	for _, cert := range in.CertificateAuthority {
		conf.AddOrSetSelfCertificateAuthority(cert.LeagueName, cert.CertName, cert.Url, cert.CaName, cert.EnrollId, cert.EnrollSecret)
	}
	defer lock.Unlock()
	lock.Lock()
	Configs[in.Client.ConfigID] = conf
}
