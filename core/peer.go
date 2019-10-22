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
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
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
