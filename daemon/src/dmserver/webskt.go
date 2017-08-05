package dmserver

import (
	"net/http"
	"github.com/gorilla/websocket"
	"strconv"
	"colorlog"
	"time"
)

var upgrader = websocket.Upgrader{}
func Webskt(){
	http.HandleFunc("/",handler)
	http.ListenAndServe(":" + strconv.Itoa(config.ServerManager.WebSocketPort),nil)

}
func handler(w http.ResponseWriter,r *http.Request){
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		colorlog.ErrorPrint(err)
		return
	}
	defer c.Close()
	//var req InterfaceRequest
	//c.ReadJSON(&req)
	//if auth.IsVerifiedValidationKeyPair(req.Req.OperateID,req.Auth){
	//	index := searchRunningServerByID(req.Req.OperateID)
	//	if index < 0 {
	//		c.WriteJSON(Response{-1,"Invalid server id"})
	//		return
	//	}
		servers[0].ToOutput.IsOutput = true
		for {
			c.WriteMessage(websocket.TextMessage,[]byte("12123"))
			time.Sleep(1 * time.Second)
			//if err != nil {
			//	servers[0].ToOutput.IsOutput = false
			//	break
			//}
		}
	//}

}