package conf

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)
const(
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
	HDM_COPY = "Copy"
	HDM_LINK = "Link"
)
type Config struct {
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

func GetConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return GenerateConfig("../conf/fg.json"), nil

	}
	var v Config
	b, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		return Config{}, err2
	}
	json.Unmarshal(b, &v)
	return v, nil
}

func GenerateConfig(filepath string) Config {
	file, err := os.Create(filepath)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	var v Config = Config{
		serverManager{"../data/servers.json", "../data/modules.json", 52024, "8:0"},
		DaemonServer{52023,
			RandString(64),
			256,
			20,
			100000,
		HDM_LINK}, // 为何选择52023？俺觉得23号这个妹纸很可爱啊
		FileTransportServer{52025},
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
