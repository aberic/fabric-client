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

type Organization struct {
	MspID                  string           `yaml:"mspid"`
	CryptoPath             string           `yaml:"cryptoPath"`
	Users                  map[string]*User `yaml:"users,omitempty"`
	Peers                  []string         `yaml:"peers,omitempty"`
	CertificateAuthorities []string         `yaml:"certificateAuthorities,omitempty"`
}

type User struct {
	Cert *Cert `yaml:"cert"`
}

type Cert struct {
	Path string `yaml:"path"`
}
