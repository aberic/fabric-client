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

import (
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
)

type Peer struct {
	ConfigID              string `json:"configID"` // ConfigID 配置唯一ID
	PeerName              string `json:"peerName"`
	URL                   string `json:"url"`
	EventUrl              string `json:"eventUrl"`
	SSLTargetNameOverride string `json:"sslTargetNameOverride"`
	KeepAliveTime         string `json:"keepAliveTime"`
	KeepAliveTimeout      string `json:"keepAliveTimeout"`
	TLSCACerts            string `json:"tlsCACerts"`
	KeepAlivePermit       bool   `json:"keepAlivePermit"`
	FailFast              bool   `json:"failFast"`
	AllowInsecure         bool   `json:"allowInsecure"`
}

func (p *Peer) Trans2pb() *pb.ReqPeer {
	return &pb.ReqPeer{
		ConfigID:              p.ConfigID,
		PeerName:              p.PeerName,
		Url:                   p.URL,
		EventUrl:              p.EventUrl,
		SslTargetNameOverride: p.SSLTargetNameOverride,
		KeepAliveTime:         p.KeepAliveTime,
		KeepAliveTimeout:      p.KeepAliveTimeout,
		TlsCACerts:            p.TLSCACerts,
		KeepAlivePermit:       p.KeepAlivePermit,
		FailFast:              p.FailFast,
		AllowInsecure:         p.AllowInsecure,
	}
}

type PeerSelf struct {
	ConfigID         string `json:"configID"` // ConfigID 配置唯一ID
	LeagueName       string `json:"leagueName"`
	PeerName         string `json:"peerName"`
	URL              string `json:"url"`
	EventUrl         string `json:"eventUrl"`
	KeepAliveTime    string `json:"keepAliveTime"`
	KeepAliveTimeout string `json:"keepAliveTimeout"`
	KeepAlivePermit  bool   `json:"keepAlivePermit"`
	FailFast         bool   `json:"failFast"`
	AllowInsecure    bool   `json:"allowInsecure"`
}

func (p *PeerSelf) Trans2pb() *pb.ReqPeerSelf {
	return &pb.ReqPeerSelf{
		ConfigID:         p.ConfigID,
		LeagueName:       p.LeagueName,
		PeerName:         p.PeerName,
		Url:              p.URL,
		EventUrl:         p.EventUrl,
		KeepAliveTime:    p.KeepAliveTime,
		KeepAliveTimeout: p.KeepAliveTimeout,
		KeepAlivePermit:  p.KeepAlivePermit,
		FailFast:         p.FailFast,
		AllowInsecure:    p.AllowInsecure,
	}
}
