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

func TestNewClient(t *testing.T) {
	client := TGetClient()
	clientData, err := yaml.Marshal(&client)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("--- dump:\n%s\n\n", string(clientData))
}

func TGetClient() *Client {
	client := Client{}
	client.initClient(true, "Org1", "debug",
		"/crypto-config",
		"/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.key",
		"/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.crt")
	return &client
}
