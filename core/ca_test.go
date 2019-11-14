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
	"github.com/ennoo/fabric-client/config"
	"github.com/ennoo/rivet/utils/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
	"testing"
)

//func TestEnrollCustom(t *testing.T) {
//	caClientConfig := &ca.ClientCAConfig{
//		Url:        "http://10.0.61.22:7054",
//		AdminKey:   "/Users/aberic/Documents/path/go/src/github.com/aberic/gnomon/tmp/example/ca/pemp256/rootCA.key",
//		AdminCert:  "/Users/aberic/Documents/path/go/src/github.com/aberic/gnomon/tmp/example/ca/pemp256/rootCA.crt",
//		CryptoConfig: &ca.CryptoConfig{
//			Family:    "ecdsa",
//			Algorithm: "P256-SHA256",
//			Hash:      "SHA2-256",
//		},
//	}
//
//	caClient := ca.ClientCA{}
//	if err := caClient.Init(caClientConfig); nil != err {
//		t.Error(err)
//	}
//	result, err := caClient.Enroll(&ca.EnrollRequest{
//		EnrollID: "admin",
//		Secret:   "adminpw",
//		//Profile:  "tls",
//		Hosts: []string{"hello.cn"},
//		Name: pkix.Name{
//			Country:            []string{"CN"},
//			Province:           []string{"Beijing"},
//			Locality:           []string{"Beijing"},
//			Organization:       []string{"world"},
//			OrganizationalUnit: []string{"BlockChain"},
//			CommonName:         "admin",
//		},
//		NotAfter:  time.Now().Add(50000 * 24 * time.Hour),
//		NotBefore: time.Now(),
//	})
//	if nil != err {
//		t.Error(err)
//	}
//	t.Log("result = ", result)
//}

func TestEnroll(t *testing.T) {
	conf := TGetCAConfigCustom()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	//attrReqs := []*msp.AttributeRequest{{Name: "test", Optional: true}}
	err = Enroll("league", "admin", confData, []msp.EnrollmentOption{
		msp.WithSecret("adminpw"),
		msp.WithType("x509" /*or idemix, which is not support now*/),
		//msp.WithProfile("tls"),
		//msp.WithLabel("ForFabric"),
		//msp.WithAttributeRequests(attrReqs),
	})
	if err != nil {
		t.Error("Enroll", err)
	}
	t.Log("test ca Enroll")
}

func TestReenroll(t *testing.T) {
	conf := TGetCAConfigCustom()
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
	conf := TGetCAConfigCustom()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := Register("league", &msp.RegistrationRequest{
		Name:           "test3",
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
		CAName: "rootCA",
		Secret: "adminpw",
	}, confData)
	if err != nil {
		t.Error("Register", err)
	}
	t.Log("test ca Register, result = ", result)
}

func TestAddAffiliation(t *testing.T) {
	conf := TGetCAConfigCustom()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := AddAffiliation("league", &msp.AffiliationRequest{
		Name:   "org100.peer1", // Name of the affiliation
		Force:  true,           // Creates parent affiliations if they do not exist
		CAName: "rootCA",       // Name of the CA
	}, confData)
	if err != nil {
		t.Error("TestAddAffiliation", err)
	}
	t.Log("test ca TestAddAffiliation, result = ", result)
}

func TestRemoveAffiliation(t *testing.T) {
	conf := TGetCAConfigCustom()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := RemoveAffiliation("league", &msp.AffiliationRequest{
		Name:   "orderer.orderer02", // Name of the affiliation
		Force:  true,                // Creates parent affiliations if they do not exist
		CAName: "rootCA",            // Name of the CA
	}, confData)
	if err != nil {
		t.Error("GetAllAffiliations", err)
	}
	t.Log("test ca GetAllAffiliations, result = ", result)
}

func TestModifyAffiliation(t *testing.T) {
	conf := TGetCAConfigCustom()
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
	conf := TGetCAConfigCustom()
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
	conf := TGetCAConfigCustom()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := GetAllAffiliationsByCaName("league", "rootCA", confData)
	if err != nil {
		t.Error("GetAllAffiliationsByCaName", err)
	}
	t.Log("test ca GetAllAffiliationsByCaName, result = ", result)
}

func TestGetAffiliation(t *testing.T) {
	conf := TGetCAConfigCustom()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := GetAffiliation("org1.department1", "league", confData)
	if err != nil {
		t.Error("GetAffiliation", err)
	}
	t.Log("test ca GetAffiliation, result = ", result)
}

func TestGetAffiliationByCaName(t *testing.T) {
	conf := TGetCAConfigCustom()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	result, err := GetAffiliationByCaName("org1.department1", "league", "rootCA", confData)
	if err != nil {
		t.Error("GetAffiliationByCaName", err)
	}
	t.Log("test ca GetAffiliationByCaName, result = ", result)
}

func TestCreateIdentity(t *testing.T) {
	conf := TGetCAConfigCustom()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	identity, err := CreateIdentity("league", &msp.IdentityRequest{
		ID:          "test4", // The enrollment ID which uniquely identifies an identity (required)
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
		Secret:         "adminpw",
		CAName:         "rootCA",
	}, confData)
	if err != nil {
		t.Error("CreateIdentity", err)
	}
	t.Log("test ca CreateIdentity, identity = ", identity)
}

func TestModifyIdentity(t *testing.T) {
	conf := TGetCAConfigCustom()
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
	conf := TGetCAConfigCustom()
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
	conf := TGetCAConfigCustom()
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}
	identity, err := GetIdentity("test3", "league", confData)
	if err != nil {
		t.Error("GetIdentity", err)
	}
	t.Log("test ca GetIdentity, identity = ", identity)
}

func TestGetIdentityByCaName(t *testing.T) {
	conf := TGetCAConfigCustom()
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
	conf := TGetCAConfigCustom()
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
	conf := TGetCAConfigCustom()
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
	privateKey := `-----BEGIN PRIVATE KEY-----
MIGkAgEBBDCAiXrmR0mhOpysFw8+K12tipf+2d3Oyc0lSQKBUv44KrkBkHcBFkrC
jLyLpYLR/OygBwYFK4EEACKhZANiAAQNerMkKw2cTcMSpLi5LLmwzBeaeLQc8is3
8ldsgzbawWR+/PEAaUJ/O/Ko4xTAqOBf3ZAHAvJ18U1ULJJULozueXnnh8rl0YCS
z8SnQjQNBtlhtPwNtaG+NjVh91iqMgw=
-----END PRIVATE KEY-----`
	cert := `-----BEGIN CERTIFICATE-----
MIICVDCCAdugAwIBAgIIEY6bhng8FTEwCgYIKoZIzj0EAwIwaTELMAkGA1UEBhMC
Q04xEDAOBgNVBAgTB0JlaWppbmcxEDAOBgNVBAcTB0JlaWppbmcxDzANBgNVBAoT
Bkdub21vbjERMA8GA1UECxMIR25vbW9uUkQxEjAQBgNVBAMTCWFiZXJpYy5jbjAe
Fw0xOTEwMjIwNjI5MTVaFw0zMzA2MzAwNjI5MTVaMGkxCzAJBgNVBAYTAkNOMRAw
DgYDVQQIEwdCZWlqaW5nMRAwDgYDVQQHEwdCZWlqaW5nMQ8wDQYDVQQKEwZHbm9t
b24xETAPBgNVBAsTCEdub21vblJEMRIwEAYDVQQDEwlhYmVyaWMuY24wdjAQBgcq
hkjOPQIBBgUrgQQAIgNiAAQNerMkKw2cTcMSpLi5LLmwzBeaeLQc8is38ldsgzba
wWR+/PEAaUJ/O/Ko4xTAqOBf3ZAHAvJ18U1ULJJULozueXnnh8rl0YCSz8SnQjQN
BtlhtPwNtaG+NjVh91iqMgyjUDBOMA4GA1UdDwEB/wQEAwIBFjAdBgNVHSUEFjAU
BggrBgEFBQcDAgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUwAwEB/zAMBgNVHQ4EBQQD
AQIDMAoGCCqGSM49BAMCA2cAMGQCMEqMc8XPPFJCV1JJmu0r7kjY2dTtcE+wJebu
f7vd3UOz24Y/Ci/1sUhW8930P5k8QQIwS1qv367kJazqv4juISFPYh6yAEr/+gE4
mhQItsh2GGpUyCfh/gjig/XnSCcMTPuR
-----END CERTIFICATE-----`
	conf := TGetCAConfigCustom()
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
	conf := TGetCAConfigCustom()
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
	conf := TGetCAConfigCustom()
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

func TestCAInfo(t *testing.T) {
	conf := TGetCAConfigCustom()
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

//func TGetCAConfig() *config.Config {
//	rootPath := "/Users/aberic/Documents/path/go/src/github.com/ennoo/fabric-client/example"
//	//rootPath := "/Users/admin/Documents/code/git/go/src/github.com/ennoo/fabric-client/example"
//	conf := config.Config{}
//	conf.InitClient(false, "league", "debug",
//		rootPath+"/config/crypto-config",
//		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.key",
//		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.crt")
//	conf.AddOrSetOrgForOrganizations("league", "Org1MSP",
//		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/msp",
//		map[string]string{},
//		[]string{},
//		[]string{"caRoot"},
//	)
//	conf.AddOrSetCertificateAuthority("caRoot", "https://10.10.203.51:30454",
//		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/tlsca/tlsca.20de78630ef6a411-org1-cert.pem",
//		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.key",
//		rootPath+"/config/crypto-config/peerOrganizations/20de78630ef6a411-org1/users/Admin@20de78630ef6a411-org1/tls/client.crt",
//		"ca.20de78630ef6a411-org1", "admin", "adminpw")
//	return &conf
//}

func TGetCAConfigCustom() *config.Config {
	conf := config.Config{}
	conf.InitCustomClient(false, "league", "debug", "", "", "",
		nil, nil, nil, nil,
		&config.ClientCredentialStore{
			Path:        strings.Join([]string{"/Users/aberic/Documents/code/ca/demo", "league", "crypto-config"}, "/"),
			CryptoStore: &config.ClientCredentialStoreCryptoStore{Path: strings.Join([]string{"/Users/aberic/Documents/code/ca/demo", "league", "crypto-config"}, "/")},
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
	conf.AddOrSetOrgForOrganizations("league", "leagueMSP", "/Users/aberic/Documents/code/ca/demo/msp", map[string]string{}, []string{}, []string{"rootCA"})
	conf.AddOrSetCertificateAuthority("rootCA", "http://10.0.61.22:7054",
		"/Users/aberic/Documents/path/go/src/github.com/aberic/gnomon/tmp/example/ca/pemp384/rootCA.crt",
		"/Users/aberic/Documents/path/go/src/github.com/aberic/gnomon/tmp/example/ca/pemp256/rootCA.key",
		"/Users/aberic/Documents/path/go/src/github.com/aberic/gnomon/tmp/example/ca/pemp256/rootCA.crt",
		"rootCA", "admin", "adminpw")
	return &conf
}

func TestCALong(t *testing.T) {
	orgName := "league"
	caName := "rootCA"
	conf := TCAConfigCustom(orgName)
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}

	CAAddAffiliation(orgName, caName, "orderer", confData, t)
	CARegister(orgName, caName, "orderer", "orderer", "adminpw", "orderer", confData, t)

	CAAddAffiliation(orgName, caName, "org1", confData, t)
	CARegister(orgName, caName, "org1", "peer", "adminpw", "org1", confData, t)

	CAAddAffiliation(orgName, caName, "org2", confData, t)
	CARegister(orgName, caName, "org2", "peer", "adminpw", "org2", confData, t)

	// 生成默认组织内容
	for i := 0; i < 10; i++ {
		affiliationName := strings.Join([]string{"orderer.orderer", strconv.Itoa(i)}, "")
		CAAddAffiliation(orgName, caName, affiliationName, confData, t)
		CARegister(orgName, caName, affiliationName, "orderer", "adminpw", affiliationName, confData, t)
	}
	for i := 0; i < 10; i++ {
		affiliationName := strings.Join([]string{"org1.peer", strconv.Itoa(i)}, "")
		CAAddAffiliation(orgName, caName, affiliationName, confData, t)
		CARegister(orgName, caName, affiliationName, "peer", "adminpw", affiliationName, confData, t)
	}
	for i := 0; i < 10; i++ {
		affiliationName := strings.Join([]string{"org2.peer", strconv.Itoa(i)}, "")
		CAAddAffiliation(orgName, caName, affiliationName, confData, t)
		CARegister(orgName, caName, affiliationName, "peer", "adminpw", affiliationName, confData, t)
	}

	//CALong("league", t)
}

func CAAddAffiliation(orgName, caName, affiliationName string, confData []byte, t *testing.T) {
	addAffiliationResult, err := AddAffiliation(orgName, &msp.AffiliationRequest{
		Name:   affiliationName, // Name of the affiliation
		Force:  true,            // Creates parent affiliations if they do not exist
		CAName: caName,          // Name of the CA
	}, confData)
	if err != nil {
		t.Error("GetAllAffiliations", err)
	}
	t.Log("test ca GetAllAffiliations, result = ", addAffiliationResult)
}

func CARegister(orgName, caName, idName, regType, secret, affiliationName string, confData []byte, t *testing.T) {
	registerResult, err := Register(orgName, &msp.RegistrationRequest{
		Name:           idName,
		Type:           regType,         // (e.g. "client, orderer, peer, app, user")
		MaxEnrollments: -1,              // if omitted, this defaults to max_enrollments configured on the server
		Affiliation:    affiliationName, // The identity's affiliation e.g. org1.department1
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
		CAName: caName,
		Secret: secret,
	}, confData)
	if err != nil {
		t.Error("Register", err)
	}
	t.Log("test ca Register, result = ", registerResult)
}

func CALong(orgName string, t *testing.T) {
	conf := TCAConfigCustom(orgName)
	confData, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("yaml", err)
	}

	//attrReqs := []*msp.AttributeRequest{{Name: "test", Optional: true}}
	err = Enroll(orgName, "admin", confData, []msp.EnrollmentOption{
		msp.WithSecret("adminpw"),
		msp.WithType("x509" /*or idemix, which is not support now*/),
		//msp.WithProfile("tls"),
		//msp.WithLabel("ForFabric"),
		//msp.WithAttributeRequests(attrReqs),
	})
	if err != nil {
		t.Error("Enroll", err)
	}
	t.Log("test ca Enroll")
}

func TCAConfigCustom(orgName string) *config.Config {
	rootPath := strings.Join([]string{"../example/ca/league", "crypto-config"}, "/")
	conf := config.Config{}
	conf.InitCustomClient(false, orgName, "debug", "", "", "",
		nil, nil, nil, nil,
		&config.ClientCredentialStore{
			Path:        strings.Join([]string{rootPath, orgName, "crypto-config"}, "/"),
			CryptoStore: &config.ClientCredentialStoreCryptoStore{Path: strings.Join([]string{rootPath, orgName, "crypto-config"}, "/")},
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
	conf.AddOrSetOrgForOrganizations(orgName, "league.cn", strings.Join([]string{rootPath, "/ordererOrganizations/league.cn/users/Admin@league.cn/msp"}, "/"), map[string]string{}, []string{}, []string{"rootCA"})
	conf.AddOrSetCertificateAuthority("rootCA", "http://10.0.61.22:7054",
		"/Users/aberic/Documents/path/go/src/github.com/aberic/gnomon/tmp/example/ca/pemp384/rootCA.crt",
		"/Users/aberic/Documents/path/go/src/github.com/aberic/gnomon/tmp/example/ca/pemp256/rootCA.key",
		"/Users/aberic/Documents/path/go/src/github.com/aberic/gnomon/tmp/example/ca/pemp256/rootCA.crt",
		"rootCA", "admin", "adminpw")
	return &conf
}
