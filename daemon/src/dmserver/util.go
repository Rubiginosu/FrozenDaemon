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
	cmd := exec.Command("/bin/bash", args...)
	colorlog.LogPrint("Running command:" + dumpCommand(args))
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		colorlog.ErrorPrint(errors.New("Error with init cgroups:" + err.Error()))
		OutputErrReason(output)
		colorlog.PromptPrint("This server's source may not valid")
	}
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
	cmd := exec.Command("/bin/bash", args...)
	colorlog.LogPrint("Running command:" + dumpCommand(args))
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		colorlog.ErrorPrint(errors.New("Error with del network:" + err.Error()))
		OutputErrReason(output)
	}
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
	cmd := exec.Command("/bin/bash", args...)
	colorlog.LogPrint("Running command:" + dumpCommand(args))
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		colorlog.ErrorPrint(errors.New("Error with flush network:" + err.Error()))
		OutputErrReason(output)
	}
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
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		colorlog.ErrorPrint(errors.New("Error with del cgroup:" + err.Error()))
		OutputErrReason(output)
	}
	colorlog.LogPrint("Del cgroup done.")
}

func dumpCommand(args []string) string {
	s := ""
	for _, v := range args {
		s += fmt.Sprint(v + " ")
	}
	return s
}
func OutputErrReason(output []byte) {
	colorlog.LogPrint("Reason:")
	fmt.Println(colorlog.ColorSprint("-----ERROR_MESSAGE-----", colorlog.FR_RED))
	fmt.Println(colorlog.ColorSprint(string(output), colorlog.FR_RED))
	fmt.Println(colorlog.ColorSprint("-----ERROR_MESSAGE-----", colorlog.FR_RED))
}

func AutoRunCmdAndOutputErr(cmd *exec.Cmd,errorAt string) bool{
	cmd.Env = os.Environ()
	out,err := cmd.CombinedOutput()
	if err != nil {
		colorlog.ErrorPrint(errors.New("Error occurred at " + errorAt + ": " + err.Error()))
		OutputErrReason(out)
		return false
	}
	return true
}