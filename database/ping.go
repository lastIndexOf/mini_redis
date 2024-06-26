package database

import (
	"github.com/lastIndexOf/mini_redis/interface/resp"
	"github.com/lastIndexOf/mini_redis/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

func init() {
	RegisterCommand("ping", Ping, 1)
}
