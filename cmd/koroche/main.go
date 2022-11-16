package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	impl "github.com/alexadastra/koroche/internal/koroche"

	"github.com/alexadastra/ramme/config"
	"github.com/alexadastra/ramme/service"
	"github.com/alexadastra/ramme/system"

	koroche_service "github.com/alexadastra/koroche/internal/app/service"
	inmemory "github.com/alexadastra/koroche/internal/app/storage/in_memory"
	"github.com/alexadastra/koroche/internal/app/storage/mongodb"
	advanced "github.com/alexadastra/koroche/internal/config"
	"github.com/alexadastra/koroche/internal/swagger"

	"github.com/alexadastra/koroche/pkg/api"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Fetch flags configuration
	args := advanced.ParseFlags()
	config.ServiceName = args.ServiceName
	config.File = args.ConfigPath

	// Fetch basic config
	basicConfManager, basicConfWatcher, err := config.InitBasicConfig()
	if err != nil {
		panic(err)
	}
	basicConfig := basicConfManager.GetBasic()

	// Fetch advanced config
	advancedConfManager, advancedConfWatcher, err := advanced.InitAdvancedConfig()
	if err != nil {
		panic(err)
	}
	advancedConfig := advancedConfManager.Get()

	// logger.Info(advancedConfig)
	mongodb.Connect(ctx, advancedConfig.MongoDBDSN)

	// Serve
	g := system.NewGroupOperator()

	g.Add(func() error {
		return basicConfWatcher.Run()
	}, func(err error) {
		_ = basicConfWatcher.Close()
	})
	g.Add(func() error {
		return advancedConfWatcher.Run()
	}, func(err error) {
		_ = advancedConfWatcher.Close()
	})

	userGrpcServer := impl.NewKoroche(
		advancedConfig.PingMessage,
		koroche_service.NewService(
			"", // TODO: fetch from config
			inmemory.NewStorage(),
		),
	)

	baseGrpcServer := grpc.NewServer()
	api.RegisterKorocheServer(baseGrpcServer, userGrpcServer)

	mux := setupGRPCGateway(ctx, userGrpcServer)
	mux = setupSwagger(mux, basicConfig)

	run(
		ctx,
		g,
		baseGrpcServer,
		mux,
		basicConfig,
	)
}

func setupGRPCGateway(ctx context.Context, userGrpcServer api.KorocheServer) *http.ServeMux {
	mux := http.NewServeMux()
	rmux := runtime.NewServeMux()
	mux.Handle("/", rmux)

	err := api.RegisterKorocheHandlerServer(ctx, rmux, userGrpcServer)
	if err != nil {
		log.Fatal(err)
	}

	return mux
}

func setupSwagger(mux *http.ServeMux, basicConfig *config.BasicConfig) *http.ServeMux {
	// TODO: this prefix workaround should be solved better
	if basicConfig.IsLocalEnvironment {
		mux.Handle(swagger.Pattern, swagger.HandlerLocal)
	} else {
		mux.Handle(swagger.Pattern, swagger.HandlerK8S)
	}

	return mux
}

// EVERYTHING HERE IS PURE PLATFORM
func run(ctx context.Context, g *system.GroupOperator, baseGrpcServer *grpc.Server, mux http.Handler, basicConfig *config.BasicConfig) {
	// Configure service and get router
	router, logger, err := service.Setup(basicConfig)
	if err != nil {
		log.Fatal(err)
	}

	grpcStart, grpcStop := setupGRPC(baseGrpcServer, basicConfig)
	g.Add(grpcStart, grpcStop)
	logger.Warnf("Serving grpc address %s", fmt.Sprintf("%s:%d", basicConfig.Host, basicConfig.GRPCPort))

	httpStart, httpStop := setupHTTP(mux, &httpServerConfig{
		WriteTimeOut: 15 * time.Second,
		ReadTimeOut:  15 * time.Second,
		Host:         basicConfig.Host,
		Port:         basicConfig.HTTPPort,
	})
	g.Add(httpStart, httpStop)
	logger.Warnf("Serving http address %s", fmt.Sprintf("%s:%d", basicConfig.Host, basicConfig.HTTPPort))

	httpSecStart, httpSecStop := setupHTTP(router, &httpServerConfig{
		WriteTimeOut: 15 * time.Second,
		ReadTimeOut:  15 * time.Second,
		Host:         basicConfig.Host,
		Port:         basicConfig.HTTPSecondaryPort,
	})
	g.Add(httpSecStart, httpSecStop)

	signals := system.NewSignals()
	g.Add(func() error {
		return signals.Wait(logger, g)
	}, func(error) {})

	if err := g.Run(); err != nil {
		logger.Fatal(err)
	}
}

type httpServerConfig struct {
	WriteTimeOut time.Duration
	ReadTimeOut  time.Duration
	Host         string
	Port         int
}

func setupHTTP(handler http.Handler, conf *httpServerConfig) (func() error, func(error)) {
	newSrv := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		WriteTimeout: conf.WriteTimeOut,
		ReadTimeout:  conf.ReadTimeOut,
	}
	return newSrv.ListenAndServe, func(err error) { _ = newSrv.Close() }
}

// setupGRPC sets up gRPC server
func setupGRPC(baseGrpcServer *grpc.Server, basicConfig *config.BasicConfig) (func() error, func(error)) {
	grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", basicConfig.Host, basicConfig.GRPCPort))
	if err != nil {
		log.Fatal(err)
	}
	return func() error { return baseGrpcServer.Serve(grpcListener) }, func(err error) { _ = grpcListener.Close() }
}
