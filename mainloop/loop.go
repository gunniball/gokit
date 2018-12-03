package mainloop

import (
	"github.com/finalist736/gokit/config"
	"github.com/finalist736/gokit/logger"
	"os"
	"os/signal"
	"syscall"
)

type StopFunction func()
type GraceFunction func()

/*
os.Interrupt, os.Kill, syscall.SIGTERM
for break loop

syscall.SIGUSR1
for logs reload

syscall.SIGUSR2
for grace restarting

*/
func Loop(stopfunc StopFunction, graceFunc GraceFunction, c *config.Config) {
	usr1sigChannel := make(chan os.Signal)
	usr2sigChannel := make(chan os.Signal)
	interruptChannel := make(chan os.Signal)
	signal.Notify(usr1sigChannel, syscall.SIGUSR1)
	signal.Notify(usr2sigChannel, syscall.SIGUSR2)
	signal.Notify(interruptChannel, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case killSignal := <-interruptChannel:
			logger.StdOut().Debug("Got signal:", killSignal)
			if stopfunc != nil {
				stopfunc()
			}
			if killSignal == os.Interrupt {
				logger.StdOut().Warn("Daemon was interrupted by system signal")
				return
			}
			logger.StdOut().Debug("Daemon was killed")
			return
		case <-usr1sigChannel:
			logger.ReloadLogs(c)
		case <-usr2sigChannel:
			if graceFunc != nil {
				graceFunc()
			}
			return
			// grace restarting...
		}
	}
}
