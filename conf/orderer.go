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
	"errors"
	"github.com/aberic/fabric-client/geneses"
	"github.com/aberic/fabric-client/grpc/proto/chain"
	"github.com/aberic/gnomon"
	"path/filepath"
)

// Orderer 发送事务和通道创建/更新请求
type Orderer struct {
	URL string `yaml:"url"` // URL grpcs://127.0.0.1:7050
	// GRPCOptions 这些是由gRPC库定义的标准属性，它们将按原样传递给gRPC客户端构造函数
	GRPCOptions *OrdererGRPCOptions `yaml:"grpcOptions"`
	TLSCACerts  *OrdererTLSCACerts  `yaml:"tlsCACerts"`
}

type OrdererGRPCOptions struct {
	// SSLTargetNameOverride orderer.example.com
	SSLTargetNameOverride string `yaml:"ssl-target-name-override"`
	// 这些参数应该与服务器上的keepalive策略协调设置，因为不兼容的设置可能导致连接关闭
	//
	// 当“keep-alive-time”的持续时间设置为0或更少时，将禁用keep alive客户端参数
	KeepAliveTime string `yaml:"keep-alive-time"`
	// 这些参数应该与服务器上的keepalive策略协调设置，因为不兼容的设置可能导致连接关闭
	//
	// 当“keep-alive-time”的持续时间设置为0或更少时，将禁用keep alive客户端参数
	KeepAliveTimeout string `yaml:"keep-alive-timeout"`
	// 这些参数应该与服务器上的keepalive策略协调设置，因为不兼容的设置可能导致连接关闭
	//
	// 当“keep-alive-time”的持续时间设置为0或更少时，将禁用keep alive客户端参数
	KeepAlivePermit bool `yaml:"keep-alive-permit"`
	FailFast        bool `yaml:"fail-fast"`
	// AllowInsecure 如果地址没有定义协议，则考虑允许不安全;如果为true，则考虑grpc或其他grpc
	AllowInsecure bool `yaml:"allow-insecure"`
}

type OrdererTLSCACerts struct {
	// Path 证书位置绝对路径
	Path string `yaml:"path"` // /fabric/crypto-config/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
}

func (o *Orderer) set(leagueDomain string, in *chain.Orderer) error {
	if gnomon.String().IsNotEmpty(in.Url) {
		o.URL = in.Url
	} else {
		return errors.New("url can't be empty")
	}
	if nil != in.GrpcOptions {
		if err := o.setOrdererGRPCOptions(in.GrpcOptions); nil != err {
			return err
		}
	}
	if nil != in.TlsCACerts && gnomon.String().IsNotEmpty(in.TlsCACerts.Path) {
		o.TLSCACerts.Path = in.TlsCACerts.Path
	} else {
		o.TLSCACerts.Path = filepath.Join(geneses.CryptoRootTLSCAPath(leagueDomain), geneses.CertRootTLSCAName(leagueDomain))
	}
	return nil
}

func (o *Orderer) setOrdererGRPCOptions(in *chain.OrdererGRPCOptions) error {
	if gnomon.String().IsNotEmpty(in.SslTargetNameOverride) {
		o.GRPCOptions.SSLTargetNameOverride = in.SslTargetNameOverride
	} else {
		return errors.New("ssl-target-name-override can't be empty")
	}
	if gnomon.String().IsNotEmpty(in.KeepAliveTime) {
		o.GRPCOptions.KeepAliveTime = in.KeepAliveTime
	} else {
		o.GRPCOptions.KeepAliveTime = "0s"
	}
	if gnomon.String().IsNotEmpty(in.KeepAliveTimeout) {
		o.GRPCOptions.KeepAliveTimeout = in.KeepAliveTimeout
	} else {
		o.GRPCOptions.KeepAliveTimeout = "20s"
	}
	o.GRPCOptions.KeepAlivePermit = in.KeepAlivePermit
	o.GRPCOptions.FailFast = in.FailFast
	o.GRPCOptions.AllowInsecure = in.AllowInsecure
	return nil
}
