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

package ca

import (
	"crypto/elliptic"
	"errors"
	"github.com/aberic/fabric-client/grpc/proto/generate"
	"github.com/aberic/gnomon"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

const (
	pemPriKeyFileDefaultName = "pri.key"
	pemPubKeyFileDefaultName = "pub.key"
)

type PemConfig struct {
	KeyConfig *generate.ReqKeyConfig
}

func (pc *PemConfig) GenerateCrypto() *generate.RespKeyConfig {
	switch pc.KeyConfig.CryptoType {
	default:
		return &generate.RespKeyConfig{Code: generate.Code_Fail, ErrMsg: "generate crypto type error"}
	case generate.CryptoType_RSA:
		return pc.cryptoRSA()
	case generate.CryptoType_ECDSA:
		return pc.cryptoECC()
	}
}

func (pc *PemConfig) cryptoRSA() *generate.RespKeyConfig {
	var (
		bits int
		err  error
	)
	if bits, err = pc.cryptoRSABits(); nil != err {
		return &generate.RespKeyConfig{Code: generate.Code_Fail, ErrMsg: err.Error()}
	}
	storePath := path.Join("/tmp", strconv.Itoa(time.Now().Nanosecond()))
	priFilePath := filepath.Join(storePath, pemPriKeyFileDefaultName)
	pubFilePath := filepath.Join(storePath, pemPubKeyFileDefaultName)
	if err = gnomon.CryptoRSA().GenerateKey(bits, storePath, pemPriKeyFileDefaultName, pemPubKeyFileDefaultName, gnomon.CryptoRSA().PKSC8()); nil != err {
		return &generate.RespKeyConfig{Code: generate.Code_Fail, ErrMsg: err.Error()}
	}
	return &generate.RespKeyConfig{Code: generate.Code_Success, PriKeyFilePath: priFilePath, PubKeyFilePath: pubFilePath}
}

func (pc *PemConfig) cryptoRSABits() (bits int, err error) {
	switch pc.KeyConfig.GetRsaAlgorithm() {
	default:
		return 0, errors.New("rsa algorithm type error")
	case generate.RsaAlgorithm_r2048:
		return 2048, nil
	case generate.RsaAlgorithm_r4096:
		return 4096, nil
	}
}

func (pc *PemConfig) cryptoECC() *generate.RespKeyConfig {
	var (
		curve elliptic.Curve
		err   error
	)
	if curve, err = pc.cryptoECCCurve(); nil != err {
		return &generate.RespKeyConfig{Code: generate.Code_Fail, ErrMsg: err.Error()}
	}
	storePath := path.Join("/tmp", strconv.Itoa(time.Now().Nanosecond()))
	priFilePath := filepath.Join(storePath, pemPriKeyFileDefaultName)
	pubFilePath := filepath.Join(storePath, pemPubKeyFileDefaultName)
	if err = gnomon.CryptoECC().GeneratePemKey(storePath, pemPriKeyFileDefaultName, pemPubKeyFileDefaultName, curve); nil != err {
		return &generate.RespKeyConfig{Code: generate.Code_Fail, ErrMsg: err.Error()}
	}
	return &generate.RespKeyConfig{Code: generate.Code_Success, PriKeyFilePath: priFilePath, PubKeyFilePath: pubFilePath}
}

func (pc *PemConfig) cryptoECCCurve() (curve elliptic.Curve, err error) {
	switch pc.KeyConfig.GetEccAlgorithm() {
	default:
		return nil, errors.New("ecc algorithm type error")
	case generate.EccAlgorithm_p256:
		return elliptic.P256(), nil
	case generate.EccAlgorithm_p384:
		return elliptic.P384(), nil
	case generate.EccAlgorithm_p521:
		return elliptic.P521(), nil
	}
}
