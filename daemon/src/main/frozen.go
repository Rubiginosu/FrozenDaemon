package main

import (
	"colorlog"
	"conf"
	"dmserver"
	"encoding/json"
	"fmt"
	"ftrans"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"time"
	"fgplugin"
	"C"
)

/*
VERSION定义了FrozenGo的当前版本
 */
const VERSION string = "v1.0.1"

/*
定义了项目的配置文件目录
 */
const FILE_CONFIGURATION string = "../conf/fg.json"



/*
全局变量conf.
 */
var config conf.Cnf
func main() {
	// 检测是否有定义FGO_DEBUG变量,定义了就直接跳过Banner打印
	if os.Getenv("FGO_DEBUG") != "Yes" {
		banner()
	}
	// 检查是否为root
	if !isRoot() {
		fmt.Println(colorlog.ColorSprint("Need root permission.", colorlog.FR_RED))
		return
	}
	checkEnv()
	if _,err := os.Stat("../cgroup");err != nil{
		b := []byte(`
case $1 in
	"cg") case $2 in
			"init")
				mkdir /sys/fs/cgroup/cpu/${3} /sys/fs/cgroup/memory/${3} /sys/fs/cgroup/blkio/${3} /sys/fs/cgroup/net_cls/${3}
				# ${4} cpu $5 memmax $6 $7 $8 blkio blkio.throttle.read_bps_device $9 netcls
				tmp=$(cat /sys/fs/cgroup/cpu/${3}/cpu.cfs_period_us)
				cpux=$4
				rmb=$6
				wmb=$7
				mmb=$5
				let rmb=1024*1024*rmb
				let wmb=1024*1024*wmb
				let mmb=1024*1024*mmb
				let tmp=tmp*cpux/100
				echo $tmp > /sys/fs/cgroup/cpu/${3}/cpu.cfs_quota_us
				echo $mmb > /sys/fs/cgroup/memory/${3}/memory.max_usage_in_bytes
				echo "0x0001${9}" > /sys/fs/cgroup/net_cls/${3}/net_cls.classid
				echo "${8} ${rmb}" > /sys/fs/cgroup/blkio/${3}/blkio.throttle.read_bps_device
				echo "${8} ${wmb}" > /sys/fs/cgroup/blkio/${3}/blkio.throttle.write_bps_device
				;;
			"del")
				rmdir /sys/fs/cgroup/cpu/${3} /sys/fs/cgroup/memory/${3} /sys/fs/cgroup/blkio/${3} /sys/fs/cgroup/net_cls/${3}
				;;
			"run")
				/bin/echo ${4} |tee /sys/fs/cgroup/cpu/${3}/tasks /sys/fs/cgroup/memory/${3}/tasks /sys/fs/cgroup/blkio/${3}/tasks /sys/fs/cgroup/net_cls/${3}/tasks
				;;
			esac;;
	"net") DEV=$6;
			case $2 in
				"add")
				tc class add dev $DEV parent 1: classid 1:${3} htb rate ${4}mbit ceil ${5}mbit;
				tc filter add dev $DEV protocol ip parent 1:0 prio 1 handle 1:${3} cgroup;;
				"change")
				tc class change dev $DEV parent 1: classid 1:${3} htb rate ${4}mbit ceil ${5}mbit;;
				"del")
				tc class del dev $DEV parent 1: classid 1:${3};
				tc filter del dev $DEV protocol ip parent 1:0 prio 1 handle 1:${3} cgroup;;

			esac;;
	"init")
	DEV=$2;
	#tc qdisc del dev $DEV root
	tc qdisc add dev $DEV root handle 1: htb;
	tc class add dev $DEV parent 1: classid 1: htb rate 10000mbit ceil 10000mbit;
	service cgconfig restart;
	;;
	esac`)
		os.MkdirAll("../cgroup",755)
		ioutil.WriteFile("../cgroup/cg.sh",b,755)
	}
	colorlog.LogPrint("Reading config file")
	config, _ = conf.GetConfig(FILE_CONFIGURATION)
	if config.DaemonServer.HardDiskMethod == conf.HDM_MOUNT {
		colorlog.WarningPrint("You are running in MOUNT HardDisk Method!")
		colorlog.LogPrint("You must know its risk and willing to take responsibility for incorrect use")
		colorlog.WarningPrint(fmt.Sprint("Please type:", colorlog.ColorSprint("I_KNOW", colorlog.FR_PURPLE)))
		check := ""
		fmt.Scanf("%s", &check)
		if check != "I_KNOW" {
			colorlog.WarningPrint("You must know the warning")
			return
		}
	}
	colorlog.LogPrint("Configuration file got.")
	colorlog.LogPrint("This version is:"+VERSION)
	//去除了版本更新，不再会强制要求更新
	//2018年1月进行了修改
	colorlog.PointPrint("Loading plugins...")
	fgplugin.LoadPlugin(config.DaemonServer.PluginPath)
	colorlog.PointPrint("Starting Server Manager...")
	go dmserver.StartDaemonServer(config)
	go ftrans.Start(config)
	colorlog.PointPrint("Starting websocket server...")
	go dmserver.Webskt()
	colorlog.PointPrint("Starting ValidationKeyUpdater...")
	colorlog.LogPrint("Done,type \"?\" for help. ")
	// 处理一些非常非常基本的指令,因为基本不用,所以并不是很想写这一块的内容
	for {
		var s string
		fmt.Scanf("%s", &s)
		processLocalCommand(s)
	}
}

func banner() {
	fmt.Println(colorlog.ColorSprint(`

    ______                                ______
   / ____/_____ ____  ____  ___   ____   / ____/____
  / /_   / ___// __ \/_  / / _ \ / __ \ / / __ / __ \
 / __/  / /   / /_/ / / /_/  __// / / // /_/ // /_/ /
/_/    /_/    \____/ /___/\___//_/ /_/ \____/ \____/


	`, colorlog.FR_CYAN))
	time.Sleep(2 * time.Second)
	fmt.Println("---------------------")
	fmt.Println("Compliance with the MIT open source protocol...")
	time.Sleep(100 * time.Microsecond)
	fmt.Print("Powered by ")
	for _, v := range []byte("Axoford12") {
		time.Sleep(240 * time.Millisecond)
		fmt.Print(colorlog.ColorSprint(string(v), colorlog.FR_GREEN))
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
/*func checkUpdate() (int, error) {
	colorlog.LogPrint("Starting Version check...")
	colorlog.LogPrint("This may take more time...")
	resp, err := http.Get(UPDATE_CURRENT_VERSION + "?v=" + VERSION)
	if err != nil {
		return -2, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	result,err := strconv.Atoi(string(body))
	if err != nil {
		return -2,err
	}
	return result,nil
	//return -2,errors.New("Unexpected error")
}
*/
/*
 更新部分，1.0.1以后不再需要
*/

func checkEnv(){
	os.MkdirAll("../plugins",755)
	os.MkdirAll("../exec",755)
	os.MkdirAll("../data",755)
	os.MkdirAll("../conf",755)
}
