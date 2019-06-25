package raft

import (
	pb "github.com/ennoo/fabric-client/grpc/proto/raft"
	"github.com/ennoo/rivet/utils/log"
	"github.com/panjf2000/ants"
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
	pools     map[string]*ants.PoolWithFunc
	Time      int64 // 最后一次心跳时间戳ms
)

func init() {
	checkCron = cron.New()
	checkErr = make(chan error)
	Time = time.Now().UnixNano() / 1e6
}

// ReStartCheck 重启定时检查raft计时任务可用性任务
func reStartCheck() {
	for _, pool := range pools {
		_ = pool.Release()
	}
	initPools()
	go check()
}

func initPools() {
	pools["hb"], _ = ants.NewPoolWithFunc(len(Nodes), func(i interface{}) {
		hb(i)
	})
	pools["rv"], _ = ants.NewPoolWithFunc(len(Nodes), func(i interface{}) {
		rv(i)
	})
	pools["fm"], _ = ants.NewPoolWithFunc(len(Nodes), func(i interface{}) {
		fm(i)
	})
	pools["lm"], _ = ants.NewPoolWithFunc(len(Nodes), func(i interface{}) {
		lm(i)
	})
	pools["sn"], _ = ants.NewPoolWithFunc(len(Nodes), func(i interface{}) {
		sn(i)
	})
}

// Start 启动定时检查raft计时任务可用性任务
func Start() {
	reStartCheck()
}

// check 定时调用此方法检查raft计时任务可用性
func check() {
	go task()
	for {
		select {
		case err := <-checkErr:
			log.Self.Panic("scheduled", log.Error(err))
			reStartCheck()
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
	ticker = time.NewTicker(timeOut * time.Millisecond / 3)
	go func() {
		for {
			select {
			case <-ticker.C:
				if Leader.BrokerID == ID { // 如果相等，则说明自身即为 Leader 节点
					pools["hb"].Tune(len(Nodes))
					// 遍历发送心跳
					for _, node := range Nodes {
						if err := pools["hb"].Invoke(&HB{
							URL: strings.Join([]string{node.Addr, ":", node.Rpc}, ""),
							Req: &pb.Beat{Beat: []byte(ID)},
						}); nil != err {
							reStartCheck()
						}
					}
				} else if time.Now().UnixNano()/1e6-Time > 350 { // 如果超时没有收到 Leader 信息
					// 切换自身为 CANDIDATE 状态
					// 发起索要投票请求
				}
				log.Self.Debug("scheduled", log.Reflect("ticker", time.Now().UnixNano()/1e6))
			}
		}
	}()
}
