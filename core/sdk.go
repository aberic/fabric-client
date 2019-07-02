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
	"errors"
	"fmt"
	config2 "github.com/ennoo/fabric-client/config"
	"github.com/ennoo/fabric-client/service"
	"github.com/ennoo/rivet/trans/response"
	"github.com/ennoo/rivet/utils/log"
	"github.com/ennoo/rivet/utils/string"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	ctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"strings"
)

// setupAndRun enables testing an end-to-end scenario against the supplied SDK options
// the createChannel flag will be used to either create a channel and the example CC or not(ie run the tests with existing ch and CC)
func Create(orderOrgName, orderOrgUser, orderURL, orgName, orgUser, channelID, channelConfigPath string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, sdk, err := resMgmtOrdererClient(orderOrgName, orderOrgUser, configBytes, sdkOpts...)
	if err != nil {
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return createChannel(orderURL, orgName, orgUser, channelID, channelConfigPath, sdk, resMgmtClient)
}

func Join(orderURL, orgName, orgUser, channelID string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	resMgmtClient, sdk, err := resMgmtOrgClient(orgName, orgUser, configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return joinChannel(orderURL, channelID, resMgmtClient)
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

func QueryLedgerInfo(configID, channelID string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	if orgName, orgUser, err := get(configID, channelID); nil != err {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	} else {
		return QueryLedgerInfoSpec(channelID, orgName, orgUser, configBytes, sdkOpts...)
	}
}

func QueryLedgerBlockByHeight(configID, channelID string, height uint64, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	if orgName, orgUser, err := get(configID, channelID); nil != err {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	} else {
		return QueryLedgerBlockByHeightSpec(channelID, orgName, orgUser, height, configBytes, sdkOpts...)
	}
}

func QueryLedgerBlockByHash(configID, channelID, hash string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	if orgName, orgUser, err := get(configID, channelID); nil != err {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	} else {
		return QueryLedgerBlockByHashSpec(channelID, orgName, orgUser, hash, configBytes, sdkOpts...)
	}
}

func QueryLedgerBlockByTxID(configID, channelID, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	if orgName, orgUser, err := get(configID, channelID); nil != err {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	} else {
		return QueryLedgerBlockByTxIDSpec(channelID, orgName, orgUser, txID, configBytes, sdkOpts...)
	}
}

func QueryLedgerTransaction(configID, channelID, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	if orgName, orgUser, err := get(configID, channelID); nil != err {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	} else {
		return QueryLedgerTransactionSpec(channelID, orgName, orgUser, txID, configBytes, sdkOpts...)
	}
}

func QueryLedgerConfig(configID, channelID string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	if orgName, orgUser, err := get(configID, channelID); nil != err {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	} else {
		return QueryLedgerConfigSpec(channelID, orgName, orgUser, configBytes, sdkOpts...)
	}
}

func QueryLedgerInfoSpec(channelID, orgName, orgUser string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryLedgerInfo(channelProvider(orgName, orgUser, channelID, sdk))
}

func QueryLedgerBlockByHeightSpec(channelID, orgName, orgUser string, height uint64, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryLedgerBlockByHeight(height, channelProvider(orgName, orgUser, channelID, sdk))
}

func QueryLedgerBlockByHashSpec(channelID, orgName, orgUser, hash string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryLedgerBlockByHash(hash, channelProvider(orgName, orgUser, channelID, sdk))
}

func QueryLedgerBlockByTxIDSpec(channelID, orgName, orgUser, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryLedgerBlockByTxID(txID, channelProvider(orgName, orgUser, channelID, sdk))
}

func QueryLedgerTransactionSpec(channelID, orgName, orgUser, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryLedgerTransaction(txID, channelProvider(orgName, orgUser, channelID, sdk))
}

func QueryLedgerConfigSpec(channelID, orgName, orgUser string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryLedgerConfig(channelProvider(orgName, orgUser, channelID, sdk))
}

func QueryChannelInfo(channelID, orgName, orgUser, peerName string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
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
	return queryChannelInfo(channelID, peerName, client)
}

func QueryChannelBlockByHeight(channelID, orgName, orgUser, peerName string, height uint64, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
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
	return queryChannelBlockByHeight(channelID, peerName, height, client)
}

func QueryChannelBlockByHash(channelID, orgName, orgUser, peerName, hash string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
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
	return queryChannelBlockByHash(channelID, peerName, hash, client)
}

func QueryChannelBlockByTxID(channelID, orgName, orgUser, peerName, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
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
	return queryChannelBlockByTxID(channelID, peerName, txID, client)
}

func QueryChannelTransaction(channelID, orgName, orgUser, peerName, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
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
	return queryChannelTransaction(channelID, peerName, txID, client)
}

func Install(orgName, orgUser, name, goPath, chainCodePath, version string, configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, sdk, err := resMgmtOrgClient(orgName, orgUser, configBytes, sdkOpts...)
	if err != nil {
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return install(name, goPath, chainCodePath, version, resMgmtClient)
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

func Instantiate(orgName, orgUser, channelID, name, path, version string, orgPolicies []string, args [][]byte,
	configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, _, err := resMgmtOrdererClient(orgName, orgUser, configBytes, sdkOpts...)
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

func Upgrade(orgName, orgUser, channelID, name, path, version string, orgPolicies []string, args [][]byte,
	configBytes []byte, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, _, err := resMgmtOrdererClient(orgName, orgUser, configBytes, sdkOpts...)
	if err != nil {
		result.Fail(err.Error())
		return &result
	}
	return upgrade(channelID, name, path, version, orgPolicies, args, resMgmtClient)
}

func Invoke(chaincodeID, orgName, orgUser, channelID, fcn string, args [][]byte, targetEndpoints []string, configBytes []byte,
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
	return invoke(chaincodeID, fcn, args, channelClient, targetEndpoints...)
}

func InvokeAsync(chaincodeID, orgName, orgUser, channelID, fcn string, args [][]byte, targetEndpoints []string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	go Invoke(chaincodeID, orgName, orgUser, channelID, fcn, args, targetEndpoints, configBytes, sdkOpts...)
	result.Success("commit success")
	return &result
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

func resMgmtOrdererClient(ordererOrgName, ordererUser string, configBytes []byte, sdkOpts ...fabsdk.Option) (*resmgmt.Client,
	*fabsdk.FabricSDK, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, nil, err
	}
	//clientContext allows creation of transactions using the supplied identity as the credential.
	clientContext := sdk.Context(fabsdk.WithUser(ordererUser), fabsdk.WithOrg(ordererOrgName))

	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Self.Error("Failed to create channel management client: " + err.Error())
		return nil, nil, fmt.Errorf("Failed to create channel management client: " + err.Error())
	}
	return resMgmtClient, sdk, nil
}

func resMgmtOrgClient(orgName, orgUser string, configBytes []byte, sdkOpts ...fabsdk.Option) (*resmgmt.Client,
	*fabsdk.FabricSDK, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, nil, err
	}
	//clientContext allows creation of transactions using the supplied identity as the credential.
	clientContext := sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))

	// Org resource management client
	orgResMgmt, err := resmgmt.New(clientContext)
	if err != nil {
		log.Self.Error("Failed to create new resource management client: " + err.Error())
		return nil, nil, fmt.Errorf("Failed to create new resource management client: " + err.Error())
	}
	return orgResMgmt, sdk, nil
}

func get(configID, channelID string) (orgName, orgUser string, err error) {
	var (
		peerName string
		conf     *config2.Config
		channel  *config2.Channel
	)
	if conf = service.Configs[configID]; nil == conf {
		err = errors.New("config must init first")
		return
	}
	if channel = conf.Channels[channelID]; nil == channel {
		err = errors.New("channel must init first")
		return
	}

	for name, peer := range channel.Peers {
		if !peer.LedgerQuery {
			continue
		}
		peerName = name
		break
	}
	if str.IsEmpty(peerName) {
		err = errors.New("peer is nil")
		return
	}

	for oName, orgItf := range conf.Organizations {
		if oName == config2.OrderOrgKey {
			continue
		}
		have := false
		for _, pName := range orgItf.Peers {
			if peerName == pName {
				orgName = oName
				have = true
				break
			}
		}
		if have {
			str1 := strings.Split(orgItf.CryptoPath, "@")
			str2 := strings.Split(str1[0], "/")
			orgUser = str2[len(str2)-1]
			break
		}
	}

	if str.IsEmpty(orgName) {
		err = errors.New("org is nil")
		return
	}
	return
}
