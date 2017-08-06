package dmserver

import "net/http"

func fuckPdcPanelHttp(){
	http.HandleFunc("/fuck",handle)
	http.ListenAndServe(":2333",nil)
}
func handle(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	switch r.Form["Method"][0]{
	case "start":
		serverSaved[0].Start()
	case "stop":
		servers[searchRunningServerByID(serverSaved[0].ID)].inputLine("stop")
	}
}