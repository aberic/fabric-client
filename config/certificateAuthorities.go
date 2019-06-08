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

type CertificateAuthority struct {
	URL        string                          `yaml:"url"`
	CAName     string                          `yaml:"caName"`
	TLSCACerts *CertificateAuthorityTLSCACerts `yaml:"tlsCACerts"`
	Registrar  *CertificateAuthorityRegistrar  `yaml:"registrar"`
}

type CertificateAuthorityTLSCACerts struct {
	Path   string                                `yaml:"path"`
	Client *CertificateAuthorityTLSCACertsClient `yaml:"client"`
}

type CertificateAuthorityTLSCACertsClient struct {
	Key  *CertificateAuthorityTLSCACertsClientKey  `yaml:"key"`
	Cert *CertificateAuthorityTLSCACertsClientCert `yaml:"cert"`
}

type CertificateAuthorityTLSCACertsClientKey struct {
	Path string `yaml:"path"` // /fabric/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.key
}

type CertificateAuthorityTLSCACertsClientCert struct {
	Path string `yaml:"path"` // /fabric/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.crt
}

type CertificateAuthorityRegistrar struct {
	EnrollId     string `yaml:"enrollId"`
	EnrollSecret string `yaml:"enrollSecret"`
}
