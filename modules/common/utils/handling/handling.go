package handling

import (
	"fmt"
	"runtime"

	"github.com/angelfluffyookami/247BVR/modules/common/global"
	"github.com/angelfluffyookami/247BVR/modules/common/utils/logger"
)

// Defer at start of functions for panic "handling".
func RecoverPanic(channelID string) {

	if r := recover(); r != nil {

		s := global.Session

		// get the stack trace of the panic
		tempbuf := make([]byte, 10000)
		buflength := runtime.Stack(tempbuf, false)
		var buf []byte
		if buflength >= 1900 {
			buf = make([]byte, 1900)
		} else {
			buf = make([]byte, buflength)
		}
		runtime.Stack(buf, false)

		log.Fatal().Message(fmt.Sprintf("Recovering from panic: %v\n Stack trace: %s", r, buf)).Alert().Add()
		if channelID != "" {
			s.ChannelMessageSend(channelID, "Error processing command.\nBug report sent to developers.")
		}

	}

}

var log = logger.Log{}
