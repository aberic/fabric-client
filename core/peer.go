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
	"github.com/hyperledger/fabric-sdk-go/pkg/client/common/discovery/dynamicdiscovery"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	contextApi "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/comm"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/discovery"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/mocks"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	discclient "github.com/hyperledger/fabric/discovery/client"
	"time"
)

func discoveryService(sdk *fabsdk.FabricSDK) *response.Result {
	result := response.Result{}
	var (
		contextClient context.Client
		err           error
	)
	//prepare context
	ctx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))
	if contextClient, err = getClient(ctx); nil != err {
		goto ERR
	}

	if disClient, err := discovery.New(contextClient); nil != err {
		goto ERR
	} else {
		reqCtx, cancel := context.NewRequest(contextClient, context.WithTimeout(10*time.Second))
		defer cancel()

		req := discclient.NewRequest().OfChannel("mychannel").AddPeersQuery()

		peerCfg1, err := comm.NetworkPeerConfig(contextClient.EndpointConfig(), "peer0.league01-org1-vh-cn")
		if nil != err {
			goto ERR
		}
		responses, err := disClient.Send(reqCtx, req, peerCfg1.PeerConfig)
		if nil != err {
			goto ERR
		}
		resp := responses[0]
		chanResp := resp.ForChannel("mychannel")

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

func getClient(ctx contextApi.ClientProvider) (context.Client, error) {
	return ctx()
}

func discoveryServiceOld(sdk *fabsdk.FabricSDK) *response.Result {
	result := response.Result{}
	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))
	//ctx := mocks.NewMockContext(mspmocks.NewMockSigningIdentity("test", mspID1))
	mspClient, err := msp.New(adminContext)
	signingIdentity, err := mspClient.GetSigningIdentity("Admin")
	ctx := mocks.NewMockContext(signingIdentity)

	var service *dynamicdiscovery.ChannelService
	service, err = dynamicdiscovery.NewChannelService(
		ctx, mocks.NewMockMembership(), "mychannel",
		dynamicdiscovery.WithRefreshInterval(10*time.Millisecond),
		dynamicdiscovery.WithResponseTimeout(100*time.Millisecond),
		dynamicdiscovery.WithErrorHandler(
			func(ctxt fab.ClientContext, channelID string, err error) {
				derr, ok := err.(dynamicdiscovery.DiscoveryError)
				if ok && derr.Error() == dynamicdiscovery.AccessDenied {
					service.Close()
				}
			},
		),
	)
	if nil != err {
		result.FailErr(err)
	} else {
		defer service.Close()
		if peers, err := service.GetPeers(); err == nil {
			result.Success(peers)
		} else {
			result.FailErr(err)
		}
	}
	return &result
}
