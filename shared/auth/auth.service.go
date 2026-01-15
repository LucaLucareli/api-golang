package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shared/helpers"
	repository "shared/repository/users"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	AccessSecret      string
	RefreshSecret     string
	AccessExpiryHours int
	RefreshExpiryDays int
	UserRepo          repository.UsersRepository
}

func (s *AuthService) Login(
	ctx context.Context,
	document string,
	password string,
) (*AuthResponse, error) {

	userModel, err := s.UserRepo.FindUserToLogin(ctx, document)

	if err != nil {
		return nil, fmt.Errorf("erro ao acessar o banco de dados: %w", err)
	}

	if userModel == nil || !helpers.CheckPassword(password, userModel.Password) {
		return nil, errors.New("usuário ou senha inválidos")
	}

	user := User{
		ID:           userModel.ID.String(),
		Document:     document,
		Name:         userModel.Name,
		AccessGroups: userModel.AccessGroupIds,
	}

	accessToken, err := s.GenerateAccessToken(&user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) ValidateAccessToken(tokenStr string) (*TokenInfo, error) {

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.AccessSecret), nil
		},
	)
	if err != nil || !token.Valid {
		return nil, errors.New("token inválido")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || claims.TokenType != "access" {
		return nil, errors.New("claims inválidas")
	}

	return &TokenInfo{
		ID:           claims.Subject,
		Document:     claims.Document,
		Name:         claims.Name,
		AccessGroups: claims.AccessGroupIds,
	}, nil
}

func (s *AuthService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*AuthResponse, error) {

	token, err := jwt.ParseWithClaims(
		refreshToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.RefreshSecret), nil
		},
	)
	if err != nil || !token.Valid {
		return nil, errors.New("refresh token inválido ou expirado")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || claims.TokenType != "refresh" {
		return nil, errors.New("refresh token inválido")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errors.New("id de usuário inválido no token")
	}

	userModel, err := s.UserRepo.FindByID(ctx, userID)

	if err != nil {
		return nil, fmt.Errorf("erro ao acessar o banco de dados: %w", err)
	}
	if userModel == nil {
		return nil, errors.New("usuário não encontrado")
	}

	user := User{
		ID:           userModel.ID.String(),
		Document:     userModel.Document,
		Name:         userModel.Name,
		AccessGroups: userModel.AccessGroupIds,
	}

	accessToken, err := s.GenerateAccessToken(&user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.GenerateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) GenerateAccessToken(user *User) (string, error) {

	now := time.Now()
	exp := now.Add(time.Duration(s.AccessExpiryHours) * time.Hour)

	claims := Claims{
		Document:       user.Document,
		Name:           user.Name,
		AccessGroupIds: user.AccessGroups,
		TokenType:      "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.AccessSecret))
}

func (s *AuthService) GenerateRefreshToken(user *User) (string, error) {

	now := time.Now()
	exp := now.Add(time.Duration(s.RefreshExpiryDays) * 24 * time.Hour)

	claims := Claims{
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.RefreshSecret))
}
