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

// Package env 环境变量工具
package env

import (
	"github.com/ennoo/rivet/utils/string"
	"os"
	"strings"
)

const (
	// LogPath 日志文件输出路径
	LogPath = "LOG_PATH"
	// ServiceName 当前启动服务名
	ServiceName = "SERVICE_NAME"
	// Port 当前服务启动端口号
	Port = "PORT"
	// HealthCheck 是否开启健康检查
	HealthCheck = "HEALTH_CHECK"
	// ServerManager 是否启用服务管理功能
	ServerManager = "SERVER_MANAGER"
	// LoadBalance 是否启用负载均衡
	LoadBalance = "LOAD_BALANCE"
	// OpenTLS 是否开启 TLS
	OpenTLS = "OPEN_TLS"
	// ConfigPath Bow配置文件路径
	ConfigPath = "CONFIG_PATH"
	// DiscoveryInit 是否启用发现服务
	DiscoveryInit = "DISCOVERY_INIT"
	// DiscoveryComponent 所启用发现服务组件名
	DiscoveryComponent = "DISCOVERY_COMPONENT"
	// DiscoveryURL 当前服务注册的发现服务地址
	DiscoveryURL = "DISCOVERY_URL"
	// DiscoveryReceiveHost 发现服务收到当前注册服务的地址
	DiscoveryReceiveHost = "DISCOVERY_RECEIVE_HOST"
	// GOPath Go工作路径
	GOPath = "GOPATH"
	// DBUrl 数据库 URL
	DBUrl = "DB_URL"
	// DBName 数据库名称
	DBName = "DB_NAME"
	// DBUser 数据库用户名
	DBUser = "DB_USER"
	// DBPass 数据库用户密码
	DBPass = "DB_PASS"
)

// GetEnv 获取环境变量 envName 的值
//
// envName 环境变量名称
func GetEnv(envName string) string {
	return os.Getenv(envName)
}

// GetEnvBool 获取环境变量 envName 的 bool 值
//
// envName 环境变量名称
func GetEnvBool(envName string) bool {
	return strings.EqualFold(os.Getenv(envName), "true")
}

// GetEnvDefault 获取环境变量 envName 的值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func GetEnvDefault(envName string, defaultValue string) string {
	env := GetEnv(envName)
	if str.IsEmpty(env) {
		return defaultValue
	}
	return env
}

// GetEnvBoolDefault 获取环境变量 envName 的 bool 值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func GetEnvBoolDefault(envName string, defaultValue bool) bool {
	env := GetEnv(envName)
	if str.IsEmpty(env) {
		return defaultValue
	}
	return strings.EqualFold(os.Getenv(envName), "true")
}
