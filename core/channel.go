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
	"github.com/ennoo/rivet/utils/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	ctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/context"
	ch "github.com/hyperledger/fabric-sdk-go/pkg/fab/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/comm"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	peer2 "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"time"
)

func channelProvider(orgName, orgUser, channelID string, sdk *fabsdk.FabricSDK) ctx.ChannelProvider {
	//prepare channel client context using client context
	return sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))
}

func channelClient(orgName, orgUser, channelID string, sdk *fabsdk.FabricSDK) *channel.Client {
	//prepare channel client context using client context
	clientChannelContext := channelProvider(orgName, orgUser, channelID, sdk)
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		log.Self.Error("Failed to create new channel client:" + err.Error())
		return nil
	}
	return client
}

// channelConfigPath mychannel.tx
func createChannel(orderURL, orgName, orgUser, channelID, channelConfigPath string, sdk *fabsdk.FabricSDK,
	client *resmgmt.Client) *response.Result {
	result := response.Result{}
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
	} else {
		adminIdentity, err := mspClient.GetSigningIdentity(orgUser)
		if err != nil {
			log.Self.Error(err.Error())
			result.Fail(err.Error())
		} else {
			req := resmgmt.SaveChannelRequest{ChannelID: channelID,
				ChannelConfigPath: channelConfigPath,
				SigningIdentities: []msp.SigningIdentity{adminIdentity}}
			txID, err := client.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(orderURL))
			if err != nil {
				log.Self.Error("error should be nil. " + err.Error())
				result.Fail("error should be nil. " + err.Error())
			} else {
				result.Success(string(txID.TransactionID))
			}
		}
	}
	return &result
}

// ordererUrl "orderer.example.com"
func joinChannel(orderURL, channelID string, client *resmgmt.Client) *response.Result {
	result := response.Result{}
	// Org peers join channel
	if err := client.JoinChannel(
		channelID,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(orderURL)); err != nil {
		log.Self.Error("Org peers failed to JoinChannel: " + err.Error())
		result.Fail("Org peers failed to JoinChannel: " + err.Error())
	} else {
		result.Success("success")
	}
	return &result
}

type ChannelArr struct {
	Channels []*peer2.ChannelInfo `json:"channels"`
}

// peer 参见peer.go PeerChannel
func queryChannels(orgName, orgUser, peerName string, sdk *fabsdk.FabricSDK) *response.Result {
	result := response.Result{}
	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		log.Self.Error("Failed to query channels: " + err.Error())
		result.Fail("Failed to query channels: " + err.Error())
	} else {
		if nil != orgResMgmt {
			qcResponse, err := orgResMgmt.QueryChannels(resmgmt.WithTargetEndpoints(peerName))
			if err != nil {
				log.Self.Error("Failed to query channels: peer cannot be nil", log.Error(err))
				result.Fail("Failed to query channels: peer cannot be nil")
			}
			if nil == qcResponse {
				log.Self.Error("qcResponse error should be nil. ")
				result.Fail("qcResponse error should be nil. ")
			} else {
				result.Success(&ChannelArr{qcResponse.Channels})
			}
		} else {
			log.Self.Error("orgResMgmt error should be nil. ")
			result.Fail("orgResMgmt error should be nil. ")
		}
	}
	return &result
}

func queryChannelInfo(channelID, peerName string, client ctx.Client) *response.Result {
	result := response.Result{}
	var (
		ledger1 *ch.Ledger
		peerCfg *fab.NetworkPeer
		peerFab fab.Peer
		res     []*fab.BlockchainInfoResponse
		err     error
	)
	if ledger1, err = ch.NewLedger(channelID); nil != err {
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
		if res, err = ledger1.QueryInfo(reqCtx, []fab.ProposalProcessor{peerFab}, nil); nil != err {
			goto ERR
		}
		result.Success(res)
		return &result
	}
ERR:
	result.Fail(err.Error())
	return &result
}

func queryChannelBlockByHeight(channelID, peerName string, height uint64, client ctx.Client) *response.Result {
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

func queryChannelBlockByHash(channelID, peerName, hash string, client ctx.Client) *response.Result {
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
		//resInfo, _ := ledger.QueryInfo(reqCtx, []fab.ProposalProcessor{peerFab}, nil)
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

func queryChannelBlockByTxID(channelID, peerName, txID string, client ctx.Client) *response.Result {
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

func queryChannelTransaction(channelID, peerName, txID string, client ctx.Client) *response.Result {
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
