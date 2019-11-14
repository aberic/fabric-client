/*
 * Copyright (c) 2019. aberic - All Rights Reserved.
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

package gnomon

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// debugLevel logs are typically voluminous, and are usually disabled in
	// production.
	debugLevel Level = iota - 1
	// infoLevel is the default logging priority.
	infoLevel
	// warnLevel logs are more important than Info, but don't need individual
	// human review.
	warnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	errorLevel
	// panicLevel logs a message, then panics.
	panicLevel
	// fatalLevel logs a message, then calls os.Exit(1).
	fatalLevel
	allLevel
)

const (
	logNameDebug = "DEBUG"
	logNameInfo  = "INFO "
	logNameWarn  = "WARN "
	logNameError = "ERROR"
	logNamePanic = "PANIC"
	logNameFatal = "FATAL"
)

// LogCommon 日志工具
type LogCommon struct {
	logDir           string           // logDir 日志文件目录
	maxSizeByte      int64            // maxSizeByte 每个日志文件保存的最大尺寸 单位：byte
	maxAge           int              // maxAge 文件最多保存多少天
	files            map[Level]*filed // files 日志文件输入io对象集合
	level            Level            // level 日志级别
	production       bool             // 生产环境，该模式下控制台不会输出任何日志
	utc              bool             // CST & UTC 时间
	date             string           // date 当前日志文件后缀日期
	mkRootDirSuccess bool             // mkRootDirSuccess 是否成功初始化log对象
	job              *cron.Cron       // job 日志定时清理任务
	once             sync.Once        // once log对象只会被初始化一次
}

// Init log初始化
//
// logDir 日志文件目录
//
// maxSize 每个日志文件保存的最大尺寸 单位：M
//
// maxAge 文件最多保存多少天
//
// utc CST & UTC 时间
func (l *LogCommon) Init(logDir string, maxSize, maxAge int, utc bool) error {
	l.Info("log service init")
	var errInit error
	l.once.Do(func() {
		if String().IsEmpty(logDir) {
			logDir = "./tmp/log"
		}
		if err := os.MkdirAll(logDir, os.ModePerm); nil != err {
			l.Error("log service init error", Log().Err(err))
			errInit = err
			return
		}
		l.mkRootDirSuccess = true
		l.logDir = logDir
		l.utc = utc
		l.maxSizeByte = int64(maxSize * 1024 * 1024)
		l.maxAge = maxAge
		l.level = debugLevel
		l.production = false
		l.files = map[Level]*filed{
			debugLevel: {fileIndex: "0", tasks: make(chan string, 1000)},
			infoLevel:  {fileIndex: "0", tasks: make(chan string, 1000)},
			warnLevel:  {fileIndex: "0", tasks: make(chan string, 1000)},
			errorLevel: {fileIndex: "0", tasks: make(chan string, 1000)},
			panicLevel: {fileIndex: "0", tasks: make(chan string, 1000)},
			fatalLevel: {fileIndex: "0", tasks: make(chan string, 1000)},
			allLevel:   {fileIndex: "0", tasks: make(chan string, 1000)},
		}
		if utc {
			l.date = time.Now().UTC().Format("20060102")
		} else {
			l.date = time.Now().Local().Format("20060102")
		}
		l.job = cron.New()
		go l.checkMaxAge()
	})
	return errInit
}

// checkMaxAge 遍历并检查文件是否达到保存天数，达到则删除
func (l *LogCommon) checkMaxAge() {
	// 每隔5秒执行一次：*/5 * * * * ?
	// 每隔1分钟执行一次：0 */1 * * * ?
	// 每天23点执行一次：0 0 23 * * ?
	// 每天凌晨1点执行一次：0 0 1 * * ?
	// 每月1号凌晨1点执行一次：0 0 1 1 * ?
	// 在26分、29分、33分执行一次：0 26,29,33 * * * ?
	// 每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
	err := l.job.AddFunc(strings.Join([]string{"0 0 0 */", strconv.Itoa(l.maxAge), " * ?"}, ""), func() {
		var timeDate string
		if l.utc {
			timeDate = time.Now().UTC().Format("20060102")
		} else {
			timeDate = time.Now().Local().Format("20060102")
		}
		logDirs, _ := File().LoopDirs(l.logDir)
		for _, dirName := range logDirs {
			if strings.Contains(dirName, timeDate) {
				_ = os.RemoveAll(dirName)
			}
		}
	})
	if nil != err {
		time.Sleep(time.Second)
		l.checkMaxAge()
	} else {
		l.job.Start()
	}
}

// Set 设置日志可变属性
//
// level 日志级别(debugLevel/infoLevel/warnLevel/ErrorLevel/panicLevel/fatalLevel)
//
// production 是否生产环境，在生产环境下控制台不会输出任何日志
func (l *LogCommon) Set(level Level, production bool) {
	l.level = level
	l.production = production
}

// DebugLevel logs are typically voluminous, and are usually disabled in production.
func (l *LogCommon) DebugLevel() Level {
	return debugLevel
}

// InfoLevel is the default logging priority.
func (l *LogCommon) InfoLevel() Level {
	return infoLevel
}

// WarnLevel logs are more important than Info, but don't need individual human review.
func (l *LogCommon) WarnLevel() Level {
	return warnLevel
}

// ErrorLevel logs are high-priority. If an application is running smoothly,
// it shouldn't generate any error-level logs.
func (l *LogCommon) ErrorLevel() Level {
	return errorLevel
}

// PanicLevel logs a message, then panics.
func (l *LogCommon) PanicLevel() Level {
	return panicLevel
}

// FatalLevel logs a message, then panics.
func (l *LogCommon) FatalLevel() Level {
	return fatalLevel
}

// Debug 输出指定级别日志
func (l *LogCommon) Debug(msg string, fields ...*Field) {
	if l.level > debugLevel {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		l.logStandard(file, logNameDebug, msg, line, ok, debugLevel, fields...)
	} else {
		l.Warn("log recovery fail")
	}
}

// Info 输出指定级别日志
func (l *LogCommon) Info(msg string, fields ...*Field) {
	if l.level > infoLevel {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		l.logStandard(file, logNameInfo, msg, line, ok, infoLevel, fields...)
	} else {
		l.Warn("log recovery fail")
	}
}

// Warn 输出指定级别日志
func (l *LogCommon) Warn(msg string, fields ...*Field) {
	if l.level > warnLevel {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		l.logStandard(file, logNameWarn, msg, line, ok, warnLevel, fields...)
	} else {
		l.Warn("log recovery fail")
	}
}

// Error 输出指定级别日志
func (l *LogCommon) Error(msg string, fields ...*Field) {
	if l.level > errorLevel {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		l.logStandard(file, logNameError, msg, line, ok, errorLevel, fields...)
	} else {
		l.Warn("log recovery fail")
	}
}

// Panic 输出指定级别日志
func (l *LogCommon) Panic(msg string, fields ...*Field) {
	if l.level > panicLevel {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		l.logStandard(file, logNamePanic, msg, line, ok, panicLevel, fields...)
	} else {
		l.Warn("log recovery fail")
	}
}

// Fatal 输出指定级别日志
func (l *LogCommon) Fatal(msg string, fields ...*Field) {
	if l.level > fatalLevel {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		l.logStandard(file, logNameFatal, msg, line, ok, fatalLevel, fields...)
	} else {
		l.Warn("log recovery fail")
	}
}

// Field 自定义输出KV对象
func (l *LogCommon) Field(key string, value interface{}) *Field {
	return &Field{key: key, value: value}
}

// Err 自定义输出错误
func (l *LogCommon) Err(err error) *Field {
	if nil != err {
		return &Field{key: "error", value: err.Error()}
	}
	return &Field{key: "error", value: nil}
}

// logStandard 将日志输出到控制台
//
// file 日志触发所在文件地址
//
// levelName 日志级别名称
//
// msg 日志默认输出信息
//
// line 日志触发所在文件的行号
//
// ok 如果无法恢复信息，则为false
//
// level 日志级别
//
// fields 日志输出对象子集
func (l *LogCommon) logStandard(file, levelName, msg string, line int, ok bool, level Level, fields ...*Field) {
	var (
		fileString  string
		timeString  string
		zoneName    string
		stackString string
	)
	timeNow := time.Now()
	if l.utc {
		timeString = timeNow.UTC().Format("2006-01-02 15:04:05.000000")
		zoneName, _ = timeNow.UTC().Zone()
	} else {
		timeString = timeNow.Local().Format("2006-01-02 15:04:05.000000")
		zoneName, _ = timeNow.Local().Zone()
	}
	timeString = strings.Join([]string{timeString, zoneName}, " ")
	logArr := strings.Split(strings.Join([]string{file, strconv.Itoa(line)}, ":"), "/go/src/")
	if len(logArr) > 1 {
		fileString = logArr[1]
	} else {
		fileString = logArr[0]
	}
	if !l.production {
		var (
			commandJSON []byte
			err         error
		)
		logCommand := make(map[string]interface{})
		logCommand["msg"] = msg
		for _, field := range fields {
			if nil == field {
				continue
			}
			logCommand[field.key] = field.value
		}
		if commandJSON, err = json.Marshal(logCommand); nil != err {
			l.Error("json Marshal error", Log().Err(err))
			return
		}
		commandString := string(commandJSON)
		fmt.Println(timeString, levelName, fileString, commandString)
		switch levelName {
		case logNameError:
			stackString = string(debug.Stack())
			fmt.Println(stackString)
		case logNamePanic:
			stackString = string(debug.Stack())
			fmt.Println(stackString)
			if nil == l.files {
				panic(commandString)
			}
		case logNameFatal:
			stackString = string(debug.Stack())
			fmt.Println(stackString)
			if nil == l.files {
				os.Exit(1)
			}
		}
	}
	if nil == l.files {
		return
	}
	go l.logFile(timeString, fileString, stackString, levelName, msg, level, fields...)
}

// logFile 将日志内容输入文件中存储
//
// timeString 日志时间
//
// fileString 日志触发所在文件所在行信息
//
// stackString 日志堆栈信息
//
// levelName 日志级别名称
//
// msg 日志默认输出信息
//
// level 日志级别
//
// fields 日志输出对象子集
func (l *LogCommon) logFile(timeString, fileString, stackString, levelName, msg string, level Level, fields ...*Field) {
	var (
		mapJSON     []byte
		printString string
		err         error
		fd          *filed
	)
	logMap := make(map[string]interface{})
	logMap["level"] = strings.ToLower(levelName)
	logMap["time"] = timeString
	logMap["file"] = fileString
	logMap["msg"] = msg
	for _, field := range fields {
		if nil == field {
			continue
		}
		logMap[field.key] = field.value
	}
	if mapJSON, err = json.Marshal(logMap); nil != err {
		l.Error("json Marshal error", Log().Err(err))
		return
	}
	switch levelName {
	case logNameError, logNamePanic, logNameFatal:
		if String().IsEmpty(stackString) {
			stackString = string(debug.Stack())
		}
		printString = strings.Join([]string{string(mapJSON), stackString}, "\n")
	default:
		printString = strings.Join([]string{string(mapJSON), "\n"}, "")
	}
	if fd, err = l.useFiled(level, printString); nil == err {
		fd.tasks <- printString
	}
	if fd, err = l.useFiled(allLevel, printString); nil == err {
		fd.tasks <- printString
	}
}

// useFiled 使用日志文件
//
// level 日志级别
//
// printString 输出字符串
func (l *LogCommon) useFiled(level Level, printString string) (fd *filed, err error) {
	if fd = l.files[level]; fd.file == nil {
		defer fd.lock.Unlock()
		fd.lock.Lock()
		if fd.file == nil {
			var f *os.File
			if f, err = os.OpenFile(l.logFilePath(fd, level), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); nil != err {
				return
			}
			fd.file = f
			if err = l.checkFiled(level, fd, int64(len(printString)), false); nil != err {
				return
			}
			go fd.running()
			return
		}
	}
	if err = l.checkFiled(level, fd, int64(len(printString)), true); nil != err {
		return
	}
	return
}

// checkFiled 检出日志文件
//
// 如果当前正在使用的日志文件已经达到单个文件大小上限，则通过后缀++的方式将内容写入新的文件中
//
// level 日志级别
//
// fd 日志文件操作对象
//
// printStringLength 输出到文件中字节数长度
//
// lock 该操作是否需要给filed文件对象上锁。如果是复用对象，则需要上锁；如果是新建对象，则新建过程中本身就已经上锁，此处无需锁定
func (l *LogCommon) checkFiled(level Level, fd *filed, printStringLength int64, lock bool) (err error) {
	var ret int64
	if ret, err = fd.file.Seek(0, io.SeekEnd); nil != err {
		return
	}
	if l.maxSizeByte-ret-printStringLength < 0 {
		if lock {
			defer fd.lock.Unlock()
			fd.lock.Lock()
			return l.findAvailableFile(level, fd, printStringLength)
		}
		return l.findAvailableFile(level, fd, printStringLength)
	}
	return
}

// findAvailableFile 查找可用日志文件
//
// 如果当前正在使用的日志文件已经达到单个文件大小上限，则通过后缀++的方式将内容写入新的文件中
//
// level 日志级别
//
// fd 日志文件操作对象
//
// printStringLength 输出到文件中字节数长度
func (l *LogCommon) findAvailableFile(level Level, fd *filed, printStringLength int64) (err error) {
	var (
		ret  int64
		pass bool
	)
	for !pass {
		if ret, err = fd.file.Seek(0, io.SeekEnd); nil != err {
			return
		}
		if l.maxSizeByte-ret-printStringLength < 0 {
			index, _ := strconv.Atoi(fd.fileIndex)
			fd.fileIndex = strconv.Itoa(index + 1)
			if fd.file, err = os.OpenFile(l.logFilePath(fd, level), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); nil != err {
				return
			}
			continue
		}
		pass = true
	}
	return nil
}

// path 日志文件路径
func (l *LogCommon) logFilePath(fd *filed, level Level) string {
	parentPath := filepath.Join(l.logDir, l.date)
	if exist := File().PathExists(parentPath); !exist {
		if err := os.MkdirAll(parentPath, os.ModePerm); nil != err {
			l.Error("path mkdirAll error", Log().Err(err))
			return ""
		}
	}
	return filepath.Join(parentPath, l.levelFileName(fd, level))
}

// levelFileName 包含日志类型的日志文件名称
func (l *LogCommon) levelFileName(fd *filed, level Level) string {
	switch level {
	case debugLevel:
		return l.logFileName("debug_", fd.fileIndex)
	case infoLevel:
		return l.logFileName("info_", fd.fileIndex)
	case warnLevel:
		return l.logFileName("warn_", fd.fileIndex)
	case errorLevel:
		return l.logFileName("error_", fd.fileIndex)
	case panicLevel:
		return l.logFileName("panic_", fd.fileIndex)
	case fatalLevel:
		return l.logFileName("fatal_", fd.fileIndex)
	}
	return l.logFileName("log_", fd.fileIndex)
}

// logFileName 不包含日志类型的日志文件名称
func (l *LogCommon) logFileName(name, index string) string {
	return strings.Join([]string{name, l.date, "-", index, ".log"}, "")
}

// filed 日志文件操作对象
type filed struct {
	fileIndex string // fileIndex 日志文件相同日期编号，根据文件新建规则确定
	file      *os.File
	tasks     chan string // 任务队列，默认1000个缓存
	lock      sync.Mutex  // lock 每次做io开销的安全锁
}

// running 循环执行文件写入，默认60秒超时
func (f *filed) running() {
	to := time.NewTimer(60 * time.Second)
	for {
		select {
		case task := <-f.tasks:
			to.Reset(time.Second)
			if _, err := f.file.WriteString(task); nil != err {
				panic(err)
			}
		case <-to.C:
			_ = f.file.Close()
			f.file = nil
			return
		}
	}
}

// Field 日志输出子集对象
type Field struct {
	key   string
	value interface{}
}
