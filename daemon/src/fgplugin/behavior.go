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

func (b *Behaviors) handle(p plugin.Plugin) {
	enabled, err := p.Lookup(b.OnEnabled)
	colorlog.LogPrint(b.OnEnabled)
	if err != nil {
		colorlog.ErrorPrint(errors.New(ErrGlobal + ErrNoEnabledFuncFind))
	}
	if f, ok := enabled.(func()); ok {
		f()
	}
	for _,v := range b.RequestHandler{
		funcSymbol,err := p.Lookup(v.FunctionName)
		if err != nil {
			colorlog.ErrorPrint(errors.New("Load behavior: " + v.RequestName + " Error!"))
			continue
		}
		if f,ok := funcSymbol.(func([]byte)[]byte);ok {
			dmserver.HandleRequestIntoFunc(f,v.RequestName)
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
