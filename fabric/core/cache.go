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
 */

package sdk

import (
	"github.com/ennoo/rivet/utils/env"
	"github.com/ennoo/rivet/utils/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"strings"
	"sync"
)

const (
	ConfigYmlPath = "CONFIG_YML_PATH"
)

var (
	sdkPool = sync.Pool{
		New: func() interface{} {
			log.Self.Debug("make new sdk in pool")
			return obtainSDK(env.GetEnvDefault(ConfigYmlPath, "/Users/aberic/Documents/path/go/src/github.com/ennoo/fom/config_e2e.yaml"))
		}}
	channelClientPools = map[string]chan *channel.Client{}
)

func obtainSDK(configYmlPath string, sdkOpts ...fabsdk.Option) *fabsdk.FabricSDK {
	sdk, err := sdk(configYmlPath, sdkOpts...)
	if err != nil {
		log.Self.Error(err.Error())
		return nil
	}
	return sdk
}

func ccpChan(orgName, channelID, username string, configYmlPath string, sdkOpts ...fabsdk.Option) chan *channel.Client {
	ccpStr := strings.Join([]string{orgName, channelID, username}, "-")
	if nil == channelClientPools[ccpStr] {
		ccpChan := make(chan *channel.Client, 100)
		sdk, err := sdk(configYmlPath, sdkOpts...)
		if err != nil {
			log.Self.Error(err.Error())
			return nil
		}
		//defer sdk.Close()
		//prepare channel client context using client context
		clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(username), fabsdk.WithOrg(orgName))
		go func() {
			for i := 0; i < 100; i++ {
				client, err := channel.New(clientChannelContext)
				if err != nil {
					log.Self.Error("Failed to create new channel client", log.Int("i", i), log.Error(err))
					continue
				}
				log.Self.Debug("make client", log.Int("i", i))
				ccpChan <- client
			}
		}()
		channelClientPools[ccpStr] = ccpChan
	}
	return channelClientPools[ccpStr]
}

// SyncPoolGetSDK
func SyncPoolGetSDK() *fabsdk.FabricSDK {
	return sdkPool.Get().(*fabsdk.FabricSDK)
}
