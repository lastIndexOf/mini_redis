package main

import (
	"fmt"
	"os"

	"github.com/lastIndexOf/mini_redis/config"
	"github.com/lastIndexOf/mini_redis/lib/logger"
	"github.com/lastIndexOf/mini_redis/tcp"
)

const configFile string = "redis.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6377,
}

func checkNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "mini_redis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if !checkNotExist(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}

	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Addr: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, tcp.MakeEchoHandler())

	if err != nil {
		logger.Error("Error starting server: ", err)
	}
}
