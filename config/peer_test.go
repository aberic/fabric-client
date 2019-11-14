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

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestNewPeers(t *testing.T) {
	peers := TGetPeers()

	peersData, err := yaml.Marshal(&peers)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("--- kfk dump:\n%s\n\n", string(peersData))
}

func TGetPeers() map[string]*Peer {
	return map[string]*Peer{
		"peer0.org1.example.com": {
			URL:      "grpc://peer0.org1.example.com:7051",
			EventURL: "grpc://peer0.org1.example.com:7053",
			GRPCOptions: &PeerGRPCOptions{
				SSLTargetNameOverride: "peer0.org1.example.com",
				KeepAliveTime:         "0s",
				KeepAliveTimeout:      "20s",
				KeepAlivePermit:       false,
				FailFast:              false,
				AllowInsecure:         false,
			},
			TLSCACerts: &PeerTLSCACerts{
				Path: "peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem",
			},
		},
		"peer1.org1.example.com": {
			URL:      "grpc://peer1.org1.example.com:7051",
			EventURL: "grpc://peer1.org1.example.com:7053",
			GRPCOptions: &PeerGRPCOptions{
				SSLTargetNameOverride: "peer1.org1.example.com",
				KeepAliveTime:         "0s",
				KeepAliveTimeout:      "20s",
				KeepAlivePermit:       false,
				FailFast:              false,
				AllowInsecure:         false,
			},
			TLSCACerts: &PeerTLSCACerts{
				Path: "peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem",
			},
		},
	}
}
