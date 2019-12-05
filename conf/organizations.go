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
	"github.com/aberic/fabric-client/geneses"
	"github.com/aberic/fabric-client/grpc/proto/chain"
	"github.com/aberic/gnomon"
)

// Organization 此网络的参与机构
type Organization struct {
	MspID      string           `yaml:"mspid"`
	CryptoPath string           `yaml:"cryptoPath"` // CryptoPath 这个组织的MSP存储(绝对路径或相对于client.cryptoconfig)
	Users      map[string]*User `yaml:"users,omitempty"`
	Peers      []string         `yaml:"peers,omitempty"`
	// CertificateAuthorities
	//
	// 在基于Fabric的网络中，证书颁发机构颁发证书用于身份验证
	//
	// 通常，证书供应是在运行时网络之外的独立进程中完成的
	//
	// ca是一个特殊的证书颁发机构，它为动态证书管理(注册、撤销、重新注册)提供了REST api
	//
	// CertificateAuthorities只针对fabric-ca服务
	CertificateAuthorities []string `yaml:"certificateAuthorities,omitempty"`
}

type User struct {
	Cert *Cert `yaml:"cert"`
}

type Cert struct {
	Path string `yaml:"path"`
}

func (o *Organization) set(leagueName, orgName string, in *chain.Organization) {
	if gnomon.String().IsNotEmpty(in.MspID) {
		o.MspID = in.MspID
	} else {
		o.MspID = geneses.MspID(orgName)
	}
	if gnomon.String().IsNotEmpty(in.CryptoPath) {
		o.CryptoPath = in.CryptoPath
	} else {
		o.CryptoPath = geneses.CryptoConfigPath(leagueName)
	}
	if len(in.Peers) > 0 {
		o.Peers = in.Peers
	}
	if len(in.CertificateAuthorities) > 0 {
		o.CertificateAuthorities = in.CertificateAuthorities
	}
}
