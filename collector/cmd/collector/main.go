package main

import (
	"fmt"
	"metrics_collector/collector/internal/config"
	"metrics_collector/collector/pkg/models/grpc"
	"metrics_collector/pkg/logging"
)

func main() {
	globalConfig := config.GetConfig("collector/internal/config", "config")
	fmt.Println(globalConfig)

	logger := logging.GetLogger(globalConfig.AppConf.LogLevel)

	serverConf := grpc.ServerConfig{
		Listen: struct {
			Type       string
			BindIp     string
			Port       string
			SocketFile string
		}(globalConfig.Grpc.Listen),
	}
}
