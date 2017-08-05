package dmserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"colorlog"
	"github.com/gorilla/websocket"
)

// 服务器状态码
// 已经关闭
const(
	SERVER_STATUS_CLOSED = iota
	SERVER_STATUS_RUNNING
	SERVER_STATUS_STATING
)


func (s *ServerRun) Close() {

}

func (server *ServerLocal) Start() error {

	//server.EnvPrepare()
	//execConf, err0 := server.loadExecutableConfig()
	//if err0 != nil {
	//	return err0
	//}

	//cmd := exec.Command("./server", "-uid="+strconv.Itoa(config.DaemonServer.UserId), "-mem="+strconv.Itoa(server.MaxMem), "-chr="+"../servers/server"+strconv.Itoa(server.ID), "-cmd="+execConf.Command)

	//#########Testing###########
	cmd := exec.Command("/usr/bin/java","-jar","/root/test/server/mc.jar")
	cmd.Dir = "/root/test/server/"
	stdoutPipe, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	stdinPipe,err2 := cmd.StdinPipe()
	if err2 != nil {
		return err2
	}
	servers = append(servers,ServerRun{
		server.ID,
		nil,
		cmd,
		&stdinPipe,
		&stdoutPipe,
	})
	err3 := cmd.Start()
	server.Status = SERVER_STATUS_STATING
	if err3 != nil {
		return err3
	}
	go servers[len(servers)-1].ProcessOutput()
	return nil
}

func (s *ServerRun)ProcessOutput() {
	fmt.Println(s.Cmd.Process.Pid)
	buf := bufio.NewReader(*s.StdoutPipe)
	index := searchServerByID(s.ID)
	go s.getServerStopped()
	for {
		if serverSaved[index].Status == 0 {
			break
		}
		line, err := buf.ReadBytes('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		//fmt.Printf("%s",line)
		s.processOutputLine(string(line)) // string对与正则更加友好吧
		//s.ToOutput.IsOutput = true
		if isOut,to := IsOutput(s.ID);isOut{
			//colorlog.LogPrint("Trying to send line to channel.")
			to.WriteMessage(websocket.TextMessage,line)
		}
	}
	colorlog.LogPrint("Break for loop,server stopped or EOF. ")

}

func outputListOfServers() Response {
	b, _ := json.Marshal(serverSaved)
	return Response{0, string(b)}
}

// 删除服务器
func (server *ServerLocal) Delete() {

	if server.Status == SERVER_STATUS_RUNNING {
		servers[server.ID].Close()
	}
	// 如果服务器仍然开启则先关闭服务器。
	// 成功关闭后，请Golang拆迁队清理违章建筑
	nowPath, _ := filepath.Abs(".")
	serverRunPath := filepath.Clean(nowPath + "/../servers/server" + strconv.Itoa(server.ID))
	os.RemoveAll(serverRunPath)
	// 清理服务器所占的储存空间
	// 违章搭建搞定以后，把这个记账本的东东也删掉
	id := searchServerByID(server.ID)
	serverSaved = append(serverSaved[:id], serverSaved[id+1:]...)
	// go这个切片是[,)左闭右开的区间，应该这么写吧~
	// 保存服务器信息。
	saveServerInfo()
}

// 搜索服务器的ID..返回index索引
// 返回-1代表没找到
func searchServerByID(id int) int {
	for i := 0; i < len(serverSaved); i++ {
		if serverSaved[i].ID == id {
			return i
		}
	}
	return -1
}
func GetServerSaved() []ServerLocal {
	return serverSaved
}

// 搜索服务器的ID..返回index索引
// 返回-1代表没找到
func searchRunningServerByID(id int) int {
	for i := 0; i < len(servers); i++ {
		if servers[i].ID == id {
			return i
		}
	}
	return -1
}

func (s *ServerRun)getServerStopped(){
	s.Cmd.Wait()
	colorlog.PointPrint("Server Stopped")
	serverSaved[searchServerByID(s.ID)].Status = 0
}