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
	pc "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/comm"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/discovery"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	fabdiscovery "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/discovery"
	"time"
)

func discoveryClientPeers(channelID, orgName, orgUser, peerName string, sdk *fabsdk.FabricSDK) *response.Result {
	result := response.Result{}
	var (
		contextClient pc.Client
		err           error
	)
	//prepare context
	ctx := sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))
	if contextClient, err = ctx(); nil != err {
		goto ERR
	}

	if disClient, err := discovery.New(contextClient); nil != err {
		goto ERR
	} else {
		reqCtx, cancel := context.NewRequest(contextClient, context.WithTimeout(10*time.Second))
		defer cancel()

		//req := discclient.NewRequest().OfChannel("mychannel").AddPeersQuery()

		peerCfg1, err := comm.NetworkPeerConfig(contextClient.EndpointConfig(), peerName)
		if nil != err {
			goto ERR
		}

		//func (c *Client) Req(channelID string) *discclient.Request {
		//	return discclient.NewRequest().OfChannel(channelID).AddPeersQuery()
		//}
		//
		//func (c *Client) ReqLocal() *discclient.Request {
		//	return discclient.NewRequest().AddLocalPeersQuery()
		//}
		//
		//func (c *Client) ReqConfig(channelID string) *discclient.Request {
		//	return discclient.NewRequest().OfChannel(channelID).AddConfigQuery()
		//}
		responses, err := disClient.Send(reqCtx, disClient.Req(channelID), peerCfg1.PeerConfig)
		if nil != err {
			goto ERR
		}
		resp := responses[0]
		chanResp := resp.ForChannel(channelID)

		peers, err := chanResp.Peers()
		if nil != err {
			goto ERR
		}

		result.Success(peers)
		return &result
	}

ERR:
	result.FailErr(err)
	return &result
}

func discoveryClientLocalPeers(orgName, orgUser, peerName string, sdk *fabsdk.FabricSDK) *response.Result {
	result := response.Result{}
	var (
		contextClient pc.Client
		err           error
	)
	//prepare context
	ctx := sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))
	if contextClient, err = ctx(); nil != err {
		goto ERR
	}

	if disClient, err := discovery.New(contextClient); nil != err {
		goto ERR
	} else {
		reqCtx, cancel := context.NewRequest(contextClient, context.WithTimeout(10*time.Second))
		defer cancel()

		//req := discclient.NewRequest().OfChannel("mychannel").AddPeersQuery()

		peerCfg1, err := comm.NetworkPeerConfig(contextClient.EndpointConfig(), peerName)
		if nil != err {
			goto ERR
		}
		responses, err := disClient.Send(reqCtx, disClient.ReqLocal(), peerCfg1.PeerConfig)
		if nil != err {
			goto ERR
		}
		resp := responses[0]
		localResp := resp.ForLocal()

		peers, err := localResp.Peers()
		if nil != err {
			goto ERR
		}

		result.Success(peers)
		return &result
	}

ERR:
	result.FailErr(err)
	return &result
}

// 未调通
func discoveryClientEndorsersPeers(channelID, orgName, orgUser, peerName, chainCodeID string, sdk *fabsdk.FabricSDK) *response.Result {
	result := response.Result{}
	var (
		contextClient pc.Client
		err           error
	)
	//prepare context
	ctx := sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))
	if contextClient, err = ctx(); nil != err {
		goto ERR
	}

	if disClient, err := discovery.New(contextClient); nil != err {
		goto ERR
	} else {
		reqCtx, cancel := context.NewRequest(contextClient, context.WithTimeout(10*time.Second))
		defer cancel()

		//req := discclient.NewRequest().OfChannel("mychannel").AddPeersQuery()

		peerCfg1, err := comm.NetworkPeerConfig(contextClient.EndpointConfig(), peerName)
		if nil != err {
			goto ERR
		}
		responses, err := disClient.Send(reqCtx, disClient.ReqEndorsers(channelID, newInterest(newCCCall(chainCodeID))), peerCfg1.PeerConfig)
		if nil != err {
			goto ERR
		}
		resp := responses[0]
		localResp := resp.ForLocal()

		peers, err := localResp.Peers()
		if nil != err {
			goto ERR
		}

		result.Success(peers)
		return &result
	}

ERR:
	result.FailErr(err)
	return &result
}

// 未调通
func discoveryClientConfigPeers(channelID, orgName, orgUser, peerName string, sdk *fabsdk.FabricSDK) *response.Result {
	result := response.Result{}
	var (
		contextClient pc.Client
		err           error
	)
	//prepare context
	ctx := sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))
	if contextClient, err = ctx(); nil != err {
		goto ERR
	}

	if disClient, err := discovery.New(contextClient); nil != err {
		goto ERR
	} else {
		reqCtx, cancel := context.NewRequest(contextClient, context.WithTimeout(10*time.Second))
		defer cancel()

		//req := discclient.NewRequest().OfChannel("mychannel").AddPeersQuery()

		peerCfg1, err := comm.NetworkPeerConfig(contextClient.EndpointConfig(), peerName)
		if nil != err {
			goto ERR
		}
		responses, err := disClient.Send(reqCtx, disClient.ReqConfig(channelID), peerCfg1.PeerConfig)
		if nil != err {
			goto ERR
		}
		resp := responses[0]
		localResp := resp.ForLocal()

		peers, err := localResp.Peers()
		if nil != err {
			goto ERR
		}

		result.Success(peers)
		return &result
	}

ERR:
	result.FailErr(err)
	return &result
}

func newInterest(ccCalls ...*fabdiscovery.ChaincodeCall) *fabdiscovery.ChaincodeInterest {
	return &fabdiscovery.ChaincodeInterest{Chaincodes: ccCalls}
}

func newCCCall(ccID string, collections ...string) *fabdiscovery.ChaincodeCall {
	return &fabdiscovery.ChaincodeCall{
		Name:            ccID,
		CollectionNames: collections,
	}
}
