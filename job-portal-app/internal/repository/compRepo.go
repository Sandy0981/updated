package repository

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"job-portal-api/internal/models"
)

func (r *Repo) InsertCompany(ctx context.Context, companyData models.Company) (models.Company, error) {
	result := r.DB.Create(&companyData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.Company{}, errors.New("failed to create company")
	}
	return companyData, nil
}

func (r *Repo) FetchAllCompanies(ctx context.Context) ([]models.Company, error) {
	var userDetails []models.Company
	result := r.DB.Find(&userDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, errors.New("failed to fetch the companies")
	}
	return userDetails, nil
}

func (r *Repo) FetchCompanyByID(ctx context.Context, cid uint64) (models.Company, error) {
	var companyData models.Company
	result := r.DB.Where("id = ?", cid).First(&companyData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.Company{}, errors.New("failed to fetch the company")
	}
	return companyData, nil
}
