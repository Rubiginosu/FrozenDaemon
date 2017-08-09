package ftrans

import (
	"auth"
	"colorlog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"github.com/gorilla/websocket"
	"os"
	"fmt"
	"regexp"
	"errors"
	"os/exec"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if authInfo, ok := r.Form["auth"]; ok {
		sid := auth.VerifyKey(authInfo[0])
		if sid < 0 {
			w.Write([]byte("Key Verify Failed!"))
			return
		}
		// 验证成功！
		// 此时的sid 储存server id
		path := "../servers/server" + strconv.Itoa(sid) + "/"
		if reqFile, ok := r.Form["req"]; ok {
			current := filepath.Clean(path + reqFile[0])
			colorlog.LogPrint("Request file:" + current)
			if strings.Index(current, path) < 0 {
				w.Write([]byte("Permission denied."))
				return
			}
			http.ServeFile(w, r, current)
		}
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println(colorlog.ColorSprint("[Websocket]",
		colorlog.BK_CYAN), "New Websocket UPLOAD client connected" + r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		colorlog.ErrorPrint(err)
		return
	}
	_,message,err2 := conn.ReadMessage()
	if err2 != nil {
		colorlog.ErrorPrint(err2)
		conn.Close()
		return
	}

	sid := auth.VerifyKey(string(message))
	if sid < 0 {
		colorlog.LogPrint("Websocket client auth failed")
		conn.WriteMessage(websocket.TextMessage,[]byte("Key Verified failed"))
	}
	colorlog.LogPrint("Websocket client auth ok.sid:" + strconv.Itoa(sid))
	conn.WriteMessage(websocket.TextMessage,[]byte("Verified key,sid:" + strconv.Itoa(sid)))
	receiveWriteUploadFile(conn,sid)


}

/*
解析要求
 =>  test.txt | 777
 <=
 */
func receiveWriteUploadFile(conn *websocket.Conn,sid int){
	_,message,err := conn.ReadMessage()
	if err != nil {
		colorlog.ErrorPrint(err)
		conn.Close()
		return
	}
	name,mode,err2 := parseUploadMessage(message)
	if err2 != nil {
		colorlog.ErrorPrint(err2)
		conn.WriteMessage(websocket.TextMessage,[]byte(err2.Error()))
		conn.Close()
		return
	} else {
		colorlog.LogPrint("Client request to upload file: " + name + " and mode:" + mode)
		conn.WriteMessage(websocket.TextMessage,[]byte("OK"))
	}
	current := filepath.Clean("../servers/server" + strconv.Itoa(sid) + "/" + name) // 过滤检查
	if strings.Index(current,"../servers/server" + strconv.Itoa(sid)) <= 0 {
		conn.WriteMessage(websocket.TextMessage,[]byte("Permission denied."))
	}

	file,err3 := os.OpenFile(current,os.O_TRUNC | os.O_WRONLY | os.O_CREATE | os.O_SYNC,777)
	if err3 != nil {
		colorlog.ErrorPrint(err3)
		conn.WriteMessage(websocket.TextMessage,[]byte(err3.Error()))
		conn.Close()
		return
	}
	for {
		//conn.ReadMessage()
		_,message,err := conn.ReadMessage()
		if err != nil {
			colorlog.ErrorPrint(err)
			return
		}
		_,err2 := file.Write(message)
		if err != nil {
			colorlog.ErrorPrint(err2)
			conn.WriteMessage(websocket.TextMessage,[]byte(err2.Error()))
			conn.Close()
			return
		}
		conn.WriteMessage(websocket.TextMessage,[]byte("OK"))
		conn.WriteMessage(websocket.TextMessage,[]byte("OK"))
	}
	file.Close()
	cmd := exec.Command("/bin/chmod",string(mode),current)
	err4 := cmd.Run()
	if err4 != nil{
		colorlog.ErrorPrint(err3)
		conn.WriteMessage(websocket.TextMessage,[]byte(err3.Error()))
		conn.Close()
	}
}
func parseUploadMessage(message []byte) (string,string,error){
	res := regexp.MustCompile(" *| *").Split(string(message),2)
	if len(res) < 2 {
		return "","",errors.New("Invalid file info format.")
	}
	return res[0],res[1],nil
}