package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"job-portal-api/internal/models"
)

func (r *Repo) FetchJobPostingByID(ctx context.Context, jid uint64) (models.Jobs, error) {
	var job models.Jobs
	result := r.DB.Where("id = ?", jid).Preload("Locations").Preload("TechnologyStacks").Preload("WorkModes").Preload("Qualifications").Preload("Shifts").Preload("JobTypes").Find(&job)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.Jobs{}, errors.New("failed to create job postings")
	}
	return job, nil
}

func (r *Repo) InsertJobPosting(ctx context.Context, jobData models.NewJobRequest) (models.NewJobResponse, error) {
	job := models.Jobs{
		Cid:              jobData.Cid,
		MinNP:            jobData.MinNP,
		MaxNP:            jobData.MaxNP,
		Budget:           jobData.Budget,
		Locations:        getLocations(jobData.Locations),
		TechnologyStacks: getTechnologyStacks(jobData.TechnologyStacks),
		WorkModes:        getWorkModes(jobData.WorkModes),
		Description:      jobData.Description,
		MinExp:           jobData.MinExp,
		MaxExp:           jobData.MaxExp,
		Qualifications:   getQualifications(jobData.Qualifications),
		Shifts:           getShifts(jobData.Shifts),
		JobTypes:         getJobTypes(jobData.JobTypes),
	}
	result := r.DB.Create(&job)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.NewJobResponse{}, errors.New("failed to insert job posting")
	}
	res := models.NewJobResponse{
		ID: job.ID,
	}
	return res, nil
}

func getLocations(locationsIds []uint) (locations []models.Locations) {
	for _, id := range locationsIds {
		locations = append(locations, models.Locations{Model: gorm.Model{ID: id}})
	}
	return
}

func getTechnologyStacks(technologyIds []uint) (technologyStacks []models.TechnologyStacks) {
	for _, id := range technologyIds {
		technologyStacks = append(technologyStacks, models.TechnologyStacks{Model: gorm.Model{ID: id}})
	}
	return
}

func getWorkModes(workModeIds []uint) (workModes []models.WorkModes) {
	for _, id := range workModeIds {
		workModes = append(workModes, models.WorkModes{Model: gorm.Model{ID: id}})
	}
	return
}

func getQualifications(qualificationIds []uint) (qualifications []models.Qualifications) {
	for _, id := range qualificationIds {
		qualifications = append(qualifications, models.Qualifications{Model: gorm.Model{ID: id}})
	}
	return
}

func getShifts(shiftIds []uint) (shifts []models.Shifts) {
	for _, id := range shiftIds {
		shifts = append(shifts, models.Shifts{Model: gorm.Model{ID: id}})
	}
	return
}

func getJobTypes(jobTypesIds []uint) (jobTypes []models.JobTypes) {
	for _, id := range jobTypesIds {
		jobTypes = append(jobTypes, models.JobTypes{Model: gorm.Model{ID: id}})
	}
	return
}

func (r *Repo) FetchAllJobPostings(ctx context.Context) ([]models.Jobs, error) {
	var jobDatas []models.Jobs
	result := r.DB.Preload("Locations").Find(&jobDatas)
	fmt.Println("DB::", jobDatas)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, errors.New("failed to fetch all job postings")
	}
	return jobDatas, nil

}

func (r *Repo) FetchJobsForCompany(ctx context.Context, cid uint64) ([]models.Jobs, error) {
	var jobData []models.Jobs
	result := r.DB.Where("cid = ?", cid).Preload("Locations").Preload("Qualification").Preload("TechnologyStack").Preload("WorkMode").Preload("Shift").Preload("JobType").Find(&jobData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, errors.New("failed to fetch jobs for company")
	}
	return jobData, nil
}
