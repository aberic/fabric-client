/*
 * Copyright (c) 2019. Aberic - All Rights Reserved.
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
	"github.com/aberic/gnomon"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"net/http"
)

// install 安装智能合约
func install(peerName, name, goPath, chainCodePath, version string, client *resmgmt.Client) *Result {
	result := Result{}
	ccPkg, err := gopackager.NewCCPackage(chainCodePath, goPath)
	if err != nil {
		gnomon.Log().Error("install", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: name, Path: chainCodePath, Version: version, Package: ccPkg}
	respList, err := client.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(peerName))
	if err != nil {
		gnomon.Log().Error("install", gnomon.Log().Err(err))
		result.Fail(err.Error())
		return &result
	}
	for index := range respList {
		resp := respList[index]
		if resp.Status == http.StatusOK {
			result.Success(resp.Info)
		} else {
			result.Fail(resp.Info)
		}
	}
	return &result
}

type ChainCodeInfoArr struct {
	ChainCodes []*peer.ChaincodeInfo `json:"chaincodes"`
}

// peer 参见peer.go Peer
func queryInstalled(orgName, orgUser, peerName string, sdk *fabsdk.FabricSDK) *Result {
	result := Result{}
	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		gnomon.Log().Error("queryInstalled", gnomon.Log().Err(err))
		result.Fail("Failed to query installed: " + err.Error())
	} else {
		if nil != orgResMgmt {
			qiResponse, err := orgResMgmt.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(peerName))
			if err != nil {
				gnomon.Log().Error("queryInstalled", gnomon.Log().Err(err))
				result.Fail("Failed to query installed: " + err.Error())
			} else {
				result.Success(&ChainCodeInfoArr{qiResponse.Chaincodes})
			}
		} else {
			gnomon.Log().Error("queryInstalled", gnomon.Log().Err(err))
			result.Fail("orgResMgmt error should be nil. ")
		}
	}
	return &result
}

// args [][]byte{[]byte(coll1), []byte("key"), []byte("value")}
func instantiate(peerName, channelID, name, path, version string, orgPolicies []string, args [][]byte, client *resmgmt.Client) *Result {
	result := Result{}
	ccPolicy := cauthdsl.SignedByAnyMember(orgPolicies)
	// Org resource manager will instantiate 'example_cc' on channel
	resp, err := client.InstantiateCC(
		channelID,
		resmgmt.InstantiateCCRequest{Name: name, Path: path, Version: version, Args: args, Policy: ccPolicy},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(peerName),
	)
	if err != nil {
		gnomon.Log().Error("instantiate", gnomon.Log().Err(err))
		result.Fail(err.Error())
	} else {
		result.Success(string(resp.TransactionID))
	}
	return &result
}

// peer 参见peer.go Peer
func queryInstantiate(orgName, orgUser, channelID, peerName string, sdk *fabsdk.FabricSDK) *Result {
	result := Result{}
	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		gnomon.Log().Error("queryInstantiate", gnomon.Log().Err(err))
		result.Fail("Failed to query instantiate: " + err.Error())
	} else {
		if nil != orgResMgmt {
			qiResponse, err := orgResMgmt.QueryInstantiatedChaincodes(channelID, resmgmt.WithTargetEndpoints(peerName))
			if err != nil {
				gnomon.Log().Error("queryInstantiate", gnomon.Log().Err(err))
				result.Fail("Failed to query instantiate: " + err.Error())
			} else {
				result.Success(&ChainCodeInfoArr{qiResponse.Chaincodes})
			}
		} else {
			gnomon.Log().Error("queryInstantiate", gnomon.Log().Err(err))
			result.Fail("orgResMgmt error should be nil. ")
		}
	}
	return &result
}

// args [][]byte{[]byte(coll1), []byte("key"), []byte("value")}
func upgrade(peerName, channelID, name, path, version string, orgPolicies []string, args [][]byte, client *resmgmt.Client) *Result {
	result := Result{}
	ccPolicy := cauthdsl.SignedByAnyMember(orgPolicies)
	// Org resource manager will instantiate 'example_cc' on channel
	resp, err := client.UpgradeCC(
		channelID,
		resmgmt.UpgradeCCRequest{Name: name, Path: path, Version: version, Args: args, Policy: ccPolicy},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(peerName),
	)
	if err != nil {
		gnomon.Log().Error("upgrade", gnomon.Log().Err(err))
		result.Fail(err.Error())
	} else {
		result.Success(string(resp.TransactionID))
	}
	return &result
}

// fcn invoke
// args [][]byte{[]byte(coll1), []byte("key"), []byte("value")}
func invoke(chaincodeID, fcn string, args [][]byte, client *channel.Client, targetEndpoints ...string) *Result {
	result := Result{}
	resp, err := client.Execute(channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(targetEndpoints...))
	if err != nil {
		gnomon.Log().Error("invoke", gnomon.Log().Err(err))
		result.Fail(err.Error())
	} else {
		result.Success(string(resp.Payload))
		result.Msg = string(string(resp.TransactionID))
	}
	return &result
}

// fcn query
// args [][]byte{[]byte(coll1), []byte("key"), []byte("value")}
func query(chaincodeID, fcn string, args [][]byte, client *channel.Client, targetEndpoints ...string) *Result {
	result := Result{}
	resp, err := client.Query(channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(targetEndpoints...))
	if err != nil {
		gnomon.Log().Error("query", gnomon.Log().Err(err))
		result.Fail(err.Error())
	} else {
		result.Success(string(resp.Payload))
		result.Msg = string(string(resp.TransactionID))
	}
	return &result
}

func queryCollectionsConfig(orgName, orgUser, peerName, channelID, chaincodeID string, sdk *fabsdk.FabricSDK) *Result {
	result := Result{}
	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		gnomon.Log().Error("queryCollectionsConfig", gnomon.Log().Err(err))
		result.Fail("Failed to query channels: " + err.Error())
	} else {
		if nil != orgResMgmt {
			coll, err := orgResMgmt.QueryCollectionsConfig(channelID, chaincodeID, resmgmt.WithTargetEndpoints(peerName))
			if err != nil {
				gnomon.Log().Error("queryCollectionsConfig", gnomon.Log().Err(err))
				result.Fail("Failed to query channels: peer cannot be nil")
			}
			result.Success(coll)
		} else {
			gnomon.Log().Error("queryCollectionsConfig", gnomon.Log().Err(err))
			result.Fail("orgResMgmt error should be nil. ")
		}
	}
	return &result
}
