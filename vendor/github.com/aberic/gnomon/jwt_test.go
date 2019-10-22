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
 */

package gnomon

import (
	"testing"
	"time"
)

func TestJwtCommon_Build(t *testing.T) {
	key := []byte("Hello World！This is secret!")
	tokenString1, _ := JWT().Build(signingMethodHS256, key, "1", "rivet", "userMD5", time.Now().Unix(), time.Now().Unix(), time.Now().Unix()+1000)
	t.Log("tokenString1", tokenString1)
	tokenString2, _ := JWT().Build(signingMethodHS384, key, "1", "rivet", "userMD5", time.Now().Unix(), time.Now().Unix(), time.Now().Unix()+1000)
	t.Log("tokenString2", tokenString2)
	tokenString3, _ := JWT().Build(signingMethodHS512, key, "1", "rivet", "userMD5", time.Now().Unix(), time.Now().Unix(), time.Now().Unix()+1000)
	t.Log("tokenString3", tokenString3)
}

func TestJwtCommon_Check(t *testing.T) {
	key := []byte("Hello World！This is secret!")
	tokenString1, _ := JWT().Build(signingMethodHS256, key, "1", "rivet", "userMD5", time.Now().Unix(), time.Now().Unix(), time.Now().Unix()+1000)
	t.Log("tokenString1", tokenString1)
	bo1 := JWT().Check(key, tokenString1)
	t.Log("bo1", bo1)
	tokenString2, _ := JWT().Build(signingMethodHS384, key, "1", "rivet", "userMD5", time.Now().Unix(), time.Now().Unix(), time.Now().Unix()+1000)
	t.Log("tokenString2", tokenString2)
	bo2 := JWT().Check(key, tokenString2)
	t.Log("bo3", bo2)
	tokenString3, _ := JWT().Build(signingMethodHS512, key, "1", "rivet", "userMD5", time.Now().Unix(), time.Now().Unix(), time.Now().Unix()+1000)
	t.Log("tokenString3", tokenString3)
	bo3 := JWT().Check(key, tokenString3)
	t.Log("bo3", bo3)

	bo4 := JWT().Check(key, tokenString3+"1")
	t.Log("bo4", bo4)
}
