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

package request

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/ennoo/rivet/utils/log"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

var instance *TPort
var once sync.Once

// TPort HTTPS请求管理对象
type TPort struct {
	Pool           *x509.CertPool
	ConnectTimeout time.Duration
	KeepAlive      time.Duration
	Transport      *http.Transport
}

// GetTPInstance 获取HTTPS请求管理对象 TPort 单例
func GetTPInstance() *TPort {
	once.Do(func() {
		instance = &TPort{
			Pool:           x509.NewCertPool(),
			ConnectTimeout: 3 * time.Second,
			KeepAlive:      30 * time.Second,
		}
	})
	return instance
}

// Timeout 设置 Transport 作为客户端请求超时设定
//
// connectTimeout 拨号等待连接完成的最长时间，TCP超时的时间一般在3s，默认3s
//
// keepAlive 指定保持活动网络连接的时间，如果为0，则不启用keep-alive，默认30s
func (t *TPort) Timeout(connectTimeout, keepAlive time.Duration) *TPort {
	t.ConnectTimeout = connectTimeout
	t.KeepAlive = keepAlive
	return t
}

// RootCACerts 初始化TLS作为客户端可用HTTPS证书池，用于对服务端验证
//
// caCertPaths 作为客户端发起 HTTPS 请求时所需客户端证书路径数组
func (t *TPort) RootCACerts(caCertPaths []string) *TPort {
	if len(caCertPaths) > 0 {
		for index := range caCertPaths {
			caCertPath := caCertPaths[index]
			if caCrt, err := ioutil.ReadFile(caCertPath); err == nil {
				//将生成的数字证书添加到数字证书集合中
				t.Pool.AppendCertsFromPEM(caCrt)
			} else {
				log.Trans.Fatal("can't read file with path",
					zap.String("path", caCertPath),
					zap.Error(err),
					zap.Int("code", log.ReadFileForCACertFail))
			}
		}
	}
	return t
}

// Instantiate 实例化 Transport 作为客户端可用HTTP/HTTPS对象
func (t *TPort) Instantiate() {
	t.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: t.Pool,
		},
		//TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, # 略过对服务端TLS的校验
		DialContext: (&net.Dialer{
			// Timeout是拨号等待连接完成的最长时间，如果还设置了Deadline，那么可能会提前失败，默认是没有超时。
			// 使用TCP拨号到多个IP地址的主机名时，可以在它们之间划分超时。
			// 即使有或没有Timeout，操作系统也会强加自己的超时时间，例如，TCP超时的时间一般在3s
			Timeout: t.ConnectTimeout,
			// KeepAlive指定保持活动网络连接的时间，如果为0，则不启用keep-alive，不支持keep-alive的网络协议会忽略此字段
			KeepAlive: t.KeepAlive,
		}).DialContext,
	}
}
