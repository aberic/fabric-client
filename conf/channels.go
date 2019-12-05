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

package conf

import (
	"github.com/aberic/fabric-client/grpc/proto/chain"
	"github.com/aberic/gnomon"
)

type Channel struct {
	Peers    map[string]*ChannelPeer `yaml:"peers"`    // Peers 可选参数，参与组织的节点列表
	Policies *Policy                 `yaml:"policies"` // Policies 可选参数，应用程序可以使用这些选项来执行通道操作，如检索通道配置等。
}

// ChannelPeer 可选参数，参与组织的节点列表
type ChannelPeer struct {
	// EndorsingPeer
	//
	// 可选参数
	//
	// 是否会向该节点发送交易以供其背书，节点必须安装chaincode。应用程序也可以使用这个属性来决定发送chaincode安装请求的节点。
	//
	// 默认值:true
	EndorsingPeer bool `yaml:"endorsingPeer"`
	// ChaincodeQuery
	//
	// 可选参数
	//
	// 这个节点是否可被发送查询建议，节点必须安装chaincode。应用程序也可以使用这个属性来决定发送chaincode安装请求的节点。
	//
	// 默认值:true
	ChaincodeQuery bool `yaml:"chaincodeQuery"`
	// LedgerQuery
	//
	// 可选参数
	//
	// 这个节点是否可发送不需要链码的查询建议，如queryBlock()、queryTransaction()等。
	//
	// 默认值:true
	LedgerQuery bool `yaml:"ledgerQuery"`
	// EventSource
	//
	// 可选参数
	//
	// 这个节点是否是SDK侦听器注册的目标，所有的对等点都可以产生事件，但应用程序通常只需要连接一个来监听事件。
	//
	// 默认值:true
	EventSource bool `yaml:"eventSource"`
}

// Policy 可选参数，应用程序可以使用这些选项来执行通道操作，如检索通道配置等
type Policy struct {
	QueryChannelConfig *PolicyQueryChannelConfig `yaml:"queryChannelConfig"` // PolicyQueryChannelConfig 可选参数，用于检索通道配置块的选项
	Discovery          *PolicyDiscovery          `yaml:"discovery"`          // Discovery 可选参数，检索发现信息的选项
	EventService       *PolicyEventService       `yaml:"eventService"`       // EventService 可选参数，事件服务的选项
}

// PolicyQueryChannelConfig 可选参数，用于检索通道配置块的选项
type PolicyQueryChannelConfig struct {
	MinResponses int32                  `yaml:"minResponses"` // MinResponses 可选参数，最小成功响应数(来自目标/节点)
	MaxTargets   int32                  `yaml:"maxTargets"`   // MaxTargets 可选参数，通道配置将为这些数目的随机目标检索
	RetryOpts    *PolicyCommonRetryOpts `yaml:"retryOpts"`    // RetryOpts 可选参数，查询配置块的重试选项
}

// PolicyCommonRetryOpts 可选参数，查询配置块的重试选项
type PolicyCommonRetryOpts struct {
	Attempts       int32   `yaml:"attempts"`       // Attempts 可选参数，number of retry attempts
	InitialBackOff string  `yaml:"initialBackoff"` // InitialBackOff 可选参数，第一次重试尝试的回退间隔
	MaxBackOff     string  `yaml:"maxBackoff"`     // MaxBackOff 可选参数，任何重试尝试的最大回退间隔
	BackOffFactor  float32 `yaml:"backoffFactor"`  // BackOffFactor 可选参数，该因子使初始回退期呈指数递增
}

// PolicyDiscovery 可选参数，检索发现信息的选项
type PolicyDiscovery struct {
	MaxTargets int32                  `yaml:"maxTargets"` // MaxTargets 可选参数，发现信息将检索这些随机目标的数量
	RetryOpts  *PolicyCommonRetryOpts `yaml:"retryOpts"`  // RetryOpts 可选参数，检索发现信息的重试选项
}

// PolicyEventService 可选参数，事件服务的选项
type PolicyEventService struct {
	// ResolverStrategy
	//
	// 可选参数
	//
	// PreferOrg:
	// 根据块高度滞后阈值确定哪些对等点是合适的，尽管它们更适用当前组织中的对等点(只要它们的块高度高于配置的阈值)。如果当前组织中的对等点都不合适，则选择另一个组织中的对等点
	//
	// MinBlockHeight:
	// 根据块高度滞后阈值选择最佳的对等点。所有对等点的最大块高度被确定，那些块高度低于最大高度但高于规定的“滞后”阈值的对等点被负载均衡。不考虑其他节点
	//
	// Balanced:
	// 使用配置的平衡器选择对等点
	ResolverStrategy string `yaml:"resolverStrategy"`
	// Balancer
	//
	// 可选参数
	//
	// 当选择一个对等点连接到可能的值时使用的负载均衡[Random (default), RoundRobin]
	Balancer string `yaml:"balancer"`
	// BlockHeightLagThreshold
	//
	// 可选参数
	//
	// 设置块高度滞后阈值。此值用于选择要连接的对等点。如果一个节点落后于最新的节点超过给定的块数，那么它将被排除在选择之外
	// 注意，此参数仅适用于minBlockHeightResolverMode设置为ResolveByThreshold时
	// 默认值:5
	BlockHeightLagThreshold int64 `yaml:"blockHeightLagThreshold"`
	// ReconnectBlockHeightLagThreshold
	//
	// 可选参数
	//
	// reconnectBlockHeightLagThreshold—如果对等方的块高度低于指定的块数，则事件客户机将断开与对等方的连接，并重新连接到性能更好的对等方
	//
	// 注意，此参数仅适用于peerMonitor设置为Enabled(默认)的情况
	//
	// 默认值:10
	//
	// 注意:设置此值过低可能会导致事件客户端过于频繁地断开/重新连接，从而影响性能
	ReconnectBlockHeightLagThreshold int64 `yaml:"reconnectBlockHeightLagThreshold"`
	// PeerMonitorPeriod
	//
	// 可选参数
	//
	// peerMonitorPeriod是监视连接的对等点以查看事件客户端是否应该断开连接并重新连接到另一个对等点的时间段
	//
	// 默认:0(禁用)用于平衡冲突解决策略;优先级和MinBlockHeight策略的5s
	PeerMonitorPeriod string `yaml:"peerMonitorPeriod"`
}

func NewConfigChannel() *Channel {
	return &Channel{
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
				ResolverStrategy:                 "Balanced",
				Balancer:                         "RoundRobin",
				BlockHeightLagThreshold:          5,
				ReconnectBlockHeightLagThreshold: 10,
				PeerMonitorPeriod:                "6s",
			},
		},
		Peers: map[string]*ChannelPeer{},
	}
}

func (c *Channel) set(in *chain.Channel) {
	if nil != in.Policies {
		c.setPolicy(in.Policies)
	}
	if len(in.Peers) > 0 {
		c.setPeers(in)
	}
}

func (c *Channel) setPolicy(in *chain.Policy) {
	if nil != in.QueryChannelConfig {
		c.setQueryChannelConfig(in.QueryChannelConfig)
	}
	if nil != in.Discovery {
		c.setDiscovery(in.Discovery)
	}
	if nil != in.EventService {
		c.setEventService(in.EventService)
	}
}

func (c *Channel) setQueryChannelConfig(in *chain.PolicyQueryChannelConfig) {
	if in.MinResponses > 0 {
		c.Policies.QueryChannelConfig.MinResponses = in.MinResponses
	}
	if in.MaxTargets > 0 {
		c.Policies.QueryChannelConfig.MaxTargets = in.MaxTargets
	}
	if nil != in.RetryOpts {
		if in.RetryOpts.Attempts > 0 {
			c.Policies.QueryChannelConfig.RetryOpts.Attempts = in.RetryOpts.Attempts
		}
		if in.RetryOpts.BackoffFactor > 0 {
			c.Policies.QueryChannelConfig.RetryOpts.BackOffFactor = in.RetryOpts.BackoffFactor
		}
		if gnomon.String().IsNotEmpty(in.RetryOpts.InitialBackoff) {
			c.Policies.QueryChannelConfig.RetryOpts.InitialBackOff = in.RetryOpts.InitialBackoff
		}
		if gnomon.String().IsNotEmpty(in.RetryOpts.MaxBackoff) {
			c.Policies.QueryChannelConfig.RetryOpts.MaxBackOff = in.RetryOpts.MaxBackoff
		}
	}
}

func (c *Channel) setDiscovery(in *chain.PolicyDiscovery) {
	if in.MaxTargets > 0 {
		c.Policies.Discovery.MaxTargets = in.MaxTargets
	}
	if nil != in.RetryOpts {
		if in.RetryOpts.Attempts > 0 {
			c.Policies.Discovery.RetryOpts.Attempts = in.RetryOpts.Attempts
		}
		if in.RetryOpts.BackoffFactor > 0 {
			c.Policies.Discovery.RetryOpts.BackOffFactor = in.RetryOpts.BackoffFactor
		}
		if gnomon.String().IsNotEmpty(in.RetryOpts.InitialBackoff) {
			c.Policies.Discovery.RetryOpts.InitialBackOff = in.RetryOpts.InitialBackoff
		}
		if gnomon.String().IsNotEmpty(in.RetryOpts.MaxBackoff) {
			c.Policies.Discovery.RetryOpts.MaxBackOff = in.RetryOpts.MaxBackoff
		}
	}
}

func (c *Channel) setEventService(in *chain.PolicyEventService) {
	if gnomon.String().IsNotEmpty(in.ResolverStrategy) {
		c.Policies.EventService.ResolverStrategy = in.ResolverStrategy
	}
	if gnomon.String().IsNotEmpty(in.Balancer) {
		c.Policies.EventService.Balancer = in.Balancer
	}
	if in.BlockHeightLagThreshold > 0 {
		c.Policies.EventService.BlockHeightLagThreshold = in.BlockHeightLagThreshold
	}
	if in.ReconnectBlockHeightLagThreshold > 0 {
		c.Policies.EventService.ReconnectBlockHeightLagThreshold = in.ReconnectBlockHeightLagThreshold
	}
	if gnomon.String().IsNotEmpty(in.PeerMonitorPeriod) {
		c.Policies.EventService.PeerMonitorPeriod = in.PeerMonitorPeriod
	}
}

func (c *Channel) setPeers(in *chain.Channel) {
	for peerName, peer := range in.Peers {
		c.Peers[peerName] = &ChannelPeer{
			EndorsingPeer:  peer.EndorsingPeer,
			ChaincodeQuery: peer.ChaincodeQuery,
			LedgerQuery:    peer.LedgerQuery,
			EventSource:    peer.EventSource,
		}
	}
}
