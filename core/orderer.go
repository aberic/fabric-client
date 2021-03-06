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
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func queryConfigFromOrderer(orgName, orgUser, channelID, orderURL string, sdk *fabsdk.FabricSDK) *Result {
	result := Result{}
	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		gnomon.Log().Error("queryConfigFromOrderer", gnomon.Log().Err(err))
		result.Fail("Failed to query config from orderer: " + err.Error())
	} else {
		channelCfg, err := orgResMgmt.QueryConfigFromOrderer(channelID, resmgmt.WithOrdererEndpoint(orderURL), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			gnomon.Log().Error("queryConfigFromOrderer", gnomon.Log().Err(err))
			result.Fail("QueryConfig return error: " + err.Error())
		} else {
			result.Success(channelCfg.Orderers())
		}
	}
	return &result
}
