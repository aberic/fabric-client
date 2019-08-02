package keeps

import (
	"github.com/ennoo/rivet/utils/log"
	"github.com/gorilla/websocket"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"sync"
)

// Keeper websocket 对象
type Keeper struct {
	ID        string
	conn      *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	mutex     sync.Mutex
	closed    bool
}

// Start 启用 websocket
func Start(id string, conn *websocket.Conn) *Keeper {
	keeper := &Keeper{
		ID:        id,
		conn:      conn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}
	go keeper.readLoop()
	go keeper.writeLoop()
	return keeper
}

// Read 读取 websocket 数据
func (keeper *Keeper) Read() (data []byte, err error) {
	select {
	case data = <-keeper.inChan:
		log.Self.Debug("data = " + string(data))
	case <-keeper.closeChan:
		err = errors.New("id " + keeper.ID + "connection is closed")
	}
	return
}

// Write 写入数据到 websocket
func (keeper *Keeper) Write(data []byte) (err error) {
	select {
	case keeper.outChan <- data:
	case <-keeper.closeChan:
		err = errors.New("id " + keeper.ID + "connection is closed")
	}
	return
}

// Close 关闭当前 websocket
func (keeper *Keeper) Close() {
	_ = keeper.conn.Close()

	keeper.mutex.Lock()
	if !keeper.closed {
		close(keeper.closeChan)
		keeper.closed = true
	}
	keeper.mutex.Unlock()
}

func (keeper *Keeper) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = keeper.conn.ReadMessage(); nil != err {
			goto ERR
		}
		select {
		case keeper.inChan <- data:
			log.Self.Debug("data = " + string(data))
		case <-keeper.closeChan:
			goto ERR
		}
	}
ERR:
	keeper.Close()
}

func (keeper *Keeper) writeLoop() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-keeper.outChan:
			if err = keeper.conn.WriteMessage(websocket.TextMessage, data); nil != err {
				goto ERR
			}
		case <-keeper.closeChan:
			goto ERR
		}
	}
ERR:
	keeper.Close()
}
