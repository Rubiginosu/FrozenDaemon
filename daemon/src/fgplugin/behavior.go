package fgplugin

import (
	"colorlog"
	"errors"
	"plugin"
	"dmserver"
)

const (
	ErrNoEnabledFuncFind  = " No enabled func find "
	ErrNoDisabledFuncFind = " No disabled func find"
)

var disabled = make([]func(), 0)

func (b *Behavior) handle(p plugin.Plugin) {
	enabled, err := p.Lookup(b.OnEnabled)
	if err != nil {
		colorlog.ErrorPrint(errors.New(ErrGlobal + ErrNoEnabledFuncFind))
	}
	if f, ok := enabled.(func()); ok {
		f()
	}
	for _,v := range b.RequestHandler{
		function,err := p.Lookup(v.functionName)
		if err != nil {
			colorlog.LogPrint("Load behavior: " + v.requestName + " Error!")
			continue
		}
		if f,ok := function.(func(dmserver.Request) dmserver.Response);ok {
			dmserver.HandleRequestIntoFunc(f,v.requestName)
		}

	}
	disabledSymbol, err := p.Lookup(b.OnDisabled)
	if err != nil {
		colorlog.ErrorPrint(errors.New(ErrGlobal + ErrNoDisabledFuncFind))
	}
	if f, ok := disabledSymbol.(func()); ok {
		disabled = append(disabled, f)
	}
}
