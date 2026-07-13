package services

import (
	"errors"

	"cms-blog-api/models"
	"cms-blog-api/repositories"
	"cms-blog-api/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Register validates user registration input, hashes password, and creates new user.
func (s *AuthService) Register(input models.RegisterInput) (*models.User, error) {
	// Check if user with same email already exists
	existingEmail, _ := s.repo.FindByEmail(input.Email)
	if existingEmail != nil {
		return nil, errors.New("email is already registered")
	}

	// Check if user with same username already exists
	existingUsername, _ := s.repo.FindByUsername(input.Username)
	if existingUsername != nil {
		return nil, errors.New("username is already taken")
	}

	// Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to encrypt password")
	}

	role := input.Role
	if role == "" {
		role = "author"
	}

	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login verifies user credentials and generates a signed JWT token.
func (s *AuthService) Login(input models.LoginInput) (string, *models.User, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	tokenString, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, errors.New("failed to generate access token")
	}

	return tokenString, user, nil
}

// GetUserByID retrieves user profile by ID.
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}
