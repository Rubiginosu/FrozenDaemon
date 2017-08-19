package fgplugin

import (
	"colorlog"
	"errors"
	"io/ioutil"
	"os"
)

const (
	ErrOpenPath   = "opening path "
	ErrPathNotDir = "checking plugin path :not a directory"
	ErrLoadPlugin = "loading plugins"
)

func LoadPlugin(path string) {
	colorlog.LogPrint("Loading plugins from " + path)
	dir, err := os.Open(path)

	if err != nil {
		// 打开目录出现错误，报错并结束程序
		colorlog.ErrorPrint(ErrOpenPath,err)
		dir.Close()
		return
	}
	if info, err := dir.Stat(); err != nil {
		if !info.IsDir() {
			colorlog.LogPrint("Loading " + info.Name())
			colorlog.ErrorPrint(ErrPathNotDir,err)
			dir.Close()
			return
		}
	}
	infoSlice, _ := ioutil.ReadDir(path)
	for _, v := range infoSlice {
		if !v.IsDir() {
			if !loader(v, path) {
				colorlog.ErrorPrint(ErrLoadPlugin,errors.New(""))
			}
		}
	}
}
