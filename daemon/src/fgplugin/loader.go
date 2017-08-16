package fgplugin

import (
	"colorlog"
	"errors"
	"os"
	"plugin"
)

const (
	ErrPlugOpenErr                       = "open plugin file "
	ErrPluginBehaviorNotFound            = "plugin behavior not defined at "
	ErrPluginBehaviorNotDefinedCorrectly = "plugin behavior not defined correctly "
)

func loader(info os.FileInfo, path string) bool {
	pluginPath := path + "/" + info.Name()
	plg, err := plugin.Open(pluginPath)
	if err != nil {
		colorlog.ErrorPrint(errors.New(ErrGlobal + ErrPlugOpenErr + pluginPath))
		return false
	}
	behavior, err := plg.Lookup("Behavior")
	if err != nil {
		colorlog.ErrorPrint(errors.New(ErrGlobal + ErrPluginBehaviorNotFound + pluginPath))
		return false
	}
	if f, ok := behavior.(func() Behavior); ok {
		if !ok {
			f()
			return true
		}
		colorlog.ErrorPrint(errors.New(ErrGlobal + ErrPluginBehaviorNotDefinedCorrectly + pluginPath))
		return false
	}
	return false
}
