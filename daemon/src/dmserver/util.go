package dmserver

import (
	"colorlog"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

/**
构建一个Infos
*/
func buildServerInfos(infos []os.FileInfo) []ServerPathFileInfo {
	spfInfos := make([]ServerPathFileInfo, len(infos))
	for k, v := range infos {
		spfInfos[k] = __buildServerInfo(v)
	}
	return spfInfos
}

// 上面函数的辅助方法
func __buildServerInfo(info os.FileInfo) ServerPathFileInfo {
	return ServerPathFileInfo{
		Name:    info.Name(),
		Mode:    strings.Replace(fmt.Sprintf("%32b", uint32(info.Mode())), " ", "0", -1),
		ModTime: info.ModTime().Unix(),
	}
}

/**
检查目录是否合法，避免出现../../../../etc/passwd之类似小天才行为
*/
func validateOperateDir(upload string, path string) bool {
	return strings.Index(filepath.Clean(upload+path), upload) >= 0
}

func (server *ServerLocal) initCgroup() {
	// 当检测到cg目录不存在时，启动cg.init
	cmd := exec.Command("/bin/bash",
		"../cgroup/cg.sh",
		"cg",
		"init",
		"server"+strconv.Itoa(server.ID),
		strconv.Itoa(server.MaxCpuUtilizatioRate),
		strconv.Itoa(server.MaxMem),
		strconv.Itoa(server.MaxHardDiskReadSpeed),
		strconv.Itoa(server.MaxHardDiskWriteSpeed),
		config.DaemonServer.BlockDeviceMajMim,
		strings.Replace(fmt.Sprintf("%4x", server.ID), " ", "0", -1))
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		colorlog.ErrorPrint(errors.New("Error with init cgroups:" + err.Error()))
		colorlog.LogPrint("Reaseon:" + string(output))
		colorlog.PromptPrint("This server's source may not valid")
	}
}

func (server *ServerLocal) networkDel() {
	cmd := exec.Command("/bin/bash",
		"../cgroup/cg.sh",
		"net",
		"del",
		strings.Replace(fmt.Sprintf("%4x", server.ID), " ", "0", -1),
		"1",
		"1",
		config.DaemonServer.NetworkCardName)
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		colorlog.ErrorPrint(errors.New("Error with del network:" + err.Error()))
		colorlog.LogPrint("Reaseon:" + string(output))
	}
}

func (server *ServerLocal) networkFlush() {
	cmd := exec.Command("/bin/bash",
		"../cgroup/cg.sh",
		"net",
		"change",
		strings.Replace(fmt.Sprintf("%4x", server.ID), " ", "0", -1),
		strconv.Itoa(server.MaxUsingUpBandwidth),
		strconv.Itoa(server.MaxUnusedUpBandwidth),
		config.DaemonServer.NetworkCardName)
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		colorlog.ErrorPrint(errors.New("Error with flush network:" + err.Error()))
		colorlog.LogPrint("Reaseon:" + string(output))
	}
}
func (server *ServerLocal) performanceFlush() {
	server.cgroupDel()
	server.initCgroup()
}
func (server *ServerLocal) cgroupDel() {
	cmd := exec.Command("/bin/bash",
		"../cgroup/cg.sh",
		"cg",
		"del",
		"server"+strconv.Itoa(server.ID))
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		colorlog.ErrorPrint(errors.New("Error with del cgroup:" + err.Error()))
		colorlog.LogPrint("Reaseon:" + string(output))
	}
}
