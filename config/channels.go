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

type Channel struct {
	Peers    map[string]*ChannelPeer `yaml:"peers"`
	Policies *Policy                 `yaml:"policies"`
}

type ChannelPeer struct {
	EndorsingPeer  bool `yaml:"endorsingPeer"`
	ChaincodeQuery bool `yaml:"chaincodeQuery"`
	LedgerQuery    bool `yaml:"ledgerQuery"`
	EventSource    bool `yaml:"eventSource"`
}

type Policy struct {
	QueryChannelConfig *PolicyQueryChannelConfig `yaml:"queryChannelConfig"`
	Discovery          *PolicyDiscovery          `yaml:"discovery"`
	EventService       *PolicyEventService       `yaml:"eventService"`
}

type PolicyQueryChannelConfig struct {
	MinResponses int32                  `yaml:"minResponses"`
	MaxTargets   int32                  `yaml:"maxTargets"`
	RetryOpts    *PolicyCommonRetryOpts `yaml:"retryOpts"`
}

type PolicyCommonRetryOpts struct {
	Attempts       int32   `yaml:"attempts"`
	InitialBackOff string  `yaml:"initialBackoff"`
	MaxBackOff     string  `yaml:"maxBackoff"`
	BackOffFactor  float32 `yaml:"backoffFactor"`
}

type PolicyDiscovery struct {
	MaxTargets int32                  `yaml:"maxTargets"`
	RetryOpts  *PolicyCommonRetryOpts `yaml:"retryOpts"`
}

type PolicyEventService struct {
	ResolverStrategy                 string `yaml:"resolverStrategy"`
	Balancer                         string `yaml:"balancer"`
	BlockHeightLagThreshold          int64  `yaml:"blockHeightLagThreshold"`
	ReconnectBlockHeightLagThreshold int64  `yaml:"reconnectBlockHeightLagThreshold"`
	PeerMonitorPeriod                string `yaml:"peerMonitorPeriod"`
}
