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

// Package response 请求返回处理包
package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var resp = sync.Pool{
	New: func() interface{} {
		return &Response{}
	}}

// Response 提供实例化调用 Do 方法，并内置返回策略
type Response struct {
	result Result
}

// SyncPoolGetResponse 提供实例化调用 Do 方法，并内置返回策略
func SyncPoolGetResponse() *Response {
	return resp.Get().(*Response)
}

// Do 处理 request 请求
//
// context：请求上下文
//
// obj：请求中 body 中内容期望转换的对象并做空实例化，如 new(Type)
//
// objBlock：obj 对象的回调方法，最终调用 Do 函数的方法会接收到返回值
//
// objBlock interface{}：obj 对象的回调方法所返回的最终交由 response 输出的对象
//
// objBlock error：obj 对象的回调方法所返回的错误对象
//
// 如未出现err，且无可描述返回内容，则返回值可为 (nil, nil)
func (response *Response) Do(context *gin.Context, objBlock func(result *Result)) {
	defer catchErr(context, &response.result)
	objBlock(&response.result)
}

// DoSelf 处理 request 请求
//
// context：请求上下文
//
// obj：请求中 body 中内容期望转换的对象并做空实例化，如 new(Type)
//
// objBlock：obj 对象的回调方法，最终调用 Do 函数的方法会接收到返回值
//
// objBlock interface{}：obj 对象的回调方法所返回的最终交由 response 输出的对象
//
// objBlock error：obj 对象的回调方法所返回的错误对象
//
// 如未出现err，且无可描述返回内容，则返回值可为 (nil, nil)
func (response *Response) DoSelf(context *gin.Context, self func(writer http.ResponseWriter, request *http.Request)) {
	defer catchErr(context, &response.result)
	self(context.Writer, context.Request)
}
