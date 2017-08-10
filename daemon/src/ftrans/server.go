package ftrans

import (
	"net/http"
	"conf"
)

func Start(configure conf.Config) {
	config = configure
	http.HandleFunc("/download", handleDownload)
	http.HandleFunc("/upload",handleUpload)
}
