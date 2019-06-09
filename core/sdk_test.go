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
	"github.com/ennoo/fabric-go-client/config"
	"github.com/ennoo/rivet/utils/log"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	fmt.Printf("--- dump:\n%s\n\n", string(confData))
}

func TestInstalled(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Installed("Org1", "Admin", "peer0.league01-org1-vh-cn", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestInstantiated(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Instantiated("Org1", "Admin", "mychannel", "peer0.league01-org1-vh-cn", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestChannels(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Channels("Org1", "Admin", "peer0.league01-org1-vh-cn", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestOrderConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := OrderConfig("Org1", "Admin", "mychannel", "grpc://10.10.203.51:30054", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQuery(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Query("care", "Org1", "Admin", "mychannel",
		"query", [][]byte{[]byte("A")}, []string{}, confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryCollectionsConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryCollectionsConfig("care", "Org1", "Admin", "mychannel",
		"peer0.league01-org1-vh-cn", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestDiscoveryService(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := DiscoveryService(confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TGetConfig() *config.Config {
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
		[]string{"peer0.league01-org1-vh-cn"},
		[]string{"ca0.league01-org1-vh-cn"},
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
	conf.AddOrSetCertificateAuthority("ca.league01-vh-cn", "https://10.10.203.51:30059",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/tlsca/tlsca.league01-org1-vh-cn-cert.pem",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.key",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.crt",
		"admin", "adminpw", "ca.league01-vh-cn")
	return &conf
}
