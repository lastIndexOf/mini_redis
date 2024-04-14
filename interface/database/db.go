package database

import "github.com/lastIndexOf/mini_redis/interface/resp"

type CmdLine = [][]byte

type Database interface {
	Exec(client resp.Connection, args CmdLine) resp.Reply
	Close()
	AfterClientClose(c resp.Connection)
}

type DataEntity struct {
	Data any
}
