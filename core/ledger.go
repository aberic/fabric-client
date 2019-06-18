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
	"github.com/ennoo/rivet/trans/response"
	ctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/context"
	ch "github.com/hyperledger/fabric-sdk-go/pkg/fab/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/comm"
	"time"
)

type Info struct {
	BCI      BCI
	Endorser string
	Status   int32
}

type BCI struct {
	height            int64
	currentBlockHash  string
	previousBlockHash string
}

func queryInfo(channelID, peerName string, client ctx.Client) *response.Result {
	result := response.Result{}
	var (
		ledger  *ch.Ledger
		peerCfg *fab.NetworkPeer
		peerFab fab.Peer
		res     []*fab.BlockchainInfoResponse
		err     error
	)
	if ledger, err = ch.NewLedger(channelID); nil != err {
		goto ERR
	} else {
		reqCtx, cancel := context.NewRequest(client, context.WithTimeout(10*time.Second))
		defer cancel()
		if peerCfg, err = comm.NetworkPeerConfig(client.EndpointConfig(), peerName); nil != err {
			goto ERR
		}
		if peerFab, err = client.InfraProvider().CreatePeerFromConfig(peerCfg); nil != err {
			goto ERR
		}
		if res, err = ledger.QueryInfo(reqCtx, []fab.ProposalProcessor{peerFab}, nil); nil != err {
			goto ERR
		}
		result.Success(res)
		return &result
	}
ERR:
	result.Fail(err.Error())
	return &result
}

func queryBlockByHeight(channelID, peerName string, height uint64, client ctx.Client) *response.Result {
	result := response.Result{}
	var (
		ledger  *ch.Ledger
		peerCfg *fab.NetworkPeer
		peerFab fab.Peer
		err     error
	)
	if ledger, err = ch.NewLedger(channelID); nil != err {
		goto ERR
	} else {
		reqCtx, cancel := context.NewRequest(client, context.WithTimeout(10*time.Second))
		defer cancel()
		if peerCfg, err = comm.NetworkPeerConfig(client.EndpointConfig(), peerName); nil != err {
			goto ERR
		}
		if peerFab, err = client.InfraProvider().CreatePeerFromConfig(peerCfg); nil != err {
			goto ERR
		}
		res, err := ledger.QueryBlock(reqCtx, height, []fab.ProposalProcessor{peerFab}, nil)
		if nil != err {
			result.Fail(err.Error())
			return &result
		}
		//for _, d := range res[0].Data.Data {
		//	if nil == d {
		//		continue
		//	}
		//	if envelope, err := utils.GetEnvelopeFromBlock(d);nil!=err {
		//		log.Self.Error("error", log.Error(err))
		//	} else {
		//		log.Self.Info("envelope", log.Reflect("envelope", envelope))
		//	}
		//
		//}
		result.Success(res)
		return &result
	}
ERR:
	result.Fail(err.Error())
	return &result
}

func queryBlockByHash(channelID, peerName, hash string, client ctx.Client) *response.Result {
	result := response.Result{}
	var (
		ledger  *ch.Ledger
		peerCfg *fab.NetworkPeer
		peerFab fab.Peer
		err     error
	)
	if ledger, err = ch.NewLedger(channelID); nil != err {
		goto ERR
	} else {
		reqCtx, cancel := context.NewRequest(client, context.WithTimeout(10*time.Second))
		defer cancel()
		if peerCfg, err = comm.NetworkPeerConfig(client.EndpointConfig(), peerName); nil != err {
			goto ERR
		}
		if peerFab, err = client.InfraProvider().CreatePeerFromConfig(peerCfg); nil != err {
			goto ERR
		}
		res, err := ledger.QueryBlockByHash(reqCtx, []byte(hash), []fab.ProposalProcessor{peerFab}, nil)
		result.Success(res)
		if nil != err {
			result.Fail(err.Error())
		}
		return &result
	}
ERR:
	result.Fail(err.Error())
	return &result
}

func queryBlockByTxID(channelID, peerName, txID string, client ctx.Client) *response.Result {
	result := response.Result{}
	var (
		ledger  *ch.Ledger
		peerCfg *fab.NetworkPeer
		peerFab fab.Peer
		res     interface{}
		err     error
	)
	if ledger, err = ch.NewLedger(channelID); nil != err {
		goto ERR
	} else {
		reqCtx, cancel := context.NewRequest(client, context.WithTimeout(10*time.Second))
		defer cancel()
		if peerCfg, err = comm.NetworkPeerConfig(client.EndpointConfig(), peerName); nil != err {
			goto ERR
		}
		if peerFab, err = client.InfraProvider().CreatePeerFromConfig(peerCfg); nil != err {
			goto ERR
		}
		if res, err = ledger.QueryBlockByTxID(reqCtx, fab.TransactionID(txID), []fab.ProposalProcessor{peerFab}, nil); nil != err {
			goto ERR
		}
		result.Success(res)
		return &result
	}
ERR:
	result.Fail(err.Error())
	return &result
}

func queryTransaction(channelID, peerName, txID string, client ctx.Client) *response.Result {
	result := response.Result{}
	var (
		ledger  *ch.Ledger
		peerCfg *fab.NetworkPeer
		peerFab fab.Peer
		res     interface{}
		err     error
	)
	if ledger, err = ch.NewLedger(channelID); nil != err {
		goto ERR
	} else {
		reqCtx, cancel := context.NewRequest(client, context.WithTimeout(10*time.Second))
		defer cancel()
		if peerCfg, err = comm.NetworkPeerConfig(client.EndpointConfig(), peerName); nil != err {
			goto ERR
		}
		if peerFab, err = client.InfraProvider().CreatePeerFromConfig(peerCfg); nil != err {
			goto ERR
		}
		if res, err = ledger.QueryTransaction(reqCtx, fab.TransactionID(txID), []fab.ProposalProcessor{peerFab}, nil); nil != err {
			goto ERR
		}
		result.Success(res)
		return &result
	}
ERR:
	result.Fail(err.Error())
	return &result
}
