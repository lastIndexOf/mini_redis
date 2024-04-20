package database

import (
	db "github.com/lastIndexOf/mini_redis/interface/database"
	"github.com/lastIndexOf/mini_redis/interface/resp"
	"github.com/lastIndexOf/mini_redis/resp/reply"
)

type EchoDatabase struct{}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}

func (e *EchoDatabase) Exec(client resp.Connection, args db.CmdLine) resp.Reply {
	return reply.MakeMultiBulkReply(args)
}

func (e *EchoDatabase) Close() {}

func (e *EchoDatabase) AfterClientClose(c resp.Connection) {}
