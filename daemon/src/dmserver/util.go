package dmserver

import (
	"os/exec"
	"fmt"
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
func (s *ServerRun)testChannel(){
	for {
		fmt.Printf("%s",<-s.ToOutput.To)
	}
}