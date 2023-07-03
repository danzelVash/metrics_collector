package config

import (
	"github.com/pkg/errors"
)

var (
	InvalidIpAddress      = errors.New("ip address is invalid")
	InvalidPort           = errors.New("port is invalid")
	InvalidSocketFileName = errors.New("socket file name is invalid")
	ListenTypeUndefined   = errors.New("listen type was not found")
)
