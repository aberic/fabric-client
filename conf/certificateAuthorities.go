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
)

type CertificateAuthority struct {
	URL        string                          `yaml:"url"`    // URL https://ca.org1.example.com:7054
	CAName     string                          `yaml:"caName"` // CAName 可选参数，name of the CA
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

func (c *CertificateAuthority) set(leagueDomain, caName string, in *chain.CertificateAuthority) error {
	if gnomon.String().IsNotEmpty(in.Url) {
		c.URL = in.Url
	} else {
		return errors.New("url can't be empty")
	}
	if gnomon.String().IsNotEmpty(in.CaName) {
		c.CAName = in.CaName
	} else {
		c.CAName = caName
	}
	if len(in.TlsCACerts.Cert) > 0 {
		tlsCACertFilePath := geneses.CertificateAuthorityFilePath(leagueDomain, caName)
		if _, err := gnomon.File().Append(tlsCACertFilePath, in.TlsCACerts.Cert, true); nil != err {
			return err
		}
		c.TLSCACerts.Path = tlsCACertFilePath
	} else {
		return errors.New("tls ca cert can't be empty")
	}
	if nil != in.TlsCACerts.Client {
		if nil != in.TlsCACerts.Client.Key && len(in.TlsCACerts.Client.Key.Key) > 0 {
			tlsCACertClientKeyFilePath := geneses.CertificateAuthorityClientKeyFilePath(leagueDomain, caName)
			if _, err := gnomon.File().Append(tlsCACertClientKeyFilePath, in.TlsCACerts.Client.Key.Key, true); nil != err {
				return err
			}
			c.TLSCACerts.Client.Key.Path = tlsCACertClientKeyFilePath
		}
		if nil != in.TlsCACerts.Client.Cert && len(in.TlsCACerts.Client.Cert.Cert) > 0 {
			tlsCACertClientCertFilePath := geneses.CertificateAuthorityClientCertFilePath(leagueDomain, caName)
			if _, err := gnomon.File().Append(tlsCACertClientCertFilePath, in.TlsCACerts.Client.Cert.Cert, true); nil != err {
				return err
			}
			c.TLSCACerts.Client.Cert.Path = tlsCACertClientCertFilePath
		}
	}
	return nil
}
