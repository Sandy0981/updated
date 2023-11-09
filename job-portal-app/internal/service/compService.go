package service

import (
	"context"

	"job-portal-api/internal/models"
)

func (s *Service) CreateCompanyService(ctx context.Context, companyData models.Company) (models.Company, error) {
	companyData, err := s.UserRepo.InsertCompany(ctx, companyData)
	if err != nil {
		return models.Company{}, err
	}
	return companyData, nil
}

func (s *Service) ListCompaniesService(ctx context.Context) ([]models.Company, error) {
	companyDetails, err := s.UserRepo.FetchAllCompanies(ctx)
	if err != nil {
		return nil, err
	}
	return companyDetails, nil
}

func (s *Service) GetCompanyService(ctx context.Context, cid uint64) (models.Company, error) {
	companyData, err := s.UserRepo.FetchCompanyByID(ctx, cid)
	if err != nil {
		return models.Company{}, err
	}
	return companyData, nil
}
