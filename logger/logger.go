package logger

import (
	"encoding/json"
	"logger_test/core/config"
	"os"
	"io"
	"github.com/sirupsen/logrus"
	"syscall"
)

var stdout, stderr logrus.FieldLogger
var socketLogger *SocketLogger

func JsonStdOut(name string, jsn interface{}) {
	ba, _ := json.MarshalIndent(jsn, "", " ")
	StdOut().Debug(name+" %v", ba)
}

func StdOut() logrus.FieldLogger {
	return stdout
}

func StdErr() logrus.FieldLogger {
	return stderr
}

func ReloadLogs() {
	logLevel := config.MustString("loglevel")
	logType := config.MustString("logtype")
	switch logType {
	case "std":
		stdout = NewLogger("out", logType, logLevel)
		stderr = NewLogger("err", logType, logLevel)
		// selected std outs
	case "file":
		stdout = NewLogger("out", logType, logLevel)
		stderr = NewLogger("err", logType, logLevel)
		logPath := config.MustString("logpath")
		if logPath == "" {
			panic("logger: no log path specified")
		}
		err := os.MkdirAll(logPath, 0764)
		if err != nil {
			panic(err)
		}

		stdoutfile, err := os.OpenFile(logPath+"/stdout.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0764)
		if err != nil {
			panic(err)
		}
		_, err = stdoutfile.Seek(0, io.SeekEnd)
		if err != nil {
			panic(err)
		}
		err = syscall.Dup2(int(stdoutfile.Fd()), int(os.Stdout.Fd()))
		if err != nil {
			panic(err)
		}

		stderrfile, err := os.OpenFile(logPath+"/stderr.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0764)
		if err != nil {
			panic(err)
		}
		_, err = stderrfile.Seek(0, io.SeekEnd)
		if err != nil {
			panic(err)
		}
		err = syscall.Dup2(int(stderrfile.Fd()), int(os.Stderr.Fd()))
		if err != nil {
			panic(err)
		}
	case "socket":
		logPath := config.MustString("logpath")
		if logPath == "" {
			panic("logger: no log socket specified")
		}
		stdout = NewLogger("out", logType, logLevel)
		stderr = NewLogger("err", logType, logLevel)
		socketLogger = NewSocketLogger()
	default:
		panic("logger: incorrect logging type")
	}
}
