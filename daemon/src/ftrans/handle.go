package ftrans

import (
	"auth"
	"colorlog"
	"conf"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"utils"
)

var config conf.Cnf

/*
[WS统一注释原则]
1.协议说明，（该区域仅包含认证协议）
	a.  <- 代表输出至浏览器websocket
	b.  -> 代表从浏览器websocket读取数据
	c.  <= 代表打印Log信息
 	d.  ||  或者,仅包含其一
	e.  && 附加动作
	f.  CLOSE 关闭连接结束函数
*/
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

/*

+--------------------------------------------------------------------------
|   -> Key
|   <- Verified key || ( Key Verified failed && CLOSE )
|   剩下部分转交receiveWriteUploadFile(conn,sid)
+--------------------------------------------------------------------------


*/
func handleUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println(colorlog.ColorSprint("[Websocket]",
		colorlog.FR_CYAN), "New Websocket UPLOAD client connected"+r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		colorlog.ErrorPrint("upgrading func",err)
		return
	}
	_, message, err2 := conn.ReadMessage()
	if err2 != nil {
		colorlog.ErrorPrint("reading message",err2)
		conn.Close()
		return
	}

	sid := auth.VerifyKey(string(message))
	if sid < 0 {
		colorlog.LogPrint("Websocket client auth failed")
		conn.WriteMessage(websocket.TextMessage, []byte("Key Verified failed"))
		conn.Close()
		return
	}
	colorlog.LogPrint("Websocket client auth ok.sid:" + strconv.Itoa(sid))
	conn.WriteMessage(websocket.TextMessage, []byte("Verified key"))
	receiveWriteUploadFile(conn, sid)

}

/*
-> 文件名|文件权限 #  | 符号左右两边的空格会包含进去
<= 客户端请求信息
<-
	{
		Invalid file info format. # 错误的请求格式
		Permission denied.        # 你想通过我日站？想多了
		OK						  # 小学生鉴定完毕，你不是小学生
	}
-> 文件流........
-> OK # 对比以前的版本，魔法两次上传还是被我一不小心碰巧解决了，真是大好大好
-> CLOSE
## Close掉链接以后，服务端跳出写包循环，然后进行chmod 以及chown操作

*/
func receiveWriteUploadFile(conn *websocket.Conn, sid int) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		colorlog.ErrorPrint("read message",err)
		conn.Close()
		return
	}
	name, mode, err2 := parseUploadMessage(message)
	if err2 != nil {
		colorlog.ErrorPrint("parsing message",err2)
		conn.WriteMessage(websocket.TextMessage, []byte(err2.Error()))
		conn.Close()
		return
	}
	current := filepath.Clean("../servers/server" + strconv.Itoa(sid) + "/" + name) // 过滤检查
	if strings.Index(current, "../servers/server"+strconv.Itoa(sid)) < 0 {
		conn.WriteMessage(websocket.TextMessage, []byte("Permission denied."))
		conn.Close()
		return
	}
	colorlog.LogPrint("Client request to upload file: " + name + " and mode:" + mode)
	conn.WriteMessage(websocket.TextMessage, []byte("Ready to upload"))

	file, err3 := os.OpenFile(current, os.O_TRUNC|os.O_WRONLY|os.O_CREATE|os.O_SYNC, 777)
	if err3 != nil {
		colorlog.ErrorPrint("writing file",err3)
		conn.WriteMessage(websocket.TextMessage, []byte(err3.Error()))
		conn.Close()
		return
	}
	for {
		//conn.ReadMessage()
		_, message, err := conn.ReadMessage()
		if err != nil {
			if err.Error() != "websocket: close 1005 (no status)" {
				colorlog.ErrorPrint("reading message",err)
				return
			} else {
				colorlog.PointPrint("Websocket upload finished.")
				break
			}

		}
		_, err2 := file.Write(message)
		if err != nil {
			colorlog.ErrorPrint("writing message",err2)
			conn.WriteMessage(websocket.TextMessage, []byte(err2.Error()))
			conn.Close()
			return
		}
		conn.WriteMessage(websocket.TextMessage, []byte("OK"))
	}
	file.Close()
	cmd := exec.Command("chmod", string(mode), current)
	colorlog.LogPrint("Running :chmod" + " " + string(mode) + " " + current)

	if !utils.AutoRunCmdAndOutputErr(cmd,"run chmod") {
		conn.WriteMessage(websocket.TextMessage, []byte(err3.Error()))
		conn.Close()
	}
	os.Chown(current, config.DaemonServer.UserId, config.DaemonServer.UserId)
}
func parseUploadMessage(message []byte) (string, string, error) {
	res := strings.Split(string(message), "|")
	if len(res) < 2 {
		return "", "", errors.New("Invalid file info format.")
	}
	return res[0], res[1], nil
}
