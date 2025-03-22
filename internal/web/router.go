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
	adminServicePkg "gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"
	auditLoggerStoragePkg "gitlab.ozon.dev/alexplay1224/homework/internal/service/auditlogger"
	orderServicePkg "gitlab.ozon.dev/alexplay1224/homework/internal/service/order"
	adminHandlerPkg "gitlab.ozon.dev/alexplay1224/homework/internal/web/admin"
	orderHandlerPkg "gitlab.ozon.dev/alexplay1224/homework/internal/web/order"
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

type auditLoggerStorage interface {
	CreateLog(context.Context, models.Log) error
}

type App struct {
	orderService       orderServicePkg.Service
	adminService       adminServicePkg.Service
	auditLoggerService auditLoggerStoragePkg.Service
	Router             *mux.Router
}

func NewApp(ctx context.Context, orders orderStorage, admins adminStorage, logs auditLoggerStorage, workerCount int, batchSize int, timeout time.Duration) *App {
	return &App{
		orderService:       *orderServicePkg.NewService(orders),
		adminService:       *adminServicePkg.NewService(admins),
		auditLoggerService: *auditLoggerStoragePkg.NewService(ctx, logs, workerCount, batchSize, timeout),
		Router:             mux.NewRouter(),
	}
}

func (a *App) SetupRoutes(ctx context.Context) {
	impl := server{
		orders: *orderHandlerPkg.NewHandler(&a.orderService),
		admins: *adminHandlerPkg.NewHandler(&a.adminService),
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

	a.Router.HandleFunc(fmt.Sprintf("/orders/{%s:[0-9]+}", orderHandlerPkg.OrderIDParam),
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
		adminHandlerPkg.AdminUsernameParam), a.wrapHandler(ctx, impl.admins.UpdateAdmin)).
		Methods(http.MethodPost)

	a.Router.HandleFunc(fmt.Sprintf("/admins/{%s:[a-zA-Z0-9]+}",
		adminHandlerPkg.AdminUsernameParam), a.wrapHandler(ctx, impl.admins.DeleteAdmin)).
		Methods(http.MethodDelete)
}

func (a *App) wrapHandler(ctx context.Context, handler func(context.Context, http.ResponseWriter,
	*http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(ctx, w, r)
	}
}

type server struct {
	orders orderHandlerPkg.Handler
	admins adminHandlerPkg.Handler
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
