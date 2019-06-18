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
	"fmt"
	"github.com/ennoo/rivet/trans/response"
	"github.com/ennoo/rivet/utils/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	ctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// setupAndRun enables testing an end-to-end scenario against the supplied SDK options
// the createChannel flag will be used to either create a channel and the example CC or not(ie run the tests with existing ch and CC)
func Create(orderOrgName, orgName, orgUser, channelID, channelConfigPath string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, sdk, err := resMgmtClient(orderOrgName, orgUser, configBytes, sdkOpts...)
	if err != nil {
		result.Fail(err.Error())
		return &result
	}
	return createChannel(orgName, orgUser, channelID, channelConfigPath, sdk, resMgmtClient)
}

func Join(orderName, orgName, orgUser, channelID, peerUrl string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return joinChannel(orderName, orgName, orgUser, channelID, peerUrl, sdk)
}

func Channels(orgName, orgUser, peerName string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryChannels(orgName, orgUser, peerName, sdk)
}

func QueryInfo(channelID, orgName, orgUser, peerName string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	var (
		client ctx.Client
		err    error
	)
	clientContext, sdk := clientContext(orgName, orgUser, configBytes, sdkOpts...)
	defer sdk.Close()
	if client, err = clientContext(); nil != err {
		result.Fail(err.Error())
		return &result
	}
	return queryInfo(channelID, peerName, client)
}

func QueryBlockByHeight(channelID, orgName, orgUser, peerName string, height uint64, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	var (
		client ctx.Client
		err    error
	)
	clientContext, sdk := clientContext(orgName, orgUser, configBytes, sdkOpts...)
	defer sdk.Close()
	if client, err = clientContext(); nil != err {
		result.Fail(err.Error())
		return &result
	}
	return queryBlockByHeight(channelID, peerName, height, client)
}

func QueryBlockByHash(channelID, orgName, orgUser, peerName, hash string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	var (
		client ctx.Client
		err    error
	)
	clientContext, sdk := clientContext(orgName, orgUser, configBytes, sdkOpts...)
	defer sdk.Close()
	if client, err = clientContext(); nil != err {
		result.Fail(err.Error())
		return &result
	}
	return queryBlockByHash(channelID, peerName, hash, client)
}

func QueryBlockByTxID(channelID, orgName, orgUser, peerName, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	var (
		client ctx.Client
		err    error
	)
	clientContext, sdk := clientContext(orgName, orgUser, configBytes, sdkOpts...)
	defer sdk.Close()
	if client, err = clientContext(); nil != err {
		result.Fail(err.Error())
		return &result
	}
	return queryBlockByTxID(channelID, peerName, txID, client)
}

func QueryTransaction(channelID, orgName, orgUser, peerName, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	var (
		client ctx.Client
		err    error
	)
	clientContext, sdk := clientContext(orgName, orgUser, configBytes, sdkOpts...)
	defer sdk.Close()
	if client, err = clientContext(); nil != err {
		result.Fail(err.Error())
		return &result
	}
	return queryTransaction(channelID, peerName, txID, client)
}

func Install(orderOrgName, orgUser, name, source, path, version string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, _, err := resMgmtClient(orderOrgName, orgUser, configBytes, sdkOpts...)
	if err != nil {
		result.Fail(err.Error())
		return &result
	}
	return install(name, source, path, version, resMgmtClient)
}

func OrderConfig(orgName, orgUser, channelID, orderURL string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryConfigFromOrderer(orgName, orgUser, channelID, orderURL, sdk)
}

func Installed(orgName, orgUser, peerName string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryInstalled(orgName, orgUser, peerName, sdk)
}

func Instantiate(orderOrgName, orgUser, channelID, name, path, version string, orgPolicies []string, args [][]byte,
	configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, _, err := resMgmtClient(orderOrgName, orgUser, configBytes, sdkOpts...)
	if err != nil {
		result.Fail(err.Error())
		return &result
	}
	return instantiate(channelID, name, path, version, orgPolicies, args, resMgmtClient)
}

func Instantiated(orgName, orgUser, channelID, peerName string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryInstantiate(orgName, orgUser, channelID, peerName, sdk)
}

func Upgrade(ordererOrgName, orgUser, channelID, name, path, version string, orgPolicies []string, args [][]byte,
	configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, _, err := resMgmtClient(ordererOrgName, orgUser, configBytes, sdkOpts...)
	if err != nil {
		result.Fail(err.Error())
		return &result
	}
	return upgrade(channelID, name, path, version, orgPolicies, args, resMgmtClient)
}

func Invoke(chaincodeID, orgName, orgUser, channelID, fcn string, args [][]byte, configBytes []byte,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	channelClient := channelClient(orgName, orgUser, channelID, sdk)
	return invoke(chaincodeID, fcn, args, channelClient)
}

func Query(chaincodeID, orgName, orgUser, channelID, fcn string, args [][]byte, targetEndpoints []string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	channelClient := channelClient(orgName, orgUser, channelID, sdk)
	return query(chaincodeID, fcn, args, channelClient, targetEndpoints...)
}

func QueryCollectionsConfig(chaincodeID, orgName, orgUser, channelID, peerName string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryCollectionsConfig(orgName, orgUser, peerName, channelID, chaincodeID, sdk)
}

func DiscoveryClientPeers(channelID, orgName, orgUser, peerName string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	return discoveryClientPeers(channelID, orgName, orgUser, peerName, sdk)
}

func DiscoveryClientLocalPeers(orgName, orgUser, peerName string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	return discoveryClientLocalPeers(orgName, orgUser, peerName, sdk)
}

func DiscoveryClientConfigPeers(channelID, orgName, orgUser, peerName string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	return discoveryClientConfigPeers(channelID, orgName, orgUser, peerName, sdk)
}

func DiscoveryClientEndorsersPeers(channelID, orgName, orgUser, peerName, chainCodeID string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	return discoveryClientEndorsersPeers(channelID, orgName, orgUser, peerName, chainCodeID, sdk)
}

func sdk(configBytes []byte, sdkOpts ...fabsdk.Option) (*fabsdk.FabricSDK, error) {
	configOpt := config.FromRaw(configBytes, "yaml")
	sdk, err := fabsdk.New(configOpt, sdkOpts...)
	if err != nil {
		return nil, err
	}
	if nil == sdk {
		return nil, fmt.Errorf("sdk error should be nil")
	}
	return sdk, nil
}

func clientContext(orgName, orgUser string, configBytes []byte, sdkOpts ...fabsdk.Option) (context.ClientProvider,
	*fabsdk.FabricSDK) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, nil
	}
	//clientContext allows creation of transactions using the supplied identity as the credential.
	return sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName)), sdk
}

func resMgmtClient(ordererOrgName, orgUser string, configBytes []byte, sdkOpts ...fabsdk.Option) (*resmgmt.Client,
	*fabsdk.FabricSDK, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, nil, err
	}
	defer sdk.Close()
	//clientContext allows creation of transactions using the supplied identity as the credential.
	clientContext := sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(ordererOrgName))

	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Self.Error("Failed to create channel management client: " + err.Error())
		return nil, nil, fmt.Errorf("Failed to create channel management client: " + err.Error())
	}
	return resMgmtClient, sdk, nil
}
