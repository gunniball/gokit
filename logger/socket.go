package logger

import (
	"github.com/sirupsen/logrus"
	"time"
	"fmt"
	"os"
	"net"
	"github.com/finalist736/gokit/config"
)


type SocketLogger struct {
	net.Conn
	sendchannel chan *Logbody
	timeout     int
}


type (
	Logger struct {
		Type 	string
		Level	logrus.Level
	}
	Loghead struct {
		Name	string	`json:"name"`
		Version	string	`json:"version"`
		Id	int	`json:"id"`
	}
	Logbody struct {
		Type	string	`json:"type"`
		Logs	string	`json:"logs"`
	}
)


func NewSocketLogger() *SocketLogger {
	s := &SocketLogger{sendchannel: make(chan*Logbody, 4096 * 3)}
	go s.Listen()
	return s
}


func NewLogger (stdtype, logtype, logLevel string) logrus.FieldLogger {
	var err error
	if logtype == "socket" {
		lg := &Logger{}
		lg.Type = stdtype
		lg.Level, err = logrus.ParseLevel(logLevel)
		if err != nil {
			lg.Level = logrus.DebugLevel
		}
		return lg
	} else {
		lg := logrus.New()
		lg.Formatter = &logrus.TextFormatter{
			ForceColors:false,
			TimestampFormat:time.RFC3339Nano,
			DisableTimestamp:false,
			FullTimestamp:true,
		}
		lg.Level, err = logrus.ParseLevel(logLevel)
		if err != nil {
			lg.Level = logrus.InfoLevel
		}
		if stdtype == "out" {
			lg.Out = os.Stdout
		}
		return lg
	}
}

func (s *SocketLogger) Connect (path string) error {
	var err error
	s.Conn, err = net.Dial("tcp", path)
	if err != nil {
		return err
	}
	return nil
}

func (s *SocketLogger) Auth() error {
	auth := &Loghead{Name: config.MustString("server_name"), Id: config.MustInt("server_id"), Version: config.MustString("server_version")}
	return s.Send(auth)
}


func (s *SocketLogger) Listen () {
	for {
		select {
		case log := <-s.sendchannel:
			s.SendLog(log)
		}
	}
}

func (s* Logger) toChan (log string) {
	data :=  &Logbody{ Type: s.Type, Logs: log }
	socketLogger.sendchannel <- data
}

func (s *SocketLogger) SendLog (log *Logbody) error {
	var err error
	for {
		err = s.Send(log)
		if err != nil {
			for {
				time.Sleep(time.Millisecond * time.Duration(s.timeout))
				err = s.Connect(config.MustString("logpath"))
				if err == nil {
					err = s.Auth()
					if err == nil {
						s.timeout = 0
						break
					}
				}
				s.timeout += 10
			}
		} else {
			break
		}
	}
	return nil
}


func (logger *Logger) WithField(key string, value interface{}) (l *logrus.Entry) {
	return
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *Logger) WithFields(fields logrus.Fields) (l *logrus.Entry) {
	return
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (logger *Logger) WithError(err error) (l *logrus.Entry) {
	return
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.DebugLevel, fmt.Sprintf(format, args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.InfoLevel, fmt.Sprintf(format, args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Printf(format string, args ...interface{}) {
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.WarnLevel, fmt.Sprintf(format, args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.WarnLevel, fmt.Sprintf(format, args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.ErrorLevel, fmt.Sprintf(format, args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.FatalLevel, fmt.Sprintf(format, args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.PanicLevel, fmt.Sprintf(format, args...))
	logger.toChan(log)
}

func (logger *Logger) Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.DebugLevel, fmt.Sprint(args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.InfoLevel, fmt.Sprint(args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Print(args ...interface{}) {
}

func (logger *Logger) Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.WarnLevel, fmt.Sprint(args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Warning(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.WarnLevel, fmt.Sprint(args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.ErrorLevel, fmt.Sprint(args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.FatalLevel, fmt.Sprint(args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Panic(args ...interface{}) {
	log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.PanicLevel, fmt.Sprint(args...))
	logger.toChan(log)
}

func (logger *Logger) Debugln(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.DebugLevel, fmt.Sprint(args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Infoln(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.InfoLevel, fmt.Sprint(args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Println(args ...interface{}) {
}

func (logger *Logger) Warnln(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.WarnLevel, fmt.Sprint(args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Warningln(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.WarnLevel, fmt.Sprint(args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Errorln(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.ErrorLevel, fmt.Sprint(args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Fatalln(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.FatalLevel, fmt.Sprint(args...))
		logger.toChan(log)
	}
}

func (logger *Logger) Panicln(args ...interface{}) {
	log := fmt.Sprintf("time=\"%s\" level=%s msg=%s\n", time.Now().Format(time.RFC3339), logrus.PanicLevel, fmt.Sprint(args...))
	logger.toChan(log)
}
