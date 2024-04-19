package handler

import (
	"context"
	"errors"
	"github.com/lastIndexOf/mini_redis/database"
	db "github.com/lastIndexOf/mini_redis/interface/database"
	"github.com/lastIndexOf/mini_redis/lib/logger"
	"github.com/lastIndexOf/mini_redis/lib/sync/atomic"
	"github.com/lastIndexOf/mini_redis/resp/connection"
	"github.com/lastIndexOf/mini_redis/resp/parser"
	"github.com/lastIndexOf/mini_redis/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-Err unknown\r\n")
)

type RespHandler struct {
	activeConn sync.Map
	db         db.Database
	closing    atomic.Boolean
}

func MakeRespHandler(db db.Database) *RespHandler {
	return &RespHandler{
		db: db,
	}
}

func MakeEchoHandler() *RespHandler {
	return &RespHandler{
		db: database.NewEchoDatabase(),
	}
}

func (handler *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	handler.db.AfterClientClose(client)
	handler.activeConn.Delete(client)
}

func (handler *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if handler.closing.Get() {
		conn.Close()
		return
	}

	client := connection.NewConn(conn)
	handler.activeConn.Store(client, struct{}{})

	ch := parser.ParseStream(conn)

	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				errors.Is(payload.Err, io.ErrUnexpectedEOF) ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				handler.closeClient(client)
				logger.Info("Connection closed by client" + client.RemoteAddr().String())
				return
			}

			errReply := reply.MakeStandardErrReply(payload.Err.Error())
			err := client.Write(errReply.Bytes())

			if err != nil {
				handler.closeClient(client)
				logger.Info("Connection closed by client" + client.RemoteAddr().String())
				return
			}

			continue
		}

		//	exec db command
		if payload.Data == nil {
			continue
		}

		params, ok := payload.Data.(*reply.MultiBulkReply)

		if !ok {
			logger.Error("Required multi bulk reply type")
			continue
		}

		reply := handler.db.Exec(client, params.Args)

		if reply != nil {
			client.Write(reply.Bytes())
		} else {
			client.Write(unknownErrReplyBytes)
		}
	}
}

func (handler *RespHandler) Close() error {
	logger.Info("Closing resp handler")

	handler.closing.Set(true)
	handler.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})

	handler.db.Close()

	return nil
}
