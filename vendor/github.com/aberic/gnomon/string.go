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
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// StringCommon 字符串工具
type StringCommon struct{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// IsEmpty 判断字符串是否为空，是则返回true，否则返回false
func (s *StringCommon) IsEmpty(str string) bool {
	return len(str) == 0
}

// IsNotEmpty 和 IsEmpty 的语义相反
func (s *StringCommon) IsNotEmpty(str string) bool {
	return !s.IsEmpty(str)
}

// Convert 下划线转换，首字母小写变大写，
// 下划线去掉并将下划线后的首字母大写
func (s *StringCommon) Convert(oriString string) string {
	cb := []byte(oriString)
	em := make([]byte, 0, 10)
	b := false
	for i, by := range cb {
		// 首字母如果是小写，则转换成大写
		if i == 0 && (97 <= by && by <= 122) {
			by = by - 32
		} else if by == 95 {
			// 下一个单词要变成大写
			b = true
			continue
		}
		if b {
			if 97 <= by && by <= 122 {
				by = by - 32
			}
			b = false
		}
		em = append(em, by)
	}
	return string(em)
}

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandSeq 创建指定长度的随机字符串
func (s *StringCommon) RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// RandSeq16 创建长度为16的随机字符串
func (s *StringCommon) RandSeq16() string {
	return s.RandSeq(16)
}

// Trim 去除字符串中的空格和换行符
func (s *StringCommon) Trim(str string) string {
	// 去除空格
	str = strings.Replace(str, " ", "", -1)
	return s.TrimN(str)
}

// TrimN 去除字符串中的换行符
func (s *StringCommon) TrimN(str string) string {
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)
	return str
}

// ToString 将对象格式化成字符串
func (s *StringCommon) ToString(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("%+v", i)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", i)
	}
	return out.String()
}

// SingleSpace 将字符串内所有连续空格替换为单个空格
func (s *StringCommon) SingleSpace(res string) string {
	for skip := false; !skip; {
		resNew := strings.Replace(res, "  ", " ", -1)
		if res == resNew {
			skip = true
		}
		res = resNew
	}
	return res
}

// PrefixSupplementZero 当字符串长度不满足时，将字符串前几位补充0
//
// str 字符串内容
//
// offset 字符串期望的长度
func (s *StringCommon) PrefixSupplementZero(str string, offset int) string {
	backZero := offset - len(str)
	if backZero <= 0 {
		return str
	}
	for i := 0; i < backZero; i++ {
		str = strings.Join([]string{"0", str}, "")
	}
	return str
}
