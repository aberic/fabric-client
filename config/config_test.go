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
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
	"testing"
)

func TestNewConfig1(t *testing.T) {
	config := Config{
		Version:                "1.0.0",
		Client:                 TGetClient(),
		Channels:               TGetChannels(),
		Organizations:          TGetOrganizations(),
		Orderers:               TGetOrderers(),
		Peers:                  TGetPeers(),
		CertificateAuthorities: TGetCertificateAuthorities(),
	}
	configData, err := yaml.Marshal(&config)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("--- dump:\n%s\n\n", string(configData))
}

func TestNewConfig2(t *testing.T) {
	config := TGetConfig()
	configData, err := yaml.Marshal(&config)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("--- dump:\n%s\n\n", string(configData))
}

func TestSubString(t *testing.T) {
	s := "league-org1"
	t.Log(strings.Split(s, "-org"))
}

func TGetConfig() *Config {
	config := Config{}
	config.InitClient(true, "Org1", "debug",
		"Users/aberic/Documents/path/go/src/github.com/ennoo/fabric-client/fabric/example/config/crypto-config",
		"/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.key",
		"/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.crt")
	config.AddOrSetPeerForChannel("mychannel1", "peer0.org1.example.com",
		true, true, true, true)
	config.AddOrSetPeerForChannel("mychannel1", "peer1.org1.example.com",
		true, true, true, true)
	config.AddOrSetQueryChannelPolicyForChannel("mychannel1", "500ms", "5s",
		1, 1, 5, 2.0)
	config.AddOrSetDiscoveryPolicyForChannel("mychannel1", "500ms", "5s",
		2, 4, 2.0)
	config.AddOrSetEventServicePolicyForChannel("mychannel1", "PreferOrg", "RoundRobin",
		"6s", 5, 8)
	config.AddOrSetOrdererForOrganizations("OrdererMSP",
		"/fabric/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/msp",
		map[string]string{
			"Admin": "/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/msp/signcerts/Admin@example.com-cert.pem",
		})
	config.AddOrSetOrgForOrganizations("Org1", "Org1MSP",
		"/fabric/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp",
		map[string]string{
			"Admin": "/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem",
			"User1": "/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem",
		},
		[]string{"peer0.org1.example.com", "peer1.org1.example.com"},
		[]string{"ca.org1.example.com"},
	)
	config.AddOrSetOrderer("orderer0.example.com", "grpc://orderer0.example.com:7050",
		"orderer0.example.com", "0s", "20s",
		"ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem",
		false, false, false)
	config.AddOrSetOrderer("orderer1.example.com", "grpc://orderer1.example.com:7050",
		"orderer1.example.com", "0s", "20s",
		"ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem",
		false, false, false)
	config.AddOrSetPeer("peer0.org1.example.com", "grpc://peer0.org1.example.com:7051",
		"grpc://peer0.org1.example.com:7053", "peer0.org1.example.com",
		"0s", "20s",
		"peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem",
		false, false, false)
	config.AddOrSetPeer("peer1.org1.example.com", "grpc://peer1.org1.example.com:7051",
		"grpc://peer1.org1.example.com:7053", "peer1.org1.example.com",
		"0s", "20s",
		"peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem",
		false, false, false)
	config.AddOrSetCertificateAuthority("ca.org1.example.com", "https://ca.org1.example.com:7054",
		"peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem",
		"peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.key",
		"peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.crt",
		"admin", "adminpw", "ca.org1.example.com")
	config.AddOrSetCertificateAuthority("ca.org2.example.com", "https://ca.org2.example.com:7054",
		"peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem",
		"peerOrganizations/org2.example.com/users/User1@org2.example.com/tls/client.key",
		"peerOrganizations/org2.example.com/users/User1@org2.example.com/tls/client.crt",
		"ca.org2.example.com", "admin", "adminpw")
	return &config
}
