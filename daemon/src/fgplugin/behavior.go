package fgplugin

import (
	"colorlog"
	"errors"
	"plugin"
	"dmserver"
)

const (
	ErrNoEnabledFuncFind  = "checking func: No enabled func find "
	ErrNoDisabledFuncFind = "checking func: No disabled func find"
	ErrLoadBehavior = "Loading behavior"
)

var disabled = make([]func(), 0)

func (b *Behaviors) handle(p plugin.Plugin) {
	enabled, err := p.Lookup(b.OnEnabled)
	colorlog.LogPrint(b.OnEnabled)
	if err != nil {
		colorlog.ErrorPrint(ErrNoEnabledFuncFind,errors.New(""))
	}
	if f, ok := enabled.(func()); ok {
		f()
	}
	for _,v := range b.RequestHandler{
		funcSymbol,err := p.Lookup(v.FunctionName)
		if err != nil {
			colorlog.ErrorPrint(ErrLoadBehavior + v.RequestName,errors.New(""))
			continue
		}
		if f,ok := funcSymbol.(func([]byte)[]byte);ok {
			dmserver.HandleRequestIntoFunc(f,v.RequestName)
		}

	}
	disabledSymbol, err := p.Lookup(b.OnDisabled)
	if err != nil {
		colorlog.ErrorPrint(ErrNoDisabledFuncFind,errors.New(""))
	}
	if f, ok := disabledSymbol.(func()); ok {
		disabled = append(disabled, f)
	}
}
