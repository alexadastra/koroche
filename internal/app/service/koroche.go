package service

import (
	"context"
	"fmt"

	"github.com/alexadastra/koroche/pkg/api"
)

// Koroche handles stuff
type Koroche struct{}

// NewKoroche creates new server
func NewKoroche() api.KorocheServer {
	return &Koroche{}
}

// Ping returns "pong" if ping in pinged
func (rt *Koroche) Ping(ctx context.Context, request *api.PingRequest) (*api.PingResponse, error) {
	if request.Value == "ping" {
		return &api.PingResponse{
			Code:  200,
			Value: "pong",
		}, nil
	}
	return nil, fmt.Errorf("unknown request message: %s", request.Value)
}
