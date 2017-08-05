package dmserver

import (
	"regexp"
	"colorlog"
)

/*
本文件用于处理一行蜜汁代码
 */
func (s *ServerRun)processOutputLine(line string){
	startReg := regexp.MustCompile("Done \\(.+s\\)!")
	if startReg.MatchString(line) {
		colorlog.PointPrint("Server Started!")
		if index := searchServerByID(s.ID);index >= 0{
			serverSaved[index].Status = SERVER_STATUS_RUNNING
		}
	}

}

func (s *ServerRun)inputLine(line string) error{
	_,err :=(*s.StdinPipe).Write([]byte(line))
	return err
}
