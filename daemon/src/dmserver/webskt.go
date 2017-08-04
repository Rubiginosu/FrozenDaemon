package dmserver

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"auth"
	"strconv"
)

var upgrader = websocket.Upgrader{}
func webskt(){
	http.HandleFunc("/",handler)
	http.ListenAndServe(":" + strconv.Itoa(config.ServerManager.WebSocketPort),nil)


}
func handler(w http.ResponseWriter,r *http.Request){
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	var req InterfaceRequest
	c.ReadJSON(&req)
	if auth.IsVerifiedValidationKeyPair(req.Req.OperateID,req.Auth){
		index := searchRunningServerByID(req.Req.OperateID)
		if index < 0 {
			c.WriteJSON(Response{-1,"Invalid server id"})
			return
		}
		servers[index].ToOutput.IsOutput = true
		for {
			err := c.WriteMessage(websocket.TextMessage,<-servers[index].ToOutput.To)
			if err != nil {
				servers[index].ToOutput.IsOutput = false
				break
			}
		}
	}

}