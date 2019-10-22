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
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	pc "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/comm"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/discovery"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	fabdiscovery "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/discovery"
	"github.com/pkg/errors"
	"time"
)

func discoveryChannelPeers(channelID, orgName, orgUser string, sdk *fabsdk.FabricSDK) ([]fab.Peer, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, errors.Errorf("error creating MSP client: %s", err)
	}
	user, err := mspClient.GetSigningIdentity(orgUser)
	if err != nil {
		return nil, errors.Errorf("GetSigningIdentity returned error: %v", err)
	}
	clientContext := sdk.Context(fabsdk.WithIdentity(user))
	channelProvider := func() (*context.Channel, error) {
		return context.NewChannel(clientContext, channelID)
	}
	chContext, err := channelProvider()
	if err != nil {
		return nil, err
	}
	discoveryDo, err := chContext.ChannelService().Discovery()
	if err != nil {
		return nil, err
	}
	peers, err := discoveryDo.GetPeers()
	if err != nil {
		return nil, err
	}
	return peers, nil
}

func discoveryLocalPeers(orgName, orgUser string, sdk *fabsdk.FabricSDK) ([]fab.Peer, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, errors.Errorf("error creating MSP client: %s", err)
	}
	user, err := mspClient.GetSigningIdentity(orgUser)
	if err != nil {
		return nil, errors.Errorf("GetSigningIdentity returned error: %v", err)
	}
	clientContext := sdk.Context(fabsdk.WithIdentity(user))
	localContext, err := context.NewLocal(clientContext)
	if err != nil {
		return nil, err
	}
	peers, err := localContext.LocalDiscoveryService().GetPeers()
	if err != nil {
		return nil, err
	}
	return peers, nil
}

func discoveryClientPeers(channelID, orgName, orgUser, peerName string, sdk *fabsdk.FabricSDK) *Result {
	result := Result{}
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

func discoveryClientLocalPeers(orgName, orgUser, peerName string, sdk *fabsdk.FabricSDK) *Result {
	result := Result{}
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
func discoveryClientEndorsersPeers(channelID, orgName, orgUser, peerName, chainCodeID string, sdk *fabsdk.FabricSDK) *Result {
	result := Result{}
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
func discoveryClientConfigPeers(channelID, orgName, orgUser, peerName string, sdk *fabsdk.FabricSDK) *Result {
	result := Result{}
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
