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
	"github.com/aberic/gnomon"
	"path/filepath"
	"strings"
)

var (
	// WorkDir 项目工作目录
	dataPath string
)

// 环境变量
const (
	// WorkPath 项目工作目录 [template]
	WorkPath       = "WORK_PATH"
	OrdererOrgName = "OrdererOrg"
	GOPath         = "GOPATH"
)

type ClientCANode int

const (
	CcnNode ClientCANode = iota
	CcnAdmin
	CcnUser
)

func init() {
	dataPath = strings.Join([]string{
		gnomon.Env().GetD(
			WorkPath,
			strings.Join([]string{gnomon.Env().Get(GOPath), "src/github.com/aberic/fabric-client/geneses/example"}, "/"),
		),
		"data",
	}, "/")
}

func MspID(orgName string) string {
	return strings.Join([]string{orgName, "MSP"}, "")
}

func NodeDomain(orgName, orgDomain, nodeName string) string {
	return strings.Join([]string{nodeName, orgName, orgDomain}, ".")
}

// CertificateAuthorityPath CertificateAuthorityPath
func CertificateAuthorityFilePath(leagueDomain, caName string) string {
	return filepath.Join(dataPath, leagueDomain, "ca", caName, "cert.pem")
}

// CertificateAuthorityClientKeyFilePath CertificateAuthorityPath
func CertificateAuthorityClientKeyFilePath(leagueDomain, caName string) string {
	return filepath.Join(dataPath, leagueDomain, "ca", caName, "client.key")
}

// CertificateAuthorityClientCertFilePath CertificateAuthorityPath
func CertificateAuthorityClientCertFilePath(leagueDomain, caName string) string {
	return filepath.Join(dataPath, leagueDomain, "ca", caName, "client.crt")
}

// CryptoRootCAPath CryptoCAPath
func CryptoRootCAPath(leagueDomain string) string {
	return filepath.Join(dataPath, leagueDomain, "crypto-config", "root", "ca")
}

// CryptoRootTLSCAPath CryptoCAPath
func CryptoRootTLSCAPath(leagueDomain string) string {
	return filepath.Join(dataPath, leagueDomain, "crypto-config", "root", "tlsca")
}

func CertRootCAName(leagueDomain string) string {
	return strings.Join([]string{"ca.", leagueDomain, "-cert.pem"}, "")
}

func CertRootTLSCAName(leagueDomain string) string {
	return strings.Join([]string{"tlsca.", leagueDomain, "-cert.pem"}, "")
}

func CertNodeCAName(orgName, orgDomain, nodeName string) string {
	return strings.Join([]string{nodeName, ".", orgName, ".", orgDomain, "-cert.pem"}, "")
}

func CertUserCAName(orgName, orgDomain, userName string) string {
	return strings.Join([]string{userName, "@", orgName, ".", orgDomain, "-cert.pem"}, "")
}

func CsrPath(leagueDomain, orgName, orgDomain string) string {
	return filepath.Join(dataPath, leagueDomain, "csr", strings.Join([]string{orgName, orgDomain}, "."))
}

func CsrFilePath(leagueDomain, orgName, orgDomain, commonName string) string {
	fileName := strings.Join([]string{commonName, "csr"}, ".")
	return filepath.Join(dataPath, leagueDomain, "csr", strings.Join([]string{orgName, orgDomain}, "."), fileName)
}

// CryptoOrgAndNodePath CryptoOrgAndNodePath
func CryptoOrgAndNodePath(leagueDomain, orgDomain, orgName, nodeName string, isPeer bool, node ClientCANode) (orgPath, nodePath string) {
	var orgsName, orgPathName, nodesName, nodePathName string
	if isPeer {
		orgsName = "peerOrganizations/"
		if node == CcnNode {
			nodesName = "peers"
			nodePathName = strings.Join([]string{nodeName, orgName, orgDomain}, ".")
		} else {
			nodesName = "users"
			nodePathName = strings.Join([]string{nodeName, "@", orgName, ".", orgDomain}, "")
		}
	} else {
		orgsName = "ordererOrganizations/"
		if node == CcnNode {
			nodesName = "orderers"
			nodePathName = strings.Join([]string{nodeName, orgName, orgDomain}, ".")
		} else {
			nodesName = "users"
			nodePathName = strings.Join([]string{nodeName, "@", orgName, ".", orgDomain}, "")
		}
	}
	orgPathName = strings.Join([]string{orgsName, orgName, ".", orgDomain}, "")
	orgPath = filepath.Join(dataPath, leagueDomain, "crypto-config", orgPathName)
	nodePath = filepath.Join(orgPath, nodesName, nodePathName)
	return
}

// CryptoOrgMspPath CryptoOrgMspPath
func CryptoOrgMspPath(leagueDomain, orgDomain, orgName string, isPeer bool) (mspPath string) {
	var orgsName, orgPathName string
	if isPeer {
		orgsName = "peerOrganizations/"
	} else {
		orgsName = "ordererOrganizations/"
	}
	orgPathName = strings.Join([]string{orgsName, orgName, ".", orgDomain}, "")
	return filepath.Join(dataPath, leagueDomain, "crypto-config", orgPathName, "msp")
}

// CryptoUserTmpPath CryptoUserTempPath
func CryptoUserTmpPath(leagueDomain, orgDomain, orgName string) string {
	tmpPath := strings.Join([]string{"tmp/", orgName, ".", orgDomain, "/users"}, "")
	return filepath.Join(dataPath, leagueDomain, "crypto-config", tmpPath)
}

// ChainCodePath code目录
func ChainCodePath(leagueName, chainCodeName, version string) (source, path, zipPath string) {
	source = filepath.Join(dataPath, leagueName, "code/go")
	path = filepath.Join(chainCodeName, version, chainCodeName)
	zipPath = strings.Join([]string{source, "/src/", path, ".zip"}, "")
	return
}

// CryptoConfigPath crypto-config目录
func CryptoConfigPath(leagueName string) string {
	return filepath.Join(dataPath, leagueName, "crypto-config")
}

// ChannelArtifactsPath channel-artifacts目录
func ChannelArtifactsPath(leagueName string) string {
	return filepath.Join(dataPath, leagueName, "channel-artifacts")
}

// GenesisBlockFilePath orderer.genesis.block路径
func GenesisBlockFilePath(leagueName string) string {
	return filepath.Join(dataPath, leagueName, "channel-artifacts/orderer.genesis.block")
}

// ChannelTXFilePath 通道tx文件路径
func ChannelTXFilePath(leagueName, channelName string) string {
	return strings.Join([]string{ChannelArtifactsPath(leagueName), "/", channelName, ".tx"}, "")
}

// ChannelUpdateTXFilePath 通道tx文件路径
func ChannelUpdateTXFilePath(leagueName, channelName string) string {
	return strings.Join([]string{ChannelArtifactsPath(leagueName), "/", channelName, "_update.pb"}, "")
}
