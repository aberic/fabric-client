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
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// HashCommon hash/散列工具
type HashCommon struct{}

// MD5 MD5
func (h *HashCommon) MD5(text string) string {
	hash := md5.New()
	_, _ = hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// MD516 MD516
func (h *HashCommon) MD516(text string) string {
	md516 := string([]rune(h.MD5(text))[8:24])
	return md516
}

// Sha1 Sha1
func (h *HashCommon) Sha1(text string) string {
	hash := sha1.New()
	_, _ = hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha224 Sha224
func (h *HashCommon) Sha224(text string) string {
	hash := crypto.SHA224.New()
	_, _ = hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha256 Sha256
func (h *HashCommon) Sha256(text string) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha384 Sha384
func (h *HashCommon) Sha384(text string) string {
	hash := sha512.New384()
	_, _ = hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha512 Sha512
func (h *HashCommon) Sha512(text string) string {
	hash := sha512.New()
	_, _ = hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
