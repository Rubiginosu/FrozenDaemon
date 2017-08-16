package fgplugin

type Behavior struct {
	OnEnabled      string
	OnDisabled     string
	RequestHandler []RequestHandler
}

type eventHandler struct {
	event   int
	handler string
}

type RequestHandler struct {
	requestName  string
	functionName string
}
