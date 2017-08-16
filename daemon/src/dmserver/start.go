/**
本包提供了对于Daemon服务器的通信协议并且对于请求=进行处理.
*/
package dmserver

import (
	"auth"
	"colorlog"
	"conf"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

func StartDaemonServer(conf conf.Cnf) {
	config = conf
	b, _ := ioutil.ReadFile(config.ServerManager.Servers)
	go serverOutDateClearer()
	go auth.Timer()
	err2 := json.Unmarshal(b, &serverSaved)
	if err2 != nil {
		fmt.Println(err2)
		os.Exit(-2)
	}
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(config.DaemonServer.Port)) // 默认使用tcp连接
	if err != nil {
		panic(err)
	} else {
		for {
			conn, err := ln.Accept()
			fmt.Println(colorlog.ColorSprint("[Daemon]", colorlog.FR_CYAN), "New Client Request send.From "+conn.LocalAddr().String())
			if err != nil {
				continue
			}
			go handleConnection(conn)
		}
	}

}
func StopDaemonServer() error {
	for _, v := range serverSaved {
		if v.Status != 0 {
			if server, ok := servers[v.ID]; ok {
				if server.Cmd != nil && server.Cmd.Process != nil {
					server.Cmd.Process.Kill()
				}
			}
			v.Status = 0
		}
	}
	return saveServerInfo()
}
