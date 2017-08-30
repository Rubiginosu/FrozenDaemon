package dmserver

import (
	"archive/zip"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func install(execConfig ExecInstallConfig) {
	modulesData, err := ioutil.ReadFile(config.ServerManager.Modules)
	if err != nil {
		fmt.Println(err)
		return
	}

	toBeInstall := needInstallModules(string(modulesData), execConfig.Rely)
	for i := 0; i < len(toBeInstall); i++ {
		fmt.Println("Attemp to install module:" + toBeInstall[i].Name)
		toBeInstall[i].install()
	}

	execConfig.downloadExecAndConf()
	fmt.Println("Congratulations!Your module:" + execConfig.StartConf.Name + " has been successfully installed!")
}

func needInstallModules(installed string, rely []Module) []Module {
	fmt.Print("Finding modules to be installed...")
	var need []Module
	for i := 0; i < len(rely); i++ {
		if !isRelyInModules(installed, rely[i]) {
			need = append(need, rely[i])
		}
	}
	fmt.Println("Done")
	return need
}
func isRelyInModules(installed string, rely Module) bool {
	return strings.Index(installed, rely.Name) >= 0
}
func (m *Module) install() error {
	fmt.Println("Installing module " + m.Name)
	conn, err := http.Get(m.Download)
	if err != nil {
		return err
	}

	file, err1 := os.Create("../exec/~temp.zip")
	b, _ := ioutil.ReadAll(conn.Body)
	fmt.Println("Download " + m.Download + " OK.sud")

	file.Write(b)
	file.Close()
	conn.Body.Close()

	isMatch, err2 := md5Check("../exec/~temp.zip", m.Md5)
	if err2 != nil {
		fmt.Println("Error with MD5 check:" + err2.Error())
		return err2
	} else if !isMatch {
		fmt.Println("MD5 mismatch")
		return errors.New("MD5 mismatch.")
	} else {
		fmt.Println("MD5 check done.")
	}

	r, err := zip.OpenReader("../exec/~temp.zip")
	dir := "../exec/" + m.Name + "/"
	fmt.Println("Extracting archive")
	for _, file := range r.File {
		rc, err := file.Open()
		if err != nil {
			fmt.Println(err)
		}
		os.MkdirAll(filepath.Dir(dir+file.Name), 755)
		f, _ := os.Create(dir + file.Name)
		io.Copy(f, rc)
		f.Close()
	}
	fmt.Println("Extract done.")
	fmt.Print("Removing temp file....")
	cmd := exec.Command("rm", "-rf", "../exec/~temp.zip")
	cmd.Run()
	fmt.Println("OK")
	fmt.Print("Changing file mode...")
	cmd2 := exec.Command("chmod", strings.Split(m.Chmod, ",")[0], "-R", "../exec/"+m.Name)
	cmd3 := exec.Command("chown", strings.Split(m.Chmod, ",")[1], "-R", "../exec/"+m.Name)
	cmd2.Run()
	cmd3.Run()
	fmt.Println("OK")
	if err1 != nil {
		return err1
	}
	modulesFile, _ := os.OpenFile(config.ServerManager.Modules, os.O_APPEND|os.O_WRONLY, 0777)
	modulesFile.Write([]byte(m.Name + ","))
	modulesFile.Close()
	return nil
}

func (e *ExecInstallConfig) downloadExecAndConf() {
	fmt.Println("Downloading file and conf...")
	conn, err := http.Get(e.Url)
	if err != nil {
		fmt.Println("Download form " + e.Url + "error!")
		return
	}
	elements := strings.Split(e.Url, "/")
	file, err1 := os.Create("../exec/" + elements[len(elements)-1])
	if err1 != nil {
		fmt.Println("Create exec file error!")
		return
	}
	io.Copy(file, conn.Body)
	file.Close()
	res, err3 := md5Check("../exec/"+elements[len(elements)-1], e.Md5)
	if err3 != nil {
		fmt.Println("Error compute MD5:" + err3.Error())
	}
	if !res {
		fmt.Println("MD5 check failed")
		return
	} else {
		fmt.Println("MD5 check ok")
	}
	conn.Body.Close()
	execConfFile, err2 := os.Create("../exec/" + e.StartConf.Name + ".json")
	if err2 != nil {
		fmt.Println("Create exec config file error!")
	}
	b, _ := json.Marshal(e.StartConf)
	execConfFile.Write(b)
	execConfFile.Close()
	fmt.Println("Done")
}
func md5Check(name string, sum string) (bool, error) {
	fmt.Println("Checking MD5 sum " + name)
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return false, err
	}
	md5bytes := md5.Sum(data)
	fmt.Println("MD5:" + fmt.Sprintf("%x", md5bytes))
	return fmt.Sprintf("%x", md5bytes) == sum, nil
}
