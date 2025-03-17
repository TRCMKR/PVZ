package web

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
	adminServicePkg "gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"
	orderServicePkg "gitlab.ozon.dev/alexplay1224/homework/internal/service/order"
	adminHandlerPkg "gitlab.ozon.dev/alexplay1224/homework/internal/web/admin"
	orderHandlerPkg "gitlab.ozon.dev/alexplay1224/homework/internal/web/order"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "gitlab.ozon.dev/alexplay1224/homework/docs"
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

type App struct {
	orderService orderServicePkg.Service
	adminService adminServicePkg.Service
	router       *mux.Router
}

func NewApp(orders orderStorage, admins adminStorage) *App {
	return &App{
		orderService: *orderServicePkg.NewService(orders),
		adminService: *adminServicePkg.NewService(admins),
		router:       mux.NewRouter(),
	}
}

func (a *App) SetupRoutes(ctx context.Context) {
	impl := server{
		orders: *orderHandlerPkg.NewHandler(&a.orderService),
		admins: *adminHandlerPkg.NewHandler(&a.adminService),
	}
	a.router.Use(FieldLogger)

	authMiddleware := AuthMiddleware{
		adminService: a.adminService,
	}

	a.router.HandleFunc("/orders", authMiddleware.BasicAuthChecker(ctx,
		a.wrapHandler(ctx, impl.orders.CreateOrder)).ServeHTTP).
		Methods(http.MethodPost)
	//	@Summary		Get Orders
	//	@Description	Fetches all orders.
	//	@Tags			orders
	//	@Accept			json
	//	@Produce		json
	//	@Success		200	{array}		models.Order
	//	@Failure		400	{object}	models.ErrorResponse
	//	@Router			/orders [get]
	a.router.HandleFunc("/orders", authMiddleware.BasicAuthChecker(ctx,
		a.wrapHandler(ctx, impl.orders.GetOrders)).ServeHTTP).
		Methods(http.MethodGet)

	//	@Summary		Delete an Order
	//	@Description	Deletes an order by ID.
	//	@Tags			orders
	//	@Accept			json
	//	@Produce		json
	//	@Param			id	path		int		true	"Order ID"
	//	@Success		200	{string}	string	"Order deleted successfully"
	//	@Failure		404	{object}	models.ErrorResponse
	//	@Router			/orders/{id} [delete]
	a.router.HandleFunc(fmt.Sprintf("/orders/{%s:[0-9]+}", orderHandlerPkg.OrderIDParam),
		authMiddleware.BasicAuthChecker(ctx, a.wrapHandler(ctx, impl.orders.DeleteOrder)).ServeHTTP).
		Methods(http.MethodDelete)

	//	@Summary		Process Orders
	//	@Description	Updates the status of orders.
	//	@Tags			orders
	//	@Accept			json
	//	@Produce		json
	//	@Param			orders	body		[]models.Order	true	"List of Orders"
	//	@Success		200		{string}	string			"Orders processed successfully"
	//	@Failure		400		{object}	models.ErrorResponse
	//	@Router			/orders/process [post]
	a.router.HandleFunc("/orders/process", authMiddleware.BasicAuthChecker(ctx,
		a.wrapHandler(ctx, impl.orders.UpdateOrders)).ServeHTTP).
		Methods(http.MethodPost)

	//	@Summary		Create Admin
	//	@Description	Creates a new admin user.
	//	@Tags			admins
	//	@Accept			json
	//	@Produce		json
	//	@Param			admin	body		models.Admin	true	"Admin"
	//	@Success		200		{object}	models.Admin
	//	@Failure		400		{object}	models.ErrorResponse
	//	@Router			/admins [post]
	a.router.HandleFunc("/admins", a.wrapHandler(ctx, impl.admins.CreateAdmin)).
		Methods(http.MethodPost)

	//	@Summary		Update Admin
	//	@Description	Updates an existing admin user by username.
	//	@Tags			admins
	//	@Accept			json
	//	@Produce		json
	//	@Param			admin_username	path		string			true	"Admin Username"
	//	@Param			admin			body		models.Admin	true	"Admin"
	//	@Success		200				{object}	models.Admin
	//	@Failure		400				{object}	models.ErrorResponse
	//	@Router			/admins/{admin_username} [post]
	a.router.HandleFunc(fmt.Sprintf("/admins/{%s:[a-zA-Z0-9]+}",
		adminHandlerPkg.AdminUsernameParam), a.wrapHandler(ctx, impl.admins.UpdateAdmin)).
		Methods(http.MethodPost)

	//	@Summary		Delete Admin
	//	@Description	Deletes an admin user by username.
	//	@Tags			admins
	//	@Accept			json
	//	@Produce		json
	//	@Param			admin_username	path		string	true	"Admin Username"
	//	@Success		200				{string}	string	"Admin deleted successfully"
	//	@Failure		404				{object}	models.ErrorResponse
	//	@Router			/admins/{admin_username} [delete]
	a.router.HandleFunc(fmt.Sprintf("/admins/{%s:[a-zA-Z0-9]+}",
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

// @title			PVZ API Documentation
// @version		1.0
// @description	This is a sample server for Swagger in Go.
// @host			localhost:9000
// @BasePath		/
func (a *App) Run(ctx context.Context) error {
	a.SetupRoutes(ctx)

	// Путь для отображения Swagger UI
	a.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	go func() {
		if err := http.ListenAndServe("localhost:9000", a.router); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
