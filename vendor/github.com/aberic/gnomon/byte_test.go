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

import "testing"

type bTest struct {
	name string
	age  uint8
	male bool
}

type BTest struct {
	Name string
	Age  uint8
	Male bool
}

func TestByteCommon_GetBytes(t *testing.T) {
	if data, err := Byte().GetBytes(&bTest{name: "test", age: 18, male: true}); nil != err {
		t.Log(err)
	} else {
		t.Log(data)
	}
	if data, err := Byte().GetBytes(&BTest{Name: "test", Age: 18, Male: true}); nil != err {
		t.Log(err)
	} else {
		t.Log(data)
	}
	if data, err := Byte().GetBytes(100); nil != err { // [4 4 0 255 200]
		t.Log(err)
	} else {
		t.Log(data)
	}
	if data, err := Byte().GetBytes(true); nil != err {
		t.Log(err)
	} else {
		t.Log(data)
	}
	if data, err := Byte().GetBytes("100"); nil != err {
		t.Log(err)
	} else {
		t.Log(data)
	}
}

func TestByteCommon_IntToBytes(t *testing.T) {
	if data, err := Byte().IntToBytes(100); nil != err { // [0 0 0 100]
		t.Log(err)
	} else {
		t.Log(data)
	}
}

func TestByteCommon_BytesToInt(t *testing.T) {
	if data, err := Byte().IntToBytes(100); nil != err { // [0 0 0 100]
		t.Log(err)
	} else {
		t.Log(data)
		if dataInt, err := Byte().BytesToInt(data); nil != err { // [0 0 0 100]
			t.Log(err)
		} else {
			t.Log(dataInt)

		}
	}
}

func TestByteCommon_Uint16ToBytes(t *testing.T) {
	t.Log(Byte().Uint16ToBytes(100))
}

func TestByteCommon_BytesToUint16(t *testing.T) {
	data := Byte().Uint16ToBytes(100)
	t.Log(Byte().BytesToUint16(data))
}

func TestByteCommon_Uint32ToBytes(t *testing.T) {
	t.Log(Byte().Uint32ToBytes(100))
}

func TestByteCommon_BytesToUint32(t *testing.T) {
	data := Byte().Uint32ToBytes(100)
	t.Log(Byte().BytesToUint32(data))
}

func TestByteCommon_Uint64ToBytes(t *testing.T) {
	t.Log(Byte().Uint64ToBytes(100))
}

func TestByteCommon_BytesToUint64(t *testing.T) {
	data := Byte().Uint64ToBytes(100)
	t.Log(Byte().BytesToUint64(data))
}
