package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Admin represents an admin user in the system
// @Description Admin structure represents an administrator user with authentication details
// @Properties
type Admin struct {
	// @Description Unique ID of the admin user
	ID int `json:"id"`

	// @Description Username of the admin user
	Username string `json:"username"`

	// @Description Password of the admin user (consider securing this in production)
	Password string `json:"password"`

	// @Description Time when the admin user was created
	CreatedAt time.Time `json:"created_at"`
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// NewAdmin ...
func NewAdmin(id int, username string, password string) *Admin {
	hashedPassword, _ := hashPassword(password)

	return &Admin{
		ID:        id,
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}
}

// CheckPassword ...
func (admin *Admin) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)) == nil
}
