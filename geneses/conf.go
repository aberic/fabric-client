/*
 * Copyright (c) 2019.. ENNOO - All Rights Reserved.
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

package geneses

// Generate 请求生成区块链配置对象
type Generate struct {
	LedgerName      string `json:"ledgerName"`
	OrderCount      int    `json:"orderCount"`
	PeerCount       int    `json:"peerCount"`
	TemplateCount   int    `json:"templateCount"`
	UserCount       int    `json:"userCount"`
	BatchTimeout    int    `json:"batchTimeout"`
	MaxMessageCount int    `json:"maxMessageCount"`
	Force           bool   `json:"force"`
}

// GenerateCustom 请求生成区块链自定义配置对象
type GenerateCustom struct {
	LedgerName string `json:"ledgerName"`
	Force      bool   `json:"force"`
}

// Crypto 请求生成区块链配置文件集合对象
type Crypto struct {
	LedgerName string `json:"ledgerName"`
	Force      bool   `json:"force"`
}

// ChannelTX 请求生成区块链配置文件集合对象
type ChannelTX struct {
	LedgerName  string `json:"ledgerName"`
	ChannelName string `json:"channelName"`
	Force       bool   `json:"force"`
}
