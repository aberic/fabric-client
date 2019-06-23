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
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	"github.com/ennoo/rivet/trans/response"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	ctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

func queryLedgerInfo(channelProvider ctx.ChannelProvider) *response.Result {
	result := response.Result{}
	// Ledger client
	if client, err := ledger.New(channelProvider); nil != err {
		result.FailErr(err)
	} else {
		if ledgerInfo, err := client.QueryInfo(); nil != err {
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

func queryLedgerBlockByHeight(height uint64, channelProvider ctx.ChannelProvider) *response.Result {
	result := response.Result{}
	// Ledger client
	if client, err := ledger.New(channelProvider); nil != err {
		result.FailErr(err)
	} else {
		if commonBlock, err := client.QueryBlock(height); nil != err {
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

func queryLedgerBlockByHash(hash string, channelProvider ctx.ChannelProvider) *response.Result {
	result := response.Result{}
	// Ledger client
	if client, err := ledger.New(channelProvider); nil != err {
		result.FailErr(err)
	} else {
		if realHash, err := hex.DecodeString(hash); nil != err {
			result.FailErr(err)
		} else {
			if commonBlock, err := client.QueryBlockByHash(realHash); nil != err {
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

func queryLedgerBlockByTxID(txID string, channelProvider ctx.ChannelProvider) *response.Result {
	result := response.Result{}
	// Ledger client
	if client, err := ledger.New(channelProvider); nil != err {
		result.FailErr(err)
	} else {
		if commonBlock, err := client.QueryBlockByTxID(fab.TransactionID(txID)); nil != err {
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

func queryLedgerTransaction(txID string, channelProvider ctx.ChannelProvider) *response.Result {
	result := response.Result{}
	// Ledger client
	if client, err := ledger.New(channelProvider); nil != err {
		result.FailErr(err)
	} else {
		if processedTransaction, err := client.QueryTransaction(fab.TransactionID(txID)); nil != err {
			result.FailErr(err)
		} else {
			result.Success(processedTransaction)
		}
	}
	return &result
}

func queryLedgerConfig(channelProvider ctx.ChannelProvider) *response.Result {
	result := response.Result{}
	// Ledger client
	if client, err := ledger.New(channelProvider); nil != err {
		result.FailErr(err)
	} else {
		if channelCfg, err := client.QueryConfig(); nil != err {
			result.FailErr(err)
		} else {
			result.Success(channelCfg)
		}
	}
	return &result
}
