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

// Package log 日志操作工具
package log

import (
	"github.com/ennoo/rivet/utils/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"sync"
)

var (
	// Log 日志全局对象
	Log Logger
	// Common 通用包日志对象
	Common, _ = zap.NewDevelopment()
	// Discovery 发现服务包日志对象
	Discovery, _ = zap.NewDevelopment()
	// Examples 案例日志对象
	Examples, _ = zap.NewDevelopment()
	// Rivet 框架日志对象
	Rivet, _ = zap.NewDevelopment()
	// Server 关联接口服务日志对象
	Server, _ = zap.NewDevelopment()
	// Bow 网关日志对象
	Bow, _ = zap.NewDevelopment()
	// Shunt 负载均衡日志对象
	Shunt, _ = zap.NewDevelopment()
	// Trans 请求处理日志对象
	Trans, _ = zap.NewDevelopment()
	// Scheduled 定时任务日志对象
	Scheduled, _ = zap.NewDevelopment()
	// Self 当前使用该框架服务日志对象
	Self, _ = zap.NewDevelopment()
)

var (
	instance   *Logger
	once       sync.Once
	initialize sync.Once
	logPath    string
)

// Logger 日志入口对象
type Logger struct {
	Config *Config
}

// GetLogInstance 获取日志管理对象 Log 单例
func GetLogInstance() *Logger {
	once.Do(func() {
		logPath = env.GetEnvDefault(env.LogPath, "./logs")
		instance = &Logger{
			&Config{
				Level:      DebugLevel,
				MaxSize:    128,
				MaxBackups: 30,
				MaxAge:     30,
				Compress:   true,
			},
		}
		Common = instance.New(strings.Join([]string{logPath, "common.log"}, "/"), "common")
		Discovery = instance.New(strings.Join([]string{logPath, "discovery.log"}, "/"), "discovery")
		Examples = instance.New(strings.Join([]string{logPath, "examples.log"}, "/"), "examples")
		Rivet = instance.New(strings.Join([]string{logPath, "rivet.log"}, "/"), "rivet")
		Server = instance.New(strings.Join([]string{logPath, "server.log"}, "/"), "server")
		Bow = instance.New(strings.Join([]string{logPath, "bow.log"}, "/"), "bow")
		Shunt = instance.New(strings.Join([]string{logPath, "shunt.log"}, "/"), "shunt")
		Trans = instance.New(strings.Join([]string{logPath, "trans.log"}, "/"), "trans")
		Scheduled = instance.New(strings.Join([]string{logPath, "scheduled.log"}, "/"), "scheduled")
		Self = instance.New(strings.Join([]string{logPath, "self.log"}, "/"), "self")
	})
	return instance
}

// Conf 配置日志基本信息
func (log *Logger) Conf(config *Config) {
	log.Config = config
}

// Init 日志初始化操作
func (log *Logger) Init(logPath, serviceName string, config *Config, dev bool) {
	initialize.Do(func() {
		instance = &Logger{config}
		if dev {
			Common = instance.New(strings.Join([]string{logPath, "common.log"}, "/"), "common")
			Discovery = instance.New(strings.Join([]string{logPath, "discovery.log"}, "/"), "discovery")
			Examples = instance.New(strings.Join([]string{logPath, "examples.log"}, "/"), "examples")
			Rivet = instance.New(strings.Join([]string{logPath, "rivet.log"}, "/"), "rivet")
			Server = instance.New(strings.Join([]string{logPath, "server.log"}, "/"), "server")
			Bow = instance.New(strings.Join([]string{logPath, "bow.log"}, "/"), "bow")
			Shunt = instance.New(strings.Join([]string{logPath, "shunt.log"}, "/"), "shunt")
			Trans = instance.New(strings.Join([]string{logPath, "trans.log"}, "/"), "trans")
			Scheduled = instance.New(strings.Join([]string{logPath, "scheduled.log"}, "/"), "scheduled")
		}
		Self = instance.New(strings.Join([]string{logPath, "/", serviceName, ".log"}, ""), serviceName)
	})
}

func (log *Logger) Self(serviceName string) {
	Self = instance.New(strings.Join([]string{logPath, "/", serviceName, ".log"}, ""), "serviceName")
}

// New 新建日志对象
func (log *Logger) New1(name string, serviceName string) *zap.Logger {
	return log.NewCustom(strings.Join([]string{logPath, name, ".log"}, "/"), log.Config.Level, log.Config.MaxSize, log.Config.MaxBackups, log.Config.MaxAge, log.Config.Compress, serviceName)
}

// New 新建日志对象
func (log *Logger) New(filePath string, serviceName string) *zap.Logger {
	return log.NewCustom(filePath, log.Config.Level, log.Config.MaxSize, log.Config.MaxBackups, log.Config.MaxAge, log.Config.Compress, serviceName)
}

// NewCustom 新建自定义日志对象
func (log *Logger) NewCustom(filePath string, level Level, maxSize int, maxBackups int, maxAge int, compress bool, serviceName string) *zap.Logger {
	var zapLevel zapcore.Level
	switch log.Config.Level {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	case DPanicLevel:
		zapLevel = zapcore.DPanicLevel
	case PanicLevel:
		zapLevel = zapcore.PanicLevel
	case FatalLevel:
		zapLevel = zapcore.FatalLevel
	}
	core := newCore(filePath, zapLevel, maxSize, maxBackups, maxAge, compress)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", serviceName))
	// 返回构造日志
	return zap.New(core, caller, development, filed)
}

func newCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) zapcore.Core {
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
}
