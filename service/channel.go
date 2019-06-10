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

package service

type ChannelPeer struct {
	ConfigID       string `json:"configID"` // ConfigID 配置唯一ID
	ChannelName    string `json:"channelName"`
	PeerName       string `json:"peerName"`
	EndorsingPeer  bool   `json:"endorsingPeer"`
	ChainCodeQuery bool   `json:"chainCodeQuery"`
	LedgerQuery    bool   `json:"ledgerQuery"`
	EventSource    bool   `json:"eventSource"`
}

type ChannelPolicyQuery struct {
	ConfigID       string  `json:"configID"` // ConfigID 配置唯一ID
	ChannelName    string  `json:"channelName"`
	InitialBackOff string  `json:"initialBackOff"`
	MaxBackOff     string  `json:"maxBackOff"`
	MinResponses   int     `json:"minResponses"`
	MaxTargets     int     `json:"maxTargets"`
	Attempts       int     `json:"attempts"`
	BackOffFactor  float32 `json:"backOffFactor"`
}

type ChannelPolicyDiscovery struct {
	ConfigID       string `json:"configID"` // ConfigID 配置唯一ID
	ChannelName    string
	InitialBackOff string
	MaxBackOff     string
	MaxTargets     int
	Attempts       int
	BackOffFactor  float32
}

type ChannelPolicyEvent struct {
	ConfigID                         string `json:"configID"` // ConfigID 配置唯一ID
	ChannelName                      string `json:"channelName"`
	ResolverStrategy                 string `json:"resolverStrategy"`
	Balance                          string `json:"balance"`
	PeerMonitorPeriod                string `json:"peerMonitorPeriod"`
	BlockHeightLagThreshold          int    `json:"blockHeightLagThreshold"`
	ReconnectBlockHeightLagThreshold int    `json:"reconnectBlockHeightLagThreshold"`
}

type ChannelCreate struct {
	ConfigID          string `json:"configID"` // ConfigID 配置唯一ID
	OrderOrgName      string `json:"orderOrgName"`
	OrgName           string `json:"orgName"`
	OrgUser           string `json:"orgUser"`
	ChannelID         string `json:"channelID"`
	ChannelConfigPath string `json:"channelConfigPath"`
}

type ChannelJoin struct {
	ConfigID  string `json:"configID"` // ConfigID 配置唯一ID
	OrderName string `json:"orderName"`
	OrgName   string `json:"orgName"`
	OrgUser   string `json:"orgUser"`
	ChannelID string `json:"channelID"`
	PeerUrl   string `json:"peerUrl"`
}

type ChannelList struct {
	ConfigID string `json:"configID"` // ConfigID 配置唯一ID
	OrgName  string `json:"orgName"`
	OrgUser  string `json:"orgUser"`
	PeerName string `json:"peerName"`
}
