package colorlog

import (
	"fmt"
	"strconv"
)

// 让终端有颜色
const (
	FR_RED = 31 + iota
	FR_GREEN
	FR_YELLOW
	FR_BLUE
	FR_PURPLE
	FR_CYAN
)

/**
返回一个字符串，包含以color为颜色的Linux字串
*/
func ColorSprint(message string, color int) string {
	return "\033[1;" + strconv.Itoa(color) + "m" + message + "\033[0m"
}

// 错误
func ErrorPrint(err error) {
	fmt.Println(ColorSprint("[Error]"+err.Error(), FR_RED))
}

// 信息
func LogPrint(message string) {
	fmt.Print(ColorSprint("[Info]", FR_BLUE))
	fmt.Println(message)
}

// 警告
func WarningPrint(message string) {
	fmt.Println(ColorSprint("[Warning]", FR_YELLOW), message)
}

// 操作点打印
func PointPrint(message string) {
	fmt.Print(ColorSprint("[Point]", FR_PURPLE))
	fmt.Println(message)
}

// 提示信息打印
func PromptPrint(message string) {
	fmt.Print(ColorSprint("[Prompt]", FR_GREEN))
	fmt.Println(message)
}
