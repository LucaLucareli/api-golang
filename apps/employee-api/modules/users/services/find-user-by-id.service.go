package services

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetGreeting() string {
	return "Oi"
}
