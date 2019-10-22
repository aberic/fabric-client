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
	"archive/zip"
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FileCommon 文件操作工具
type FileCommon struct{}

// PathExists 判断路径是否存在
func (f *FileCommon) PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ReadFirstLine 从文件中读取第一行并返回字符串数组
func (f *FileCommon) ReadFirstLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()
	finReader := bufio.NewReader(file)
	inputString, _ := finReader.ReadString('\n')
	return String().TrimN(inputString), nil
}

// ReadPointLine 从文件中读取指定行并返回字符串数组
func (f *FileCommon) ReadPointLine(filePath string, line int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()
	finReader := bufio.NewReader(file)
	lineCount := 1
	for {
		inputString, err := finReader.ReadString('\n')
		//fmt.Println(inputString)
		if err == io.EOF {
			if lineCount == line {
				return inputString, nil
			}
			return "", errors.New("index out of line count")
		}
		if lineCount == line {
			return inputString, nil
		}
		lineCount++
	}
}

// ReadLines 从文件中逐行读取并返回字符串数组
func (f *FileCommon) ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	finReader := bufio.NewReader(file)
	var fileList []string
	for {
		inputString, err := finReader.ReadString('\n')
		//fmt.Println(inputString)
		if err == io.EOF {
			fileList = append(fileList, String().TrimN(inputString))
			break
		}
		fileList = append(fileList, String().TrimN(inputString))
	}
	//fmt.Println("fileList",fileList)
	return fileList, nil
}

// ParentPath 文件父路径
func (f *FileCommon) ParentPath(filePath string) string {
	return filePath[0:strings.LastIndex(filePath, "/")]
}

// Append 追加内容到文件中
//
// filePath 文件地址
//
// data 内容
//
// force 如果文件已存在，会将文件清空
//
// It returns the number of bytes written and an error
func (f *FileCommon) Append(filePath string, data []byte, force bool) (int, error) {
	var (
		file *os.File
		n    int
		err  error
	)
	exist := f.PathExists(filePath)
	if exist {
		if force {
			// 创建文件，如果文件已存在，会将文件清空
			if file, err = os.Create(filePath); err != nil {
				return 0, err
			}
		} else {
			if file, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0644); nil != err {
				return 0, err
			}
		}
	} else {
		parentPath := f.ParentPath(filePath)
		if err = os.MkdirAll(parentPath, os.ModePerm); nil != err {
			return 0, err
		}
		if file, err = os.Create(filePath); err != nil {
			return 0, err
		}
	}
	defer func() {
		_ = file.Close()
	}()
	// 将数据写入文件中
	//file.WriteString(string(data)) //写入字符串
	if n, err = file.Write(data); nil != err { // 写入byte的slice数据
		return 0, err
	}
	return n, nil
}

// Modify 修改文件中指定位置的内容
//
// filePath 文件地址
//
// offset 以0为起始坐标的偏移量
//
// data 内容
//
// force 如果文件已存在，会将文件清空
//
// It returns the number of bytes written and an error
func (f *FileCommon) Modify(filePath string, offset int64, data []byte, force bool) (int, error) {
	var (
		file *os.File
		n    int
		err  error
	)
	exist := f.PathExists(filePath)
	if exist {
		if force {
			// 创建文件，如果文件已存在，会将文件清空
			if file, err = os.Create(filePath); err != nil {
				return 0, err
			}
		} else {
			if file, err = os.OpenFile(filePath, os.O_RDWR, 0644); nil != err {
				return 0, err
			}
		}
	} else {
		parentPath := f.ParentPath(filePath)
		if err = os.MkdirAll(parentPath, os.ModePerm); nil != err {
			return 0, err
		}
		if file, err = os.Create(filePath); err != nil {
			return 0, err
		}
	}
	defer func() {
		_ = file.Close()
	}()
	// 以0为起始坐标偏移指标到指定位置
	if _, err = file.Seek(offset, io.SeekStart); nil != err {
		return 0, err
	}
	// 将数据写入文件中
	//file.WriteString(string(data)) //写入字符串
	if n, err = file.Write(data); nil != err { // 写入byte的slice数据
		return 0, err
	}
	return n, nil
}

// LoopDirs 遍历文件夹下的所有子文件夹
func (f *FileCommon) LoopDirs(pathname string) ([]string, error) {
	var s []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		Log().Debug("read dir fail", Log().Err(err))
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

// LoopFiles 遍历文件夹及子文件夹下的所有文件
func (f *FileCommon) LoopFiles(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		Log().Debug("read dir fail", Log().Err(err))
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = f.LoopFiles(fullDir, s)
			if err != nil {
				Log().Debug("read dir fail", Log().Err(err))
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

// Compress 压缩文件
// files 文件数组，可以是不同dir下的文件或者文件夹
// dest 压缩文件存放地址
func (f *FileCommon) Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer func() { _ = d.Close() }()
	w := zip.NewWriter(d)
	defer func() { _ = w.Close() }()
	for _, file := range files {
		err := f.compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileCommon) compress(file *os.File, prefix string, zw *zip.Writer) error {
	var (
		info   os.FileInfo
		header *zip.FileHeader
		writer io.Writer
		err    error
	)
	defer func() { _ = file.Close() }()
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
			fil, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = f.compress(fil, prefix, zw)
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
		if err != nil {
			return err
		}
	}
	return nil
}

// DeCompressTar 压缩文件
// 压缩文件路径
// 解压文件夹
func (f *FileCommon) DeCompressTar(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		Log().Error("DeCompressTar", Log().Err(err))
		return err
	}
	defer func() { _ = srcFile.Close() }()
	reader, err := zip.OpenReader(srcFile.Name())
	if nil != err {
		return err
	}
	return f.deCompress(reader, dest)
}

// DeCompressZip 解压
func (f *FileCommon) DeCompressZip(zipFile, dest string) error {
	var (
		reader *zip.ReadCloser
		err    error
	)
	if reader, err = zip.OpenReader(zipFile); nil != err {
		Log().Error("DeCompressZip", Log().Err(err))
		return err
	}
	return f.deCompress(reader, dest)
}

// deCompress 压缩文件
func (f *FileCommon) deCompress(reader *zip.ReadCloser, dest string) error {
	defer func() { _ = reader.Close() }()
	for _, innerFile := range reader.File {
		info := innerFile.FileInfo()
		if info.IsDir() {
			err := os.MkdirAll(innerFile.Name, os.ModePerm)
			if err != nil {
				Log().Error("deCompress1", Log().Err(err))
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
			Log().Error("deCompress2", Log().Err(err))
			return err
		}
		filePath := filepath.Join(dest, innerFile.Name)
		if exist := f.PathExists(filePath); !exist {
			lastIndex := strings.LastIndex(filePath, "/")
			parentPath := filePath[0:lastIndex]
			if err := os.MkdirAll(parentPath, os.ModePerm); nil != err {
				Log().Error("deCompress3", Log().Err(err))
				return err
			}
		}
		newFile, err := os.Create(filePath)
		if err != nil {
			Log().Error("deCompress4", Log().Err(err))
			_ = srcFile.Close()
			continue
		}
		if _, err = io.Copy(newFile, srcFile); nil != err {
			_ = newFile.Close()
			_ = srcFile.Close()
			Log().Error("deCompress5", Log().Err(err))
			return err
		}
		_ = newFile.Close()
		_ = srcFile.Close()
	}
	return nil
}
