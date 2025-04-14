package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	order_Handler "gitlab.ozon.dev/alexplay1224/homework/internal/web/http/order"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"
	"gitlab.ozon.dev/alexplay1224/homework/internal/service/auditlogger"
)

var (
	errUnauthorized    = errors.New("unauthorized")
	errInvalidEncoding = errors.New("invalid encoding")
	errInvalidFormat   = errors.New("invalid format")
	errNoSuchUser      = errors.New("no such user")
	errWrongPassword   = errors.New("wrong password")
)

// FieldLogger logs fields of passed request body
func FieldLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Println("Error reading r body:", err)
				http.Error(w, "can't read body", http.StatusInternalServerError)

				return
			}
			log.Printf("%s, %s, %s\n", r.Method, r.URL.Path, body)

			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		handler.ServeHTTP(w, r)
	})
}

// AuthMiddleware is a structure for auth middleware
type AuthMiddleware struct {
	adminService admin.Service
}

func (a *AuthMiddleware) parseHeader(request *http.Request) (string, error) {
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errUnauthorized
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", errUnauthorized
	}

	return parts[1], nil
}

// BasicAuthChecker is a function that checks request for basic auth
func (a *AuthMiddleware) BasicAuthChecker(ctx context.Context, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credsStr, err := a.parseHeader(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)

			return
		}

		decoded, err := base64.StdEncoding.DecodeString(credsStr)
		if err != nil {
			http.Error(w, errInvalidEncoding.Error(), http.StatusUnauthorized)

			return
		}

		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			http.Error(w, errInvalidFormat.Error(), http.StatusUnauthorized)

			return
		}

		username, password := creds[0], creds[1]

		admin, err := a.adminService.GetAdminByUsername(ctx, username)
		if err != nil {
			http.Error(w, errNoSuchUser.Error(), http.StatusUnauthorized)

			return
		}

		if !admin.CheckPassword(password) {
			http.Error(w, errWrongPassword.Error(), http.StatusUnauthorized)

			return
		}

		handler.ServeHTTP(w, r)
	})
}

// AuditLoggerMiddleware is a structure for audit logger middleware
type AuditLoggerMiddleware struct {
	adminService       admin.Service
	auditLoggerService auditlogger.Service
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	rw.body.Write(b)

	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

type requestBody struct {
	ID int `json:"id"`
}

// AuditLogger is a function that logs all requests and responses
func (a *AuditLoggerMiddleware) AuditLogger(ctx context.Context, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request requestBody
		body, err := io.ReadAll(r.Body)
		if err != nil {
			request.ID, _ = strconv.Atoi(mux.Vars(r)[order_Handler.OrderIDParam])
		} else {
			err = json.Unmarshal(body, &request)
			if err != nil || request.ID == 0 {
				request.ID = -1
			}
		}

		r.Body = io.NopCloser(bytes.NewReader(body))

		username, _, _ := r.BasicAuth()

		rw := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		handler.ServeHTTP(rw, r)

		select {
		case <-ctx.Done():
			return
		default:
			someAdmin, _ := a.adminService.GetAdminByUsername(ctx, username)
			responseText := strings.TrimSpace(rw.body.String())
			currentLog := *models.NewLog(request.ID, someAdmin.ID, responseText, r.URL.Path, r.Method, rw.statusCode)
			a.auditLoggerService.CreateLog(ctx, currentLog)
		}
	})
}
