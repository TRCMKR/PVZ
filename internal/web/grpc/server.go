package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
	admin_service "gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"
	order_service "gitlab.ozon.dev/alexplay1224/homework/internal/service/order"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web/grpc/admin"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web/grpc/order"
	admin_proto "gitlab.ozon.dev/alexplay1224/homework/pkg/api/admin/proto"
	order_proto "gitlab.ozon.dev/alexplay1224/homework/pkg/api/order/proto"
	"gitlab.ozon.dev/alexplay1224/homework/pkg/monitoring"
)

// Server is a struct for a grpc server
type Server struct {
	orderHandler order.Handler
	adminHandler admin.Handler
}

type orderStorage interface {
	AddOrder(context.Context, pgx.Tx, models.Order) error
	RemoveOrder(context.Context, pgx.Tx, int) error
	UpdateOrder(context.Context, pgx.Tx, int, models.Order) error
	GetByID(context.Context, pgx.Tx, int) (models.Order, error)
	GetByUserID(context.Context, pgx.Tx, int, int) ([]models.Order, error)
	GetReturns(context.Context, pgx.Tx) ([]models.Order, error)
	GetOrders(context.Context, pgx.Tx, []query.Cond, int, int) ([]models.Order, error)
	Contains(context.Context, pgx.Tx, int) (bool, error)
}

type adminStorage interface {
	CreateAdmin(context.Context, models.Admin) error
	GetAdminByUsername(context.Context, string) (models.Admin, error)
	UpdateAdmin(context.Context, int, models.Admin) error
	DeleteAdmin(context.Context, string) error
	ContainsUsername(context.Context, string) (bool, error)
	ContainsID(context.Context, int) (bool, error)
}

type txManager interface {
	RunSerializable(context.Context, func(context.Context, pgx.Tx) error) error
	RunRepeatableRead(context.Context, func(context.Context, pgx.Tx) error) error
	RunReadCommitted(context.Context, func(context.Context, pgx.Tx) error) error
}

// NewServer creates instance of a grpc server
func NewServer(logger *zap.Logger, orders orderStorage, admins adminStorage, txManager txManager) *Server {
	orderHandler := order.NewHandler(logger.With(
		zap.String("layer", "handler"),
		zap.String("domain", "orders"),
	), *order_service.NewService(logger.With(
		zap.String("layer", "service"),
		zap.String("domain", "orders"),
	), orders, txManager))
	adminHandler := admin.NewHandler(logger.With(
		zap.String("layer", "handler"),
		zap.String("domain", "admins"),
	), *admin_service.NewService(logger.With(
		zap.String("layer", "service"),
		zap.String("domain", "admins"),
	), admins))

	return &Server{
		orderHandler: *orderHandler,
		adminHandler: *adminHandler,
	}
}

// Run runs a grpc server
func (s *Server) Run(ctx context.Context, cfg config.Config, logger *zap.Logger) error {
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort())
	if err != nil {
		logger.Fatal("failed to listen",
			zap.String("GRPC_PORT", cfg.GRPCPort()),
			zap.Error(err),
		)
	}

	errCh := make(chan error)
	monitoring.StartMetricsServer(errCh)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(MetricsInterceptor()),
	)

	order_proto.RegisterOrderServiceServer(grpcServer, &s.orderHandler)
	admin_proto.RegisterAdminServiceServer(grpcServer, &s.adminHandler)

	logger.Info(fmt.Sprintf("server listening at %v", lis.Addr()))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()

		return nil
	case err := <-errCh:
		return err
	}
}
