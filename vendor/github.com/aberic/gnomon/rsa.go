/*
 * Copyright (c) 2019. aberic - All Rights Reserved.
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

package gnomon

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	pksC1                PKSCType = "PKCS1"
	pksC8                PKSCType = "PKCS8"
	signPss              SignMode = "pss"
	signPKCS1v15         SignMode = "pkcs#1 v1.5"
	publicRSAKeyPemType           = "PUBLIC KEY"
	privateRSAKeyPemType          = "PRIVATE KEY"
)

// PKSCType 私钥格式，默认提供PKCS1和PKCS8
//
// 通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
type PKSCType string

// SignMode RSA签名模式，默认提供PSS和PKCS1v15
//
// 通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
type SignMode string

// RSACommon RSA工具
type RSACommon struct{}

// GenerateKey RSA公钥私钥产生
//
// bits 指定生成位大小
//
// path 指定公私钥所在生成目录
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) GenerateKey(bits int, path, priFileName, pubFileName string, pks PKSCType) error {
	return r.GenerateKeyWithPass(bits, path, priFileName, pubFileName, "", -1, pks)
}

// GenerateKeyWithPass RSA公钥私钥产生
//
// bits 指定生成位大小
//
// path 指定公私钥所在生成目录
//
// passwd 生成密码
//
// alg der编码数据指定算法，如：x509.PEMCipher3DES
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) GenerateKeyWithPass(bits int, path, priFileName, pubFileName, passwd string, alg x509.PEMCipher, pks PKSCType) error {
	privateKey, err := r.GeneratePriKeyWithPass(bits, path, priFileName, passwd, alg, pks)
	if nil != err {
		return err
	}
	return r.GeneratePubKey(privateKey, path, pubFileName, pksC8)
}

// GeneratePriKey RSA私钥产生
//
// bits 指定生成位大小
//
// path 指定私钥所在生成目录
//
// fileName 指定私钥的文件名称，如'private.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) GeneratePriKey(bits int, path, fileName string, pks PKSCType) (*rsa.PrivateKey, error) {
	return r.GeneratePriKeyWithPass(bits, path, fileName, "", -1, pks)
}

// GeneratePriKeyWithPass RSA私钥产生
//
// bits 指定生成位大小
//
// path 指定私钥所在生成目录
//
// fileName 指定私钥的文件名称，如'private.pem'
//
// passwd 生成密码
//
// alg der编码数据指定算法，如：x509.PEMCipher3DES
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) GeneratePriKeyWithPass(bits int, path, fileName, passwd string, alg x509.PEMCipher, pks PKSCType) (*rsa.PrivateKey, error) {
	var (
		privateKey *rsa.PrivateKey
		err        error
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return nil, err
		}
	}
	// 生成私钥文件
	if privateKey, err = rsa.GenerateKey(rand.Reader, bits); nil != err {
		return nil, err
	}
	if err = r.SavePriPemWithPass(privateKey, path, fileName, passwd, alg, pks); nil != err {
		return nil, err
	}
	return privateKey, nil
}

// GeneratePubKey RSA公钥产生
//
// privateKey 私钥
//
// path 指定公私钥所在生成目录
//
// fileName 指定公钥的文件名称，如'public.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) GeneratePubKey(privateKey *rsa.PrivateKey, path, fileName string, pks PKSCType) error {
	var (
		fileIO  *os.File
		derPkiX []byte
		err     error
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	// 将公钥序列化为der编码的PKIX格式
	if derPkiX, err = x509.MarshalPKIXPublicKey(publicKey); nil != err {
		return err
	}
	block := &pem.Block{
		Type:  publicRSAKeyPemType,
		Bytes: derPkiX,
	}
	defer func() { _ = fileIO.Close() }()
	if fileIO, err = os.Create(filepath.Join(path, fileName)); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return err
	}
	return nil
}

// GeneratePubKeyBytes RSA公钥产生
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// path 指定公私钥所在生成目录
//
// fileName 指定公钥的文件名称，如'public.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) GeneratePubKeyBytes(privateKey []byte, path, fileName string, pks PKSCType) error {
	return r.GeneratePubKeyBytesWithPass(privateKey, "", path, fileName, pks)
}

// GeneratePubKeyBytesWithPass RSA公钥产生
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// passwd 生成privateKey时输入密码
//
// path 指定公私钥所在生成目录
//
// fileName 指定公钥的文件名称，如'public.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) GeneratePubKeyBytesWithPass(privateKey []byte, passwd, path, fileName string, pks PKSCType) error {
	pri, err := r.LoadPriWithPass(privateKey, passwd, pks)
	if err != nil {
		return err
	}
	return r.GeneratePubKey(pri, path, fileName, pks)
}

// GeneratePubKeyFP RSA公钥产生
//
// privateKeyFilePath 私钥地址
//
// path 指定公私钥所在生成目录
//
// fileName 指定公钥的文件名称，如'public.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) GeneratePubKeyFP(privateKeyFilePath, path, fileName string, pks PKSCType) error {
	bs, err := ioutil.ReadFile(privateKeyFilePath)
	if nil != err {
		return err
	}
	return r.GeneratePubKeyBytes(bs, path, fileName, pks)
}

// GeneratePubKeyFPWithPass RSA公钥产生
//
// privateKeyFilePath 私钥地址
//
// passwd 生成privateKey时输入密码
//
// path 指定公私钥所在生成目录
//
// fileName 指定公钥的文件名称，如'public.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) GeneratePubKeyFPWithPass(privateKeyFilePath, passwd, path, fileName string, pks PKSCType) error {
	bs, err := ioutil.ReadFile(privateKeyFilePath)
	if nil != err {
		return err
	}
	return r.GeneratePubKeyBytesWithPass(bs, passwd, path, fileName, pks)
}

// SavePriPem 将私钥保存到给定文件
//
// alg der编码数据指定算法，如：x509.PEMCipher3DES
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) SavePriPem(privateKey *rsa.PrivateKey, path, fileName string, alg x509.PEMCipher, pks PKSCType) error {
	return r.SavePriPemWithPass(privateKey, path, fileName, "", alg, pks)
}

// SavePriPemWithPass 将私钥保存到给定文件
//
// alg der编码数据指定算法，如：x509.PEMCipher3DES
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) SavePriPemWithPass(privateKey *rsa.PrivateKey, path, fileName, passwd string, alg x509.PEMCipher, pks PKSCType) error {
	var (
		derStream []byte
		fileIO    *os.File
		block     *pem.Block
		err       error
	)
	// 将私钥转换为ASN.1 DER编码的形式
	switch pks {
	default:
		derStream = x509.MarshalPKCS1PrivateKey(privateKey)
	case pksC8:
		if derStream, err = x509.MarshalPKCS8PrivateKey(privateKey); nil != err {
			return err
		}
	}
	// block表示PEM编码的结构
	if String().IsEmpty(passwd) {
		block = &pem.Block{Type: privateRSAKeyPemType, Bytes: derStream}
	} else {
		block, err = x509.EncryptPEMBlock(rand.Reader, privateRSAKeyPemType, derStream, []byte(passwd), alg)
		if nil != err {
			return err
		}
	}
	defer func() { _ = fileIO.Close() }()
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, 0755); nil != err {
			return err
		}
	}
	if fileIO, err = os.Create(filepath.Join(path, fileName)); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return err
	}
	return nil
}

// Encrypt 公钥加密
//
// publicKey 公钥内容，如取出字符串'pubData'，则传入'string(pubData)'即可
//
// data 待加密数据
func (r *RSACommon) Encrypt(publicKey, data []byte) ([]byte, error) {
	pub, err := r.LoadPub(publicKey)
	if nil != err {
		return nil, err
	}
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// EncryptFP 公钥加密
//
// publicKeyFilePath 公钥地址
//
// data 待加密数据
func (r *RSACommon) EncryptFP(publicKeyFilePath string, data []byte) ([]byte, error) {
	pub, err := r.LoadPubFP(publicKeyFilePath)
	if nil != err {
		return nil, err
	}
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// Decrypt 私钥解密
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// data 待解密数据
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) Decrypt(privateKey, data []byte, pks PKSCType) ([]byte, error) {
	return r.DecryptWithPass(privateKey, data, "", pks)
}

// DecryptWithPass 私钥解密
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// data 待解密数据
//
// passwd 生成privateKey时输入密码
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) DecryptWithPass(privateKey, data []byte, passwd string, pks PKSCType) ([]byte, error) {
	pri, err := r.LoadPriWithPass(privateKey, passwd, pks)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, pri, data)
}

// DecryptFP 私钥解密
//
// privateKeyFilePath 私钥地址
//
// data 待解密数据
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) DecryptFP(privateKeyFilePath string, data []byte, pks PKSCType) ([]byte, error) {
	return r.DecryptFPWithPass(privateKeyFilePath, "", data, pks)
}

// DecryptFPWithPass 私钥解密
//
// privateKeyFilePath 私钥地址
//
// passwd 生成privateKey时输入密码
//
// data 待解密数据
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) DecryptFPWithPass(privateKeyFilePath, passwd string, data []byte, pks PKSCType) ([]byte, error) {
	pri, err := r.LoadPriFPWithPass(privateKeyFilePath, passwd, pks)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, pri, data)
}

// Sign 签名：采用RSA-PKCS#1 v1.5模式
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
//
// mod RSA签名模式，默认提供PSS和PKCS1v15，通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
func (r *RSACommon) Sign(privateKey, data []byte, hash crypto.Hash, pks PKSCType, mod SignMode) ([]byte, error) {
	return r.SignWithPass(privateKey, data, "", hash, pks, mod)
}

// SignWithPass 签名：采用RSA-PKCS#1 v1.5模式
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// data 待签名数据
//
// passwd 生成privateKey时输入密码
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
//
// mod RSA签名模式，默认提供PSS和PKCS1v15，通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
func (r *RSACommon) SignWithPass(privateKey, data []byte, passwd string, hash crypto.Hash, pks PKSCType, mod SignMode) ([]byte, error) {
	pri, err := r.LoadPriWithPass(privateKey, passwd, pks)
	if err != nil {
		return nil, err
	}
	switch mod {
	default:
		return r.signPSS(pri, data, hash)
	case signPKCS1v15:
		return r.signPKCS1v15(pri, data, hash)
	}
}

// SignFP 签名：采用RSA-PKCS#1 v1.5模式
//
// privateKeyPath 私钥文件存储路径
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
//
// mod RSA签名模式，默认提供PSS和PKCS1v15，通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
func (r *RSACommon) SignFP(privateKeyPath string, data []byte, hash crypto.Hash, pks PKSCType, mod SignMode) ([]byte, error) {
	pri, err := r.LoadPriFP(privateKeyPath, pks)
	if err != nil {
		return nil, err
	}
	switch mod {
	default:
		return r.signPSS(pri, data, hash)
	case signPKCS1v15:
		return r.signPKCS1v15(pri, data, hash)
	}
}

// SignFPWithPass 签名：采用RSA-PKCS#1 v1.5模式
//
// privateKeyPath 私钥文件存储路径
//
// passwd 生成privateKey时输入密码
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
//
// mod RSA签名模式，默认提供PSS和PKCS1v15，通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
func (r *RSACommon) SignFPWithPass(privateKeyPath, passwd string, data []byte, hash crypto.Hash, pks PKSCType, mod SignMode) ([]byte, error) {
	pri, err := r.LoadPriFPWithPass(privateKeyPath, passwd, pks)
	if err != nil {
		return nil, err
	}
	switch mod {
	default:
		return r.signPSS(pri, data, hash)
	case signPKCS1v15:
		return r.signPKCS1v15(pri, data, hash)
	}
}

// signPSS 签名：采用PSS模式
//
// privateKey 私钥
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
func (r *RSACommon) signPSS(privateKey *rsa.PrivateKey, data []byte, hash crypto.Hash) ([]byte, error) {
	h := hash.New()
	if _, err := h.Write(data); nil != err {
		return nil, err
	}
	hashed := h.Sum(nil)
	return rsa.SignPSS(rand.Reader, privateKey, hash, hashed, nil)
}

// signPKCS1v15 签名：采用RSA-PKCS#1 v1.5模式
//
// privateKey 私钥
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
func (r *RSACommon) signPKCS1v15(privateKey *rsa.PrivateKey, data []byte, hash crypto.Hash) ([]byte, error) {
	h := hash.New()
	if _, err := h.Write(data); nil != err {
		return nil, err
	}
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, hash, hashed)
}

// Verify 验签：采用RSA-PKCS#1 v1.5模式
//
// publicKey 公钥内容，如取出字符串'pubData'，则传入'string(pubData)'即可
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// mod RSA签名模式，默认提供PSS和PKCS1v15，通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
func (r *RSACommon) Verify(publicKey, data, signData []byte, hash crypto.Hash, mod SignMode) error {
	pub, err := r.LoadPub(publicKey)
	if nil != err {
		return err
	}
	switch mod {
	default:
		return r.verifyPSS(pub, data, signData, hash)
	case signPKCS1v15:
		return r.verifyPKCS1v15(pub, data, signData, hash)
	}
}

// VerifyFP 验签：采用RSA-PKCS#1 v1.5模式
//
// publicKeyPath 公钥文件存储路径
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// mod RSA签名模式，默认提供PSS和PKCS1v15，通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
func (r *RSACommon) VerifyFP(publicKeyPath string, data, signData []byte, hash crypto.Hash, mod SignMode) error {
	pub, err := r.LoadPubFP(publicKeyPath)
	if nil != err {
		return err
	}
	switch mod {
	default:
		return r.verifyPSS(pub, data, signData, hash)
	case signPKCS1v15:
		return r.verifyPKCS1v15(pub, data, signData, hash)
	}
}

// verifyPSS 验签：采用PSS模式
//
// publicKey 公钥
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
func (r *RSACommon) verifyPSS(publicKey *rsa.PublicKey, data, signData []byte, hash crypto.Hash) error {
	h := hash.New()
	if _, err := h.Write(data); nil != err {
		return err
	}
	hashed := h.Sum(nil)
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: hash}
	return rsa.VerifyPSS(publicKey, hash, hashed, signData, opts)
}

// verifyPKCS1v15 验签：采用RSA-PKCS#1 v1.5模式
//
// publicKey 公钥
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
func (r *RSACommon) verifyPKCS1v15(publicKey *rsa.PublicKey, data, signData []byte, hash crypto.Hash) error {
	h := hash.New()
	if _, err := h.Write(data); nil != err {
		return err
	}
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(publicKey, hash, hashed, signData)
}

// LoadPriFP 加载私钥
//
// privateKeyFilePath 私钥地址
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) LoadPriFP(privateKeyFilePath string, pks PKSCType) (*rsa.PrivateKey, error) {
	return r.LoadPriFPWithPass(privateKeyFilePath, "", pks)
}

// LoadPriFPWithPass 加载私钥
//
// privateKeyFilePath 私钥地址
//
// passwd 生成privateKey时输入密码
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) LoadPriFPWithPass(privateKeyFilePath, passwd string, pks PKSCType) (*rsa.PrivateKey, error) {
	bs, err := ioutil.ReadFile(privateKeyFilePath)
	if nil != err {
		return nil, err
	}
	return r.LoadPriWithPass(bs, passwd, pks)
}

// LoadPubFP 加载公钥
//
// publicKeyFilePath 公钥地址
func (r *RSACommon) LoadPubFP(publicKeyFilePath string) (*rsa.PublicKey, error) {
	bs, err := ioutil.ReadFile(publicKeyFilePath)
	if nil != err {
		return nil, err
	}
	return r.LoadPub(bs)
}

// LoadPri 解析私钥
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) LoadPri(privateKey []byte, pks PKSCType) (*rsa.PrivateKey, error) {
	return r.LoadPriWithPass(privateKey, "", pks)
}

// LoadPriWithPass 解析私钥
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// passwd 生成privateKey时输入密码
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (r *RSACommon) LoadPriWithPass(privateKey []byte, passwd string, pks PKSCType) (*rsa.PrivateKey, error) {
	var (
		pemData []byte
		err     error
	)
	if String().IsEmpty(passwd) {
		pemData, err = r.pemParse(privateKey, privateRSAKeyPemType)
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
	switch pks {
	default:
		pri, err := x509.ParsePKCS8PrivateKey(pemData)
		if nil != err {
			return nil, err
		}
		return pri.(*rsa.PrivateKey), nil
	case pksC1:
		return x509.ParsePKCS1PrivateKey(pemData)
	}
}

// LoadPub 加载公钥
//
// publicKey 公钥内容，如取出字符串'pubData'，则传入'string(pubData)'即可
func (r *RSACommon) LoadPub(publicKey []byte) (*rsa.PublicKey, error) {
	pemData, err := r.pemParse(publicKey, publicRSAKeyPemType)
	if err != nil {
		return nil, err
	}
	keyInterface, err := x509.ParsePKIXPublicKey(pemData)
	if err != nil {
		return nil, err
	}
	pubKey, ok := keyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("could not cast parsed key to *rsa.PublicKey")
	}
	return pubKey, nil
}

// pemParse 解密pem格式密钥并验证pem类型
func (r *RSACommon) pemParse(key []byte, pemType string) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("no pem block found")
	}
	if pemType != "" && block.Type != pemType {
		return nil, errors.New(strings.Join([]string{"Key's type is ", block.Type, ", expected ", pemType}, ""))
	}
	return block.Bytes, nil
}

// PKSC1 私钥格PKCS1
func (r *RSACommon) PKSC1() PKSCType {
	return pksC1
}

// PKSC8 私钥格式PKCS8
func (r *RSACommon) PKSC8() PKSCType {
	return pksC8
}

// SignPSS RSA签名模式，PSS
func (r *RSACommon) SignPSS() SignMode {
	return signPss
}

// SignPKCS RSA签名模式，PKCS1v15
func (r *RSACommon) SignPKCS() SignMode {
	return signPKCS1v15
}
