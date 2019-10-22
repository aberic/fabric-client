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
 *
 */

package sdk

import (
	"encoding/hex"
	"github.com/aberic/gnomon"
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	ctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
)

func queryLedgerInfo(peerName string, channelProvider ctx.ChannelProvider) *Result {
	result := Result{}
	// Ledger client
	if client, err := ledger.New(channelProvider); nil != err {
		result.FailErr(err)
	} else {
		var ledgerInfo *fab.BlockchainInfoResponse
		if gnomon.String().IsEmpty(peerName) {
			ledgerInfo, err = client.QueryInfo()
		} else {
			ledgerInfo, err = client.QueryInfo(ledger.WithTargetEndpoints(peerName))
		}
		if nil != err {
			result.FailErr(err)
		} else {
			info := &pb.ChannelInfo{
				Endorser: ledgerInfo.Endorser,
				Status:   ledgerInfo.Status,
				Bci: &pb.BCI{
					Height:            ledgerInfo.BCI.Height,
					CurrentBlockHash:  hex.EncodeToString(ledgerInfo.BCI.CurrentBlockHash),
					PreviousBlockHash: hex.EncodeToString(ledgerInfo.BCI.PreviousBlockHash),
				},
			}
			result.Success(info)
		}
	}
	return &result
}

func queryLedgerBlockByHeight(peerName string, height uint64, channelProvider ctx.ChannelProvider) *Result {
	result := Result{}
	// Ledger client
	if client, err := ledger.New(channelProvider); nil != err {
		result.FailErr(err)
	} else {
		var commonBlock *common.Block
		if gnomon.String().IsEmpty(peerName) {
			commonBlock, err = client.QueryBlock(height)
		} else {
			commonBlock, err = client.QueryBlock(height, ledger.WithTargetEndpoints(peerName))
		}
		if nil != err {
			result.FailErr(err)
		} else {
			if block, err := parseBlock(commonBlock); nil != err {
				result.FailErr(err)
			} else {
				es := block.Envelopes
				for i := 0; i < len(es); {
					if nil == es[i] {
						es = append(es[:i], es[i+1:]...)
					} else {
						i++
					}
				}
				block.Envelopes = es
				result.Success(block)
			}
		}
	}
	return &result
}

func queryLedgerBlockByHash(peerName string, hash string, channelProvider ctx.ChannelProvider) *Result {
	result := Result{}
	// Ledger client
	if client, err := ledger.New(channelProvider); nil != err {
		result.FailErr(err)
	} else {
		if realHash, err := hex.DecodeString(hash); nil != err {
			result.FailErr(err)
		} else {
			var commonBlock *common.Block
			if gnomon.String().IsEmpty(peerName) {
				commonBlock, err = client.QueryBlockByHash(realHash)
			} else {
				commonBlock, err = client.QueryBlockByHash(realHash, ledger.WithTargetEndpoints(peerName))
			}
			if nil != err {
				result.FailErr(err)
			} else {
				if block, err := parseBlock(commonBlock); nil != err {
					result.FailErr(err)
				} else {
					result.Success(block)
				}
			}
		}
	}
	return &result
}

func queryLedgerBlockByTxID(peerName string, txID string, channelProvider ctx.ChannelProvider) *Result {
	result := Result{}
	// Ledger client
	if client, err := ledger.New(channelProvider); nil != err {
		result.FailErr(err)
	} else {
		var commonBlock *common.Block
		if gnomon.String().IsEmpty(peerName) {
			commonBlock, err = client.QueryBlockByTxID(fab.TransactionID(txID))
		} else {
			commonBlock, err = client.QueryBlockByTxID(fab.TransactionID(txID), ledger.WithTargetEndpoints(peerName))
		}
		if nil != err {
			result.FailErr(err)
		} else {
			if block, err := parseBlock(commonBlock); nil != err {
				result.FailErr(err)
			} else {
				result.Success(block)
			}
		}
	}
	return &result
}

func queryLedgerTransaction(peerName string, txID string, channelProvider ctx.ChannelProvider) *Result {
	result := Result{}
	// Ledger client
	if client, err := ledger.New(channelProvider); nil != err {
		result.FailErr(err)
	} else {
		if processedTransaction, err := client.QueryTransaction(fab.TransactionID(txID), ledger.WithTargetEndpoints(peerName)); nil != err {
			result.FailErr(err)
		} else {
			result.Success(processedTransaction)
		}
	}
	return &result
}

func queryLedgerConfig(peerName string, channelProvider ctx.ChannelProvider) *Result {
	result := Result{}
	// Ledger client
	if client, err := ledger.New(channelProvider); nil != err {
		result.FailErr(err)
	} else {
		if channelCfg, err := client.QueryConfig(ledger.WithTargetEndpoints(peerName)); nil != err {
			result.FailErr(err)
		} else {
			result.Success(channelCfg)
		}
	}
	return &result
}
