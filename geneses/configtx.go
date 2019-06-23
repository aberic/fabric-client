/*
 * Copyright (c) 2019.. ENNOO - All Rights Reserved.
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
	"errors"
	"github.com/ennoo/rivet/utils/string"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
)

// ConfigTX 区块链生成创世区块、通道创世区块等所用配置对象
type ConfigTX struct {
	Organizations []Organization   `yaml:"Organizations" json:"organizations"`
	Application   *Application     `yaml:"Application" json:"application"`
	Capabilities  *CapabilitiesAll `yaml:"Capabilities" json:"capabilities"`
	Channel       *Channel         `yaml:"Channel" json:"channel"`
	Orderer       *Orderer         `yaml:"Orderer" json:"order"`
	Profiles      *Profiles        `yaml:"Profiles" json:"profiles"`
}

// Organization 定义了不同的组织标识，这些标识将在稍后的配置中引用
type Organization interface {
	mspDir()
}

// OrdererOrg 定义了一个 Orderer Org MSP
type OrdererOrg struct {
	Name     string    `yaml:"Name" json:"name"`
	ID       string    `yaml:"ID" json:"id"`
	MSPDir   string    `yaml:"MSPDir" json:"mspDir"`
	Policies *Policies `yaml:"Policies" json:"policies"`
}

// Org 定义了一个 Member Org MSP
type Org struct {
	Name        string        `yaml:"Name" json:"name"`
	ID          string        `yaml:"ID" json:"id"`
	MSPDir      string        `yaml:"MSPDir" json:"mspDir"`
	Policies    *Policies     `yaml:"Policies" json:"policies"`
	AnchorPeers []*AnchorPeer `yaml:"AnchorPeers" json:"anchorPeers"`
}

// AnchorPeer 锚节点
type AnchorPeer struct {
	Host string `yaml:"Host" json:"host"`
	Port int    `yaml:"Port" json:"port"`
}

func (o *OrdererOrg) mspDir() {}

func (o *Org) mspDir() {}

// Application 定义要编码到配置事务或应用程序相关参数的genesis块中的值
type Application struct {
	Organizations []Organization `yaml:"Organizations" json:"organizations"`
	Capabilities  *Capabilities  `yaml:"Capabilities" json:"capabilities"`
	Policies      *Policies      `yaml:"Policies" json:"policies"`
}

// Capabilities 应用程序功能只应用于对等网络，并且可以安全地与以前的版本订购者一起使用。将功能的值设置为true以满足需要。
type Capabilities struct {
	// V11 应用程序支持fabric v1.1新的非向后兼容特性和补丁(注意，如果设置了较晚版本的功能，则不需要设置此功能)
	V11 bool `yaml:"V1_1" json:"v11"`
	// V12 应用程序支持fabric v1.2新的非向后兼容特性和补丁(注意，如果设置了较晚版本的功能，则不需要设置此功能)
	V12 bool `yaml:"V1_2" json:"v12"`
	// V13 应用程序支持fabric v1.3新的非向后兼容特性和补丁
	V13 bool `yaml:"V1_3" json:"v13"`
}

// Policies 在组织策略的配置树的这个级别定义了一组策略，它们的规范路径通常是/Channel/<Application|Orderer>/<OrgName>/<PolicyName>
type Policies struct {
	Admins  *Policy `yaml:"Admins" json:"admins"`
	Readers *Policy `yaml:"Readers" json:"readers"`
	Writers *Policy `yaml:"Writers" json:"writers"`
}

// Policy Policies 等策略详情
type Policy struct {
	Rule string `yaml:"Rule" json:"rule"`
	Type string `yaml:"Type" json:"type"`
}

// CapabilitiesAll 定义fabric network的功能。这是v1.1.0的一个新概念，不应该在v1.0的混合网络中使用。
//
// x个对等点和顺序。功能定义了fabric二进制文件中必须提供的特性，以便该二进制文件安全地参与fabric网络。
//
// 例如，如果添加了新的MSP类型，较新的二进制文件可能会识别并验证来自该类型的签名，而没有此支持的较老的二进制文件将无法验证这些事务。
//
// 这可能导致不同版本的fabric二进制文件具有不同的世界状态。
//
// 相反，为通道定义一个功能会通知那些没有这个功能的二进制文件，它们必须停止处理事务，直到它们被升级。
//
// v1.0。如果定义了任何功能(包括关闭所有功能的映射)，则使用v1.0。x peer会故意崩溃。
type CapabilitiesAll struct {
	Application *Capabilities `yaml:"Application" json:"application"`
	Channel     *V13          `yaml:"Channel" json:"channel"`
	Orderer     *V11          `yaml:"Orderer" json:"order"`
}

// V13 通道功能同时适用于订购方和对等方，并且必须得到双方的支持。将功能的值设置为true以满足需要
type V13 struct {
	V13 bool `yaml:"V1_3" json:"v13"`
}

// V11 Orderer功能只适用于Orderer，并且可以安全地与以前的版本对等点一起使用。将功能的值设置为true以满足需要。
type V11 struct {
	V11 bool `yaml:"V1_1" json:"v11"`
}

// Channel 通道配置
type Channel struct {
	Capabilities *V13      `yaml:"Capabilities" json:"capabilities"`
	Policies     *Policies `yaml:"Policies" json:"policies"`
}

// Orderer 排序服务配置
type Orderer struct {
	Addresses []string   `yaml:"Addresses" json:"addresses"`
	BatchSize *BatchSize `yaml:"BatchSize" json:"batchSize"`
	// BatchTimeout 批处理超时:在创建批处理之前等待的时间
	BatchTimeout  string           `yaml:"BatchTimeout" json:"batchTimeout"`
	Capabilities  *V11             `yaml:"Capabilities" json:"capabilities"`
	Kafka         *Kafka           `yaml:"Kafka" json:"kafka"`
	OrdererType   string           `yaml:"OrdererType" json:"ordererType"`
	Organizations []Organization   `yaml:"Organizations" json:"organizations"`
	Policies      *PoliciesOrderer `yaml:"Policies" json:"policies"`
}

// BatchSize 打包策略配置
type BatchSize struct {
	// AbsoluteMaxBytes 绝对最大字节数:批处理中允许序列化消息的绝对最大字节数。
	AbsoluteMaxBytes string `yaml:"AbsoluteMaxBytes" json:"absoluteMaxBytes"`
	// MaxMessageCount 最大消息数:批处理中允许的最大消息数
	MaxMessageCount int `yaml:"MaxMessageCount" json:"maxMessageCount"`
	// PreferredMaxBytes 首选最大字节数:批处理中序列化消息所允许的首选最大字节数。大于首选最大字节的消息将导致批处理大于首选最大字节。
	PreferredMaxBytes string `yaml:"PreferredMaxBytes" json:"preferredMaxBytes"`
}

// Kafka kfk配置
type Kafka struct {
	Brokers []string `yaml:"Brokers" json:"brokers"`
}

// PoliciesOrderer 排序服务组策略配置
type PoliciesOrderer struct {
	Admins          *Policy `yaml:"Admins" json:"admins"`
	Readers         *Policy `yaml:"Readers" json:"readers"`
	Writers         *Policy `yaml:"Writers" json:"writers"`
	BlockValidation *Policy `yaml:"BlockValidation" json:"blockValidation"`
}

// Profiles 配置文件输出策略
type Profiles struct {
	HBaaSChannel      *HBaaSChannel      `yaml:"HBaaSChannel" json:"hbaasChannel"`
	HBaaSOrderGenesis *HBaaSOrderGenesis `yaml:"HBaaSOrderGenesis" json:"hbaasOrderGenesis"`
}

// HBaaSChannel 配置文件通道创世区块输出策略
type HBaaSChannel struct {
	Application *Application `yaml:"Application" json:"application"`
	Consortium  string       `yaml:"Consortium" json:"consortium"`
}

// HBaaSOrderGenesis 配置文件联盟创世区块输出策略
type HBaaSOrderGenesis struct {
	Capabilities *V13         `yaml:"Capabilities" json:"capabilities"`
	Consortiums  *Consortiums `yaml:"Consortiums" json:"consortiums"`
	Orderer      *Orderer     `yaml:"Orderer" json:"orderer"`
	Policies     *Policies    `yaml:"Policies" json:"policies"`
}

// Consortiums 配置指定组织集合对象
type Consortiums struct {
	HBaaSConsortium *HBaaSConsortium `yaml:"HBaaSConsortium" json:"hbaasConsortium"`
}

// HBaaSConsortium 配置指定组织集合
type HBaaSConsortium struct {
	Organizations []Organization `yaml:"Organizations" json:"organizations"`
}

// generateConfigTXYml 生成模板配置Yml文件
func generateConfigTXYml(leagueComment string, orderCount, peerCount, batchTimeout, maxMessageCount int) ([]byte, error) {
	var (
		configTx *ConfigTX
		err      error
		data     []byte
	)
	if configTx, err = generateConfigTX(leagueComment, orderCount, peerCount, batchTimeout, maxMessageCount); nil != err {
		return nil, err
	}
	if data, err = yaml.Marshal(&configTx); err != nil {
		return nil, err
	}
	return data, nil
}

// generateConfigTXYmlCustom 生成模板自定义配置Yml文件
func generateConfigTXCustomYml(organizations []Organization, application *Application, capabilities *CapabilitiesAll,
	channel *Channel, orderer *Orderer, profiles *Profiles) ([]byte, error) {
	var (
		data []byte
		err  error
	)
	configTx := generateConfigTXCustom(organizations, application, capabilities, channel, orderer, profiles)
	if data, err = yaml.Marshal(&configTx); err != nil {
		return nil, err
	}
	return data, nil
}

// generateConfigTXCustom 生成自定义配置对象
func generateConfigTXCustom(organizations []Organization, application *Application, capabilities *CapabilitiesAll,
	channel *Channel, orderer *Orderer, profiles *Profiles) *ConfigTX {
	return &ConfigTX{
		Organizations: organizations,
		Application:   application,
		Capabilities:  capabilities,
		Channel:       channel,
		Orderer:       orderer,
		Profiles:      profiles,
	}
}

// generateConfigTX 生成模板配置对象
func generateConfigTX(leagueComment string, orderCount, peerCount, batchTimeout, maxMessageCount int) (*ConfigTX, error) {
	if str.IsEmpty(leagueComment) || orderCount <= 0 || peerCount <= 0 || batchTimeout <= 0 || maxMessageCount <= 0 {
		return nil, errors.New("config params exception")
	}
	cryptoConfigPath := CryptoConfigPath(leagueComment)
	organizations := getOrganizations(leagueComment, cryptoConfigPath, peerCount)
	application := getApplication(organizations)
	capabilities := getCapabilities(application)
	channel := getChannel()
	orderer := getOrderer(organizations, leagueComment, orderCount, batchTimeout, maxMessageCount)
	profiles := getProfiles(organizations, application, orderer)
	return generateConfigTXCustom(organizations, application, capabilities, channel, orderer, profiles), nil
}

func getOrganizations(leagueComment, cryptoGenFilesPath string, peerCount int) []Organization {
	organizations := make([]Organization, peerCount+1)
	for index := range organizations {
		if index == 0 {
			organizations[0] = &OrdererOrg{
				Name:   "OrdererOrg",
				ID:     "OrdererMSP",
				MSPDir: strings.Join([]string{cryptoGenFilesPath, "ordererOrganizations", leagueComment, "msp"}, "/"),
				Policies: &Policies{
					Admins: &Policy{
						Rule: "OR('OrdererMSP.admin')",
						Type: "Signature",
					},
					Readers: &Policy{
						Rule: "OR('OrdererMSP.member')",
						Type: "Signature",
					},
					Writers: &Policy{
						Rule: "OR('OrdererMSP.member')",
						Type: "Signature",
					},
				},
			}
		} else {
			mspName := strings.Join([]string{"Org", strconv.Itoa(index), "MSP"}, "")
			peerHost := strings.Join([]string{"peer0.", leagueComment, "-org", strconv.Itoa(index)}, "")
			peerLeagueComment := strings.Join([]string{leagueComment, "-org", strconv.Itoa(index)}, "")
			organizations[index] = &Org{
				Name:     mspName,
				ID:       mspName,
				MSPDir:   strings.Join([]string{cryptoGenFilesPath, "peerOrganizations", peerLeagueComment, "msp"}, "/"),
				Policies: getPoliciesForOrg(index),
				AnchorPeers: []*AnchorPeer{
					{
						Host: peerHost,
						Port: 7051,
					},
				},
			}
		}
	}
	return organizations
}

func getPoliciesForOrg(index int) *Policies {
	is := strconv.Itoa(index)
	return &Policies{
		Admins: &Policy{
			Rule: strings.Join([]string{"OR('Org", is, "MSP.admin')"}, ""),
			Type: "Signature",
		},
		Readers: &Policy{
			Rule: strings.Join([]string{"OR('Org", is, "MSP.admin', 'Org", is, "MSP.peer', 'Org", is, "MSP.client')"}, ""),
			Type: "Signature",
		},
		Writers: &Policy{
			Rule: strings.Join([]string{"OR('Org", is, "MSP.admin', 'Org", is, "MSP.client')"}, ""),
			Type: "Signature",
		},
	}
}

func getPolicies() *Policies {
	return &Policies{
		Admins: &Policy{
			Rule: "MAJORITY Admins",
			Type: "ImplicitMeta",
		},
		Readers: &Policy{
			Rule: "ANY Readers",
			Type: "ImplicitMeta",
		},
		Writers: &Policy{
			Rule: "ANY Writers",
			Type: "ImplicitMeta",
		},
	}
}

func getApplication(organizations []Organization) *Application {
	application := &Application{
		Organizations: organizations[1:],
		Capabilities: &Capabilities{
			V11: false,
			V12: false,
			V13: true,
		},
		Policies: getPolicies(),
	}
	return application
}

func getCapabilities(application *Application) *CapabilitiesAll {
	capabilities := &CapabilitiesAll{
		Application: application.Capabilities,
		Channel:     getV13(),
		Orderer:     getV11(),
	}
	return capabilities
}

func getV13() *V13 {
	return &V13{
		V13: true,
	}
}

func getV11() *V11 {
	return &V11{
		V11: true,
	}
}

func getChannel() *Channel {
	return &Channel{
		Capabilities: getV13(),
		Policies:     getPolicies(),
	}
}

func getOrderer(organizations []Organization, leagueComment string, orderCount, batchTimeout, maxMessageCount int) *Orderer {
	addresses := make([]string, orderCount)
	for index := range addresses {
		addresses[index] = strings.Join([]string{OrderPrefix, strconv.Itoa(index), ".", leagueComment, ":7050"}, "")
	}
	return &Orderer{
		Addresses: addresses,
		BatchSize: &BatchSize{
			AbsoluteMaxBytes:  "98 MB",
			MaxMessageCount:   maxMessageCount,
			PreferredMaxBytes: "512 KB",
		},
		BatchTimeout: strings.Join([]string{strconv.Itoa(batchTimeout), "s"}, ""),
		Capabilities: getV11(),
		Kafka: &Kafka{
			Brokers: getKafkaBrokers(leagueComment),
		},
		OrdererType:   "kafka",
		Organizations: organizations[0:1],
		Policies: &PoliciesOrderer{
			Admins: &Policy{
				Rule: "MAJORITY Admins",
				Type: "ImplicitMeta",
			},
			Readers: &Policy{
				Rule: "ANY Readers",
				Type: "ImplicitMeta",
			},
			Writers: &Policy{
				Rule: "ANY Writers",
				Type: "ImplicitMeta",
			},
			BlockValidation: &Policy{
				Rule: "ANY Writers",
				Type: "ImplicitMeta",
			},
		},
	}
}

func getKafkaBrokers(leagueComment string) []string {
	brokers := make([]string, 7)
	for i := 0; i < 7; i++ {
		brokers[i] = strings.Join([]string{"kfk", strconv.Itoa(i + 1), ".", leagueComment, ":9092"}, "")
	}
	return brokers
}

func getProfiles(organizations []Organization, application *Application, orderer *Orderer) *Profiles {
	profiles := &Profiles{
		HBaaSChannel: &HBaaSChannel{
			Application: application,
			Consortium:  "HBaaSConsortium",
		},
		HBaaSOrderGenesis: &HBaaSOrderGenesis{
			Capabilities: getV13(),
			Consortiums: &Consortiums{
				HBaaSConsortium: &HBaaSConsortium{
					Organizations: organizations[1:],
				},
			},
			Orderer:  orderer,
			Policies: getPolicies(),
		},
	}
	return profiles
}
