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
	"crypto/cipher"
	"crypto/des"
)

// DESCommon DES工具
type DESCommon struct{}

//--------------------------------------------------------------------------------------------------------------------

// EncryptCBC CBC加密
//
// data 待加密数据
//
// key 自定义密钥，如：'[]byte("12345678")'，长度必须是8位
func (d *DESCommon) EncryptCBC(data, key []byte) []byte {
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	data = d.pkcs5Padding(data, block.BlockSize())
	//获取CBC加密模式
	iv := key //用密钥作为向量(不建议这样使用)
	mode := cipher.NewCBCEncrypter(block, iv)
	out := make([]byte, len(data))
	mode.CryptBlocks(out, data)
	return out
}

// DecryptCBC CBC解密
//
// data 待加密数据
//
// key 自定义密钥，如：'[]byte("12345678")'，长度必须是8位
func (d *DESCommon) DecryptCBC(data, key []byte) []byte {
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := key //用密钥作为向量(不建议这样使用)
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(data))
	mode.CryptBlocks(plaintext, data)
	plaintext = d.pkcs5UnPadding(plaintext)
	return plaintext
}

//--------------------------------------------------------------------------------------------------------------------

// EncryptECB ECB加密
//
// data 待加密数据
//
// key 自定义密钥，如：'[]byte("12345678")'，长度必须是8位
func (d *DESCommon) EncryptECB(data, key []byte) []byte {
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	bs := block.BlockSize()
	//对明文数据进行补码
	data = d.pkcs5Padding(data, bs)
	if len(data)%bs != 0 {
		panic("Need a multiple of the block size")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//对明文按照blocksize进行分块加密
		//必要时可以使用go关键字进行并行加密
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out
}

// DecryptECB ECB解密
//
// data 待加密数据
//
// key 自定义密钥，如：'[]byte("12345678")'，长度必须是8位
func (d *DESCommon) DecryptECB(data, key []byte) []byte {
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = d.pkcs5UnPadding(out)
	return out
}

//--------------------------------------------------------------------------------------------------------------------

// pkcs5Padding 明文补码算法
func (d *DESCommon) pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// pkcs5UnPadding 明文减码算法
func (d *DESCommon) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
