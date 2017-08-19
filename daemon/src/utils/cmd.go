package utils

import (
	"os/exec"
	"os"
	"colorlog"
	"fmt"
)

func AutoRunCmdAndOutputErr(cmd *exec.Cmd,errorAt string) bool{
	cmd.Env = os.Environ()
	out,err := cmd.CombinedOutput()
	if err != nil {
		colorlog.ErrorPrint(errorAt,err)
		OutputErrReason(out)
		return false
	}
	return true
}
func OutputErrReason(output []byte) {
	colorlog.LogPrint("Reason:")
	fmt.Println(colorlog.ColorSprint("-----ERROR_MESSAGE-----", colorlog.FR_RED))
	fmt.Println(colorlog.ColorSprint(string(output), colorlog.FR_RED))
	fmt.Println(colorlog.ColorSprint("-----ERROR_MESSAGE-----", colorlog.FR_RED))
}