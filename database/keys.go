package database

import (
	"github.com/lastIndexOf/mini_redis/interface/resp"
	"github.com/lastIndexOf/mini_redis/lib/wildcard"
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
	key := string(args[0])

	entity, exists := db.GetEntity(key)

	if !exists {
		return reply.MakeStatusReply("none")
	}

	switch entity.Data.(type) {
	// support string only now
	case []byte:
		return reply.MakeStatusReply("string")
	}

	return reply.MakeUnknownErrReply()
}

func Rename(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	target := string(args[1])

	val, exists := db.GetEntity(key)

	if !exists {
		return reply.MakeStandardErrReply("no such key " + key)
	}

	db.PutEntity(target, val)
	db.Remove(key)

	return reply.MakeOkReply()
}

func Renamenx(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	target := string(args[1])

	if _, exists := db.GetEntity(target); exists {
		return reply.MakeIntReply(0)
	}

	val, exists := db.GetEntity(key)

	if !exists {
		return reply.MakeStandardErrReply("no such key " + key)
	}

	db.PutEntity(target, val)
	db.Remove(key)

	return reply.MakeIntReply(1)
}

func Keys(db *DB, args [][]byte) resp.Reply {
	pattern, err := wildcard.CompilePattern(string(args[0]))

	if err != nil {
		return reply.MakeStandardErrReply("invalid key")
	}

	ret := make([][]byte, 0)

	db.data.ForEach(func(key string, val any) bool {
		if pattern.IsMatch(key) {
			ret = append(ret, []byte(key))
		}

		return true
	})

	return reply.MakeMultiBulkReply(ret)
}

func init() {
	RegisterCommand("del", Del, -2)
	RegisterCommand("exists", Exits, -2)
	RegisterCommand("flushdb", Flush, -1)
	RegisterCommand("type", Type, 2)
	RegisterCommand("rename", Rename, 3)
	RegisterCommand("renamenx", Renamenx, 3)
	RegisterCommand("keys", Keys, -2)
}
