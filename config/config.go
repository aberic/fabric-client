/*
 * Copyright (c) 2019. Aberic - All Rights Reserved.
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
	"github.com/aberic/fabric-client/geneses"
	pb "github.com/aberic/fabric-client/grpc/proto/chain"
	"strings"
)

const (
	OrderOrgKey = "OrdererOrg"
)

type Config struct {
	Version                string                           `yaml:"version"`
	Client                 *Client                          `yaml:"client"`
	Channels               map[string]*Channel              `yaml:"channels"`
	Organizations          map[string]*Organization         `yaml:"organizations"`
	Orderers               map[string]*Orderer              `yaml:"orderers"`
	Peers                  map[string]*Peer                 `yaml:"peers"`
	CertificateAuthorities map[string]*CertificateAuthority `yaml:"certificateAuthorities"`
}

func (c *Config) InitClient(tls bool, orgName, level, cryptoConfig, keyPath, certPath string) {
	c.initClient()
	c.Client.initClient(tls, orgName, level, cryptoConfig, keyPath, certPath)
}

func (c *Config) InitSelfClient(tls bool, leagueName, orgName, userName, level string) {
	c.initClient()
	c.Client.initSelfClient(tls, leagueName, orgName, userName, level)
}

func (c *Config) InitCustomClient(tls bool, orgName, level, cryptoConfig, keyPath, certPath string,
	peer *ClientPeer, eventService *ClientEventService, order *ClientOrder, global *ClientGlobal,
	ccs *ClientCredentialStore, bccsp *ClientBCCSP) {
	c.initClient()
	c.Client.initCustomClient(tls, orgName, level, cryptoConfig, keyPath, certPath, peer, eventService,
		order, global, ccs, bccsp)
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
	minResponses, maxTargets, attempts int32, backOffFactor float32) {
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
	maxTargets, attempts int32, backOffFactor float32) {
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
	blockHeightLagThreshold, reconnectBlockHeightLagThreshold int64) {
	c.initChannel(channelName)
	c.Channels[channelName].Policies.EventService = &PolicyEventService{
		ResolverStrategy:                 resolverStrategy,
		BlockHeightLagThreshold:          blockHeightLagThreshold,
		Balancer:                         balancer,
		ReconnectBlockHeightLagThreshold: reconnectBlockHeightLagThreshold,
		PeerMonitorPeriod:                peerMonitorPeriod,
	}
}

func (c *Config) AddOrSetOrdererForOrganizations(orderName, mspID, cryptoPath string, users map[string]string) {
	c.initOrganizations()
	userMap := map[string]*User{}
	for username, path := range users {
		userMap[username] = &User{Cert: &Cert{Path: path}}
	}
	c.Organizations[orderName] = &Organization{
		MspID:      mspID,
		CryptoPath: cryptoPath,
		Users:      userMap,
	}
}

func (c *Config) AddOrSetSelfOrdererForOrganizations(leagueName string) {
	c.initOrganizations()
	cryptoPath := strings.Join([]string{
		geneses.CryptoConfigPath(leagueName), "/ordererOrganizations/", leagueName, "/users/Admin@", leagueName, "/msp"},
		"")
	userCertPath := strings.Join([]string{cryptoPath, "/signcerts/Admin@", leagueName, "-cert.pem"}, "")
	c.Organizations[OrderOrgKey] = &Organization{
		MspID:      "OrdererMSP",
		CryptoPath: cryptoPath,
		Users:      map[string]*User{"Admin": {Cert: &Cert{Path: userCertPath}}},
	}
}

func (c *Config) AddOrSetOrgForOrganizations(orgName, mspid, cryptoPath string, users map[string]string,
	peers, certificateAuthorities []string) {
	c.initOrganizations()
	userMap := map[string]*User{}
	for username, path := range users {
		userMap[username] = &User{Cert: &Cert{Path: path}}
	}
	c.Organizations[orgName] = &Organization{
		MspID:                  mspid,
		CryptoPath:             cryptoPath,
		Users:                  userMap,
		Peers:                  peers,
		CertificateAuthorities: certificateAuthorities,
	}
}

// AddOrSetSelfOrgForOrganizations
//
// peers peer0 peer1
func (c *Config) AddOrSetSelfOrgForOrganizations(leagueName string, peers, certificateAuthorities []string) {
	c.initOrganizations()
	orgName, userName := c.getOrgAndUserName()
	mspid := strings.Join([]string{orgName, "MSP"}, "")
	peerOrg := strings.Join([]string{leagueName, strings.ToLower(orgName)}, "-")
	cryptoPath := strings.Join([]string{
		geneses.CryptoConfigPath(leagueName),
		"/peerOrganizations/", peerOrg, "/users/", userName, "@", peerOrg, "/msp"}, "")
	userCertPath := strings.Join([]string{cryptoPath, "/signcerts/", userName, "@", peerOrg, "-cert.pem"}, "")
	c.Organizations[orgName] = &Organization{
		MspID:                  mspid,
		CryptoPath:             cryptoPath,
		Users:                  map[string]*User{userName: {Cert: &Cert{Path: userCertPath}}},
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

// AddOrSetSelfOrderer
//
// ordererName order0
//
// url grpc://10.10.203.51:30054
func (c *Config) AddOrSetSelfOrderer(leagueName, ordererName, url, keepAliveTime, keepAliveTimeout string,
	keepAlivePermit, failFast, allowInsecure bool) {
	orderTargetName := strings.Join([]string{ordererName, leagueName}, ".")
	orderName := strings.Join([]string{orderTargetName, "7050"}, ":")
	c.initOrderers(orderName)
	tlsCACerts := strings.Join([]string{
		geneses.CryptoConfigPath(leagueName),
		"/ordererOrganizations/", leagueName, "/tlsca/tlsca.", leagueName, "-cert.pem"}, "")
	c.Orderers[orderName] = &Orderer{
		URL: url,
		GRPCOptions: &OrdererGRPCOptions{
			SSLTargetNameOverride: orderTargetName,
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

// AddOrSetSelfPeer
//
// peerName peer0
//
// Url:                   "grpc://10.10.203.51:30056",
//
// EventUrl:              "grpc://10.10.203.51:30058",
func (c *Config) AddOrSetSelfPeer(leagueName, peerName, url, eventUrl, keepAliveTime, keepAliveTimeout string,
	keepAlivePermit, failFast, allowInsecure bool) {
	orgName, _ := c.getOrgAndUserName()
	peerOrg := strings.Join([]string{leagueName, strings.ToLower(orgName)}, "-")
	tlsCACerts := strings.Join([]string{
		geneses.CryptoConfigPath(leagueName),
		"/peerOrganizations/", peerOrg, "/tlsca/tlsca.", peerOrg, "-cert.pem"}, "")
	c.initPeers(peerName)
	c.Peers[peerName] = &Peer{
		URL:      url,
		EventURL: eventUrl,
		GRPCOptions: &PeerGRPCOptions{
			SSLTargetNameOverride: peerName,
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

// AddOrSetSelfCertificateAuthority
//
// url https://10.10.203.51:30059
func (c *Config) AddOrSetSelfCertificateAuthority(leagueName, certName, url, caName, enrollId, enrollSecret string) {
	orgName, _ := c.getOrgAndUserName()
	peerOrg := strings.Join([]string{leagueName, strings.ToLower(orgName)}, "-")
	tlsCACerts := strings.Join([]string{
		geneses.CryptoConfigPath(leagueName),
		"/peerOrganizations/", peerOrg, "/tlsca/tlsca.", peerOrg, "-cert.pem"}, "")
	tlsCACertClientKey := strings.Join([]string{
		geneses.CryptoConfigPath(leagueName),
		"/peerOrganizations/", peerOrg, "/users/Admin@", peerOrg, "/tls/client.key"}, "")
	tlsCACertClientCert := strings.Join([]string{
		geneses.CryptoConfigPath(leagueName),
		"/peerOrganizations/", peerOrg, "/users/Admin@", peerOrg, "/tls/client.crt"}, "")
	c.initCertificateAuthorities(certName)
	c.CertificateAuthorities[certName] = &CertificateAuthority{
		URL: url,
		TLSCACerts: &CertificateAuthorityTLSCACerts{
			Path: tlsCACerts,
			Client: &CertificateAuthorityTLSCACertsClient{
				Key: &CertificateAuthorityTLSCACertsClientKey{
					Path: tlsCACertClientKey,
				},
				Cert: &CertificateAuthorityTLSCACertsClientCert{
					Path: tlsCACertClientCert,
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
		c.Organizations = map[string]*Organization{}
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

func (c *Config) GetPBConfig() *pb.Config {
	return &pb.Config{
		Version:                c.Version,
		Client:                 c.getPBClient(),
		Channels:               c.getPBChannels(),
		Organizations:          c.getPBOrganizations(),
		Orderers:               c.getPBOrders(),
		Peers:                  c.getPBPeers(),
		CertificateAuthorities: c.getPBCertificateAuthorities(),
	}
}

func (c *Config) getPBCertificateAuthorities() map[string]*pb.CertificateAuthority {
	cs := map[string]*pb.CertificateAuthority{}
	for cName, c := range c.CertificateAuthorities {
		cs[cName] = &pb.CertificateAuthority{
			Url:        c.URL,
			CaName:     c.CAName,
			TlsCACerts: &pb.CertificateAuthorityTLSCACerts{
				//Path: c.TLSCACerts.Path,
				//Client: &pb.CertificateAuthorityTLSCACertsClient{
				//	Key: &pb.CertificateAuthorityTLSCACertsClientKey{
				//		Path: c.TLSCACerts.Client.Key.Path,
				//	},
				//	Cert: &pb.CertificateAuthorityTLSCACertsClientCert{
				//		Path: c.TLSCACerts.Client.Cert.Path,
				//	},
				//},
			},
			Registrar: &pb.CertificateAuthorityRegistrar{
				EnrollId:     c.Registrar.EnrollId,
				EnrollSecret: c.Registrar.EnrollSecret,
			},
		}
	}
	return cs
}

func (c *Config) getPBPeers() map[string]*pb.Peer {
	ps := map[string]*pb.Peer{}
	for pName, p := range c.Peers {
		ps[pName] = &pb.Peer{
			Url:      p.URL,
			EventUrl: p.EventURL,
			GrpcOptions: &pb.PeerGRPCOptions{
				SslTargetNameOverride: p.GRPCOptions.SSLTargetNameOverride,
				KeepAliveTime:         p.GRPCOptions.KeepAliveTime,
				KeepAliveTimeout:      p.GRPCOptions.KeepAliveTimeout,
				KeepAlivePermit:       p.GRPCOptions.KeepAlivePermit,
				FailFast:              p.GRPCOptions.FailFast,
				AllowInsecure:         p.GRPCOptions.AllowInsecure,
			},
			TlsCACerts: &pb.PeerTLSCACerts{
				Path: p.TLSCACerts.Path,
			},
		}
	}
	return ps
}

func (c *Config) getPBOrders() map[string]*pb.Orderer {
	os := map[string]*pb.Orderer{}
	for oName, o := range c.Orderers {
		os[oName] = &pb.Orderer{
			Url: o.URL,
			GrpcOptions: &pb.OrdererGRPCOptions{
				SslTargetNameOverride: o.GRPCOptions.SSLTargetNameOverride,
				KeepAliveTime:         o.GRPCOptions.KeepAliveTime,
				KeepAliveTimeout:      o.GRPCOptions.KeepAliveTimeout,
				KeepAlivePermit:       o.GRPCOptions.KeepAlivePermit,
				FailFast:              o.GRPCOptions.FailFast,
				AllowInsecure:         o.GRPCOptions.AllowInsecure,
			},
			TlsCACerts: &pb.OrdererTLSCACerts{
				Path: o.TLSCACerts.Path,
			},
		}
	}
	return os
}

func (c *Config) getPBOrganizations() map[string]*pb.Organization {
	os := map[string]*pb.Organization{}
	for oName, o := range c.Organizations {
		os[oName] = &pb.Organization{
			MspID:                  o.MspID,
			CryptoPath:             o.CryptoPath,
			Peers:                  o.Peers,
			CertificateAuthorities: o.CertificateAuthorities,
		}
	}
	return os
}

func (c *Config) getPBChannels() map[string]*pb.Channel {
	chs := map[string]*pb.Channel{}
	for chName, ch := range c.Channels {
		peers := map[string]*pb.ChannelPeer{}
		for pName, p := range ch.Peers {
			peers[pName] = &pb.ChannelPeer{
				EndorsingPeer:  p.EndorsingPeer,
				ChaincodeQuery: p.ChaincodeQuery,
				EventSource:    p.EventSource,
				LedgerQuery:    p.LedgerQuery,
			}
		}
		chs[chName] = &pb.Channel{
			Peers: peers,
			Policies: &pb.Policy{
				QueryChannelConfig: &pb.PolicyQueryChannelConfig{
					MinResponses: ch.Policies.QueryChannelConfig.MinResponses,
					MaxTargets:   ch.Policies.QueryChannelConfig.MaxTargets,
					RetryOpts: &pb.PolicyCommonRetryOpts{
						Attempts:       ch.Policies.QueryChannelConfig.RetryOpts.Attempts,
						InitialBackoff: ch.Policies.QueryChannelConfig.RetryOpts.InitialBackOff,
						MaxBackoff:     ch.Policies.QueryChannelConfig.RetryOpts.MaxBackOff,
						BackoffFactor:  ch.Policies.QueryChannelConfig.RetryOpts.BackOffFactor,
					},
				},
				Discovery: &pb.PolicyDiscovery{
					MaxTargets: ch.Policies.Discovery.MaxTargets,
					RetryOpts: &pb.PolicyCommonRetryOpts{
						Attempts:       ch.Policies.Discovery.RetryOpts.Attempts,
						InitialBackoff: ch.Policies.Discovery.RetryOpts.InitialBackOff,
						MaxBackoff:     ch.Policies.Discovery.RetryOpts.MaxBackOff,
						BackoffFactor:  ch.Policies.Discovery.RetryOpts.BackOffFactor,
					},
				},
				EventService: &pb.PolicyEventService{
					ResolverStrategy:                 ch.Policies.EventService.ResolverStrategy,
					Balancer:                         ch.Policies.EventService.Balancer,
					BlockHeightLagThreshold:          ch.Policies.EventService.BlockHeightLagThreshold,
					ReconnectBlockHeightLagThreshold: ch.Policies.EventService.ReconnectBlockHeightLagThreshold,
					PeerMonitorPeriod:                ch.Policies.EventService.PeerMonitorPeriod,
				},
			},
		}
	}
	return chs
}

func (c *Config) getPBClient() *pb.Client {
	return &pb.Client{
		Organization: c.Client.Organization,
		Logging: &pb.ClientLogging{
			Level: c.Client.Logging.Level,
		},
		Peer: &pb.ClientPeer{
			Timeout: &pb.ClientPeerTimeout{
				Connection: c.Client.Peer.Timeout.Connection,
				Response:   c.Client.Peer.Timeout.Response,
				Discovery: &pb.ClientPeerTimeoutDiscovery{
					GreyListExpiry: c.Client.Peer.Timeout.Discovery.GreyListExpiry,
				},
			},
		},
		EventService: &pb.ClientEventService{
			Timeout: &pb.ClientEventServiceTimeout{
				RegistrationResponse: c.Client.EventService.Timeout.RegistrationResponse,
			},
		},
		Order: &pb.ClientOrder{
			Timeout: &pb.ClientOrderTimeout{
				Connection: c.Client.Order.Timeout.Connection,
				Response:   c.Client.Order.Timeout.Response,
			},
		},
		Global: &pb.ClientGlobal{
			Timeout: &pb.ClientGlobalTimeout{
				Query:   c.Client.Global.Timeout.Query,
				Execute: c.Client.Global.Timeout.Execute,
				Resmgmt: c.Client.Global.Timeout.Resmgmt,
			},
			Cache: &pb.ClientGlobalCache{
				ConnectionIdle:    c.Client.Global.Cache.ConnectionIdle,
				EventServiceIdle:  c.Client.Global.Cache.EventServiceIdle,
				ChannelMembership: c.Client.Global.Cache.ChannelMembership,
				ChannelConfig:     c.Client.Global.Cache.ChannelConfig,
				Discovery:         c.Client.Global.Cache.Discovery,
				Selection:         c.Client.Global.Cache.Selection,
			},
		},
		CryptoConfig: &pb.ClientCryptoConfig{
			Path: c.Client.CryptoConfig.Path,
		},
		CredentialStore: &pb.ClientCredentialStore{
			Path: c.Client.CredentialStore.Path,
			CryptoStore: &pb.ClientCredentialStoreCryptoStore{
				Path: c.Client.CredentialStore.CryptoStore.Path,
			},
		},
		BCCSP: &pb.ClientBCCSP{
			Security: &pb.ClientBCCSPSecurity{
				Enabled: c.Client.BCCSP.Security.Enabled,
				Default: &pb.ClientBCCSPSecurityDefault{
					Provider: c.Client.BCCSP.Security.Default.Provider,
				},
				HashAlgorithm: c.Client.BCCSP.Security.HashAlgorithm,
				SoftVerify:    c.Client.BCCSP.Security.SoftVerify,
				Level:         c.Client.BCCSP.Security.Level,
			},
		},
		TlsCerts: &pb.ClientTLSCerts{
			SystemCertPool: c.Client.TLSCerts.SystemCertPool,
			Client: &pb.ClientTLSCertsClient{
				Key: &pb.ClientTLSCertsClientKey{
					Path: c.Client.TLSCerts.Client.Key.Path,
				},
				Cert: &pb.ClientTLSCertsClientCert{
					Path: c.Client.TLSCerts.Client.Cert.Path,
				},
			},
		},
	}
}

func (c *Config) getOrgAndUserName() (orgName, userName string) {
	client := c.Client
	orgName = client.Organization
	userName = strings.Split(strings.Split(client.TLSCerts.Client.Cert.Path, "@")[0], "/users/")[1]
	return
}
