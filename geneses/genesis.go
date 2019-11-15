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

package geneses

import (
	"github.com/aberic/fabric-client/grpc/proto/generate"
	"github.com/aberic/gnomon"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource/genesisconfig"
	"strings"
	"time"
)

type Genesis struct {
	Info               *generate.ReqGenesis
	orderOrganizations []*genesisconfig.Organization
	peerOrganizations  []*genesisconfig.Organization
	allOrganizations   []*genesisconfig.Organization
}

func (g *Genesis) Init() {
	g.orderOrganizations, g.peerOrganizations, g.allOrganizations = g.organizations(g.Info.Orgs)
}

func (g *Genesis) CreateGenesisBlock(consortium string) error {
	data, err := resource.CreateGenesisBlock(g.genesisBlockConfigProfile(consortium), consortium)
	if nil != err {
		return err
	}
	if _, err = gnomon.File().Append(GenesisBlockFilePath(g.Info.League.Domain), data, true); nil != err {
		return err
	}
	return nil
}

func (g *Genesis) CreateChannelCreateTx(consortium, channelID string) error {
	data, err := resource.CreateChannelCreateTx(g.genesisChannelTxConfigProfile(consortium), nil, channelID)
	if nil != err {
		return err
	}
	if _, err = gnomon.File().Append(ChannelTXFilePath(g.Info.League.Domain, channelID), data, true); nil != err {
		return err
	}
	return nil
}

func (g *Genesis) orgPolicies(mspID string) map[string]*genesisconfig.Policy {
	return map[string]*genesisconfig.Policy{
		"Readers": {
			Type: "Signature",
			Rule: strings.Join([]string{"OR('", mspID, ".member')"}, ""),
		},
		"Writers": {
			Type: "Signature",
			Rule: strings.Join([]string{"OR('", mspID, ".member')"}, ""),
		},
		"Admins": {
			Type: "Signature",
			Rule: strings.Join([]string{"OR('", mspID, ".admin')"}, ""),
		},
		"Endorsement": {
			Type: "Signature",
			Rule: strings.Join([]string{"OR('", mspID, ".member')"}, ""),
		},
	}
}

func (g *Genesis) organizations(orgs []*generate.OrgInBlock) (orders, peers, all []*genesisconfig.Organization) {
	for _, org := range orgs {
		var (
			mspDir string
			mspID  = MspID(org.Name)
		)
		organization := &genesisconfig.Organization{
			Name:           org.Name,
			SkipAsForeign:  false,
			ID:             mspID,
			MSPType:        "bccsp",
			Policies:       g.orgPolicies(mspID),
			AdminPrincipal: "Role.ADMIN",
		}
		switch org.Type {
		default:
			return
		case generate.OrgType_Peer:
			var anchorPeers []*genesisconfig.AnchorPeer
			for _, peer := range org.AnchorPeers {
				anchorPeers = append(anchorPeers, &genesisconfig.AnchorPeer{Host: peer.Host, Port: int(peer.Port)})
			}
			organization.AnchorPeers = anchorPeers
			mspDir = CryptoOrgMspPath(g.Info.League.Domain, org.Domain, org.Name, true)
			organization.MSPDir = mspDir
			peers = append(peers, organization)
		case generate.OrgType_Order:
			mspDir = CryptoOrgMspPath(g.Info.League.Domain, org.Domain, org.Name, false)
			organization.MSPDir = mspDir
			orders = append(orders, organization)
		}
		all = append(all, organization)
	}
	return
}

func (g *Genesis) applicationCapabilities() map[string]bool {
	return map[string]bool{
		"V1_1": false,
		"V1_2": false,
		"V1_3": true,
	}
}

func (g *Genesis) applications() *genesisconfig.Application {
	//rule := strings.Join([]string{"OR('", adminOrgMspID, ".admin')"}, "")
	return &genesisconfig.Application{
		Organizations: g.peerOrganizations,
		Capabilities:  g.applicationCapabilities(),
		Policies: map[string]*genesisconfig.Policy{
			"LifecycleEndorsement": {
				Rule: "MAJORITY Endorsement",
				Type: "ImplicitMeta",
			},
			"Endorsement": {
				Rule: "MAJORITY Endorsement",
				Type: "ImplicitMeta",
			},
			"Readers": {
				Rule: "ANY Readers",
				Type: "ImplicitMeta",
			},
			"Writers": {
				Rule: "ANY Writers",
				Type: "ImplicitMeta",
			},
			"Admins": {
				Rule: "MAJORITY Admins",
				Type: "ImplicitMeta",
			},
			//"ChannelCreate": {
			//	Type: "Signature",
			//	Rule: rule,
			//},
		},
		ACLs: map[string]string{
			"_lifecycle/CommitChaincodeDefinition": "/Channel/Application/Writers",
			"_lifecycle/QueryChaincodeDefinition":  "/Channel/Application/Readers",
			"_lifecycle/QueryNamespaceDefinitions": "/Channel/Application/Readers",
			"lscc/ChaincodeExists":                 "/Channel/Application/Readers",
			"lscc/GetDeploymentSpec":               "/Channel/Application/Readers",
			"lscc/GetChaincodeData":                "/Channel/Application/Readers",
			"lscc/GetInstantiatedChaincodes":       "/Channel/Application/Readers",
			"qscc/GetChainInfo":                    "/Channel/Application/Readers",
			"qscc/GetBlockByNumber":                "/Channel/Application/Readers",
			"qscc/GetBlockByHash":                  "/Channel/Application/Readers",
			"qscc/GetTransactionByID":              "/Channel/Application/Readers",
			"qscc/GetBlockByTxID":                  "/Channel/Application/Readers",
			"cscc/GetConfigBlock":                  "/Channel/Application/Readers",
			"cscc/GetConfigTree":                   "/Channel/Application/Readers",
			"cscc/SimulateConfigTreeUpdate":        "/Channel/Application/Readers",
			"peer/Propose":                         "/Channel/Application/Writers",
			"peer/ChaincodeToChaincode":            "/Channel/Application/Readers",
			"event/Block":                          "/Channel/Application/Readers",
			"event/FilteredBlock":                  "/Channel/Application/Readers",
		},
	}
}

func (g *Genesis) ordererCapabilities() map[string]bool {
	return map[string]bool{
		"V1_1": true,
	}
}

func (g *Genesis) orderer() *genesisconfig.Orderer {
	return &genesisconfig.Orderer{
		OrdererType:  "kafka",
		Addresses:    g.Info.League.Addresses, // []string{"orderer.example.org:7050"}
		BatchTimeout: time.Duration(time.Duration(g.Info.League.BatchTimeout) * time.Second),
		BatchSize: genesisconfig.BatchSize{
			MaxMessageCount:   g.Info.League.BatchSize.MaxMessageCount,   // 500
			AbsoluteMaxBytes:  g.Info.League.BatchSize.AbsoluteMaxBytes,  //10 * 1024 * 1024
			PreferredMaxBytes: g.Info.League.BatchSize.PreferredMaxBytes, //2 * 1024 * 1024
		},
		Kafka: genesisconfig.Kafka{
			Brokers: g.Info.League.Kafka.Brokers, // []string{"kafka1.league01:9092", "kafka2.league01:9092"}
		},
		Organizations: g.orderOrganizations,
		MaxChannels:   g.Info.League.MaxChannels, // 1000
		// Policies defines the set of policies at this level of the config tree
		// For Orderer policies, their canonical path is
		// /Channel/Orderer/<PolicyName>
		Policies: map[string]*genesisconfig.Policy{
			"Readers": {
				Type: "ImplicitMeta",
				Rule: "ANY Readers",
			},
			"Writers": {
				Type: "ImplicitMeta",
				Rule: "ANY Writers",
			},
			"Admins": {
				Type: "ImplicitMeta",
				Rule: "MAJORITY Admins",
			},
			"BlockValidation": {
				Type: "ImplicitMeta",
				Rule: "ANY Writers",
			},
		},
		Capabilities: g.ordererCapabilities(),
	}
}

func (g *Genesis) channelDefaults() map[string]*genesisconfig.Policy {
	// Policies defines the set of policies at this level of the config tree
	// For Channel policies, their canonical path is
	// /Channel/<PolicyName>
	policies := map[string]*genesisconfig.Policy{
		"Admins": {
			Type: "ImplicitMeta",
			Rule: "MAJORITY Admins",
		},
		"Readers": {
			Type: "ImplicitMeta",
			Rule: "ANY Readers",
		},
		"Writers": {
			Type: "ImplicitMeta",
			Rule: "ANY Writers",
		},
	}
	return policies
}

func (g *Genesis) genesisBlockConfigProfile(consortium string) *genesisconfig.Profile {
	profile := &genesisconfig.Profile{
		Orderer: g.orderer(),
		Consortiums: map[string]*genesisconfig.Consortium{
			consortium: {Organizations: g.peerOrganizations},
		},
		Policies: g.channelDefaults(),
	}
	return profile
}

func (g *Genesis) genesisChannelTxConfigProfile(consortium string) *genesisconfig.Profile {
	profile := &genesisconfig.Profile{
		Consortium:  consortium,
		Application: g.applications(),
		Policies:    g.channelDefaults(),
	}
	return profile
}
