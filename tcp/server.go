package tcp

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/lastIndexOf/mini_redis/interface/tcp"
	"github.com/lastIndexOf/mini_redis/lib/logger"
)

type Config struct {
	Addr string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	listener, err := net.Listen("tcp", cfg.Addr)
	closeChan := make(chan struct{})
	signalChan := make(chan os.Signal, 2)

	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-signalChan

		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM:
			closeChan <- struct{}{}
		}
	}()

	if err != nil {
		return err
	}

	logger.Info("Listening on " + cfg.Addr)

	ListenAndServe(listener, handler, closeChan)

	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	defer func() {
		listener.Close()
		handler.Close()
	}()

	go func() {
		<-closeChan

		logger.Info("Shutting down the server...")

		listener.Close()
		handler.Close()
	}()

	ctx := context.Background()
	var waitAll sync.WaitGroup

	for {
		conn, err := listener.Accept()

		if err != nil {
			break
		}

		logger.Info("Accept connection from " + conn.RemoteAddr().String())

		waitAll.Add(1)

		go func() {
			defer waitAll.Done()
			handler.Handle(ctx, conn)
		}()
	}

	waitAll.Wait()
}
