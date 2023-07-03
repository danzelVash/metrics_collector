package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	grpc2 "metrics_collector/collector/pkg/grpc_contracts/v1/grpc"
	"metrics_collector/pkg/logging"
	"net"
)

const (
	tcp  = "tcp"
	sock = "sock"
	unix = "unix"
)

type ServerConfig struct {
	Listen struct {
		Type       string
		BindIp     string
		Port       string
		SocketFile string
	}

	Options struct {
		MaxConcurrentStreams uint32
		InitialWindowSize    int32
	}
}

type Server struct {
	grpc2.UnimplementedCollectorServer
}

func (s *Server) register(conf ServerConfig) (*grpc.Server, net.Listener, error) {
	logger := logging.GetLogger("")

	var listener net.Listener

	if conf.Listen.Type == tcp {
		lis, err := net.Listen(tcp, fmt.Sprintf(":%s", conf.Listen.Port))
		if err != nil {
			logger.Fatalf("error creating %s listener: %s", tcp, err.Error())
		}
		listener = lis
	} else if conf.Listen.Type == sock {
		lis, err := net.Listen(unix, conf.Listen.SocketFile)
		if err != nil {
			logger.Fatalf("error creating %s listener: %s", tcp, err.Error())
		}
		listener = lis
	} else {
		return nil, nil, UnknownListenType
	}

	srv := grpc.NewServer(
		grpc.MaxConcurrentStreams(conf.Options.MaxConcurrentStreams),
		grpc.InitialWindowSize(conf.Options.InitialWindowSize),
	)

	reflection.Register(srv)
	grpc2.RegisterCollectorServer(srv, &Server{})

	return srv, listener, nil
}

func (s *Server) Run(conf ServerConfig) error {
	server, lis, err := s.register(conf)
	if err != nil {
		return err
	}
	return server.Serve(lis)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Shutdown(ctx)
}
