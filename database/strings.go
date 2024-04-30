package database

import (
	"github.com/lastIndexOf/mini_redis/interface/database"
	"github.com/lastIndexOf/mini_redis/interface/resp"
	"github.com/lastIndexOf/mini_redis/resp/reply"
)

func Get(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])

	entity, exists := db.GetEntity(key)

	if !exists {
		return reply.MakeNullBulkReply()
	}

	return reply.MakeBulkReply(entity.Data.([]byte))
}

func Set(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]

	db.PutEntity(key, &database.DataEntity{Data: val})

	return reply.MakeOkReply()
}

func SetNx(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]

	modified := db.PutIfAbsent(key, &database.DataEntity{Data: val})

	return reply.MakeIntReply(int64(modified))
}

func GetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]

	entity, exists := db.GetEntity(key)
	db.PutEntity(key, &database.DataEntity{Data: val})

	if exists {
		return reply.MakeBulkReply(entity.Data.([]byte))
	}

	return reply.MakeNullBulkReply()
}

func StrLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])

	entity, exists := db.GetEntity(key)

	if !exists {
		return reply.MakeNullBulkReply()
	}

	return reply.MakeIntReply(int64(len(entity.Data.([]byte))))
}

func init() {
	RegisterCommand("get", Get, 2)
	RegisterCommand("set", Set, 3)
	RegisterCommand("setnx", Set, 3)
	RegisterCommand("getset", GetSet, 3)
	RegisterCommand("strlen", StrLen, 2)
}
