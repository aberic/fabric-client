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

package geneses

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/aberic/fabric-client/grpc/proto/generate"
	"github.com/aberic/gnomon"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	GeneratePriKeyFileName = "private.key"
)

type GenerateConfig struct{}

func (gc *GenerateConfig) SignAlgorithm(signAlgorithm generate.SignAlgorithm) x509.SignatureAlgorithm {
	switch signAlgorithm {
	default:
		return x509.ECDSAWithSHA256
	case generate.SignAlgorithm_ECDSAWithSHA256:
		return x509.ECDSAWithSHA256
	case generate.SignAlgorithm_ECDSAWithSHA384:
		return x509.ECDSAWithSHA384
	case generate.SignAlgorithm_ECDSAWithSHA512:
		return x509.ECDSAWithSHA512
	case generate.SignAlgorithm_SHA256WithRSA:
		return x509.SHA256WithRSA
	case generate.SignAlgorithm_SHA512WithRSA:
		return x509.SHA512WithRSA
	}
}

func (gc *GenerateConfig) CreateLeague(league *generate.ReqCreateLeague) error {
	var err error
	caPath, certName, tlsCaPath, tlsCertName := gc.getRootCA(league.Domain)
	if err = gc.generateCryptoLeagueCrt(league, league.PriData, caPath, certName); nil != err {
		return err
	}

	if err = gc.generateCryptoLeagueCrt(league, league.PriTlsData, tlsCaPath, tlsCertName); nil != err {
		return err
	}
	return nil
}

func (gc *GenerateConfig) CreateOrg(org *generate.ReqCreateOrg) error {
	var err error
	orgMspPath := CryptoOrgMspPath(org.LeagueDomain, org.Domain, org.Name, org.OrgType == generate.OrgType_Peer)
	if gnomon.File().PathExists(orgMspPath) {
		return errors.New("org already exist")
	}
	if err = os.MkdirAll(orgMspPath, 0755); nil != err {
		return err
	}
	adminCertsPath := path.Join(orgMspPath, "admincerts")
	caCertsPath := path.Join(orgMspPath, "cacerts")
	tlsCaCertsPath := path.Join(orgMspPath, "tlscacerts")
	if err = os.Mkdir(adminCertsPath, 0755); nil != err {
		return err
	}
	if err = os.Mkdir(caCertsPath, 0755); nil != err {
		return err
	}
	if err = os.Mkdir(tlsCaCertsPath, 0755); nil != err {
		return err
	}

	caPath, certName, tlsCaPath, tlsCertName := gc.getRootCA(org.LeagueDomain)
	rootCaCertFilePath := filepath.Join(caPath, certName)
	caCertsFilePath := filepath.Join(caCertsPath, certName)
	if _, err = gnomon.File().Copy(rootCaCertFilePath, caCertsFilePath); nil != err {
		return err
	}
	rootTlsCaCertFilePath := filepath.Join(tlsCaPath, tlsCertName)
	tlsCaCertsFilePath := filepath.Join(tlsCaCertsPath, tlsCertName)
	if _, err = gnomon.File().Copy(rootTlsCaCertFilePath, tlsCaCertsFilePath); nil != err {
		return err
	}
	return nil
}

func (gc *GenerateConfig) CreateCsr(reqCsr *generate.ReqCreateCsr) error {
	subject := gc.getSubject(reqCsr.Name)
	if subject.CommonName == "" {
		return errors.New("missing CommonName")
	}
	rawSubj := subject.ToRDNSequence()

	asn1Subj, err := asn1.Marshal(rawSubj)
	if err != nil {
		return err
	}

	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		SignatureAlgorithm: gc.SignAlgorithm(reqCsr.SignAlgorithm),
		DNSNames:           []string{subject.CommonName},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, gc.getPriKeyFromBytes(reqCsr.PriKey))
	if err != nil {
		return err
	}

	csrPath := CsrPath(reqCsr.LeagueDomain, reqCsr.OrgName, reqCsr.OrgDomain)
	if !gnomon.File().PathExists(csrPath) {
		if err = os.MkdirAll(csrPath, 0755); nil != err {
			return err
		}
	}
	fileIO, err := os.OpenFile(CsrFilePath(reqCsr.LeagueDomain, reqCsr.OrgName, reqCsr.OrgDomain, subject.CommonName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	return pem.Encode(fileIO, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
}

func (gc *GenerateConfig) CreateOrgNode(node *generate.ReqCreateOrgNode) error {
	var (
		isPeer bool
		ccn    = CcnNode
		err    error
	)
	if node.OrgType == generate.OrgType_Peer {
		isPeer = true
	} else {
		isPeer = false
	}
	orgPath, nodePath := CryptoOrgAndNodePath(node.OrgChild.LeagueDomain, node.OrgChild.OrgDomain, node.OrgChild.OrgName, node.OrgChild.Name, isPeer, ccn)
	if err = gc.orgChildExec(node.OrgChild.LeagueDomain, orgPath, nodePath, ccn); nil != err {
		return err
	}
	if err = gc.generateCryptoOrgChild(node.OrgChild, nodePath, true); nil != err {
		return err
	}
	orgMspAdminCertPath := path.Join(orgPath, "msp", "admincerts")
	adminCertFileNames, err := gnomon.File().LoopFileNames(orgMspAdminCertPath)
	if nil != err {
		return err
	}
	for _, adminCertFileName := range adminCertFileNames {
		orgMspAdminCertFilePath := filepath.Join(orgMspAdminCertPath, adminCertFileName)
		orgDirAdminCertFilePath := filepath.Join(nodePath, "msp", "admincerts", adminCertFileName)
		if _, err = gnomon.File().Copy(orgMspAdminCertFilePath, orgDirAdminCertFilePath); nil != err {
			return err
		}
	}
	return nil
}

func (gc *GenerateConfig) CreateOrgUser(user *generate.ReqCreateOrgUser) error {
	var (
		nodesName       string
		isPeer          bool
		ccn             ClientCANode
		orgsDirPathName []string
		err             error
	)
	if user.OrgType == generate.OrgType_Peer {
		nodesName = "peers"
		isPeer = true
	} else {
		nodesName = "orderers"
		isPeer = false
	}
	if user.IsAdmin {
		ccn = CcnAdmin
	} else {
		ccn = CcnUser
	}
	orgPath, nodePath := CryptoOrgAndNodePath(user.OrgChild.LeagueDomain, user.OrgChild.OrgDomain, user.OrgChild.OrgName, user.OrgChild.Name, isPeer, ccn)
	if err = gc.orgChildExec(user.OrgChild.LeagueDomain, orgPath, nodePath, ccn); nil != err {
		return err
	}
	if err = gc.generateCryptoOrgChild(user.OrgChild, nodePath, false); nil != err {
		return err
	}
	adminCertPath := path.Join(nodePath, "msp", "admincerts")
	signCertPath := path.Join(nodePath, "msp", "signcerts")
	signCertFileName := CertUserCAName(user.OrgChild.OrgName, user.OrgChild.OrgDomain, user.OrgChild.Name)
	adminCertFilePath := filepath.Join(adminCertPath, signCertFileName)
	signCertFilePath := filepath.Join(signCertPath, signCertFileName)
	if _, err = gnomon.File().Copy(signCertFilePath, adminCertFilePath); nil != err {
		return err
	}
	if !user.IsAdmin {
		return nil
	}
	// admin用户证书处理
	orgMspAdminCertPath := path.Join(orgPath, "msp", "admincerts")
	orgMspAdminCertFilePath := filepath.Join(orgMspAdminCertPath, signCertFileName)
	if _, err = gnomon.File().Copy(signCertFilePath, orgMspAdminCertFilePath); nil != err {
		return err
	}
	if orgsDirPathName, err = gnomon.File().LoopOneDirs(path.Join(orgPath, nodesName)); nil != err {
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil
		}
		return err
	}
	for _, orgDirPath := range orgsDirPathName {
		orgAdminCertPath := path.Join(orgPath, nodesName, orgDirPath, "msp", "admincerts")
		orgAdminCertFilePath := filepath.Join(orgAdminCertPath, signCertFileName)
		if _, err = gnomon.File().Copy(signCertFilePath, orgAdminCertFilePath); nil != err {
			return err
		}
	}
	return nil
}

//func (gc *GenerateConfig) UpdateChannel(sdk *fabsdk.FabricSDK, channelID, orgName, orgUser, orderURL string) error {
//	orgChannelClientContext := sdk.ChannelContext(channelID, fabsdk.WithOrg(orgName), fabsdk.WithUser(orgUser))
//	channelCtx, err := orgChannelClientContext()
//	if err != nil {
//		return err
//	}
//	reqCtx, cancel := context.NewRequest(channelCtx, context.WithTimeoutType(fab.PeerResponse))
//	clientProvider := sdk.Context()
//	ClientContext, err := clientProvider()
//	order, err := orderer.New(ClientContext.EndpointConfig(), orderer.WithURL(orderURL)) // "localhost:9999"
//	defer cancel()
//	block, err := resource.LastConfigFromOrderer(reqCtx, channelID, order)
//	if err != nil {
//		return err
//	}
//
//}

func (gc *GenerateConfig) getSubject(csr *generate.CSR) pkix.Name {
	return pkix.Name{
		Country:            csr.Country,
		Organization:       csr.Organization,
		OrganizationalUnit: csr.OrganizationalUnit,
		Locality:           csr.Locality,
		Province:           csr.Province,
		StreetAddress:      csr.StreetAddress,
		PostalCode:         csr.PostalCode,
		SerialNumber:       csr.SerialNumber,
		CommonName:         csr.CommonName,
	}
}

func (gc *GenerateConfig) generateCryptoLeagueCrt(league *generate.ReqCreateLeague, priKeyData []byte, path, certName string) error {
	priKey, err := gc.getPriKey(priKeyData, path)
	if nil != err {
		return err
	}
	if _, err := gnomon.CA().GenerateCertificateSelf(&gnomon.CertSelf{
		CertificateFilePath:   filepath.Join(path, certName),
		Subject:               gc.getSubject(league.Csr),
		ParentPrivateKey:      priKey,
		PublicKey:             priKey.Public(),
		NotAfterDays:          time.Now().Add(5000 * 24 * time.Hour),
		NotBeforeDays:         time.Now(),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		SignatureAlgorithm:    gc.SignAlgorithm(league.SignAlgorithm),
	}); nil != err {
		return err
	}
	return nil
}

func (gc *GenerateConfig) generateCryptoOrgChild(child *generate.OrgChild, nodePath string, isNode bool) error {
	var (
		signCertPath, signCertFileName, tlsCertFileName string
		err                                             error
	)
	// ca cert
	signCertPath = path.Join(nodePath, "msp", "signcerts")
	if isNode {
		signCertFileName = CertNodeCAName(child.OrgName, child.OrgDomain, child.Name)
	} else {
		signCertFileName = CertUserCAName(child.OrgName, child.OrgDomain, child.Name)
	}
	if err = gc.enroll(child, signCertPath, signCertFileName); nil != err {
		return err
	}
	// tls ca cert
	if isNode {
		tlsCertFileName = "server.crt"
	} else {
		tlsCertFileName = "client.crt"
	}
	if err = gc.generateCryptoOrgChildTlsCaCrt(child, nodePath, tlsCertFileName); nil != err {
		return err
	}
	return nil
}

func (gc *GenerateConfig) generateCryptoOrgChildTlsCaCrt(child *generate.OrgChild, nodePath, certName string) (err error) {
	var parentTlsCert *x509.Certificate
	_, _, tlsCaPath, tlsCertName := gc.getRootCA(child.LeagueDomain)
	if parentTlsCert, err = gnomon.CA().LoadCrtFromFP(filepath.Join(tlsCaPath, tlsCertName)); nil != err {
		return err
	}
	priTlsParentKey, pubTlsKey, err := gc.getCertKey(filepath.Join(tlsCaPath, GeneratePriKeyFileName), child.PubTlsData)
	if nil != err {
		return err
	}
	tlsCertPath := path.Join(nodePath, "tls")
	if _, err = gnomon.CA().GenerateCertificate(&gnomon.Cert{
		ParentCert: parentTlsCert,
		CertSelf: gnomon.CertSelf{
			CertificateFilePath:   filepath.Join(tlsCertPath, certName),
			Subject:               gc.getSubject(child.EnrollInfo.EnrollRequest.Name),
			ParentPrivateKey:      priTlsParentKey,
			PublicKey:             pubTlsKey,
			NotAfterDays:          time.Now().Add(5000 * 24 * time.Hour),
			NotBeforeDays:         time.Now(),
			BasicConstraintsValid: true,
			IsCA:                  false,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
			SignatureAlgorithm:    gc.SignAlgorithm(child.SignAlgorithm),
		},
	}); nil != err {
		return err
	}
	return nil
}

func (gc *GenerateConfig) enroll(child *generate.OrgChild, path, certFileName string) error {
	notAfter := time.Now().Add(time.Duration(child.EnrollInfo.NotAfter) * 24 * time.Hour)
	notBefore := time.Now().Add(time.Duration(child.EnrollInfo.NotBefore) * 24 * time.Hour)
	gcr := generateCertificateRequest{CR: string(child.EnrollInfo.CsrPem), EnrollRequest: *child.EnrollInfo.EnrollRequest}
	gcr.NotAfter = notAfter
	gcr.NotBefore = notBefore
	crm, err := json.Marshal(gcr)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, strings.Join([]string{child.EnrollInfo.FabricCaServerURL, "api/v1/enroll"}, "/"), bytes.NewBuffer(crm))
	if nil != err {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(child.EnrollInfo.EnrollRequest.EnrollID, child.EnrollInfo.EnrollRequest.Secret)

	httpClient := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		enrResp := new(enrollmentResponse)
		if err := json.Unmarshal(body, enrResp); err != nil {
			return err
		}
		if !enrResp.Success {
			return enrResp.error()
		}
		rawCert, err := base64.StdEncoding.DecodeString(enrResp.Result.Cert)
		if nil != err {
			return err
		}
		_, err = gnomon.File().Append(filepath.Join(path, certFileName), rawCert, true)
		return err
	}
	return fmt.Errorf("non 200 response: %v message is: %s", resp.StatusCode, string(body))
}

func (gc *GenerateConfig) stringToCert(data string) (*x509.Certificate, error) {
	rawCert, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	pemResult, _ := pem.Decode(rawCert)
	return x509.ParseCertificate(pemResult.Bytes)
}

func (gc *GenerateConfig) getPriKeyFromBytes(priKeyData []byte) (priKey interface{}) {
	var (
		priEccKey *ecdsa.PrivateKey
		err       error
	)
	if priEccKey, err = gnomon.CryptoECC().LoadPriPem(priKeyData); nil != err {
		var (
			priRsaKey *rsa.PrivateKey
			pks       gnomon.PKSCType
		)
		pks = gnomon.CryptoRSA().PKSC8()
		if priRsaKey, err = gnomon.CryptoRSA().LoadPri(priKeyData, pks); nil != err {
			pks = gnomon.CryptoRSA().PKSC1()
			if priRsaKey, err = gnomon.CryptoRSA().LoadPri(priKeyData, pks); nil != err {
				err = errors.New("private key is not support")
				return
			}
		}
		priKey = priRsaKey
	} else {
		priKey = priEccKey

	}
	return
}

func (gc *GenerateConfig) getPriKey(priKeyData []byte, path string) (priKey crypto.Signer, err error) {
	var priEccKey *ecdsa.PrivateKey
	if priEccKey, err = gnomon.CryptoECC().LoadPriPem(priKeyData); nil != err {
		var (
			priRsaKey *rsa.PrivateKey
			pks       gnomon.PKSCType
		)
		pks = gnomon.CryptoRSA().PKSC8()
		if priRsaKey, err = gnomon.CryptoRSA().LoadPri(priKeyData, pks); nil != err {
			pks = gnomon.CryptoRSA().PKSC1()
			if priRsaKey, err = gnomon.CryptoRSA().LoadPri(priKeyData, pks); nil != err {
				err = errors.New("private key is not support")
				return
			}
		}
		priKey = priRsaKey
		if err = gnomon.CryptoRSA().SavePriPem(priRsaKey, path, GeneratePriKeyFileName, x509.PEMCipherAES128, pks); nil != err {
			return
		}
	} else {
		priKey = priEccKey
		if err = gnomon.CryptoECC().SavePriPem(priEccKey, path, GeneratePriKeyFileName); nil != err {
			return
		}

	}
	return
}

func (gc *GenerateConfig) getCertKey(priParentKeyFilePath string, pubKeyData []byte) (priParentKey crypto.Signer, pubKey interface{}, err error) {
	if priParentKey, err = gnomon.CryptoECC().LoadPriPemFP(priParentKeyFilePath); nil != err {
		if priParentKey, err = gnomon.CryptoRSA().LoadPriFP(priParentKeyFilePath, gnomon.CryptoRSA().PKSC8()); nil != err {
			if priParentKey, err = gnomon.CryptoRSA().LoadPriFP(priParentKeyFilePath, gnomon.CryptoRSA().PKSC1()); nil != err {
				err = errors.New("private key is not support")
				return
			}
		}
	}
	if pubKey, err = gnomon.CryptoECC().LoadPubPem(pubKeyData); nil != err {
		if pubKey, err = gnomon.CryptoRSA().LoadPub(pubKeyData); nil != err {
			err = errors.New("public key is not support")
			return
		}
	}
	return
}

func (gc *GenerateConfig) getRootCA(leagueDomain string) (caPath, caFileName, tlsCaPath, tlsCaFileName string) {
	caPath = CryptoRootCAPath(leagueDomain)
	caFileName = CertRootCAName(leagueDomain)
	tlsCaPath = CryptoRootTLSCAPath(leagueDomain)
	tlsCaFileName = CertRootTLSCAName(leagueDomain)
	return
}

func (gc *GenerateConfig) orgChildExec(leagueDomain, orgPath, nodePath string, ccn ClientCANode) error {
	var err error
	if !gnomon.File().PathExists(orgPath) {
		return errors.New("org done't exist")
	}
	if gnomon.File().PathExists(nodePath) {
		return errors.New("node or user already exist")
	}
	if err = os.MkdirAll(nodePath, 0755); nil != err {
		return err
	}

	childMspPath := path.Join(nodePath, "msp")
	if err = os.Mkdir(childMspPath, 0755); nil != err {
		return err
	}

	childMspAdminCertsPath := path.Join(childMspPath, "admincerts")
	childMspCaCertsPath := path.Join(childMspPath, "cacerts")
	childMspKeyStorePath := path.Join(childMspPath, "keystore")
	childMspSignCertsPath := path.Join(childMspPath, "signcerts")
	childMspTlsCaCertsPath := path.Join(childMspPath, "tlscacerts")
	if err = os.Mkdir(childMspAdminCertsPath, 0755); nil != err {
		return err
	}
	if err = os.Mkdir(childMspCaCertsPath, 0755); nil != err {
		return err
	}
	if err = os.Mkdir(childMspKeyStorePath, 0755); nil != err {
		return err
	}
	if err = os.Mkdir(childMspSignCertsPath, 0755); nil != err {
		return err
	}
	if err = os.Mkdir(childMspTlsCaCertsPath, 0755); nil != err {
		return err
	}

	childTlsPath := path.Join(nodePath, "tls")
	if err = os.Mkdir(childTlsPath, 0755); nil != err {
		return err
	}

	return gc.orgChildCopy(leagueDomain, orgPath, nodePath, childMspAdminCertsPath, childMspCaCertsPath, childMspTlsCaCertsPath, ccn)
}

func (gc *GenerateConfig) orgChildCopy(leagueDomain, orgPath, nodePath, childMspAdminCertsPath, childMspCaCertsPath,
	childMspTlsCaCertsPath string, ccn ClientCANode) error {

	var err error
	caPath, certName, tlsCaPath, tlsCertName := gc.getRootCA(leagueDomain)
	rootCaCertFilePath := filepath.Join(caPath, certName)
	caCertsFilePath := filepath.Join(childMspCaCertsPath, certName)
	if _, err = gnomon.File().Copy(rootCaCertFilePath, caCertsFilePath); nil != err {
		return err
	}
	rootTlsCaCertFilePath := filepath.Join(tlsCaPath, tlsCertName)
	tlsCaCertsFilePath := filepath.Join(childMspTlsCaCertsPath, tlsCertName)
	if _, err = gnomon.File().Copy(rootTlsCaCertFilePath, tlsCaCertsFilePath); nil != err {
		return err
	}

	tlsCertFilePath := filepath.Join(nodePath, "tls", "ca.crt")
	if _, err = gnomon.File().Copy(rootTlsCaCertFilePath, tlsCertFilePath); nil != err {
		return err
	}

	switch ccn {
	default:
		return nil
	case CcnNode:
		orgAdminCertsPath := path.Join(orgPath, "msp", "admincerts")
		var fileNames []string
		if fileNames, err = gnomon.File().LoopFileNames(orgAdminCertsPath); nil != err {
			return err
		}
		for _, fileName := range fileNames {
			orgAdminCertsFilePath := filepath.Join(orgAdminCertsPath, fileName)
			childMspAdminCertsFilePath := filepath.Join(childMspAdminCertsPath, fileName)
			if _, err = gnomon.File().Copy(orgAdminCertsFilePath, childMspAdminCertsFilePath); nil != err {
				return err
			}
		}
	}
	return nil
}

// SKI returns the subject key identifier of this key.
func (gc *GenerateConfig) EccSKI(key *ecdsa.PrivateKey) []byte {
	if key == nil {
		return nil
	}

	// Marshall the public key
	raw := elliptic.Marshal(key.Curve, key.PublicKey.X, key.PublicKey.Y)

	// Hash it
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKeyASN struct {
	N *big.Int
	E int
}

// SKI returns the subject key identifier of this key.
func (gc *GenerateConfig) RsaSKI(key *rsa.PrivateKey) []byte {
	if key == nil {
		return nil
	}

	// Marshall the public key
	raw, _ := asn1.Marshal(rsaPublicKeyASN{
		N: key.N,
		E: key.E,
	})

	// Hash it
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}
