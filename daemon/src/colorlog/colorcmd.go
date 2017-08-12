package colorlog

import (
	"fmt"
	"strconv"
)

// 让终端有颜色
const (
	FR_WHITE = 30 + iota
	FR_RED
	FR_GREEN
	FR_YELLO
	FR_BLUE
	FR_PURPLE
	FR_CYAN
)
const (
	BK_WHITE = 40 + iota
	BK_RED
	BK_GREEN
	BK_YELLO
	BK_BLUE
	BK_PURPLE
	BK_CYAN
)

func ColorSprint(message string, color int) string {
	return "\033[1;" + strconv.Itoa(color) + "m" + message + "\033[0m"
}

func ErrorPrint(err error) {
	fmt.Println(ColorSprint("[Error]"+err.Error(), FR_RED))
}
func LogPrint(message string) {
	fmt.Print(ColorSprint("[Info]", FR_BLUE))
	fmt.Println(message)
}

func WarningPrint(message string) {
	fmt.Println(ColorSprint("[Warning]", FR_YELLO),message)
}
func PointPrint(message string) {
	fmt.Print(ColorSprint("[Point]", FR_PURPLE))
	fmt.Println(message)
}
func PromptPrint(message string){
	fmt.Print(ColorSprint("[Prompt]",FR_GREEN))
	fmt.Println(message)
}
