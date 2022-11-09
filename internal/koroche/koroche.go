package implementation

import (
	"context"
	"fmt"

	"github.com/alexadastra/koroche/internal/app/models"
	"github.com/alexadastra/koroche/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Koroche handles stuff
type Koroche struct {
	service Service
}

// Service is domain laywer wrapper
type Service interface {
	AddURL(ctx context.Context, url models.UserURL) (models.ShortURL, error)
	GetURL(ctx context.Context, url models.ShortURL) (models.UserURL, error)
}

// NewKoroche creates new server
func NewKoroche(service Service) api.KorocheServer {
	return &Koroche{service: service}
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

// AddURL adds short_url -> user_url pair to the system
func (impl *Koroche) AddURL(ctx context.Context, request *api.AddURLRequest) (*api.AddURLResponse, error) {
	url, err := models.NewUserURL(request.UserUrl.Value)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	newURL, err := impl.service.AddURL(ctx, url)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.AddURLResponse{ShortUrl: &api.ShortURL{Value: newURL.ToString()}}, nil
}

// GetURL fetches user url by short url
func (impl *Koroche) GetURL(ctx context.Context, request *api.GetURLRequest) (*api.GetURLResponse, error) {
	url := models.NewShortURL(request.Url.Value)

	userURL, err := impl.service.GetURL(ctx, url)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.GetURLResponse{UserUrl: &api.UserURL{Value: userURL.ToString()}}, nil
}
