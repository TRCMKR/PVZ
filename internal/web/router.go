package web

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/service"

	"github.com/gorilla/mux"
)

type orderStorage interface {
	AddOrder(context.Context, models.Order)
	RemoveOrder(context.Context, int)
	UpdateOrder(context.Context, int, models.Order)
	GetByID(context.Context, int) models.Order
	GetByUserID(context.Context, int, int) []models.Order
	GetReturns(context.Context) []models.Order
	GetOrders(context.Context, map[string]string, int, int) []models.Order
	Save(context.Context) error
	Contains(context.Context, int) bool
}

type adminStorage interface {
	CreateAdmin(context.Context, models.Admin)
	GetAdminByUsername(context.Context, string) models.Admin
	UpdateAdmin(context.Context, int, models.Admin)
	DeleteAdmin(context.Context, string)
	ContainsUsername(context.Context, string) bool
	ContainsID(context.Context, int) bool
}

type App struct {
	orderService service.OrderService
	adminService service.AdminService
	router       mux.Router
}

func NewApp(ctx context.Context, orders orderStorage, admins adminStorage) *App {
	return &App{
		orderService: service.OrderService{
			Storage: orders,
			Ctx:     ctx,
		},
		adminService: service.AdminService{
			Storage: admins,
			Ctx:     ctx,
		},
		router: *mux.NewRouter(),
	}
}

const (
	orderIDParam         = "id"
	userIDParam          = "user_id"
	weightParam          = "weight"
	priceParam           = "price"
	statusParam          = "status"
	arrivalDateFromParam = "arrival_date_from"
	arrivalDateToParam   = "arrival_date_to"
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

func (a *App) Run() error {
	impl := server{
		orderService: a.orderService,
		adminService: a.adminService,
	}
	router := mux.NewRouter()
	router.Use(FieldLogger)

	authMiddleware := AuthMiddleware{
		adminService: a.adminService,
	}

	router.HandleFunc("/orders", impl.CreateOrder).
		Methods("POST").
		Handler(authMiddleware.BasicAuthChecker(http.HandlerFunc(impl.CreateOrder)))

	router.HandleFunc("/orders", impl.GetOrders).
		Methods("GET").
		Handler(authMiddleware.BasicAuthChecker(http.HandlerFunc(impl.GetOrders)))

	router.HandleFunc(fmt.Sprintf("/orders/{%s:[0-9]+}", orderIDParam), impl.DeleteOrder).
		Methods("DELETE").
		Handler(authMiddleware.BasicAuthChecker(http.HandlerFunc(impl.DeleteOrder)))

	router.HandleFunc("/orders/process", impl.UpdateOrders).
		Methods("POST").
		Handler(authMiddleware.BasicAuthChecker(http.HandlerFunc(impl.UpdateOrders)))

	router.HandleFunc("/admins", impl.CreateAdmin).
		Methods("POST")

	router.HandleFunc(fmt.Sprintf("/admins/{%s:[a-zA-Z0-9]+}", adminUsernameParam), impl.UpdateAdmin).
		Methods("POST")

	router.HandleFunc(fmt.Sprintf("/admins/{%s:[a-zA-Z0-9]+}", adminUsernameParam), impl.DeleteAdmin).
		Methods("DELETE")

	http.Handle("/", router)
	if err := http.ListenAndServe("localhost:9000", nil); err != nil {
		log.Fatal(err)
	}

	return nil
}
