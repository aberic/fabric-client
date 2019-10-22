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

package sdk

import "github.com/ennoo/rivet/utils/log"

const (
	// Success 请求返回成功码
	Success = "200"
	// Fail 请求返回失败码
	Fail = "9999"
)

// Result 请求返回对象实体
type Result struct {
	ResultCode string `json:"code"`
	Msg        string `json:"msg"`
	// 数据接口
	Data interface{} `json:"data"`
}

// Success 默认成功返回
func (result *Result) Success(obj interface{}) {
	result.Msg = "Success!"
	result.ResultCode = Success
	result.Data = obj
}

// Fail 方法主要提供返回错误的json数据
func (result *Result) Fail(msg string) {
	result.Msg = msg
	result.ResultCode = Fail
}

// FailErr 携带error信息
func (result *Result) FailErr(err error) {
	result.ResultCode = Fail
	log.Trans.Error(err.Error())
	//result.Msg = ServiceException.Msg
	result.Msg = err.Error()
}
