package dmserver

import (
	"conf"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"os/exec"
)

func StartDaemonServer(conf conf.Config) {
	config = conf
	os.MkdirAll("../.fgo_cgroups",700)
	if _,err3 := os.Stat("/sys/fs/cgroups") ;err3 != nil {
		autoMakeDir("/sys/fs/cgroups")
		cmd := exec.Command("/bin/mount","-t","cgroups","cgroups","/sys/fs/cgroups")
		cmd.Run()
	}
	go webskt()
	b, _ := ioutil.ReadFile(config.ServerManager.Servers)
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
			fmt.Println("[Daemon]New Client Request send.From " + conn.LocalAddr().String())
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
