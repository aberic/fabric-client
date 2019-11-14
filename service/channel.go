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

package service

import (
	pb "github.com/aberic/fabric-client/grpc/proto/chain"
)

type ChannelPeer struct {
	ConfigID       string `json:"configID"` // ConfigID 配置唯一ID
	ChannelName    string `json:"channelName"`
	PeerName       string `json:"peerName"`
	EndorsingPeer  bool   `json:"endorsingPeer"`
	ChainCodeQuery bool   `json:"chainCodeQuery"`
	LedgerQuery    bool   `json:"ledgerQuery"`
	EventSource    bool   `json:"eventSource"`
}

func (c *ChannelPeer) Trans2pb() *pb.ReqChannelPeer {
	return &pb.ReqChannelPeer{
		ConfigID:       c.ConfigID,
		ChannelName:    c.ChannelName,
		PeerName:       c.PeerName,
		EndorsingPeer:  c.EndorsingPeer,
		ChainCodeQuery: c.ChainCodeQuery,
		LedgerQuery:    c.LedgerQuery,
		EventSource:    c.EventSource,
	}
}

type ChannelPolicyQuery struct {
	ConfigID       string  `json:"configID"` // ConfigID 配置唯一ID
	ChannelName    string  `json:"channelName"`
	InitialBackOff string  `json:"initialBackOff"`
	MaxBackOff     string  `json:"maxBackOff"`
	MinResponses   int32   `json:"minResponses"`
	MaxTargets     int32   `json:"maxTargets"`
	Attempts       int32   `json:"attempts"`
	BackOffFactor  float32 `json:"backOffFactor"`
}

func (c *ChannelPolicyQuery) Trans2pb() *pb.ReqChannelPolicyQuery {
	return &pb.ReqChannelPolicyQuery{
		ConfigID:       c.ConfigID,
		ChannelName:    c.ChannelName,
		InitialBackOff: c.InitialBackOff,
		MaxBackOff:     c.MaxBackOff,
		MinResponses:   c.MinResponses,
		MaxTargets:     c.MaxTargets,
		Attempts:       c.Attempts,
		BackOffFactor:  c.BackOffFactor,
	}
}

type ChannelPolicyDiscovery struct {
	ConfigID       string `json:"configID"` // ConfigID 配置唯一ID
	ChannelName    string
	InitialBackOff string
	MaxBackOff     string
	MaxTargets     int32
	Attempts       int32
	BackOffFactor  float32
}

func (c *ChannelPolicyDiscovery) Trans2pb() *pb.ReqChannelPolicyDiscovery {
	return &pb.ReqChannelPolicyDiscovery{
		ConfigID:       c.ConfigID,
		ChannelName:    c.ChannelName,
		InitialBackOff: c.InitialBackOff,
		MaxBackOff:     c.MaxBackOff,
		MaxTargets:     c.MaxTargets,
		Attempts:       c.Attempts,
		BackOffFactor:  c.BackOffFactor,
	}
}

type ChannelPolicyEvent struct {
	ConfigID                         string `json:"configID"` // ConfigID 配置唯一ID
	ChannelName                      string `json:"channelName"`
	ResolverStrategy                 string `json:"resolverStrategy"`
	Balance                          string `json:"balance"`
	PeerMonitorPeriod                string `json:"peerMonitorPeriod"`
	BlockHeightLagThreshold          int64  `json:"blockHeightLagThreshold"`
	ReconnectBlockHeightLagThreshold int64  `json:"reconnectBlockHeightLagThreshold"`
}

func (c *ChannelPolicyEvent) Trans2pb() *pb.ReqChannelPolicyEvent {
	return &pb.ReqChannelPolicyEvent{
		ConfigID:                         c.ConfigID,
		ChannelName:                      c.ChannelName,
		ResolverStrategy:                 c.ResolverStrategy,
		Balance:                          c.Balance,
		PeerMonitorPeriod:                c.PeerMonitorPeriod,
		BlockHeightLagThreshold:          c.BlockHeightLagThreshold,
		ReconnectBlockHeightLagThreshold: c.ReconnectBlockHeightLagThreshold,
	}
}

type ChannelCreate struct {
	ConfigID   string `json:"configID"` // ConfigID 配置唯一ID
	LeagueName string `json:"leagueName"`
	ChannelID  string `json:"channelID"`
}

type ChannelJoin struct {
	ConfigID  string `json:"configID"` // ConfigID 配置唯一ID
	OrgName   string `json:"orgName"`
	OrgUser   string `json:"orgUser"`
	ChannelID string `json:"channelID"`
}

type ChannelList struct {
	ConfigID string `json:"configID"` // ConfigID 配置唯一ID
	OrgName  string `json:"orgName"`
	OrgUser  string `json:"orgUser"`
	PeerName string `json:"peerName"`
}
