//go:build windows

package win32

import (
	"log"

	"golang.org/x/sys/windows/svc"
)

var svcName = "247bvr"

var Install = "install"
var Debug = "debug"
var Remove = "remove"
var Start = "start"
var Stop = "stop"
var Pause = "pause"
var Continue = "continue"

func Service(cmd string) {

	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service: %v", err)
	}
	if inService {
		runService(svcName, false)
		return
	}

	switch cmd {
	case "debug":
		runService(svcName, true)
		return
	case "install":
		err = installService(svcName, "HSVR ELO Statistics Bot service.")
	case "remove":
		err = removeService(svcName)
	case "start":
		err = startService(svcName)
	case "stop":
		err = controlService(svcName, svc.Stop, svc.Stopped)
	case "pause":
		err = controlService(svcName, svc.Pause, svc.Paused)
	case "continue":
		err = controlService(svcName, svc.Continue, svc.Running)
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
	}
}

func ContinueSvc() {

}
