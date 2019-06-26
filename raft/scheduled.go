package raft

import (
	pb "github.com/ennoo/fabric-client/grpc/proto/raft"
	"github.com/ennoo/rivet/utils/log"
	"github.com/panjf2000/ants"
	"github.com/robfig/cron"
	"strings"
	"time"
)

var (
	ticker    *time.Ticker
	checkCron *cron.Cron
	checkErr  chan error
	pools     map[string]*ants.PoolWithFunc
	Time      int64 // 最后一次心跳时间戳ms
	hb        string
	rv        string
	sn        string
)

func init() {
	hb = "hb"
	rv = "rv"
	sn = "sn"
	checkCron = cron.New()
	checkErr = make(chan error)
	pools = map[string]*ants.PoolWithFunc{}
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
	pools[hb], _ = ants.NewPoolWithFunc(len(Nodes), func(i interface{}) {
		heartBeat(i)
	})
	pools[rv], _ = ants.NewPoolWithFunc(len(Nodes), func(i interface{}) {
		RequestVote(i)
	})
	pools[sn], _ = ants.NewPoolWithFunc(len(Nodes), func(i interface{}) {
		syncNode(i)
	})
}

// Start 启动定时检查raft计时任务可用性任务
func Start() {
	log.Self.Info("scheduled", log.Reflect("start", Nodes))
	for _, node := range Nodes {
		if node.Id == ID {
			continue
		}
		leaderMe(strings.Join([]string{node.Addr, ":", node.Rpc}, ""), Nodes[ID])
	}
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
	times := 0
	checkCron.Stop()
	err := checkCron.AddFunc(strings.Join([]string{"*/1 * * * * ?"}, ""), func() {
		if nil == ticker {
			log.Self.Debug("scheduled", log.Int32("Term", Term), log.String("cron", "start ticker"))
			Time = time.Now().UnixNano() / 1e6
			tickerStart()
		}
		if Nodes[ID].Status == pb.Status_LEADER && Leader.BrokerID == ID { // 如果相等，则说明自身即为 Leader 节点
			if times >= 3 {
				log.Self.Debug("scheduled", log.Int32("Term", Term), log.String("sync", "发起同步节点信息"))
				// 遍历发起同步节点请求
				for _, node := range Nodes {
					if node.Id == ID {
						continue
					}
					if err := pools[sn].Invoke(&SN{
						URL: strings.Join([]string{node.Addr, ":", node.Rpc}, ""),
						Req: &pb.NodeMap{Nodes: Nodes},
					}); nil != err {
						reStartCheck()
					}
				}
				times = 0
			} else {
				times += 1
			}
		}
	})
	if nil != err {
		checkErr <- err
	} else {
		checkCron.Start()
	}
}

func tickerStart() {
	ticker = time.NewTicker(timeOut * time.Millisecond / 15)
	go func() {
		for {
			select {
			case <-ticker.C:
				if Nodes[ID].Status == pb.Status_LEADER && Leader.BrokerID == ID { // 如果相等，则说明自身即为 Leader 节点
					pools[hb].Tune(len(Nodes))
					// 遍历发送心跳
					for _, node := range Nodes {
						if node.Id == ID {
							continue
						}
						if err := pools[hb].Invoke(&HB{
							URL: strings.Join([]string{node.Addr, ":", node.Rpc}, ""),
							Req: &pb.ReqElection{Node: Nodes[ID], Term: Term},
						}); nil != err {
							reStartCheck()
						}
					}
				} else if time.Now().UnixNano()/1e6-Time > timeOut { // 如果超时没有收到 Leader 信息
					if Nodes[ID].Status != pb.Status_CANDIDATE {
						// 切换自身为 CANDIDATE 状态
						Nodes[ID].Status = pb.Status_CANDIDATE
						//Leader = &pb.Leader{}
						Term += 1
					}
					// 遍历发起索要投票请求
					for _, node := range Nodes {
						if node.Id == ID {
							continue
						}
						if err := pools[rv].Invoke(&RV{
							URL: strings.Join([]string{node.Addr, ":", node.Rpc}, ""),
							Req: &pb.ReqElection{
								Node: Nodes[ID],
								Term: Term,
							},
							Target: &pb.ReqElection{
								Node: node,
								Term: Term,
							},
						}); nil != err {
							reStartCheck()
						}
					}
				}
			}
		}
	}()
}

func RefreshTimeOut() {
	Time = time.Now().UnixNano() / 1e6
}
