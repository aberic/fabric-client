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

package chains

import (
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	"github.com/ennoo/rivet/utils/log"
	"gopkg.in/yaml.v3"
	"testing"
)

const (
	rootPath    = "/data"
	configID    = "7ac2ffedd58aa9219038386511b2b801"
	channelName = "mychannel"
	peerName    = "peer0"
)

func TestGetConfig(t *testing.T) {
	i, err := GetConfig("10.10.203.51:31247", &pb.ReqConfig{
		ConfigID: configID,
	})
	t.Log(err)
	t.Log(i)
	configData, err := yaml.Marshal(i)
	if err != nil {
		t.Log("client", log.Error(err))
	}
	t.Log(string(configData))
}

func TestInitClient(t *testing.T) {
	t.Log(InitClient("10.10.203.51:32262", &pb.ReqClient{
		ConfigID:     configID,
		Tls:          false,
		Organization: "Org1",
		Level:        "debug",
		CryptoConfig: rootPath + "/config/crypto-config",
		KeyPath:      rootPath + "/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.key",
		CertPath:     rootPath + "/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.crt",
	}))
}

func TestAddOrSetPeerForChannel(t *testing.T) {
	t.Log(AddOrSetPeerForChannel("127.0.0.1:19878", &pb.ReqChannelPeer{
		ConfigID:       configID,
		ChannelName:    channelName,
		PeerName:       peerName,
		EndorsingPeer:  true,
		ChainCodeQuery: true,
		LedgerQuery:    true,
		EventSource:    true,
	}))
}

func TestAddOrSetQueryChannelPolicyForChannel(t *testing.T) {
	t.Log(AddOrSetQueryChannelPolicyForChannel("127.0.0.1:19878", &pb.ReqChannelPolicyQuery{
		ConfigID:       configID,
		ChannelName:    channelName,
		InitialBackOff: "500ms",
		MaxBackOff:     "5s",
		MinResponses:   1,
		MaxTargets:     1,
		Attempts:       5,
		BackOffFactor:  2.0,
	}))
}

func TestAddOrSetDiscoveryPolicyForChannel(t *testing.T) {
	t.Log(AddOrSetDiscoveryPolicyForChannel("127.0.0.1:19877", &pb.ReqChannelPolicyDiscovery{
		ConfigID:       configID,
		ChannelName:    channelName,
		InitialBackOff: "500ms",
		MaxBackOff:     "5s",
		MaxTargets:     1,
		Attempts:       5,
		BackOffFactor:  2.0,
	}))
}

func TestAddOrSetEventServicePolicyForChannel(t *testing.T) {
	t.Log(AddOrSetEventServicePolicyForChannel("127.0.0.1:19877", &pb.ReqChannelPolicyEvent{
		ConfigID:                         configID,
		ChannelName:                      channelName,
		ResolverStrategy:                 "PreferOrg",
		Balance:                          "Random",
		PeerMonitorPeriod:                "6s",
		BlockHeightLagThreshold:          5,
		ReconnectBlockHeightLagThreshold: 8,
	}))
}

func TestAddOrSetOrdererForOrganizations(t *testing.T) {
	t.Log(AddOrSetOrdererForOrganizations("127.0.0.1:19877", &pb.ReqOrganizationsOrder{
		ConfigID:   configID,
		MspID:      "OrdererMSP",
		CryptoPath: rootPath + "/config/crypto-config/ordererOrganizations/league01-vh-cn/users/Admin@league01-vh-cn/msp",
		Users: map[string]string{
			"Admin": rootPath + "/config/crypto-config/ordererOrganizations/league01-vh-cn/users/Admin@league01-vh-cn/msp/signcerts/Admin@league01-vh-cn-cert.pem",
		},
	}))
}

func TestAddOrSetOrdererForOrganizationsSelf(t *testing.T) {
	t.Log(AddOrSetOrdererForOrganizationsSelf("127.0.0.1:19877", &pb.ReqOrganizationsOrderSelf{
		ConfigID:   configID,
		LeagueName: configID,
	}))
}

func TestAddOrSetOrgForOrganizations(t *testing.T) {
	t.Log(AddOrSetOrgForOrganizations("127.0.0.1:19877", &pb.ReqOrganizationsOrg{
		ConfigID:   configID,
		OrgName:    "Org1",
		MspID:      "Org1MSP",
		CryptoPath: rootPath + "/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/msp",
		Users: map[string]string{
			"Admin": rootPath + "/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/msp/signcerts/Admin@league01-org1-vh-cn-cert.pem",
		},
		Peers:                  []string{"peer0", "peer1"},
		CertificateAuthorities: []string{"ca0.league01-org1-vh-cn"},
	}))
}

func TestAddOrSetOrderer(t *testing.T) {
	t.Log(AddOrSetOrderer("127.0.0.1:19877", &pb.ReqOrder{
		ConfigID:              configID,
		OrderName:             "order0.league01-vh-cn:7050",
		Url:                   "grpc://10.10.203.51:30054",
		SslTargetNameOverride: "order0.league01-vh-cn",
		KeepAliveTime:         "0s",
		KeepAliveTimeout:      "20s",
		TlsCACerts:            rootPath + "/config/crypto-config/ordererOrganizations/league01-vh-cn/tlsca/tlsca.league01-vh-cn-cert.pem",
		KeepAlivePermit:       false,
		FailFast:              false,
		AllowInsecure:         false,
	}))
}

func TestAddOrSetPeer(t *testing.T) {
	t.Log(AddOrSetPeer("127.0.0.1:19877", &pb.ReqPeer{
		ConfigID:              configID,
		PeerName:              peerName,
		Url:                   "grpc://10.10.203.51:30056",
		EventUrl:              "grpc://10.10.203.51:30058",
		SslTargetNameOverride: "peer0",
		KeepAliveTime:         "0s",
		KeepAliveTimeout:      "20s",
		TlsCACerts:            rootPath + "/config/crypto-config/peerOrganizations/league01-org1-vh-cn/tlsca/tlsca.league01-org1-vh-cn-cert.pem",
		KeepAlivePermit:       false,
		FailFast:              false,
		AllowInsecure:         false,
	}))
}

func TestAddOrSetCertificateAuthority(t *testing.T) {
	t.Log(AddOrSetCertificateAuthority("127.0.0.1:19877", &pb.ReqCertificateAuthority{
		ConfigID:                configID,
		CertName:                "ca.league01-vh-cn",
		Url:                     "https://10.10.203.51:30059",
		TlsCACertPath:           rootPath + "/config/crypto-config/peerOrganizations/league01-org1-vh-cn/tlsca/tlsca.league01-org1-vh-cn-cert.pem",
		TlsCACertClientKeyPath:  rootPath + "/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.key",
		TlsCACertClientCertPath: rootPath + "/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.crt",
		CaName:                  "ca.league01-vh-cn",
		EnrollId:                "admin",
		EnrollSecret:            "adminpw",
	}))
}
