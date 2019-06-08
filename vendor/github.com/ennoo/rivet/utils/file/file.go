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

// Package file 文件操作工具
package file

import (
	"bufio"
	"github.com/ennoo/rivet/utils/string"
	"io"
	"os"
)

// PathExists 判断路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ReadFileFirstLine 从文件中逐行读取并返回字符串数组
func ReadFileFirstLine(filePath string) (string, error) {
	fileIn, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fileIn.Close()
	finReader := bufio.NewReader(fileIn)
	inputString, _ := finReader.ReadString('\n')
	return str.TrimN(inputString), nil
}

// ReadFileByLine 从文件中逐行读取并返回字符串数组
func ReadFileByLine(filePath string) ([]string, error) {
	fileIn, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fileIn.Close()
	finReader := bufio.NewReader(fileIn)
	var fileList []string
	for {
		inputString, err := finReader.ReadString('\n')
		//fmt.Println(inputString)
		if err == io.EOF {
			fileList = append(fileList, str.TrimN(inputString))
			break
		}
		fileList = append(fileList, str.TrimN(inputString))
	}
	//fmt.Println("fileList",fileList)
	return fileList, nil
}
