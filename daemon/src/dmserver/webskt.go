package dmserver

import (
	"colorlog"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
)

var OutputMaps = make(map[int]*websocket.Conn, 0)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Webskt() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+strconv.Itoa(config.ServerManager.WebSocketPort), nil)

}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(colorlog.ColorSprint("[Websocket]", colorlog.BK_CYAN), "New Websocket client connected"+r.Host)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		colorlog.ErrorPrint(err)
		return
	}
	defer c.Close()
	// TODO 鉴权

	OutputMaps[0] = c
	if len(servers) >= 1 {
		for i := 0; i < len(servers[0].BufLog); i++ {
			c.WriteMessage(websocket.TextMessage, servers[0].BufLog[i])
		}
		c.WriteMessage(websocket.TextMessage, []byte("<span class=\"am-text-success\">["+time.Now().Format("15:04:05")+"] 以上为历史信息</span>\n"))
	}
	for {
		// 心跳包
		c.WriteMessage(websocket.TextMessage, []byte("HeartPkg"))
		time.Sleep(10 * time.Second)

	}

}

func IsOutput(n int) (bool, *websocket.Conn) {
	if _, ok := OutputMaps[n]; ok {
		return true, OutputMaps[n]
	}
	return false, nil
}
