package fgplugin

import (
	"colorlog"
	"errors"
	"os"
	"plugin"
	"encoding/json"
)

const (
	ErrPlugOpenErr                       = "open plugin file "
	ErrPluginBehaviorNotFound            = "plugin behavior not defined at "
	ErrPluginBehaviorNotDefinedCorrectly = "plugin behavior not defined correctly "
	ErrPluginBehaviorJsonUnmarshal = "plugin behavior returns not correct json format bytes :"
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
	if f, ok := behavior.(func() []byte); ok {
		pluginBehave := Behaviors{}
		err := json.Unmarshal(f(),&pluginBehave)
		if err != nil {
			colorlog.ErrorPrint(errors.New(ErrGlobal + ErrPluginBehaviorJsonUnmarshal + err.Error()))
		}
		pluginBehave.handle(*plg)
		return true

	}
	colorlog.ErrorPrint(errors.New(ErrGlobal + ErrPluginBehaviorNotDefinedCorrectly + pluginPath))
	return false
}
