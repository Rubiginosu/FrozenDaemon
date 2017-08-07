package ftrans

import "net/http"

func Start() {
	http.HandleFunc("/download", handleDownload)
}
