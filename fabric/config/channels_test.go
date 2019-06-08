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

package config

import (
	"fmt"
	"github.com/ennoo/rivet/utils/log"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestNewChannels(t *testing.T) {
	channels := TGetChannels()

	channelsData, err := yaml.Marshal(&channels)
	if err != nil {
		log.Self.Error("channels", log.Error(err))
	}
	fmt.Printf("--- kfk dump:\n%s\n\n", string(channelsData))
}

func TGetChannels() map[string]*Channel {
	return map[string]*Channel{
		"mychannel1": {
			Peers: map[string]*ChannelPeer{
				"peer0.org1.example.com": {
					EndorsingPeer:  true,
					ChaincodeQuery: true,
					LedgerQuery:    true,
					EventSource:    true,
				},
				"peer1.org1.example.com": {
					EndorsingPeer:  true,
					ChaincodeQuery: true,
					LedgerQuery:    true,
					EventSource:    true,
				},
			},
			Policies: &Policy{
				QueryChannelConfig: &PolicyQueryChannelConfig{
					MinResponses: 1,
					MaxTargets:   1,
					RetryOpts: &PolicyCommonRetryOpts{
						Attempts:       5,
						InitialBackOff: "500ms",
						MaxBackOff:     "5s",
						BackOffFactor:  2.0,
					},
				},
				Discovery: &PolicyDiscovery{
					MaxTargets: 2,
					RetryOpts: &PolicyCommonRetryOpts{
						Attempts:       4,
						InitialBackOff: "500ms",
						MaxBackOff:     "5s",
						BackOffFactor:  2.0,
					},
				},
				EventService: &PolicyEventService{
					ResolverStrategy:                 "PreferOrg",
					Balancer:                         "Random",
					BlockHeightLagThreshold:          5,
					ReconnectBlockHeightLagThreshold: 8,
					PeerMonitorPeriod:                "6s",
				},
			},
		},
		"mychannel2": {
			Peers: map[string]*ChannelPeer{
				"peer0.org2.example.com": {
					EndorsingPeer:  true,
					ChaincodeQuery: true,
					LedgerQuery:    true,
					EventSource:    true,
				},
				"peer1.org2.example.com": {
					EndorsingPeer:  true,
					ChaincodeQuery: true,
					LedgerQuery:    true,
					EventSource:    true,
				},
			},
			Policies: &Policy{
				QueryChannelConfig: &PolicyQueryChannelConfig{
					MinResponses: 1,
					MaxTargets:   1,
					RetryOpts: &PolicyCommonRetryOpts{
						Attempts:       5,
						InitialBackOff: "500ms",
						MaxBackOff:     "5s",
						BackOffFactor:  2.0,
					},
				},
				Discovery: &PolicyDiscovery{
					MaxTargets: 2,
					RetryOpts: &PolicyCommonRetryOpts{
						Attempts:       4,
						InitialBackOff: "500ms",
						MaxBackOff:     "5s",
						BackOffFactor:  2.0,
					},
				},
				EventService: &PolicyEventService{
					ResolverStrategy:                 "PreferOrg",
					Balancer:                         "Random",
					BlockHeightLagThreshold:          5,
					ReconnectBlockHeightLagThreshold: 8,
					PeerMonitorPeriod:                "6s",
				},
			},
		},
	}
}
