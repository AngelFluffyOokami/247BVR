//go:build windows

package win32

import (
	"log"
	"os"

	"github.com/angelfluffyookami/HSVRUSB/modules/common/global"
	"golang.org/x/sys/windows/svc"
)

var svcName = "HSVRUSB"

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
		if err != nil {
			go func() {
				TextLog <- err.Error()
			}()
		}
		return
	case "install":
		TextLog <- "Creating user..."
		ok := AddUser()
		if !ok {
			TextLog <- "Unable to create user"
			os.Exit(1)
		}
		TextLog <- "User set up"
		TextLog <- "Preparing install destination..."

		ok = populatePaths()
		if !ok {

			TextLog <- "Could not prepare install destination"
			Installed <- false

		}
		TextLog <- "Beginning service install..."
		err = installService(svcName, "HSVR USB 2.0 service")
		if err != nil {

			TextLog <- err.Error()
			TextLog <- "Quitting installer"
			Installed <- false
			return
		} else {
			Installed <- true
		}
	case "remove":
		err = removeService(svcName)
		if err != nil {
			go func() {
				TextLog <- err.Error()
			}()
		}
	case "start":
		err = startService(svcName)
		if err != nil {
			go func() {
				TextLog <- err.Error()
			}()
		}
	case "stop":
		err = controlService(svcName, svc.Stop, svc.Stopped)
		if err != nil {
			go func() {
				TextLog <- err.Error()
			}()
		}
	case "pause":
		err = controlService(svcName, svc.Pause, svc.Paused)
		if err != nil {
			go func() {
				TextLog <- err.Error()
			}()
		}
		svcCtrl(true)
	case "continue":
		err = controlService(svcName, svc.Continue, svc.Running)
		if err != nil {
			go func() {
				TextLog <- err.Error()
			}()
		}
		svcCtrl(false)
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
	}
}

var Installed = make(chan bool)
var TextLog = make(chan string)

func svcCtrl(stop bool) {
	global.PauseThreads <- stop
}
