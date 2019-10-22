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
 */

package gnomon

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

// CommandCommon 命令行工具
type CommandCommon struct{}

// ExecCommand 执行cmd命令
//
// commandName 命令执行文件名
//
// 命令后续参数以字符串数组的方式传入
//
// line 执行命令后输出总行数
//
// cmd cmd对象
//
// contentArray 执行命令后输出内容按行放入字符串数组
func (c *CommandCommon) ExecCommand(commandName string, params ...string) (line int, cmd *exec.Cmd, contentArray []string, err error) {
	var (
		stdout   io.ReadCloser
		stderr   io.ReadCloser
		bytesErr []byte
	)
	cmd = exec.Command(commandName, params...)
	//显示运行的命令
	Log().Debug("ExecCommand", Log().Field("cmd", strings.Join([]string{commandName, strings.Join(cmd.Args[1:], " ")}, " ")))
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
		return line, cmd, contentArray, nil
	}
ERR:
	Log().Error("ExecCommand", Log().Err(err))
	return 0, nil, nil, err
}

// ExecCommandSilent 执行cmd命令
//
// commandName 命令执行文件名
//
// 命令后续参数以字符串数组的方式传入
//
// line 执行命令后输出总行数
//
// cmd cmd对象
//
// contentArray 执行命令后输出内容按行放入字符串数组
func (c *CommandCommon) ExecCommandSilent(commandName string, params ...string) (line int, cmd *exec.Cmd, contentArray []string, err error) {
	var (
		stdout   io.ReadCloser
		stderr   io.ReadCloser
		bytesErr []byte
	)
	cmd = exec.Command(commandName, params...)
	//显示运行的命令
	Log().Debug("ExecCommand", Log().Field("cmd", strings.Join([]string{commandName, strings.Join(cmd.Args[1:], " ")}, " ")))
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

		if err = cmd.Wait(); nil != err {
			goto ERR
		}
		return line, cmd, contentArray, nil
	}
ERR:
	Log().Error("ExecCommand", Log().Err(err))
	return 0, nil, nil, err
}

// ExecCommandTail 实时打印执行脚本过程中的命令
//
// 命令后续参数以字符串数组的方式传入
//
// line 执行命令后输出总行数
//
// cmd cmd对象
//
// contentArray 执行命令后输出内容按行放入字符串数组
func (c *CommandCommon) ExecCommandTail(commandName string, params ...string) (line int, cmd *exec.Cmd, contentArray []string, err error) {
	var (
		stdout   io.ReadCloser
		stderr   io.ReadCloser
		bytesErr []byte
	)
	cmd = exec.Command(commandName, params...)
	//显示运行的命令
	Log().Debug("ExecCommandTail", Log().Field("cmd", strings.Join([]string{commandName, strings.Join(cmd.Args[1:], " ")}, " ")))
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
		return line, cmd, contentArray, nil
	}
ERR:
	Log().Error("ExecCommand", Log().Err(err))
	return 0, nil, nil, err
}

// CommandAsync 异步命令执行归属通道对象
type CommandAsync struct {
	Command *exec.Cmd
	Tail    string
	Err     error
}

// ExecCommandAsync 异步执行cmd命令
//
// commandAsync CommandAsync通道对象
//
// commandName 命令执行文件名
//
// 命令后续参数以字符串数组的方式传入
func (c *CommandCommon) ExecCommandAsync(commandAsync chan *CommandAsync, commandName string, params ...string) {
	var (
		stdout   io.ReadCloser
		stderr   io.ReadCloser
		bytesErr []byte
		err      error
	)
	ca := &CommandAsync{}
	cmd := exec.Command(commandName, params...)
	ca.Command = cmd
	commandAsync <- ca
	//显示运行的命令
	Log().Debug("ExecCommandAsync", Log().Field("cmd", strings.Join([]string{commandName, strings.Join(cmd.Args[1:], " ")}, " ")))
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
				ca.Tail = lineStr
				commandAsync <- ca
				break
			}
			ca.Tail = lineStr
			commandAsync <- ca
		}
		ca.Tail = "OFF"
		commandAsync <- ca
		if err = cmd.Wait(); nil != err {
			goto ERR
		}
	}
ERR:
	Log().Error("ExecCommandAsync", Log().Err(err))
	ca.Err = err
	commandAsync <- ca
}
