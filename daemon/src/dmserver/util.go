package dmserver

import (
	"os/exec"
)

func IsServerAvaible(serverid int) bool {
	for i := 0; i < len(serverSaved); i++ {
		if serverSaved[i].ID == serverid {
			return true
		}
	}
	return false
}

func autoMakeDir(name string) {
	cmd := exec.Command("mkdir", "-p", name)
	cmd.Run()
}
