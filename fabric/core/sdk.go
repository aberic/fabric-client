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
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// setupAndRun enables testing an end-to-end scenario against the supplied SDK options
// the createChannel flag will be used to either create a channel and the example CC or not(ie run the tests with existing ch and CC)
func Create(ordererOrgName, orgName, orgAdmin, channelID, channelConfigPath string, configYmlPath string,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, sdk, err := resMgmtClient(ordererOrgName, orgAdmin, configYmlPath, sdkOpts...)
	if err != nil {
		result.Fail(err.Error())
		return &result
	}
	return createChannel(orgName, orgAdmin, channelID, channelConfigPath, sdk, resMgmtClient)
}

func Join(orgName, orgAdmin, channelID, peerUrl string, configYmlPath string, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configYmlPath, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return joinChannel(orgName, orgAdmin, channelID, peerUrl, sdk)
}

func Install(ordererOrgName, orgAdmin, name, source, path, version, configYmlPath string,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, _, err := resMgmtClient(ordererOrgName, orgAdmin, configYmlPath, sdkOpts...)
	if err != nil {
		result.Fail(err.Error())
		return &result
	}
	return install(name, source, path, version, resMgmtClient)
}

func Instantiate(ordererOrgName, orgAdmin, channelID, name, path, version string, orgPolicies []string, args [][]byte,
	configYmlPath string, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, _, err := resMgmtClient(ordererOrgName, orgAdmin, configYmlPath, sdkOpts...)
	if err != nil {
		result.Fail(err.Error())
		return &result
	}
	return instantiate(channelID, name, path, version, orgPolicies, args, resMgmtClient)
}

func Upgrade(ordererOrgName, orgAdmin, channelID, name, path, version string, orgPolicies []string, args [][]byte,
	configYmlPath string, sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, _, err := resMgmtClient(ordererOrgName, orgAdmin, configYmlPath, sdkOpts...)
	if err != nil {
		result.Fail(err.Error())
		return &result
	}
	return upgrade(channelID, name, path, version, orgPolicies, args, resMgmtClient)
}

func Invoke(chaincodeID, orgName, orgUser, channelID, fcn string, args [][]byte, configYmlPath string,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configYmlPath, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	channelClient := channelClient(orgName, orgUser, channelID, sdk)
	return invoke(chaincodeID, fcn, args, channelClient)
}

func Query(chaincodeID, orgName, orgUser, channelID, fcn string, args [][]byte, targetEndpoints []string, configYmlPath string,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdk(configYmlPath, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	channelClient := channelClient(orgName, orgUser, channelID, sdk)
	return query(chaincodeID, fcn, args, channelClient, targetEndpoints...)
}

func QueryRaw(chaincodeID, orgName, orgUser, channelID, fcn string, args [][]byte, targetEndpoints []string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *response.Result {
	result := response.Result{}
	sdk, err := sdkRaw(configBytes, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	channelClient := channelClient(orgName, orgUser, channelID, sdk)
	return query(chaincodeID, fcn, args, channelClient, targetEndpoints...)
}

func sdk(configYmlPath string, sdkOpts ...fabsdk.Option) (*fabsdk.FabricSDK, error) {
	configOpt := config.FromFile(configYmlPath)
	sdk, err := fabsdk.New(configOpt, sdkOpts...)
	if err != nil {
		return nil, err
	}
	if nil == sdk {
		return nil, fmt.Errorf("sdk error should be nil")
	}
	return sdk, nil
}

func sdkRaw(configBytes []byte, sdkOpts ...fabsdk.Option) (*fabsdk.FabricSDK, error) {
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

func resMgmtClient(ordererOrgName, orgAdmin, configYmlPath string,
	sdkOpts ...fabsdk.Option) (*resmgmt.Client, *fabsdk.FabricSDK, error) {
	sdk, err := sdk(configYmlPath, sdkOpts...)
	if err != nil {
		return nil, nil, err
	}
	defer sdk.Close()
	//clientContext allows creation of transactions using the supplied identity as the credential.
	clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(ordererOrgName))

	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Self.Error("Failed to create channel management client: " + err.Error())
		return nil, nil, fmt.Errorf("Failed to create channel management client: " + err.Error())
	}
	return resMgmtClient, sdk, nil
}
