package implementation

import (
	"context"
	"fmt"

	"github.com/alexadastra/koroche/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Koroche handles stuff
type Koroche struct{}

// NewKoroche creates new server
func NewKoroche() api.KorocheServer {
	return &Koroche{}
}

// Ping returns "pong" if ping in pinged
func (impl *Koroche) Ping(ctx context.Context, request *api.PingRequest) (*api.PingResponse, error) {
	if request.Value == "ping" {
		return &api.PingResponse{
			Code:  200,
			Value: "pong",
		}, nil
	}
	return nil, fmt.Errorf("unknown request message: %s", request.Value)
}

func (impl *Koroche) AddURL(ctx context.Context, request *api.AddURLRequest) (*api.AddURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented!")
}

func (impl *Koroche) GetURL(ctx context.Context, request *api.GetURLRequest) (*api.GetURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented!")
}
