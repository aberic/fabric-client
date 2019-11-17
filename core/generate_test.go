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
	"encoding/hex"
	"github.com/aberic/fabric-client/config"
	"github.com/aberic/fabric-client/geneses"
	"github.com/aberic/fabric-client/grpc/proto/generate"
	"github.com/aberic/gnomon"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

const (
	leagueDomain   = "league01.com"
	orderName      = "od"
	orderDomain    = "order.com"
	order0NodeName = "order0"
	order1NodeName = "order1"
	org1Name       = "org1"
	org2Name       = "org2"
	org1Domain     = "one.com"
	org2Domain     = "two.com"
	node1          = "node1"
	node2          = "node2"
	admin          = "Admin"
	user2          = "User2"
	channelID      = "mychannel02"

	cryptoType    = generate.CryptoType_ECDSA
	signAlgorithm = generate.SignAlgorithm_ECDSAWithSHA256
)

var (
	algorithm = &generate.ReqKeyConfig_EccAlgorithm{
		EccAlgorithm: generate.EccAlgorithm_p256,
	}
)

func TestPemConfig_GenerateCrypto(t *testing.T) {
	pemConfig := &geneses.PemConfig{KeyConfig: &generate.ReqKeyConfig{
		CryptoType: cryptoType,
		Algorithm:  algorithm,
	}}
	t.Log(pemConfig.GenerateCrypto())
}

func TestPemConfig_GenerateCryptos(t *testing.T) {
	for i := 0; i < 30; i++ {
		pemConfig := &geneses.PemConfig{KeyConfig: &generate.ReqKeyConfig{
			CryptoType: cryptoType,
			Algorithm:  algorithm,
		}}
		t.Log(pemConfig.GenerateCrypto())
	}
}

func TestGenerateConfig_CreateLeague(t *testing.T) {
	priData, err := ioutil.ReadFile("/tmp/354535000/pri.key")
	if nil != err {
		t.Error(err)
	}
	priTlsData, err := ioutil.ReadFile("/tmp/364725000/pri.key")
	if nil != err {
		t.Error(err)
	}
	gc := &geneses.GenerateConfig{}
	t.Log(gc.CreateLeague(&generate.ReqCreateLeague{
		Domain:     leagueDomain,
		PriData:    priData,
		PriTlsData: priTlsData,
		Csr: &generate.CSR{
			Country:      []string{"CN"},
			Organization: []string{"league01"},
			Locality:     []string{"Beijing"},
			Province:     []string{"Beijing"},
			CommonName:   leagueDomain,
		},
		SignAlgorithm: signAlgorithm,
	}))
}

func TestGenerateConfig_CreateOrg(t *testing.T) {
	gc := &geneses.GenerateConfig{}
	t.Log(gc.CreateOrg(&generate.ReqCreateOrg{
		OrgType:      generate.OrgType_Order,
		LeagueDomain: leagueDomain,
		Name:         orderName,
		Domain:       orderDomain,
	}))
	t.Log(gc.CreateOrg(&generate.ReqCreateOrg{
		OrgType:      generate.OrgType_Peer,
		LeagueDomain: leagueDomain,
		Name:         org1Name,
		Domain:       org1Domain,
	}))
	t.Log(gc.CreateOrg(&generate.ReqCreateOrg{
		OrgType:      generate.OrgType_Peer,
		LeagueDomain: leagueDomain,
		Name:         org2Name,
		Domain:       org2Domain,
	}))
}

func TestGenerateConfig_CAGenesisAddAffiliation(t *testing.T) {
	genesisAddAffiliation("admin", "adminpw", leagueDomain, orderDomain, orderName, geneses.MspID(orderName), "rootCA", strings.Join([]string{orderName, order0NodeName}, "."), "http://10.0.61.22:7054", t)
	genesisAddAffiliation("admin", "adminpw", leagueDomain, org1Domain, org1Name, geneses.MspID(org1Name), "rootCA", strings.Join([]string{org1Name, node1}, "."), "http://10.0.61.22:7054", t)
	genesisAddAffiliation("admin", "adminpw", leagueDomain, org2Domain, org2Name, geneses.MspID(org2Name), "rootCA", strings.Join([]string{org2Name, node2}, "."), "http://10.0.61.22:7054", t)
}

func TestGenerateConfig_CAGenesisRegister(t *testing.T) {
	// 获取各组织根CA用于启动fabric-ca
	genesisRegister("admin", "adminpw", leagueDomain, orderDomain, orderName, geneses.MspID(orderName), "rootCA", strings.Join([]string{orderName, order0NodeName}, "."), "orderer,user", "orderer,user", strings.Split(geneses.CertUserCAName(orderName, orderDomain, admin), "-")[0], "adminpw", "user", "http://10.0.61.22:7054", t)
	genesisRegister("admin", "adminpw", leagueDomain, orderDomain, orderName, geneses.MspID(orderName), "rootCA", strings.Join([]string{orderName, order0NodeName}, "."), "orderer,user", "orderer,user", strings.Split(geneses.CertNodeCAName(orderName, orderDomain, order0NodeName), "-")[0], "adminpw", "orderer", "http://10.0.61.22:7054", t)
	genesisRegister("admin", "adminpw", leagueDomain, orderDomain, orderName, geneses.MspID(orderName), "rootCA", strings.Join([]string{orderName, order0NodeName}, "."), "orderer,user", "orderer,user", strings.Split(geneses.CertNodeCAName(orderName, orderDomain, order1NodeName), "-")[0], "adminpw", "orderer", "http://10.0.61.22:7054", t)
	genesisRegister("admin", "adminpw", leagueDomain, org1Domain, org1Name, geneses.MspID(org1Name), "rootCA", strings.Join([]string{org1Name, node1}, "."), "client,peer,user", "peer,user", strings.Split(geneses.CertUserCAName(org1Name, org1Domain, admin), "-")[0], "adminpw", "user", "http://10.0.61.22:7054", t)
	genesisRegister("admin", "adminpw", leagueDomain, org1Domain, org1Name, geneses.MspID(org1Name), "rootCA", strings.Join([]string{org1Name, node1}, "."), "client,peer,user", "peer,user", strings.Split(geneses.CertUserCAName(org1Name, org1Domain, user2), "-")[0], "adminpw", "user", "http://10.0.61.22:7054", t)
	genesisRegister("admin", "adminpw", leagueDomain, org1Domain, org1Name, geneses.MspID(org1Name), "rootCA", strings.Join([]string{org1Name, node1}, "."), "client,peer,user", "peer,user", strings.Split(geneses.CertNodeCAName(org1Name, org1Domain, node1), "-")[0], "adminpw", "peer", "http://10.0.61.22:7054", t)
	genesisRegister("admin", "adminpw", leagueDomain, org1Domain, org1Name, geneses.MspID(org1Name), "rootCA", strings.Join([]string{org1Name, node1}, "."), "client,peer,user", "peer,user", strings.Split(geneses.CertNodeCAName(org1Name, org1Domain, node2), "-")[0], "adminpw", "peer", "http://10.0.61.22:7054", t)
	genesisRegister("admin", "adminpw", leagueDomain, org1Domain, org1Name, geneses.MspID(org1Name), "rootCA", strings.Join([]string{org1Name, node1}, "."), "client,peer,user", "peer,user", strings.Split(geneses.CertUserCAName(org2Name, org2Domain, admin), "-")[0], "adminpw", "peer", "http://10.0.61.22:7054", t)
	genesisRegister("admin", "adminpw", leagueDomain, org1Domain, org1Name, geneses.MspID(org1Name), "rootCA", strings.Join([]string{org1Name, node1}, "."), "client,peer,user", "peer,user", strings.Split(geneses.CertUserCAName(org2Name, org2Domain, user2), "-")[0], "adminpw", "user", "http://10.0.61.22:7054", t)
	genesisRegister("admin", "adminpw", leagueDomain, org1Domain, org1Name, geneses.MspID(org1Name), "rootCA", strings.Join([]string{org1Name, node1}, "."), "client,peer,user", "peer,user", strings.Split(geneses.CertNodeCAName(org2Name, org2Domain, node1), "-")[0], "adminpw", "peer", "http://10.0.61.22:7054", t)
	genesisRegister("admin", "adminpw", leagueDomain, org1Domain, org1Name, geneses.MspID(org1Name), "rootCA", strings.Join([]string{org1Name, node1}, "."), "client,peer,user", "peer,user", strings.Split(geneses.CertNodeCAName(org2Name, org2Domain, node2), "-")[0], "adminpw", "peer", "http://10.0.61.22:7054", t)
}

func TestGenerateConfig_CreateCsr(t *testing.T) {
	t.Log(createCsr("/tmp/365412000/pri.key", orderName, orderDomain, admin, false, t))
	t.Log(createCsr("/tmp/366428000/pri.key", org1Name, org1Domain, admin, false, t))
	t.Log(createCsr("/tmp/367627000/pri.key", org1Name, org1Domain, user2, false, t))
	t.Log(createCsr("/tmp/368896000/pri.key", org2Name, org2Domain, admin, false, t))
	t.Log(createCsr("/tmp/370156000/pri.key", org2Name, org2Domain, user2, false, t))

	t.Log(createCsr("/tmp/371367000/pri.key", orderName, orderDomain, order0NodeName, true, t))
	t.Log(createCsr("/tmp/372924000/pri.key", orderName, orderDomain, order1NodeName, true, t))
	t.Log(createCsr("/tmp/374165000/pri.key", org1Name, org1Domain, node1, true, t))
	t.Log(createCsr("/tmp/375201000/pri.key", org1Name, org1Domain, node2, true, t))
	t.Log(createCsr("/tmp/376241000/pri.key", org2Name, org2Domain, node1, true, t))
	t.Log(createCsr("/tmp/377321000/pri.key", org2Name, org2Domain, node2, true, t))
}

func TestGenerateConfig_CreateOrgUser(t *testing.T) {
	t.Log(createOrgUser(generate.OrgType_Order, "/tmp/365939000/pub.key", orderName, orderDomain, admin, "http://10.0.61.22:7054", true, t))
	t.Log(createOrgUser(generate.OrgType_Peer, "/tmp/367051000/pub.key", org1Name, org1Domain, admin, "http://10.0.61.22:7054", true, t))
	t.Log(createOrgUser(generate.OrgType_Peer, "/tmp/368193000/pub.key", org1Name, org1Domain, user2, "http://10.0.61.22:7054", false, t))
	t.Log(createOrgUser(generate.OrgType_Peer, "/tmp/369533000/pub.key", org2Name, org2Domain, admin, "http://10.0.61.22:7054", true, t))
	t.Log(createOrgUser(generate.OrgType_Peer, "/tmp/370743000/pub.key", org2Name, org2Domain, user2, "http://10.0.61.22:7054", false, t))
}

func TestGenerateConfig_CreateOrgNode(t *testing.T) {
	t.Log(createOrgNode(generate.OrgType_Order, "/tmp/371936000/pub.key", orderName, orderDomain, order0NodeName, "http://10.0.61.22:7054", t))
	t.Log(createOrgNode(generate.OrgType_Order, "/tmp/373643000/pub.key", orderName, orderDomain, order1NodeName, "http://10.0.61.22:7054", t))
	t.Log(createOrgNode(generate.OrgType_Peer, "/tmp/374676000/pub.key", org1Name, org1Domain, node1, "http://10.0.61.22:7054", t))
	t.Log(createOrgNode(generate.OrgType_Peer, "/tmp/375722000/pub.key", org1Name, org1Domain, node2, "http://10.0.61.22:7054", t))
	t.Log(createOrgNode(generate.OrgType_Peer, "/tmp/376760000/pub.key", org2Name, org2Domain, node1, "http://10.0.61.22:7054", t))
	t.Log(createOrgNode(generate.OrgType_Peer, "/tmp/377848000/pub.key", org2Name, org2Domain, node2, "http://10.0.61.22:7054", t))
}

func TestGenerateConfig_GenesisBlock(t *testing.T) {
	genesis := geneses.Genesis{
		Info: &generate.ReqGenesis{
			League: &generate.LeagueInBlock{
				Domain: leagueDomain,
				Addresses: []string{
					strings.Join([]string{order0NodeName, ".", orderName, ".", orderDomain, ":7050"}, ""),
					strings.Join([]string{order1NodeName, ".", orderName, ".", orderDomain, ":7050"}, ""),
				},
				BatchTimeout: 2,
				BatchSize: &generate.BatchSize{
					MaxMessageCount:   1000,
					AbsoluteMaxBytes:  10 * 1024 * 1024,
					PreferredMaxBytes: 2 * 1024 * 1024,
				},
				Kafka: &generate.Kafka{
					Brokers: []string{"kafka1:9090", "kafka2:9091", "kafka3:9092", "kafka4:9093"},
				},
				MaxChannels: 1000,
			},
			Orgs: []*generate.OrgInBlock{
				{Domain: orderDomain, Name: orderName, Type: generate.OrgType_Order},
				{Domain: org1Domain, Name: org1Name, Type: generate.OrgType_Peer, AnchorPeers: []*generate.AnchorPeer{
					{Host: strings.Join([]string{node1, org1Name, org1Domain}, "."), Port: 7051},
				}},
				{Domain: org2Domain, Name: org2Name, Type: generate.OrgType_Peer, AnchorPeers: []*generate.AnchorPeer{
					{Host: strings.Join([]string{node1, org2Name, org2Domain}, "."), Port: 7051},
				}},
			},
		},
	}
	genesis.Init()
	if err := genesis.CreateGenesisBlock("default"); nil != err {
		t.Error(err)
	} else {
		t.Log("create genesis block success")
	}
}

func TestGenerateConfig_InspectGenesisBlock(t *testing.T) {
	data, err := ioutil.ReadFile(geneses.GenesisBlockFilePath(leagueDomain))
	//data, err := ioutil.ReadFile("/Users/aberic/Documents/path/go/src/github.com/aberic/fabric-client/geneses/example/test/channel-artifacts/genesis.block")
	if nil != err {
		t.Error(err)
	}
	str, err := resource.InspectBlock(data)
	if nil != err {
		t.Error(err)
	}
	t.Log(str)
}

func TestGenerateConfig_CreateChannelTx(t *testing.T) {
	genesis := geneses.Genesis{
		Info: &generate.ReqGenesis{
			League: &generate.LeagueInBlock{
				Domain: leagueDomain,
				Addresses: []string{
					strings.Join([]string{order0NodeName, ".", orderName, ".", orderDomain, ":7050"}, ""),
					strings.Join([]string{order1NodeName, ".", orderName, ".", orderDomain, ":7050"}, ""),
				},
				BatchTimeout: 2,
				BatchSize: &generate.BatchSize{
					MaxMessageCount:   1000,
					AbsoluteMaxBytes:  10 * 1024 * 1024,
					PreferredMaxBytes: 2 * 1024 * 1024,
				},
				Kafka: &generate.Kafka{
					Brokers: []string{"kafka1:9090", "kafka2:9091", "kafka3:9092", "kafka4:9093"},
				},
				MaxChannels: 1000,
			},
			Orgs: []*generate.OrgInBlock{
				{Domain: orderDomain, Name: orderName, Type: generate.OrgType_Order},
				{Domain: org1Domain, Name: org1Name, Type: generate.OrgType_Peer, AnchorPeers: []*generate.AnchorPeer{
					{Host: strings.Join([]string{node1, org1Name, org1Domain}, "."), Port: 7051},
				}},
				{Domain: org2Domain, Name: org2Name, Type: generate.OrgType_Peer, AnchorPeers: []*generate.AnchorPeer{
					{Host: strings.Join([]string{node1, org2Name, org2Domain}, "."), Port: 7051},
				}},
			},
		},
	}
	genesis.Init()
	if err := genesis.CreateChannelCreateTx("default", channelID); nil != err {
		t.Error(err)
	} else {
		t.Log("create genesis block success")
	}
}

func TestGenerateConfig_InspectChannelTx(t *testing.T) {
	data, err := ioutil.ReadFile(geneses.ChannelTXFilePath(leagueDomain, channelID))
	//data, err := ioutil.ReadFile("/Users/aberic/Documents/path/go/src/github.com/aberic/fabric-client/geneses/example/test/channel-artifacts/channel01.tx")
	if nil != err {
		t.Error(err)
	}
	str, err := resource.InspectChannelCreateTx(data)
	if nil != err {
		t.Error(err)
	}
	t.Log(str)
}

func TestGenerateConfig_CreateChannel(t *testing.T) {
	conf := configGenerateConfig(leagueDomain, orderName, orderDomain, order0NodeName, admin, org2Name, org2Domain, node1, admin, channelID)
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error(err)
	}
	result, err := Create(orderName, admin, "grpcs://10.0.61.23:7050", org2Name, admin, channelID,
		filepath.Join(geneses.ChannelArtifactsPath(leagueDomain), strings.Join([]string{channelID, "tx"}, ".")), // "/Users/aberic/Documents/path/go/src/github.com/aberic/fabric-client/example/config/channel-artifacts/cc6519b67c4177fc1110.tx",
		confData)
	t.Log("test query result", result)
}

func TestGenerateConfig_Join(t *testing.T) {
	conf := configGenerateConfig(leagueDomain, orderName, orderDomain, order0NodeName, admin, org1Name, org1Domain, node1, admin, channelID)
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error(err)
	}
	result := Join("grpcs://10.0.61.23:7050", org1Name, admin, channelID, node1, confData)
	t.Log("test query result", result)
}

func TestGenerateConfig_Channels(t *testing.T) {
	conf := configGenerateConfig(leagueDomain, orderName, orderDomain, order0NodeName, admin, org1Name, org1Domain, node1, admin, channelID)
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error(err)
	}
	result, err := Channels(org1Name, admin, node1, confData)
	t.Log(result)
}

func genesisAddAffiliation(enrollID, enrollSecret, leagueDomain, orgDomain, orgName, orgMspID, caName, affiliationName, url string, t *testing.T) {
	conf := genesisCAConfigCustom(enrollID, enrollSecret, leagueDomain, orgDomain, orgName, orgMspID, caName, url)
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := AddAffiliation(orgName, &msp.AffiliationRequest{
		Name:   affiliationName, // Name of the affiliation
		Force:  true,            // Creates parent affiliations if they do not exist
		CAName: caName,          // Name of the CA
	}, confData)
	if err != nil {
		t.Error("TestAddAffiliation", err)
	}
	t.Log("test ca TestAddAffiliation, result = ", result)
}

func genesisRegister(enrollID, enrollSecret, leagueDomain, orgDomain, orgName, orgMspID, caName, affiliationName, roles, delegateRoles, name, secret, registerCAType, url string, t *testing.T) {
	conf := genesisCAConfigCustom(enrollID, enrollSecret, leagueDomain, orgDomain, orgName, orgMspID, caName, url)
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := Register(orgName, &msp.RegistrationRequest{
		Name:           name,
		Type:           registerCAType,  // (e.g. "client, orderer, peer, app, user")
		MaxEnrollments: -1,              // if omitted, this defaults to max_enrollments configured on the server
		Affiliation:    affiliationName, // The identity's affiliation e.g. org1.department1
		Attributes: []msp.Attribute{ // Optional attributes associated with this identity.
			// Attribute defines additional attributes that may be passed along during registration
			{Name: "hf.Registrar.Roles", Value: roles},                 // "client,orderer,peer,user" 注册者允许管理的角色列表
			{Name: "hf.Registrar.DelegateRoles", Value: delegateRoles}, // "client,orderer,peer,user" 注册者可以赋予被注册身份的hf.Registrar.Roles属性的角色列表
			{Name: "hf.Registrar.Attributes", Value: "*"},              // 注册者允许注册的属性列表
			{Name: "hf.GenCRL", Value: "true"},                         // 要注册的身份是否可以生成CRL
			{Name: "hf.Revoker", Value: "true"},                        // 要注册的身份是否可以回收证书
			{Name: "hf.AffiliationMgr", Value: "false"},                // 要注册的身份是否可以管理联盟
			{Name: "hf.IntermediateCA", Value: "true"},
		},
		CAName: caName,
		Secret: secret,
	}, confData)
	if err != nil {
		t.Error("Register", err)
	}
	t.Log("test ca Register, result = ", result)
}

func createCsr(priKeyFilePath, orgName, orgDomain, childName string, isNode bool, t *testing.T) error {
	priData, err := ioutil.ReadFile(priKeyFilePath)
	if nil != err {
		t.Error(err)
	}
	var commonName string
	if isNode {
		commonName = strings.Split(geneses.CertNodeCAName(orgName, orgDomain, childName), "-")[0]
	} else {
		commonName = strings.Split(geneses.CertUserCAName(orgName, orgDomain, childName), "-")[0]
	}
	gc := &geneses.GenerateConfig{}
	return gc.CreateCsr(&generate.ReqCreateCsr{
		PriKey:       priData,
		LeagueDomain: leagueDomain,
		OrgName:      orgName,
		OrgDomain:    orgDomain,
		Name: &generate.CSR{
			Country:      []string{"CN"},
			Organization: []string{orgName},
			Locality:     []string{childName},
			Province:     []string{orgName},
			CommonName:   commonName,
		},
		SignAlgorithm: signAlgorithm,
	})
}

func createOrgNode(orgType generate.OrgType, pubTlsKeyFilePath, orgName, orgDomain, nodeName, caURL string, t *testing.T) error {
	commonName := strings.Split(geneses.CertNodeCAName(orgName, orgDomain, nodeName), "-")[0]
	csrPem, err := ioutil.ReadFile(geneses.CsrFilePath(leagueDomain, orgName, orgDomain, commonName))
	if nil != err {
		t.Error(err)
	}
	pubTlsData, err := ioutil.ReadFile(pubTlsKeyFilePath)
	if nil != err {
		t.Error(err)
	}
	gc := &geneses.GenerateConfig{}
	return gc.CreateOrgNode(&generate.ReqCreateOrgNode{
		OrgType: orgType,
		OrgChild: &generate.OrgChild{
			LeagueDomain: leagueDomain,
			OrgName:      orgName,
			OrgDomain:    orgDomain,
			Name:         nodeName,
			PubTlsData:   pubTlsData,
			EnrollInfo: &generate.EnrollInfo{
				CsrPem:            csrPem,
				FabricCaServerURL: caURL,
				EnrollRequest: &generate.EnrollRequest{
					EnrollID: commonName,
					Secret:   "adminpw",
					//Profile:  "tls",
					Hosts: []string{commonName}, // []string{"hello.cn"}
					Name: &generate.CSR{
						Country:      []string{"CN"},
						Organization: []string{orgName},
						Locality:     []string{nodeName},
						Province:     []string{orgName},
						CommonName:   commonName,
					},
				},
				NotAfter:  50000,
				NotBefore: 0,
			},
			SignAlgorithm: signAlgorithm,
		},
	})
}

func createOrgUser(orgType generate.OrgType, pubTlsKeyFilePath, orgName, orgDomain, nodeName, caURL string, isAdmin bool, t *testing.T) error {
	commonName := strings.Split(geneses.CertUserCAName(orgName, orgDomain, nodeName), "-")[0]
	csrPem, err := ioutil.ReadFile(geneses.CsrFilePath(leagueDomain, orgName, orgDomain, commonName))
	if nil != err {
		t.Error(err)
	}
	pubTlsData, err := ioutil.ReadFile(pubTlsKeyFilePath)
	if nil != err {
		t.Error(err)
	}
	gc := &geneses.GenerateConfig{}
	return gc.CreateOrgUser(&generate.ReqCreateOrgUser{
		IsAdmin: isAdmin,
		OrgType: orgType,
		OrgChild: &generate.OrgChild{
			LeagueDomain: leagueDomain,
			OrgName:      orgName,
			OrgDomain:    orgDomain,
			Name:         nodeName,
			PubTlsData:   pubTlsData,
			EnrollInfo: &generate.EnrollInfo{
				CsrPem:            csrPem,
				FabricCaServerURL: caURL,
				EnrollRequest: &generate.EnrollRequest{
					EnrollID: commonName,
					Secret:   "adminpw",
					//Profile:  "tls",
					Hosts: []string{commonName}, // []string{"hello.cn"}
					Name: &generate.CSR{
						Country:      []string{"CN"},
						Organization: []string{orgName},
						Locality:     []string{nodeName},
						Province:     []string{orgName},
						CommonName:   commonName,
					},
				},
				NotAfter:  50000,
				NotBefore: 0,
			},
			SignAlgorithm: signAlgorithm,
		},
	})
}

func genesisCAConfigCustom(enrollID, enrollSecret, leagueDomain, orgDomain, orgName, orgMspID, caName, url string) *config.Config {
	conf := config.Config{}
	orgPath := geneses.CryptoUserTmpPath(leagueDomain, orgDomain, orgName)
	conf.InitCustomClient(false, orgName, "debug", "", "", "",
		nil, nil, nil, nil,
		&config.ClientCredentialStore{
			Path:        orgPath,
			CryptoStore: &config.ClientCredentialStoreCryptoStore{Path: orgPath},
		},
		&config.ClientBCCSP{
			Security: &config.ClientBCCSPSecurity{
				Enabled:       true,
				HashAlgorithm: "SHA2",
				SoftVerify:    true,
				Level:         256,
				Default:       &config.ClientBCCSPSecurityDefault{Provider: "SW"},
			},
		})
	conf.AddOrSetOrgForOrganizations(orgName, orgMspID, "/Users/aberic/Documents/code/ca/demo/msp", map[string]string{}, []string{}, []string{caName})
	conf.AddOrSetCertificateAuthority(caName, url,
		"/Users/aberic/Documents/path/go/src/github.com/aberic/gnomon/tmp/example/ca/pemp384/rootCA.crt",
		"/Users/aberic/Documents/path/go/src/github.com/aberic/gnomon/tmp/example/ca/pemp256/rootCA.key",
		"/Users/aberic/Documents/path/go/src/github.com/aberic/gnomon/tmp/example/ca/pemp256/rootCA.crt",
		caName, enrollID, enrollSecret)
	return &conf
}

func configGenerateConfig(leagueDomain, orderName, orderDomain, orderNodeName, orderUserName, orgName, orgDomain, peerName, orgUserName, channelID string) *config.Config {
	rootPath := geneses.CryptoConfigPath(leagueDomain)
	//rootPath := "/Users/admin/Documents/code/git/go/src/github.com/aberic/fabric-client/example"
	conf := config.Config{}
	_, orgUserPath := geneses.CryptoOrgAndNodePath(leagueDomain, orgDomain, orgName, orgUserName, true, geneses.CcnAdmin)
	conf.InitCustomClient(true, orgName, "debug", rootPath, // rootPath+"/config/crypto-config"
		filepath.Join(orgUserPath, "tls", "client.key"), // rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.key",
		filepath.Join(orgUserPath, "tls", "client.crt"), // rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.crt")
		&config.ClientPeer{
			Timeout: &config.ClientPeerTimeout{Connection: "10s", Response: "180s",
				Discovery: &config.ClientPeerTimeoutDiscovery{GreyListExpiry: "10s"},
			},
		},
		&config.ClientEventService{Timeout: &config.ClientEventServiceTimeout{RegistrationResponse: "15s"}},
		&config.ClientOrder{Timeout: &config.ClientOrderTimeout{Connection: "15s", Response: "15s"}},
		&config.ClientGlobal{
			Timeout: &config.ClientGlobalTimeout{
				Query:   "180s",
				Execute: "180s",
				Resmgmt: "180s",
			},
			Cache: &config.ClientGlobalCache{
				ConnectionIdle:    "30s",
				EventServiceIdle:  "2m",
				ChannelConfig:     "30m",
				ChannelMembership: "30s",
				Discovery:         "10s",
				Selection:         "10m",
			},
		},
		&config.ClientCredentialStore{
			Path:        path.Join(orgUserPath, "msp", "signcerts"),
			CryptoStore: &config.ClientCredentialStoreCryptoStore{Path: path.Join(orgUserPath, "msp")},
		},
		&config.ClientBCCSP{
			Security: &config.ClientBCCSPSecurity{
				Enabled:       true,
				HashAlgorithm: "SHA2",
				SoftVerify:    true,
				Level:         256,
				Default:       &config.ClientBCCSPSecurityDefault{Provider: "SW"},
			},
		},
	)
	//conf.AddOrSetPeerForChannel("cc6519b67c4177fc11", "peer0",
	//	true, true, true, true)
	conf.AddOrSetPeerForChannel(channelID, peerName,
		true, true, true, true)
	conf.AddOrSetQueryChannelPolicyForChannel(channelID, "500ms", "5s",
		1, 1, 5, 2.0)
	conf.AddOrSetDiscoveryPolicyForChannel(channelID, "500ms", "5s",
		2, 4, 2.0)
	conf.AddOrSetEventServicePolicyForChannel(channelID, "PreferOrg", "RoundRobin",
		"6s", 5, 8)
	_, orderUserPath := geneses.CryptoOrgAndNodePath(leagueDomain, orderDomain, orderName, orderUserName, false, geneses.CcnAdmin)
	conf.AddOrSetOrdererForOrganizations(orderName, strings.Join([]string{orderName, "MSP"}, ""),
		path.Join(orderUserPath, "msp"), // rootPath+"/config/crypto-config/ordererOrganizations/20de78630ef6a411/users/Admin@20de78630ef6a411/msp",
		map[string]string{
			//"Admin": rootPath + "/config/crypto-config/ordererOrganizations/20de78630ef6a411/users/Admin@20de78630ef6a411/msp/signcerts/Admin@20de78630ef6a411-cert.pem",
			orderUserName: filepath.Join(orderUserPath, "msp", "signcerts", geneses.CertUserCAName(orderName, orderDomain, orderUserName)),
		},
	)
	conf.AddOrSetOrgForOrganizations(orgName, strings.Join([]string{orgName, "MSP"}, ""),
		path.Join(orgUserPath, "msp"), // rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/msp",
		map[string]string{
			//"Admin": rootPath + "/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/msp/signcerts/Admin@20de78630ef6a411-org1-cert.pem",
			orgUserName: filepath.Join(orgUserPath, "msp", "signcerts", geneses.CertUserCAName(orgName, orgDomain, orgUserName)),
		},
		[]string{peerName}, //"peer0", "peer1",
		[]string{},
	)
	tlsRootCaFilePath := filepath.Join(geneses.CryptoRootTLSCAPath(leagueDomain), geneses.CertRootTLSCAName(leagueDomain))
	conf.AddOrSetOrderer(orderName, "grpcs://10.0.61.23:7050",
		strings.Join([]string{orderNodeName, orderName, orderDomain}, "."), "0s", "20s",
		tlsRootCaFilePath, // rootPath+"/config/crypto-config/ordererOrganizations/20de78630ef6a411/tlsca/tlsca.20de78630ef6a411-cert.pem",
		false, false, false)
	conf.AddOrSetPeer(peerName, "grpcs://10.0.61.23:7051",
		"grpcs://10.0.61.23:7052", strings.Join([]string{peerName, orgName, orgDomain}, "."),
		"0s", "20s",
		tlsRootCaFilePath,
		false, false, false)
	return &conf
}

func TestGenerateConfig_eccSKI(t *testing.T) {
	//priKey, err := gnomon.CryptoECC().LoadPriPemFP("/tmp/366428000/pri.key")
	priKey, err := gnomon.CryptoECC().LoadPriPemFP("/tmp/368896000/pri.key")
	if nil != err {
		t.Error(err)
	}
	gc := &geneses.GenerateConfig{}
	data := gc.EccSKI(priKey)
	t.Log(data)
	t.Log(hex.EncodeToString(data))
}
