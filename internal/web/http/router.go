package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	http_swagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"

	admin_handler "gitlab.ozon.dev/alexplay1224/homework/internal/web/http/admin"
	order_handler "gitlab.ozon.dev/alexplay1224/homework/internal/web/http/order"

	_ "gitlab.ozon.dev/alexplay1224/homework/docs" // docs needed for swagger
	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
	admin_service "gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"
	audit_logger_storage "gitlab.ozon.dev/alexplay1224/homework/internal/service/auditlogger"
	order_service "gitlab.ozon.dev/alexplay1224/homework/internal/service/order"
)

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

type auditLoggerStorage interface {
	GetAndMarkLogs(context.Context, int) ([]models.Log, error)
	UpdateLog(context.Context, int, int, int) error
	CreateLog(context.Context, []models.Log) error
}

// App is a structure for an app
type App struct {
	orderService       order_service.Service
	adminService       admin_service.Service
	auditLoggerService audit_logger_storage.Service
	Router             *mux.Router
}

// NewApp creates an instance of an App
func NewApp(ctx context.Context, cfg config.Config, logger *zap.Logger, orders orderStorage, admins adminStorage,
	logs auditLoggerStorage, txManager txManager, workerCount int, batchSize int, timeout time.Duration) (*App, error) {
	kafkaLogger, err := audit_logger_storage.NewService(ctx, cfg, logs, workerCount, batchSize, timeout)
	if err != nil {
		return nil, err
	}

	return &App{
		orderService:       *order_service.NewService(logger, orders, txManager),
		adminService:       *admin_service.NewService(logger, admins),
		auditLoggerService: *kafkaLogger,
		Router:             mux.NewRouter(),
	}, nil
}

// SetupRoutes setups all the routing
func (a *App) SetupRoutes(ctx context.Context) {
	impl := server{
		orders: *order_handler.NewHandler(&a.orderService),
		admins: *admin_handler.NewHandler(&a.adminService),
	}
	logger := AuditLoggerMiddleware{
		adminService:       a.adminService,
		auditLoggerService: a.auditLoggerService,
	}
	a.Router.Use(FieldLogger)

	authMiddleware := AuthMiddleware{
		adminService: a.adminService,
	}

	a.Router.HandleFunc("/orders",
		authMiddleware.BasicAuthChecker(ctx,
			logger.AuditLogger(ctx,
				a.wrapHandler(ctx, impl.orders.CreateOrder))).ServeHTTP).
		Methods(http.MethodPost)

	a.Router.HandleFunc("/orders",
		authMiddleware.BasicAuthChecker(ctx,
			a.wrapHandler(ctx, impl.orders.GetOrders)).ServeHTTP).
		Methods(http.MethodGet)

	a.Router.HandleFunc(fmt.Sprintf("/orders/{%s:[0-9]+}", order_handler.OrderIDParam),
		authMiddleware.BasicAuthChecker(ctx,
			logger.AuditLogger(ctx,
				a.wrapHandler(ctx, impl.orders.DeleteOrder))).ServeHTTP).
		Methods(http.MethodDelete)

	a.Router.HandleFunc("/orders/process",
		authMiddleware.BasicAuthChecker(ctx,
			logger.AuditLogger(ctx,
				a.wrapHandler(ctx, impl.orders.UpdateOrder))).ServeHTTP).
		Methods(http.MethodPost)

	a.Router.HandleFunc("/admins", a.wrapHandler(ctx, impl.admins.CreateAdmin)).
		Methods(http.MethodPost)

	a.Router.HandleFunc(fmt.Sprintf("/admins/{%s:[a-zA-Z0-9]+}",
		admin_handler.AdminUsernameParam), a.wrapHandler(ctx, impl.admins.UpdateAdmin)).
		Methods(http.MethodPost)

	a.Router.HandleFunc(fmt.Sprintf("/admins/{%s:[a-zA-Z0-9]+}",
		admin_handler.AdminUsernameParam), a.wrapHandler(ctx, impl.admins.DeleteAdmin)).
		Methods(http.MethodDelete)
}

func (a *App) wrapHandler(ctx context.Context, handler func(context.Context, http.ResponseWriter,
	*http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(ctx, w, r)
	}
}

type server struct {
	orders order_handler.Handler
	admins admin_handler.Handler
}

// @securityDefinitions.basic BasicAuth

// Run runs the app
// @title			PVZ API Documentation
// @version		1.0
// @description	This is a sample server for Swagger in Go.
// @host			localhost:9000
// @BasePath		/
func (a *App) Run(ctx context.Context) error {
	a.SetupRoutes(ctx)

	// Путь для отображения Swagger UI
	a.Router.PathPrefix("/swagger/").Handler(http_swagger.WrapHandler)
	if err := http.ListenAndServe("localhost:9000", a.Router); err != nil {
		log.Fatal(err)
	}

	return nil
}
