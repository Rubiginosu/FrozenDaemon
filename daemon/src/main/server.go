package main

import (
	"flag"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
	"colorlog"
	"dmserver"
)

//#include<unistd.h>
//#include<malloc.h>
//void setug(int id){
//	while(setgid(id)!=0) sleep(1);
//  while(setuid(id)!=0) sleep(1);
//}
import "C"

func main() {
	var (
		uid     int    // 要运行这个的程序的身份uid，随便指定一个.
		command string // command: ping www.baidu.com
		proc    bool   // if proc -> mount /proc with "mount -t proc none /proc" in container
		sid     int    // Server id ,write to cgroups .
	)
	flag.IntVar(&uid, "uid", 0, "uid for setuid command")
	flag.StringVar(&command, "cmd", "", "Command to be run")
	flag.BoolVar(&proc, "proc", true, " if true -> Mounting proc dir.")
	flag.IntVar(&sid, "sid", 0, "The serverid 's config will be write with cgroups config")
	flag.Parse()
	// 命名空间
	syscall.Unshare(syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_FILES | syscall.CLONE_FS)
	root := "../servers/server" + strconv.Itoa(sid)
	syscall.Chroot(root)

	if proc {

		os.Mkdir("/proc", 555)
		cmd := exec.Command("/bin/mount", "-t", "proc", "none", "/proc")
		out,err := cmd.CombinedOutput()
		if err != nil {
			colorlog.ErrorPrint(err)
			dmserver.OutputErrReason(out)
			return
		}
	}

	os.Chdir("serverData") // 都已经Chroot了，思想要开放....
	err4 := syscall.Setgroups([]int{uid})
	if err4 != nil {
		panic(err4)
	}
	// 降权
	// 以上区域拥有root权限
	///////////////////////////////////////////////////////////
	C.setug(C.int(uid))

	// 正则分离 args
	// args[0] args[1] args[2] ....
	commands := regexp.MustCompile(" +").Split(command, -1)
	err5 := syscall.Exec(commands[0], commands, os.Environ()) // 鬼畜golang之特效
	if err5 != nil {
		panic(err5)
	}
}