package dmserver

import "colorlog"

func HandleRequestIntoFunc(requestFunc func([]byte)[]byte,name string){
	if requestHandlers[name] == nil || *requestHandlers[name] == nil {
		newHandler := make([]func([]byte)[]byte,0)
		requestHandlers[name]=&newHandler
	}
	*requestHandlers[name] = append(*requestHandlers[name], requestFunc)
	colorlog.LogPrint("function " + colorlog.ColorSprint(name,colorlog.FR_GREEN) + " has ben registered")
}

