//go:build windows

package win32

import (
	"syscall"
)

var (
	beepFunc = syscall.MustLoadDLL("user32.dll").MustFindProc("MessageBeep")
)

func Beep() {
	beepFunc.Call(0xffffffff)
}
