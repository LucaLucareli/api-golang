package request

type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refreshToken" validate:"required,min=1"`
}
