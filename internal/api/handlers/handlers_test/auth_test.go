package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"OwnGameWeb/config"

	"OwnGameWeb/internal/api/handlers"

	"OwnGameWeb/internal/services"
	"OwnGameWeb/internal/services/servicemocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(authService *servicemocks.AuthServiceMock) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	cfg := config.Load()

	authHandler := handlers.NewAuthHandler(authService, cfg)

	router.POST("/auth/signin", authHandler.SignIn)
	router.POST("/auth/signup", authHandler.SignUp)
	router.POST("/auth/recover", authHandler.RecoverPassword)

	return router
}

func TestAuthHandler_SignIn(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	tests := []struct {
		name         string
		payload      map[string]interface{}
		mockSetup    func(mock *servicemocks.AuthServiceMock)
		expectedCode int
		checkCookie  bool
	}{
		{
			name: "successful signin",
			payload: map[string]interface{}{
				"email":    "handlers_test@example.com",
				"password": "validpass",
			},
			mockSetup: func(m *servicemocks.AuthServiceMock) {
				m.On("SignIn", "handlers_test@example.com", "validpass").Return(1, nil)
			},
			expectedCode: http.StatusOK,
			checkCookie:  true,
		},
		{
			name: "invalid credentials",
			payload: map[string]interface{}{
				"email":    "wrong@example.com",
				"password": "wrongpass",
			},
			mockSetup: func(m *servicemocks.AuthServiceMock) {
				m.On("SignIn", "wrong@example.com", "wrongpass").Return(0, services.ErrInvalidEmail)
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid request format",
			payload: map[string]interface{}{
				"email": 123,
			},
			mockSetup: func(m *servicemocks.AuthServiceMock) {
				m.On("SignIn", "", "").Return(0, services.ErrInvalidEmail)
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockAuth := new(servicemocks.AuthServiceMock)
			testCase.mockSetup(mockAuth)

			router := setupRouter(mockAuth)
			body, _ := json.Marshal(testCase.payload)
			req, _ := http.NewRequestWithContext(t.Context(), "POST", "/auth/signin", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)

			assert.Equal(t, testCase.expectedCode, writer.Code)

			if testCase.checkCookie {
				cookies := writer.Result().Cookies()
				assert.NotEmpty(t, cookies)
				assert.Equal(t, "token", cookies[0].Name)
			}

			mockAuth.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_SignUp(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	tests := []struct {
		name         string
		payload      map[string]interface{}
		mockSetup    func(*servicemocks.AuthServiceMock)
		expectedCode int
	}{
		{
			name: "successful signup",
			payload: map[string]interface{}{
				"name":     "Test User",
				"email":    "new@example.com",
				"password": "newpass123",
			},
			mockSetup: func(m *servicemocks.AuthServiceMock) {
				m.On("SignUp", "Test User", "new@example.com", "newpass123").Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "duplicate email",
			payload: map[string]interface{}{
				"name":     "Existing User",
				"email":    "exists@example.com",
				"password": "exists123",
			},
			mockSetup: func(m *servicemocks.AuthServiceMock) {
				m.On("SignUp", "Existing User", "exists@example.com", "exists123").Return(services.ErrUserAlreadyExists)
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid request body",
			payload: map[string]interface{}{
				"name": 123,
			},
			mockSetup:    func(_ *servicemocks.AuthServiceMock) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockAuth := new(servicemocks.AuthServiceMock)
			testCase.mockSetup(mockAuth)

			router := setupRouter(mockAuth)
			body, _ := json.Marshal(testCase.payload)
			req, _ := http.NewRequestWithContext(t.Context(), "POST", "/auth/signup", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedCode, w.Code)
			mockAuth.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_RecoverPassword(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	mockAuth := new(servicemocks.AuthServiceMock)
	mockAuth.On("RecoverPassword", "handlers_test@example.com").Return(services.ErrNotImplemented)

	router := setupRouter(mockAuth)
	payload := map[string]interface{}{"email": "handlers_test@example.com"}
	body, _ := json.Marshal(payload) //nolint:errchkjson

	req, _ := http.NewRequestWithContext(t.Context(), "POST", "/auth/recover", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockAuth.AssertExpectations(t)
}
