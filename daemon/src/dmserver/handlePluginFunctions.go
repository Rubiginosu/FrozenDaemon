package dmserver
func HandleRequestIntoFunc(requestFunc func(Request) Response,name string){
	if *requestHandlers[name] == nil {
		*requestHandlers[name]=make([]func(Request)Response,0)
	}
	*requestHandlers[name] = append(*requestHandlers[name], requestFunc)
}

