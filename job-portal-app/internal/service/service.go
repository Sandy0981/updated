package service

import (
	"context"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/models"
	"job-portal-api/internal/redis"
	"job-portal-api/internal/repository"
)

type Service struct {
	UserRepo repository.UserRepo
	auth     auth.Authentication
	rdb      *redis.RedisClient
}

//go:generate mockgen -source=service.go -destination=service_mock.go -package=service
type UserService interface {
	RegisterUserService(ctx context.Context, userData models.NewUser) (models.User, error)
	UserLoginService(ctx context.Context, userData models.NewUser) (string, error)
	CreateCompanyService(ctx context.Context, companyData models.Company) (models.Company, error)
	ListCompaniesService(ctx context.Context) ([]models.Company, error)
	GetCompanyService(ctx context.Context, cid uint64) (models.Company, error)
	ListJobsForCompanyService(ctx context.Context, cid uint64) ([]models.Jobs, error)
	CreateJobPostingService(ctx context.Context, jobData models.NewJobRequest, cid uint64) (models.NewJobResponse, error)
	GetAllJobPostingsService(ctx context.Context) ([]models.Jobs, error)
	GetJobPostingByIDService(ctx context.Context, jid uint64) (models.Jobs, error)
	ApplicationProcessor(ctx context.Context, job []models.RequestJob) ([]models.RequestJob, error)
}

// NewService creates a new UserService with the provided user repository and authentication service.
// It returns a UserService and an error if the user repository is nil.
func NewService(userRepo repository.UserRepo, a auth.Authentication, rdb *redis.RedisClient) (UserService, error) {
	return &Service{
		UserRepo: userRepo,
		auth:     a,
		rdb:      rdb,
	}, nil
}
