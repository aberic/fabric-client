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

package sdk

import (
	"fmt"
	"github.com/aberic/fabric-client/config"
	"github.com/aberic/fabric-client/service"
	"gopkg.in/yaml.v3"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	fmt.Printf("--- dump:\n%s\n\n", string(confData))

	service.Configs["test"] = conf
	t.Log(get("test", "cc6519b67c4177fc11"))
}

func TestCreate(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result, err := Create("OrdererOrg", "Admin", "grpc://10.10.203.51:30054",
		"Org1", "Admin", "cc6519b67c4177fc1110",
		"/Users/aberic/Documents/path/go/src/github.com/aberic/fabric-client/example/config/channel-artifacts/cc6519b67c4177fc1110.tx",
		confData)
	t.Log(result)
}

func TestJoin(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := Join("grpc://10.10.203.51:30054", "Org1", "Admin", "cc6519b67c4177fc112", "peer0", confData)
	t.Log(result)
}

func TestChannels(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result, err := Channels("Org1", "Admin", "peer1", confData)
	t.Log(result)
}

func TestQueryLedgerInfo(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	service.Configs["test"] = conf
	result := QueryLedgerInfo("test", "peer0", "cc6519b67c4177fc11", confData)
	t.Log(result)
}

func TestQueryLedgerBlockByHeight(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	service.Configs["test"] = conf
	result := QueryLedgerBlockByHeight("test", "peer0", "cc6519b67c4177fc11", 3, confData)
	t.Log(result)
}

func TestQueryLedgerBlockByHash(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	service.Configs["test"] = conf
	result := QueryLedgerBlockByHash("test", "peer0", "cc6519b67c4177fc11", "b949429f98d25bf58cb242b215aaf868662a6309489e5583663247ce522f2fc6", confData)
	t.Log(result)
}

func TestQueryLedgerBlockByTxID(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	service.Configs["test"] = conf
	result := QueryLedgerBlockByTxID("test", "peer0", "cc6519b67c4177fc11", "9f1090e9d1fc45f53c16394420db24b8ea2225f2e9c33717d9cf9004e31c74c4", confData)
	t.Log(result)
}

func TestQueryLedgerTransaction(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	service.Configs["test"] = conf
	result := QueryLedgerTransaction("test", "peer0", "cc6519b67c4177fc11", "9f1090e9d1fc45f53c16394420db24b8ea2225f2e9c33717d9cf9004e31c74c4", confData)
	t.Log(result)
}

func TestQueryLedgerConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	service.Configs["test"] = conf
	result := QueryLedgerConfig("test", "peer0", "cc6519b67c4177fc11", confData)
	t.Log(result)
}

func TestQueryLedgerInfoSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := QueryLedgerInfoSpec("peer0", "cc6519b67c4177fc11", "Org1", "Admin", confData)
	t.Log(result)
}

func TestQueryLedgerBlockByHeightSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	service.Configs["test"] = conf
	result := QueryLedgerBlockByHeightSpec("peer0", "cc6519b67c4177fc11", "Org1", "Admin", 2, confData)
	t.Log(result)
}

func TestQueryLedgerBlockByHashSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := QueryLedgerBlockByHashSpec("peer0", "cc6519b67c4177fc11", "Org1", "Admin", "19dce7325781ed8dc022348ee08aa7edb274a91d4d30981b886992704a25b2d4", confData)
	t.Log(result)
}

func TestQueryLedgerBlockByTxIDSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := QueryLedgerBlockByTxIDSpec("peer0", "cc6519b67c4177fc11", "Org1", "Admin", "b3712eef661af9dbd5b4144e8e6d5b106dd0cb4c1f68f3203749b6c73b04f2f6", confData)
	t.Log(result)
}

func TestQueryLedgerTransactionSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := QueryLedgerTransactionSpec("peer0", "cc6519b67c4177fc11", "Org1", "Admin", "b3712eef661af9dbd5b4144e8e6d5b106dd0cb4c1f68f3203749b6c73b04f2f6", confData)
	t.Log(result)
}

func TestQueryLedgerConfigSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := QueryLedgerConfigSpec("peer0", "cc6519b67c4177fc11", "Org1", "Admin", confData)
	t.Log(result)
}

func TestQueryChannelInfo(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := QueryChannelInfo("cc6519b67c4177fc11", "Org1", "Admin", "peer0.20de78630ef6a411-org1", confData)
	t.Log(result)
}

func TestQueryChannelBlockByHeight(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := QueryChannelBlockByHeight("cc6519b67c4177fc11", "Org1", "Admin", "peer0.20de78630ef6a411-org1", 2, confData)
	t.Log(result)
}

func TestQueryChannelBlockByHash(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := QueryChannelBlockByHash("cc6519b67c4177fc11", "Org1", "Admin", "peer0.20de78630ef6a411-org1", "", confData)
	t.Log(result)
}

func TestQueryChannelBlockByTxID(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := QueryChannelBlockByTxID("cc6519b67c4177fc11", "Org1", "Admin", "peer0.20de78630ef6a411-org1", "", confData)
	t.Log(result)
}

func TestQueryChannelTransaction(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := QueryChannelTransaction("cc6519b67c4177fc11", "Org1", "Admin", "peer0.20de78630ef6a411-org1", "", confData)
	t.Log(result)
}

func TestQueryConfigBlock(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := QueryConfigBlock("cc6519b67c4177fc11", "Org1", "Admin", "peer0.20de78630ef6a411-org1", confData)
	t.Log(result)
}

func TestInstall(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := Install("Org1", "Admin", "peer0", "medical",
		"/Users/aberic/Documents/path/go", "viewhigh.com/dams/chaincode/medical", "1.1",
		confData)
	t.Log(result)
}

func TestInstalled(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := Installed("Org1", "Admin", "peer0", confData)
	t.Log(result)
}

func TestInstantiate(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := Instantiate("Org1", "Admin", "peer0", "cc6519b67c4177fc11", "medical",
		"viewhigh.com/dams/chaincode/medical", "1.0", []string{"Org1MSP", "Org2MSP", "Org3MSP"},
		[][]byte{[]byte("init"), []byte("A"), []byte("10000"), []byte("B"), []byte("10000")}, confData)
	t.Log(result)
}

func TestUpgrade(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := Upgrade("Org1", "Admin", "peer0", "cc6519b67c4177fc11", "medical",
		"viewhigh.com/dams/chaincode/medical", "1.1", []string{"Org1MSP", "Org2MSP", "Org3MSP"},
		[][]byte{[]byte("init"), []byte("A"), []byte("10000"), []byte("B"), []byte("10000")}, confData)
	t.Log(result)
}

func TestInstantiated(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := Instantiated("Org1", "Admin", "cc6519b67c4177fc11", "peer0", confData)
	t.Log(result)
}

func TestInvoke(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := Invoke("medical", "Org1", "Admin", "cc6519b67c4177fc11",
		"invoke", [][]byte{[]byte("A"), []byte("B"), []byte("1")}, []string{}, confData)
	t.Log(result)
}

func TestInvokeAsync(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := InvokeAsync("medical", "Org1", "Admin", "cc6519b67c4177fc11", "http://localhost:8082/rivet/post",
		"invoke", [][]byte{[]byte("A"), []byte("B"), []byte("1")}, []string{"peer1"}, confData)
	t.Log(result)
	time.Sleep(time.Second * 60)
}

func TestQuery(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := Query("medical", "Org1", "Admin", "cc6519b67c4177fc11",
		"query", [][]byte{[]byte("A")}, []string{"peer0"}, confData)
	t.Log(result)
}

func TestQueryCollectionsConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := QueryCollectionsConfig("medical", "Org1", "Admin", "cc6519b67c4177fc11",
		"peer1", confData)
	t.Log(result)
}

func TestDiscoveryChannelPeers(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	peers, err := DiscoveryChannelPeers("cc6519b67c4177fc11", "Org1", "Admin", confData)
	t.Log(peers, err)
}

func TestDiscoveryLocalPeers(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	peers, err := DiscoveryLocalPeers("Org1", "Admin", confData)
	t.Log(peers, err)
}

func TestOrderConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Log(err)
	}
	result := OrderConfig("Org1", "Admin", "cc6519b67c4177fc11", "grpc://10.10.203.51:30054", confData)
	t.Log(result)
}

func TGetConfig() *config.Config {
	rootPath := "/Users/aberic/Documents/path/go/src/github.com/aberic/fabric-client/example"
	//rootPath := "/Users/admin/Documents/code/git/go/src/github.com/aberic/fabric-client/example"
	conf := config.Config{}
	conf.InitClient(true, "Org1", "debug",
		rootPath+"/config/crypto-config",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.key",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.crt")
	//conf.AddOrSetPeerForChannel("cc6519b67c4177fc11", "peer0",
	//	true, true, true, true)
	conf.AddOrSetPeerForChannel("cc6519b67c4177fc11", "peer1",
		true, true, true, true)
	conf.AddOrSetQueryChannelPolicyForChannel("cc6519b67c4177fc11", "500ms", "5s",
		1, 1, 5, 2.0)
	conf.AddOrSetDiscoveryPolicyForChannel("cc6519b67c4177fc11", "500ms", "5s",
		2, 4, 2.0)
	conf.AddOrSetEventServicePolicyForChannel("cc6519b67c4177fc11", "PreferOrg", "RoundRobin",
		"6s", 5, 8)
	conf.AddOrSetOrdererForOrganizations(config.OrderOrgKey, "OrdererMSP",
		rootPath+"/config/crypto-config/ordererOrganizations/20de78630ef6a411/users/Admin@20de78630ef6a411/msp",
		map[string]string{
			"Admin": rootPath + "/config/crypto-config/ordererOrganizations/20de78630ef6a411/users/Admin@20de78630ef6a411/msp/signcerts/Admin@20de78630ef6a411-cert.pem",
		},
	)
	conf.AddOrSetOrgForOrganizations("Org1", "Org1MSP",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/msp",
		map[string]string{
			"Admin": rootPath + "/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/msp/signcerts/Admin@20de78630ef6a411-org1-cert.pem",
		},
		[]string{"peer0", "peer1"},
		[]string{"ca"},
	)
	conf.AddOrSetOrderer("orderer0.20de78630ef6a411:7050", "grpcs://10.10.203.51:30054",
		"orderer0.20de78630ef6a411", "0s", "20s",
		rootPath+"/config/crypto-config/ordererOrganizations/20de78630ef6a411/tlsca/tlsca.20de78630ef6a411-cert.pem",
		false, false, false)
	conf.AddOrSetPeer("peer0", "grpcs://10.10.203.51:32625",
		"grpcs://10.10.203.51:30386", "peer0",
		"0s", "20s",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/tlsca/tlsca.20de78630ef6a411-org1-cert.pem",
		false, false, false)
	conf.AddOrSetPeer("peer1", "grpcs://10.10.203.51:32707",
		"grpcs://10.10.203.51:32636", "peer1",
		"0s", "20s",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/tlsca/tlsca.20de78630ef6a411-org1-cert.pem",
		false, false, false)
	conf.AddOrSetCertificateAuthority("ca", "https://10.10.203.51:31906",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/tlsca/tlsca.20de78630ef6a411-org1-cert.pem",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.key",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.crt",
		"ca.20de78630ef6a411-org1", "admin", "adminpw")
	return &conf
}
