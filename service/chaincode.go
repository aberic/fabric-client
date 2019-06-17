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

type Install struct {
	ConfigID     string `json:"configID"` // ConfigID 配置唯一ID
	OrderOrgName string `json:"orderOrgName"`
	OrgUser      string `json:"orgUser"`
	Name         string `json:"name"`
	Source       string `json:"source"`
	Path         string `json:"path"`
	Version      string `json:"version"`
}

type Installed struct {
	ConfigID string `json:"configID"` // ConfigID 配置唯一ID
	OrgName  string `json:"orgName"`
	OrgUser  string `json:"orgUser"`
	PeerName string `json:"peerName"`
}

type Instantiate struct {
	ConfigID     string   `json:"configID"` // ConfigID 配置唯一ID
	OrderOrgName string   `json:"orderOrgName"`
	OrgUser      string   `json:"orgUser"`
	ChannelID    string   `json:"channelID"`
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Version      string   `json:"version"`
	OrgPolicies  []string `json:"orgPolicies"`
	Args         [][]byte `json:"args"`
}

type Instantiated struct {
	Installed
	ChannelID string `json:"channelID"`
}

type Upgrade struct {
	Instantiate
}

type Invoke struct {
	ConfigID    string   `json:"configID"` // ConfigID 配置唯一ID
	ChannelID   string   `json:"channelID"`
	ChainCodeID string   `json:"chainCodeID"`
	OrgName     string   `json:"orgName"`
	OrgUser     string   `json:"orgUser"`
	Fcn         string   `json:"fcn"`
	Args        [][]byte `json:"args"`
}

type Query struct {
	Invoke
	TargetEndpoints []string `json:"targetEndpoints"`
}

type ChainCodeCollectionsConfig struct {
	ConfigID    string `json:"configID"` // ConfigID 配置唯一ID
	ChainCodeID string `json:"chainCodeID"`
	OrgName     string `json:"orgName"`
	OrgUser     string `json:"orgUser"`
	ChannelID   string `json:"channelID"`
	PeerName    string `json:"peerName"`
}
