package user

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sinasadeghi83/aut-grader/internal/api/config"
)

type UserService struct {
	repo UserRepo
	cfg  *config.AppConfig
}

func NewUserService(repo UserRepo, cfg *config.AppConfig) *UserService {
	return &UserService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *UserService) Login(username, password string) (*User, string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil || !s.CheckPassword(password, user.Password) {
		return nil, "", ErrInvalidCreds
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(s.cfg.SecretKey))

	return user, token, err
}

func (s *UserService) CheckPassword(password, hash string) bool {
	return password == hash
}

func (s *UserService) CheckToken(tokenStr string) (*User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignature
		}
		return []byte(s.cfg.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return nil, ErrExpiredToken
	}

	user, err := s.repo.FindById(uint(claims["id"].(float64)))
	if err != nil {
		return nil, ErrInvalidCreds
	}

	return user, nil
}
