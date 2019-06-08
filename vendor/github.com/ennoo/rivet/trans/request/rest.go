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
	"github.com/ennoo/rivet/utils/log"
	"github.com/ennoo/rivet/utils/slip"
	"github.com/ennoo/rivet/utils/string"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	// TransCallbackRequest 转发请求，需要考虑是否携带上一请求的 Header 和 Cookies
	TransCallbackRequest = iota
	// DirectJSONRequest 直接请求，自处理 Header 和 Cookies
	DirectJSONRequest
	// DirectTextRequest 直接请求，自处理 Header 和 Cookies
	DirectTextRequest
)

// RestHandler 处理请求发送内容
type RestHandler struct {
	RemoteServer string
	URI          string
	Body         io.ReadCloser
	Header       http.Header
	Cookies      []*http.Cookie
}

// Rest http 请求方法接口
type Rest interface {
	// Post 发起 Post 请求，body 为请求后的返回内容，err 指出请求出错原因
	Post(reqType int) (resp *http.Response, err error)

	// Put 发起 Put 请求，body 为请求后的返回内容，err 指出请求出错原因
	Put(reqType int) (resp *http.Response, err error)

	// Delete 发起 Delete 请求，body 为请求后的返回内容，err 指出请求出错原因
	Delete(reqType int) (resp *http.Response, err error)

	// Patch 发起 Patch 请求，body 为请求后的返回内容，err 指出请求出错原因
	Patch(reqType int) (resp *http.Response, err error)

	// Options 发起 Options 请求，body 为请求后的返回内容，err 指出请求出错原因
	Options(reqType int) (resp *http.Response, err error)

	// Head 发起 Head 请求，body 为请求后的返回内容，err 指出请求出错原因
	Head(reqType int) (resp *http.Response, err error)

	// Connect 发起 Connect 请求，body 为请求后的返回内容，err 指出请求出错原因
	Connect(reqType int) (resp *http.Response, err error)

	// Trace 发起 Trace 请求，body 为请求后的返回内容，err 指出请求出错原因
	Trace(reqType int) (resp *http.Response, err error)

	// Get 发起 Get 请求，body 为请求后的返回内容，err 指出请求出错原因
	Get(reqType int) (resp *http.Response, err error)
}

// Handler http 处理请求发送内容接口
type Handler interface {
	// ObtainRemoteServer 获取本次 http 请求服务根路径 如：localhost:8080
	ObtainRemoteServer() string

	// ObtainURI 获取本次 http 请求服务方法路径 如：/user/login
	ObtainURI() string

	// ObtainBody 获取本次 http 请求 body io
	ObtainBody() io.Reader

	// ObtainHeader 获取本次 http 请求 header
	ObtainHeader() http.Header

	// ObtainCookies 获取本次 http 请求 cookies
	ObtainCookies() []*http.Cookie
}

func addCookies(request *http.Request, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
}

func request(method string, handler Handler, reqType int) (*http.Response, error) {
	req, err := http.NewRequest(method, getFullURI(handler), handler.ObtainBody())
	if nil != err {
		slips := slip.NewSlip(slip.RestRequestError, err.Error(), nil)
		return nil, slips
	}
	switch reqType {
	case TransCallbackRequest:
		addCookies(req, handler.ObtainCookies())
		req.Header = handler.ObtainHeader()
	case DirectJSONRequest:
		req.Header.Add("Content-Type", "application/json")
	case DirectTextRequest:
		req.Header.Add("Content-Type", "text/html")
	}
	return exec(req)
}

// Get 发送get请求
func get(handler Handler, reqType int) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, getFullURI(handler), nil)
	if nil != err {
		slips := slip.NewSlip(slip.RestRequestError, err.Error(), nil)
		return nil, slips
	}
	switch reqType {
	case TransCallbackRequest:
		addCookies(req, handler.ObtainCookies())
		req.Header = handler.ObtainHeader()
	case DirectJSONRequest:
		req.Header.Add("Content-Type", "application/json")
	case DirectTextRequest:
		req.Header.Add("Content-Type", "text/html")
	}
	return exec(req)
}

func exec(req *http.Request) (resp *http.Response, err error) {
	if nil != GetTPInstance().Transport {
		client := http.Client{Transport: GetTPInstance().Transport}
		resp, err = client.Do(req)
	} else {
		resp, err = http.DefaultClient.Do(req)
	}

	if err != nil {
		slips := slip.NewSlip(slip.RestResponseError, err.Error(), resp)
		return nil, slips
	}

	return
}

func getFullURI(handler Handler) string {
	return filepath.ToSlash(strings.Join([]string{handler.ObtainRemoteServer(), filepath.Join("/", handler.ObtainURI())}, ""))
}

// GetAccessTokenFromReq 从Header或者Cookie中获取到用户的access_token
func GetAccessTokenFromReq(c *gin.Context) (token string) {
	var err error

	token = c.GetHeader("Access_token")
	if str.IsEmpty(token) {
		token = c.GetHeader("access_token")
		// 如果依然为空，则从cookie中尝试获取
		if str.IsEmpty(token) {
			token, err = c.Cookie("access_token")
			if err != nil {
				log.Trans.Error(err.Error())
				return ""
			}
		}
	}
	return token
}

// GetCookieByName 忽略大小写，找到指定的cookie
func GetCookieByName(cookies []*http.Cookie, cookieName string) *http.Cookie {
	for _, cookie := range cookies {
		if strings.EqualFold(cookie.Name, cookieName) {
			return cookie
		}
	}
	return nil
}
