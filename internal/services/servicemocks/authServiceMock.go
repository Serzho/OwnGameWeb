//nolint:wrapcheck
package servicemocks

import (
	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) SignIn(email, password string) (int, error) {
	args := m.Called(email, password)
	return args.Int(0), args.Error(1)
}

func (m *AuthServiceMock) SignUp(name, email, password string) error {
	args := m.Called(name, email, password)
	return args.Error(0)
}

func (m *AuthServiceMock) RecoverPassword(email string) error {
	args := m.Called(email)
	return args.Error(0) //nolint:wrapcheck
}
