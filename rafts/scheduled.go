/*
 * Copyright (c) 2019.. Aberic - All Rights Reserved.
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

package raft

import (
	"github.com/ennoo/rivet/utils/log"
	"github.com/robfig/cron"
	"strings"
	"time"
)

const (
	timeout = 3000 // Raft心跳超时ms
)

type scheduled struct {
	// Raft服务
	raft *Raft
	// Raft关闭心跳服务
	tickerEnd chan int8
	// 检查Raft状态定时任务
	checkCron *cron.Cron
	// 检查Raft计时任务可用性通道
	checkErr chan error
	// 释放关闭
	checkRelease chan int8
	// 心跳定时任务
	ticker *time.Ticker
	// 最后一次接收到心跳时间戳ms
	time int64
}

// start 启动定时检查raft计时任务可用性任务
func (s *scheduled) start() {
	log.Self.Info("raft", log.Reflect("start", s.raft.nodes))
	s.tickerEnd = make(chan int8, 1)
	s.checkCron = cron.New()
	s.checkErr = make(chan error, 1)
	s.checkRelease = make(chan int8, 1)
	// 如果不异步，会阻塞主线程
	go s.check()
}

// check 定时调用此方法检查raft计时任务可用性
func (s *scheduled) check() {
	go s.task()
Loop:
	for {
		select {
		case err := <-s.checkErr:
			log.Self.Error("raft", log.String("check err", "reStartCheck"), log.Error(err))
			s.check()
		case <-s.checkRelease:
			break Loop
		}
	}
}

// task 心跳及Raft状态检查定时任务方法
func (s *scheduled) task() {
	s.checkCron.Stop()
	err := s.checkCron.AddFunc(strings.Join([]string{"*/3 * * * * ?"}, ""), func() {
		if s.raft.role.role() == roleFollower && time.Now().UnixNano()/1e6-s.time > timeout { // 如果自身是follower节点
			log.Self.Debug("raft", log.Int32("Term", s.raft.term), log.String("task", "follower timeout"))
			s.raft.role.candidate()
		}
	})
	if nil != err {
		s.checkErr <- err
	} else {
		s.checkCron.Start()
	}
}

// tickerStart 启动心跳定时任务
func (s *scheduled) tickerStart() {
	log.Self.Debug("raft", log.Int32("Term", s.raft.term), log.String("cron", "start ticker"))
	s.ticker = time.NewTicker(timeout * time.Millisecond / 10)
	go func() {
	Loop:
		for {
			select {
			case <-s.ticker.C:
				if s.raft.role.role() == roleLeader { // 如果相等，则说明自身即为 Leader 节点
					s.raft.role.work()
				}
			case <-s.tickerEnd:
				break Loop
			}
		}
	}()
}

// refreshLastHeartBeatTime 更新最后一次接收到心跳时间戳ms
func (s *scheduled) refreshLastHeartBeatTime() {
	s.time = time.Now().UnixNano() / 1e6
}
