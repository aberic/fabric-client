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

package utils

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ennoo/rivet/utils/log"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

var contentArray = make([]string, 0, 5)

// ExecCommand 执行cmd命令
//
// commandName 命令执行文件名
//
// 命令后续参数以字符串数组的方式传入
//
// 返回值
//
// int 执行命令后输出总行数
//
// []string 执行命令后输出内容按行放入字符串数组
func ExecCommand(commandName string, params ...string) (int, []string, error) {
	var (
		err      error
		stdout   io.ReadCloser
		stderr   io.ReadCloser
		bytesErr []byte
		line     int
	)
	contentArray = make([]string, 0)
	cmd := exec.Command(commandName, params...)
	//显示运行的命令
	log.Self.Info("exec: ", log.String("cmd", strings.Join([]string{commandName, strings.Join(cmd.Args[1:], " ")}, " ")))
	if stdout, err = cmd.StdoutPipe(); err != nil {
		goto ERR
	} else {
		if stderr, err = cmd.StderrPipe(); err != nil {
			goto ERR
		}
		// Start开始执行c包含的命令，但并不会等待该命令完成即返回。Wait方法会返回命令的返回状态码并在命令返回后释放相关的资源。
		if err = cmd.Start(); nil != err {
			goto ERR
		}

		if bytesErr, err = ioutil.ReadAll(stderr); err != nil {
			goto ERR
		} else if len(bytesErr) != 0 {
			err = errors.New("stderr is not nil: " + string(bytesErr))
			goto ERR
		}

		reader := bufio.NewReader(stdout)

		//实时循环读取输出流中的一行内容
		for {
			lineStr, err2 := reader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				contentArray = append(contentArray, lineStr)
				break
			}
			line++
			contentArray = append(contentArray, lineStr)
		}

		fmt.Println("========================")
		fmt.Println("line = ", line)
		fmt.Println("========================")
		fmt.Println(strings.Join(contentArray, ""))

		if err = cmd.Wait(); nil != err {
			goto ERR
		}
		return line, contentArray, nil
	}
ERR:
	log.Self.Error("error", log.Error(err))
	return 0, nil, err
}

// ExecCommandTail 实时打印执行脚本过程中的命令
func ExecCommandTail(commandName string, params ...string) (int, []string, error) {
	var (
		err      error
		stdout   io.ReadCloser
		stderr   io.ReadCloser
		bytesErr []byte
		line     int
	)
	contentArray = make([]string, 0)
	cmd := exec.Command(commandName, params...)
	//显示运行的命令
	fmt.Printf("exec: %s %s\n", commandName, strings.Join(cmd.Args[1:], " "))
	if stdout, err = cmd.StdoutPipe(); err != nil {
		goto ERR
	} else {
		if stderr, err = cmd.StderrPipe(); err != nil {
			goto ERR
		}
		// Start开始执行c包含的命令，但并不会等待该命令完成即返回。Wait方法会返回命令的返回状态码并在命令返回后释放相关的资源。
		if err = cmd.Start(); nil != err {
			goto ERR
		}

		if bytesErr, err = ioutil.ReadAll(stderr); err != nil {
			goto ERR
		} else if len(bytesErr) != 0 {
			err = errors.New(string(bytesErr))
			goto ERR
		}

		reader := bufio.NewReader(stdout)

		//实时循环读取输出流中的一行内容
		for {
			lineStr, err2 := reader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				contentArray = append(contentArray, lineStr)
				break
			}
			line++
			fmt.Println(lineStr)
			contentArray = append(contentArray, lineStr)
		}
		if err = cmd.Wait(); nil != err {
			goto ERR
		}
		return line, contentArray, nil
	}
ERR:
	errStr := err.Error()
	if strings.Contains(errStr, "error") || strings.Contains(errStr, "Error") ||
		strings.Contains(errStr, "ERROR") || strings.Contains(errStr, "fail") ||
		strings.Contains(errStr, "Fail") || strings.Contains(errStr, "FAIL") {
		log.Self.Error("error", log.Error(err))
		return 0, nil, err
	} else {
		strArr := strings.Split(errStr, "\n")
		fmt.Println(errStr)
		return len(strArr), strArr, nil
	}
}
