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

package sdk

import (
	"fmt"
	"github.com/ennoo/fom/fabric/config"
	"github.com/ennoo/rivet/utils/log"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestQueryYaml(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	fmt.Printf("--- dump:\n%s\n\n", string(confData))
	result := Query("care", "Org1", "Admin", "mychannel",
		"query", [][]byte{[]byte("A")}, []string{},
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fom/config_e2e_code.yaml")
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryRaw(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	fmt.Printf("--- dump:\n%s\n\n", string(confData))
	result := QueryRaw("care", "Org1", "Admin", "mychannel",
		"query", [][]byte{[]byte("A")}, []string{}, confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TGetConfig() *config.Config {
	conf := config.Config{}
	conf.InitClient(false, "Org1", "debug",
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fom/fabric/example/config/crypto-config",
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fom/fabric/example/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.key",
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fom/fabric/example/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.crt")
	conf.AddOrSetPeerForChannel("mychannel", "peer0.league01-org1-vh-cn",
		true, true, true, true)
	conf.AddOrSetQueryChannelPolicyForChannel("mychannel", "500ms", "5s",
		1, 1, 5, 2.0)
	conf.AddOrSetDiscoveryPolicyForChannel("mychannel", "500ms", "5s",
		2, 4, 2.0)
	conf.AddOrSetEventServicePolicyForChannel("mychannel", "PreferOrg", "Random",
		"6s", 5, 8)
	conf.AddOrSetOrdererForOrganizations("OrdererMSP",
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fom/fabric/example/config/crypto-config/ordererOrganizations/league01-vh-cn/users/Admin@league01-vh-cn/msp")
	conf.AddOrSetOrgForOrganizations("Org1", "Org1MSP",
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fom/fabric/example/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/msp",
		[]string{"peer0.league01-org1-vh-cn"},
		[]string{"ca0.league01-org1-vh-cn"},
	)
	conf.AddOrSetOrderer("order0.league01-vh-cn:7050", "grpc://10.10.203.51:30054",
		"order0.league01-vh-cn", "0s", "20s",
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fom/fabric/example/config/crypto-config/ordererOrganizations/league01-vh-cn/tlsca/tlsca.league01-vh-cn-cert.pem",
		false, false, false)
	conf.AddOrSetPeer("peer0.league01-org1-vh-cn", "grpc://10.10.203.51:30056",
		"grpc://10.10.203.51:30058", "peer0.league01-org1-vh-cn",
		"0s", "20s",
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fom/fabric/example/config/crypto-config/peerOrganizations/league01-org1-vh-cn/tlsca/tlsca.league01-org1-vh-cn-cert.pem",
		false, false, false)
	conf.AddOrSetCertificateAuthority("ca.league01-vh-cn", "https://10.10.203.51:30059",
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fom/fabric/example/config/crypto-config/peerOrganizations/league01-org1-vh-cn/tlsca/tlsca.league01-org1-vh-cn-cert.pem",
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fom/fabric/example/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.key",
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fom/fabric/example/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.crt",
		"admin", "adminpw", "ca.league01-vh-cn")
	return &conf
}
