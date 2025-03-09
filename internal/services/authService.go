package services

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) SignIn(_, _ string) error {
	return nil
}

func (s *AuthService) SignUp(_, _, _ string) error {
	return nil
}

func (s *AuthService) RecoverPassword(_ string) error {
	return nil
}
