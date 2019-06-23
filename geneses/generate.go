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
	"errors"
	pb "github.com/ennoo/fabric-client/grpc/proto/geneses"
	"github.com/ennoo/rivet/utils/cmd"
	"github.com/ennoo/rivet/utils/file"
	"github.com/ennoo/rivet/utils/log"
	"github.com/ennoo/rivet/utils/string"
	"os"
	"strings"
)

const (
	// OrderPrefix 联盟配置文件对象排序服务默认前缀
	OrderPrefix = "orderer"
)

// GenerateYml 生成区块链配置YML文件
func GenerateYml(generate *pb.Generate) error {
	var (
		err              error
		cryptoGenData    []byte
		configTxData     []byte
		cryptoGenYmlPath string
		configTxYmlPath  string
	)
	if cryptoGenData, err = generateCryptoGenYml(generate.LedgerName, generate.OrderCount, generate.PeerCount,
		generate.TemplateCount, generate.UserCount); nil != err {
		goto ERR
	}
	if configTxData, err = generateConfigTXYml(generate.LedgerName, generate.OrderCount, generate.PeerCount, generate.BatchTimeout,
		generate.MaxMessageCount); nil != err {
		goto ERR
	}
	cryptoGenYmlPath = CryptoGenYmlPath(generate.LedgerName)
	configTxYmlPath = ConfigTxYmlPath(generate.LedgerName)
	if err = file.CreateAndWrite(cryptoGenYmlPath, cryptoGenData, generate.Force); nil != err {
		goto ERR
	}
	if err = file.CreateAndWrite(configTxYmlPath, configTxData, generate.Force); nil != err {
		goto ERR
	}
	return nil
ERR:
	//_ = os.Remove(cryptoGenPath)
	//_ = os.Remove(configTxPath)
	log.Self.Debug("generate", log.Error(err))
	return err
}

// GenerateCustomYml 生成区块链配置YML文件
func GenerateCustomYml(generateCustom *GenerateCustom, ordererOrgs []*Order, peerOrgs []*Peer,
	organizations []Organization, application *Application, capabilities *CapabilitiesAll, channel *Channel,
	orderer *Orderer, profiles *Profiles) error {
	var (
		err              error
		cryptoGenData    []byte
		configTxData     []byte
		cryptoGenYmlPath string
		configTxYmlPath  string
	)
	if cryptoGenData, err = generateCryptoGenCustomYml(ordererOrgs, peerOrgs); nil != err {
		goto ERR
	}
	if configTxData, err = generateConfigTXCustomYml(organizations, application, capabilities, channel, orderer,
		profiles); nil != err {
		goto ERR
	}
	cryptoGenYmlPath = CryptoGenYmlPath(generateCustom.LedgerName)
	configTxYmlPath = ConfigTxYmlPath(generateCustom.LedgerName)
	if err = file.CreateAndWrite(cryptoGenYmlPath, cryptoGenData, generateCustom.Force); nil != err {
		goto ERR
	}
	if err = file.CreateAndWrite(configTxYmlPath, configTxData, generateCustom.Force); nil != err {
		goto ERR
	}
	return nil
ERR:
	//_ = os.Remove(cryptoGenPath)
	//_ = os.Remove(configTxPath)
	log.Self.Debug("generate custom", log.Error(err))
	return err
}

// GenerateCryptoFiles 生成区块链配置文件集合
func GenerateCryptoFiles(crypto *pb.Crypto) (int, []string, error) {
	var (
		exist bool
		err   error
	)
	if str.IsEmpty(crypto.LedgerName) {
		return 0, nil, errors.New("league name can not be nil")
	}
	cryptoGenYmlPath := CryptoGenYmlPath(crypto.LedgerName)
	cryptoConfigPath := CryptoConfigPath(crypto.LedgerName)
	if exist, _ = file.PathExists(cryptoGenYmlPath); !exist {
		return 0, nil, errors.New("cryptogen.yaml have not bean created")
	}
	if exist, _ = file.PathExists(cryptoConfigPath); exist && !crypto.Force {
		return 0, nil, errors.New("crypto-config dir had already exist")
	}
	if err = os.RemoveAll(cryptoConfigPath); nil != err {
		return 0, nil, err
	}
	return utils.ExecCommandTail(
		FabricCryptoGenPathV14,
		"generate",
		strings.Join([]string{"--config=", cryptoGenYmlPath}, ""),
		strings.Join([]string{"--output=", cryptoConfigPath}, ""),
	)
}

// GenerateGenesisBlock 生成区块链创世区块
func GenerateGenesisBlock(crypto *pb.Crypto) (int, []string, error) {
	if str.IsEmpty(crypto.LedgerName) {
		return 0, nil, errors.New("league comment can not be nil")
	}
	var (
		exist bool
		err   error
	)
	cryptoConfigPath := CryptoConfigPath(crypto.LedgerName)
	confPath := ConfPath(crypto.LedgerName)
	configTxYmlPath := ConfigTxYmlPath(crypto.LedgerName)
	channelArtifactsPath := ChannelArtifactsPath(crypto.LedgerName)
	genesisBlockFilePath := GenesisBlockFilePath(crypto.LedgerName)
	if exist, _ = file.PathExists(cryptoConfigPath); !exist {
		return 0, nil, errors.New("crypto-config dir should be created first")
	}
	if exist, _ = file.PathExists(configTxYmlPath); !exist {
		return 0, nil, errors.New("configtx.yml file is not exist")
	}
	if err = os.MkdirAll(channelArtifactsPath, os.ModePerm); nil != err {
		return 0, nil, err
	}
	if exist, _ = file.PathExists(genesisBlockFilePath); exist && !crypto.Force {
		return 0, nil, errors.New("genesis block file already exist")
	} else if err = os.Remove(genesisBlockFilePath); nil != err && exist && crypto.Force {
		return 0, nil, err
	}
	return utils.ExecCommandTail(FabricConfigTXGenPathV14, "--configPath", confPath, "--profile",
		"HBaaSOrderGenesis", "--outputBlock", genesisBlockFilePath)
}

// GenerateChannelTX 生成区块链创世区块
func GenerateChannelTX(channelTX *pb.ChannelTX) (int, []string, error) {
	if str.IsEmpty(channelTX.LedgerName) {
		return 0, nil, errors.New("league comment can not be nil")
	}
	var (
		exist bool
		err   error
	)
	cryptoConfigPath := CryptoConfigPath(channelTX.LedgerName)
	confPath := ConfPath(channelTX.LedgerName)
	configTxYmlPath := ConfigTxYmlPath(channelTX.LedgerName)
	genesisBlockFilePath := GenesisBlockFilePath(channelTX.LedgerName)
	channelTXFilePath := ChannelTXFilePath(channelTX.LedgerName, channelTX.ChannelName)
	if exist, _ = file.PathExists(cryptoConfigPath); !exist {
		return 0, nil, errors.New("crypto-config dir should be created first")
	}
	if exist, _ = file.PathExists(configTxYmlPath); !exist {
		return 0, nil, errors.New("configtx.yml file is not exist")
	}
	if exist, _ = file.PathExists(genesisBlockFilePath); !exist {
		return 0, nil, errors.New("genesis block should be created first")
	}
	if exist, _ = file.PathExists(channelTXFilePath); exist && !channelTX.Force {
		return 0, nil, errors.New("channel tx file already exist")
	} else if err = os.Remove(channelTXFilePath); nil != err && exist && channelTX.Force {
		return 0, nil, err
	}
	return utils.ExecCommandTail(FabricConfigTXGenPathV14, "--configPath", confPath, "--profile",
		"HBaaSChannel", "--outputCreateChannelTx", channelTXFilePath, "-channelID", channelTX.ChannelName)
}
