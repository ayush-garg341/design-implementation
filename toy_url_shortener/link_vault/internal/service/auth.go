package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/linkvault/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// Use a secure environment variable for this in production
var jwtKey = []byte("16b05f36-349e-41bd-9be9-1b3b38da8016")

// Claims struct with standard claims
type Claims struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthResponse struct {
	User  store.User `json:"user"`
	Token string     `json:"token"`
}

type UserService struct {
	store store.UserStore
}

func NewUserService(store *store.PostgresUserStore) *UserService {
	return &UserService{store}
}

func (svc *UserService) Create(ctx context.Context, name, email, password string) (*store.User, error) {
	var user store.User
	user.Name = name
	user.Email = email
	hpwd, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	user.Password = hpwd

	u, err := svc.store.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (svc *UserService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	user, err := svc.store.Login(ctx, email)
	if err != nil {
		return nil, err
	}

	err = comparePassword(user.Password, password)
	if err != nil {
		return nil, err
	}

	jwtToken, err := createJwt(email, user.ID)
	if err != nil {
		return nil, err
	}

	authResponse := &AuthResponse{
		Token: jwtToken,
		User:  *user,
	}
	return authResponse, nil

}

func hashPassword(text string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	return string(bytes), err
}

func comparePassword(storedHash, text string) error {
	plainPassword := []byte(text)
	storedHashBytes := []byte(storedHash)
	err := bcrypt.CompareHashAndPassword(storedHashBytes, plainPassword)
	return err
}

func createJwt(username string, user_id string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: username,
		UserID:   user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)

}
