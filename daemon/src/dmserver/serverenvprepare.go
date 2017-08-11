package dmserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// 准备环境
func (server *ServerLocal) EnvPrepare() error {
	cmd := exec.Command("/bin/bash", "../cgroup/cg.sh", "cg",
		"init",
		"server"+strconv.Itoa(server.ID),
		strconv.Itoa(server.MaxCpuCores),
		strconv.Itoa(server.MaxMem), "10", "10", "",
		strings.Replace(fmt.Sprintf("%4x", server.ID), " ", "0", -1))
	cmd.Env = os.Environ()
	//  上面的替换是让服务器的id替换为四位十六进制id
	cmd.Run()
	fmt.Printf("Preparing server runtime for ServerID:%d \n", server.ID)
	serverDataDir := "../servers/server" + strconv.Itoa(server.ID) // 在一开头就把serverDir算好，增加代码重用
	// 文件夹不存在则创建文件夹
	autoMakeDir(serverDataDir + "/serverData")

	if _, err0 := os.Stat(serverDataDir + ".loop"); err0 != nil { //检查loop回环文件是否存在，如果不存在则创建
		fmt.Println("No loop file found!")
		//  新增 loop
		if server.MaxHardDisk == 0 {
			server.MaxHardDisk = 10240
		}
		cmd := exec.Command("/bin/dd", "if=/dev/zero", "bs=1024", // MaxHardDisk单位kb
			"count="+strconv.Itoa(server.MaxHardDisk), "of=../servers/server"+strconv.Itoa(server.ID)+".loop")
		fmt.Print("Writing file...")
		err := cmd.Run()
		if err != nil {
			return err
		}
		fmt.Println("Done")
		// 用mkfs格式化
		fmt.Println("Formatting...")
		cmd2 := exec.Command("/sbin/mkfs.ext4", serverDataDir+"/server"+strconv.Itoa(server.ID)+".loop")
		err2 := cmd2.Run()
		fmt.Println("Done")
		if err2 != nil {
			fmt.Println(err2)
			return err2
		}

	}
	fmt.Println("Preparing server data dir.")
	// 为挂载文件夹做好准备
	autoMakeDir(serverDataDir + "/lib")
	autoMakeDir(serverDataDir + "/execPath")
	//execPath,_ := filepath.Abs("../exec")
	//cmd2 := exec.Command("/bin/mount","-o","bind",execPath,serverDataDir + "/execPath")
	//cmd2.Run()
	if _, err := os.Stat("/lib64"); err == nil { // 32位系统貌似没有lib64,那就不新建了
		autoMakeDir(serverDataDir + "/lib64")
		// 这个谁说的准？ 哈哈～
	}
	// 挂载回环文件
	fmt.Println("Mounting loop file")
	cmd3 := exec.Command("/bin/mount", "-o", "loop", serverDataDir+"/server"+strconv.Itoa(server.ID)+".loop", serverDataDir)
	cmd3.Run()

	//err := server.mountDirs() // 挂载其他文件
	//if err != nil {
	//	fmt.Println("[ERROR]" + err.Error())
	//}
	// 挂载结束
	/////////////////////////////////////////////////////////
	return nil
}

func (server *ServerLocal) loadExecutableConfig() (ExecConf, error) {
	var newServerRuntimeConf ExecConf
	b, err := ioutil.ReadFile("../exec/" + server.Executable + ".json") // 将配置文件读入
	if err != nil {
		// 若在读文件时就有异常则停止反序列化
		return newServerRuntimeConf, err
	}
	err2 := json.Unmarshal(b, &newServerRuntimeConf) //使用自带的json库对读入的东西反序列化
	if err2 != nil {
		return newServerRuntimeConf, err
	}
	return newServerRuntimeConf, nil // 返回结果
}

func (server *ServerLocal) mountDirs() error {
	serverDataDir := "../servers/server" + strconv.Itoa(server.ID) // 在一开头就把serverDir算好，增加代码重用
	execConfig, err := server.loadExecutableConfig()
	fmt.Println(execConfig)
	if err != nil {
		return err
	}
	cmd := exec.Command("/bin/mount", "-o", "bind", "/lib", serverDataDir+"/lib")
	cmd.Run()
	cmdMountBin := exec.Command("/bin/mount", "-o", "bind", "/bin", serverDataDir+"/bin")
	cmdMountBin.Run() // 在这一版本中，将会强制挂载bin目录
	if _, err := os.Stat("/lib64"); err == nil {
		// 这里不用serverDataDir是处于安全考虑，万一小天才给我在../新建了一个lib64 那我把没有的lib64挂载过来就纯属多此一举了
		cmd := exec.Command("/bin/mount", "-o", "bind", "/lib64", serverDataDir+"/lib64")
		cmd.Run()
	}
	// TODO Link 方式实现Mount ,解决小白租服商会出现的类似于/bin目录被删掉之类的等等问题
	return nil
}
