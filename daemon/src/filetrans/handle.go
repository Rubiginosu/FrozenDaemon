package filetrans

import (
	"dmserver"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	COMMAND_AUTH     = "AUTH"
	COMMAND_LIST     = "LIST"
	COMMAND_UPLOAD   = "UPLO"
	COMMAND_DOWNLOAD = "DOLO"
)

func (c *Command) handleCommand(conn net.Conn, serverid int) {
	serverPath := "../servers/server" + strconv.Itoa(serverid) + "/"
	switch c.Command {
	case COMMAND_LIST:
		dir, err := os.Open(serverPath)
		files, err2 := dir.Readdir(-1)
		if err != nil || err2 != nil {

			sendMessage(conn, SERVER_SERVER__ERROR)
		} else {
			var localFiles []localServerFile
			for i := 0; i < len(files); i++ {
				localFiles = append(localFiles, parseFileInfoToLocalFile(files[i]))
			}
			b, _ := json.Marshal(localFiles)
			sendMessage(conn, string(b))
		}

	case COMMAND_UPLOAD:
		fixedpath := path.Clean(serverPath + c.Args)
		if strings.Index(fixedpath, serverPath) != 0 {
			sendMessage(conn, SERVER_FILE_ERROR+"Permission Denied")
		} else {
			file, err := os.Create(fixedpath)
			if err != nil {
				panic(err)
			}
			io.Copy(file, conn)
			file.Close()
			sendMessage(conn, SERVER_UPDATE_OK)
		}
	case COMMAND_DOWNLOAD:
		fixedpath := path.Clean(serverPath + c.Args)
		if strings.Index(fixedpath, serverPath) != 0 {
			sendMessage(conn, SERVER_FILE_ERROR+"Permission Denied")
		} else {
			file, err := os.Open(fixedpath)
			if err != nil {
				fmt.Println(err.Error())
				sendMessage(conn, SERVER_FILE_ERROR+err.Error())
			}
			io.Copy(conn, file)
			sendMessage(conn, SERVER_FILE_SEND_FINISHED)
		}
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	sendMessage(c, SERVER_WELCOME)
	var serverid int
	for {
		command := parseCommandArg([]byte(getMessage(c)))
		if command.Command == "" {
			return
		} else if serverid = command.auth(); serverid >= 0 {
			sendMessage(c, SERVER_AUTH_SUCCEEDED)
			break // 验证成功就跳出循环
		} else if command.Command != COMMAND_AUTH {
			fmt.Println(command.Command)
			sendMessage(c, SERVER_PLEASE_AUTH)
		} else {
			// 告诉某人认证失败。
			sendMessage(c, SERVER_AUTH_FAILED)
			b, _ := json.Marshal(command)
			sendMessage(c, string(b))
		}
	}
	for {
		command := parseCommandArg([]byte(getMessage(c)))
		command.handleCommand(c, serverid)
	}
}

func (c *Command) auth() int {
	// 需要先判断一下下
	// 把Arg分隔 按  | 成两部分
	args := strings.SplitN(c.Args, "|", 2)
	if len(args) < 2 {
		// 这个args可能并没有分隔（第三方客户端？）
	} else if serverID, err := strconv.Atoi(args[0]); err != nil {
	} else if _, ok := dmserver.GetServerSaved()[serverID]; !ok {
		// 卧槽服务器又不可用
		// 这年头考虑熊孩子要考虑的真的够多..
		//} else if auth.IsVerifiedValidationKeyPair(serverID, args[1]) {
		//	// 哎，貌似是考虑完了。。
		//	return serverID
		//} else {
	}
	return -1
}
