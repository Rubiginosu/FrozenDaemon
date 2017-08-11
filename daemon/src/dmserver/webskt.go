package dmserver

import (
	"colorlog"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"auth"
	"strconv"
)

var OutputMaps = make(map[int]*websocket.Conn, 0)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Webskt() {
	http.HandleFunc("/ws", handler)
	http.ListenAndServe(":52022", nil)

}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(colorlog.ColorSprint("[Websocket]",
		colorlog.FR_CYAN), "New Websocket OUTPUT client connected" + r.RemoteAddr)
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		colorlog.ErrorPrint(err)
		return
	}
	defer c.Close()
	// 鉴权格式
	// ->   key
	// <-   result
	// result == false -> Close
	// result == true  ->  add OutputMaps -> send messages
	_, message, err2 := c.ReadMessage()
	if err2 != nil {
		colorlog.ErrorPrint(err2)
	}
	sid := auth.VerifyKey(string(message)) // 强转String再传入Auth
	if sid < 0 {
		c.WriteMessage(websocket.TextMessage,[]byte("Key Verified failed..."))
		return
	}
	// 上面检测sid。
	defer delete(OutputMaps, sid)
	c.WriteMessage(websocket.TextMessage, []byte("Verified key,sid:"+strconv.Itoa(sid)))
	if len(servers) >= 1 {
		for i := 0; i < len(servers[sid].BufLog); i++ {
			c.WriteMessage(websocket.TextMessage, servers[0].BufLog[i])
		}
		c.WriteMessage(websocket.TextMessage, []byte("<span class=\"am-text-success\">["+time.Now().Format("15:04:05")+"] 以上为历史信息</span>\n"))
	}
	OutputMaps[sid] = c
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
