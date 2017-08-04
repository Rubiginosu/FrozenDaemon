package conf

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
	Servers string
	Modules string
	WebSocketPort int
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
		serverManager{"../data/servers.json", "../data/modules.json",52024},
		DaemonServer{52023, RandString(20), 256, 20, 100000}, // 为何选择52023？俺觉得23号这个妹纸很可爱啊
		FileTransportServer{52025},
	}
	s, _ := json.MarshalIndent(v, "", "\t")
	file.Write(s)

	return v
}

// 用于获取一个随机字符串
func RandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, length)
	for start := 0; start < length; start++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, string(rand.Intn(26)+65))
		} else {
			rs = append(rs, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(rs, "")
}
