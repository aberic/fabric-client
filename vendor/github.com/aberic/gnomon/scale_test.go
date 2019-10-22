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

import "testing"

func TestScaleHex(t *testing.T) {
	var (
		i    int
		ui   uint
		i8   int8
		ui8  uint8
		i16  int16
		ui16 uint16
		i32  int32
		ui32 uint32
		i64  int64
		ui64 uint64
	)
	i = 2999999
	ui = 2888888
	i8 = 127
	ui8 = 255
	i16 = 32767
	ui16 = 65535
	i32 = 99922299
	ui32 = 88811188
	i64 = 827639847923879
	ui64 = 92873890910928019

	iStr := Scale().IntToHexString(i)
	Log().Debug("i", Log().Field("2999999", iStr))
	Log().Debug("i", Log().Field("2999999", Scale().HexStringToUint64(iStr)))
	Log().Debug("")

	iStr = Scale().UintToHexString(uint(i))
	Log().Debug("i", Log().Field("2999999", iStr))
	Log().Debug("i", Log().Field("2999999", Scale().HexStringToUint64(iStr)))
	Log().Debug("")

	uiStr := Scale().UintToHexString(uint(ui))
	Log().Debug("ui", Log().Field("2888888", uiStr))
	Log().Debug("i", Log().Field("2888888", Scale().HexStringToUint64(uiStr)))
	Log().Debug("")

	i8Str := Scale().UintToHexString(uint(i8))
	Log().Debug("i8", Log().Field("127", i8Str))
	Log().Debug("i8", Log().Field("127", Scale().HexStringToUint64(i8Str)))
	Log().Debug("")

	ui8Str := Scale().UintToHexString(uint(ui8))
	Log().Debug("ui8", Log().Field("255", ui8Str))
	Log().Debug("ui8", Log().Field("255", Scale().HexStringToUint64(ui8Str)))
	Log().Debug("")

	i16Str := Scale().UintToHexString(uint(i16))
	Log().Debug("i16", Log().Field("32767", i16Str))
	Log().Debug("i16", Log().Field("32767", Scale().HexStringToUint64(i16Str)))
	Log().Debug("")

	ui16Str := Scale().UintToHexString(uint(ui16))
	Log().Debug("ui16", Log().Field("65535", ui16Str))
	Log().Debug("ui16", Log().Field("65535", Scale().HexStringToUint64(ui16Str)))
	Log().Debug("")

	i32Str := Scale().Int32ToHexString(i32)
	Log().Debug("i32", Log().Field("99922299", i32Str))
	Log().Debug("i32", Log().Field("99922299", Scale().HexStringToUint64(i32Str)))
	Log().Debug("")

	ui32Str := Scale().Uint32ToHexString(ui32)
	Log().Debug("ui32", Log().Field("88811188", ui32Str))
	Log().Debug("ui32", Log().Field("88811188", Scale().HexStringToUint64(ui32Str)))
	Log().Debug("")

	i64Str := Scale().Int64ToHexString(i64)
	Log().Debug("i64", Log().Field("827639847923879", i64Str))
	Log().Debug("i64", Log().Field("827639847923879", Scale().HexStringToInt64(i64Str)))
	Log().Debug("")

	ui64Str := Scale().Uint64ToHexString(ui64)
	Log().Debug("ui64", Log().Field("92873890910928019", ui64Str))
	Log().Debug("ui64", Log().Field("92873890910928019", Scale().HexStringToUint64(ui64Str)))
}

func TestScaleDuo(t *testing.T) {
	var (
		i    int
		ui   uint
		i8   int8
		ui8  uint8
		i16  int16
		ui16 uint16
		i32  int32
		ui32 uint32
		i64  int64
		ui64 uint64
	)
	i = 2999999
	ui = 2888888
	i8 = 127
	ui8 = 255
	i16 = 32767
	ui16 = 65535
	i32 = 99922299
	ui32 = 88811188
	i64 = 827639847923879
	ui64 = 92873890910928019

	iStr := Scale().IntToDuoString(i)
	Log().Debug("i", Log().Field("2999999", iStr))
	Log().Debug("i", Log().Field("2999999", Scale().DuoStringToUint64(iStr)))
	Log().Debug("")

	iStr = Scale().UintToDuoString(uint(i))
	Log().Debug("i", Log().Field("2999999", iStr))
	Log().Debug("i", Log().Field("2999999", Scale().DuoStringToUint64(iStr)))
	Log().Debug("")

	uiStr := Scale().UintToDuoString(uint(ui))
	Log().Debug("ui", Log().Field("2888888", uiStr))
	Log().Debug("i", Log().Field("2888888", Scale().DuoStringToUint64(uiStr)))
	Log().Debug("")

	i8Str := Scale().UintToDuoString(uint(i8))
	Log().Debug("i8", Log().Field("127", i8Str))
	Log().Debug("i8", Log().Field("127", Scale().DuoStringToUint64(i8Str)))
	Log().Debug("")

	ui8Str := Scale().UintToDuoString(uint(ui8))
	Log().Debug("ui8", Log().Field("255", ui8Str))
	Log().Debug("ui8", Log().Field("255", Scale().DuoStringToUint64(ui8Str)))
	Log().Debug("")

	i16Str := Scale().UintToDuoString(uint(i16))
	Log().Debug("i16", Log().Field("32767", i16Str))
	Log().Debug("i16", Log().Field("32767", Scale().DuoStringToUint64(i16Str)))
	Log().Debug("")

	ui16Str := Scale().UintToDuoString(uint(ui16))
	Log().Debug("ui16", Log().Field("65535", ui16Str))
	Log().Debug("ui16", Log().Field("65535", Scale().DuoStringToUint64(ui16Str)))
	Log().Debug("")

	i32Str := Scale().Int32ToDuoString(i32)
	Log().Debug("i32", Log().Field("99922299", i32Str))
	Log().Debug("i32", Log().Field("99922299", Scale().DuoStringToUint64(i32Str)))
	Log().Debug("")

	ui32Str := Scale().Uint32ToDuoString(ui32)
	Log().Debug("ui32", Log().Field("88811188", ui32Str))
	Log().Debug("ui32", Log().Field("88811188", Scale().DuoStringToUint64(ui32Str)))
	Log().Debug("")

	i64Str := Scale().Int64ToDuoString(i64)
	Log().Debug("i64", Log().Field("827639847923879", i64Str))
	Log().Debug("i64", Log().Field("827639847923879", Scale().DuoStringToInt64(i64Str)))
	Log().Debug("")

	ui64Str := Scale().Uint64ToDuoString(ui64)
	Log().Debug("ui64", Log().Field("92873890910928019", ui64Str))
	Log().Debug("ui64", Log().Field("92873890910928019", Scale().DuoStringToUint64(ui64Str)))
}

func TestScaleDDuo(t *testing.T) {
	var (
		i    int
		ui   uint
		i8   int8
		ui8  uint8
		i16  int16
		ui16 uint16
		i32  int32
		ui32 uint32
		i64  int64
		ui64 uint64
	)
	i = 2999999
	ui = 2888888
	i8 = 127
	ui8 = 255
	i16 = 32767
	ui16 = 65535
	i32 = 99922299
	ui32 = 88811188
	i64 = 9223372036854770018
	ui64 = 92873890910928019

	iStr := Scale().IntToDDuoString(i)
	Log().Debug("i", Log().Field("2999999", iStr))
	Log().Debug("i", Log().Field("2999999", Scale().DDuoStringToUint64(iStr)))
	Log().Debug("")

	iStr = Scale().UintToDDuoString(uint(i))
	Log().Debug("i", Log().Field("2999999", iStr))
	Log().Debug("i", Log().Field("2999999", Scale().DDuoStringToUint64(iStr)))
	Log().Debug("")

	uiStr := Scale().UintToDDuoString(uint(ui))
	Log().Debug("ui", Log().Field("2888888", uiStr))
	Log().Debug("i", Log().Field("2888888", Scale().DDuoStringToUint64(uiStr)))
	Log().Debug("")

	i8Str := Scale().UintToDDuoString(uint(i8))
	Log().Debug("i8", Log().Field("127", i8Str))
	Log().Debug("i8", Log().Field("127", Scale().DDuoStringToUint64(i8Str)))
	Log().Debug("")

	ui8Str := Scale().UintToDDuoString(uint(ui8))
	Log().Debug("ui8", Log().Field("255", ui8Str))
	Log().Debug("ui8", Log().Field("255", Scale().DDuoStringToUint64(ui8Str)))
	Log().Debug("")

	i16Str := Scale().UintToDDuoString(uint(i16))
	Log().Debug("i16", Log().Field("32767", i16Str))
	Log().Debug("i16", Log().Field("32767", Scale().DDuoStringToUint64(i16Str)))
	Log().Debug("")

	ui16Str := Scale().UintToDDuoString(uint(ui16))
	Log().Debug("ui16", Log().Field("65535", ui16Str))
	Log().Debug("ui16", Log().Field("65535", Scale().DDuoStringToUint64(ui16Str)))
	Log().Debug("")

	i32Str := Scale().Int32ToDDuoString(i32)
	Log().Debug("i32", Log().Field("99922299", i32Str))
	Log().Debug("i32", Log().Field("99922299", Scale().DDuoStringToUint64(i32Str)))
	Log().Debug("")

	ui32Str := Scale().Uint32ToDDuoString(ui32)
	Log().Debug("ui32", Log().Field("88811188", ui32Str))
	Log().Debug("ui32", Log().Field("88811188", Scale().DDuoStringToUint64(ui32Str)))
	Log().Debug("")

	i64Str := Scale().Int64ToDDuoString(i64)
	Log().Debug("i64", Log().Field("827639847923879", i64Str))
	Log().Debug("i64", Log().Field("827639847923879", Scale().DDuoStringToInt64(i64Str)))
	Log().Debug("")

	ui64Str := Scale().Uint64ToDDuoString(ui64)
	Log().Debug("ui64", Log().Field("92873890910928019", ui64Str))
	Log().Debug("ui64", Log().Field("92873890910928019", Scale().DDuoStringToUint64(ui64Str)))

	ui64 = 18446744073709551615
	ui64Str = Scale().Uint64ToDDuoString(ui64)
	Log().Debug("ui64", Log().Field("18446744073709551615", ui64Str))
	Log().Debug("ui64", Log().Field("18446744073709551615", Scale().DDuoStringToUint64(ui64Str)))

	for i := 0; i < 100000; i++ {
		ui64Str = Scale().Uint64ToDDuoString(ui64)
		Scale().DDuoStringToUint64(ui64Str)
		//Log().Debug("ui64", Log().Field("18446744073709551615", ui64Str))
		//Log().Debug("ui64", Log().Field("18446744073709551615", Scale().DDuoStringToUint64(ui64Str)))
	}
}

func TestScaleLen(t *testing.T) {
	var (
		i    int
		ui   uint
		i32  int32
		ui32 uint32
		i64  int64
		ui64 uint64
	)
	i = 2999999
	ui = 2888888
	i32 = 99922299
	ui32 = 88811188
	i64 = 827639847923879
	ui64 = 92873890910928019
	t.Log(Scale().Uint64Len(ui64))
	t.Log(Scale().Int64Len(i64))
	t.Log(Scale().Uint32Len(ui32))
	t.Log(Scale().Int32Len(i32))
	t.Log(Scale().UintLen(ui))
	t.Log(Scale().IntLen(i))
}

func TestScaleFullState(t *testing.T) {
	var (
		ui8  uint8
		ui32 uint32
	)
	ui8 = 25
	ui32 = 88811188
	t.Log(Scale().Uint8toFullState(ui8))
	t.Log(Scale().Uint32toFullState(ui32))
}

func TestScaleFloat64(t *testing.T) {
	var (
		i64 int64
		f64 float64
	)
	i64 = 87372
	f64 = 92837.87263876498
	t.Log(Scale().Int64toFloat64(i64, 1))
	t.Log(Scale().Int64toFloat64(i64, 2))
	t.Log(Scale().Int64toFloat64(i64, 3))
	t.Log(Scale().Int64toFloat64(i64, 4))
	t.Log(Scale().Int64toFloat64(i64, 5))
	t.Log()
	t.Log(Scale().Float64toInt64(f64, 1))
	t.Log(Scale().Float64toInt64(f64, 2))
	t.Log(Scale().Float64toInt64(f64, 3))
	t.Log(Scale().Float64toInt64(f64, 4))
	t.Log(Scale().Float64toInt64(f64, 5))
}
