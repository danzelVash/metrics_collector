package grpc

import (
	"context"
	contracts "metrics_collector/collector/pkg/grpc_contracts/v1/grpc"
)

func (s *Server) RegisterService(context.Context, *contracts.RegisterReq) (*contracts.RegisterResp, error) {
	return nil, nil
}

func (s *Server) RefreshToken(context.Context, *contracts.RefreshTokenReq) (*contracts.RefreshTokenResp, error) {
	return nil, nil
}
func (s *Server) SendMetrics(contracts.Collector_SendMetricsServer) error {
	return nil
}
