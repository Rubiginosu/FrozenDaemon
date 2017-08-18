package dmserver

import (
	"os"
	"colorlog"
	"os/exec"
	"strconv"
	"utils"
)

func (s *ServerLocal)prepareVirtualEnv(){
	lfs := "../servers/server" + strconv.Itoa(s.ID) + "/"
	dirs := []string{"dev","proc","sys","run"}
	for _,v := range dirs{
		makeUserOwnedDir(lfs + v)
	}
	// 创建初始化设备节点
	if _,err := os.Stat(lfs + "/dev/console");err != nil{
		colorlog.LogPrint("Checked /dev/console may need mount")
		cmdMkNodConsole := exec.Command("mknod","-m","600",lfs + "dev/console","c","5","1")
		utils.AutoRunCmdAndOutputErr(cmdMkNodConsole,"make virtual node /dev/console")
	}

	if _,err := os.Stat(lfs + "/dev/null");err != nil{
		colorlog.LogPrint("Checked /dev/null may need mount")
		cmdMkNodNull := exec.Command("mknod","-m","666",lfs + "dev/null","c","1","3")
		utils.AutoRunCmdAndOutputErr(cmdMkNodNull,"make virtual node /dev/null")
	}

	// 挂载/dev
	if _,err := os.Stat(lfs + "/dev/core");err != nil{
		colorlog.LogPrint("Checked /dev may need mount")
		cmdMountDev := exec.Command("mount","-v","--bind","/dev",lfs + "/dev")
		utils.AutoRunCmdAndOutputErr(cmdMountDev,"mount virtual file system dev")
	}

	// 挂载虚拟文件系统
	// devpts
	if _,err := os.Stat(lfs + "/dev/pts");err != nil{
		colorlog.LogPrint("Checked /dev/pts need to be mount")
		cmd := exec.Command("mount","-vt","devpts","devpts",lfs + "/dev/pts","-o","gid=" + strconv.Itoa(config.DaemonServer.UserId) + ",mode=620")
		utils.AutoRunCmdAndOutputErr(cmd,"mount devpts")
	}
	// proc
	if _,err := os.Stat(lfs + "/proc/cpuinfo");err != nil{
		colorlog.LogPrint("Checked /proc need to be mount")
		cmd := exec.Command("mount","-vt","proc","proc",lfs + "/proc")
		utils.AutoRunCmdAndOutputErr(cmd,"mount proc")
	}


	// sysfs
	if _,err := os.Stat(lfs + "/sys/fs");err != nil{
		colorlog.LogPrint("Checked /sys need to be mount")
		cmd := exec.Command("mount","-vt","sysfs","sysfs",lfs + "/sys")
		utils.AutoRunCmdAndOutputErr(cmd,"mount sys")
	}

	// tmpfs
	colorlog.LogPrint("Checked /run need to be mount")
	cmd := exec.Command("mount","-vt","tmpfs","tmpfs",lfs + "/run")
	cmd.Env = os.Environ()
	cmd.Run()

}
func makeUserOwnedDir(dirname string){
	os.MkdirAll(dirname,777)
	os.Chown(dirname,config.DaemonServer.UserId,config.DaemonServer.UserId)
}


