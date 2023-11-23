//go:build windows

package win32

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/angelfluffyookami/HSVRUSB/modules/common/global"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log

type bvr struct{}

func (m *bvr) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {

		c := <-r
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
			// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
			time.Sleep(100 * time.Millisecond)
			changes <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			// golang.org/x/sys/windows/svc.TestExample is verifying this output.
			testOutput := strings.Join(args, "-")
			testOutput += fmt.Sprintf("-%d", c.Context)
			elog.Info(1, testOutput)
			break loop
		case svc.Pause:
			global.PauseThreads <- true
			changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			elog.Info(2, "Pausing operations.")
		case svc.Continue:
			global.PauseThreads <- false
			changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			elog.Info(3, "Resuming operations.")

		default:
			elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
		}

	}
	changes <- svc.Status{State: svc.StopPending}
	defer handleExit()
	defer fmt.Println("cleaning...")
	panic(Exit{3})

}

type Exit struct{ Code int }

// exit code handler
func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok {
			os.Exit(exit.Code)
		}
		panic(e) // not an Exit, bubble up
	}
}

func runService(name string, isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(name, &bvr{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}
