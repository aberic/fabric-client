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
	"sync"
)

var (
	bc          *ByteCommon
	cc          *CommandCommon
	ec          *EnvCommon
	fc          *FileCommon
	ic          *IPCommon
	jc          *JWTCommon
	sc          *StringCommon
	hc          *HashCommon
	rc          *RSACommon
	ac          *AESCommon
	dc          *DESCommon
	ecc         *ECCCommon
	cac         *CACommon
	lc          *LogCommon
	scc         *ScaleCommon
	tc          *TimeCommon
	onceByte    sync.Once
	onceCommand sync.Once
	onceEnv     sync.Once
	onceFile    sync.Once
	onceIP      sync.Once
	onceJwt     sync.Once
	onceString  sync.Once
	onceHash    sync.Once
	onceRSA     sync.Once
	onceAES     sync.Once
	onceDES     sync.Once
	onceECC     sync.Once
	onceCAC     sync.Once
	onceLog     sync.Once
	onceScale   sync.Once
	onceTime    sync.Once
)

// Byte 字节工具
func Byte() *ByteCommon {
	onceByte.Do(func() {
		bc = &ByteCommon{}
	})
	return bc
}

// Command 命令行工具
func Command() *CommandCommon {
	onceCommand.Do(func() {
		cc = &CommandCommon{}
	})
	return cc
}

// Env 环境变量工具
func Env() *EnvCommon {
	onceEnv.Do(func() {
		ec = &EnvCommon{}
	})
	return ec
}

// File 文件操作工具
func File() *FileCommon {
	onceFile.Do(func() {
		fc = &FileCommon{}
	})
	return fc
}

// IP IP工具
func IP() *IPCommon {
	onceIP.Do(func() {
		ic = &IPCommon{}
	})
	return ic
}

// JWT JWT工具
func JWT() *JWTCommon {
	onceJwt.Do(func() {
		jc = &JWTCommon{}
	})
	return jc
}

// String 字符串工具
func String() *StringCommon {
	onceString.Do(func() {
		sc = &StringCommon{}
	})
	return sc
}

// CryptoHash Hash/散列工具
func CryptoHash() *HashCommon {
	onceHash.Do(func() {
		hc = &HashCommon{}
	})
	return hc
}

// CryptoRSA RSA工具
func CryptoRSA() *RSACommon {
	onceRSA.Do(func() {
		rc = &RSACommon{}
	})
	return rc
}

// CryptoAES AES工具
func CryptoAES() *AESCommon {
	onceAES.Do(func() {
		ac = &AESCommon{}
	})
	return ac
}

// CryptoDES DES工具
func CryptoDES() *DESCommon {
	onceDES.Do(func() {
		dc = &DESCommon{}
	})
	return dc
}

// CryptoECC ECC工具
func CryptoECC() *ECCCommon {
	onceECC.Do(func() {
		ecc = &ECCCommon{}
	})
	return ecc
}

// CA CA工具
func CA() *CACommon {
	onceCAC.Do(func() {
		cac = &CACommon{}
	})
	return cac
}

// Log 日志工具
func Log() *LogCommon {
	onceLog.Do(func() {
		lc = &LogCommon{level: debugLevel, production: false}
	})
	return lc
}

// Scale 算数/转换工具
func Scale() *ScaleCommon {
	onceScale.Do(func() {
		scc = &ScaleCommon{}
	})
	return scc
}

// Time 时间工具
func Time() *TimeCommon {
	onceTime.Do(func() {
		tc = &TimeCommon{}
	})
	return tc
}
