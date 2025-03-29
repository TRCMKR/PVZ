package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "gitlab.ozon.dev/alexplay1224/homework/docs"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
	admin_Service "gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"
	audit_Logger_Storage "gitlab.ozon.dev/alexplay1224/homework/internal/service/auditlogger"
	order_Service "gitlab.ozon.dev/alexplay1224/homework/internal/service/order"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/tx_manager"
	admin_Handler "gitlab.ozon.dev/alexplay1224/homework/internal/web/admin"
	order_Handler "gitlab.ozon.dev/alexplay1224/homework/internal/web/order"
)

type orderStorage interface {
	AddOrder(context.Context, models.Order) error
	RemoveOrder(context.Context, int) error
	UpdateOrder(context.Context, int, models.Order) error
	GetByID(context.Context, int) (models.Order, error)
	GetByUserID(context.Context, int, int) ([]models.Order, error)
	GetReturns(context.Context) ([]models.Order, error)
	GetOrders(context.Context, []query.Cond, int, int) ([]models.Order, error)
	Contains(context.Context, int) (bool, error)
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
	RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunReadCommitted(ctx context.Context, fn func(ctxTx context.Context) error) error
	GetQueryEngine(ctx context.Context) tx_manager.Database
}

type auditLoggerStorage interface {
	CreateLog(context.Context, []models.Log) error
}

type App struct {
	orderService       order_Service.Service
	adminService       admin_Service.Service
	auditLoggerService audit_Logger_Storage.Service
	Router             *mux.Router
}

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
