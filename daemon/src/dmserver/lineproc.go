package dmserver

import (
	"colorlog"
	"regexp"
)

/*
本文件用于处理一行蜜汁代码
*/
func (s *ServerRun) processOutputLine(line string /*,startReg *regexp.Regexp,joinReg *regexp.Regexp,leftReg *regexp.Regexp*/) {
	startReg := regexp.MustCompile("Done \\(.+s\\)!")
	joinReg := regexp.MustCompile("(\\w+)\\[.+\\] logged in")
	leftReg := regexp.MustCompile("(\\w+) left the game.")
	if startReg.MatchString(line) {
		colorlog.PointPrint("Server Started!")
		if server, ok := serverSaved[s.ID]; ok {
			server.Status = SERVER_STATUS_RUNNING
		}
	} else if sub := joinReg.FindStringSubmatch(line); len(sub) > 0 {
		// 加入服务器
		s.Players = append(s.Players, sub[1])
	} else if sub := leftReg.FindStringSubmatch(line); len(sub) > 0 {
		// 退出服务器
		if index := s.findPlayByName(sub[1]); index >= 0 {
			s.Players = append(s.Players[:index], s.Players[index+1:]...)
		}

	}

}
func (s *ServerRun) findPlayByName(name string) int {
	for i := 0; i < len(s.Players); i++ {
		if s.Players[i] == name {
			return i
		}
	}
	return -1
}

func (s *ServerRun) inputLine(line string) error {
	_, err := (*s.StdinPipe).Write([]byte(line))
	return err
}
