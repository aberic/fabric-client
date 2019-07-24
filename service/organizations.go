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

type OrganizationsOrder struct {
	ConfigID   string            `json:"configID"` // ConfigID 配置唯一ID
	MspID      string            `json:"mspID"`
	CryptoPath string            `json:"cryptoPath"`
	Users      map[string]string `json:"users"`
}

func (o *OrganizationsOrder) Trans2pb() *pb.ReqOrganizationsOrder {
	return &pb.ReqOrganizationsOrder{
		ConfigID:   o.ConfigID,
		MspID:      o.MspID,
		CryptoPath: o.CryptoPath,
		Users:      o.Users,
	}
}

type OrganizationsOrderSelf struct {
	ConfigID   string `json:"configID"` // ConfigID 配置唯一ID
	LeagueName string `json:"leagueName"`
}

func (o *OrganizationsOrderSelf) Trans2pb() *pb.ReqOrganizationsOrderSelf {
	return &pb.ReqOrganizationsOrderSelf{
		ConfigID:   o.ConfigID,
		LeagueName: o.LeagueName,
	}
}

type OrganizationsOrg struct {
	ConfigID               string            `json:"configID"` // ConfigID 配置唯一ID
	OrgName                string            `json:"orgName"`
	MspID                  string            `json:"mspID"`
	CryptoPath             string            `json:"cryptoPath"`
	Users                  map[string]string `json:"users"`
	Peers                  []string          `json:"peers"`
	CertificateAuthorities []string          `json:"certificateAuthorities"`
}

func (o *OrganizationsOrg) Trans2pb() *pb.ReqOrganizationsOrg {
	return &pb.ReqOrganizationsOrg{
		ConfigID:               o.ConfigID,
		OrgName:                o.OrgName,
		MspID:                  o.MspID,
		CryptoPath:             o.CryptoPath,
		Users:                  o.Users,
		Peers:                  o.Peers,
		CertificateAuthorities: o.CertificateAuthorities,
	}
}

type OrganizationsOrgSelf struct {
	ConfigID               string   `json:"configID"` // ConfigID 配置唯一ID
	LeagueName             string   `json:"leagueName"`
	Peers                  []string `json:"peers"`
	CertificateAuthorities []string `json:"certificateAuthorities"`
}

func (o *OrganizationsOrgSelf) Trans2pb() *pb.ReqOrganizationsOrgSelf {
	return &pb.ReqOrganizationsOrgSelf{
		ConfigID:               o.ConfigID,
		LeagueName:             o.LeagueName,
		Peers:                  o.Peers,
		CertificateAuthorities: o.CertificateAuthorities,
	}
}
