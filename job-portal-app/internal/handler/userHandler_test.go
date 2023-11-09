package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"job-portal-api/internal/middleware"
	"job-portal-api/internal/models"
	"job-portal-api/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_RegisterUser(t *testing.T) {
	type fields struct {
		service service.UserService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name               string
		setup              func() (*gin.Context, *httptest.ResponseRecorder, service.UserService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "missing trace id",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com", nil)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
		{
			name: "invalid user data",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"please provide all details"}`,
		},
		{
			name: "error while registering user",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{"username": "testuser","email": "testuser@example.com","password": "password123"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				ms.EXPECT().RegisterUserService(c.Request.Context(), gomock.Any()).Return(models.User{}, errors.New("test service error")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"test service error"}`,
		},

		{
			name: "success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{"Username": "sandeep","Email": "sandeep@gmail.com","Password": "sandy"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				userData := models.NewUser{
					Username: "sandeep",
					Email:    "sandeep@gmail.com",
					Password: "sandy",
				}

				ms.EXPECT().RegisterUserService(c.Request.Context(), gomock.Any()).Return(models.User{
					Model:        gorm.Model{ID: 1},
					Username:     userData.Username,
					Email:        userData.Email,
					PasswordHash: "hashedpassword",
				}, nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"username":"sandeep","email":"sandeep@gmail.com"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()

			h := &handler{
				service: ms,
			}
			h.RegisterUser(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_UserLogin(t *testing.T) {
	tests := []struct {
		name               string
		setup              func() (*gin.Context, *httptest.ResponseRecorder, service.UserService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "valid user login",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)

				// Create a request with valid user login data
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{
					"email": "sandeep@gmail.com",
					"password": "sandy"
				}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				// Create a mock UserService
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				// Expect the UserLoginService to be called and return a valid token
				ms.EXPECT().UserLoginService(c.Request.Context(), gomock.Any()).Return("validtoken", nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"token":"validtoken"}`,
		},
		{
			name: "missing trace id",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", nil)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
		{
			name: "error during user login",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)

				// Create a request with valid user login data
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{
					"email": "testuser@example.com",
					"password": "password123"
				}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				// Create a mock UserService
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				// Expect the UserLoginService to be called and return an error
				ms.EXPECT().UserLoginService(c.Request.Context(), gomock.Any()).Return("", errors.New("test service error")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"test service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()

			h := &handler{
				service: ms,
			}
			h.UserLogin(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}
