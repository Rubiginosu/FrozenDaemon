package dmserver

import (
	"bufio"
	"colorlog"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"
	"time"
	"conf"
	"utils"
)

// 服务器状态码
// 已经关闭
const (
	SERVER_STATUS_CLOSED = iota
	SERVER_STATUS_RUNNING
	SERVER_STATUS_STATING
)

func (s *ServerRun) Close() {
	colorlog.LogPrint("Closing server...")
	var execConf ExecConf
	if server, ok := serverSaved[s.ID]; ok {
		var err error
		execConf, err = server.loadExecutableConfig()
		if err != nil {
			colorlog.ErrorPrint("loading exec config",err)
			return
		}
	}
	colorlog.LogPrint("Closing command is" + execConf.StoppedServerCommand)
	go time.AfterFunc(20*time.Second, func() {
		// 杀死进程组.
		colorlog.PointPrint("Timeout,kill them.")
		if serverSaved[s.ID].Status != 0 {
			if s.Cmd.Process != nil {
				syscall.Kill(s.Cmd.Process.Pid, syscall.SIGKILL)
			}
		}
	})
	s.inputLine(execConf.StoppedServerCommand)
}

func (server *ServerLocal) Start() error {

	server.EnvPrepare()
	colorlog.LogPrint("Done.")
	execConf, err0 := server.loadExecutableConfig()
	if err0 != nil {
		return err0
	}
	os.Mkdir("../servers/server" + strconv.Itoa(server.ID) + "/serverData",775)
	os.Chown("../servers/server" + strconv.Itoa(server.ID) + "/serverData",config.DaemonServer.UserId,0)
	command := regexp.MustCompile(" +").Split(execConf.Command,-1)
	commandArgs := []string{
		strconv.Itoa(config.DaemonServer.UserId),
		"../servers/server" + strconv.Itoa(server.ID),
		"/serverData",
	}
	commandArgs = append(commandArgs,command...)
	cmd := exec.Command("./server", commandArgs...)
	//#########Testing###########
	stdoutPipe, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	stdinPipe, err2 := cmd.StdinPipe()
	if err2 != nil {
		return err2
	}
	servers[server.ID] = &ServerRun{
		server.ID,
		[]string{},
		cmd,
		make([][]byte, 50),
		&stdinPipe,
		&stdoutPipe,
	}
	err3 := cmd.Start()
	server.Status = SERVER_STATUS_STATING
	if err3 != nil {
		return err3
	}
	start, join, left := execConf.getRegexps()
	cmdCgroup := exec.Command("/bin/bash",
		"../cgroup/cg.sh",
		"cg",
		"run",
		"server"+strconv.Itoa(server.ID),
		strconv.Itoa(cmd.Process.Pid))
	cmdCgroup.Env = os.Environ()
	output, err4 := cmdCgroup.CombinedOutput()
	if err4 != nil {
		colorlog.ErrorPrint("initialing cgroups", err4)
		colorlog.LogPrint("Reaseon:" + string(output))
		colorlog.PromptPrint("This server's source may not valid")
	}
	go servers[server.ID].ProcessOutput(start, join, left) // 将三个参数传递
	return nil
}

func (s *ServerRun) ProcessOutput(start, join, left *regexp.Regexp) {
	if s.Cmd == nil || s.Cmd.Process == nil {
		return
	}
	colorlog.LogPrint(fmt.Sprintf("PID: %d", s.Cmd.Process.Pid))
	buf := bufio.NewReader(*s.StdoutPipe)
	if _, ok := serverSaved[s.ID]; !ok {
		delete(servers, s.ID)
		return
	}
	defer colorlog.LogPrint("Break for loop,server stopped or EOF. ")
	go s.getServerStopped()
	for {
		if serverSaved[s.ID].Status == 0 {
			break
		}
		line, err := buf.ReadBytes('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		fmt.Printf("%s", line)
		s.BufLog = append(s.BufLog[1:], line)
		s.processOutputLine(string(line), start, join, left) // string对与正则更加友好吧
		//s.ToOutput.IsOutput = true
		if isOut, to := IsOutput(s.ID); isOut {
			// 向ws客户端输出.
			to.WriteMessage(websocket.TextMessage, line)
		}
	}

	//delete(serverSaved,s.ID)

}

// 删除服务器
func (server *ServerLocal) Delete() {

	if server.Status == SERVER_STATUS_RUNNING {
		servers[server.ID].Close()
	}
	server.cgroupDel()
	server.networkDel()
	// 如果服务器仍然开启则先关闭服务器。
	// 成功关闭后，请Golang拆迁队清理违章建筑
	nowPath, _ := filepath.Abs(".")



	serverRunPath := filepath.Clean(nowPath + "/../servers/server" + strconv.Itoa(server.ID))

	if config.DaemonServer.HardDiskMethod == conf.HDM_MOUNT {
		colorlog.LogPrint("Umounting dirs.")
		cmd := exec.Command("umount","-f",serverRunPath + "/*")
		utils.AutoRunCmdAndOutputErr(cmd,"umount dirs")
	}


	os.RemoveAll(serverRunPath)
	if config.DaemonServer.HardDiskMethod != conf.HDM_LINK {
		os.Remove(serverRunPath + ".loop")
	}
	// 清理服务器所占的储存空间
	// 违章搭建搞定以后，把这个记账本的东东也删掉
	delete(serverSaved, server.ID)
	// 保存服务器信息。
	saveServerInfo()
}

func GetServerSaved() map[int]*ServerLocal {
	return serverSaved
}

func (s *ServerRun) getServerStopped() {
	s.Cmd.Wait()
	serverSaved[s.ID].Status = 0
	delete(servers, s.ID)
	colorlog.PointPrint("Server Stopped")
}

// 获取那些正则表达式
func (e *ExecConf) getRegexps() (*regexp.Regexp, *regexp.Regexp, *regexp.Regexp) {
	startReg, err := regexp.Compile(e.StartServerRegexp)
	if err != nil {
		colorlog.ErrorPrint("compiling start regexp",err)
		startReg = regexp.MustCompile("Done \\(.+s\\)!") // 用户自己的表达式骚写时,打印错误信息并使用系统默认表达式(1.7.2 spigot)
	}
	joinReg, err2 := regexp.Compile(e.NewPlayerJoinRegexp)
	if err2 != nil {
		colorlog.ErrorPrint("compiling NewPlayerJoin regexp",err2)
		joinReg = regexp.MustCompile("(\\w+)\\[.+\\] logged in")
	}
	exitReg, err3 := regexp.Compile(e.PlayerExitRegexp)
	if err3 != nil {
		colorlog.ErrorPrint("compiling player exit regexp",err3)
		exitReg = regexp.MustCompile("(\\w+) left the game")
	}
	return startReg, joinReg, exitReg
}
