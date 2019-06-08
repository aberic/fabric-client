package rivet

import (
	"github.com/ennoo/rivet/keeps"
	"github.com/ennoo/rivet/trans/response"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	upgrade = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// WS websocket 启动
func WS(router *response.Router) {
	response.SyncPoolGetResponse().DoSelf(router.Context, func(writer http.ResponseWriter, request *http.Request) {
		id := router.Context.Param("id")
		var (
			conn *websocket.Conn
			err  error
		)
		if conn, err = upgrade.Upgrade(writer, request, nil); nil != err {
			return
		}
		Keepers = append(Keepers, keeps.Start(id, conn))
	})
}

// remove Service 服务器对象集合内移除
//func remove(position int) {
//	Keepers = append(Keepers[:position], Keepers[position+1:]...)
//}
