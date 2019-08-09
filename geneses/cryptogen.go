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
	str "github.com/ennoo/rivet/utils/string"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
)

// CryptoGen 联盟证书文件对象
type CryptoGen struct {
	OrdererOrgs []*Order `yaml:"OrdererOrgs" json:"orderOrgArr"`
	PeerOrgs    []*Peer  `yaml:"PeerOrgs" json:"peerOrgArr"`
}

// Order 通道配置
type Order struct {
	Name   string      `yaml:"Name" json:"name"`
	Domain string      `yaml:"Domain" json:"domain"`
	CA     *CA         `yaml:"CA" json:"ca"`
	Specs  []*Hostname `yaml:"Specs" json:"specs"`
}

// Hostname 主机配置
type Hostname struct {
	Hostname string `yaml:"Hostname" json:"hostname"`
}

// Peer 节点配置
type Peer struct {
	Name          string    `yaml:"Name" json:"name"`
	Domain        string    `yaml:"Domain" json:"domain"`
	CA            *CA       `yaml:"CA" json:"ca"`
	Template      *Template `yaml:"Template" json:"template"`
	Users         *Users    `yaml:"Users" json:"users"`
	EnableNodeOUs bool      `yaml:"EnableNodeOUs" json:"enableNodeOUs"`
}

// Template 模板属性配置
type Template struct {
	Count int32 `yaml:"Count" json:"count"`
}

// Users 用户属性配置
type Users struct {
	Count int32 `yaml:"Count" json:"count"`
}

// CA 证书配置
type CA struct {
	Country  string `yaml:"Country" json:"country"`
	Province string `yaml:"Province" json:"province"`
	Locality string `yaml:"Locality" json:"locality"`
}

func generateCryptoGenYml(leagueComment string, orderCount, peerCount, templateCount, userCount int32) ([]byte, error) {
	var (
		cryptoGen *CryptoGen
		err       error
		data      []byte
	)
	if cryptoGen, err = generateCryptoGen(leagueComment, orderCount, peerCount, templateCount, userCount); nil != err {
		return nil, err
	}
	if data, err = yaml.Marshal(&cryptoGen); err != nil {
		return nil, err
	}
	return data, nil
}

func generateCryptoGenCustomYml(ordererOrgs []*Order, peerOrgs []*Peer) ([]byte, error) {
	var (
		cryptoGen *CryptoGen
		err       error
		data      []byte
	)
	if cryptoGen, err = generateCryptoGenCustom(ordererOrgs, peerOrgs); nil != err {
		return nil, err
	}
	if data, err = yaml.Marshal(&cryptoGen); err != nil {
		return nil, err
	}
	return data, nil
}

// generateCryptoGen 生成联盟证书文件对象
func generateCryptoGen(ledgerName string, orderCount, peerCount, templateCount, userCount int32) (*CryptoGen, error) {
	if str.IsEmpty(ledgerName) || orderCount <= 0 || peerCount <= 0 {
		return nil, errors.New("crypto params exception")
	}
	ca := &CA{
		Country:  "CN",
		Locality: "Beijing",
		Province: "Beijing",
	}
	hostnames := make([]*Hostname, orderCount)
	for index := range hostnames {
		hostnames[index] = &Hostname{strings.Join([]string{OrderPrefix, strconv.Itoa(index)}, "")}
	}
	ordererOrgs := []*Order{
		{
			Name:   "Orderer",
			Domain: ledgerName,
			CA:     ca,
			Specs:  hostnames,
		},
	}
	peerOrgs := make([]*Peer, peerCount)
	for index := range peerOrgs {
		peerOrgs[index] = &Peer{
			Name:   strings.Join([]string{"Org", strconv.Itoa(index + 1)}, ""),
			Domain: strings.Join([]string{ledgerName, "-org", strconv.Itoa(index + 1)}, ""),
			CA:     ca,
			Template: &Template{
				Count: templateCount,
			},
			Users: &Users{
				Count: userCount,
			},
			EnableNodeOUs: true,
		}
	}
	return generateCryptoGenCustom(ordererOrgs, peerOrgs)
}

// generateCryptoGenCustom 生成自定义联盟证书文件对象
func generateCryptoGenCustom(ordererOrgs []*Order, peerOrgs []*Peer) (*CryptoGen, error) {
	return &CryptoGen{
		OrdererOrgs: ordererOrgs,
		PeerOrgs:    peerOrgs,
	}, nil
}
