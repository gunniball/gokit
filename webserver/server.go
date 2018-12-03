package webserver

import (
	"flag"
	"fmt"
	"github.com/finalist736/gokit/logger"
	"github.com/gocraft/web"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

var file1 *os.File = nil
var listener1 net.Listener
var listen_file_descriptor *int = flag.Int("fd", 0, "Server socket descriptor")

func Start(router *web.Router, port string) {
	var err error
	serv := &http.Server{Addr: fmt.Sprintf("%s", port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serv.SetKeepAlivesEnabled(false)

	logger.StdOut().Infof("Listen on %s port", serv.Addr)
	if *listen_file_descriptor != 0 {
		logger.StdOut().Info("Starting with FD ", *listen_file_descriptor)
		file1 = os.NewFile(uintptr(*listen_file_descriptor), "/tmp/socket-"+strconv.FormatInt(rand.Int63(), 16))
		listener1, err = net.FileListener(file1)
		if err != nil {
			panic(fmt.Sprintf("fd listener failed: %s", err))
		}
	} else {
		logger.StdOut().Info("Virgin Start")
		listener1, err = net.Listen("tcp", serv.Addr)
		if err != nil {
			panic(fmt.Sprintf("listener failed: %s", err))
		}
	}
	go serv.Serve(listener1)
}

func Stop() {
	if listener1 != nil {
		listener1.Close()
	}

	if file1 != nil {
		file1.Close()
	}
}

func Grace(configPath *string) {
	if configPath == nil {
		panic("config can not be empty")
	}
	listener2 := listener1.(*net.TCPListener)
	file2, err := listener2.File()
	if err != nil {
		logger.StdErr().Warn("get file2 from listener", err)
	}
	fd1 := int(file2.Fd())
	fd2, err := syscall.Dup(fd1)
	if err != nil {
		logger.StdErr().Warn("Dup error:", err)
	}
	listener1.Close()
	if file1 != nil {
		file1.Close()
	}
	commandLine := fmt.Sprintf("%s", os.Args[0])
	fdParam := fmt.Sprintf("%d", fd2)
	cmd := exec.Command(commandLine, "-config", *configPath, "-fd", fdParam)
	err = cmd.Start()
	if err != nil {
		panic("sub process run error: " + err.Error())
	}
}
