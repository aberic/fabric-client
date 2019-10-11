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
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
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
	t.Log(get("test", "cc6519b67c4177fc11"))
}

func TestCreate(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result, err := Create("OrdererOrg", "Admin", "grpc://10.10.203.51:30054",
		"Org1", "Admin", "cc6519b67c4177fc1110",
		"/Users/aberic/Documents/path/go/src/github.com/ennoo/fabric-client/example/config/channel-artifacts/cc6519b67c4177fc1110.tx",
		confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestJoin(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Join("grpc://10.10.203.51:30054", "Org1", "Admin", "cc6519b67c4177fc112", "peer0", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestChannels(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result, err := Channels("Org1", "Admin", "peer1", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerInfo(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerInfo("test", "peer0", "cc6519b67c4177fc11", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerBlockByHeight(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerBlockByHeight("test", "peer0", "cc6519b67c4177fc11", 3, confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerBlockByHash(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerBlockByHash("test", "peer0", "cc6519b67c4177fc11", "b949429f98d25bf58cb242b215aaf868662a6309489e5583663247ce522f2fc6", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerBlockByTxID(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerBlockByTxID("test", "peer0", "cc6519b67c4177fc11", "9f1090e9d1fc45f53c16394420db24b8ea2225f2e9c33717d9cf9004e31c74c4", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerTransaction(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerTransaction("test", "peer0", "cc6519b67c4177fc11", "9f1090e9d1fc45f53c16394420db24b8ea2225f2e9c33717d9cf9004e31c74c4", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerConfig("test", "peer0", "cc6519b67c4177fc11", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerInfoSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryLedgerInfoSpec("peer0", "cc6519b67c4177fc11", "Org1", "Admin", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerBlockByHeightSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	service.Configs["test"] = conf
	result := QueryLedgerBlockByHeightSpec("peer0", "cc6519b67c4177fc11", "Org1", "Admin", 2, confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerBlockByHashSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryLedgerBlockByHashSpec("peer0", "cc6519b67c4177fc11", "Org1", "Admin", "19dce7325781ed8dc022348ee08aa7edb274a91d4d30981b886992704a25b2d4", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerBlockByTxIDSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryLedgerBlockByTxIDSpec("peer0", "cc6519b67c4177fc11", "Org1", "Admin", "b3712eef661af9dbd5b4144e8e6d5b106dd0cb4c1f68f3203749b6c73b04f2f6", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerTransactionSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryLedgerTransactionSpec("peer0", "cc6519b67c4177fc11", "Org1", "Admin", "b3712eef661af9dbd5b4144e8e6d5b106dd0cb4c1f68f3203749b6c73b04f2f6", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryLedgerConfigSpec(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryLedgerConfigSpec("peer0", "cc6519b67c4177fc11", "Org1", "Admin", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryChannelInfo(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryChannelInfo("cc6519b67c4177fc11", "Org1", "Admin", "peer0.20de78630ef6a411-org1", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryChannelBlockByHeight(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryChannelBlockByHeight("cc6519b67c4177fc11", "Org1", "Admin", "peer0.20de78630ef6a411-org1", 2, confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryChannelBlockByHash(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryChannelBlockByHash("cc6519b67c4177fc11", "Org1", "Admin", "peer0.20de78630ef6a411-org1", "", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryChannelBlockByTxID(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryChannelBlockByTxID("cc6519b67c4177fc11", "Org1", "Admin", "peer0.20de78630ef6a411-org1", "", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryChannelTransaction(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryChannelTransaction("cc6519b67c4177fc11", "Org1", "Admin", "peer0.20de78630ef6a411-org1", "", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestInstall(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Install("Org1", "Admin", "peer0", "medical",
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
	result := Instantiate("Org1", "Admin", "peer0", "cc6519b67c4177fc11", "medical",
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
	result := Upgrade("Org1", "Admin", "peer0", "cc6519b67c4177fc11", "medical",
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
	result := Instantiated("Org1", "Admin", "cc6519b67c4177fc11", "peer0", confData)
	log.Self.Debug("test instantiated", log.Reflect("result", result))
}

func TestInvoke(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := Invoke("medical", "Org1", "Admin", "cc6519b67c4177fc11",
		"invoke", [][]byte{[]byte("A"), []byte("B"), []byte("1")}, []string{}, confData)
	log.Self.Debug("test invoke", log.Reflect("result", result))
}

func TestInvokeAsync(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := InvokeAsync("medical", "Org1", "Admin", "cc6519b67c4177fc11", "http://localhost:8082/rivet/post",
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
	result := Query("medical", "Org1", "Admin", "cc6519b67c4177fc11",
		"query", [][]byte{[]byte("A")}, []string{"peer0"}, confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestQueryCollectionsConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := QueryCollectionsConfig("medical", "Org1", "Admin", "cc6519b67c4177fc11",
		"peer1", confData)
	log.Self.Debug("test query", log.Reflect("result", result))
}

func TestDiscoveryClientPeers(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := DiscoveryClientPeers("cc6519b67c4177fc11", "Org1", "Admin", "peer1", confData)
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
	result := DiscoveryClientConfigPeers("cc6519b67c4177fc11", "Org1", "Admin", "peer1", confData)
	log.Self.Debug("test discovery client config peers", log.Reflect("result", result))
}

func TestDiscoveryClientEndorsersPeers(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := DiscoveryClientEndorsersPeers("cc6519b67c4177fc11", "Org1", "Admin", "peer1", "care", confData)
	log.Self.Debug("test discovery client endorsers peers", log.Reflect("result", result))
}

func TestDiscoveryChannelPeers(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	peers, err := DiscoveryChannelPeers("cc6519b67c4177fc11", "Org1", "Admin", confData)
	log.Self.Debug("test discovery channel peers", log.Reflect("peers", peers), log.Error(err))
}

func TestDiscoveryLocalPeers(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	peers, err := DiscoveryLocalPeers("Org1", "Admin", confData)
	log.Self.Debug("test discovery local peers", log.Reflect("peers", peers), log.Error(err))
}

func TestOrderConfig(t *testing.T) {
	conf := TGetConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Debug("client", log.Error(err))
	}
	result := OrderConfig("Org1", "Admin", "cc6519b67c4177fc11", "grpc://10.10.203.51:30054", confData)
	log.Self.Debug("test order config", log.Reflect("result", result))
}

func TestCAInfo(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		log.Self.Error("client", log.Error(err))
	}
	result, err := CAInfo("league", confData)
	if err != nil {
		log.Self.Error("ca info", log.Error(err))
	}
	log.Self.Debug("test ca info", log.Reflect("result", result))
}

func TestEnroll(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	//attrReqs := []*msp.AttributeRequest{{Name: "test", Optional: true}}
	err = Enroll("league", "admin", confData, []msp.EnrollmentOption{
		msp.WithSecret("adminpw"),
		msp.WithType("x509" /*or idemix, which is not support now*/),
		msp.WithProfile("tls"),
		//msp.WithLabel("ForFabric"),
		//msp.WithAttributeRequests(attrReqs),
	})
	if err != nil {
		t.Error("Enroll", err)
	}
	t.Log("test ca Enroll")
}

func TestReenroll(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	err = Reenroll("league", "admin", confData, []msp.EnrollmentOption{
		msp.WithSecret("adminpw"),
		msp.WithType("x509" /*or idemix, which is not support now*/),
		// msp.WithProfile("tls"),
		// msp.WithLabel("ForFabric"),
		//msp.WithAttributeRequests(attrReqs),
	})
	if err != nil {
		t.Error("Enroll", err)
	}
	t.Log("test ca Enroll")
}

func TestRegister(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := Register("league", &msp.RegistrationRequest{
		Name:           "test2",
		Type:           "client", // (e.g. "client, orderer, peer, app, user")
		MaxEnrollments: -1,       // if omitted, this defaults to max_enrollments configured on the server
		Affiliation:    "org1",   // The identity's affiliation e.g. org1.department1
		Attributes: []msp.Attribute{ // Optional attributes associated with this identity.
			// Attribute defines additional attributes that may be passed along during registration
			{Name: "hf.Registrar.Roles", Value: "client,orderer,peer,user"},
			{Name: "hf.Registrar.DelegateRoles", Value: "client,orderer,peer,user"},
			{Name: "hf.Registrar.Attributes", Value: "*"},
			{Name: "hf.GenCRL", Value: "true"},
			{Name: "hf.Revoker", Value: "true"},
			{Name: "hf.AffiliationMgr", Value: "true"},
			{Name: "hf.IntermediateCA", Value: "true"},
			{Name: "role", Value: "admin", ECert: true},
		},
		CAName: "league",
		Secret: "test2",
	}, confData)
	if err != nil {
		t.Error("Register", err)
	}
	t.Log("test ca Register, result = ", result)
}

func TestAddAffiliation(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := AddAffiliation("league", &msp.AffiliationRequest{
		Name:   "org100.peer1", // Name of the affiliation
		Force:  true,           // Creates parent affiliations if they do not exist
		CAName: "league",       // Name of the CA
	}, confData)
	if err != nil {
		t.Error("GetAllAffiliations", err)
	}
	t.Log("test ca GetAllAffiliations, result = ", result)
}

func TestRemoveAffiliation(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := RemoveAffiliation("league", &msp.AffiliationRequest{
		Name:   "org100.peer1", // Name of the affiliation
		Force:  true,           // Creates parent affiliations if they do not exist
		CAName: "league",       // Name of the CA
	}, confData)
	if err != nil {
		t.Error("GetAllAffiliations", err)
	}
	t.Log("test ca GetAllAffiliations, result = ", result)
}

func TestModifyAffiliation(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := ModifyAffiliation("league", &msp.ModifyAffiliationRequest{
		NewName: "org1.department3",
		AffiliationRequest: msp.AffiliationRequest{
			Name:   "org101.peer2", // Name of the affiliation
			Force:  true,           // Creates parent affiliations if they do not exist
			CAName: "league",       // Name of the CA
		},
	}, confData)
	if err != nil {
		t.Error("GetAllAffiliations", err)
	}
	t.Log("test ca GetAllAffiliations, result = ", result)
}

func TestGetAllAffiliations(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := GetAllAffiliations("league", confData)
	if err != nil {
		t.Error("GetAllAffiliations", err)
	}
	t.Log("test ca GetAllAffiliations, result = ", result)
}

func TestGetAllAffiliationsByCaName(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := GetAllAffiliationsByCaName("league", "league", confData)
	if err != nil {
		t.Error("GetAllAffiliationsByCaName", err)
	}
	t.Log("test ca GetAllAffiliationsByCaName, result = ", result)
}

func TestGetAffiliation(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := GetAffiliation("org100.peer0", "league", confData)
	if err != nil {
		t.Error("GetAffiliation", err)
	}
	t.Log("test ca GetAffiliation, result = ", result)
}

func TestGetAffiliationByCaName(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := GetAffiliationByCaName("org100.peer0", "league", "league", confData)
	if err != nil {
		t.Error("GetAffiliationByCaName", err)
	}
	t.Log("test ca GetAffiliationByCaName, result = ", result)
}

func TestCreateIdentity(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	identity, err := CreateIdentity("league", &msp.IdentityRequest{
		ID:          "test3", // The enrollment ID which uniquely identifies an identity (required)
		Affiliation: "org1",  // The identity's affiliation e.g. org1.department1
		Attributes: []msp.Attribute{ // Optional attributes associated with this identity.
			// Attribute defines additional attributes that may be passed along during registration
			{Name: "hf.Registrar.Roles", Value: "client,orderer,peer,user"},
			{Name: "hf.Registrar.DelegateRoles", Value: "client,orderer,peer,user"},
			{Name: "hf.Registrar.Attributes", Value: "*"},
			{Name: "hf.GenCRL", Value: "true"},
			{Name: "hf.Revoker", Value: "true"},
			{Name: "hf.AffiliationMgr", Value: "true"},
			{Name: "hf.IntermediateCA", Value: "true"},
			{Name: "role", Value: "admin", ECert: true},
		},
		Type:           "client", // (e.g. "client, orderer, peer, app, user")
		MaxEnrollments: -1,       // if omitted, this defaults to max_enrollments configured on the server
		Secret:         "test3",
		CAName:         "league",
	}, confData)
	if err != nil {
		t.Error("CreateIdentity", err)
	}
	t.Log("test ca CreateIdentity, identity = ", identity)
}

func TestModifyIdentity(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	identity, err := ModifyIdentity("league", &msp.IdentityRequest{
		ID:          "test3", // The enrollment ID which uniquely identifies an identity (required)
		Affiliation: "org2",  // The identity's affiliation e.g. org1.department1
		Attributes: []msp.Attribute{ // Optional attributes associated with this identity.
			// Attribute defines additional attributes that may be passed along during registration
			{Name: "hf.Registrar.Roles", Value: "client,orderer,peer,user"},
			{Name: "hf.Registrar.DelegateRoles", Value: "client,orderer,peer,user"},
			{Name: "hf.Registrar.Attributes", Value: "*"},
			{Name: "hf.GenCRL", Value: "true"},
			{Name: "hf.Revoker", Value: "true"},
			{Name: "hf.AffiliationMgr", Value: "true"},
			{Name: "hf.IntermediateCA", Value: "true"},
			{Name: "role", Value: "admin", ECert: true},
		},
		Type:           "peer", // (e.g. "client, orderer, peer, app, user")
		MaxEnrollments: -1,     // if omitted, this defaults to max_enrollments configured on the server
		Secret:         "test3",
		CAName:         "league",
	}, confData)
	if err != nil {
		t.Error("ModifyIdentity", err)
	}
	t.Log("test ca ModifyIdentity, identity = ", identity)
}

func TestRemoveIdentity(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	identity, err := RemoveIdentity("league", &msp.RemoveIdentityRequest{
		ID:     "test3", // The enrollment ID which uniquely identifies an identity (required)
		CAName: "league",
		Force:  true, //  Force delete
	}, confData)
	if err != nil {
		t.Error("RemoveIdentity", err)
	}
	t.Log("test ca RemoveIdentity, identity = ", identity)
}

func TestGetIdentity(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	identity, err := GetIdentity("test1", "league", confData)
	if err != nil {
		t.Error("GetIdentity", err)
	}
	t.Log("test ca GetIdentity, identity = ", identity)
}

func TestGetIdentityByCaName(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	identity, err := GetIdentityByCaName("test2", "league", "league", confData)
	if err != nil {
		t.Error("GetIdentity", err)
	}
	t.Log("test ca GetIdentity, identity = ", identity)
}

func TestGetAllIdentities(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	identities, err := GetAllIdentities("league", confData)
	if err != nil {
		t.Error("GetAllIdentities", err)
	}
	for _, identity := range identities {
		t.Log("test ca GetAllIdentities, identity = ", identity)
	}
}

func TestGetAllIdentitiesByCaName(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	identities, err := GetAllIdentitiesByCaName("league", "league", confData)
	if err != nil {
		t.Error("GetAllIdentities", err)
	}
	for _, identity := range identities {
		t.Log("test ca GetAllIdentities, identity = ", identity)
	}
}

func TestCreateSigningIdentity(t *testing.T) {
	cert := `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgTgaXa53bD8bs4HVD
CRQmrz9y/aXksSglI05MHKNNWWihRANCAARmsno2FlgD+qRtlV7pxUn5YJEkLOBb
WawnrkyK8pCWFDbEPr2h1oof1jmRTLaY2GyraMjp0OSJIoRO+gtErTyP
-----END PRIVATE KEY-----`
	privateKey := `-----BEGIN CERTIFICATE-----
MIICHDCCAcOgAwIBAgIUYQr4HjSiartZqKYYTSy8P6XNRqcwCgYIKoZIzj0EAwIw
cjELMAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxEDAOBgNVBAcTB0JlaWpp
bmcxETAPBgNVBAoTCFZpZXdoaWdoMRMwEQYDVQQLEwpCbG9ja2NoYWluMRcwFQYD
VQQDEw5jYS5sZWFndWU6NzA1NDAeFw0xOTA4MTkwMzQ4MDBaFw0yMDA4MTgwMzUz
MDBaMCExDzANBgNVBAsTBmNsaWVudDEOMAwGA1UEAxMFYWRtaW4wWTATBgcqhkjO
PQIBBggqhkjOPQMBBwNCAARmsno2FlgD+qRtlV7pxUn5YJEkLOBbWawnrkyK8pCW
FDbEPr2h1oof1jmRTLaY2GyraMjp0OSJIoRO+gtErTyPo4GHMIGEMA4GA1UdDwEB
/wQEAwIHgDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBTI3duAPmpSY4kU8zhXm2NX
PShvSjAfBgNVHSMEGDAWgBR0B6m5gg9lN/p7de2M0BKo9TZAajAkBgNVHREEHTAb
ghlBYmVyaWNkZU1hY0Jvb2stUHJvLmxvY2FsMAoGCCqGSM49BAMCA0cAMEQCIAvf
ZUYJid9dzEcOei1+i13S+B2HhUuY0048xnEpANPoAiAkuN5DzQEt/8/4YCq4xjEm
IZk4cTJDwbIjGemfgi6PKg==
-----END CERTIFICATE-----`
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	identity, err := CreateSigningIdentity(
		"league",
		confData,
		[]mspctx.SigningIdentityOption{
			mspctx.WithPrivateKey([]byte(privateKey)),
			mspctx.WithCert([]byte(cert)),
		},
	)
	if err != nil {
		t.Error("CreateSigningIdentity", err)
	}
	if string(identity.EnrollmentCertificate()) != cert {
		t.Log("certificate mismatch\n")
		return
	}
	t.Log("test ca CreateSigningIdentity, identity = ", identity)
}

func TestGetSigningIdentity(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	identity, err := GetSigningIdentity("admin", "league", confData)
	if err != nil {
		t.Error("GetSigningIdentity", err)
	}
	t.Log("test ca GetSigningIdentity, identity = ", identity)
}

func TestRevoke(t *testing.T) {
	conf := TGetCAConfig()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	identity, err := Revoke("league", &msp.RevocationRequest{
		Name: "test1",
	}, confData)
	if err != nil {
		t.Error("GetSigningIdentity", err)
	}
	t.Log("test ca GetSigningIdentity, identity = ", identity)
}

func TGetCAConfig() *config.Config {
	rootPath := "/Users/aberic/Documents/path/go/src/github.com/ennoo/fabric-client/example"
	//rootPath := "/Users/admin/Documents/code/git/go/src/github.com/ennoo/fabric-client/example"
	conf := config.Config{}
	conf.InitClient(false, "league", "debug",
		rootPath+"/config/crypto-config",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.key",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.crt")
	conf.AddOrSetOrgForOrganizations("league", "Org1MSP",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/msp",
		map[string]string{},
		[]string{},
		[]string{"caRoot"},
	)
	conf.AddOrSetCertificateAuthority("caRoot", "https://10.10.203.51:30454",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/tlsca/tlsca.20de78630ef6a411-org1-cert.pem",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.key",
		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.crt",
		"ca.20de78630ef6a411-org1", "admin", "adminpw")
	return &conf
}

func TGetConfig() *config.Config {
	rootPath := "/Users/aberic/Documents/path/go/src/github.com/ennoo/fabric-client/example"
	//rootPath := "/Users/admin/Documents/code/git/go/src/github.com/ennoo/fabric-client/example"
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
	conf.AddOrSetOrdererForOrganizations("OrdererMSP",
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
