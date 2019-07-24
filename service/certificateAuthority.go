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
 *
 */

package service

import (
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
)

type CertificateAuthority struct {
	ConfigID                string `json:"configID"` // ConfigID 配置唯一ID
	CertName                string `json:"certName"`
	URL                     string `json:"url"`
	TLSCACertPath           string `json:"tlsCACertPath"`
	TLSCACertClientKeyPath  string `json:"tlsCACertClientKeyPath"`
	TLSCACertClientCertPath string `json:"tlsCACertClientCertPath"`
	CAName                  string `json:"caName"`
	EnrollId                string `json:"enrollId"`
	EnrollSecret            string `json:"enrollSecret"`
}

func (c *CertificateAuthority) Trans2pb() *pb.ReqCertificateAuthority {
	return &pb.ReqCertificateAuthority{
		ConfigID:                c.ConfigID,
		CertName:                c.CertName,
		Url:                     c.URL,
		TlsCACertPath:           c.TLSCACertPath,
		TlsCACertClientKeyPath:  c.TLSCACertClientKeyPath,
		TlsCACertClientCertPath: c.TLSCACertClientCertPath,
		CaName:                  c.CAName,
		EnrollId:                c.EnrollId,
		EnrollSecret:            c.EnrollSecret,
	}
}

type CertificateAuthoritySelf struct {
	ConfigID     string `json:"configID"` // ConfigID 配置唯一ID
	LeagueName   string `json:"leagueName"`
	CertName     string `json:"certName"`
	URL          string `json:"url"`
	CAName       string `json:"caName"`
	EnrollId     string `json:"enrollId"`
	EnrollSecret string `json:"enrollSecret"`
}

func (c *CertificateAuthoritySelf) Trans2pb() *pb.ReqCertificateAuthoritySelf {
	return &pb.ReqCertificateAuthoritySelf{
		ConfigID:     c.ConfigID,
		LeagueName:   c.LeagueName,
		CertName:     c.CertName,
		Url:          c.URL,
		CaName:       c.CAName,
		EnrollId:     c.EnrollId,
		EnrollSecret: c.EnrollSecret,
	}
}
