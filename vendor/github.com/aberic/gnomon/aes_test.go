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
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestCryptoAES(t *testing.T) {
	data := []byte("Hello World")     // 待加密的数据
	key := []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
	t.Log("原文：", string(data))

	t.Log("------------------ CBC模式 --------------------")
	encrypted := CryptoAES().EncryptCBC(data, key)
	t.Log("密文(hex)：", hex.EncodeToString(encrypted))
	t.Log("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := CryptoAES().DecryptCBC(encrypted, key)
	t.Log("解密结果：", string(decrypted))

	t.Log("------------------ ECB模式 --------------------")
	encrypted = CryptoAES().EncryptECB(data, key)
	t.Log("密文(hex)：", hex.EncodeToString(encrypted))
	t.Log("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = CryptoAES().DecryptECB(encrypted, key)
	t.Log("解密结果：", string(decrypted))

	t.Log("------------------ CFB模式 --------------------")
	encrypted = CryptoAES().EncryptCFB(data, key)
	t.Log("密文(hex)：", hex.EncodeToString(encrypted))
	t.Log("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = CryptoAES().DecryptCFB(encrypted, key)
	t.Log("解密结果：", string(decrypted))

	t.Log("------------------ OFB模式 --------------------")
	encrypted = CryptoAES().EncryptOFB(data, key)
	t.Log("密文(hex)：", hex.EncodeToString(encrypted))
	t.Log("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = CryptoAES().DecryptOFB(encrypted, key)
	t.Log("解密结果：", string(decrypted))
}
