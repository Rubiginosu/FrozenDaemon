package fgplugin

import (
	"colorlog"
	"errors"
	"io/ioutil"
	"os"
)

const (
	ErrGlobal     = "Error occurred at package fgplugin : "
	ErrOpenPath   = "opening path "
	ErrPathNotDir = "path is not a directory"
)

func LoadPlugin(path string) {
	colorlog.LogPrint("Loading plugins from " + path)
	dir, err := os.Open(path)

	if err != nil {
		// 打开目录出现错误，报错并结束程序
		colorlog.ErrorPrint(errors.New(ErrGlobal + ErrOpenPath + err.Error()))
		dir.Close()
		return
	}
	if info, err := dir.Stat(); err != nil {
		if !info.IsDir() {
			colorlog.LogPrint("Loading " + info.Name())
			colorlog.ErrorPrint(errors.New(ErrGlobal + ErrPathNotDir))
			dir.Close()
			return
		}
	}
	infoSlice, _ := ioutil.ReadDir(path)
	for _, v := range infoSlice {
		if !v.IsDir() {
			loader(v, path)
		}
	}
}
