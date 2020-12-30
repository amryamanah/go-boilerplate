package exithandler

import (
	"github.com/amryamanah/go-boilerplate/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func Init(cb func()) {
	sigs := make(chan os.Signal, 1)
	terminate := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logger.Info.Println("exit reason: ", sig)
		terminate <- true
	}()

	<-terminate
	cb()
	logger.Info.Println("exiting program")
}
