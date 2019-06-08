/*
 * Copyright (c) 2019. ENNOO - All Rights Reserved.
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

// Package slip 自定义错误信息
package slip

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	// StartError 启动有误
	StartError = iota
	// StartSuccess 启动成功
	StartSuccess
	// DiscoveryStart 启动发现服务
	DiscoveryStart
	// RestRequestError Rest 请求调用有误
	RestRequestError
	// RestResponseError Rest 请求返回有误
	RestResponseError
	// JSONUnmarshalError JSON 逆向解析有误
	JSONUnmarshalError
)

// Slip 自定义 error 对象
type Slip struct {
	Code     int
	Msg      string
	Response *http.Response
}

// Error 实现 error 接口
//
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
func (slip *Slip) Error() string {
	return strings.Join([]string{slip.Msg, ",error code is:", strconv.Itoa(slip.Code)}, "")
}

// FormatError 将 error 对象格式化成 slip 对象
func (slip *Slip) FormatError(errCode int, err error) {
	slip.Code = errCode
	slip.Msg = err.Error()
}

// NewSlip 新建 slip，slip 实现了 error 接口
func NewSlip(errCode int, errMsg string, response *http.Response) *Slip {
	return &Slip{
		Code: errCode,
		Msg:  errMsg,
		// 默认正常
		Response: response,
	}
}
