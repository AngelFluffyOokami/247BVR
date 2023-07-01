package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/angelfluffyookami/247BVR/modules/common"
	"github.com/angelfluffyookami/247BVR/modules/common/global"
)

type logEntry struct {
	time    time.Time
	message string
	level   string
}

type Log struct {
	entry logEntry
	alert bool
	panic bool
}

func (l *Log) Update() *Log {
	l.entry.level = "UPDATE"
	return l
}

func (l *Log) Info() *Log {
	l.entry.level = "INFO"
	return l
}

func (l *Log) Warn() *Log {
	l.entry.level = "WARN"
	return l
}

func (l *Log) Err() *Log {
	l.entry.level = "ERROR"
	return l
}

func (l *Log) Fatal() *Log {
	l.entry.level = "FATAL"
	return l
}

func (l *Log) Panic() *Log {
	l.panic = true
	return l
}

func (l *Log) Message(m string) *Log {
	l.entry.message = m
	return l
}

// If called, then sends log alert to discord channel.
func (l *Log) Alert() *Log {

	l.alert = true

	return l

}

// Adds log to logfile.json
func (l *Log) Add() {

	// open the log file
	logFile, err := os.OpenFile("logfile.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// marshal the log entry
	entryJSON, err := json.MarshalIndent(l.entry, "\n", "\n")
	if err != nil {
		panic(err)
	}

	// write the log entry to the file
	if _, err := logFile.Write(entryJSON); err != nil {
		panic(err)
	}
	if err := logFile.Sync(); err != nil {
		panic(err)
	}
	if l.alert {
		l.sendAlert()
	}
	if l.panic {
		log.Panic(l.entry.message)
	}

}

func (l *Log) sendAlert() {
	config := common.Config
	s := global.Session
	switch l.entry.level {
	case "INFO":
		s.ChannelMessageSend(config.InfoChannel, l.entry.level+": \n"+l.entry.message+"\n"+fmt.Sprint(l.entry.time))
	case "WARN":
		s.ChannelMessageSend(config.WarnChannel, l.entry.level+": \n"+l.entry.message+"\n"+fmt.Sprint(l.entry.time))
	case "ERR":
		s.ChannelMessageSend(config.ErrChannel, l.entry.level+": \n"+l.entry.message+"\n"+fmt.Sprint(l.entry.time))
	case "FATAL":
		s.ChannelMessageSend(config.ErrChannel, l.entry.level+": \n"+l.entry.message+"\n"+fmt.Sprint(l.entry.time))
	case "UPDATE":
		s.ChannelMessageSend(config.UpdateChannel, l.entry.level+": \n"+l.entry.message+"\n"+fmt.Sprint(l.entry.time))
	}
}
