package filetrans

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func getMessage(c net.Conn) string {
	bufc := bufio.NewReader(c)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			c.Close()
			break
		}
		res := strings.Trim(strings.TrimRight(line, "\r"), "\n")
		fmt.Printf("Receved: %s\n", res)
		return res
	}
	return ""
}

func sendMessage(c net.Conn, message string) bool {
	_, err := io.WriteString(c, message)
	return err == nil
}

func parseCommandArg(data []byte) *Command {
	// 前四位设置为Command ,第六位到最后是Arg
	//   AAAA      BBBBBBBBBB
	// Command       Arg
	if len(data) < 6 {
		return &Command{"", ""}
	}
	return &Command{string(data[:4]), string(data[5:])}
}

func parseFileInfoToLocalFile(f os.FileInfo) localServerFile {
	return localServerFile{
		f.Name(),
		f.Mode().String(),
		f.IsDir(),
		f.Size(),
		f.ModTime().Unix(),
	}
}
