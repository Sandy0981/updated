package repository

import (
	"context"
	"errors"
	"job-portal-api/internal/models"

	"gorm.io/gorm"
	//"job-portal/internal/models"
)

type Repo struct {
	DB *gorm.DB
}

//go:generate mockgen -source=repo.go -destination=repo_mock.go -package=repository

type UserRepo interface {
	InsertUser(ctx context.Context, userData models.User) (models.User, error)
	VerifyUserCredentials(ctx context.Context, email string) (models.User, error)

	InsertCompany(ctx context.Context, companyData models.Company) (models.Company, error)
	FetchAllCompanies(ctx context.Context) ([]models.Company, error)
	FetchCompanyByID(ctx context.Context, cid uint64) (models.Company, error)

	InsertJobPosting(ctx context.Context, jobData models.NewJobRequest) (models.NewJobResponse, error)
	FetchJobsForCompany(ctx context.Context, cid uint64) ([]models.Jobs, error)
	FetchAllJobPostings(ctx context.Context) ([]models.Jobs, error)
	FetchJobPostingByID(ctx context.Context, jid uint64) (models.Jobs, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

func NewRepository(db *gorm.DB) (UserRepo, error) {
	if db == nil {
		return nil, errors.New("db cannot be null")
	}
	return &Repo{
		DB: db,
	}, nil
}
