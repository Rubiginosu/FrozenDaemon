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
	"os"
	"strconv"
	"time"
	"utils"
	"regexp"
	"os/exec"
)

var config conf.Cnf
var serverSaved = make(map[int]*ServerLocal, 0)
var servers = make(map[int]*ServerRun, 0)
var requestHandlers = make(map[string]*[]func([]byte) []byte,0)
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
	// Received ? Who cares?
	colorlog.PointPrint("Received " + colorlog.ColorSprint(request.Method, colorlog.FR_GREEN) + " Command!")
	if functions,ok := requestHandlers[request.Method];ok {
		colorlog.LogPrint("Calling plugin functions")
		for _,function := range *functions{
			resp := Response{}
			b,_ := json.Marshal(request)
			err := json.Unmarshal(function(b),&resp)
			if err != nil {
				colorlog.ErrorPrint("dmserver.handle during unmarshall json : ",err)
				return Response{-1,"Plugin override error."}
			}
			return resp
		}
	}
	switch request.Method {

	case "List":
		return outputListOfServers()
	case "Create":

		if _, ok := serverSaved[request.OperateID]; ok {
			return Response{-1, "Server id: " + strconv.Itoa(request.OperateID) + "has been token"}
		}

		serverSaved[request.OperateID] = &ServerLocal{
			ID:                    request.OperateID,
			Name:                  request.Message,
			Executable:            "",
			Status:                0,
			MaxCpuUtilizatioRate:  1,
			MaxMem:                1024,
			MaxHardDiskCapacity:   0,
			MaxHardDiskReadSpeed:  50,
			MaxHardDiskWriteSpeed: 50,
			MaxUnusedUpBandwidth:  100,
			MaxUsingUpBandwidth:   100,
			Expire:                time.Now().Unix() + 3600,
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
			if server.MaxHardDiskCapacity == 0 {
				return Response{-1, "Please set MaxHardDiskCapacity！"}
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
		colorlog.LogPrint("Stop command sent.")
		if server, ok := servers[request.OperateID]; ok {

			server.Close()
		} else {
			return Response{-1, "Invalid server id"}
		}
		return Response{0, "OK"}
	case "ExecInstall":

		colorlog.LogPrint("Try to auto install " + request.Message)
		conn, err := http.Get( "http://vae.fg.mcpe.cc/exec?name=" + request.Message + "&arch=" + getArch())
		if err != nil {
			return Response{-1, err.Error()}
		}
		defer conn.Body.Close()
		respData, err := ioutil.ReadAll(conn.Body)
		if err != nil {
			fmt.Println("Read body error")
			return Response{-1, err.Error()}
		}
		var resp Response
		err = json.Unmarshal(respData, &config)
		if err != nil {
			fmt.Println("Json Unmarshal error!")
			return Response{-1, err.Error()}
		}
		if resp.Status != 0 {
			return Response{-1, "return error:" + resp.Message}
		}
		var config ExecInstallConfig
		err = json.Unmarshal([]byte(resp.Message),&config)
		// 解析成功且没有错误
		go install(config)
		return Response{0, "OK,installing"}

	case "SetServerConfig":
		var elements []ServerAttrElement
		err := json.Unmarshal([]byte(request.Message), &elements)
		if err != nil {
			return Response{-1, "Json decoding error:" + err.Error()}
		}
		nums, err2 := setServerConfigAll(elements, request.OperateID)
		if err2 != nil {
			return Response{-1, err2.Error()}
		}
		return Response{0, fmt.Sprintf("OK,Setted %d element(s)", nums)}
	case "InputLineToServer":
		// 此方法将一行指令输入至服务端
		if server, ok := servers[request.OperateID]; ok {
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
		if server, ok := servers[request.OperateID]; ok {
			b, _ := json.Marshal(server.Players)
			return Response{-1, fmt.Sprintf("%s", b)}
		} else {
			return Response{-1, "Invalid server id"}
		}
	case "KeyRegister":
		auth.KeyRegister(request.Message, request.OperateID)
		return Response{0, "OK"}
	case "GetServerDir":
		if server, ok := serverSaved[request.OperateID]; ok {
			if validateOperateDir("../servers/server"+strconv.Itoa(server.ID)+"/serverData/", request.Message) {
				infos, err := ioutil.ReadDir("../servers/server" + strconv.Itoa(server.ID) + "/serverData/" + request.Message)
				if err != nil {
					colorlog.ErrorPrint("Server Reading dir: ID:" + strconv.Itoa(server.ID) ,err)
					return Response{-1, "Reading Dir :" + err.Error()}
				}
				b, _ := json.Marshal(buildServerInfos(infos))
				return Response{0, string(b)}
			} else {
				return Response{-1, "Permission denied."}
			}

		}
		return Response{-1, "Invalid server id"}
	case "DeleteServerFile":
		if server, ok := serverSaved[request.OperateID]; ok {
			if validateOperateDir("../servers/server"+strconv.Itoa(server.ID)+"/serverData/", request.Message) {
				err := os.RemoveAll("../servers/server" + strconv.Itoa(server.ID) + "/serverData/" + request.Message)
				if err != nil {
					colorlog.ErrorPrint("Server Reading dir error: ID:" + strconv.Itoa(server.ID) ,err)
					return Response{-1, "Reading Dir :" + err.Error()}
				}
				return Response{0, "Deleted dir."}
			} else {
				return Response{-1, "Permission denied."}
			}

		}
	case "GetExecList":
		info,err := ioutil.ReadDir("../exec")
		if err != nil {
			colorlog.ErrorPrint("read exec path",err)
			return Response{-1,"Reading dir :" + err.Error()}
		}
		result := make([]string,0)
		for _,v := range info {
			s := v.Name()
			cstr := utils.CString(s)
			if !v.IsDir() && cstr.Contains(".json"){
				result = append(result,regexp.MustCompile("\\.json$").ReplaceAllString(v.Name(),""))
			}
		}
		lang, err := json.Marshal(result)
		if err == nil {
			return Response{0, string(lang)}
		}else{
			fmt.Println("Json Unmarshal error!")
			return Response{-1, err.Error()}
		}

	}
	return Response{
		-1, "Unexpected err",
	}
}

func setServerConfigAll(attrs []ServerAttrElement, index int) (int, error) {
	res := 0
	// 设置该设置的Attrs
	if server, ok := serverSaved[index]; ok {
		// 判断被设置那个服务器是否存在于映射
		for i := 0; i < len(attrs); i++ {
			colorlog.LogPrint("Attempt to set " + colorlog.ColorSprint(attrs[i].AttrName, colorlog.FR_CYAN) + "=" + colorlog.ColorSprint(attrs[i].AttrValue, colorlog.FR_GREEN) + " to server" + strconv.Itoa(index))
			switch attrs[i].AttrName {
			case "MaxCpuRate":
				rate, err := strconv.Atoi(attrs[i].AttrValue)
				if err != nil {
					return -1, err
				}
				server.MaxCpuUtilizatioRate = rate
			case "MaxMem":
				mem, err := strconv.Atoi(attrs[i].AttrValue)
				if err != nil {
					return -1, err
				}
				server.MaxMem = mem
			case "Executable":
				server.Executable = attrs[i].AttrValue
			case "MaxHardDiskCapacity":
				disk, err := strconv.Atoi(attrs[i].AttrValue)
				if err != nil {
					return -1, err
				}
				server.MaxHardDiskCapacity = disk
			case "Name":
				server.Name = attrs[i].AttrValue
			case "Expire":
				expire, err := strconv.Atoi(attrs[i].AttrValue)
				if err != nil {
					return -1, err
				}
				server.Expire = server.Expire + int64(expire) - 3600
			case "MaxHardDiskWriteSpeed":
				speed, err := strconv.Atoi(attrs[i].AttrValue)
				if err != nil {
					return -1, err
				}
				server.MaxHardDiskWriteSpeed = speed
			case "MaxHardDiskReadSpeed":
				speed, err := strconv.Atoi(attrs[i].AttrValue)
				if err != nil {
					return -1, err
				}
				server.MaxHardDiskReadSpeed = speed
			case "MaxUnusedUpBandwidth":
				width, err := strconv.Atoi(attrs[i].AttrValue)
				if err != nil {
					return -1, err
				}
				server.MaxUnusedUpBandwidth = width
			case "MaxUsingUpBandwidth":
				width, err := strconv.Atoi(attrs[i].AttrValue)
				if err != nil {
					return -1, err
				}
				server.MaxUsingUpBandwidth = width
			default:
				colorlog.WarningPrint("Attr " + colorlog.ColorSprint(attrs[i].AttrName, colorlog.FR_RED) + " not found or cannot be set.")
			}
			res++
		}

		server.networkFlush()
		server.performanceFlush()
		return res, nil
	}

	return -1, errors.New("Err with invalid server id.")
}
func outputListOfServers() Response {
	b, _ := json.Marshal(serverSaved)
	return Response{0, string(b)}
}

func getArch() string {
	cmd := exec.Command("uname","-a")
	out,err := cmd.CombinedOutput()
	if err != nil {
		colorlog.ErrorPrint("run [uname -a]",err)
		utils.OutputErrReason(out)
	}
	cstr := utils.CString(string(out))
	if cstr.Contains("86_64"){
		return "64"
	} else {
		return "32"
	}
}