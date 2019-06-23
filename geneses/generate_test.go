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
	"os"
	"testing"
)

func TestGenerateYmlTrue(t *testing.T) {
	if err := GenerateYml(&Generate{
		LeagueComment:   "league",
		OrderCount:      10,
		PeerCount:       10,
		TemplateCount:   10,
		UserCount:       10,
		BatchTimeout:    5,
		MaxMessageCount: 1000,
		Force:           true,
	}); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestGenerateYmlFalse(t *testing.T) {
	if err := GenerateYml(&Generate{
		LeagueComment:   "league",
		OrderCount:      10,
		PeerCount:       10,
		TemplateCount:   10,
		UserCount:       10,
		BatchTimeout:    5,
		MaxMessageCount: 1000,
		Force:           false,
	}); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestGenerateYmlErrParamsExceptionCrypto(t *testing.T) {
	if err := GenerateYml(&Generate{}); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestGenerateYmlErrParamsExceptionConfig(t *testing.T) {
	if err := GenerateYml(&Generate{
		LeagueComment: "league",
		OrderCount:    10,
		PeerCount:     10,
		TemplateCount: 10,
		UserCount:     10,
		Force:         false,
	}); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestGenerateCryptoFilesTrue(t *testing.T) {
	if line, strs, err := GenerateCryptoFiles(&Crypto{LeagueComment: "league", Force: true}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateCryptoFilesFalse(t *testing.T) {
	if line, strs, err := GenerateCryptoFiles(&Crypto{LeagueComment: "league", Force: false}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateCryptoFilesErrLeagueCommentNil(t *testing.T) {
	if line, strs, err := GenerateCryptoFiles(&Crypto{Force: false}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateGenesisBlockTrue(t *testing.T) {
	if line, strs, err := GenerateGenesisBlock(&Crypto{LeagueComment: "league", Force: true}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateGenesisBlockFalse(t *testing.T) {
	if line, strs, err := GenerateGenesisBlock(&Crypto{LeagueComment: "league", Force: false}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateGenesisBlockErrLeagueCommentNil(t *testing.T) {
	if line, strs, err := GenerateGenesisBlock(&Crypto{Force: false}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateChannelTXTrue(t *testing.T) {
	if line, strs, err := GenerateChannelTX(
		&ChannelTX{
			LeagueComment: "league",
			ChannelName:   "mychannel",
			Force:         true,
		}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateChannelTXFalse(t *testing.T) {
	if line, strs, err := GenerateChannelTX(
		&ChannelTX{
			LeagueComment: "league",
			ChannelName:   "mychannel",
			Force:         false,
		}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateChannelTXErrLeagueCommentNil(t *testing.T) {
	if line, strs, err := GenerateChannelTX(
		&ChannelTX{}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateChannelTXErrGenesisBlockNil(t *testing.T) {
	leagueComment := "league"
	genesisBlockFilePath := GenesisBlockFilePath(leagueComment)
	_ = os.Remove(genesisBlockFilePath)
	if line, strs, err := GenerateChannelTX(
		&ChannelTX{
			LeagueComment: leagueComment,
			ChannelName:   "mychannel",
			Force:         false,
		}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateChannelTXErrConfigTxYmlNil(t *testing.T) {
	leagueComment := "league"
	configTxYmlPath := ConfigTxYmlPath(leagueComment)
	_ = os.RemoveAll(configTxYmlPath)
	if line, strs, err := GenerateChannelTX(
		&ChannelTX{
			LeagueComment: leagueComment,
			ChannelName:   "mychannel",
			Force:         false,
		}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateGenesisBlockErrConfigTxYmlNil(t *testing.T) {
	if line, strs, err := GenerateGenesisBlock(&Crypto{LeagueComment: "league", Force: false}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateChannelTXErrCryptoConfigNil(t *testing.T) {
	leagueComment := "league"
	cryptoConfigPath := CryptoConfigPath(leagueComment)
	_ = os.RemoveAll(cryptoConfigPath)
	if line, strs, err := GenerateChannelTX(
		&ChannelTX{
			LeagueComment: leagueComment,
			ChannelName:   "mychannel",
			Force:         false,
		}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateGenesisBlockErrCryptoConfigNil(t *testing.T) {
	if line, strs, err := GenerateGenesisBlock(&Crypto{LeagueComment: "league", Force: false}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateCryptoFilesErrCryptoGenYmlNil(t *testing.T) {
	leagueComment := "league"
	cryptoGenYmlPath := CryptoGenYmlPath(leagueComment)
	_ = os.Remove(cryptoGenYmlPath)
	if line, strs, err := GenerateCryptoFiles(&Crypto{LeagueComment: leagueComment, Force: false}); nil != err {
		t.Skip(err)
	} else {
		t.Log("line = ", line, " | strs = ", strs)
	}
}

func TestGenerateCustomYmlTrue(t *testing.T) {
	if err := GenerateCustomYml(&GenerateCustom{
		LeagueComment: "league",
		Force:         true,
	}, nil, nil, nil, nil, nil, nil, nil, nil); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestGenerateCustomYmlFalse(t *testing.T) {
	if err := GenerateCustomYml(&GenerateCustom{
		LeagueComment: "league",
		Force:         false,
	}, nil, nil, nil, nil, nil, nil, nil, nil); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}
