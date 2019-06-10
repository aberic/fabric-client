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
	"time"
)

func discoveryService(channelID, orgName, orgUser, peerName string, sdk *fabsdk.FabricSDK) *response.Result {
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
