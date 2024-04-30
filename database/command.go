package database

import "strings"

var commands = make(map[string]*command)

type command struct {
	executor ExecFunc
	argsLen  int // 如果为负数表示不定长，绝对值表示至少要有几个参数
}

func RegisterCommand(name string, executor ExecFunc, argsLen int) {
	name = strings.ToLower(name)
	commands[name] = &command{
		executor: executor,
		argsLen:  argsLen,
	}
}
