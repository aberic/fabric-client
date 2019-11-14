/*
 *
 *  * Copyright (c) 2019. aberic - All Rights Reserved.
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  * http://www.apache.org/licenses/LICENSE-2.0
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 *
 */

package gnomon

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	rd "math/rand"
	"net"
	"net/url"
	"os"
	"time"
)

const (
	certificateRequestType = "CERTIFICATE REQUEST"
	certificateType        = "CERTIFICATE"
)

// CACommon CA工具
type CACommon struct{}

// GenerateRSACertificateRequest 生成证书签名请求文件
//
// cert 证书生成请求对象
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (ca *CACommon) GenerateRSACertificateRequest(cert *CertRequest, pks PKSCType) (csr []byte, err error) {
	priRSAKey, err := CryptoRSA().LoadPri(cert.PrivateKeyData, pks)
	if err != nil {
		return nil, err
	}
	certModel := cert.init()
	certModel.PrivateKey = priRSAKey
	return ca.GenerateCertificateRequest(certModel)
}

// GenerateRSACertificateRequestWithPass 生成证书签名请求文件
//
// cert 证书生成请求对象
//
// password 生成时输入的密码
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (ca *CACommon) GenerateRSACertificateRequestWithPass(cert *CertRequest, password string, pks PKSCType) (csr []byte, err error) {
	priRSAKey, err := CryptoRSA().LoadPriWithPass(cert.PrivateKeyData, password, pks)
	if err != nil {
		return nil, err
	}
	certModel := cert.init()
	certModel.PrivateKey = priRSAKey
	return ca.GenerateCertificateRequest(certModel)
}

// GenerateRSACertificateRequestFP 生成证书签名请求文件
//
// cert 证书生成请求对象
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (ca *CACommon) GenerateRSACertificateRequestFP(cert *CertRequestFP, pks PKSCType) (csr []byte, err error) {
	priRSAKey, err := CryptoRSA().LoadPriFP(cert.PrivateKeyFilePath, pks)
	if err != nil {
		return nil, err
	}
	certModel := cert.init()
	certModel.PrivateKey = priRSAKey
	return ca.GenerateCertificateRequest(certModel)
}

// GenerateRSACertificateRequestFPWithPass 生成证书签名请求文件
//
// cert 证书生成请求对象
//
// password 生成时输入的密码
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func (ca *CACommon) GenerateRSACertificateRequestFPWithPass(cert *CertRequestFP, password string, pks PKSCType) (csr []byte, err error) {
	priRSAKey, err := CryptoRSA().LoadPriFPWithPass(cert.PrivateKeyFilePath, password, pks)
	if err != nil {
		return nil, err
	}
	certModel := cert.init()
	certModel.PrivateKey = priRSAKey
	return ca.GenerateCertificateRequest(certModel)
}

// GenerateECCCertificateRequest 生成证书签名请求文件
//
// cert 证书生成请求对象
func (ca *CACommon) GenerateECCCertificateRequest(cert *CertRequest) (csr []byte, err error) {
	return ca.GenerateECCCertificateRequestWithPass(cert, "")
}

// GenerateECCCertificateRequestWithPass 生成证书签名请求文件
//
// cert 证书生成请求对象
//
// password 生成时输入的密码
func (ca *CACommon) GenerateECCCertificateRequestWithPass(cert *CertRequest, password string) (csr []byte, err error) {
	priECCKey, err := CryptoECC().LoadPriPemWithPass(cert.PrivateKeyData, password)
	if err != nil {
		return nil, err
	}
	certModel := cert.init()
	certModel.PrivateKey = priECCKey
	return ca.GenerateCertificateRequest(certModel)
}

// GenerateECCCertificateRequestFP 生成证书签名请求文件
//
// cert 证书生成请求对象
func (ca *CACommon) GenerateECCCertificateRequestFP(cert *CertRequestFP) (csr []byte, err error) {
	return ca.GenerateECCCertificateRequestFPWithPass(cert, "")
}

// GenerateECCCertificateRequestFPWithPass 生成证书签名请求文件
//
// cert 证书生成请求对象
//
// password 生成时输入的密码
func (ca *CACommon) GenerateECCCertificateRequestFPWithPass(cert *CertRequestFP, password string) (csr []byte, err error) {
	priECCKey, err := CryptoECC().LoadPriPemFPWithPass(cert.PrivateKeyFilePath, password)
	if err != nil {
		return nil, err
	}
	certModel := cert.init()
	certModel.PrivateKey = priECCKey
	return ca.GenerateCertificateRequest(certModel)
}

// GenerateCertificateRequest 生成证书签名请求文件
//
// cert 证书生成请求对象
func (ca *CACommon) GenerateCertificateRequest(cert *CertRequestModel) (csr []byte, err error) {
	csrData, err := x509.CreateCertificateRequest(rand.Reader, cert.Template, cert.PrivateKey)
	if nil != err {
		return nil, err
	}
	fileIO, err := os.OpenFile(cert.CertificateRequestFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if nil != err {
		return nil, err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, &pem.Block{Type: certificateRequestType, Bytes: csrData}); nil != err {
		return nil, err
	}
	return csrData, nil
}

func (ca *CACommon) LoadCsrPemFromFP(csrFilePath string) (cert *x509.CertificateRequest, err error) {
	data, err := ioutil.ReadFile(csrFilePath)
	if nil != err {
		return nil, err
	}
	csrData, _ := pem.Decode(data)
	return x509.ParseCertificateRequest(csrData.Bytes)
}

// GenerateCertificateSelf 对签名请求进行处理并生成自签名数字证书
//
// cert 签名数字证书对象
func (ca *CACommon) GenerateCertificateSelf(cert *CertSelf) (certData []byte, err error) {
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(rd.Int63()), // 证书序列号
		Subject:               cert.Subject,
		NotBefore:             cert.NotBeforeDays,
		NotAfter:              cert.NotAfterDays,
		BasicConstraintsValid: cert.BasicConstraintsValid,
		IsCA:                  cert.IsCA,
		SignatureAlgorithm:    cert.SignatureAlgorithm,
		ExtKeyUsage:           cert.ExtKeyUsage,
		KeyUsage:              cert.KeyUsage,
		SubjectKeyId:          []byte{1, 2, 3},
	}
	certData, err = x509.CreateCertificate(rand.Reader, template, template, cert.PublicKey, cert.ParentPrivateKey)
	if err != nil {
		return nil, err
	}
	path := File().ParentPath(cert.CertificateFilePath)
	// 创建生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return nil, err
		}
	}
	fileIO, err := os.OpenFile(cert.CertificateFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if nil != err {
		return nil, err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, &pem.Block{Type: certificateType, Bytes: certData}); nil != err {
		return nil, err
	}
	return certData, nil
}

// GenerateCertificate 对签名请求进行处理并生成签名数字证书
//
// cert 签名数字证书对象
func (ca *CACommon) GenerateCertificate(cert *Cert) (certData []byte, err error) {
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(rd.Int63()), // 证书序列号
		Subject:               cert.Subject,
		NotBefore:             cert.NotBeforeDays,
		NotAfter:              cert.NotAfterDays,
		BasicConstraintsValid: cert.BasicConstraintsValid,
		IsCA:                  cert.IsCA,
		SignatureAlgorithm:    cert.SignatureAlgorithm,
		ExtKeyUsage:           cert.ExtKeyUsage,
		KeyUsage:              cert.KeyUsage,
		SubjectKeyId:          []byte{1, 2, 3},
	}
	certData, err = x509.CreateCertificate(rand.Reader, template, cert.ParentCert, cert.PublicKey, cert.ParentPrivateKey)
	if err != nil {
		return nil, err
	}
	path := File().ParentPath(cert.CertificateFilePath)
	// 创建生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return nil, err
		}
	}
	fileIO, err := os.OpenFile(cert.CertificateFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if nil != err {
		return nil, err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, &pem.Block{Type: certificateType, Bytes: certData}); nil != err {
		return nil, err
	}
	return certData, nil
}

// LoadCrtFromFP 从文件中加载Crt对象
func (ca *CACommon) LoadCrtFromFP(crtFilePath string) (certificate *x509.Certificate, err error) {
	data, err := ioutil.ReadFile(crtFilePath)
	if nil != err {
		return nil, err
	}
	certData, _ := pem.Decode(data)
	return x509.ParseCertificate(certData.Bytes)
}

// CertSelf 自签名数字证书对象
type CertSelf struct {
	CertificateFilePath         string                  // 签名后数字证书文件存储路径
	Subject                     pkix.Name               // Subject 签名信息
	ParentPrivateKey, PublicKey interface{}             // 公私钥
	BasicConstraintsValid       bool                    // 基本的有效性约束
	IsCA                        bool                    // 是否是根证书
	NotBeforeDays, NotAfterDays time.Time               // 在指定时间之后生效及之前失效
	ExtKeyUsage                 []x509.ExtKeyUsage      // ExtKeyUsage表示对给定键有效的扩展操作集。每个ExtKeyUsage*常量定义一个惟一的操作。
	KeyUsage                    x509.KeyUsage           // KeyUsage表示对给定密钥有效的操作集。它是KeyUsage*常量的位图。
	SignatureAlgorithm          x509.SignatureAlgorithm // signatureAlgorithm 生成证书时候采用的签名算法
}

// Cert 签名数字证书对象
type Cert struct {
	ParentCert *x509.Certificate // 父证书对象
	CertSelf
}

// CertRequest 证书生成请求对象
type CertRequest struct {
	PrivateKeyData             []byte                  // privateKeyData 私钥字节数组
	CertificateRequestFilePath string                  // certificateFilePath 指定生成的证书签名请求文件路径，如'/etc/rootCA.csr'
	SignatureAlgorithm         x509.SignatureAlgorithm // signatureAlgorithm 生成证书时候采用的签名算法
	Subject                    pkix.Name               // Subject 签名信息
	DNSNames                   []string                // DNSNames DNS限制
	EmailAddresses             []string                // EmailAddresses 邮箱地址限制
	IPAddresses                []net.IP                // IPAddresses IP地址限制
	URIs                       []*url.URL              // URIs URL地址限制
}

func (cert *CertRequest) init() *CertRequestModel {
	temp := &x509.CertificateRequest{
		SignatureAlgorithm: cert.SignatureAlgorithm,
		Subject:            cert.Subject,
		DNSNames:           cert.DNSNames,
		EmailAddresses:     cert.EmailAddresses,
		IPAddresses:        cert.IPAddresses,
		URIs:               cert.URIs,
	}
	return &CertRequestModel{
		CertificateRequestFilePath: cert.CertificateRequestFilePath,
		Template:                   temp,
	}
}

// CertRequestFP 证书生成请求对象
type CertRequestFP struct {
	PrivateKeyFilePath         string                  // privateKeyFilePath 私钥文件存储路径
	CertificateRequestFilePath string                  // certificateFilePath 指定生成的证书签名请求文件路径，如'/etc/rootCA.csr'
	SignatureAlgorithm         x509.SignatureAlgorithm // signatureAlgorithm 生成证书时候采用的签名算法
	Subject                    pkix.Name               // Subject 签名信息
	DNSNames                   []string                // DNSNames DNS限制
	EmailAddresses             []string                // EmailAddresses 邮箱地址限制
	IPAddresses                []net.IP                // IPAddresses IP地址限制
	URIs                       []*url.URL              // URIs URL地址限制
}

func (cert *CertRequestFP) init() *CertRequestModel {
	temp := &x509.CertificateRequest{
		SignatureAlgorithm: cert.SignatureAlgorithm,
		Subject:            cert.Subject,
		DNSNames:           cert.DNSNames,
		EmailAddresses:     cert.EmailAddresses,
		IPAddresses:        cert.IPAddresses,
		URIs:               cert.URIs,
	}
	return &CertRequestModel{
		CertificateRequestFilePath: cert.CertificateRequestFilePath,
		Template:                   temp,
	}
}

// CertRequestModel 证书生成请求对象
type CertRequestModel struct {
	PrivateKey                 crypto.Signer            // privateKey 私钥字节数组
	CertificateRequestFilePath string                   // certificateFilePath 指定生成的证书签名请求文件路径，如'/etc/rootCA.csr'
	Template                   *x509.CertificateRequest // template 生成证书时候采用的请求模板
}
