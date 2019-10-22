[![GoDoc](https://godoc.org/github.com/aberic/gnomon?status.svg)](https://godoc.org/github.com/aberic/gnomon)
[![Go Report Card](https://goreportcard.com/badge/github.com/aberic/gnomon)](https://goreportcard.com/report/github.com/aberic/gnomon)
[![GolangCI](https://golangci.com/badges/github.com/aberic/gnomon.svg)](https://golangci.com/r/github.com/aberic/gnomon)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/4f11995425294f42aec6a207b8aab367)](https://www.codacy.com/manual/aberic/gnomon?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=aberic/gnomon&amp;utm_campaign=Badge_Grade)
[![Travis (.org)](https://img.shields.io/travis/aberic/gnomon.svg?label=build)](https://www.travis-ci.org/aberic/gnomon)
[![Coveralls github](https://img.shields.io/coveralls/github/aberic/gnomon)](https://coveralls.io/github/aberic/gnomon?branch=master)

# Gnomon
通用编写go应用的公共库。

## 开发环境
* Go 1.12+
* Darwin/amd64

## 测试环境
* Go 1.11+
* Linux/x64

### 安装
``go get github.com/aberic/gnomon``

### 使用
```go
gnomon.Byte(). … // 字节
gnomon.Command(). … // 命令行
gnomon.Env(). … // 环境变量
gnomon.File(). … // 文件操作
gnomon.IP(). … // IP
gnomon.JWT(). … // JWT
gnomon.String(). … // 字符串
gnomon.CryptoHash(). … // Hash/散列
gnomon.CryptoRSA(). … // RSA
gnomon.CryptoECC(). … // ECC
gnomon.CryptoAES(). … // AES
gnomon.CryptoDES(). … // DES
gnomon.CA(). … // CA
gnomon.Log(). … // 日志
gnomon.Scale(). … // 算数/转换
gnomon.Time(). … // 时间
```

### 文档
参考 https://godoc.org/github.com/aberic/gnomon

<br><br>