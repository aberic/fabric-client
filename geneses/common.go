/*
 * Copyright (c) 2019.. ENNOO - All Rights Reserved.
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
	"fmt"
	pb "github.com/ennoo/fabric-client/grpc/proto/geneses"
	"github.com/ennoo/rivet/utils/env"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	// WorkDir 项目工作目录
	dataPath string
	// FabricCryptoGenPathV14 1.4版本cryptogen二进制文件路径
	fabricCryptoGenPathV14 string
	// FabricConfigTXGenPathV14 1.4版本configtxgen二进制文件路径
	fabricConfigTXGenPathV14 string
)

// 环境变量
const (
	// Fabric 二进制文件目录 [template]
	BinPath = "BIN_PATH"
	// WorkPath 项目工作目录 [template]
	WorkPath       = "WORK_PATH"
	OrdererOrgName = "OrdererOrg"
)

func init() {
	goPath := env.GetEnv(env.GOPath)
	dataPath = strings.Join([]string{
		env.GetEnvDefault(
			WorkPath,
			strings.Join([]string{goPath, "src/github.com/ennoo/fabric-client/geneses/example"}, "/"),
		),
		"data",
	}, "/")
	binDir := env.GetEnvDefault(BinPath,
		strings.Join([]string{goPath, "src/github.com/ennoo/fabric-client/bin"}, "/"))
	fabricCryptoGenPathV14 = strings.Join([]string{binDir, "1.4", systemDir(), "cryptogen"}, "/")
	fabricConfigTXGenPathV14 = strings.Join([]string{binDir, "1.4", systemDir(), "configtxgen"}, "/")
}

func systemDir() string {
	osStr := runtime.GOOS
	osArch := runtime.GOARCH
	fmt.Println(osStr, "-", osArch)
	if osArch != "amd64" {
		return ""
	}
	if osStr == "darwin" {
		return "darwin-amd64"
	} else if osStr == "linux" {
		return "linux-amd64"
	} else if osStr == "windows" {
		return "windows-amd64"
	} else {
		return ""
	}
}

func FabricCryptoGenPath(version pb.Version) string {
	switch version {
	case pb.Version_V14:
		return fabricCryptoGenPathV14
	default:
		return ""
	}
}

func FabricConfigTXGenPath(version pb.Version) string {
	switch version {
	case pb.Version_V14:
		return fabricConfigTXGenPathV14
	default:
		return ""
	}
}

// CryptoCAPath CryptoCAPath
//
// caType 0-ordererOrganizations;1-peerOrganizations
func CryptoCAPath(ledgerName, orgName string, caType int) string {
	var cPath string
	switch caType {
	default:
		cPath = strings.Join([]string{"peerOrganizations/", orgName, ".", ledgerName}, "")
	case 0:
		cPath = strings.Join([]string{"ordererOrganizations/", ledgerName}, "")
	}
	return filepath.Join(dataPath, ledgerName, "crypto-config", cPath, "ca")
}

// CryptoGenYmlPath cryptogen.yaml
func CryptoGenYmlPath(ledgerName string) string {
	return filepath.Join(dataPath, ledgerName, "config/cryptogen.yaml")
}

// ConfigTxYmlPath configtx.yaml
func ConfigTxYmlPath(ledgerName string) string {
	return filepath.Join(dataPath, ledgerName, "config/configtx.yaml")
}

// ConfPath crypto-config和channel-artifacts的上一级目录
func ConfPath(ledgerName string) string {
	return filepath.Join(dataPath, ledgerName, "config")
}

// ChainCodePath code目录
func ChainCodePath(ledgerName, chainCodeName, version string) (source, path, zipPath string) {
	source = filepath.Join(dataPath, ledgerName, "code/go")
	path = filepath.Join(chainCodeName, version, chainCodeName)
	zipPath = strings.Join([]string{source, "/src/", path, ".zip"}, "")
	return
}

// CryptoConfigPath crypto-config目录
func CryptoConfigPath(ledgerName string) string {
	return filepath.Join(dataPath, ledgerName, "crypto-config")
}

// ChannelArtifactsPath channel-artifacts目录
func ChannelArtifactsPath(ledgerName string) string {
	return filepath.Join(dataPath, ledgerName, "channel-artifacts")
}

// GenesisBlockFilePath orderer.genesis.block路径
func GenesisBlockFilePath(ledgerName string) string {
	return filepath.Join(dataPath, ledgerName, "channel-artifacts/orderer.genesis.block")
}

// ChannelTXFilePath 通道tx文件路径
func ChannelTXFilePath(ledgerName, channelName string) string {
	return strings.Join([]string{ChannelArtifactsPath(ledgerName), "/", channelName, ".tx"}, "")
}
