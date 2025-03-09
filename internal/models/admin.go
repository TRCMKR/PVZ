package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID        int
	Username  string
	Password  string
	CreatedAt time.Time
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func NewAdmin(id int, username string, password string) *Admin {
	hashedPassword, _ := hashPassword(password)

	return &Admin{
		ID:        id,
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}
}

func (admin *Admin) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)) == nil
}
