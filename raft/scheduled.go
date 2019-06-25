package raft

import (
	"github.com/ennoo/rivet/utils/log"
	"github.com/robfig/cron"
	"strings"
	"time"
)

const (
	timeOut = 350 // raft心跳超时ms
)

var (
	ticker    *time.Ticker
	checkCron *cron.Cron
	checkErr  chan error
)

func init() {
	checkCron = cron.New()
	checkErr = make(chan error)
}

// ReStartCheck 重启定时检查raft计时任务可用性任务
func ReStartCheck() {
	time.Sleep(time.Second * 1)
	go check()
}

// Start 启动定时检查raft计时任务可用性任务
func Start() {
	go check()
}

// check 定时调用此方法检查raft计时任务可用性
func check() {
	go task()
	for {
		select {
		case err := <-checkErr:
			log.Self.Panic("check", log.Error(err))
			ReStartCheck()
			return
		}
	}
}

func task() {
	checkCron.Stop()
	err := checkCron.AddFunc(strings.Join([]string{"*/1 * * * * ?"}, ""), func() {
		if nil == ticker {
			log.Self.Debug("scheduled", log.String("cron", "restart ticker"))
			tickerStart()
		}
	})
	if nil != err {
		checkErr <- err
	} else {
		checkCron.Start()
	}
}

func tickerStart() {
	ticker = time.NewTicker(timeOut * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Self.Debug("scheduled", log.Reflect("ticker", time.Now().UnixNano()/1e6))
			}
		}
	}()
}
