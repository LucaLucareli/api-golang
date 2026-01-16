package services

import (
	"auth-api/modules/auth/dto/io"
	"context"
	"shared/auth"
)

type RefreshTokenService struct {
	authService *auth.AuthService
}

func NewRefreshTokenService(authService *auth.AuthService) *RefreshTokenService {
	return &RefreshTokenService{
		authService: authService,
	}
}

func (s *RefreshTokenService) RefreshTokenService(
	ctx context.Context,
	input io.RefreshTokenInputDTO,
) (*io.RefreshTokenOutputDTO, error) {

	response, err := s.authService.RefreshToken(
		ctx,
		input.RefreshToken,
	)
	if err != nil {
		return nil, err
	}

	return &io.RefreshTokenOutputDTO{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}
