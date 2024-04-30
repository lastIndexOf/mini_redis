package database

import (
	"strconv"
	"strings"

	"github.com/lastIndexOf/mini_redis/config"
	db "github.com/lastIndexOf/mini_redis/interface/database"
	"github.com/lastIndexOf/mini_redis/interface/resp"
	"github.com/lastIndexOf/mini_redis/lib/logger"
	"github.com/lastIndexOf/mini_redis/resp/reply"
)

type Database struct {
	dbs []*DB
}

func NewDatabase() *Database {
	database := &Database{}

	if config.Properties.Databases <= 0 {
		config.Properties.Databases = 16
	}

	database.dbs = make([]*DB, 0, config.Properties.Databases)
	for index := range config.Properties.Databases {
		database.dbs = append(database.dbs, MakeDB(index))
	}

	return database
}

func (db *Database) Exec(client resp.Connection, args db.CmdLine) resp.Reply {
	go func() {
		if r := recover(); r != nil {
			logger.Error(r)
		}
	}()

	command := strings.ToLower(string(args[0]))

	if command == "select" {
		if len(args) != 2 {
			return reply.MakeArgNumErrReply("select")
		}

		return Select(client, db, args[1:])
	}

	return db.dbs[client.GetDBIndex()].Exec(client, args)
}

func (db *Database) Close() {}

func (db *Database) AfterClientClose(c resp.Connection) {}

func Select(conn resp.Connection, database *Database, args [][]byte) resp.Reply {
	index, err := strconv.Atoi(string(args[0]))

	if err != nil {
		return reply.MakeStandardErrReply("ERR invalid DB index")
	}

	if index >= len(database.dbs) {
		return reply.MakeStandardErrReply("DB index is out of range")
	}

	conn.SelectDB(index)

	return reply.MakeOkReply()
}
