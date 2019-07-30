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

//import (
//	"fmt"
//	"github.com/ennoo/rivet/utils/log"
//	"gopkg.in/yaml.v3"
//	"testing"
//)
//
//func TestConfig(t *testing.T) {
//	var (
//		channelName = "test"
//	)
//	TGetConfig(channelName)
//
//	configData, err := yaml.Marshal(Configs[channelName])
//	if err != nil {
//		log.Self.Debug("client", log.Error(err))
//	}
//	fmt.Printf("--- dump:\n%s\n\n", string(configData))
//}
//
//func TGetConfig(channelName string) {
//	_ = InitClient(&Client{
//		ConfigID:     channelName,
//		TlS:          true,
//		Organization: "Org1",
//		Level:        "debug",
//		CryptoConfig: "/crypto-config",
//		KeyPath:      "/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.key",
//		CertPath:     "/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.crt",
//	})
//	_ = AddOrSetPeerForChannel(&ChannelPeer{
//		ConfigID:       channelName,
//		ChannelName:    "mychannel1",
//		PeerName:       "peer0.org1.example.com",
//		EndorsingPeer:  true,
//		ChainCodeQuery: true,
//		LedgerQuery:    true,
//		EventSource:    true,
//	})
//	_ = AddOrSetPeerForChannel(&ChannelPeer{
//		ConfigID:       channelName,
//		ChannelName:    "mychannel1",
//		PeerName:       "peer1.org1.example.com",
//		EndorsingPeer:  true,
//		ChainCodeQuery: true,
//		LedgerQuery:    true,
//		EventSource:    true,
//	})
//	_ = AddOrSetQueryChannelPolicyForChannel(&ChannelPolicyQuery{
//		ConfigID:       channelName,
//		ChannelName:    "mychannel1",
//		InitialBackOff: "500ms",
//		MaxBackOff:     "5s",
//		MinResponses:   1,
//		MaxTargets:     1,
//		Attempts:       5,
//		BackOffFactor:  2.0,
//	})
//	_ = AddOrSetDiscoveryPolicyForChannel(&ChannelPolicyDiscovery{
//		ConfigID:       channelName,
//		ChannelName:    "mychannel1",
//		InitialBackOff: "500ms",
//		MaxBackOff:     "5s",
//		MaxTargets:     2,
//		Attempts:       4,
//		BackOffFactor:  2.0,
//	})
//	_ = AddOrSetEventServicePolicyForChannel(&ChannelPolicyEvent{
//		ConfigID:                         channelName,
//		ChannelName:                      "mychannel1",
//		ResolverStrategy:                 "PreferOrg",
//		Balance:                          "Random",
//		PeerMonitorPeriod:                "6s",
//		BlockHeightLagThreshold:          5,
//		ReconnectBlockHeightLagThreshold: 8,
//	})
//	_ = AddOrSetOrdererForOrganizations(&OrganizationsOrder{
//		ConfigID:   channelName,
//		MspID:      "OrdererMSP",
//		CryptoPath: "/fabric/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/msp",
//	})
//	_ = AddOrSetOrgForOrganizations(&OrganizationsOrg{
//		ConfigID:               channelName,
//		OrgName:                "Org1",
//		MspID:                  "Org1MSP",
//		CryptoPath:             "/fabric/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp",
//		Peers:                  []string{"peer0.org1.example.com", "peer1.org1.example.com"},
//		CertificateAuthorities: []string{"ca.org1.example.com"},
//	})
//	_ = AddOrSetOrderer(&Order{
//		ConfigID:              channelName,
//		OrderName:             "orderer0.example.com",
//		URL:                   "grpc://orderer0.example.com:7050",
//		SSLTargetNameOverride: "orderer0.example.com",
//		KeepAliveTime:         "0s",
//		KeepAliveTimeout:      "20s",
//		TLSCACerts:            "ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem",
//	})
//	_ = AddOrSetOrderer(&Order{
//		ConfigID:              channelName,
//		OrderName:             "orderer1.example.com",
//		URL:                   "grpc://orderer1.example.com:7050",
//		SSLTargetNameOverride: "orderer1.example.com",
//		KeepAliveTime:         "0s",
//		KeepAliveTimeout:      "20s",
//		TLSCACerts:            "ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem",
//	})
//	_ = AddOrSetPeer(&Peer{
//		ConfigID:              channelName,
//		PeerName:              "peer0.org1.example.com",
//		URL:                   "grpc://peer0.org1.example.com:7051",
//		EventUrl:              "grpc://peer0.org1.example.com:7053",
//		SSLTargetNameOverride: "peer0.org1.example.com",
//		KeepAliveTime:         "0s",
//		KeepAliveTimeout:      "20s",
//		TLSCACerts:            "peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem",
//	})
//	_ = AddOrSetPeer(&Peer{
//		ConfigID:              channelName,
//		PeerName:              "peer1.org1.example.com",
//		URL:                   "grpc://peer1.org1.example.com:7051",
//		EventUrl:              "grpc://peer1.org1.example.com:7053",
//		SSLTargetNameOverride: "peer1.org1.example.com",
//		KeepAliveTime:         "0s",
//		KeepAliveTimeout:      "20s",
//		TLSCACerts:            "peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem",
//	})
//	_ = AddOrSetCertificateAuthority(&CertificateAuthority{
//		ConfigID:                channelName,
//		CertName:                "ca.org1.example.com",
//		URL:                     "https://ca.org1.example.com:7054",
//		TLSCACertPath:           "peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem",
//		TLSCACertClientKeyPath:  "peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.key",
//		TLSCACertClientCertPath: "peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.crt",
//		CAName:                  "admin",
//		EnrollId:                "adminpw",
//		EnrollSecret:            "ca.org1.example.com",
//	})
//	_ = AddOrSetCertificateAuthority(&CertificateAuthority{
//		ConfigID:                channelName,
//		CertName:                "ca.org2.example.com",
//		URL:                     "https://ca.org2.example.com:7054",
//		TLSCACertPath:           "peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem",
//		TLSCACertClientKeyPath:  "peerOrganizations/org2.example.com/users/User1@org2.example.com/tls/client.key",
//		TLSCACertClientCertPath: "peerOrganizations/org2.example.com/users/User1@org2.example.com/tls/client.crt",
//		CAName:                  "admin",
//		EnrollId:                "adminpw",
//		EnrollSecret:            "ca.org2.example.com",
//	})
//}
