/*
 *  Copyright (c) 2019. aberic - All Rights Reserved.
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

package gnomon

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

const (
	privateECCKeyPemType = "PRIVATE KEY"
	publicECCKeyPemType  = "PUBLIC KEY"
)

// ECCCommon ECC椭圆加密工具，依赖ETH的包
//
// ECC，全称椭圆曲线密码学（英语：Elliptic curve cryptography，缩写为 ECC），主要是指相关数学原理
//
// ECIES，在ECC原理的基础上实现的一种公钥加密方法，和RSA类似
//
// ECDSA，在ECC原理上实现的签名方法
//
// ECDH在ECC和DH的基础上实现的密钥交换算法
type ECCCommon struct{}

// Generate 生成公私钥对
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (e *ECCCommon) Generate(curve elliptic.Curve) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	var (
		privateKey *ecdsa.PrivateKey
		err        error
	)
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// GenerateKey 生成公私钥对
//
// path 指定公私钥所在生成目录
//
// priFileName 指定生成的密钥名称
//
// pubFileName 指定生成的密钥名称
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (e *ECCCommon) GenerateKey(path, priFileName, pubFileName string, curve elliptic.Curve) error {
	var (
		privateKey *ecdsa.PrivateKey
		err        error
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return err
	}

	if err = e.SavePri(filepath.Join(path, priFileName), privateKey); nil != err {
		return err
	}
	if err = e.SavePub(filepath.Join(path, pubFileName), &privateKey.PublicKey, curve); nil != err {
		return err
	}
	return nil
}

// GeneratePemKey 生成公私钥对
//
// path 指定公私钥所在生成目录
//
// priFileName 指定生成的密钥名称
//
// pubFileName 指定生成的密钥名称
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (e *ECCCommon) GeneratePemKey(path, priFileName, pubFileName string, curve elliptic.Curve) error {
	return e.GeneratePemKeyWithPass(path, priFileName, pubFileName, "", curve)
}

// GeneratePemKeyWithPass 生成公私钥对
//
// path 指定公私钥所在生成目录
//
// priFileName 指定生成的密钥名称
//
// pubFileName 指定生成的密钥名称
//
// passwd 生成密码
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (e *ECCCommon) GeneratePemKeyWithPass(path, priFileName, pubFileName, passwd string, curve elliptic.Curve) error {
	var (
		privateKey *ecdsa.PrivateKey
		err        error
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return err
	}

	if err = e.SavePriPemWithPass(privateKey, passwd, filepath.Join(path, priFileName)); nil != err {
		return err
	}
	if err = e.SavePubPem(filepath.Join(path, pubFileName), &privateKey.PublicKey); nil != err {
		return err
	}
	return nil
}

// GeneratePriKey 生成私钥
//
// path 指定私钥所在生成目录
//
// priFileName 指定生成的密钥名称
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (e *ECCCommon) GeneratePriKey(path, priFileName string, curve elliptic.Curve) error {
	var (
		privateKey *ecdsa.PrivateKey
		err        error
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return err
	}
	if err = e.SavePri(filepath.Join(path, priFileName), privateKey); nil != err {
		return err
	}
	return nil
}

// GeneratePemPriKey 生成私钥
//
// path 指定私钥所在生成目录
//
// priFileName 指定生成的密钥名称
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (e *ECCCommon) GeneratePemPriKey(path, priFileName string, curve elliptic.Curve) error {
	return e.GeneratePemPriKeyWithPass(path, priFileName, "", curve)
}

// GeneratePemPriKeyWithPass 生成私钥
//
// path 指定私钥所在生成目录
//
// priFileName 指定生成的密钥名称
//
// passwd 生成时输入的密码
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (e *ECCCommon) GeneratePemPriKeyWithPass(path, priFileName, passwd string, curve elliptic.Curve) error {
	var (
		privateKey *ecdsa.PrivateKey
		err        error
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return err
	}
	if err = e.SavePriPemWithPass(privateKey, passwd, filepath.Join(path, priFileName)); nil != err {
		return err
	}
	return nil
}

// GeneratePubKey 生成公钥
//
// path 指定公钥所在生成目录
func (e *ECCCommon) GeneratePubKey(privateKey *ecdsa.PrivateKey, path, pubFileName string, curve elliptic.Curve) error {
	if err := e.SavePub(filepath.Join(path, pubFileName), &privateKey.PublicKey, curve); nil != err {
		return err
	}
	return nil
}

// GeneratePemPubKey 生成公钥
//
// path 指定公钥所在生成目录
func (e *ECCCommon) GeneratePemPubKey(privateKey *ecdsa.PrivateKey, path, pubFileName string) error {
	if err := e.SavePubPem(filepath.Join(path, pubFileName), &privateKey.PublicKey); nil != err {
		return err
	}
	return nil
}

// SavePri 将私钥保存到给定文件，密钥数据保存为hex编码
func (e *ECCCommon) SavePri(file string, privateKey *ecdsa.PrivateKey) error {
	k := hex.EncodeToString(e.PriKey2Bytes(privateKey))
	return ioutil.WriteFile(file, []byte(k), 0600)
	//return crypto.SaveECDSA(file, privateKey)
}

// LoadPri 从文件中加载私钥
//
// file 文件路径
func (e *ECCCommon) LoadPri(file string, curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	bs, err := ioutil.ReadFile(file)
	if nil != err {
		return nil, err
	}
	data, err := hex.DecodeString(string(bs))
	if err != nil {
		return nil, err
	}
	return e.Bytes2PriKey(data, curve), nil
	//return crypto.LoadECDSA(file)
}

// SavePriPem 将私钥保存到给定文件
func (e *ECCCommon) SavePriPem(privateKey *ecdsa.PrivateKey, file string) error {
	return e.SavePriPemWithPass(privateKey, "", file)
}

// SavePriPemWithPass 将私钥保存到给定文件
func (e *ECCCommon) SavePriPemWithPass(privateKey *ecdsa.PrivateKey, passwd, file string) error {
	var (
		fileIO *os.File
		block  *pem.Block
	)
	// 将私钥转换为ASN.1 DER编码的形式
	derStream, err := x509.MarshalECPrivateKey(privateKey)
	if nil != err {
		return err
	}
	// block表示PEM编码的结构
	if String().IsEmpty(passwd) {
		block = &pem.Block{Type: privateECCKeyPemType, Bytes: derStream}
	} else {
		block, err = x509.EncryptPEMBlock(rand.Reader, privateECCKeyPemType, derStream, []byte(passwd), x509.PEMCipher3DES)
		if nil != err {
			return err
		}
	}
	defer func() { _ = fileIO.Close() }()
	if fileIO, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return err
	}
	return nil
}

// LoadPriPem 解析私钥
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
func (e *ECCCommon) LoadPriPem(privateKey []byte) (*ecdsa.PrivateKey, error) {
	return e.LoadPriPemWithPass(privateKey, "")
}

// LoadPriPemWithPass 解析私钥
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// passwd 生成privateKey时输入密码
func (e *ECCCommon) LoadPriPemWithPass(privateKey []byte, passwd string) (*ecdsa.PrivateKey, error) {
	var (
		pemData []byte
		err     error
	)
	if String().IsEmpty(passwd) {
		pemData, err = e.pemParse(privateKey, privateECCKeyPemType)
		if err != nil {
			return nil, err
		}
	} else {
		block, _ := pem.Decode(privateKey)
		pemData, err = x509.DecryptPEMBlock(block, []byte(passwd))
		if err != nil {
			return nil, err
		}
	}
	pri, err := x509.ParseECPrivateKey(pemData)
	if nil != err {
		return nil, err
	}
	return pri, nil
}

// LoadPriPemFP 从文件中加载私钥
//
// file 文件路径
func (e *ECCCommon) LoadPriPemFP(file string) (*ecdsa.PrivateKey, error) {
	return e.LoadPriPemFPWithPass(file, "")
}

// LoadPriPemFPWithPass 从文件中加载私钥
//
// file 文件路径
//
// passwd 生成privateKey时输入密码
func (e *ECCCommon) LoadPriPemFPWithPass(file, passwd string) (*ecdsa.PrivateKey, error) {
	var (
		pemData []byte
		err     error
	)
	keyData, err := ioutil.ReadFile(file)
	if nil != err {
		return nil, err
	}
	if String().IsEmpty(passwd) {
		pemData, err = e.pemParse(keyData, privateECCKeyPemType)
		if err != nil {
			return nil, err
		}
	} else {
		block, _ := pem.Decode(keyData)
		pemData, err = x509.DecryptPEMBlock(block, []byte(passwd))
		if err != nil {
			return nil, err
		}
	}
	pri, err := x509.ParseECPrivateKey(pemData)
	if nil != err {
		return nil, err
	}
	return pri, nil
}

// SavePubPem 将公钥保存到给定文件
//
// file 文件路径
func (e *ECCCommon) SavePubPem(file string, publicKey *ecdsa.PublicKey) error {
	var fileIO *os.File
	// 将公钥序列化为der编码的PKIX格式
	derPkiX, err := x509.MarshalPKIXPublicKey(publicKey)
	if nil != err {
		return err
	}
	block := &pem.Block{
		Type:  publicECCKeyPemType,
		Bytes: derPkiX,
	}
	defer func() { _ = fileIO.Close() }()
	if fileIO, err = os.Create(file); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return err
	}
	return nil
}

// LoadPubPem 从文件中加载公钥
//
// file 文件路径
func (e *ECCCommon) LoadPubPem(file string) (*ecdsa.PublicKey, error) {
	pubData, err := ioutil.ReadFile(file)
	if nil != err {
		return nil, err
	}
	pemData, err := e.pemParse(pubData, publicECCKeyPemType)
	if err != nil {
		return nil, err
	}
	keyInterface, err := x509.ParsePKIXPublicKey(pemData)
	if err != nil {
		return nil, err
	}
	pubKey, ok := keyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("could not cast parsed key to *rsa.PublicKey")
	}
	return pubKey, nil
}

// SavePub 将公钥保存到给定文件，密钥数据保存为hex编码
//
// file 文件路径
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (e *ECCCommon) SavePub(file string, publicKey *ecdsa.PublicKey, curve elliptic.Curve) error {
	k := hex.EncodeToString(e.PubKey2Bytes(publicKey, curve))
	return ioutil.WriteFile(file, []byte(k), 0600)
}

// LoadPub 从文件中加载公钥
//
// file 文件路径
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (e *ECCCommon) LoadPub(file string, curve elliptic.Curve) (*ecdsa.PublicKey, error) {
	data, err := ioutil.ReadFile(file)
	if nil != err {
		return nil, err
	}
	key, err := hex.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	return e.Bytes2PubKey(key, curve)
}

// PriKey2Bytes 私钥转[]byte
func (e *ECCCommon) PriKey2Bytes(privateKey *ecdsa.PrivateKey) []byte {
	return privateKey.D.Bytes()
}

// Bytes2PriKey []byte转私钥
func (e *ECCCommon) Bytes2PriKey(data []byte, curve elliptic.Curve) *ecdsa.PrivateKey {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = new(big.Int).SetBytes(data)
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(data)
	return priv
}

// PubKey2Bytes 公钥转[]byte
//
// pub 公钥
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (e *ECCCommon) PubKey2Bytes(publicKey *ecdsa.PublicKey, curve elliptic.Curve) []byte {
	if publicKey == nil || publicKey.X == nil || publicKey.Y == nil {
		return nil
	}
	return elliptic.Marshal(curve, publicKey.X, publicKey.Y)
}

// Bytes2PubKey []byte转公钥
func (e *ECCCommon) Bytes2PubKey(data []byte, curve elliptic.Curve) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(curve, data)
	if x == nil {
		return nil, errors.New("invalid public key")
	}
	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
}

// Encrypt 加密
func (e *ECCCommon) Encrypt(data []byte, publicKey *ecies.PublicKey) ([]byte, error) {
	if publicKey.Curve.Params().BitSize != 256 {
		return nil, errors.New("just support P256 and S256")
	}
	return ecies.Encrypt(rand.Reader, publicKey, data, nil, nil)
}

//// EncryptFP 加密
////
//// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
//func (e *ECCCommon) EncryptFP(data []byte, publicKeyFilePath string, curve elliptic.Curve) ([]byte, error) {
//	publicKey, err := e.LoadPub(publicKeyFilePath, curve)
//	if nil != err {
//		return nil, err
//	}
//	return e.Encrypt(data, ecies.ImportECDSAPublic(publicKey))
//}

// Decrypt 解密
func (e *ECCCommon) Decrypt(data []byte, privateKey *ecies.PrivateKey) ([]byte, error) {
	if privateKey.Curve.Params().BitSize != 256 {
		return nil, errors.New("just support P256 and S256")
	}
	return privateKey.Decrypt(data, nil, nil)
}

//// DecryptFP 解密
//func (e *ECCCommon) DecryptFP(data []byte, privateKeyFilePath string) ([]byte, error) {
//	privateKey, err := e.LoadPri(privateKeyFilePath)
//	if nil != err {
//		return nil, err
//	}
//	return e.Decrypt(data, ecies.ImportECDSA(privateKey))
//}

// Sign 签名
func (e *ECCCommon) Sign(privateKey *ecdsa.PrivateKey, data []byte) (sign []byte, err error) {
	// 根据明文plaintext和私钥，生成两个big.Ing
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, data)
	if err != nil {
		return nil, err
	}
	rs, err := r.MarshalText()
	if err != nil {
		return nil, err
	}
	ss, err := s.MarshalText()
	if err != nil {
		return nil, err
	}
	// 将r，s合并（以“+”分割），作为签名返回
	var b bytes.Buffer
	b.Write(rs)
	b.Write([]byte(`+`))
	b.Write(ss)
	return b.Bytes(), nil
}

// Verify 验签
func (e *ECCCommon) Verify(publicKey *ecdsa.PublicKey, data, sign []byte) (bool, error) {
	var rint, sint big.Int
	// 根据sign，解析出r，s
	rs := bytes.Split(sign, []byte("+"))
	if err := rint.UnmarshalText(rs[0]); nil != err {
		return false, err
	}
	if err := sint.UnmarshalText(rs[1]); nil != err {
		return false, err
	}
	// 根据公钥，明文，r，s验证签名
	v := ecdsa.Verify(publicKey, data, &rint, &sint)
	return v, nil
}

// pemParse 解密pem格式密钥并验证pem类型
func (e *ECCCommon) pemParse(key []byte, pemType string) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("no pem block found")
	}
	if pemType != "" && block.Type != pemType {
		return nil, errors.New(strings.Join([]string{"Key's type is ", block.Type, ", expected ", pemType}, ""))
	}
	return block.Bytes, nil
}
