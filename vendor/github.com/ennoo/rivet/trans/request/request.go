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

// Package request 请求接受处理包
package request

import (
	"encoding/json"
	"errors"
	"github.com/ennoo/rivet/server"
	"github.com/ennoo/rivet/trans/response"
	"github.com/ennoo/rivet/utils/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var (
	// LB 是否开启负载均衡
	LB = false
	// Route 是否开启网关路由
	Route = false
	req   = sync.Pool{
		New: func() interface{} {
			return &Request{}
		},
	}
)

// Request 提供实例化调用请求方法，并内置返回策略
type Request struct {
	result response.Result
}

// SyncPoolGetRequest 提供实例化调用请求方法，并内置返回策略
func SyncPoolGetRequest() *Request {
	return req.Get().(*Request)
}

// RestJSONByURL JSON 请求
//
// method：请求方法
//
// url：完整请求路径
//
// param 请求对象
func (request *Request) RestJSONByURL(method string, url string, param interface{}) ([]byte, error) {
	remote, uri := remoteURI(url)
	return request.RestJSON(method, remote, uri, param)
}

// RestJSON JSON 请求
//
// method：请求方法
//
// remote：请求主体域名
//
// uri：请求主体方法路径
//
// param 请求对象
func (request *Request) RestJSON(method string, remote string, uri string, param interface{}) ([]byte, error) {
	restJSONHandler := RestJSONHandler{
		RestHandler: RestHandler{
			RemoteServer: remote,
			URI:          uri,
		},
		Param: param,
	}
	var body []byte
	var resp *http.Response
	var err error

	switch method {
	case http.MethodGet:
		resp, err = restJSONHandler.Get(DirectJSONRequest)
	case http.MethodHead:
		resp, err = restJSONHandler.Head(DirectJSONRequest)
	case http.MethodPost:
		resp, err = restJSONHandler.Post(DirectJSONRequest)
	case http.MethodPut:
		resp, err = restJSONHandler.Put(DirectJSONRequest)
	case http.MethodPatch:
		resp, err = restJSONHandler.Patch(DirectJSONRequest)
	case http.MethodDelete:
		resp, err = restJSONHandler.Delete(DirectJSONRequest)
	case http.MethodConnect:
		resp, err = restJSONHandler.Connect(DirectJSONRequest)
	case http.MethodOptions:
		resp, err = restJSONHandler.Options(DirectJSONRequest)
	case http.MethodTrace:
		resp, err = restJSONHandler.Trace(DirectJSONRequest)
	}
	return restDone(body, resp, err)
}

// RestTextByURL TEXT 请求
//
// method：请求方法
//
// url：完整请求路径
//
// param 请求对象
func (request *Request) RestTextByURL(method string, url string, values url.Values) ([]byte, error) {
	remote, uri := remoteURI(url)
	return request.RestText(method, remote, uri, values)
}

// RestText TEXT 请求
//
// method：请求方法
//
// remote：请求主体域名
//
// uri：请求主体方法路径
//
// values 请求参数
func (request *Request) RestText(method string, remote string, uri string, values url.Values) ([]byte, error) {
	restTextHandler := RestTextHandler{
		Values: values,
	}
	var body []byte
	var resp *http.Response
	var err error

	switch method {
	case http.MethodGet:
		resp, err = restTextHandler.Get(DirectTextRequest)
	case http.MethodHead:
		resp, err = restTextHandler.Head(DirectTextRequest)
	case http.MethodPost:
		resp, err = restTextHandler.Post(DirectTextRequest)
	case http.MethodPut:
		resp, err = restTextHandler.Put(DirectTextRequest)
	case http.MethodPatch:
		resp, err = restTextHandler.Patch(DirectTextRequest)
	case http.MethodDelete:
		resp, err = restTextHandler.Delete(DirectTextRequest)
	case http.MethodConnect:
		resp, err = restTextHandler.Connect(DirectTextRequest)
	case http.MethodOptions:
		resp, err = restTextHandler.Options(DirectTextRequest)
	case http.MethodTrace:
		resp, err = restTextHandler.Trace(DirectTextRequest)
	}
	return restDone(body, resp, err)
}

func restDone(body []byte, resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	if nil != resp {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
	} else {
		err = errors.New("response is nil")
	}
	return body, err
}

// CallByURL 请求转发处理方案
//
// context：原请求上下文
//
// method：即将转发的请求方法
//
// url：完整转发路径
func (request *Request) CallByURL(context *gin.Context, method string, url string) {
	remote, uri := remoteURI(url)
	request.Call(context, method, remote, uri)
}

// Call 请求转发处理方案
//
// context：原请求上下文
//
// method：即将转发的请求方法
//
// remote：请求转发主体域名
//
// uri：请求转发主体方法路径
func (request *Request) Call(context *gin.Context, method string, remote string, uri string) {
	request.call(context, method, remote, uri, nil)
}

// CallbackByURL 请求转发处理方案
//
// context：原请求上下文
//
// method：即将转发的请求方法
//
// url：完整转发路径
//
// callback：请求转发失败后回调降级策略
//
// callback *response.Result 请求转发降级后返回请求方结果对象
func (request *Request) CallbackByURL(context *gin.Context, method string, url string, callback func() *response.Result) {
	remote, uri := remoteURI(url)
	request.Callback(context, method, remote, uri, callback)
}

// Callback 请求转发处理方案
//
// context：原请求上下文
//
// method：即将转发的请求方法
//
// remote：请求转发主体域名
//
// uri：请求转发主体方法路径
//
// callback：请求转发失败后回调降级策略
//
// callback *response.Result 请求转发降级后返回请求方结果对象
func (request *Request) Callback(context *gin.Context, method string, remote string, uri string, callback func() *response.Result) {
	request.call(context, method, remote, uri, callback)
}

// call 请求转发处理方案
//
// context：原请求上下文
//
// method：即将转发的请求方法
//
// remote：请求转发主体域名
//
// uri：请求转发主体方法路径
//
// callback：请求转发失败后回调降级策略
//
// callback *response.Result 请求转发降级后返回请求方结果对象
func (request *Request) call(context *gin.Context, method string, remote string, uri string, callback func() *response.Result) {
	request.callReal(context, method, remote, uri, callback)
}

// callReal 请求转发处理方案
//
// context：原请求上下文
//
// method：即将转发的请求方法
//
// remote：请求转发主体域名
//
// uri：请求转发主体方法路径
//
// callback：请求转发失败后回调降级策略
//
// callback *response.Result 请求转发降级后返回请求方结果对象
func (request *Request) callReal(context *gin.Context, method string, remote string, uri string, callback func() *response.Result) {
	req := context.Request
	cookies := req.Cookies()
	restTransHandler := RestTransHandler{
		RestHandler: RestHandler{
			RemoteServer: remote,
			URI:          uri,
			Body:         req.Body,
			Header:       req.Header,
			Cookies:      cookies}}
	var body []byte
	var resp *http.Response
	var err error

	switch method {
	case http.MethodGet:
		resp, err = restTransHandler.Get(TransCallbackRequest)
	case http.MethodHead:
		resp, err = restTransHandler.Head(TransCallbackRequest)
	case http.MethodPost:
		resp, err = restTransHandler.Post(TransCallbackRequest)
	case http.MethodPut:
		resp, err = restTransHandler.Put(TransCallbackRequest)
	case http.MethodPatch:
		resp, err = restTransHandler.Patch(TransCallbackRequest)
	case http.MethodDelete:
		resp, err = restTransHandler.Delete(TransCallbackRequest)
	case http.MethodConnect:
		resp, err = restTransHandler.Connect(TransCallbackRequest)
	case http.MethodOptions:
		resp, err = restTransHandler.Options(TransCallbackRequest)
	case http.MethodTrace:
		resp, err = restTransHandler.Trace(TransCallbackRequest)
	}
	request.callDone(context, body, resp, err, callback)
}

func (request *Request) callDone(context *gin.Context, body []byte, resp *http.Response, err error, callback func() *response.Result) {
	if err != nil {
		request.result.Fail(err.Error())
		context.JSON(http.StatusOK, request.result)
		return
	}

	if nil != resp {
		defer resp.Body.Close()
		bodyRead, err := ioutil.ReadAll(resp.Body)
		done(context, resp, request, bodyRead, err, callback)
	} else {
		request.result.Fail("Response is nil")
		context.JSON(http.StatusOK, request.result)
	}
}

// done 请求转发处理结果
//
// 转发请求或降级回调
func done(context *gin.Context, resp *http.Response, request *Request, body []byte, err error, callback func() *response.Result) {
	if err != nil {
		request.result.Callback(callback, err)
	} else {
		log.Trans.Debug("body = " + string(body))

		if err := json.Unmarshal(body, &request.result); nil != err {
			request.result.Fail(err.Error())
		}
	}

	//for k := range resp.Header {
	//	context.Writer.Header().Add(k, resp.Header.Get(k))
	//}
	for index := range resp.Cookies() {
		cookie := resp.Cookies()[index]
		context.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	}
	context.JSON(http.StatusOK, request.result)
}

// FormatURL 将注册到服务的 service host 和 port 组装成完成的 URL
func FormatURL(context *gin.Context, service *server.Service) string {
	switch context.Request.Proto {
	case "HTTP/1.1":
		return strings.Join([]string{"http://", service.Host, ":", strconv.Itoa(service.Port)}, "")
	case "HTTP/2.0":
		return strings.Join([]string{"https://", service.Host, ":", strconv.Itoa(service.Port)}, "")
	default:
		return ""
	}
}

func remoteURI(url string) (remote, uri string) {
	urlTmp := url
	if strings.Contains(urlTmp, "//") {
		urlTmp = strings.Split(urlTmp, "//")[1]
	}
	size := len(strings.Split(urlTmp, "/")[0]) + 1
	urlTmp = urlTmp[size:]
	remote = url[0:(len(url) - len(urlTmp) - 1)]
	uri = urlTmp
	log.Trans.Debug("remoteURI", zap.String("url", url), zap.String("remote", remote), zap.String("uri", uri))
	return
}
