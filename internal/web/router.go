package web

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
	"gitlab.ozon.dev/alexplay1224/homework/internal/service"

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
	orderService service.OrderService
	adminService service.AdminService
	router       *mux.Router
}

func NewApp(orders orderStorage, admins adminStorage) *App {
	return &App{
		orderService: service.OrderService{
			Storage: orders,
		},
		adminService: service.AdminService{
			Storage: admins,
		},
		router: mux.NewRouter(),
	}
}

const (
	orderIDParam         = "id"
	userIDParam          = "user_id"
	weightParam          = "weight"
	priceParam           = "price"
	statusParam          = "status"
	arrivalDateParam     = "arrival_date"
	arrivalDateFromParam = "arrival_date_from"
	arrivalDateToParam   = "arrival_date_to"
	expiryDateParam      = "expiry_date"
	expiryDateFromParam  = "expiry_date_from"
	expiryDateToParam    = "expiry_date_to"
	weightFromParam      = "weight_from"
	weightToParam        = "weight_to"
	priceFromParam       = "price_from"
	priceToParam         = "price_to"
	countParam           = "count"
	pageParam            = "page"

	adminUsernameParam = "admin_username"
)

func (a *App) wrapHandler(ctx context.Context, handler func(context.Context, http.ResponseWriter,
	*http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(ctx, w, r)
	}
}

// @title			PVZ API Documentation
// @version		1.0
// @description	This is a sample server for Swagger in Go.
// @host			localhost:9000
// @BasePath		/
func (a *App) Run(ctx context.Context) error {
	impl := server{
		orderService: a.orderService,
		adminService: a.adminService,
	}
	router := mux.NewRouter()
	router.Use(FieldLogger)

	authMiddleware := AuthMiddleware{
		adminService: a.adminService,
	}

	router.HandleFunc("/orders", authMiddleware.BasicAuthChecker(ctx, a.wrapHandler(ctx, impl.CreateOrder)).ServeHTTP).
		Methods(http.MethodPost)
	//	@Summary		Get Orders
	//	@Description	Fetches all orders.
	//	@Tags			orders
	//	@Accept			json
	//	@Produce		json
	//	@Success		200	{array}		models.Order
	//	@Failure		400	{object}	models.ErrorResponse
	//	@Router			/orders [get]
	router.HandleFunc("/orders", authMiddleware.BasicAuthChecker(ctx, a.wrapHandler(ctx, impl.GetOrders)).ServeHTTP).
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
	router.HandleFunc(fmt.Sprintf("/orders/{%s:[0-9]+}", orderIDParam),
		authMiddleware.BasicAuthChecker(ctx, a.wrapHandler(ctx, impl.DeleteOrder)).ServeHTTP).
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
	router.HandleFunc("/orders/process", authMiddleware.BasicAuthChecker(ctx,
		a.wrapHandler(ctx, impl.UpdateOrders)).ServeHTTP).
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
	router.HandleFunc("/admins", a.wrapHandler(ctx, impl.CreateAdmin)).
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
	router.HandleFunc(fmt.Sprintf("/admins/{%s:[a-zA-Z0-9]+}", adminUsernameParam), a.wrapHandler(ctx, impl.UpdateAdmin)).
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
	router.HandleFunc(fmt.Sprintf("/admins/{%s:[a-zA-Z0-9]+}", adminUsernameParam), a.wrapHandler(ctx, impl.DeleteAdmin)).
		Methods(http.MethodDelete)

	// Путь для отображения Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	if err := http.ListenAndServe("localhost:9000", router); err != nil {
		log.Fatal(err)
	}

	return nil
}
