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

package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// RestJSONHandler 处理 json 请求发送内容
type RestJSONHandler struct {
	RestHandler
	Param interface{}
}

// ObtainRemoteServer 获取本次 http 请求服务根路径 如：localhost:8080
func (handler *RestJSONHandler) ObtainRemoteServer() string {
	return handler.RemoteServer
}

// ObtainURI 获取本次 http 请求服务方法路径 如：/user/login
func (handler *RestJSONHandler) ObtainURI() string {
	return handler.URI
}

// ObtainBody 获取本次 http 请求 body io
func (handler *RestJSONHandler) ObtainBody() io.Reader {
	jsonByte, _ := json.Marshal(handler.Param)
	return bytes.NewReader(jsonByte)
}

// ObtainHeader 获取本次 http 请求 header
func (handler *RestJSONHandler) ObtainHeader() http.Header {
	handler.Header.Add("Content-Type", "application/json")
	return handler.Header
}

// ObtainCookies 获取本次 http 请求 cookies
func (handler *RestJSONHandler) ObtainCookies() []*http.Cookie {
	return handler.Cookies
}

// Post 发起 Post 请求，body 为请求后的返回内容，err 指出请求出错原因
func (handler *RestJSONHandler) Post(reqType int) (resp *http.Response, err error) {
	return request(http.MethodPost, handler, reqType)
}

// Put 发起 Put 请求，body 为请求后的返回内容，err 指出请求出错原因
func (handler *RestJSONHandler) Put(reqType int) (resp *http.Response, err error) {
	return request(http.MethodPut, handler, reqType)
}

// Delete 发起 Delete 请求，body 为请求后的返回内容，err 指出请求出错原因
func (handler *RestJSONHandler) Delete(reqType int) (resp *http.Response, err error) {
	return request(http.MethodDelete, handler, reqType)
}

// Patch 发起 Patch 请求，body 为请求后的返回内容，err 指出请求出错原因
func (handler *RestJSONHandler) Patch(reqType int) (resp *http.Response, err error) {
	return request(http.MethodPatch, handler, reqType)
}

// Options 发起 Options 请求，body 为请求后的返回内容，err 指出请求出错原因
func (handler *RestJSONHandler) Options(reqType int) (resp *http.Response, err error) {
	return request(http.MethodOptions, handler, reqType)
}

// Head 发起 Head 请求，body 为请求后的返回内容，err 指出请求出错原因
func (handler *RestJSONHandler) Head(reqType int) (resp *http.Response, err error) {
	return request(http.MethodHead, handler, reqType)
}

// Connect 发起 Connect 请求，body 为请求后的返回内容，err 指出请求出错原因
func (handler *RestJSONHandler) Connect(reqType int) (resp *http.Response, err error) {
	return request(http.MethodConnect, handler, reqType)
}

// Trace 发起 Trace 请求，body 为请求后的返回内容，err 指出请求出错原因
func (handler *RestJSONHandler) Trace(reqType int) (resp *http.Response, err error) {
	return request(http.MethodTrace, handler, reqType)
}

// Get 发起 Get 请求，body 为请求后的返回内容，err 指出请求出错原因
func (handler *RestJSONHandler) Get(reqType int) (resp *http.Response, err error) {
	return get(handler, reqType)
}
