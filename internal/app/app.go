package app

import (
	"fmt"
	"github.com/Verce11o/yata-comments/config"
	commentGRPC "github.com/Verce11o/yata-comments/internal/handler/grpc"
	"github.com/Verce11o/yata-comments/internal/lib/logger"
	"github.com/Verce11o/yata-comments/internal/metrics/trace"
	"github.com/Verce11o/yata-comments/internal/repository/minio"
	"github.com/Verce11o/yata-comments/internal/repository/postgres"
	"github.com/Verce11o/yata-comments/internal/repository/redis"
	"github.com/Verce11o/yata-comments/internal/service"
	pb "github.com/Verce11o/yata-protos/gen/go/comments"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	log := logger.NewLogger()
	cfg := config.LoadConfig()

	tracer := trace.InitTracer("yata-comments")

	// Init repos
	db := postgres.NewPostgres(cfg)
	repo := postgres.NewCommentsPostgres(db, tracer.Tracer)

	rdb := redis.NewRedis(cfg)
	redisRepo := redis.NewCommentsRedis(rdb, tracer.Tracer)

	minioClient := minio.NewMinio(cfg)
	minioRepo := minio.NewCommentMinio(minioClient, tracer.Tracer)

	s := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor(
		otelgrpc.WithTracerProvider(tracer.Provider),
		otelgrpc.WithPropagators(propagation.TraceContext{}),
	)))

	commentService := service.NewCommentService(log, tracer.Tracer, repo, redisRepo, minioRepo)

	pb.RegisterCommentsServer(s, commentGRPC.NewCommentGRPC(log, tracer.Tracer, commentService))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.App.Port))

	if err != nil {
		log.Info("failed to listen: %v", err)
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Infof("error while listen server: %s", err)
		}
	}()

	log.Info(fmt.Sprintf("server listening at %s", lis.Addr().String()))

	defer log.Sync()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.GracefulStop()

	if err := db.Close(); err != nil {
		log.Infof("error while close db: %s", err)
	}

}
