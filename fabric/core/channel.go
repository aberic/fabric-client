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
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	fcmocks "github.com/hyperledger/fabric-sdk-go/pkg/fab/mocks"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"strings"
)

func channelClient(orgName, orgUser, channelID string, sdk *fabsdk.FabricSDK) *channel.Client {
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		log.Self.Error("Failed to create new channel client:" + err.Error())
		return nil
	}
	return client
}

// channelConfigPath mychannel.tx
func createChannel(orgName, orgAdmin, channelID, channelConfigPath string, sdk *fabsdk.FabricSDK, client *resmgmt.Client) *response.Result {
	result := response.Result{}
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	adminIdentity, err := mspClient.GetSigningIdentity(orgAdmin)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: channelConfigPath,
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := client.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(orgName))
	if err != nil {
		log.Self.Error("error should be nil. " + err.Error())
		result.Fail("error should be nil. " + err.Error())
	} else {
		result.Success(txID.TransactionID)
	}
	return &result
}

// ordererUrl "orderer.example.com"
// peerUrl grpc://peerUrl or grpcs://peerUrl
func joinChannel(orgName, orgAdmin, channelID, peerUrl string, sdk *fabsdk.FabricSDK) *response.Result {
	result := response.Result{}
	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		log.Self.Error("Failed to join channel: " + err.Error())
		result.Fail("Failed to join channel: " + err.Error())
	} else {
		if nil != orgResMgmt {
			peerNew, _ := peer.New(fcmocks.NewMockEndpointConfig(), peer.WithURL(peerUrl))
			// Org peers join channel
			if err = orgResMgmt.JoinChannel(
				channelID,
				resmgmt.WithRetry(retry.DefaultResMgmtOpts),
				resmgmt.WithOrdererEndpoint(orgName),
				resmgmt.WithTargets(peerNew)); err != nil {
				log.Self.Error("Org peers failed to JoinChannel: " + err.Error())
				result.Fail("Org peers failed to JoinChannel: " + err.Error())
			} else {
				result.Success("success")
			}
		} else {
			log.Self.Error("orgResMgmt error should be nil. ")
			result.Fail("orgResMgmt error should be nil. ")
		}
	}
	return &result
}

// peer 参见peer.go PeerChannel
func queryChannels(orgName, orgAdmin, channelID string, peer *fcmocks.MockPeer, sdk *fabsdk.FabricSDK) *response.Result {
	result := response.Result{}
	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		log.Self.Error("Failed to query channels: " + err.Error())
		result.Fail("Failed to query channels: " + err.Error())
	} else {
		if nil != orgResMgmt {
			found := false
			qcResponse, err := orgResMgmt.QueryChannels(resmgmt.WithTargets(peer))
			if err == nil {
				log.Self.Error("Failed to query channels: peer cannot be nil")
				result.Fail("Failed to query channels: peer cannot be nil")
			}
			if nil == qcResponse {
				log.Self.Error("qcResponse error should be nil. ")
				result.Fail("qcResponse error should be nil. ")
			} else {
				for _, responseChannel := range qcResponse.Channels {
					if responseChannel.ChannelId == channelID {
						found = true
						break
					}
				}
				if !found {
					result.Fail(strings.Join([]string{"Peer has not joined ", channelID, " channel"}, ""))
				} else {
					result.Success(strings.Join([]string{"Peer has joined ", channelID, " channel"}, ""))
				}
			}
		} else {
			log.Self.Error("orgResMgmt error should be nil. ")
			result.Fail("orgResMgmt error should be nil. ")
		}
	}
	return &result
}
