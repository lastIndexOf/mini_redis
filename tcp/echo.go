package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/lastIndexOf/mini_redis/lib/logger"
	"github.com/lastIndexOf/mini_redis/lib/sync/atomic"
	"github.com/lastIndexOf/mini_redis/lib/sync/wait"
)

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (client *EchoClient) Close() (_ error) {
	client.Waiting.WaitWithTimeout(time.Second * 10)
	client.Conn.Close()
	return nil
}

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if handler.closing.Get() {
		conn.Close()
		return
	}

	client := &EchoClient{Conn: conn}
	handler.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString(byte('\n'))

		if err != nil {
			if err == io.EOF {
				logger.Info("Connection closed by client")
				handler.activeConn.Delete(client)
			} else {
				logger.Warn("Error reading from connection: ", err)
			}
			return
		}

		client.Waiting.Add(1)
		conn.Write([]byte(msg))
		client.Waiting.Done()
	}
}

func (handler *EchoHandler) Close() (_ error) {
	logger.Info("Closing echo handler")

	handler.closing.Set(true)
	handler.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		client.Conn.Close()
		return true
	})

	return nil
}
