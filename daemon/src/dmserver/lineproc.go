package dmserver

import "bufio"

/*
本文件用于处理一行蜜汁代码
 */
func (s *ServerRun)processOutputLine(line string){

}

func (s *ServerRun)inputLine(line string){
	writer := bufio.NewWriter(*s.StdinPipe)
	writer.WriteString(line)
}
