package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/angelfluffyookami/HSVRUSB/modules/common/global"
)

type logEntry struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Level   string    `json:"level"`
}

type Log struct {
	entry logEntry
	alert bool
}

func (l *Log) Update() *Log {
	l.entry.Level = "UPDATE"
	return l
}

func (l *Log) Info() *Log {
	l.entry.Level = "INFO"
	return l
}

func (l *Log) Warn() *Log {
	l.entry.Level = "WARN"
	return l
}

func (l *Log) Err() *Log {
	l.entry.Level = "ERROR"
	return l
}

func (l *Log) Fatal() *Log {
	l.entry.Level = "FATAL"
	return l
}

func (l *Log) Message(m string) *Log {
	l.entry.Message = m
	return l
}

// If called, then sends log alert to discord channel.
func (l *Log) Alert() *Log {

	l.alert = true

	return l

}

// Adds log to logfile.json
func (l *Log) Add() {

	l.entry.Time = time.Now()

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

}

func (l *Log) sendAlert() {
	config := global.Config
	s := global.Session
	switch l.entry.Level {
	case "INFO":
		s.ChannelMessageSend(config.InfoChannel, l.entry.Level+": \n"+l.entry.Message+"\n"+fmt.Sprint(l.entry.Time))
	case "WARN":
		s.ChannelMessageSend(config.WarnChannel, l.entry.Level+": \n"+l.entry.Message+"\n"+fmt.Sprint(l.entry.Time))
	case "ERR":
		s.ChannelMessageSend(config.ErrChannel, l.entry.Level+": \n"+l.entry.Message+"\n"+fmt.Sprint(l.entry.Time))
	case "FATAL":
		s.ChannelMessageSend(config.ErrChannel, l.entry.Level+": \n"+l.entry.Message+"\n"+fmt.Sprint(l.entry.Time))
	case "UPDATE":
		s.ChannelMessageSend(config.UpdateChannel, l.entry.Level+": \n"+l.entry.Message+"\n"+fmt.Sprint(l.entry.Time))
	}
}
