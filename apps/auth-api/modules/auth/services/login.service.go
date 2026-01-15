package services

import (
	"auth-api/modules/auth/dto/io"
	"context"
	"shared/auth"
)

type LoginService struct {
	authService *auth.AuthService
}

func NewLoginService(authService *auth.AuthService) *LoginService {
	return &LoginService{
		authService: authService,
	}
}

func (s *LoginService) LoginService(
	ctx context.Context,
	input io.LoginInputDTO,
) (*io.LoginOutputDTO, error) {

	response, err := s.authService.Login(
		ctx,
		input.Document,
		input.Password,
	)
	if err != nil {
		return nil, err
	}

	return &io.LoginOutputDTO{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}
