package database

import (
	"github.com/lastIndexOf/mini_redis/interface/resp"
	"github.com/lastIndexOf/mini_redis/resp/reply"
)

func Del(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))

	for i, keyBytes := range keys {
		keys[i] = string(keyBytes)
	}

	deleted := db.RemoveMany(keys...)

	return reply.MakeIntReply(int64(deleted))
}

func Exits(db *DB, args [][]byte) resp.Reply {
	ret := 0

	for _, keyByte := range args {
		key := string(keyByte)

		_, exists := db.GetEntity(key)

		if exists {
			ret += 1
		}
	}

	return reply.MakeIntReply(int64(ret))
}

func Flush(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return reply.MakeOkReply()
}

func Type(db *DB, args [][]byte) resp.Reply {
	panic("implement me")
}

func Keys(db *DB, args [][]byte) resp.Reply {
	panic("implement me")
}

func init() {
	RegisterCommand("del", Del, -2)
	RegisterCommand("exists", Exits, -2)
	RegisterCommand("flushdb", Flush, -1)
}
