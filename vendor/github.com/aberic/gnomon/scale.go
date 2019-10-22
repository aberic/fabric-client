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
	"math"
	"strconv"
	"strings"
)

// ScaleCommon 运算/转换工具
type ScaleCommon struct{}

var (
	hexIntMap = map[int]string{ // hexIntMap 十六进制对应十进制映射
		0: "0", 1: "1", 2: "2", 3: "3",
		4: "4", 5: "5", 6: "6", 7: "7",
		8: "8", 9: "9", 10: "a", 11: "b",
		12: "c", 13: "d", 14: "e", 15: "f",
	}
	intHexMap = map[string]int{ // intHexMap 十进制对应十六进制映射
		"0": 0, "1": 1, "2": 2, "3": 3,
		"4": 4, "5": 5, "6": 6, "7": 7,
		"8": 8, "9": 9, "a": 10, "b": 11,
		"c": 12, "d": 13, "e": 14, "f": 15,
	}
	duoIntMap = map[int]string{ // duoIntMap 三十二进制对应十进制映射
		0: "0", 1: "1", 2: "2", 3: "3",
		4: "4", 5: "5", 6: "6", 7: "7",
		8: "8", 9: "9", 10: "a", 11: "b",
		12: "c", 13: "d", 14: "e", 15: "f",
		16: "g", 17: "h", 18: "i", 19: "j",
		20: "k", 21: "l", 22: "m", 23: "n",
		24: "o", 25: "p", 26: "q", 27: "r",
		28: "s", 29: "t", 30: "u", 31: "v",
	}
	intDuoMap = map[string]int{ // intDuoMap 十进制对应三十二进制映射
		"0": 0, "1": 1, "2": 2, "3": 3,
		"4": 4, "5": 5, "6": 6, "7": 7,
		"8": 8, "9": 9, "a": 10, "b": 11,
		"c": 12, "d": 13, "e": 14, "f": 15,
		"g": 16, "h": 17, "i": 18, "j": 19,
		"k": 20, "l": 21, "m": 22, "n": 23,
		"o": 24, "p": 25, "q": 26, "r": 27,
		"s": 28, "t": 29, "u": 30, "v": 31,
	}
	dDuoIntMap = map[int]string{ // dDuoIntMap 64进制对应十进制映射
		0: "0", 1: "1", 2: "2", 3: "3",
		4: "4", 5: "5", 6: "6", 7: "7",
		8: "8", 9: "9", 10: "a", 11: "b",
		12: "c", 13: "d", 14: "e", 15: "f",
		16: "g", 17: "h", 18: "i", 19: "j",
		20: "k", 21: "l", 22: "m", 23: "n",
		24: "o", 25: "p", 26: "q", 27: "r",
		28: "s", 29: "t", 30: "u", 31: "v",
		32: "w", 33: "x", 34: "y", 35: "z",
		36: "A", 37: "B", 38: "C", 39: "D",
		40: "E", 41: "F", 42: "G", 43: "H",
		44: "I", 45: "J", 46: "K", 47: "L",
		48: "M", 49: "N", 50: "O", 51: "P",
		52: "Q", 53: "R", 54: "S", 55: "T",
		56: "U", 57: "V", 58: "W", 59: "X",
		60: "Y", 61: "Z", 62: "+", 63: "-",
	}
	intDDuoMap = map[string]int{ // intDDuoMap 十进制对应64进制映射
		"0": 0, "1": 1, "2": 2, "3": 3,
		"4": 4, "5": 5, "6": 6, "7": 7,
		"8": 8, "9": 9, "a": 10, "b": 11,
		"c": 12, "d": 13, "e": 14, "f": 15,
		"g": 16, "h": 17, "i": 18, "j": 19,
		"k": 20, "l": 21, "m": 22, "n": 23,
		"o": 24, "p": 25, "q": 26, "r": 27,
		"s": 28, "t": 29, "u": 30, "v": 31,
		"w": 32, "x": 33, "y": 34, "z": 35,
		"A": 36, "B": 37, "C": 38, "D": 39,
		"E": 40, "F": 41, "G": 42, "H": 43,
		"I": 44, "J": 45, "K": 46, "L": 47,
		"M": 48, "N": 49, "O": 50, "P": 51,
		"Q": 52, "R": 53, "S": 54, "T": 55,
		"U": 56, "V": 57, "W": 58, "X": 59,
		"Y": 60, "Z": 61, "+": 62, "-": 63,
	}
)

// Uint64ToHexString uint64转十六进制字符串
func (s *ScaleCommon) Uint64ToHexString(i uint64) string {
	iSt := ""
	for i > 0 {
		if i >= 16 {
			iSt = strings.Join([]string{hexIntMap[int(i%16)], iSt}, "")
			i /= 16
		} else if i > 0 && i < 16 {
			iSt = strings.Join([]string{hexIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// Int64ToHexString int64转十六进制字符串
func (s *ScaleCommon) Int64ToHexString(i int64) string {
	iSt := ""
	for i > 0 {
		if i >= 16 {
			iSt = strings.Join([]string{hexIntMap[int(i%16)], iSt}, "")
			i /= 16
		} else if i > 0 && i < 16 {
			iSt = strings.Join([]string{hexIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// Uint32ToHexString uint32转十六进制字符串
func (s *ScaleCommon) Uint32ToHexString(i uint32) string {
	iSt := ""
	for i > 0 {
		if i >= 16 {
			iSt = strings.Join([]string{hexIntMap[int(i%16)], iSt}, "")
			i /= 16
		} else if i > 0 && i < 16 {
			iSt = strings.Join([]string{hexIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// Int32ToHexString int32转十六进制字符串
func (s *ScaleCommon) Int32ToHexString(i int32) string {
	iSt := ""
	for i > 0 {
		if i >= 16 {
			iSt = strings.Join([]string{hexIntMap[int(i%16)], iSt}, "")
			i /= 16
		} else if i > 0 && i < 16 {
			iSt = strings.Join([]string{hexIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// UintToHexString uint转十六进制字符串
func (s *ScaleCommon) UintToHexString(i uint) string {
	iSt := ""
	for i > 0 {
		if i >= 16 {
			iSt = strings.Join([]string{hexIntMap[int(i%16)], iSt}, "")
			i /= 16
		} else if i > 0 && i < 16 {
			iSt = strings.Join([]string{hexIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// IntToHexString int转十六进制字符串
func (s *ScaleCommon) IntToHexString(i int) string {
	iSt := ""
	for i > 0 {
		if i >= 16 {
			iSt = strings.Join([]string{hexIntMap[int(i%16)], iSt}, "")
			i /= 16
		} else if i > 0 && i < 16 {
			iSt = strings.Join([]string{hexIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// HexStringToUint64 int字符串转int
func (s *ScaleCommon) HexStringToUint64(hex string) uint64 {
	hexLen := len(hex)
	var uint64Hex uint64
	for i := 0; i < hexLen; i++ {
		uint64Hex += uint64(intHexMap[hex[i:i+1]]) * uint64(math.Pow(16, float64(hexLen-i-1)))
	}
	return uint64Hex
}

// HexStringToInt64 int字符串转int
func (s *ScaleCommon) HexStringToInt64(hex string) int64 {
	hexLen := len(hex)
	var int64Hex int64
	for i := 0; i < hexLen; i++ {
		int64Hex += int64(intHexMap[hex[i:i+1]]) * int64(math.Pow(16, float64(hexLen-i-1)))
	}
	return int64Hex
}

// Uint64ToDuoString uint64转十六进制字符串
func (s *ScaleCommon) Uint64ToDuoString(i uint64) string {
	iSt := ""
	for i > 0 {
		if i >= 32 {
			iSt = strings.Join([]string{duoIntMap[int(i%32)], iSt}, "")
			i /= 32
		} else if i > 0 && i < 32 {
			iSt = strings.Join([]string{duoIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// Int64ToDuoString int64转十六进制字符串
func (s *ScaleCommon) Int64ToDuoString(i int64) string {
	iSt := ""
	for i > 0 {
		if i >= 32 {
			iSt = strings.Join([]string{duoIntMap[int(i%32)], iSt}, "")
			i /= 32
		} else if i > 0 && i < 32 {
			iSt = strings.Join([]string{duoIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// Uint32ToDuoString uint32转十六进制字符串
func (s *ScaleCommon) Uint32ToDuoString(i uint32) string {
	iSt := ""
	for i > 0 {
		if i >= 32 {
			iSt = strings.Join([]string{duoIntMap[int(i%32)], iSt}, "")
			i /= 32
		} else if i > 0 && i < 32 {
			iSt = strings.Join([]string{duoIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// Int32ToDuoString int32转十六进制字符串
func (s *ScaleCommon) Int32ToDuoString(i int32) string {
	iSt := ""
	for i > 0 {
		if i >= 32 {
			iSt = strings.Join([]string{duoIntMap[int(i%32)], iSt}, "")
			i /= 32
		} else if i > 0 && i < 32 {
			iSt = strings.Join([]string{duoIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// UintToDuoString uint转十六进制字符串
func (s *ScaleCommon) UintToDuoString(i uint) string {
	iSt := ""
	for i > 0 {
		if i >= 32 {
			iSt = strings.Join([]string{duoIntMap[int(i%32)], iSt}, "")
			i /= 32
		} else if i > 0 && i < 32 {
			iSt = strings.Join([]string{duoIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// IntToDuoString int转十六进制字符串
func (s *ScaleCommon) IntToDuoString(i int) string {
	iSt := ""
	for i > 0 {
		if i >= 32 {
			iSt = strings.Join([]string{duoIntMap[int(i%32)], iSt}, "")
			i /= 32
		} else if i > 0 && i < 32 {
			iSt = strings.Join([]string{duoIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// DuoStringToUint64 int字符串转int
func (s *ScaleCommon) DuoStringToUint64(duo string) uint64 {
	duoLen := len(duo)
	var uint64Duo uint64
	for i := 0; i < duoLen; i++ {
		uint64Duo += uint64(intDuoMap[duo[i:i+1]]) * uint64(math.Pow(32, float64(duoLen-i-1)))
	}
	return uint64Duo
}

// DuoStringToInt64 int字符串转int
func (s *ScaleCommon) DuoStringToInt64(duo string) int64 {
	duoLen := len(duo)
	var int64Duo int64
	for i := 0; i < duoLen; i++ {
		int64Duo += int64(intDuoMap[duo[i:i+1]]) * int64(math.Pow(32, float64(duoLen-i-1)))
	}
	return int64Duo
}

// Uint64ToDDuoString uint64转十六进制字符串
func (s *ScaleCommon) Uint64ToDDuoString(i uint64) string {
	iSt := ""
	for i > 0 {
		if i >= 64 {
			iSt = strings.Join([]string{dDuoIntMap[int(i%64)], iSt}, "")
			i /= 64
		} else if i > 0 && i < 64 {
			iSt = strings.Join([]string{dDuoIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// Int64ToDDuoString int64转十六进制字符串
func (s *ScaleCommon) Int64ToDDuoString(i int64) string {
	iSt := ""
	for i > 0 {
		if i >= 64 {
			iSt = strings.Join([]string{dDuoIntMap[int(i%64)], iSt}, "")
			i /= 64
		} else if i > 0 && i < 64 {
			iSt = strings.Join([]string{dDuoIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// Uint32ToDDuoString uint32转十六进制字符串
func (s *ScaleCommon) Uint32ToDDuoString(i uint32) string {
	iSt := ""
	for i > 0 {
		if i >= 64 {
			iSt = strings.Join([]string{dDuoIntMap[int(i%64)], iSt}, "")
			i /= 64
		} else if i > 0 && i < 64 {
			iSt = strings.Join([]string{dDuoIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// Int32ToDDuoString int32转十六进制字符串
func (s *ScaleCommon) Int32ToDDuoString(i int32) string {
	iSt := ""
	for i > 0 {
		if i >= 64 {
			iSt = strings.Join([]string{dDuoIntMap[int(i%64)], iSt}, "")
			i /= 64
		} else if i > 0 && i < 64 {
			iSt = strings.Join([]string{dDuoIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// UintToDDuoString int转64进制字符串
func (s *ScaleCommon) UintToDDuoString(i uint) string {
	iSt := ""
	for i > 0 {
		if i >= 64 {
			iSt = strings.Join([]string{dDuoIntMap[int(i%64)], iSt}, "")
			i /= 64
		} else if i > 0 && i < 64 {
			iSt = strings.Join([]string{dDuoIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// IntToDDuoString int转64进制字符串
func (s *ScaleCommon) IntToDDuoString(i int) string {
	iSt := ""
	for i > 0 {
		if i >= 64 {
			iSt = strings.Join([]string{dDuoIntMap[int(i%64)], iSt}, "")
			i /= 64
		} else if i > 0 && i < 64 {
			iSt = strings.Join([]string{dDuoIntMap[int(i)], iSt}, "")
			i = 0
		}
	}
	return iSt
}

// DDuoStringToUint64 int字符串转int
func (s *ScaleCommon) DDuoStringToUint64(dDuo string) uint64 {
	dDuoLen := len(dDuo)
	var uint64DDuo uint64
	for i := 0; i < dDuoLen; i++ {
		uint64DDuo += uint64(intDDuoMap[dDuo[i:i+1]]) * uint64(math.Pow(64, float64(dDuoLen-i-1)))
	}
	return uint64DDuo
}

// DDuoStringToInt64 int字符串转int
func (s *ScaleCommon) DDuoStringToInt64(dDuo string) int64 {
	dDuoLen := len(dDuo)
	var int64DDuo int64
	for i := 0; i < dDuoLen; i++ {
		int64DDuo += int64(intDDuoMap[dDuo[i:i+1]]) * int64(math.Pow(64, float64(dDuoLen-i-1)))
	}
	return int64DDuo
}

// Uint64Len 计算整型字符串长度
func (s *ScaleCommon) Uint64Len(i uint64) int {
	iLen := 1
	for i >= 10 {
		i /= 10
		iLen++
	}
	return iLen
}

// Int64Len 计算整型字符串长度
func (s *ScaleCommon) Int64Len(i int64) int {
	iLen := 1
	for i >= 10 {
		i /= 10
		iLen++
	}
	return iLen
}

// Uint32Len 计算整型字符串长度
func (s *ScaleCommon) Uint32Len(i uint32) int {
	iLen := 1
	for i >= 10 {
		i /= 10
		iLen++
	}
	return iLen
}

// Int32Len 计算整型字符串长度
func (s *ScaleCommon) Int32Len(i int32) int {
	iLen := 1
	for i >= 10 {
		i /= 10
		iLen++
	}
	return iLen
}

// UintLen 计算整型字符串长度
func (s *ScaleCommon) UintLen(i uint) int {
	iLen := 1
	for i >= 10 {
		i /= 10
		iLen++
	}
	return iLen
}

// IntLen 计算整型字符串长度
func (s *ScaleCommon) IntLen(i int) int {
	iLen := 1
	for i >= 10 {
		i /= 10
		iLen++
	}
	return iLen
}

// Uint8toFullState 补全不满三位数状态，如1->001、34->034、215->215
func (s *ScaleCommon) Uint8toFullState(index uint8) string {
	result := strconv.Itoa(int(index))
	if index < 10 {
		return strings.Join([]string{"00", result}, "")
	} else if index < 100 {
		return strings.Join([]string{"0", result}, "")
	}
	return result
}

// Uint32toFullState 补全不满十位数状态，如1->0000000001、34->0000000034、215->0000000215
func (s *ScaleCommon) Uint32toFullState(index uint32) string {
	indexStr := strconv.Itoa(int(index))
	if s.Uint32Len(index) >= 10 {
		return indexStr
	}
	pos := 0
	for index > 1 {
		index /= 10
		pos++
	}
	backZero := 10 - pos
	for i := 0; i < backZero; i++ {
		indexStr = strings.Join([]string{"0", indexStr}, "")
	}
	return indexStr
}

// Float64toInt64 将float64转成精确的int64
func (s *ScaleCommon) Float64toInt64(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}

// Int64toFloat64 将int64恢复成正常的float64
func (s *ScaleCommon) Int64toFloat64(num int64, retain int) float64 {
	return float64(num) / math.Pow10(retain)
}
