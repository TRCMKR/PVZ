package web

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"gitlab.ozon.dev/alexplay1224/homework/internal/service"
)

var (
	errUnauthorized    = errors.New("unauthorized")
	errInvalidEncoding = errors.New("invalid encoding")
	errInvalidFormat   = errors.New("invalid format")
	errNoSuchUser      = errors.New("no such user")
	errWrongPassword   = errors.New("wrong password")
)

func FieldLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost || request.Method == http.MethodPut || request.Method == http.MethodDelete {
			body, err := io.ReadAll(request.Body)
			if err != nil {
				log.Println("Error reading request body:", err)
				http.Error(writer, "can't read body", http.StatusInternalServerError)

				return
			}
			log.Printf("%s, %s, %s\n", request.Method, request.URL.Path, body)

			request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		handler.ServeHTTP(writer, request)
	})
}

type AuthMiddleware struct {
	adminService service.AdminService
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

func (a *AuthMiddleware) BasicAuthChecker(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		credsStr, err := a.parseHeader(request)
		if err != nil {
			log.Println("Error parsing auth header:", err)
		}

		decoded, err := base64.StdEncoding.DecodeString(credsStr)
		if err != nil {
			http.Error(writer, errInvalidEncoding.Error(), http.StatusUnauthorized)

			return
		}

		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			http.Error(writer, errInvalidFormat.Error(), http.StatusUnauthorized)

			return
		}

		username, password := creds[0], creds[1]

		admin, err := a.adminService.GetAdminByUsername(username)
		if err != nil {
			http.Error(writer, errNoSuchUser.Error(), http.StatusUnauthorized)

			return
		}

		if !admin.CheckPassword(password) {
			http.Error(writer, errWrongPassword.Error(), http.StatusUnauthorized)

			return
		}

		handler.ServeHTTP(writer, request)
	})
}
