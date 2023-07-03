package main

import (
	"metrics_collector/collector/internal/config"
	"metrics_collector/collector/pkg/models/grpc"
	"metrics_collector/pkg/logging"
)

func main() {
	globalConfig := config.GetConfig("collector/internal/config", "config")
	//fmt.Println(globalConfig)

	logger := logging.GetLogger(globalConfig.AppConf.LogLevel)

	serverConf := grpc.ServerConfig{
		Listen: struct {
			Type       string
			BindIp     string
			Port       string
			SocketFile string
		}(globalConfig.Grpc.Listen),

		Options: struct {
			MaxConcurrentStreams uint32
			InitialWindowSize    int32
		}(globalConfig.Grpc.Options),
	}

	srv := new(grpc.Server)
	if err := srv.Run(serverConf); err != nil {
		logger.Fatalf("error running server: %s", err.Error())
	}
}
