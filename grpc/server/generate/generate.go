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

package generate

import (
	"context"
	"errors"
	"github.com/aberic/fabric-client/geneses"
	"github.com/aberic/fabric-client/grpc/proto/generate"
)

type CreationServer struct{}

func (cs *CreationServer) GenerateCrypto(ctx context.Context, in *generate.ReqKeyConfig) (*generate.RespKeyConfig, error) {
	pemConfig := &geneses.PemConfig{KeyConfig: in}
	kcr := pemConfig.GenerateCrypto()
	if kcr.Code == generate.Code_Success {
		return kcr, nil
	}
	return kcr, errors.New(kcr.ErrMsg)
}

func (cs *CreationServer) GenerateLeague(ctx context.Context, in *generate.ReqCreateLeague) (*generate.RespCreateLeague, error) {
	gc := &geneses.GenerateConfig{}
	if err := gc.CreateLeague(in); nil != err {
		return &generate.RespCreateLeague{Code: generate.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &generate.RespCreateLeague{Code: generate.Code_Success}, nil
}

func (cs *CreationServer) GenerateCsr(ctx context.Context, in *generate.ReqCreateCsr) (*generate.RespCreateCsr, error) {
	gc := &geneses.GenerateConfig{}
	if err := gc.CreateCsr(in); nil != err {
		return &generate.RespCreateCsr{Code: generate.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &generate.RespCreateCsr{Code: generate.Code_Success}, nil
}

func (cs *CreationServer) GenerateOrg(ctx context.Context, in *generate.ReqCreateOrg) (*generate.RespCreateOrg, error) {
	gc := &geneses.GenerateConfig{}
	if err := gc.CreateOrg(in); nil != err {
		return &generate.RespCreateOrg{Code: generate.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &generate.RespCreateOrg{Code: generate.Code_Success}, nil
}

func (cs *CreationServer) GenerateOrgNode(ctx context.Context, in *generate.ReqCreateOrgNode) (*generate.RespCreateOrgNode, error) {
	gc := &geneses.GenerateConfig{}
	if err := gc.CreateOrgNode(in); nil != err {
		return &generate.RespCreateOrgNode{Code: generate.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &generate.RespCreateOrgNode{Code: generate.Code_Success}, nil
}

func (cs *CreationServer) GenerateOrgUser(ctx context.Context, in *generate.ReqCreateOrgUser) (*generate.RespCreateOrgUser, error) {
	gc := &geneses.GenerateConfig{}
	if err := gc.CreateOrgUser(in); nil != err {
		return &generate.RespCreateOrgUser{Code: generate.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &generate.RespCreateOrgUser{Code: generate.Code_Success}, nil
}

func (cs *CreationServer) GenerateGenesisBlock(ctx context.Context, in *generate.ReqGenesis) (*generate.RespGenesis, error) {
	genesis := geneses.Genesis{Info: in}
	genesis.Init()
	if err := genesis.CreateGenesisBlock("default"); nil != err {
		return &generate.RespGenesis{Code: generate.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &generate.RespGenesis{Code: generate.Code_Success}, nil
}

func (cs *CreationServer) GenerateChannelTx(ctx context.Context, in *generate.ReqChannelTx) (*generate.RespChannelTx, error) {
	genesis := geneses.Genesis{Info: in.Genesis}
	genesis.Init()
	if err := genesis.CreateChannelCreateTx("default", in.ChannelID); nil != err {
		return &generate.RespChannelTx{Code: generate.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &generate.RespChannelTx{Code: generate.Code_Success}, nil
}
