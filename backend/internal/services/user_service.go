package services

import (
	"anime-tanyaayomi/internal/database"
	"anime-tanyaayomi/internal/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(req models.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2)`
	_, err = database.DB.Exec(query, req.Username, string(hashedPassword))
	return err
}

func (s *UserService) Login(req models.LoginRequest) (*models.User, error) {
	query := `SELECT id, username, password_hash, created_at FROM users WHERE username = $1`
	row := database.DB.QueryRow(query, req.Username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}
