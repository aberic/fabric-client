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
	"github.com/aberic/gnomon"
	config2 "github.com/ennoo/fabric-client/config"
	"github.com/ennoo/fabric-client/service"
	"github.com/ennoo/rivet"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	ctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"net/http"
	"strings"
)

// setupAndRun enables testing an end-to-end scenario against the supplied SDK options
// the createChannel flag will be used to either create a channel and the example CC or not(ie run the tests with existing ch and CC)
func Create(orderOrgName, orderOrgUser, orderURL, orgName, orgUser, channelID, channelConfigPath string, configBytes []byte,
	sdkOpts ...fabsdk.Option) (txID string, err error) {
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, sdk, err := resMgmtOrdererClient(orderOrgName, orderOrgUser, configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("Create", gnomon.Log().Err(err))
		return "", err
	}
	defer sdk.Close()
	return createChannel(orderURL, orgName, orgUser, channelID, channelConfigPath, sdk, resMgmtClient)
}

func Join(orderURL, orgName, orgUser, channelID, peerName string, configBytes []byte, sdkOpts ...fabsdk.Option) error {
	resMgmtClient, sdk, err := resMgmtOrgClient(orgName, orgUser, configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("Join", gnomon.Log().Err(err))
		return err
	}
	defer sdk.Close()
	return joinChannel(orderURL, channelID, peerName, resMgmtClient)
}

func Channels(orgName, orgUser, peerName string, configBytes []byte, sdkOpts ...fabsdk.Option) ([]*peer.ChannelInfo, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("Channels", gnomon.Log().Err(err))
		return nil, err
	}
	defer sdk.Close()
	return queryChannels(orgName, orgUser, peerName, sdk)
}

func QueryLedgerInfo(configID, peerName, channelID string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	if orgName, orgUser, err := get(configID, channelID); nil != err {
		gnomon.Log().Error("QueryLedgerInfo", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	} else {
		return QueryLedgerInfoSpec(peerName, channelID, orgName, orgUser, configBytes, sdkOpts...)
	}
}

func QueryLedgerBlockByHeight(configID, peerName, channelID string, height uint64, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	if orgName, orgUser, err := get(configID, channelID); nil != err {
		gnomon.Log().Error("QueryLedgerBlockByHeight", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	} else {
		return QueryLedgerBlockByHeightSpec(peerName, channelID, orgName, orgUser, height, configBytes, sdkOpts...)
	}
}

func QueryLedgerBlockByHash(configID, peerName, channelID, hash string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	if orgName, orgUser, err := get(configID, channelID); nil != err {
		gnomon.Log().Error("QueryLedgerBlockByHash", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	} else {
		return QueryLedgerBlockByHashSpec(peerName, channelID, orgName, orgUser, hash, configBytes, sdkOpts...)
	}
}

func QueryLedgerBlockByTxID(configID, peerName, channelID, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	if orgName, orgUser, err := get(configID, channelID); nil != err {
		gnomon.Log().Error("QueryLedgerBlockByTxID", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	} else {
		return QueryLedgerBlockByTxIDSpec(peerName, channelID, orgName, orgUser, txID, configBytes, sdkOpts...)
	}
}

func QueryLedgerTransaction(configID, peerName, channelID, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	if orgName, orgUser, err := get(configID, channelID); nil != err {
		gnomon.Log().Error("QueryLedgerTransaction", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	} else {
		return QueryLedgerTransactionSpec(peerName, channelID, orgName, orgUser, txID, configBytes, sdkOpts...)
	}
}

func QueryLedgerConfig(configID, peerName, channelID string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	if orgName, orgUser, err := get(configID, channelID); nil != err {
		gnomon.Log().Error("QueryLedgerConfig", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	} else {
		return QueryLedgerConfigSpec(peerName, channelID, orgName, orgUser, configBytes, sdkOpts...)
	}
}

func QueryLedgerInfoSpec(peerName, channelID, orgName, orgUser string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("QueryLedgerInfoSpec", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryLedgerInfo(peerName, channelProvider(orgName, orgUser, channelID, sdk))
}

func QueryLedgerBlockByHeightSpec(peerName, channelID, orgName, orgUser string, height uint64, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("QueryLedgerBlockByHeightSpec", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryLedgerBlockByHeight(peerName, height, channelProvider(orgName, orgUser, channelID, sdk))
}

func QueryLedgerBlockByHashSpec(peerName, channelID, orgName, orgUser, hash string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("QueryLedgerBlockByHashSpec", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryLedgerBlockByHash(peerName, hash, channelProvider(orgName, orgUser, channelID, sdk))
}

func QueryLedgerBlockByTxIDSpec(peerName, channelID, orgName, orgUser, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("QueryLedgerBlockByTxIDSpec", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryLedgerBlockByTxID(peerName, txID, channelProvider(orgName, orgUser, channelID, sdk))
}

func QueryLedgerTransactionSpec(peerName, channelID, orgName, orgUser, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("QueryLedgerTransactionSpec", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryLedgerTransaction(peerName, txID, channelProvider(orgName, orgUser, channelID, sdk))
}

func QueryLedgerConfigSpec(peerName, channelID, orgName, orgUser string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("QueryLedgerConfigSpec", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryLedgerConfig(peerName, channelProvider(orgName, orgUser, channelID, sdk))
}

func QueryChannelInfo(channelID, orgName, orgUser, peerName string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	var (
		client ctx.Client
		err    error
	)
	clientContext, sdk := clientContext(orgName, orgUser, configBytes, sdkOpts...)
	defer sdk.Close()
	if client, err = clientContext(); nil != err {
		gnomon.Log().Error("QueryChannelInfo", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	return queryChannelInfo(channelID, peerName, client)
}

func QueryChannelBlockByHeight(channelID, orgName, orgUser, peerName string, height uint64, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	var (
		client ctx.Client
		err    error
	)
	clientContext, sdk := clientContext(orgName, orgUser, configBytes, sdkOpts...)
	defer sdk.Close()
	if client, err = clientContext(); nil != err {
		gnomon.Log().Error("QueryChannelBlockByHeight", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	return queryChannelBlockByHeight(channelID, peerName, height, client)
}

func QueryChannelBlockByHash(channelID, orgName, orgUser, peerName, hash string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	var (
		client ctx.Client
		err    error
	)
	clientContext, sdk := clientContext(orgName, orgUser, configBytes, sdkOpts...)
	defer sdk.Close()
	if client, err = clientContext(); nil != err {
		gnomon.Log().Error("QueryChannelBlockByHash", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	return queryChannelBlockByHash(channelID, peerName, hash, client)
}

func QueryChannelBlockByTxID(channelID, orgName, orgUser, peerName, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	var (
		client ctx.Client
		err    error
	)
	clientContext, sdk := clientContext(orgName, orgUser, configBytes, sdkOpts...)
	defer sdk.Close()
	if client, err = clientContext(); nil != err {
		gnomon.Log().Error("QueryChannelBlockByTxID", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	return queryChannelBlockByTxID(channelID, peerName, txID, client)
}

func QueryChannelTransaction(channelID, orgName, orgUser, peerName, txID string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	var (
		client ctx.Client
		err    error
	)
	clientContext, sdk := clientContext(orgName, orgUser, configBytes, sdkOpts...)
	defer sdk.Close()
	if client, err = clientContext(); nil != err {
		gnomon.Log().Error("QueryChannelTransaction", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	return queryChannelTransaction(channelID, peerName, txID, client)
}

func Install(orgName, orgUser, peerName, name, goPath, chainCodePath, version string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, sdk, err := resMgmtOrgClient(orgName, orgUser, configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("Install", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return install(peerName, name, goPath, chainCodePath, version, resMgmtClient)
}

func OrderConfig(orgName, orgUser, channelID, orderURL string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("OrderConfig", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryConfigFromOrderer(orgName, orgUser, channelID, orderURL, sdk)
}

func Installed(orgName, orgUser, peerName string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("Installed", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryInstalled(orgName, orgUser, peerName, sdk)
}

func Instantiate(orgName, orgUser, peerName, channelID, name, path, version string, orgPolicies []string, args [][]byte,
	configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, _, err := resMgmtOrdererClient(orgName, orgUser, configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("Instantiate", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	return instantiate(peerName, channelID, name, path, version, orgPolicies, args, resMgmtClient)
}

func Instantiated(orgName, orgUser, channelID, peerName string, configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("Instantiated", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryInstantiate(orgName, orgUser, channelID, peerName, sdk)
}

func Upgrade(orgName, orgUser, peerName, channelID, name, path, version string, orgPolicies []string, args [][]byte,
	configBytes []byte, sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, _, err := resMgmtOrdererClient(orgName, orgUser, configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("Upgrade", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	return upgrade(peerName, channelID, name, path, version, orgPolicies, args, resMgmtClient)
}

func Invoke(chaincodeID, orgName, orgUser, channelID, fcn string, args [][]byte, targetEndpoints []string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("Invoke", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	channelClient := channelClient(orgName, orgUser, channelID, sdk)
	return invoke(chaincodeID, fcn, args, channelClient, targetEndpoints...)
}

func InvokeAsync(chaincodeID, orgName, orgUser, channelID, callback, fcn string, args [][]byte, targetEndpoints []string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	go func() {
		res := Invoke(chaincodeID, orgName, orgUser, channelID, fcn, args, targetEndpoints, configBytes, sdkOpts...)
		if gnomon.String().IsNotEmpty(callback) {
			gnomon.Log().Debug("InvokeAsync", gnomon.Log().Field("callback", callback))
			_, err := rivet.Request().RestJSONByURL(http.MethodPost, callback, res)
			if nil != err {
				gnomon.Log().Error("InvokeAsync", gnomon.Log().Err(err))
			}
		}
	}()
	result.Success("commit success")
	return &result
}

func Query(chaincodeID, orgName, orgUser, channelID, fcn string, args [][]byte, targetEndpoints []string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("Query", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	channelClient := channelClient(orgName, orgUser, channelID, sdk)
	return query(chaincodeID, fcn, args, channelClient, targetEndpoints...)
}

func QueryCollectionsConfig(chaincodeID, orgName, orgUser, channelID, peerName string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("QueryCollectionsConfig", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	defer sdk.Close()
	return queryCollectionsConfig(orgName, orgUser, peerName, channelID, chaincodeID, sdk)
}

func DiscoveryChannelPeers(channelID, orgName, orgUser string, configBytes []byte, sdkOpts ...fabsdk.Option) ([]fab.Peer, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("DiscoveryChannelPeers", gnomon.Log().Err(err))
		return nil, err
	}
	return discoveryChannelPeers(channelID, orgName, orgUser, sdk)
}

func DiscoveryLocalPeers(orgName, orgUser string, configBytes []byte, sdkOpts ...fabsdk.Option) ([]fab.Peer, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("DiscoveryLocalPeers", gnomon.Log().Err(err))
		return nil, err
	}
	return discoveryLocalPeers(orgName, orgUser, sdk)
}

func DiscoveryClientPeers(channelID, orgName, orgUser, peerName string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("DiscoveryClientPeers", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	return discoveryClientPeers(channelID, orgName, orgUser, peerName, sdk)
}

func DiscoveryClientLocalPeers(orgName, orgUser, peerName string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("DiscoveryClientLocalPeers", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	return discoveryClientLocalPeers(orgName, orgUser, peerName, sdk)
}

func DiscoveryClientConfigPeers(channelID, orgName, orgUser, peerName string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("DiscoveryClientConfigPeers", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	return discoveryClientConfigPeers(channelID, orgName, orgUser, peerName, sdk)
}

func DiscoveryClientEndorsersPeers(channelID, orgName, orgUser, peerName, chainCodeID string, configBytes []byte,
	sdkOpts ...fabsdk.Option) *Result {
	result := Result{}
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		gnomon.Log().Error("DiscoveryClientEndorsersPeers", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	return discoveryClientEndorsersPeers(channelID, orgName, orgUser, peerName, chainCodeID, sdk)
}

func CAInfo(orgName string, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.GetCAInfoResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return caInfo(orgName, sdk)
}

func Enroll(orgName, enrollmentID string, configBytes []byte, opts []msp.EnrollmentOption, sdkOpts ...fabsdk.Option) error {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return err
	}
	defer sdk.Close()
	return enroll(orgName, enrollmentID, sdk, opts...)
}

func Reenroll(orgName, enrollmentID string, configBytes []byte, opts []msp.EnrollmentOption, sdkOpts ...fabsdk.Option) error {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return err
	}
	defer sdk.Close()
	return reenroll(orgName, enrollmentID, sdk, opts...)
}

func Register(orgName string, registerReq *msp.RegistrationRequest, configBytes []byte, sdkOpts ...fabsdk.Option) (string, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return "", err
	}
	defer sdk.Close()
	return register(orgName, registerReq, sdk)
}

func AddAffiliation(orgName string, affReq *msp.AffiliationRequest, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.AffiliationResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return addAffiliation(orgName, affReq, sdk)
}

func RemoveAffiliation(orgName string, affReq *msp.AffiliationRequest, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.AffiliationResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return removeAffiliation(orgName, affReq, sdk)
}

func ModifyAffiliation(orgName string, affReq *msp.ModifyAffiliationRequest, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.AffiliationResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return modifyAffiliation(orgName, affReq, sdk)
}

func GetAllAffiliations(orgName string, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.AffiliationResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return getAllAffiliations(orgName, sdk)
}

func GetAllAffiliationsByCaName(orgName, caName string, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.AffiliationResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return getAllAffiliationsByCaName(orgName, caName, sdk)
}

func GetAffiliation(affiliation, orgName string, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.AffiliationResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return getAffiliation(affiliation, orgName, sdk)
}

func GetAffiliationByCaName(affiliation, orgName, caName string, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.AffiliationResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return getAffiliationByCaName(affiliation, orgName, caName, sdk)
}

func GetAllIdentities(orgName string, configBytes []byte, sdkOpts ...fabsdk.Option) ([]*msp.IdentityResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return getAllIdentities(orgName, sdk)
}

func GetAllIdentitiesByCaName(orgName, caName string, configBytes []byte, sdkOpts ...fabsdk.Option) ([]*msp.IdentityResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return getAllIdentitiesByCaName(orgName, caName, sdk)
}

func CreateIdentity(orgName string, req *msp.IdentityRequest, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.IdentityResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return createIdentity(orgName, req, sdk)
}

func ModifyIdentity(orgName string, req *msp.IdentityRequest, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.IdentityResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return modifyIdentity(orgName, req, sdk)
}

func GetIdentity(id, orgName string, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.IdentityResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return getIdentity(id, orgName, sdk)
}

func GetIdentityByCaName(id, caName, orgName string, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.IdentityResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return getIdentityByCaName(id, caName, orgName, sdk)
}

func RemoveIdentity(orgName string, req *msp.RemoveIdentityRequest, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.IdentityResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return removeIdentity(orgName, req, sdk)
}

func CreateSigningIdentity(orgName string, configBytes []byte, opts []mspctx.SigningIdentityOption, sdkOpts ...fabsdk.Option) (mspctx.SigningIdentity, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return createSigningIdentity(orgName, sdk, opts...)
}

func GetSigningIdentity(id, orgName string, configBytes []byte, sdkOpts ...fabsdk.Option) (mspctx.SigningIdentity, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return getSigningIdentity(id, orgName, sdk)
}

func Revoke(orgName string, req *msp.RevocationRequest, configBytes []byte, sdkOpts ...fabsdk.Option) (*msp.RevocationResponse, error) {
	sdk, err := sdk(configBytes, sdkOpts...)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()
	return revoke(orgName, req, sdk)
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
		gnomon.Log().Error("resMgmtOrdererClient", gnomon.Log().Err(err))
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
		gnomon.Log().Error("resMgmtOrgClient", gnomon.Log().Err(err))
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

	for name, chPeer := range channel.Peers {
		if !chPeer.LedgerQuery {
			continue
		}
		peerName = name
		break
	}
	if gnomon.String().IsEmpty(peerName) {
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

	if gnomon.String().IsEmpty(orgName) {
		err = errors.New("org is nil")
		return
	}
	return
}
