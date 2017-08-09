package dmserver

import (
	"auth"
	"colorlog"
	"conf"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

var config conf.Config
var serverSaved = make(map[int]*ServerLocal, 0)
var servers = make(map[int]*ServerRun, 0)

func connErrorToExit(errorInfo string, c net.Conn) {
	res, _ := json.Marshal(Response{-1, errorInfo})
	c.Write(res)
	c.Close()
}

// 保存服务器信息
func saveServerInfo() error {
	b, err := json.Marshal(serverSaved)
	if err != nil {
		return err
	}
	ioutil.WriteFile(config.ServerManager.Servers, b, 0664)
	return nil
}

// 处理本地命令

func handleConnection(c net.Conn) {
	buf := make([]byte, config.DaemonServer.DefaultBufLength)
	length, _ := c.Read(buf)
	request := InterfaceRequest{}
	err := json.Unmarshal(buf[:length], &request)
	if err != nil {
		connErrorToExit(err.Error(), c)
	}
	if request.Auth == config.DaemonServer.VerifyCode {
		res, _ := json.Marshal(handleRequest(request.Req))
		c.Write(res)
		c.Close()
	} else {
		connErrorToExit("No command or Auth error.", c)
	}

}

// 命令处理器
func handleRequest(request Request) Response {
	colorlog.PointPrint("Recevied " + colorlog.ColorSprint(request.Method, colorlog.FR_GREEN) + " Command!")
	switch request.Method {

	case "List":
		return outputListOfServers()
	case "Create":

		if _, ok := serverSaved[request.OperateID]; ok {
			return Response{-1, "Server id:" + strconv.Itoa(request.OperateID) + "has been token"}
		}

		serverSaved[request.OperateID] = &ServerLocal{
			request.OperateID,
			"test",
			"",
			0,
			1,
			1024,
			0,
			time.Now().Unix() + 3600,
		}
		// 序列化b来储存。
		b, err := json.MarshalIndent(serverSaved, "", "\t")

		// 新创建的服务器写入data文件
		err2 := ioutil.WriteFile(config.ServerManager.Servers, b, 0666)
		if err2 != nil {
			return Response{
				-1,
				err2.Error(),
			}
		}
		if err != nil {
			return Response{
				-1,
				err.Error(),
			}
		}

		return Response{
			0,
			"OK",
		}
		serverSaved[0].MaxMem = 1
	case "Delete":
		if _, ok := serverSaved[request.OperateID]; ok {
			serverSaved[request.OperateID].Delete()
		} else {
			return Response{-1, "Invalid server id."}
		}

		return Response{0, "OK"}
	case "Start":

		// 运行这个服务器
		if server, ok := serverSaved[request.OperateID]; ok {
			if server.Status != SERVER_STATUS_CLOSED {
				return Response{-1, "Server Running or staring"}
			}
			if server.MaxHardDisk == 0 {
				return Response{-1,"Please set MaxHardDisk！"}
			}
			err2 := server.EnvPrepare()
			if err2 != nil {
				colorlog.ErrorPrint(err2)
				return Response{
					-1,"Env prepare error",
				}
			}
			err := server.Start()
			if err == nil {
				return Response{
					0, "OK",
				}
			} else {
				return Response{-1, err.Error()}
			}
		} else {
			return Response{-1, "Invalid server id."}
		}
	case "Stop":
		if server, ok := servers[request.OperateID]; ok {
			server.Close()
		} else {
			return Response{-1, "Invalid server id"}
		}
		return Response{0, "OK"}
	case "ExecInstall":

		colorlog.LogPrint("Try to auto install id:" + strconv.Itoa(request.OperateID))
		colorlog.LogPrint("From " + request.Message)
		conn, err := http.Get(request.Message + "?action=GetJar&JarID=" + strconv.Itoa(request.OperateID))
		if err != nil {
			fmt.Println("Get ExecInstallConfig error!")
			return Response{-1, err.Error()}
		}
		defer conn.Body.Close()
		respData, err2 := ioutil.ReadAll(conn.Body)
		if err2 != nil {
			fmt.Println("Read body error")
			return Response{-1, err2.Error()}
		}
		var config ExecInstallConfig
		err3 := json.Unmarshal(respData, &config)
		if err2 != nil {
			fmt.Println("Json Unmarshal error!")
			return Response{-1, err3.Error()}
		}
		if !config.Success {
			return Response{-1, "Get exec data error:" + config.Message}
		}
		// 解析成功且没有错误
		go install(config)
		return Response{0, "OK,Installing"}

	case "SetServerConfig":
		var elements []ServerAttrElement
		err := json.Unmarshal([]byte(request.Message), &elements)
		if err != nil {
			return Response{-1, "Json decoding error:" + err.Error()}
		}
		err2 := setServerConfigAll(elements, request.OperateID)
		if err2 != nil {
			return Response{-1, err2.Error()}
		}
		return Response{0, fmt.Sprintf("OK,Setted %d element(s)", len(elements))}
	case "GetServerConfig":
		// 获取服务器信息（已保存信息）
		if server, ok := serverSaved[request.OperateID]; ok {
			b, _ := json.Marshal(server)
			return Response{0, fmt.Sprintf("%s", b)}
		}

	case "InputLineToServer":
		// 此方法将一行指令输入至服务端
		if server,ok := servers[request.OperateID]; ok  {
			err := server.inputLine(request.Message)
			if err != nil {
				return Response{-1, err.Error()}
			} else {
				return Response{0, "Send OK"}
			}
		} else {
			return Response{-1, "Invalid Server id"}
		}
	case "GetServerPlayers":
		if server,ok := servers[request.OperateID]; ok {
			b, _ := json.Marshal(server.Players)
			return Response{-1, fmt.Sprintf("%s", b)}
		} else {
			return Response{-1, "Invalid server id"}
		}
	case "KeyRegister":
		auth.KeyRigist(request.Message, request.OperateID)
		return Response{0, "OK"}
	}
	return Response{
		-1, "Unexpected err",
	}
}

func setServerConfigAll(attrs []ServerAttrElement, index int) error {
	// 设置该设置的Attrs
	if server, ok := serverSaved[index]; ok {
		// 判断被设置那个服务器是否存在于映射
		for i := 0; i < len(attrs); i++ {
			switch attrs[i].AttrName {
			case "MaxMemory":
				mem, err := strconv.Atoi(attrs[i].AttrValue)
				if err != nil {
					return err
				}
				server.MaxMem = mem
			case "Executable":
				server.Executable = attrs[i].AttrValue
			case "MaxHardDisk":
				disk, err := strconv.Atoi(attrs[i].AttrValue)
				if err != nil {
					return err
				}
				server.MaxHardDisk = disk
			case "Name":
				server.Name = attrs[i].AttrValue
			case "Expire":
				expire, err := strconv.Atoi(attrs[i].AttrValue)
				if err != nil {
					return err
				}
				server.Expire = server.Expire + int64(expire) - 3600
			}
		}
	}

	return errors.New("Err with invalid server id.")
}
func outputListOfServers() Response {
	b, _ := json.Marshal(serverSaved)
	return Response{0, string(b)}
}
