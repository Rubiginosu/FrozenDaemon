package fgplugin

type Behaviors struct {
	OnEnabled      string
	OnDisabled     string
	RequestHandler []RequestHandler
}


type RequestHandler struct {
	RequestName  string
	FunctionName string
}
