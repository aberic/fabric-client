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
	"github.com/ennoo/fabric-client/config"
	"github.com/ennoo/fabric-client/service"
	"github.com/ennoo/rivet/utils/log"
	"gopkg.in/yaml.v2"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	fmt.Printf("--- dump:\n%s\n\n", string(confData))

	service.Configs["test"] = conf
	t.Log(get("test", "mychannel"))
}

func TestCreate(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Create("OrdererOrg", "Admin", "grpc://10.10.203.51:30054",
		"Org1", "Admin", "mychannel10",
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fabric-client/example/config/channel-artifacts/mychannel10.tx",
		confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestJoin(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Join("grpc://10.10.203.51:30054", "Org1", "Admin", "mychannel2", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestChannels(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Channels("Org1", "Admin", "peer0", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerInfo(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerInfo("test", "mychannel", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerBlockByHeight(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerBlockByHeight("test", "mychannel", 17, confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerBlockByHash(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerBlockByHash("test", "mychannel", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerBlockByTxID(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerBlockByTxID("test", "mychannel", "b3712eef661af9dbd5b4144e8e6d5b106dd0cb4c1f68f3203749b6c73b04f2f6", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerTransaction(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerTransaction("test", "mychannel", "b3712eef661af9dbd5b4144e8e6d5b106dd0cb4c1f68f3203749b6c73b04f2f6", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerConfig("test", "mychannel", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerInfoSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryLedgerInfoSpec("mychannel", "Org1", "Admin", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerBlockByHeightSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerBlockByHeightSpec("mychannel", "Org1", "Admin", 2, confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerBlockByHashSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryLedgerBlockByHashSpec("mychannel", "Org1", "Admin", "19dce7325781ed8dc022348ee08aa7edb274a91d4d30981b886992704a25b2d4", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerBlockByTxIDSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryLedgerBlockByTxIDSpec("mychannel", "Org1", "Admin", "b3712eef661af9dbd5b4144e8e6d5b106dd0cb4c1f68f3203749b6c73b04f2f6", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerTransactionSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryLedgerTransactionSpec("mychannel", "Org1", "Admin", "b3712eef661af9dbd5b4144e8e6d5b106dd0cb4c1f68f3203749b6c73b04f2f6", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerConfigSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryLedgerConfigSpec("mychannel", "Org1", "Admin", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryChannelInfo(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryChannelInfo("mychannel", "Org1", "Admin", "peer0.league01-org1-vh-cn", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryChannelBlockByHeight(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryChannelBlockByHeight("mychannel", "Org1", "Admin", "peer0.league01-org1-vh-cn", 2, confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryChannelBlockByHash(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryChannelBlockByHash("mychannel", "Org1", "Admin", "peer0.league01-org1-vh-cn", "", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryChannelBlockByTxID(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryChannelBlockByTxID("mychannel", "Org1", "Admin", "peer0.league01-org1-vh-cn", "", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryChannelTransaction(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryChannelTransaction("mychannel", "Org1", "Admin", "peer0.league01-org1-vh-cn", "", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestInstall(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Install("Org1", "Admin", "medical",
		"/Users/aberic/Documents/path/go", "viewhigh.com/dams/chaincode/medical", "1.1",
		confData)
	log.Self.Debug("test install", log.Reflect("result", result))
}

func TestInstalled(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Installed("Org1", "Admin", "peer0", confData)
	log.Self.Debug("test installed", log.Reflect("result", result))
}

func TestInstantiate(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Instantiate("Org1", "Admin", "mychannel", "medical",
		"viewhigh.com/dams/chaincode/medical", "1.0", []string{"Org1MSP", "Org2MSP", "Org3MSP"},
		[][]byte{[]byte("init"), []byte("A"), []byte("10000"), []byte("B"), []byte("10000")}, confData)
	log.Self.Debug("test instantiate", log.Reflect("result", result))
}

func TestUpgrade(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Upgrade("Org1", "Admin", "mychannel", "medical",
		"viewhigh.com/dams/chaincode/medical", "1.1", []string{"Org1MSP", "Org2MSP", "Org3MSP"},
		[][]byte{[]byte("init"), []byte("A"), []byte("10000"), []byte("B"), []byte("10000")}, confData)
	log.Self.Debug("test upgrade", log.Reflect("result", result))
}

func TestInstantiated(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Instantiated("Org1", "Admin", "mychannel", "peer1", confData)
	log.Self.Debug("test instantiated", log.Reflect("result", result))
}

func TestInvoke(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Invoke("medical", "Org1", "Admin", "mychannel",
		"invoke", [][]byte{[]byte("A"), []byte("B"), []byte("1")}, []string{}, confData)
	log.Self.Debug("test invoke", log.Reflect("result", result))
}

func TestInvokeAsync(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := InvokeAsync("medical", "Org1", "Admin", "mychannel", "http://localhost:8082/rivet/post",
		"invoke", [][]byte{[]byte("A"), []byte("B"), []byte("1")}, []string{"peer1"}, confData)
	log.Self.Debug("test invoke", log.Reflect("result", result))
	time.Sleep(time.Second * 60)
}

func TestQuery(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Query("medical", "Org1", "Admin", "mychannel",
		"query", [][]byte{[]byte("A")}, []string{"peer0"}, confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryCollectionsConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryCollectionsConfig("medical", "Org1", "Admin", "mychannel",
		"peer1", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestDiscoveryClientPeers(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := DiscoveryClientPeers("mychannel", "Org1", "Admin", "peer1", confData)
	log.Self.Debug("test discovery client peers", log.Reflect("result", result))
}

func TestDiscoveryClientLocalPeers(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := DiscoveryClientLocalPeers("Org1", "Admin", "peer1", confData)
	log.Self.Debug("test discovery client local peers", log.Reflect("result", result))
}

func TestDiscoveryClientConfigPeers(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := DiscoveryClientConfigPeers("mychannel", "Org1", "Admin", "peer1", confData)
	log.Self.Debug("test discovery client config peers", log.Reflect("result", result))
}

func TestDiscoveryClientEndorsersPeers(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := DiscoveryClientEndorsersPeers("mychannel", "Org1", "Admin", "peer1", "care", confData)
	log.Self.Debug("test discovery client endorsers peers", log.Reflect("result", result))
}

func TestOrderConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := OrderConfig("Org1", "Admin", "mychannel", "grpc://10.10.203.51:30054", confData)
	log.Self.Debug("test order config", log.Reflect("result", result))
}

func TGetConfig() *config.Config {
	rootPath := "/Users/aberic/Documents/path/go/src/github.com/ennoo/fabric-client/example"
	//rootPath := "/Users/admin/Documents/code/git/go/src/github.com/ennoo/fabric-client/example"
	conf := config.Config{}
	_ = conf.InitClient(false, "Org1", "debug",
		rootPath+"/config/crypto-config",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.key",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/tls/client.crt")
	conf.AddOrSetPeerForChannel("mychannel", "peer0",
		true, true, true, true)
	conf.AddOrSetPeerForChannel("mychannel", "peer1",
		true, true, true, true)
	conf.AddOrSetQueryChannelPolicyForChannel("mychannel", "500ms", "5s",
		1, 1, 5, 2.0)
	conf.AddOrSetDiscoveryPolicyForChannel("mychannel", "500ms", "5s",
		2, 4, 2.0)
	conf.AddOrSetEventServicePolicyForChannel("mychannel", "PreferOrg", "RoundRobin",
		"6s", 5, 8)
	conf.AddOrSetOrdererForOrganizations("OrdererMSP",
		rootPath+"/config/crypto-config/ordererOrganizations/league01-vh-cn/users/Admin@league01-vh-cn/msp",
		map[string]string{
			"Admin": rootPath + "/config/crypto-config/ordererOrganizations/league01-vh-cn/users/Admin@league01-vh-cn/msp/signcerts/Admin@league01-vh-cn-cert.pem",
		},
	)
	conf.AddOrSetOrgForOrganizations("Org1", "Org1MSP",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/msp",
		map[string]string{
			"Admin": rootPath + "/config/crypto-config/peerOrganizations/league01-org1-vh-cn/users/Admin@league01-org1-vh-cn/msp/signcerts/Admin@league01-org1-vh-cn-cert.pem",
		},
		[]string{"peer0", "peer1"},
		[]string{"ca0.league01-org1-vh-cn"},
	)
	conf.AddOrSetOrderer("order0.league01-vh-cn:7050", "grpc://10.10.203.51:30054",
		"order0.league01-vh-cn", "0s", "20s",
		rootPath+"/config/crypto-config/ordererOrganizations/league01-vh-cn/tlsca/tlsca.league01-vh-cn-cert.pem",
		false, false, false)
	conf.AddOrSetPeer("peer0", "grpc://10.10.203.51:30056",
		"grpc://10.10.203.51:30058", "peer0",
		"0s", "20s",
		rootPath+"/config/crypto-config/peerOrganizations/league01-org1-vh-cn/tlsca/tlsca.league01-org1-vh-cn-cert.pem",
		false, false, false)
	conf.AddOrSetPeer("peer1", "grpc://10.10.203.51:30061",
		"grpc://10.10.203.51:30063", "peer1",
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
