package ftrans

import (
	"net/http"
	"auth"
	"strconv"
	"path/filepath"
	"strings"
	"colorlog"
)

func handleDownload(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	if authInfo,ok := r.Form["auth"];ok {
		sid := auth.VerifyKey(authInfo[0])
		if sid == auth.KEY_VERIFY_FAILED {
			w.Write([]byte("Key Verify Failed!"))
			return
		} else if sid == auth.KEY_OUT_DATE{
			w.Write([]byte("Key out of date"))
			return
		}
		// 验证成功！
		// 此时的sid 储存server id
		path := "../servers/server" + strconv.Itoa(sid) + "/"
		if reqFile,ok := r.Form["req"];ok{
			current := filepath.Clean(path + reqFile[0])
			colorlog.LogPrint("Request file:" + current)
			if strings.Index(current,path) < 0 {
				w.Write([]byte("Permission denied."))
				return
			}
			http.ServeFile(w,r,current)
		}
	}
}