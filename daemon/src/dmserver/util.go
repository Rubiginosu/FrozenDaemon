package dmserver

import (
	"colorlog"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"utils"
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
	args := []string{"../cgroup/cg.sh",
		"cg",
		"init",
		"server" + strconv.Itoa(server.ID),
		strconv.Itoa(server.MaxCpuUtilizatioRate),
		strconv.Itoa(server.MaxMem),
		strconv.Itoa(server.MaxHardDiskReadSpeed),
		strconv.Itoa(server.MaxHardDiskWriteSpeed),
		config.DaemonServer.BlockDeviceMajMim,
		strings.Replace(fmt.Sprintf("%4x", server.ID), " ", "0", -1)}
	// 当检测到cg目录不存在时，启动cg.init
	cmd := exec.Command("bash", args...)
	colorlog.LogPrint("Running command:" + dumpCommand(args))
	utils.AutoRunCmdAndOutputErr(cmd,"initial cgroup")
	colorlog.LogPrint("Init cgroup done.")
}

func (server *ServerLocal) networkDel() {

	args := []string{"../cgroup/cg.sh",
		"net",
		"del",
		strings.Replace(fmt.Sprintf("%4x", server.ID), " ", "0", -1),
		"1",
		"1",
		config.DaemonServer.NetworkCardName}
	cmd := exec.Command("bash", args...)
	colorlog.LogPrint("Running command:" + dumpCommand(args))
	utils.AutoRunCmdAndOutputErr(cmd,"delete network")
}

func (server *ServerLocal) networkFlush() {
	args := []string{
		"../cgroup/cg.sh",
		"net",
		"change",
		strings.Replace(fmt.Sprintf("%4x", server.ID), " ", "0", -1),
		strconv.Itoa(server.MaxUsingUpBandwidth),
		strconv.Itoa(server.MaxUnusedUpBandwidth),
		config.DaemonServer.NetworkCardName}
	cmd := exec.Command("bash", args...)
	colorlog.LogPrint("Running command:" + dumpCommand(args))
	utils.AutoRunCmdAndOutputErr(cmd,"flush network")
}
func (server *ServerLocal) performanceFlush() {
	server.cgroupDel()
	server.initCgroup()
}
func (server *ServerLocal) cgroupDel() {

	if _, err := os.Stat("../sys/fs/cgroup/cpu/server" + strconv.Itoa(server.ID)); err != nil {
		colorlog.LogPrint("Cgroup not exists.No need to del. ")
		return
	}
	args := []string{
		"../cgroup/cg.sh",
		"cg",
		"del",
		"server" + strconv.Itoa(server.ID),
	}
	cmd := exec.Command("/bin/bash", args...)
	colorlog.LogPrint("Running command:" + dumpCommand(args))
	utils.AutoRunCmdAndOutputErr(cmd,"delete cgroup")
	colorlog.LogPrint("Del cgroup done.")
}

func dumpCommand(args []string) string {
	s := ""
	for _, v := range args {
		s += fmt.Sprint(v + " ")
	}
	return s
}

