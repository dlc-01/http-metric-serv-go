package app

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/encryption"
	pb "github.com/dlc-01/http-metric-serv-go/internal/protobuf"
	"github.com/dlc-01/http-metric-serv-go/internal/server/grpcserver"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/decryptor"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/subnet"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/all"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/db"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/json"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/jsonbatch"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/url"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/checkinghash"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/storagesync"
)

// Run — function that set up and run sever.
func Run(cfg *config.ServerConfig) {
	router, err := setupRouter(cfg)
	if err != nil {
		logging.Fatalf("error while starting server: %s", err)
	}
	if cfg.GRPC {
		listen, err := net.Listen("tcp", cfg.GRPCAddress)
		if err != nil {
			logging.Fatalf("cannot set up grpc server: %w", err)
		}
		grpcServer := grpc.NewServer()
		grpcMetric := grpcserver.CreateGrpcServer()
		pb.RegisterMetricsServiceServer(grpcServer, grpcMetric)
		go func() {
			if err := grpcServer.Serve(listen); err != nil {
				logging.Fatalf("cannot run grpc server: %w", err)
			}
		}()
		logging.Infof("start grpc at address: %s", cfg.GRPCAddress)
	}

	err = router.Run(cfg.ServerAddress)
	if err != nil {
		logging.Fatalf("error while starting server: %s", err)
	}

}

func setupRouter(cfg *config.ServerConfig) (*gin.Engine, error) {
	router := gin.Default()
	if cfg.TrustedSubnet != "" {
		router.Use(subnet.CheckSubnet(cfg.TrustedSubnet))
	}
	if cfg.PathCryptoKey != "" {
		if err := encryption.InitDecryptor(cfg.PathCryptoKey); err != nil {
			return nil, fmt.Errorf("cannot create decryptor: %w", err)
		}
		router.Use(decryptor.DecryptMiddleware())
	}
	router.Use(logging.GetMiddlewareLogger(), gzip.Gzip(gzip.BestSpeed))
	if cfg.HashKey != "" {
		router.Use(checkinghash.CheckHash(cfg.HashKey))
	}
	router.POST("/value/", json.ValueJSONHandler)
	router.GET("/value/:types/:name", url.ValueHandler)
	router.GET("/", all.ShowMetrics)
	router.GET("/ping", db.PingDB)
	updateRouterGroup := router.Group("/")
	if cfg.DatabaseAddress == "" {
		updateRouterGroup.Use(storagesync.GetSyncMiddleware())
	}

	{
		updateRouterGroup.POST("/update", json.UpdateJSONHandler)
		updateRouterGroup.POST("/update/:types/:name/:value", url.UpdateHandler)
		updateRouterGroup.POST("/updates", jsonbatch.UpdatesButchJSONHandler)

	}
	return router, nil
}
