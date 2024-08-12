package shutdown

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Graceful(closers ...io.Closer) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM)

	signal := <-sig

	log.Printf("received signal %s, shutting down", signal.String())
	for _, c := range closers {
		if err := c.Close(); err != nil {
			log.Printf("error closing: %v", err)
		}
	}
	os.Exit(0)
}
