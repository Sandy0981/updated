package service

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"job-portal-api/internal/models"
	"sync"
)

func (s *Service) GetJobPostingByIDService(ctx context.Context, jid uint64) (models.Jobs, error) {
	if jid < uint64(0) {
		return models.Jobs{}, errors.New("id cannot be 0")
	}
	jobData, err := s.UserRepo.FetchJobPostingByID(ctx, jid)
	if err != nil {
		return models.Jobs{}, err
	}
	return jobData, nil
}

func (s *Service) GetAllJobPostingsService(ctx context.Context) ([]models.Jobs, error) {
	jobDatas, err := s.UserRepo.FetchAllJobPostings(ctx)
	if err != nil {
		return nil, err
	}
	return jobDatas, nil
}

func (s *Service) CreateJobPostingService(ctx context.Context, jobData models.NewJobRequest, cid uint64) (models.NewJobResponse, error) {
	jobData.Cid = uint(cid)
	jobDatas, err := s.UserRepo.InsertJobPosting(ctx, jobData)
	if err != nil {
		return models.NewJobResponse{}, err
	}
	return jobDatas, nil
}

func (s *Service) ListJobsForCompanyService(ctx context.Context, cid uint64) ([]models.Jobs, error) {
	jobData, err := s.UserRepo.FetchJobsForCompany(ctx, cid)
	if err != nil {
		return nil, err
	}
	return jobData, nil
}

func (s *Service) ApplicationProcessor(ctx context.Context, jobApplications []models.RequestJob) ([]models.RequestJob, error) {
	ch := make(chan models.RequestJob)
	var wg sync.WaitGroup

	for _, application := range jobApplications {
		wg.Add(1)

		go func(application models.RequestJob) {
			defer wg.Done()
			job, err := s.UserRepo.FetchJobPostingByID(ctx, application.Jid)
			if err != nil {
				return
			}
			if validateJobApplication(application, job) {
				ch <- application
			}
		}(application)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := []models.RequestJob{}
	for app := range ch {
		result = append(result, app)
	}

	return result, nil
}

func validateJobApplication(application models.RequestJob, job models.Jobs) bool {

	// Check Notice Period
	if application.NoticePeriod < job.MinNP || application.NoticePeriod > job.MaxNP {
		log.Print("Notice period criteria not matched :-", application.Name)
		return false
	}
	// Compare Budget
	if application.Budget > job.Budget {
		log.Print("Budget criteria not matched :-", application.Name)
		return false
	}
	// Check for matching location IDs
	foundLocation := false
	for _, v := range application.LocationsIDs {
		for _, c := range job.Locations {
			if v == c.ID {
				foundLocation = true
				break
			}
		}
	}
	if !foundLocation {
		log.Print("location criteria not matched :-", application.Name)
		return false
	}

	requiredTechCount := len(job.TechnologyStacks) / 2 // 50% of required technologies

	// Count how many technologies in the application match the required ones
	matchingTechCount := 0
	for _, techID := range application.TechnologyStackIDs {
		for _, reqTechID := range job.TechnologyStacks {
			if techID == reqTechID.ID {
				matchingTechCount++
				break
			}
		}
	}
	if matchingTechCount < requiredTechCount {
		log.Print("Skills criteria not matched :-", application.Name)
		return false
	}
	// Check Work Mode
	workModeSatisfied := false
	for _, appWorkMode := range application.WorkModeIDs {
		for _, jobWorkMode := range job.WorkModes {
			if appWorkMode == jobWorkMode.ID {
				workModeSatisfied = true
				break
			}
		}
		if workModeSatisfied {
			break
		}
	}

	if !workModeSatisfied {
		log.Print("workMode criteria not matched :-", application.Name)
		return false

	}

	// Check Experience
	if application.Experience < job.MinExp || application.Experience > job.MaxExp {
		log.Print("experience criteria not match :-", application.Name)
		return false
	}

	// Check for at least 2 qualifications
	// Check Qualifications
	requiredQualifications := job.Qualifications // List of required qualifications

	qualificationSatisfied := false
	for _, reqQualificationID := range requiredQualifications {
		for _, appQualificationID := range application.QualificationIDs {
			if reqQualificationID.ID == appQualificationID {
				qualificationSatisfied = true
				break
			}
		}
		if qualificationSatisfied {
			break
		}
	}

	if !qualificationSatisfied {
		log.Print("qualification criteria not match :-", application.Name)
		return false
	}
	// Check Shifts
	shiftSatisfied := false
	for _, appShift := range application.ShiftIDs {
		for _, jobShift := range job.Shifts {
			if appShift == jobShift.ID {
				shiftSatisfied = true
				break
			}
		}
		if shiftSatisfied {
			break
		}
	}
	if !shiftSatisfied {
		log.Print("shift criteria not matched :-", application.Name)
		return false
	}
	// Check Job Type
	jobTypeSatisfied := false
	for _, appJobType := range application.JobTypeIDs {
		for _, jobJobType := range job.JobTypes {
			if appJobType == jobJobType.ID {
				jobTypeSatisfied = true
				break
			}
		}
		if jobTypeSatisfied {
			break
		}
	}
	if !jobTypeSatisfied {
		log.Print("jobType criteria not match :-", application.Name)
		return false
	}
	return true
}
