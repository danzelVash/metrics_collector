package config

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"metrics_collector/pkg/logging"
	"net"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	confFilePathDefault = "./"
	confFileNameDefault = "config"
	maxNumOfPorts       = 65535
	bindIpDefault       = "127.0.0.1"
	portDefault         = "8080"
	logLevelDefault     = "trace"
	socketFileDefault   = "app.sock"
	socketFileExt       = ".sock"
	tcp                 = "tcp"
	sock                = "sock"
)

type GlobalConfig struct {
	Grpc struct {
		Listen struct {
			Type       string `mapstructure:"type"`
			BindIp     string `mapstructure:"bind_ip"`
			Port       string `mapstructure:"port"`
			SocketFile string `mapstructure:"socket_file"`
		} `mapstructure:"listen"`
	} `mapstructure:"grpc"`

	AppConf struct {
		LogLevel string `mapstructure:"log_level"`
	} `mapstructure:"app_config"`
}

func (c *GlobalConfig) validate() error {
	if c.Grpc.Listen.Type == tcp {
		if c.Grpc.Listen.BindIp != "" {
			if net.ParseIP(c.Grpc.Listen.BindIp) == nil {
				return InvalidIpAddress
			}
		} else {
			c.Grpc.Listen.BindIp = bindIpDefault
		}

		if c.Grpc.Listen.Port != "" {
			port, err := strconv.Atoi(c.Grpc.Listen.Port)
			if err != nil {
				return errors.WithMessagef(InvalidPort, "error converting port into int: %s", err.Error())
			}
			if port < 0 || port > maxNumOfPorts {
				return errors.WithMessage(InvalidPort, "port number is out of range")
			}
		} else {
			c.Grpc.Listen.Port = portDefault
		}
	} else if c.Grpc.Listen.Type == sock {
		if c.Grpc.Listen.SocketFile == "" {
			c.Grpc.Listen.SocketFile = socketFileDefault
		} else {
			if filepath.Ext(c.Grpc.Listen.SocketFile) != socketFileExt {
				return InvalidSocketFileName
			}
		}
	} else {
		return ListenTypeUndefined
	}

	if _, ok := logging.LogLevels[c.AppConf.LogLevel]; !ok {
		c.AppConf.LogLevel = logLevelDefault
	}

	return nil
}

var (
	instance *GlobalConfig
	once     sync.Once
)

func GetConfig(pathToConf, confFileName string) *GlobalConfig {
	once.Do(func() {
		viperInst := viper.GetViper()

		if pathToConf != "" {
			viperInst.AddConfigPath(pathToConf)
		} else {
			viperInst.AddConfigPath(confFilePathDefault)
		}

		if confFileName != "" {
			viperInst.SetConfigName(confFileName)
		} else {
			viperInst.SetConfigName(confFileNameDefault)
		}

		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("error load env files: %s", err.Error())
		}

		viperInst.AutomaticEnv()

		if err := viperInst.ReadInConfig(); err != nil {
			log.Fatalf("error reading in config: %s", err.Error())
		}

		cfg := &GlobalConfig{}

		if err := viperInst.Unmarshal(cfg); err != nil {
			log.Fatalf("error reading conf: %s", err.Error())
		}

		instance = cfg
	})

	if err := instance.validate(); err != nil {
		log.Fatalf("invalid cfg: %s", err.Error())
	}
	return instance
}
