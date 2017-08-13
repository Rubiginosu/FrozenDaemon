package dmserver

import (
	"os"
	"strings"
	"fmt"
	"path/filepath"
)
/**
构建一个Infos
 */
func buildServerInfos(infos []os.FileInfo) []ServerPathFileInfo{
	spfInfos := make([]ServerPathFileInfo,len(infos))
	for k,v := range infos {
		spfInfos[k] = __buildServerInfo(v)
	}
	return spfInfos
}
// 上面函数的辅助方法
func __buildServerInfo(info os.FileInfo) ServerPathFileInfo{
	return ServerPathFileInfo{
		Name:info.Name(),
		Mode:strings.Replace(fmt.Sprintf("%32b",uint32(info.Mode()))," ","0",-1),
		ModTime:info.ModTime().Unix(),
	}
}
/**
检查目录是否合法，避免出现../../../../etc/passwd之类似小天才行为
 */
func validateOperateDir(upload string,path string) bool{
	return strings.Index(filepath.Clean(upload + path),upload)  >= 0
}