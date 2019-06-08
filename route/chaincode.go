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

package route

import (
	"github.com/ennoo/fabric-go-client/core"
	"github.com/ennoo/rivet"
	"github.com/ennoo/rivet/trans/response"
)

const (
	channelID      = "mychannel"
	ChaincodeID    = "testcc"
	orgName        = "Org1"
	orgUser        = "Admin"
	ordererOrgName = "OrdererOrg"
)

func ChainCode(router *response.Router) {
	// 仓库相关路由设置
	router.Group = router.Engine.Group("/code")
	router.GET("/invoke", invoke)
	router.GET("/query", query)
}

func invoke(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		sdk.Invoke(ChaincodeID, orgName, orgUser, channelID, "invoke", [][]byte{[]byte("A"), []byte("B"), []byte("10")},
			"/Users/aberic/Documents/path/go/src/github.com/ennoo/fabric-go-client/config_e2e.yaml").Say(router.Context)
	})
}

func query(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		sdk.Query(ChaincodeID, orgName, orgUser, channelID, "query", [][]byte{[]byte("A")}, []string{},
			"/Users/aberic/Documents/path/go/src/github.com/ennoo/fabric-go-client/config_e2e.yaml").Say(router.Context)
	})
}
