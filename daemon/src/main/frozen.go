package main

import (
	"auth"
	"conf"
	"dmserver"
	"encoding/json"
	"filetrans"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
	"colorlog"
)

const VERSION string = "v0.3.1"
const FILE_CONFIGURATION string = "../conf/fg.json"
const UPDATE_CURRENT_VERSION = "https://raw.githubusercontent.com/Rubiginosu/frozen-go/master/VERSION"

var config conf.Config

func main() {
	if !(len(os.Args) > 1 && os.Args[1] == "-jump") {
		printInfo()
	} // 如果需要调试本程序，那么加上-jump参数可以跳过打印.
	if !isRoot() {
		fmt.Println(colorlog.ColorSprint("Need root permission.",colorlog.FR_RED))
		return
	}
	colorlog.LogPrint("Reading config file")
	config, _ = conf.GetConfig(FILE_CONFIGURATION)
	colorlog.LogPrint("Configuration get done")
	colorlog.LogPrint("Checking Update")
	if versionCode, err := checkUpdate(); err != nil {
		colorlog.ErrorPrint(err)
	} else {
		colorlog.LogPrint("Version Check done:")
		if versionCode > 1 {
			colorlog.WarningPrint("|---Daemon out of date")
			colorlog.WarningPrint("|---Your daemon need to be updated!")
			return
		} else if versionCode == 1 {
			colorlog.WarningPrint("Small bugs fixed,You choose to updated it or not.")
		} else {
			colorlog.LogPrint("Lastest Version")
		}
	}
	colorlog.PointPrint("Starting Server Manager.")
	go dmserver.StartDaemonServer(config)
	go filetrans.ListenAndServe(config)
	colorlog.PointPrint("Starting websocket server")
	go dmserver.Webskt()
	colorlog.PointPrint("Starting ValidationKeyUpdater.")
	go auth.ValidationKeyUpdate(config.DaemonServer.ValidationKeyOutDateTimeSeconds)
	colorlog.LogPrint("Done,type \"?\" for help. ")
	for {
		var s string
		fmt.Scanf("%s", &s)
		processLocalCommand(s)
	}
}

func printInfo() {
	fmt.Println(colorlog.ColorSprint(`

    ______                                ______
   / ____/_____ ____  ____  ___   ____   / ____/____
  / /_   / ___// __ \/_  / / _ \ / __ \ / / __ / __ \
 / __/  / /   / /_/ / / /_/  __// / / // /_/ // /_/ /
/_/    /_/    \____/ /___/\___//_/ /_/ \____/ \____/


	`,colorlog.FR_CYAN))
	time.Sleep(2 * time.Second)
	fmt.Println("---------------------")
	time.Sleep(100 * time.Microsecond)
	fmt.Print("Powered by ")
	for _, v := range []byte("Axoford12") {
		time.Sleep(240 * time.Millisecond)
		fmt.Print(colorlog.ColorSprint(string(v),colorlog.BK_GREEN))
	}
	fmt.Println()
	time.Sleep(1000 * time.Millisecond)
	time.Sleep(100 * time.Microsecond)
	fmt.Println("---------------------")
	time.Sleep(300 * time.Millisecond)
	colorlog.LogPrint("version:" + VERSION)
	time.Sleep(1 * time.Second)
}

func processLocalCommand(c string) {
	switch c {
	case "stop":
		fmt.Println("Stopping...")
		dmserver.StopDaemonServer()
		os.Exit(0)
	case "?":
		fmt.Println("FrozenGo" + VERSION + " Help Manual -- by Axoford12")
		fmt.Println("stop: Stop the daemon.save server changes.")
		fmt.Println("status: Echo server status.")
		return
	case "status":
		b, _ := json.Marshal(dmserver.GetServerSaved())
		fmt.Println(string(b))
		return
	}
}
func isRoot() bool {
	nowUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	userId, err2 := strconv.Atoi(nowUser.Uid)
	if err2 != nil {
		panic(err)
	}
	return userId == 0
}
func checkUpdate() (int, error) {
	colorlog.LogPrint("Starting Version check...")
	colorlog.LogPrint("This may take more time..")
	resp, err := http.Get(UPDATE_CURRENT_VERSION)
	if err != nil {
		return -2, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	body = []byte(strings.TrimRight(string(body), "\n\r"))
	nowVersion := []byte(VERSION)
	if body[1] != nowVersion[1] {
		return 3, nil
	} else if body[3] != nowVersion[3] {
		return 2, nil
	} else if body[5] != body[5] {
		return 1, nil
	} else {
		return 0, nil
	}
	//return -2,errors.New("Unexpected error")
}
