package dmserver

import (
	"net/http"
	"github.com/gorilla/websocket"
	"strconv"
	"colorlog"
	"time"
)
var OutputMaps = make(map[int]*websocket.Conn,0)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Webskt() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+strconv.Itoa(config.ServerManager.WebSocketPort), nil)

}
func handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		colorlog.ErrorPrint(err)
		return
	}
	defer c.Close()
	// TODO 鉴权
	for i:=0;i<len(servers[0].BufLog);i++{
		c.WriteMessage(websocket.TextMessage,servers[0].BufLog[i])
	}
	OutputMaps[0] = c
	for {
		// 心跳包
		c.WriteMessage(websocket.TextMessage,[]byte("HeartPkg"))
		time.Sleep(10 * time.Second)

	}


}

func IsOutput(n int) (bool,*websocket.Conn){
	if _,ok := OutputMaps[n];ok {
		return true,OutputMaps[n]
	}
	return false,nil
}