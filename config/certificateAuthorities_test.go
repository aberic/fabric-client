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

func TestNewCertificateAuthorities(t *testing.T) {
	certificateAuthorities := TGetCertificateAuthorities()

	certificateAuthoritiesData, err := yaml.Marshal(&certificateAuthorities)
	if err != nil {
		log.Self.Error("channels", log.Error(err))
	}
	fmt.Printf("--- kfk dump:\n%s\n\n", string(certificateAuthoritiesData))
}

func TGetCertificateAuthorities() map[string]*CertificateAuthority {
	return map[string]*CertificateAuthority{
		"ca.org1.example.com": {
			URL: "https://ca.org1.example.com:7054",
			TLSCACerts: &CertificateAuthorityTLSCACerts{
				Path: "peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem",
				Client: &CertificateAuthorityTLSCACertsClient{
					Key: &CertificateAuthorityTLSCACertsClientKey{
						Path: "peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.key",
					},
					Cert: &CertificateAuthorityTLSCACertsClientCert{
						Path: "peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/client.crt",
					},
				},
			},
			Registrar: &CertificateAuthorityRegistrar{
				EnrollId:     "admin",
				EnrollSecret: "adminpw",
			},
			CAName: "ca.org1.example.com",
		},
		"ca.org2.example.com": {
			URL: "https://ca.org2.example.com:7054",
			TLSCACerts: &CertificateAuthorityTLSCACerts{
				Path: "peerOrganizations/org1.example.com/tlsca/tlsca.org2.example.com-cert.pem",
				Client: &CertificateAuthorityTLSCACertsClient{
					Key: &CertificateAuthorityTLSCACertsClientKey{
						Path: "peerOrganizations/org1.example.com/users/User1@org2.example.com/tls/client.key",
					},
					Cert: &CertificateAuthorityTLSCACertsClientCert{
						Path: "peerOrganizations/org1.example.com/users/User1@org2.example.com/tls/client.crt",
					},
				},
			},
			Registrar: &CertificateAuthorityRegistrar{
				EnrollId:     "admin",
				EnrollSecret: "adminpw",
			},
			CAName: "ca.org2.example.com",
		},
	}
}
