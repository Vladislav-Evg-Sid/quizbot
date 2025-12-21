package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/config"
	server "github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/api/admin_service_api"
	admininfoupsertconsumer "github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/consumer/admin_Info_upsert_consumer"
	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/pb/admins_api"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AppRun(api server.AdminServiceAPI, adminInfoUpsertConsumer *admininfoupsertconsumer.AdminInfoUpsertConsumer, cfg *config.Config) {
	go adminInfoUpsertConsumer.Consume(context.Background())
	go func() {
		if err := runGRPCServer(api, cfg); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %v", err))
		}
	}()

	if err := runGatewayServer(cfg); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %v", err))
	}
}

func runGRPCServer(api server.AdminServiceAPI, cfg *config.Config) error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Service.GRPCPort))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	admins_api.RegisterAdminsServiceServer(s, &api)

	slog.Info("gRPC-server server listening on :" + strconv.Itoa(cfg.Service.GRPCPort))
	return s.Serve(lis)
}

func runGatewayServer(cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	swaggerPath := cfg.Service.SwaggerPath
	if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
		panic(fmt.Errorf("swagger file not found: %s", swaggerPath))
	}

	r := chi.NewRouter()
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerPath)
	})

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := admins_api.RegisterAdminsServiceHandlerFromEndpoint(ctx, mux, ":"+strconv.Itoa(cfg.Service.GRPCPort), opts)
	if err != nil {
		panic(err)
	}

	r.Mount("/", mux)

	slog.Info("gRPC-Gateway server listening on :" + strconv.Itoa(cfg.Service.RESTPort))
	return http.ListenAndServe(":"+strconv.Itoa(cfg.Service.RESTPort), r)
}
