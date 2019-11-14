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
 */

package config

type Peer struct {
	URL         string           `yaml:"url"`
	EventURL    string           `yaml:"eventUrl"`
	GRPCOptions *PeerGRPCOptions `yaml:"grpcOptions"`
	TLSCACerts  *PeerTLSCACerts  `yaml:"tlsCACerts"`
}

type PeerGRPCOptions struct {
	SSLTargetNameOverride string `yaml:"ssl-target-name-override"`
	KeepAliveTime         string `yaml:"keep-alive-time"`
	KeepAliveTimeout      string `yaml:"keep-alive-timeout"`
	KeepAlivePermit       bool   `yaml:"keep-alive-permit"`
	FailFast              bool   `yaml:"fail-fast"`
	AllowInsecure         bool   `yaml:"allow-insecure"`
}

type PeerTLSCACerts struct {
	Path string `yaml:"path"` // /fabric/crypto-config/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
}
