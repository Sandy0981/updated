package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/middleware"
	"job-portal-api/internal/models"
	"job-portal-api/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_GetJobPostingByID(t *testing.T) {
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
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
		{
			name: "missing jwt claims",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Unauthorized"}`,
		},
		{
			name: "invalid job id",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Params = append(c.Params, gin.Param{Key: "jobID", Value: "1ab"})

				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Bad Request"}`,
		},
		{
			name: "error while fetching jobs from service",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "123"})
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				ms.EXPECT().GetJobPostingByIDService(c.Request.Context(), gomock.Any()).Return(models.Jobs{}, errors.New("test service error")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Bad Request"}`,
		},
		{
			name: "error from db",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "jobID", Value: "123"})
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				ms.EXPECT().GetJobPostingByIDService(c.Request.Context(), gomock.Any()).Return(models.Jobs{}, errors.New("test service error")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"test service error"}`,
		},
		{
			name: "success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "jobID", Value: "123"})
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				ms.EXPECT().GetJobPostingByIDService(c.Request.Context(), gomock.Any()).Return(models.Jobs{}, nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"cid":0,"minNp":0,"maxNp":0,"budget":0,"Locations":null,"TechnologyStacks":null,"WorkModes":null,"description":"","minExp":0,"maxExp":0,"Qualifications":null,"Shifts":null,"JobTypes":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()

			h := &handler{
				service: ms,
			}
			h.GetJobPostingByID(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_GetAllJobPostings(t *testing.T) {
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
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
		{
			name: "missing jwt claims",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Unauthorized"}`,
		},
		{
			name: "error while fetching jobs from service",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				ms.EXPECT().GetAllJobPostingsService(c.Request.Context()).Return([]models.Jobs{}, errors.New("test service error")).AnyTimes()

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
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				ms.EXPECT().GetAllJobPostingsService(c.Request.Context()).Return([]models.Jobs{}, nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()

			h := &handler{
				service: ms,
			}
			h.GetAllJobPostings(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_ListJobsForCompany(t *testing.T) {
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
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
		{
			name: "missing jwt claims",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Unauthorized"}`,
		},
		{
			name: "invalid company ID",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "companyID", Value: "abc"})

				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Bad Request"}`,
		},
		{
			name: "error from db",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "jobID", Value: "123"})
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				ms.EXPECT().GetJobPostingByIDService(c.Request.Context(), gomock.Any()).Return(models.Jobs{}, errors.New("test service error")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Bad Request"}`,
		},
		{
			name: "error from service",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "companyID", Value: "123"})
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)
				expectedJobData := []models.Jobs{}
				ms.EXPECT().ListJobsForCompanyService(c.Request.Context(), gomock.Any()).Return(expectedJobData, errors.New("error from service")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"error from service"}`,
		},
		{
			name: "success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "companyID", Value: "123"})
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)
				expectedJobData := []models.Jobs{}
				ms.EXPECT().ListJobsForCompanyService(c.Request.Context(), gomock.Any()).Return(expectedJobData, nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()

			h := &handler{
				service: ms,
			}
			h.ListJobsForCompany(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_CreateJobPosting(t *testing.T) {
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
			name: "missing jwt claims",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Unauthorized"}`,
		},
		{
			name: "invalid company ID",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "companyID", Value: "invalid"})

				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Bad Request"}`,
		},
		{
			name: "invalid request body",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{"invalid`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "companyID", Value: "123"})

				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"please provide valid name, location"}`,
		},
		{
			name: "error while creating job posting",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{"name":"Software Engineer","salary":"$100,000","location":"San Francisco"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "companyID", Value: "123"})
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				ms.EXPECT().CreateJobPostingService(c.Request.Context(), gomock.Any(), gomock.Any()).Return(models.NewJobResponse{}, errors.New("test service error")).AnyTimes()

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
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "companyID", Value: "123"})
				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				ms.EXPECT().CreateJobPostingService(c.Request.Context(), gomock.Any(), gomock.Any()).Return(models.NewJobResponse{}, nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"ID":0}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()

			h := &handler{
				service: ms,
			}
			h.CreateJobPosting(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_ProcessJobApplication(t *testing.T) {
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
			name: "invalid request",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com", bytes.NewBufferString(`invalid`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"please provide valid application data"}`,
		},
		{
			name: "error from ApplicationProcessor",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				// Setup
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com", bytes.NewBufferString(`[]`))

				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				httpRequest = httpRequest.WithContext(ctx)

				c.Request = httpRequest

				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				ms.EXPECT().ApplicationProcessor(gomock.Any(), gomock.Any()).Return(nil, errors.New("test error")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"application processing error"}`,
		},
		{
			name: "valid response from ApplicationProcessor",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				// Setup
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com", bytes.NewBufferString(`[]`))

				// Add trace ID to context
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "123")
				httpRequest = httpRequest.WithContext(ctx)

				// Set request
				c.Request = httpRequest

				mc := gomock.NewController(t)
				ms := service.NewMockUserService(mc)

				// Expect ApplicationProcessor to return a valid response
				ms.EXPECT().ApplicationProcessor(gomock.Any(), gomock.Any()).Return([]models.RequestJob{}, nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()

			h := &handler{
				service: ms,
			}
			h.ProcessJobApplication(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}
