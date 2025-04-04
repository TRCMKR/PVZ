package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	// docs ...
	_ "gitlab.ozon.dev/alexplay1224/homework/docs"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
	admin_Service "gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"
	audit_Logger_Storage "gitlab.ozon.dev/alexplay1224/homework/internal/service/auditlogger"
	order_Service "gitlab.ozon.dev/alexplay1224/homework/internal/service/order"
	admin_Handler "gitlab.ozon.dev/alexplay1224/homework/internal/web/admin"
	order_Handler "gitlab.ozon.dev/alexplay1224/homework/internal/web/order"
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
	CreateLog(context.Context, []models.Log) error
}

// App ...
type App struct {
	orderService       order_Service.Service
	adminService       admin_Service.Service
	auditLoggerService audit_Logger_Storage.Service
	Router             *mux.Router
}

// NewApp ...
func NewApp(ctx context.Context, orders orderStorage, admins adminStorage, logs auditLoggerStorage, txManager txManager,
	workerCount int, batchSize int, timeout time.Duration) (*App, error) {
	logger, err := audit_Logger_Storage.NewService(ctx, logs, workerCount, batchSize, timeout)
	if err != nil {
		return nil, err
	}

	return &App{
		orderService:       *order_Service.NewService(orders, txManager),
		adminService:       *admin_Service.NewService(admins),
		auditLoggerService: *logger,
		Router:             mux.NewRouter(),
	}, nil
}

// SetupRoutes ...
func (a *App) SetupRoutes(ctx context.Context) {
	impl := server{
		orders: *order_Handler.NewHandler(&a.orderService),
		admins: *admin_Handler.NewHandler(&a.adminService),
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

	a.Router.HandleFunc(fmt.Sprintf("/orders/{%s:[0-9]+}", order_Handler.OrderIDParam),
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
		admin_Handler.AdminUsernameParam), a.wrapHandler(ctx, impl.admins.UpdateAdmin)).
		Methods(http.MethodPost)

	a.Router.HandleFunc(fmt.Sprintf("/admins/{%s:[a-zA-Z0-9]+}",
		admin_Handler.AdminUsernameParam), a.wrapHandler(ctx, impl.admins.DeleteAdmin)).
		Methods(http.MethodDelete)
}

func (a *App) wrapHandler(ctx context.Context, handler func(context.Context, http.ResponseWriter,
	*http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(ctx, w, r)
	}
}

type server struct {
	orders order_Handler.Handler
	admins admin_Handler.Handler
}

// @securityDefinitions.basic BasicAuth

// Run ...
// @title			PVZ API Documentation
// @version		1.0
// @description	This is a sample server for Swagger in Go.
// @host			localhost:9000
// @BasePath		/
func (a *App) Run(ctx context.Context) error {
	a.SetupRoutes(ctx)

	// Путь для отображения Swagger UI
	a.Router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	if err := http.ListenAndServe("localhost:9000", a.Router); err != nil {
		log.Fatal(err)
	}

	return nil
}
