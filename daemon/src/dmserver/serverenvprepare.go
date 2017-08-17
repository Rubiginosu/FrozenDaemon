package dmserver

import (
	"colorlog"
	"conf"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// 准备环境
func (server *ServerLocal) EnvPrepare() error {
	defer os.Chown("../servers/server"+strconv.Itoa(server.ID), config.DaemonServer.UserId, 0) // 让他有权限访问.
	networkID := strings.Replace(fmt.Sprintf("%4x", server.ID), " ", "0", -1)
	//  上面的替换是让服务器的id替换为四位十六进制id
	colorlog.PointPrint("You are running in " + config.DaemonServer.HardDiskMethod + " HardDisk method.")
	if _, err := os.Stat("/sys/fs/cgroup/cpu/server" + strconv.Itoa(server.ID)); err != nil {
		server.initCgroup()
	}
	colorlog.PointPrint("Preparing server runtime for ServerID:" + strconv.Itoa(server.ID))
	serverDataDir := "../servers/server" + strconv.Itoa(server.ID) // 在一开头就把serverDir算好，增加代码重用
	// 文件夹不存在则创建文件夹,并且将所有者设置为用户
	os.MkdirAll(serverDataDir+"/serverData", 0750)
	os.Chown(serverDataDir+"/serverData", config.DaemonServer.UserId, 0)
	if _, err0 := os.Stat(serverDataDir + ".loop"); err0 != nil && config.DaemonServer.HardDiskMethod != conf.HDM_LINK { //检查loop回环文件是否存在，如果不存在则创建
		colorlog.PointPrint("No loop file...")
		colorlog.LogPrint("Frozen Go Daemon will just make a new loop file")
		//  新增 loop
		if server.MaxHardDiskCapacity == 0 {
			server.MaxHardDiskCapacity = 10240
		}
		cmd := exec.Command("/bin/dd", "if=/dev/zero", "bs=1024", // MaxHardDisk单位kb
			"count="+strconv.Itoa(server.MaxHardDiskCapacity), "of=../servers/server"+strconv.Itoa(server.ID)+".loop")
		colorlog.LogPrint("Writing File with dd")
		output, err := cmd.CombinedOutput()
		if err != nil {
			colorlog.ErrorPrint(errors.New("Error with init cgroups:" + err.Error()))
			OutputErrReason(output)
			return errors.New("Error with dd output loop file." + err.Error())
		}
		colorlog.LogPrint("Done.")
		// 用mkfs格式化
		colorlog.LogPrint("Formatting loop file. Using mkfs.ext4")
		cmd2 := exec.Command("/sbin/mkfs.ext4", serverDataDir+".loop")
		colorlog.LogPrint("Done")
		output, err2 := cmd2.CombinedOutput()
		if err != nil {
			colorlog.ErrorPrint(errors.New("Error with init cgroups:" + err2.Error()))
			OutputErrReason(output)

			return errors.New("Error with mkfs.ext4:" + err2.Error())
		}

	}
	colorlog.LogPrint("Preparing server data dir.")
	// 为挂载文件夹做好准备
	//autoMakeDir(serverDataDir + "/execPath")
	//execPath,_ := filepath.Abs("../exec")
	//cmd2 := exec.Command("/bin/mount","-o","bind",execPath,serverDataDir + "/execPath")
	//cmd2.Run()
	// 挂载回环文件
	if config.DaemonServer.HardDiskMethod != conf.HDM_LINK {
		colorlog.LogPrint("The [" + config.DaemonServer.HardDiskMethod + "] method will limit the hardDisk space,so mounting loop file now.")
		colorlog.LogPrint("Mounting loop file")
		cmd3 := exec.Command("/bin/mount", "-o", "loop", serverDataDir+".loop", serverDataDir)
		//output3, err3 := cmd3.CombinedOutput()
		//if err3 != nil && strings.Index(string(output3), "is already mounted") <= 0 {
		//	colorlog.ErrorPrint(errors.New("Error with init cgroups:" + err3.Error()))
		//	OutputErrReason(output3)
		//	//return errors.New("Error with mounting loop file:"+err3.Error())
		//}
		AutoRunCmdAndOutputErr(cmd3,"initial cgroups")
	}
	networkArgs := []string{
		"../cgroup/cg.sh",
		"net",
		"add",
		networkID,
		strconv.Itoa(server.MaxUsingUpBandwidth),
		strconv.Itoa(server.MaxUnusedUpBandwidth),
		config.DaemonServer.NetworkCardName,
	}
	colorlog.LogPrint("Running command to add network: " + dumpCommand(networkArgs))
	cmdNetwork := exec.Command("/bin/bash", networkArgs...)
	if !AutoRunCmdAndOutputErr(cmdNetwork,"add network(tc)"){
		colorlog.PromptPrint("Server network bandwidth limit may invalid.")
	}

	colorlog.LogPrint("HardDisk Method is " + config.DaemonServer.HardDiskMethod)
	execConfig, err3 := server.loadExecutableConfig()
	if err3 != nil {
		colorlog.ErrorPrint(err3)
		return err3
	}

	switch config.DaemonServer.HardDiskMethod {

	case conf.HDM_MOUNT:
		for _,v := range execConfig.Link {
			colorlog.LogPrint("Mounting Dir : " + v)
			server.mountDir(v)
		}
		makeUserOwnedDir(serverDataDir + "/exec")
		cmd := exec.Command("mount","-o","bind","../exec",serverDataDir+"/exec")
		AutoRunCmdAndOutputErr(cmd,"mount dir file")
	case conf.HDM_LINK:
		err := server.linkDirs(execConfig)
		return err
	case conf.HDM_COPY:
		for _, v := range execConfig.Link {
			if _, dirExists := os.Stat(serverDataDir + "/" + v); dirExists != nil {
				cmd := exec.Command("cp", "-R", v, serverDataDir+"/"+v)
				if !AutoRunCmdAndOutputErr(cmd,"Copy files error"){
					return errors.New("Copy files error")
				}
			}

		}
		return nil
	}
	server.prepareVirtualEnv()
	return errors.New("Unexpected error")

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

func (server *ServerLocal) makeDirInServerPath(path string, mode os.FileMode) error {
	if filepath.IsAbs(path) {
		return os.MkdirAll(server.getLinkDir()+path, mode)
	} else {
		re, _ := filepath.Abs("../")
		current, _ := filepath.Abs(path)
		relPath, _ := filepath.Rel(re, current)
		return os.MkdirAll(server.getLinkDir()+"/"+relPath, mode)
	}
}

func (server *ServerLocal) getLinkDir() string {
	return "../servers/server" + strconv.Itoa(server.ID) + "/"
}

func (server *ServerLocal) linkDirFile(oldName string) error {
	if filepath.IsAbs(oldName) {
		return os.Link(oldName, "../servers/server"+strconv.Itoa(server.ID)+"/"+oldName)
	} else {
		re, _ := filepath.Abs("../")
		current, _ := filepath.Abs(oldName)
		relPath, _ := filepath.Rel(re, current)
		return os.Link(oldName, server.getLinkDir()+"/"+relPath)
	}
}
func (server *ServerLocal) linkDirs(conf ExecConf) error {
	links := append(conf.Link, "../exec")
	for _, v := range links {
		err := filepath.Walk(v, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return errors.New("Not a correct path in server execFile.")
			}
			if info.IsDir() {
				err := server.makeDirInServerPath(path, info.Mode())
				if err != nil {
					return err
				}
			} else {
				err := server.linkDirFile(path)
				if err != nil {
					return err
				}
			}
			return nil
		}))
		if err != nil {
			if err.Error() == "Not a correct path in server execFile." {
				return err
			} else {
				colorlog.WarningPrint(err.Error())
			}
		}
	}
	return nil
}
func (s *ServerLocal)mountDir(dirname string){
	target := "../servers/server" + strconv.Itoa(s.ID) + dirname
	if _,err := os.Stat(target);err != nil{
		if info,err := os.Stat(dirname);err == nil {
			os.MkdirAll(target, info.Mode())
		} else {
			return
		}
	}
	cmd := exec.Command("mount","-o","bind",dirname,target)
	AutoRunCmdAndOutputErr(cmd,"mount dirs")
}
