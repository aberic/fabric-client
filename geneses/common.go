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
	"github.com/ennoo/rivet/utils/env"
	"strings"
)

var (
	// WorkDir 项目工作目录
	WorkDir string
	// FabricCryptoGenPathV14 1.4版本cryptogen二进制文件路径
	FabricCryptoGenPathV14 string
	// FabricConfigTXGenPathV14 1.4版本configtxgen二进制文件路径
	FabricConfigTXGenPathV14 string
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
	WorkDir = env.GetEnvDefault(WorkPath, strings.Join([]string{goPath, "src/github.com/ennoo/fabric-client/geneses/example"}, "/"))
	binDir := env.GetEnvDefault(BinPath, strings.Join([]string{WorkDir, "../../example/bin"}, "/"))
	FabricCryptoGenPathV14 = strings.Join([]string{binDir, "1.4/cryptogen"}, "/")
	FabricConfigTXGenPathV14 = strings.Join([]string{binDir, "1.4/configtxgen"}, "/")
}

func FabricCryptoGenPath(version string) string {
	switch version {
	case "1.4":
		return FabricCryptoGenPathV14
	default:
		return ""
	}
}

func FabricConfigTXGenPath(version string) string {
	switch version {
	case "1.4":
		return FabricConfigTXGenPathV14
	default:
		return ""
	}
}

// CryptoGenYmlPath cryptogen.yaml
func CryptoGenYmlPath(leagueComment string) string {
	return strings.Join([]string{WorkDir, "chain", leagueComment, "config/cryptogen.yaml"}, "/")
}

// ConfigTxYmlPath configtx.yaml
func ConfigTxYmlPath(leagueComment string) string {
	return strings.Join([]string{WorkDir, "chain", leagueComment, "config/configtx.yaml"}, "/")
}

// ConfPath crypto-config和channel-artifacts的上一级目录
func ConfPath(leagueComment string) string {
	return strings.Join([]string{WorkDir, "chain", leagueComment, "config"}, "/")
}

// CryptoConfigPath crypto-config目录
func CryptoConfigPath(leagueComment string) string {
	return strings.Join([]string{ConfPath(leagueComment), "crypto-config"}, "/")
}

// ChannelArtifactsPath channel-artifacts目录
func ChannelArtifactsPath(leagueComment string) string {
	return strings.Join([]string{ConfPath(leagueComment), "channel-artifacts"}, "/")
}

// GenesisBlockFilePath orderer.genesis.block路径
func GenesisBlockFilePath(leagueComment string) string {
	return strings.Join([]string{ChannelArtifactsPath(leagueComment), "orderer.genesis.block"}, "/")
}

// ChannelTXFilePath 通道tx文件路径
func ChannelTXFilePath(leagueComment, channelName string) string {
	return strings.Join([]string{ChannelArtifactsPath(leagueComment), "/", channelName, ".tx"}, "")
}
