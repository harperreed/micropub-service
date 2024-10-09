package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harperreed/micropub-service/internal/config"
	"github.com/harperreed/micropub-service/internal/middleware"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPocketBase struct {
	mock.Mock
}

func (m *MockPocketBase) OnBeforeServe() *pocketbase.Hook {
	args := m.Called()
	return args.Get(0).(*pocketbase.Hook)
}

func (m *MockPocketBase) Start() error {
	args := m.Called()
	return args.Error(0)
}

func TestServerInitialization(t *testing.T) {
	cfg := &config.Config{
		GitRepoPath: "/tmp/test-repo",
	}

	app, err := initializeApp(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, app)

	mockApp := new(MockPocketBase)
	mockApp.On("OnBeforeServe").Return(pocketbase.NewHook())
	mockApp.On("Start").Return(nil)

	setupRoutes(mockApp)

	mockApp.AssertCalled(t, "OnBeforeServe")
	mockApp.AssertCalled(t, "Start")
}

func TestRoleAuthorizationMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test with no user
	handler := middleware.RoleAuthorization("admin")(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)
	assert.Error(t, err)
	assert.Equal(t, "You must be logged in to access this resource", err.Error())

	// Test with user but wrong role
	c.Set("user", &models.Record{Id: "user123"})
	err = handler(c)
	assert.Error(t, err)
	assert.Equal(t, "You are not authorized to access this resource", err.Error())

	// Test with correct role
	middleware.GetUserRole = func(userId string) string {
		return "admin"
	}
	err = handler(c)
	assert.NoError(t, err)
}

func TestErrorHandling(t *testing.T) {
	// Test with invalid configuration
	_, err := config.Load()
	assert.Error(t, err)

	// Test with invalid Git repository path
	cfg := &config.Config{
		GitRepoPath: "/invalid/path",
	}
	_, err = initializeApp(cfg)
	assert.Error(t, err)
}

func TestLoginFunctionality(t *testing.T) {
	e := echo.New()

	// Test successful login
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("email=test@example.com&password=password"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handleLogin(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	// Test login failure
	req = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("email=wrong@example.com&password=wrongpassword"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err = handleLogin(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
