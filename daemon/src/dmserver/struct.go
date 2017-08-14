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
	Link                 []string
	ProcDir              bool // 需不需要特殊挂载proc 即在容器环境中，是否需要mount -t proc none /proc
}

type Module struct {
	Name     string
	Download string
	Chmod    string
	Md5      string
}

type ServerLocal struct {
	//  --- 表示该属性与的AttrName与原属性的名称相同
	ID         int            // 不可设置，创建时指定
	Name       string        // 可设置 ----
	Executable string        // 可设置 ----
	Status     int
	/*
		Status  enum 应该只有
		0 - 正常未运行
		1 - 正常运行
		2 - 正在开启
	*/
	MaxCpuUtilizatioRate  int // CPU使用率                                   可设置     名称：MaxCpuRate  因为名字太长所以简写
	MaxMem                int // 最大内存                                    可设置     ------
	MaxHardDiskCapacity   int // 磁盘空间，开服时就必须设定好，以后不允许改变.  仅第一次   ------
	MaxHardDiskReadSpeed  int // 磁盘最大读速率                              可设置      ------
	MaxHardDiskWriteSpeed int // 磁盘最大写速率                               可设置     ------
	// 读写速率单位均为 M/s
	MaxUpBandwidth int   // 最大上行带宽 单位 Mb/s // b ： bit.               可设置     ------
	Expire         int64 // 过期时间，Unix时间戳                              可设置    ------
	// ########### 设置 expire ，服务端会将其加入到开服时间，单位秒，
	//             Example 设置 3600 服务端将在服务器被创建一小时后直接删除...该删除没有提示
	//             Daemon 每隔五秒进行一次过期检测，清理所有过期服务器并给出Log信息
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
type ServerPathFileInfo struct {
	Name    string // 文件名
	Mode    string // 文件Unix模式位
	ModTime int64  // 文件修改时间
}
