/*
本包提供了一些关于FrozenGo文件传输协议的实现
FrozenGo的文件传输协议采用的是ws分块上传模式，将用户面板的FTP简洁化地放到了网页上，以及
实现了ws slice 上传算法，让文件的上传变得更加灵活，可控，但有所不足的是，进度条功能由于Axoford12只学了半个小时的js
实在是无力实现。
*/
package ftrans

import (
	"conf"
	"net/http"
)

func Start(configure conf.Config) {
	config = configure
	// 配置conf

	// 为http加上处理器
	http.HandleFunc("/download", handleDownload)
	http.HandleFunc("/upload", handleUpload)
	// TODO Handler : Delete
}
