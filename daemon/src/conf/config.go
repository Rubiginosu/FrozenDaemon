/*
Powered by Axoford12
Rubiginosu | Freeze Team
本报提供了一些对于配置文件的生成和读取函数
*/
package conf

import (
	"colorlog"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	/**
	三种方法分别有各自的长处和缺点
	下面给出详细解释：
	1.Mount法:  .高效，带磁盘配额限制，但是假设你不小心误删掉服务器数据，会炸掉你的系统并且丢失所有数据，只能重装
	2.Copy法：   效率较低，但是设备安全，就算误删也无所谓，就算服务器被攻破也影响不了主机，原理是把所需要的目录全部复制一份...所以占用炸，效率炸
	3.Link法     效率高且安全，但是因为其不能用loop作为磁盘配额限制，如果您的机器并不想限制租用方的磁盘资源，请选用该方法
	如果您对Linux十分了解，知道lsof等使用方法，请采用mount 法，服务器所有数据丢失我们概不负责。
	若你还是锑一样的删掉了/bin /lib 等目录，然后跑来挂婊作者，那你还真是个十足的小学生。
	如果您不是很了解，甚至是Linux小白，请采用2,3方法。
	2,3方法如果需要限制磁盘要求，请选用Copy,如果不需要，讲求性能，请选用Link。它们都是安全的方法
	*/
	HDM_MOUNT = "Mount"
	HDM_COPY  = "Copy"
	HDM_LINK  = "Link"
)

type Cnf struct {
	ServerManager       serverManager
	DaemonServer        DaemonServer
	FileTransportServer FileTransportServer
}

type DaemonServer struct {
	Port                            int
	VerifyCode                      string
	DefaultBufLength                int
	ValidationKeyOutDateTimeSeconds float64
	UserId                          int
	HardDiskMethod                  string
	BlockDeviceMajMim               string
	NetworkCardName                 string
	PluginPath                      string
}

type serverManager struct {
	Servers       string
	Modules       string
	WebSocketPort int
	HardDisk      string
}

type FileTransportServer struct {
	Port int
}

func GetConfig(filename string) (Cnf, error) {
	file, err := os.Open(filename)
	if err != nil {
		return GenerateConfig("../conf/fg.json"), nil

	}
	var v Cnf
	b, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		return Cnf{}, err2
	}
	json.Unmarshal(b, &v)
	return v, nil
}

/**
生成一个配置文件，并提示用户输入磁盘信息及网卡信息
*/
func GenerateConfig(filepath string) Cnf {
	var v Cnf = Cnf{
		serverManager{"../data/servers.json", "../data/modules.json", 52024, "8:0"},
		DaemonServer{52023,
			RandString(64),
			256,
			20,
			100000,
			HDM_LINK, "", "","../plugins"}, // 为何选择52023？俺觉得23号这个妹纸很可爱啊
		FileTransportServer{52025},
	}
	file, err := os.Create(filepath)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	colorlog.WarningPrint("Your hardDisk Maj:Min not configured.")
	colorlog.PointPrint("This will help you input a correct Maj:Min information.")
	colorlog.PromptPrint("This program will show a lsblk info and you should choose a info like this:")
	colorlog.PromptPrint("")
	fmt.Println(
		`
NAME   MAJ:MIN RM   SIZE RO TYPE MOUNTPOINT
sda      ` + colorlog.ColorSprint("8:0", colorlog.FR_CYAN) + `    0 931.5G  0 disk
├─sda2   8:2    0 927.1G  0 part /
├─sda3   8:3    0   3.9G  0 part [SWAP]
└─sda1   8:1    0   512M  0 part /boot/efi
		`)
	colorlog.PromptPrint("And you should choose " + colorlog.ColorSprint("8:0", colorlog.FR_CYAN))
	for {
		majMin := ""

		cmd := exec.Command("lsblk")
		cmd.Env = os.Environ()
		output, err := cmd.CombinedOutput()
		if err != nil {
			colorlog.ErrorPrint(errors.New("Error occurred while run command lsblk. Error info:" + err.Error()))
			colorlog.LogPrint("Reason:" + string(output))
			os.Exit(-2) // 退出程序
		}
		colorlog.PromptPrint("Please choose: Your main hard disk Maj:Min\n")
		fmt.Println(string(output))
		colorlog.PromptPrint("Please input your hardDisk Maj:Min number.")
		fmt.Scanf("%s", &majMin)
		if err := validate(majMin, output, regexp.MustCompile("\\d+:\\d+")); err == nil {
			v.DaemonServer.BlockDeviceMajMim = majMin
			break
		} else {
			colorlog.PromptPrint("Your majMin may not valid.")
			colorlog.PromptPrint("Reason" + err.Error())
			colorlog.PromptPrint("the correct majMin likes this : ")
			colorlog.PromptPrint("8:0")
		}
	}
	colorlog.PointPrint("Maj:Min checked ok.")
	colorlog.PromptPrint("Please set network card name")
	for {
		name := ""

		cmd := exec.Command("ip", "a")
		cmd.Env = os.Environ()
		output, err := cmd.CombinedOutput()
		if err != nil {
			colorlog.ErrorPrint(errors.New("Error occurred while run command lsblk. Error info:" + err.Error()))
			colorlog.LogPrint("Reason:" + string(output))
			os.Exit(-2) // 退出程序
		}
		colorlog.PromptPrint("Please choose: Your main networkCard mame\n")
		fmt.Println(string(output))
		colorlog.PromptPrint("Please input:")
		fmt.Scanf("%s", &name)
		if err := validate(name, output, regexp.MustCompile(".+")); err == nil {
			v.DaemonServer.NetworkCardName = name
			break
		} else {
			colorlog.PromptPrint("Your networkCard name may not valid.")
			colorlog.PromptPrint("Reason" + err.Error())
			colorlog.PromptPrint("the correct networkCard name likes this : ")
			colorlog.PromptPrint("eth0")
		}
	}
	s, _ := json.MarshalIndent(v, "", "\t")
	file.Write(s)

	return v
}

// 用于获取一个随机字符串
func RandString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func validate(input string, output []byte, reg *regexp.Regexp) error {
	if !reg.Match([]byte(input)) {
		return errors.New("Not a correct format.")
	}
	if strings.Index(string(output), input) >= 0 {
		return nil
	} else {
		return errors.New("Must contain in output.")
	}
}
