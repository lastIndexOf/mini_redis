package database

import "strings"

var commands = make(map[string]*command)

type command struct {
	executor ExecFunc
	argsLen  int
}

func RegisterCommand(name string, executor ExecFunc, argsLen int) {
	name = strings.ToLower(name)
	commands[name] = &command{
		executor: executor,
		argsLen:  argsLen,
	}
}
