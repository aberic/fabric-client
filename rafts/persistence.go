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

// role
//
// 所有节点初始状态都是Follower角色
//
// 超时时间内没有收到Leader的请求则转换为Candidate进行选举
//
// Candidate收到大多数节点的选票则转换为Leader；发现Leader或者收到更高任期的请求则转换为Follower
//
// Leader在收到更高任期的请求后转换为Follower
//
// Raft把时间切割为任意长度的任期（term），每个任期都有一个任期号，采用连续的整数

package raft

import (
	"github.com/ennoo/fabric-client/config"
)

// persistence 所有角色都拥有的持久化的状态（在响应RPC请求之前变更且持久化的状态）
type persistence struct {
	leaderID    string                    // 当前任务Leader ID
	currentTerm int32                     // 服务器的任期，初始为0，递增
	votedFor    *votedFor                 // 在当前获得选票的候选人的 Id
	version     int32                     // 当前配置版本 index 递增
	configs     map[string]*config.Config // 当前term同步配置信息
}

type votedFor struct {
	id        string // 在当前获得选票的候选人的 Id
	term      int32  // 在当前获得选票的候选人的任期
	timestamp int64  // 在当前获取选票的候选人时间戳
}

//// storage 存储本地文件
//func (p *persistence) storage() {
//	filePath, _ := p.fileStorage()
//	if exist, err := file.PathExists(filePath); nil != err {
//		log.Self.Error("raft", log.Error(err))
//	} else if !exist {
//		lastIndex := strings.LastIndex(filePath, "/")
//		parentPath := filePath[0:lastIndex]
//		if err := os.MkdirAll(parentPath, os.ModePerm); nil != err {
//			log.Self.Error("raft", log.Error(err))
//			return
//		}
//		// 创建文件，如果文件已存在，会将文件清空
//		if f, err := os.Create(filePath); nil != err {
//			log.Self.Error("raft", log.Error(err))
//			return
//		} else {
//			p.save(f)
//		}
//	} else {
//		// //表示最佳的方式打开文件，如果不存在就创建，打开的模式是可读可写，权限是644
//		if f, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0644); nil != err {
//			log.Self.Error("raft", log.Error(err))
//			return
//		} else {
//			p.save(f)
//		}
//	}
//}
//
//func (p *persistence) loadStorage() (configs map[string]*config.Config, err error) {
//	var (
//		exist bool
//	)
//	filePath, line := p.fileStorage()
//	if exist, err = file.PathExists(filePath); nil != err {
//		return
//	} else if !exist {
//		err = errors.New("file not exist")
//
//		return nil, nil
//	}
//	data := p.readLine(filePath, line)
//	//byteLineChange := []byte("\n")
//	//data = data[0 : len(data)-len(byteLineChange)]
//	if nil == data {
//		err = errors.New("file read nil")
//		return
//	}
//	err = yaml.Unmarshal(data, configs)
//	if nil == err {
//		p.configs = configs
//	}
//	return
//}
//
//func (p *persistence) readLine(filename string, ling int32) []byte {
//	f, _ := os.Open(filename)
//	defer p.fileClose(f)
//	r := bufio.NewReader(f)
//	var lineLocal int32
//	lineLocal = 0
//	for {
//		if lineLocal == ling {
//			data, err := r.ReadBytes('\n')
//			if err == io.EOF {
//				return data
//			}
//			return data
//		} else {
//			_, err := r.ReadBytes('\n')
//			if err == io.EOF {
//				break
//			}
//			lineLocal++
//		}
//	}
//	return nil
//}
//
//func (p *persistence) save(f *os.File) {
//	defer p.fileClose(f)
//	data, err := yaml.Marshal(p.configs)
//	if nil != err {
//		log.Self.Error("raft", log.Error(err))
//		return
//	}
//	if ret, err := f.Seek(io.SeekStart, io.SeekEnd); nil != err {
//		log.Self.Error("raft", log.Error(err))
//		return
//	} else {
//		byteLineChange := []byte("\n")
//		for _, b := range byteLineChange {
//			data = append(data, b)
//		}
//		if _, err := f.WriteAt(data, ret); nil != err {
//			log.Self.Error("raft", log.Error(err))
//			return
//		}
//	}
//}
//
//func (p *persistence) fileClose(f *os.File) {
//	if err := f.Close(); nil != err {
//		log.Self.Error("raft", log.Error(err))
//		return
//	}
//}
//
//func (p *persistence) fileStorage() (filePath string, line int32) {
//	mod := p.version % 1000
//	return strings.Join(
//		[]string{
//			logPathLocal,
//			"/",
//			strconv.Itoa(int(p.currentTerm)),
//			"/",
//			strconv.Itoa(int(mod + 1)),
//			".raft",
//		},
//		"",
//	), p.version - mod*1000
//}
