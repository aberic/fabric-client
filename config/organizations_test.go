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

func TestNewOrganizations(t *testing.T) {
	organizations := TGetOrganizations()

	organizationsData, err := yaml.Marshal(&organizations)
	if err != nil {
		log.Self.Error("channels", log.Error(err))
	}
	fmt.Printf("--- kfk dump:\n%s\n\n", string(organizationsData))
}

func TGetOrganizations() map[string]interface{} {
	return map[string]interface{}{
		"ordererorg": &OrdererOrg{
			MspID:      "OrdererMSP",
			CryptoPath: "/fabric/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/msp",
		},
		"Org1": &Org{
			MspID:                  "Org1MSP",
			CryptoPath:             "/fabric/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp",
			Peers:                  []string{"peer0.org1.example.com", "peer1.org1.example.com"},
			CertificateAuthorities: []string{"ca.org1.example.com"},
		},
		"Org2": &Org{
			MspID:                  "Org1MSP",
			CryptoPath:             "/fabric/crypto-config/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp",
			Peers:                  []string{"peer0.org2.example.com", "peer1.org2.example.com"},
			CertificateAuthorities: []string{"ca.org2.example.com"},
		},
	}
}
