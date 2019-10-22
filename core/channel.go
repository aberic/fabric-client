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
	"github.com/aberic/gnomon"
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
	"github.com/pkg/errors"
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
		gnomon.Log().Error("channelClient", gnomon.Log().Err(err))
		return nil
	}
	return client
}

// channelConfigPath mychannel.tx
func createChannel(orderURL, orgName, orgUser, channelID, channelConfigPath string, sdk *fabsdk.FabricSDK,
	client *resmgmt.Client) (txID string, err error) {
	var (
		mspClient     *mspclient.Client
		adminIdentity msp.SigningIdentity
		scResp        resmgmt.SaveChannelResponse
	)
	mspClient, err = mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
	if err != nil {
		return
	}
	adminIdentity, err = mspClient.GetSigningIdentity(orgUser)
	if err != nil {
		return
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: channelConfigPath,
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	scResp, err = client.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(orderURL))
	if err != nil {
		gnomon.Log().Error("createChannel", gnomon.Log().Err(err))
		return "", errors.Errorf("error should be nil. %v", err)
	}
	return string(scResp.TransactionID), nil
}

// ordererUrl "orderer.example.com"
func joinChannel(orderURL, channelID, peerName string, client *resmgmt.Client) error {
	// Org peers join channel
	if err := client.JoinChannel(
		channelID,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(peerName),
		resmgmt.WithOrdererEndpoint(orderURL)); err != nil {
		gnomon.Log().Error("joinChannel", gnomon.Log().Err(err))
		return errors.Errorf("Org peers failed to JoinChannel:  %v", err)
	}
	return nil
}

type ChannelArr struct {
	Channels []*peer2.ChannelInfo `json:"channels"`
}

// peer 参见peer.go PeerChannel
func queryChannels(orgName, orgUser, peerName string, sdk *fabsdk.FabricSDK) ([]*peer2.ChannelInfo, error) {
	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		gnomon.Log().Error("queryChannels", gnomon.Log().Err(err))
		return nil, errors.Errorf("Failed to query channels:  %v", err)
	} else {
		if nil != orgResMgmt {
			qcResponse, err := orgResMgmt.QueryChannels(resmgmt.WithTargetEndpoints(peerName))
			if err != nil {
				gnomon.Log().Error("queryChannels", gnomon.Log().Err(err))
				return nil, errors.Errorf("Failed to query channels: peer cannot be nil.  %v", err)
			}
			if nil == qcResponse {
				gnomon.Log().Error("queryChannels", gnomon.Log().Err(err))
				return nil, errors.Errorf("qcResponse error should be nil. ")
			} else {
				return qcResponse.Channels, nil
			}
		} else {
			gnomon.Log().Error("queryChannels", gnomon.Log().Err(err))
			return nil, errors.Errorf("orgResMgmt error should be nil. ")
		}
	}
}

func queryChannelInfo(channelID, peerName string, client ctx.Client) *Result {
	result := Result{}
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

func queryChannelBlockByHeight(channelID, peerName string, height uint64, client ctx.Client) *Result {
	result := Result{}
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
		//		gnomon.Log().Error("error", log.Error(err))
		//	} else {
		//		gnomon.Log().Info("envelope", log.Reflect("envelope", envelope))
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

func queryChannelBlockByHash(channelID, peerName, hash string, client ctx.Client) *Result {
	result := Result{}
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

func queryChannelBlockByTxID(channelID, peerName, txID string, client ctx.Client) *Result {
	result := Result{}
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

func queryChannelTransaction(channelID, peerName, txID string, client ctx.Client) *Result {
	result := Result{}
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
