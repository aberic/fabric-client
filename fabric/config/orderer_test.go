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

package config

import (
	"fmt"
	"github.com/ennoo/rivet/utils/log"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestNewOrderers(t *testing.T) {
	orderers := TGetOrderers()

	orderersData, err := yaml.Marshal(&orderers)
	if err != nil {
		log.Self.Error("channels", log.Error(err))
	}
	fmt.Printf("--- kfk dump:\n%s\n\n", string(orderersData))
}

func TGetOrderers() map[string]*Orderer {
	return map[string]*Orderer{
		"orderer0.example.com": {
			URL: "grpc://orderer0.example.com:7050",
			GRPCOptions: &OrdererGRPCOptions{
				SSLTargetNameOverride: "orderer0.example.com",
				KeepAliveTime:         "0s",
				KeepAliveTimeout:      "20s",
				KeepAlivePermit:       false,
				FailFast:              false,
				AllowInsecure:         false,
			},
			TLSCACerts: &OrdererTLSCACerts{
				Path: "ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem",
			},
		},
		"orderer1.example.com": {
			URL: "grpc://orderer1.example.com:7050",
			GRPCOptions: &OrdererGRPCOptions{
				SSLTargetNameOverride: "orderer1.example.com",
				KeepAliveTime:         "0s",
				KeepAliveTimeout:      "20s",
				KeepAlivePermit:       false,
				FailFast:              false,
				AllowInsecure:         false,
			},
			TLSCACerts: &OrdererTLSCACerts{
				Path: "ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem",
			},
		},
	}
}
