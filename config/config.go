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

const (
	OrderOrgKey = "ordererorg"
)

type Config struct {
	Version                string                           `yaml:"version"`
	Client                 *Client                          `yaml:"client"`
	Channels               map[string]*Channel              `yaml:"channels"`
	Organizations          map[string]interface{}           `yaml:"organizations"`
	Orderers               map[string]*Orderer              `yaml:"orderers"`
	Peers                  map[string]*Peer                 `yaml:"peers"`
	CertificateAuthorities map[string]*CertificateAuthority `yaml:"certificateAuthorities"`
}

func (c *Config) InitClient(tls bool, organization, level, cryptoConfig, keyPath, certPath string) error {
	c.initClient()
	return c.Client.initClient(tls, organization, level, cryptoConfig, keyPath, certPath)
}

func (c *Config) InitCustomClient(tls bool, organization, level, cryptoConfig, keyPath, certPath string,
	peer *ClientPeer, eventService *ClientEventService, order *ClientOrder, global *ClientGlobal, bccsp *ClientBCCSP) error {
	c.initClient()
	return c.Client.initCustomClient(tls, organization, level, cryptoConfig, keyPath, certPath, peer, eventService,
		order, global, bccsp)
}

func (c *Config) AddOrSetPeerForChannel(channelName, peerName string, endorsingPeer, chaincodeQuery, ledgerQuery, eventSource bool) {
	c.initChannel(channelName, peerName)
	c.Channels[channelName].Peers[peerName] = &ChannelPeer{
		EndorsingPeer:  endorsingPeer,
		ChaincodeQuery: chaincodeQuery,
		LedgerQuery:    ledgerQuery,
		EventSource:    eventSource,
	}
}

func (c *Config) AddOrSetQueryChannelPolicyForChannel(channelName, initialBackOff, maxBackOff string,
	minResponses, maxTargets, attempts int, backOffFactor float32) {
	c.initChannel(channelName)
	c.Channels[channelName].Policies.QueryChannelConfig = &PolicyQueryChannelConfig{
		MinResponses: minResponses,
		MaxTargets:   maxTargets,
		RetryOpts: &PolicyCommonRetryOpts{
			Attempts:       attempts,
			InitialBackOff: initialBackOff,
			MaxBackOff:     maxBackOff,
			BackOffFactor:  backOffFactor,
		},
	}
}

func (c *Config) AddOrSetDiscoveryPolicyForChannel(channelName, initialBackOff, maxBackOff string,
	maxTargets, attempts int, backOffFactor float32) {
	c.initChannel(channelName)
	c.Channels[channelName].Policies.Discovery = &PolicyDiscovery{
		MaxTargets: maxTargets,
		RetryOpts: &PolicyCommonRetryOpts{
			Attempts:       attempts,
			InitialBackOff: initialBackOff,
			MaxBackOff:     maxBackOff,
			BackOffFactor:  backOffFactor,
		},
	}
}

func (c *Config) AddOrSetEventServicePolicyForChannel(channelName, resolverStrategy, balancer, peerMonitorPeriod string,
	blockHeightLagThreshold, reconnectBlockHeightLagThreshold int) {
	c.initChannel(channelName)
	c.Channels[channelName].Policies.EventService = &PolicyEventService{
		ResolverStrategy:                 resolverStrategy,
		BlockHeightLagThreshold:          blockHeightLagThreshold,
		Balancer:                         balancer,
		ReconnectBlockHeightLagThreshold: reconnectBlockHeightLagThreshold,
		PeerMonitorPeriod:                peerMonitorPeriod,
	}
}

func (c *Config) AddOrSetOrdererForOrganizations(mspID, cryptoPath string) {
	c.initOrganizations()
	c.Organizations[OrderOrgKey] = &OrdererOrg{
		MspID:      mspID,
		CryptoPath: cryptoPath,
	}
}

func (c *Config) AddOrSetOrgForOrganizations(orgName, mspid, cryptoPath string, peers, certificateAuthorities []string) {
	c.initOrganizations()
	c.Organizations[orgName] = &Org{
		MspID:                  mspid,
		CryptoPath:             cryptoPath,
		Peers:                  peers,
		CertificateAuthorities: certificateAuthorities,
	}
}

func (c *Config) AddOrSetOrderer(ordererName, url, sslTargetNameOverride, keepAliveTime, keepAliveTimeout,
	tlsCACerts string, keepAlivePermit, failFast, allowInsecure bool) {
	c.initOrderers(ordererName)
	c.Orderers[ordererName] = &Orderer{
		URL: url,
		GRPCOptions: &OrdererGRPCOptions{
			SSLTargetNameOverride: sslTargetNameOverride,
			KeepAliveTime:         keepAliveTime,
			KeepAliveTimeout:      keepAliveTimeout,
			KeepAlivePermit:       keepAlivePermit,
			FailFast:              failFast,
			AllowInsecure:         allowInsecure,
		},
		TLSCACerts: &OrdererTLSCACerts{
			Path: tlsCACerts,
		},
	}
}

func (c *Config) AddOrSetPeer(peerName, url, eventUrl, sslTargetNameOverride, keepAliveTime, keepAliveTimeout,
	tlsCACerts string, keepAlivePermit, failFast, allowInsecure bool) {
	c.initPeers(peerName)
	c.Peers[peerName] = &Peer{
		URL:      url,
		EventURL: eventUrl,
		GRPCOptions: &PeerGRPCOptions{
			SSLTargetNameOverride: sslTargetNameOverride,
			KeepAliveTime:         keepAliveTime,
			KeepAliveTimeout:      keepAliveTimeout,
			KeepAlivePermit:       keepAlivePermit,
			FailFast:              failFast,
			AllowInsecure:         allowInsecure,
		},
		TLSCACerts: &PeerTLSCACerts{
			Path: tlsCACerts,
		},
	}
}

func (c *Config) AddOrSetCertificateAuthority(certName, url, tlsCACertPath, tlsCACertClientKeyPath,
	tlsCACertClientCertPath, caName, enrollId, enrollSecret string) {
	c.initCertificateAuthorities(certName)
	c.CertificateAuthorities[certName] = &CertificateAuthority{
		URL: url,
		TLSCACerts: &CertificateAuthorityTLSCACerts{
			Path: tlsCACertPath,
			Client: &CertificateAuthorityTLSCACertsClient{
				Key: &CertificateAuthorityTLSCACertsClientKey{
					Path: tlsCACertClientKeyPath,
				},
				Cert: &CertificateAuthorityTLSCACertsClientCert{
					Path: tlsCACertClientCertPath,
				},
			},
		},
		Registrar: &CertificateAuthorityRegistrar{
			EnrollId:     enrollId,
			EnrollSecret: enrollSecret,
		},
		CAName: caName,
	}
}

func (c *Config) initClient() {
	c.Version = "1.0.0"
	c.Client = &Client{}
}

func (c *Config) initChannel(channelName string, peerNames ...string) {
	if nil == c.Channels {
		c.Channels = map[string]*Channel{}
		goto NewChannel
	} else {
		if nil == c.Channels[channelName] {
			goto NewChannel
		} else if len(peerNames) > 0 {
			goto NewPeer
		}
	}
	return
NewChannel:
	c.Channels[channelName] = &Channel{
		Peers: map[string]*ChannelPeer{},
		Policies: &Policy{
			QueryChannelConfig: &PolicyQueryChannelConfig{
				RetryOpts: &PolicyCommonRetryOpts{},
			},
			Discovery: &PolicyDiscovery{
				RetryOpts: &PolicyCommonRetryOpts{},
			},
			EventService: &PolicyEventService{},
		},
	}
	goto NewPeer
NewPeer:
	for index := range peerNames {
		c.Channels[channelName].Peers[peerNames[index]] = &ChannelPeer{}
	}
}

func (c *Config) initOrganizations() {
	if nil == c.Organizations {
		c.Organizations = map[string]interface{}{}
	}
}

func (c *Config) initOrderers(ordererName string) {
	if nil == c.Orderers {
		c.Orderers = map[string]*Orderer{}
		goto NewOrder
	} else {
		if nil == c.Orderers[ordererName] {
			goto NewOrder
		}
	}
	return
NewOrder:
	c.Orderers[ordererName] = &Orderer{
		GRPCOptions: &OrdererGRPCOptions{},
		TLSCACerts:  &OrdererTLSCACerts{},
	}
}

func (c *Config) initPeers(peerName string) {
	if nil == c.Peers {
		c.Peers = map[string]*Peer{}
		goto NewPeer
	} else {
		if nil == c.Peers[peerName] {
			goto NewPeer
		}
	}
	return
NewPeer:
	c.Peers[peerName] = &Peer{
		GRPCOptions: &PeerGRPCOptions{},
		TLSCACerts:  &PeerTLSCACerts{},
	}
}

func (c *Config) initCertificateAuthorities(certName string) {
	if nil == c.CertificateAuthorities {
		c.CertificateAuthorities = map[string]*CertificateAuthority{}
		goto NewCertificateAuthority
	} else {
		if nil == c.CertificateAuthorities[certName] {
			goto NewCertificateAuthority
		}
	}
	return
NewCertificateAuthority:
	c.CertificateAuthorities[certName] = &CertificateAuthority{
		TLSCACerts: &CertificateAuthorityTLSCACerts{
			Client: &CertificateAuthorityTLSCACertsClient{
				Key:  &CertificateAuthorityTLSCACertsClientKey{},
				Cert: &CertificateAuthorityTLSCACertsClientCert{},
			},
		},
		Registrar: &CertificateAuthorityRegistrar{},
	}
}
