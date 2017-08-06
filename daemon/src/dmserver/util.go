package dmserver

import (
	"os/exec"
)

func autoMakeDir(name string) {
	cmd := exec.Command("mkdir", "-p", name)
	cmd.Run()
}
