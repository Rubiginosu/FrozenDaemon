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
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if authInfo, ok := r.Form["auth"]; ok {
		sid := auth.VerifyKey(authInfo[0])
		if sid == auth.KEY_VERIFY_FAILED {
			w.Write([]byte("Key Verify Failed!"))
			return
		} else if sid == auth.KEY_OUT_DATE {
			w.Write([]byte("Key out of date"))
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		colorlog.ErrorPrint(err)
		return
	}
	/* TODO 删掉测试档,
	   TODO 加入鉴权,
	   TODO 应该会专门写个函数来执行文件上传，设计权限，文件名称
	   TODO 记得加过滤！！
	 */

	file,_ := os.Create("/home/axoford12/testupload")
	defer conn.Close()
	for {
		//conn.ReadMessage()
		_,message,err := conn.ReadMessage()
		if err != nil {
			colorlog.ErrorPrint(err)
			return
		}
		file.Write(message)

		conn.WriteMessage(websocket.TextMessage,[]byte("OK"))
		conn.WriteMessage(websocket.TextMessage,[]byte("OK"))
	}
	file.Close()
}
