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
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var (
	contentRSA       = "this is a test"
	pathrsapksc1256  = "./tmp/example/rsa/pksc1/256"
	pathrsapksc1512  = "./tmp/example/rsa/pksc1/512"
	pathrsapksc11024 = "./tmp/example/rsa/pksc1/1024"
	pathrsapksc12048 = "./tmp/example/rsa/pksc1/2048"
	pathrsapksc8256  = "./tmp/example/rsa/pksc8/256"
	pathrsapksc8512  = "./tmp/example/rsa/pksc8/512"
	pathrsapksc81024 = "./tmp/example/rsa/pksc8/1024"
	pathrsapksc82048 = "./tmp/example/rsa/pksc8/2048"
	privateRSAName   = "private.pem"
	publicRSAName    = "public.pem"

	priRSAKey *rsa.PrivateKey
	//pubRSAKey *rsa.PublicKey

	priRSAKeyData []byte
	pubRSAKeyData []byte
	dataRSA       []byte
	dataRSAEncode []byte
	signRSAResult []byte
	errRSA        error
)

func TestRSACommon_GenerateRsaKey(t *testing.T) {
	t.Log(CryptoRSA().GenerateKey(256, pathrsapksc1256, privateRSAName, publicRSAName, CryptoRSA().PKSC1()))
	t.Log(CryptoRSA().GenerateKeyWithPass(512, pathrsapksc1512, privateRSAName, publicRSAName, "123456", x509.PEMCipher3DES, CryptoRSA().PKSC1()))
	t.Log(CryptoRSA().GenerateKey(1024, pathrsapksc11024, privateRSAName, publicRSAName, CryptoRSA().PKSC1()))
	t.Log(CryptoRSA().GenerateKeyWithPass(2048, pathrsapksc12048, privateRSAName, publicRSAName, "123456", x509.PEMCipher3DES, CryptoRSA().PKSC1()))

	t.Log(CryptoRSA().GenerateKeyWithPass(256, pathrsapksc8256, privateRSAName, publicRSAName, "123456", x509.PEMCipher3DES, CryptoRSA().PKSC8()))
	t.Log(CryptoRSA().GenerateKey(512, pathrsapksc8512, privateRSAName, publicRSAName, CryptoRSA().PKSC8()))
	t.Log(CryptoRSA().GenerateKeyWithPass(1024, pathrsapksc81024, privateRSAName, publicRSAName, "123456", x509.PEMCipher3DES, CryptoRSA().PKSC8()))
	t.Log(CryptoRSA().GenerateKey(2048, pathrsapksc82048, privateRSAName, publicRSAName, CryptoRSA().PKSC8()))
}

func TestRSACommon_GenerateRsaCustomPriKey(t *testing.T) {
	t.Log(CryptoRSA().GeneratePriKeyWithPass(256, pathrsapksc1256, "private1.pem", "123456", x509.PEMCipher3DES, CryptoRSA().PKSC1()))
	t.Log(CryptoRSA().GeneratePriKey(512, pathrsapksc1512, "private1.pem", CryptoRSA().PKSC1()))
	t.Log(CryptoRSA().GeneratePriKeyWithPass(1024, pathrsapksc11024, "private1.pem", "123456", -1, CryptoRSA().PKSC1()))
	t.Log(CryptoRSA().GeneratePriKeyWithPass(1024, pathrsapksc11024, "private1.pem", "123456", x509.PEMCipher3DES, CryptoRSA().PKSC1()))
	t.Log(CryptoRSA().GeneratePriKey(2048, pathrsapksc12048, "private1.pem", CryptoRSA().PKSC1()))
	t.Log()

	t.Log(CryptoRSA().GeneratePriKey(256, pathrsapksc8256, "private1.pem", CryptoRSA().PKSC8()))
	t.Log(CryptoRSA().GeneratePriKeyWithPass(512, pathrsapksc8512, "private1.pem", "123456", x509.PEMCipher3DES, CryptoRSA().PKSC8()))
	t.Log(CryptoRSA().GeneratePriKey(1024, pathrsapksc81024, "private1.pem", CryptoRSA().PKSC8()))
	t.Log(CryptoRSA().GeneratePriKeyWithPass(2048, pathrsapksc82048, "private1.pem", "123456", -1, CryptoRSA().PKSC8()))
	t.Log(CryptoRSA().GeneratePriKeyWithPass(2048, pathrsapksc82048, "private1.pem", "123456", x509.PEMCipher3DES, CryptoRSA().PKSC8()))
	t.Log(CryptoRSA().GeneratePubKey(nil, "/etc/pub", "public1.pem", CryptoRSA().PKSC8()))
}

func TestRSACommon_GenerateRsaCustomPubKey(t *testing.T) {
	t.Log(CryptoRSA().GeneratePubKeyFPWithPass(pathrsapksc1256+"/"+"private1.pem", "123456", pathrsapksc1256, "public1.pem", CryptoRSA().PKSC1()))
	t.Log(CryptoRSA().GeneratePubKeyFP(pathrsapksc1512+"/"+"private1.pem", pathrsapksc1512, "public1.pem", CryptoRSA().PKSC1()))
	t.Log(CryptoRSA().GeneratePubKeyFPWithPass(pathrsapksc11024+"/"+"private1.pem", "123456", pathrsapksc11024, "public1.pem", CryptoRSA().PKSC1()))
	t.Log(CryptoRSA().GeneratePubKeyFP(pathrsapksc12048+"/"+"private1.pem", pathrsapksc12048, "public1.pem", CryptoRSA().PKSC1()))
	t.Log()

	t.Log(CryptoRSA().GeneratePubKeyFP(pathrsapksc8256+"/"+"private1.pem", pathrsapksc8256, "public1.pem", CryptoRSA().PKSC8()))
	t.Log(CryptoRSA().GeneratePubKeyFPWithPass(pathrsapksc8512+"/"+"private1.pem", "123456", pathrsapksc8512, "public1.pem", CryptoRSA().PKSC8()))
	t.Log(CryptoRSA().GeneratePubKeyFP(pathrsapksc81024+"/"+"private1.pem", pathrsapksc81024, "public1.pem", CryptoRSA().PKSC8()))
	t.Log(CryptoRSA().GeneratePubKeyFPWithPass(pathrsapksc82048+"/"+"private1.pem", "123456", pathrsapksc82048, "public1.pem", CryptoRSA().PKSC8()))
	t.Log()

	t.Log(CryptoRSA().GeneratePubKeyFP(pathrsapksc12048+"/"+"private100.pem", pathrsapksc12048, "public1.pem", CryptoRSA().PKSC1()))
	t.Log(CryptoRSA().GeneratePubKeyFPWithPass(pathrsapksc82048+"/"+"private100.pem", "123456", pathrsapksc82048, "public1.pem", CryptoRSA().PKSC8()))
	t.Log(CryptoRSA().GeneratePubKeyFP(pathrsapksc82048+"/"+"private1.pem", pathrsapksc82048, "public2.pem", CryptoRSA().PKSC1()))
}

func TestRSACommon_GenerateRsaKey_FailPathExists(t *testing.T) {
	t.Log(CryptoRSA().GenerateKey(256, "/etc/test", privateRSAName, publicRSAName, CryptoRSA().PKSC1()))
}

func TestRSACommon_GenerateRsaKey_FailGenerate(t *testing.T) {
	t.Log(CryptoRSA().GenerateKey(-1, pathrsapksc1256, privateRSAName, publicRSAName, CryptoRSA().PKSC1()))
}

func TestRSACommon_GenerateRsaKey_FailCreate(t *testing.T) {
	t.Log(CryptoRSA().GenerateKey(256, "/etc", privateRSAName, publicRSAName, CryptoRSA().PKSC1()))
}

func TestRSACommon_GenerateRsaPKSC8Key_FailPathExists(t *testing.T) {
	t.Log(CryptoRSA().GenerateKey(256, "/etc/test", privateRSAName, publicRSAName, CryptoRSA().PKSC8()))
}

func TestRSACommon_GenerateRsaPKSC8Key_FailGenerate(t *testing.T) {
	t.Log(CryptoRSA().GenerateKey(-1, pathrsapksc1256, privateRSAName, publicRSAName, CryptoRSA().PKSC8()))
}

func TestRSACommon_GenerateRsaPKSC8Key_FailCreate(t *testing.T) {
	t.Log(CryptoRSA().GenerateKey(256, "/etc", privateRSAName, publicRSAName, CryptoRSA().PKSC8()))
}

func TestRSACommon_RsaEncryptDecrypt_Fail(t *testing.T) {
	dataRSA, errRSA = CryptoRSA().Encrypt([]byte{}, []byte(contentRSA))
	if nil != errRSA {
		t.Log(errRSA)
	}

	dataRSA, errRSA = CryptoRSA().EncryptFP(pathrsapksc1256+"/a/"+publicRSAName, []byte(contentRSA))
	if nil != errRSA {
		t.Log(errRSA)
	}

	dataRSA, errRSA = CryptoRSA().Decrypt([]byte{}, []byte(contentRSA), CryptoRSA().PKSC8())
	if nil != errRSA {
		t.Log(errRSA)
	}
}

func TestRSACommon_RsaEncryptDecrypt(t *testing.T) {
	t.Log("加密前：", contentRSA)
	t.Log("=================================")

	pubRSAKeyData, errRSA = ioutil.ReadFile(filepath.Join(pathrsapksc1256, publicRSAName))
	if nil != errRSA {
		t.Error(errRSA)
	}
	dataRSA, errRSA = CryptoRSA().Encrypt(pubRSAKeyData, []byte(contentRSA))
	if nil != errRSA {
		t.Error(errRSA)
	}
	dataRSA, errRSA = CryptoRSA().EncryptFP(pathrsapksc1256+"/"+publicRSAName, []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后256：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = CryptoRSA().DecryptFP(pathrsapksc1256+"/"+privateRSAName, dataRSA, CryptoRSA().PKSC1())
	t.Log("解密后256：", string(dataRSAEncode))
	priRSAKeyData, errRSA = ioutil.ReadFile(filepath.Join(pathrsapksc1256, privateRSAName))
	if nil != errRSA {
		t.Error(errRSA)
	}
	dataRSAEncode, errRSA = CryptoRSA().Decrypt(priRSAKeyData, dataRSA, CryptoRSA().PKSC1())
	if nil != errRSA {
		t.Log(errRSA)
	}
	t.Log("解密后256：", string(dataRSAEncode))
	t.Log("=================================")

	dataRSA, errRSA = CryptoRSA().EncryptFP(pathrsapksc1512+"/"+publicRSAName, []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后512：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = CryptoRSA().DecryptFP(pathrsapksc1512+"/"+privateRSAName, dataRSA, CryptoRSA().PKSC1())
	t.Log("解密后512：", string(dataRSAEncode))
	t.Log("=================================")

	dataRSA, errRSA = CryptoRSA().EncryptFP(pathrsapksc11024+"/"+publicRSAName, []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后1024：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = CryptoRSA().DecryptFP(pathrsapksc11024+"/"+privateRSAName, dataRSA, CryptoRSA().PKSC1())
	t.Log("解密后1024：", string(dataRSAEncode))
	t.Log("=================================")

	dataRSA, errRSA = CryptoRSA().EncryptFP(pathrsapksc12048+"/"+publicRSAName, []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后2048：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = CryptoRSA().DecryptFP(pathrsapksc12048+"/"+privateRSAName, dataRSA, CryptoRSA().PKSC1())
	t.Log("解密后2048：", string(dataRSAEncode))
}

func TestRSACommon_RsaPKSC8EncryptDecrypt(t *testing.T) {
	t.Log("加密前：", contentRSA)
	t.Log("=================================")
	dataRSA, errRSA = CryptoRSA().EncryptFP(pathrsapksc8256+"/"+publicRSAName, []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后256：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = CryptoRSA().DecryptFP(pathrsapksc8256+"/"+privateRSAName, dataRSA, CryptoRSA().PKSC8())
	t.Log("解密后256：", string(dataRSAEncode))
	t.Log("=================================")

	dataRSA, errRSA = CryptoRSA().EncryptFP(pathrsapksc8512+"/"+publicRSAName, []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后512：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = CryptoRSA().DecryptFP(pathrsapksc8512+"/"+privateRSAName, dataRSA, CryptoRSA().PKSC8())
	t.Log("解密后512：", string(dataRSAEncode))
	t.Log("=================================")

	dataRSA, errRSA = CryptoRSA().EncryptFP(pathrsapksc81024+"/"+publicRSAName, []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后1024：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = CryptoRSA().DecryptFP(pathrsapksc81024+"/"+privateRSAName, dataRSA, CryptoRSA().PKSC8())
	t.Log("解密后1024：", string(dataRSAEncode))
	t.Log("=================================")

	dataRSA, errRSA = CryptoRSA().EncryptFP(pathrsapksc82048+"/"+publicRSAName, []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后2048：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = CryptoRSA().DecryptFP(pathrsapksc82048+"/"+privateRSAName, dataRSA, CryptoRSA().PKSC8())
	t.Log("解密后2048：", string(dataRSAEncode))
}

func TestRSACommon_RsaSign(t *testing.T) {
	t.Log("签名：", contentRSA)
	t.Log("=================================")
	priRSAKeyData, errRSA = ioutil.ReadFile(filepath.Join(pathrsapksc11024, privateRSAName))
	if nil != errRSA {
		t.Error(errRSA)
	}
	if signRSAResult, errRSA = CryptoRSA().Sign(priRSAKeyData, []byte(contentRSA), crypto.SHA256, CryptoRSA().PKSC1(), CryptoRSA().SignPKCS()); nil != errRSA {
		t.Skip("签名错误1024：", errRSA)
	} else {
		t.Log("验签1024：", signRSAResult)
		pubRSAKeyData, errRSA = ioutil.ReadFile(filepath.Join(pathrsapksc11024, publicRSAName))
		if nil != errRSA {
			t.Error(errRSA)
		}
		if errRSA = CryptoRSA().Verify(pubRSAKeyData, []byte(contentRSA), signRSAResult, crypto.SHA256, CryptoRSA().SignPKCS()); nil != errRSA {
			t.Skip("验签错误1024：", errRSA)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFPWithPass(pathrsapksc1512+"/"+privateRSAName, "123456", []byte(contentRSA), crypto.SHA256, CryptoRSA().PKSC1(), CryptoRSA().SignPKCS()); nil != errRSA {
		t.Skip("签名错误512：", errRSA)
	} else {
		t.Log("验签512：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc1512+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA256, CryptoRSA().SignPKCS()); nil != errRSA {
			t.Skip("验签错误512：", errRSA)
		} else {
			t.Log("验签通过512")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc11024+"/"+privateRSAName, []byte(contentRSA), crypto.SHA512, CryptoRSA().PKSC1(), CryptoRSA().SignPKCS()); nil != errRSA {
		t.Skip("签名错误1024：", errRSA)
	} else {
		t.Log("验签1024：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc11024+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA512, CryptoRSA().SignPKCS()); nil != errRSA {
			t.Skip("验签错误1024：", errRSA)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc12048+"/"+privateRSAName, []byte(contentRSA), crypto.SHA384, CryptoRSA().PKSC1(), CryptoRSA().SignPKCS()); nil != errRSA {
		t.Skip("签名错误2048：", errRSA)
	} else {
		t.Log("验签2048：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc12048+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA384, CryptoRSA().SignPKCS()); nil != errRSA {
			t.Skip("验签错误2048：", errRSA)
		} else {
			t.Log("验签通过2048")
		}
	}
}

func TestRSACommon_RsaSignPSS(t *testing.T) {
	t.Log("签名：", contentRSA)
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc1512+"/"+privateRSAName, []byte(contentRSA), crypto.SHA256, CryptoRSA().PKSC1(), CryptoRSA().SignPSS()); nil != errRSA {
		t.Skip("签名错误512：", errRSA)
	} else {
		t.Log("验签512：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc1512+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA256, CryptoRSA().SignPSS()); nil != errRSA {
			t.Skip("验签错误512：", errRSA)
		} else {
			t.Log("验签通过512")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc11024+"/"+privateRSAName, []byte(contentRSA), crypto.SHA512, CryptoRSA().PKSC1(), CryptoRSA().SignPSS()); nil != errRSA {
		t.Skip("签名错误1024：", errRSA)
	} else {
		t.Log("验签1024：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc11024+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA512, CryptoRSA().SignPSS()); nil != errRSA {
			t.Skip("验签错误1024：", errRSA)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc12048+"/"+privateRSAName, []byte(contentRSA), crypto.SHA384, CryptoRSA().PKSC1(), CryptoRSA().SignPSS()); nil != errRSA {
		t.Skip("签名错误2048：", errRSA)
	} else {
		t.Log("验签2048：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc12048+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA384, CryptoRSA().SignPSS()); nil != errRSA {
			t.Skip("验签错误2048：", errRSA)
		} else {
			t.Log("验签通过2048")
		}
	}
}

func TestRSACommon_RsaPKSC8Sign(t *testing.T) {
	t.Log("签名：", contentRSA)
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc8512+"/"+privateRSAName, []byte(contentRSA), crypto.SHA256, CryptoRSA().PKSC8(), CryptoRSA().SignPKCS()); nil != errRSA {
		t.Skip("签名错误512：", errRSA)
	} else {
		t.Log("验签512：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc8512+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA256, CryptoRSA().SignPKCS()); nil != errRSA {
			t.Skip("验签错误512：", errRSA)
		} else {
			t.Log("验签通过512")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc81024+"/"+privateRSAName, []byte(contentRSA), crypto.SHA512, CryptoRSA().PKSC8(), CryptoRSA().SignPKCS()); nil != errRSA {
		t.Skip("签名错误1024：", errRSA)
	} else {
		t.Log("验签1024：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc81024+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA512, CryptoRSA().SignPKCS()); nil != errRSA {
			t.Skip("验签错误1024：", errRSA)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc82048+"/"+privateRSAName, []byte(contentRSA), crypto.SHA384, CryptoRSA().PKSC8(), CryptoRSA().SignPKCS()); nil != errRSA {
		t.Skip("签名错误2048：", errRSA)
	} else {
		t.Log("验签2048：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc82048+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA384, CryptoRSA().SignPKCS()); nil != errRSA {
			t.Skip("验签错误2048：", errRSA)
		} else {
			t.Log("验签通过2048")
		}
	}
}

func TestRSACommon_RsaPKSC8SignPSS(t *testing.T) {
	t.Log("签名：", contentRSA)
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc8512+"/"+privateRSAName, []byte(contentRSA), crypto.SHA256, CryptoRSA().PKSC8(), CryptoRSA().SignPSS()); nil != errRSA {
		t.Skip("签名错误512：", errRSA)
	} else {
		t.Log("验签512：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc8512+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA256, CryptoRSA().SignPSS()); nil != errRSA {
			t.Skip("验签错误512：", errRSA)
		} else {
			t.Log("验签通过512")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc81024+"/"+privateRSAName, []byte(contentRSA), crypto.SHA512, CryptoRSA().PKSC8(), CryptoRSA().SignPSS()); nil != errRSA {
		t.Skip("签名错误1024：", errRSA)
	} else {
		t.Log("验签1024：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc81024+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA512, CryptoRSA().SignPSS()); nil != errRSA {
			t.Skip("验签错误1024：", errRSA)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc82048+"/"+privateRSAName, []byte(contentRSA), crypto.SHA384, CryptoRSA().PKSC8(), CryptoRSA().SignPSS()); nil != errRSA {
		t.Skip("签名错误2048：", errRSA)
	} else {
		t.Log("验签2048：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc82048+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA384, CryptoRSA().SignPSS()); nil != errRSA {
			t.Skip("验签错误2048：", errRSA)
		} else {
			t.Log("验签通过2048")
		}
	}
}

func TestRSACommon_RsaSign_Fail(t *testing.T) {
	t.Log("签名：", contentRSA)
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc1256+"/"+privateRSAName, []byte(contentRSA), crypto.SHA384, CryptoRSA().PKSC1(), CryptoRSA().SignPKCS()); nil != errRSA {
		t.Skip("签名错误256：", errRSA)
	} else {
		t.Log("验签256：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc1256+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA384, CryptoRSA().SignPKCS()); nil != errRSA {
			t.Skip("验签错误256：", errRSA)
		} else {
			t.Log("验签通过256")
		}
	}
}

func TestRSACommon_RsaPKSC8Sign_Fail(t *testing.T) {
	t.Log("签名：", contentRSA)
	t.Log("=================================")
	if signRSAResult, errRSA = CryptoRSA().SignFP(pathrsapksc8256+"/"+privateRSAName, []byte(contentRSA), crypto.SHA384, CryptoRSA().PKSC8(), CryptoRSA().SignPKCS()); nil != errRSA {
		t.Skip("签名错误256：", errRSA)
	} else {
		t.Log("验签256：", signRSAResult)
		if errRSA = CryptoRSA().VerifyFP(pathrsapksc8256+"/"+publicRSAName, []byte(contentRSA), signRSAResult, crypto.SHA384, CryptoRSA().SignPKCS()); nil != errRSA {
			t.Skip("验签错误256：", errRSA)
		} else {
			t.Log("验签通过256")
		}
	}
}
