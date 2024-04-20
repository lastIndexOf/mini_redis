package database

import (
	"strings"

	"github.com/lastIndexOf/mini_redis/datastruct/dict"
	"github.com/lastIndexOf/mini_redis/interface/database"
	"github.com/lastIndexOf/mini_redis/interface/resp"
	"github.com/lastIndexOf/mini_redis/resp/reply"
)

// Redis 内核的每个数据库实例
type DB struct {
	index int
	data  dict.Dict
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

func (db *DB) Exec(client resp.Connection, args database.CmdLine) resp.Reply {
	name := strings.ToLower(string(args[0]))
	cmd, ok := commands[name]

	if !ok {
		return reply.MakeUnknownErrReply()
	}

	if !validateArgs(args, cmd.argsLen) {
		return reply.MakeArgNumErrReply(name)
	}

	return cmd.executor(db, args[1:])
}

func (db *DB) Close() {}

func (db *DB) AfterClientClose(c resp.Connection) {}

func validateArgs(args [][]byte, expected int) bool {
	argsLen := len(args)

	if argsLen >= 0 {
		return argsLen == expected
	}

	return argsLen >= -expected
}

func (db *DB) GetEntity(key string) (data *database.DataEntity, exists bool) {
	val, exists := db.data.Get(key)

	if !exists {
		return nil, false
	}

	return val.(*database.DataEntity), true
}

func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

func (db *DB) Remove(key string) int {
	return db.data.Remove(key)
}

func (db *DB) RemoveMany(keys ...string) int {
	count := 0

	for _, key := range keys {
		count += db.Remove(key)
	}

	return count
}

func (db *DB) Clear() {
	db.data.Clear()
}

func MakeDB() *DB {
	return &DB{
		data: dict.MakeSyncDict(),
	}
}
