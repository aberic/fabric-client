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
	"archive/zip"
	"bufio"
	"errors"
	"github.com/ennoo/rivet/utils/log"
	"github.com/ennoo/rivet/utils/string"
	"io"
	"os"
	"path/filepath"
	"strings"
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

// CreateAndWrite 创建并写入内容到文件中
func CreateAndWrite(filePath string, data []byte, force bool) error {
	if exist, _ := PathExists(filePath); exist && !force {
		return errors.New("file exist")
	}
	lastIndex := strings.LastIndex(filePath, "/")
	parentPath := filePath[0:lastIndex]
	if err := os.MkdirAll(parentPath, os.ModePerm); nil != err {
		return err
	}
	// 创建文件，如果文件已存在，会将文件清空
	if file, err := os.Create(filePath); err != nil {
		return err
	} else {
		defer file.Close()
		// 将数据写入文件中
		//file.WriteString(string(data)) //写入字符串
		if n, err := file.Write(data); nil != err { // 写入byte的slice数据
			return err
		} else {
			log.Rivet.Debug("write", log.Int("byte count", n))
			return nil
		}
	}
}

//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer func() { _ = w.Close() }()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	var (
		info   os.FileInfo
		header *zip.FileHeader
		writer io.Writer
		err    error
	)
	info, err = file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		if header, err = zip.FileInfoHeader(info); nil != err {
			return err
		}
		header.Name = prefix + "/" + header.Name
		if writer, err = zw.CreateHeader(header); nil != err {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

/**
@tarFile：压缩文件路径
@dest：解压文件夹
*/
func DeCompressTar(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		log.Self.Error("DeCompressTar", log.Error(err))
		return err
	}
	defer srcFile.Close()
	var reader *zip.ReadCloser
	if reader, err = zip.OpenReader(srcFile.Name()); nil != err {
		return err
	}
	return deCompress(reader, dest)
}

//解压
func DeCompressZip(zipFile, dest string) error {
	var (
		reader *zip.ReadCloser
		err    error
	)
	if reader, err = zip.OpenReader(zipFile); nil != err {
		log.Self.Error("DeCompressZip", log.Error(err))
		return err
	}
	return deCompress(reader, dest)
}

/**
@zipFile：压缩文件
*/
func deCompress(reader *zip.ReadCloser, dest string) error {
	defer func() { _ = reader.Close() }()
	for _, innerFile := range reader.File {
		info := innerFile.FileInfo()
		if info.IsDir() {
			err := os.MkdirAll(innerFile.Name, os.ModePerm)
			if err != nil {
				log.Self.Error("deCompress1", log.Error(err))
				return err
			}
			continue
		}
		srcFile, err := innerFile.Open()
		if err != nil {
			continue
		}
		err = os.MkdirAll(dest, 0755)
		if err != nil {
			log.Self.Error("deCompress2", log.Error(err))
			return err
		}
		filePath := filepath.Join(dest, innerFile.Name)
		if exist, err := PathExists(filePath); !exist || nil != err {
			lastIndex := strings.LastIndex(filePath, "/")
			parentPath := filePath[0:lastIndex]
			if err := os.MkdirAll(parentPath, os.ModePerm); nil != err {
				log.Self.Error("deCompress3", log.Error(err))
				return err
			}
		}
		newFile, err := os.Create(filePath)
		if err != nil {
			log.Self.Error("deCompress4", log.Error(err))
			srcFile.Close()
			continue
		}
		if _, err = io.Copy(newFile, srcFile); nil != err {
			newFile.Close()
			srcFile.Close()
			log.Self.Error("deCompress5", log.Error(err))
			return err
		}
		newFile.Close()
		srcFile.Close()
	}
	return nil
}
