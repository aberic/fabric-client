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

func (o *Order) Trans2pb() *pb.ReqOrder {
	return &pb.ReqOrder{
		ConfigID:              o.ConfigID,
		OrderName:             o.OrderName,
		Url:                   o.URL,
		SslTargetNameOverride: o.SSLTargetNameOverride,
		KeepAliveTime:         o.KeepAliveTime,
		KeepAliveTimeout:      o.KeepAliveTimeout,
		TlsCACerts:            o.TLSCACerts,
		KeepAlivePermit:       o.KeepAlivePermit,
		FailFast:              o.FailFast,
		AllowInsecure:         o.AllowInsecure,
	}
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

func (o *OrderSelf) Trans2pb() *pb.ReqOrderSelf {
	return &pb.ReqOrderSelf{
		ConfigID:         o.ConfigID,
		OrderName:        o.OrderName,
		Url:              o.URL,
		KeepAliveTime:    o.KeepAliveTime,
		KeepAliveTimeout: o.KeepAliveTimeout,
		KeepAlivePermit:  o.KeepAlivePermit,
		FailFast:         o.FailFast,
		AllowInsecure:    o.AllowInsecure,
	}
}

type OrderConfig struct {
	ConfigID  string `json:"configID"` // ConfigID 配置唯一ID
	OrgName   string `json:"orgName"`
	OrgUser   string `json:"orgUser"`
	ChannelID string `json:"channelID"`
	OrderURL  string `json:"orderURL"`
}
