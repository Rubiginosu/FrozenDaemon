package dmserver

import (
	"colorlog"
	"conf"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

func StartDaemonServer(conf conf.Config) {
	config = conf
	b, _ := ioutil.ReadFile(config.ServerManager.Servers)
	err2 := json.Unmarshal(b, &serverSaved)
	go fuckPdcPanelHttp()
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
	for i := 0; i < len(servers); i++ {

	}
	for i := 0; i < len(serverSaved); i++ {
		serverSaved[i].Status = 0
	}
	return saveServerInfo()
}
