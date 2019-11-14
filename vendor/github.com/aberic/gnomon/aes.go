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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// AESCommon AES工具
//
// 对称加密的四种模式(ECB、CBC、CFB、OFB)
//
// ECB模式——电码本模式（Electronic Codebook Book (ECB)）
//
// 优点:1.简单;2.有利于并行计算;3.误差不会被传送；
//
// 缺点:1.不能隐藏明文的模式;2.可能对明文进行主动攻击；
//
// ======================================
//
// CBC模式——密码分组链接模式（Cipher Block Chaining (CBC)）
//
// 优点:1.不容易主动攻击,安全性好于ECB,适合传输长度长的报文,是SSL、IPSec的标准。
//
// 缺点:1.不利于并行计算;2.误差传递;3.需要初始化向量IV
//
// ======================================
//
// CFB模式——密码反馈模式（Cipher FeedBack (CFB)）
//
// 优点:1.隐藏了明文模式;2.分组密码转化为流模式;3.可以及时加密传送小于分组的数据;
//
// 缺点:1.不利于并行计算;2.误差传送:一个明文单元损坏影响多个单元;3.唯一的IV;
//
// ======================================
//
// OFB模式——输出反馈模式（Output FeedBack (OFB)）
//
// 优点:1.隐藏了明文模式;2.分组密码转化为流模式;3.可以及时加密传送小于分组的数据;
//
// 缺点:1.不利于并行计算;2.对明文的主动攻击是可能的;3.误差传送：一个明文单元损坏影响多个单元;
type AESCommon struct{}

//--------------------------------------------------------------------------------------------------------------------

// EncryptCBC AES CBC模式加密
func (a *AESCommon) EncryptCBC(data []byte, key []byte) (encrypted []byte) {
	// 分组密钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	blockMode, blockSize := a.blockMode(key, true) // 加密模式
	data = a.pkcs5Padding(data, blockSize)         // 补全码
	encrypted = make([]byte, len(data))            // 创建数组
	blockMode.CryptBlocks(encrypted, data)         // 加密
	return encrypted
}

// DecryptCBC AES CBC模式解密
func (a *AESCommon) DecryptCBC(encrypted []byte, key []byte) []byte {
	blockMode, _ := a.blockMode(key, false)     // 加密模式
	decrypted := make([]byte, len(encrypted))   // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted) // 解密
	decrypted = a.pkcs5UnPadding(decrypted)     // 去除补全码
	return decrypted
}

// blockMode 加密模式
func (a *AESCommon) blockMode(key []byte, encrypt bool) (cipher.BlockMode, int) {
	block, _ := aes.NewCipher(key) // 分组密钥
	blockSize := block.BlockSize() // 获取密钥块的长度
	if encrypt {
		return cipher.NewCBCEncrypter(block, key[:blockSize]), blockSize // 加密模式
	}
	return cipher.NewCBCDecrypter(block, key[:blockSize]), blockSize // 加密模式
}

// pkcs5Padding 明文补码算法
func (a *AESCommon) pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// pkcs5UnPadding 明文减码算法
func (a *AESCommon) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

//--------------------------------------------------------------------------------------------------------------------

// EncryptECB AES ECB模式加密
func (a *AESCommon) EncryptECB(data []byte, key []byte) (encrypted []byte) {
	cpr, _ := aes.NewCipher(a.generateKey(key))
	length := (len(data) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, data)
	pad := byte(len(plain) - len(data))
	for i := len(data); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cpr.BlockSize(); bs <= len(data); bs, be = bs+cpr.BlockSize(), be+cpr.BlockSize() {
		cpr.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted
}

// DecryptECB AES ECB模式解密
func (a *AESCommon) DecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cpr, _ := aes.NewCipher(a.generateKey(key))
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, cpr.BlockSize(); bs < len(encrypted); bs, be = bs+cpr.BlockSize(), be+cpr.BlockSize() {
		cpr.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim]
}

func (a *AESCommon) generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

//--------------------------------------------------------------------------------------------------------------------

// EncryptCFB AES CFB模式加密
func (a *AESCommon) EncryptCFB(data []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(data))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], data)
	return encrypted
}

// DecryptCFB AES CFB模式解密
func (a *AESCommon) DecryptCFB(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("cipher text too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}

//--------------------------------------------------------------------------------------------------------------------

// EncryptOFB AES OFB模式加密
func (a *AESCommon) EncryptOFB(data []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(data))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], data)
	return encrypted
}

// DecryptOFB AES OFB模式解密
func (a *AESCommon) DecryptOFB(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}
