package dmserver

import (
	"io"
	"os/exec"
)

type Request struct {
	Method    string
	OperateID int
	Message   string
}

type Response struct {
	Status  int
	Message string
}

type InterfaceRequest struct {
	Auth string
	Req  Request
}

type ExecInstallConfig struct {
	Rely      []Module
	Success   bool
	Timestamp int
	Url       string
	StartConf ExecConf
	Message   string
	Md5       string
}

type ExecConf struct {
	Name                 string
	Command              string // 开服的指令
	StartServerRegexp    string // 判定服务器成功开启的正则表达式
	NewPlayerJoinRegexp  string // 判定新人加入的表达式
	PlayExitRegexp       string // 判定有人退出的表达式
	StoppedServerCommand string // 服务器软退出指令
	Mount                []string
	ProcDir              bool // 需不需要特殊挂载proc 即在容器环境中，是否需要mount -t proc none /proc
}

type Module struct {
	Name     string
	Download string
	Chmod    string
	Md5      string
}

type ServerLocal struct {
	ID         int
	Name       string
	Executable string
	Status     int
	/*
	Status  enum 应该只有
	0 - 正常未运行
	1 - 正常运行
	2 - 正在开启
	 */
	MaxCpuCores int // CPU核心数
	MaxMem      int // 最大内存
	MaxHardDisk int // 磁盘空间，开服时就必须设定好，以后不允许改变.
}

type ServerRun struct {
	ID         int
	Players    []string
	Cmd        *exec.Cmd
	BufLog     [][]byte
	StdinPipe  *io.WriteCloser
	StdoutPipe *io.ReadCloser
}
type ServerAttrElement struct {
	// Set server.AttrName = AttrValue
	// eg server.MaxMemory = 1024
	AttrName  string // Attribute Name
	AttrValue string // Attribute Value
}
