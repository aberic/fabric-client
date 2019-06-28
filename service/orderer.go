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

type Order struct {
	ConfigID              string `json:"configID"` // ConfigID 配置唯一ID
	OrderName             string `json:"orderName"`
	URL                   string `json:"url"`
	SSLTargetNameOverride string `json:"sslTargetNameOverride"`
	KeepAliveTime         string `json:"keepAliveTime"`
	KeepAliveTimeout      string `json:"keepAliveTimeout"`
	TLSCACerts            string `json:"tlsCACerts"`
	KeepAlivePermit       bool   `json:"keepAlivePermit"`
	FailFast              bool   `json:"failFast"`
	AllowInsecure         bool   `json:"allowInsecure"`
}

type OrderSelf struct {
	ConfigID         string `json:"configID"` // ConfigID 配置唯一ID
	LeagueName       string `json:"leagueName"`
	OrderName        string `json:"orderName"`
	URL              string `json:"url"`
	KeepAliveTime    string `json:"keepAliveTime"`
	KeepAliveTimeout string `json:"keepAliveTimeout"`
	KeepAlivePermit  bool   `json:"keepAlivePermit"`
	FailFast         bool   `json:"failFast"`
	AllowInsecure    bool   `json:"allowInsecure"`
}

type OrderConfig struct {
	ConfigID  string `json:"configID"` // ConfigID 配置唯一ID
	OrgName   string `json:"orgName"`
	OrgUser   string `json:"orgUser"`
	ChannelID string `json:"channelID"`
	OrderURL  string `json:"orderURL"`
}
