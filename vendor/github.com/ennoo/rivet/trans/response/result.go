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
 *
 */

package response

import (
	"fmt"
	"github.com/ennoo/rivet/utils/log"
	"github.com/ennoo/rivet/utils/string"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

// Callback 返回结果对象介入降级操作方法
func (result *Result) Callback(callback func() *Result, err error) {
	if nil == callback || str.IsEmpty(callback().ResultCode) {
		log.Trans.Info("放弃降级或降级策略有误")
		result.ResultCode = Fail
		result.Msg = err.Error()
	} else {
		log.Trans.Info("降级回调")
		result.reSet(callback())
	}
}

func (result *Result) reSet(res *Result) {
	result.ResultCode = res.ResultCode
	result.Data = res.Data
	result.Msg = res.Msg
}

// FailErr 携带error信息,如果是respError，则
// 必然存在errorCode和msg，因此进行赋值。否则不赋值
func (result *Result) FailErr(err error) {
	switch vtype := err.(type) {
	case *RespError:
		result.Msg = vtype.ErrorMsg
		result.ResultCode = vtype.ErrorCode
	default:
		result.ResultCode = Fail
		log.Trans.Error(err.Error())
		//result.Msg = ServiceException.Msg
		result.Msg = err.Error()
	}

}

// Say response 返回自身
func (result *Result) Say(context *gin.Context) {
	context.JSON(http.StatusOK, &result)
}

// SaySuccess response 返回请求成功对象
func (result *Result) SaySuccess(context *gin.Context, obj interface{}) {
	result.Success(obj)
	context.JSON(http.StatusOK, &result)
}

// SayFail response 返回请求失败对象
func (result *Result) SayFail(context *gin.Context, msg string) {
	result.Fail(msg)
	context.JSON(http.StatusOK, &result)
}

// Write response 字节
func (result *Result) Write(context *gin.Context, byte []byte) (int, error) {
	return context.Writer.Write(byte)
}

// catchErr 捕获所有异常信息并放入json到context，便于controller直接调用
func catchErr(context *gin.Context, res *Result) {
	if r := recover(); r != nil {
		value, ok := r.(Exception)
		if ok {
			res.Fail(value.Msg)
			log.Trans.Error(fmt.Sprint(r))
			context.JSON(value.code, res)
		} else {
			res.Fail(fmt.Sprintf("An error occurred:%v \n", r))
			log.Trans.Error(fmt.Sprint(r))
			context.JSON(http.StatusInternalServerError, res)
		}
		return
	}
}
