package web

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"
)

var (
	errUnauthorized    = errors.New("unauthorized")
	errInvalidEncoding = errors.New("invalid encoding")
	errInvalidFormat   = errors.New("invalid format")
	errNoSuchUser      = errors.New("no such user")
	errWrongPassword   = errors.New("wrong password")
)

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
