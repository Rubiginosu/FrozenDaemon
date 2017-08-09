package conf

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
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
		serverManager{"../data/servers.json", "../data/modules.json", 52024,"8:0"},
		DaemonServer{52023, RandString(64), 256, 20, 100000}, // 为何选择52023？俺觉得23号这个妹纸很可爱啊
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
