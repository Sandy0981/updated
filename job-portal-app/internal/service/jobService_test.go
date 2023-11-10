package service

import (
	"context"
	"errors"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/models"
	"job-portal-api/internal/repository"
	"reflect"
	"testing"
)

func TestService_GetJobPostingByIDService(t *testing.T) {
	type fields struct {
		UserRepo repository.UserRepo
		auth     auth.Authentication
	}
	type args struct {
		ctx context.Context
		jid uint64
	}
	tests := []struct {
		name             string
		args             args
		want             models.NewJobRequest
		wantErr          bool
		mockRepoResponse func() (models.NewJobRequest, error)
	}{
		{
			name: "error from db",
			args: args{
				ctx: context.Background(),
				jid: 15,
			},
			want:    models.NewJobRequest{},
			wantErr: true,
			mockRepoResponse: func() (models.NewJobRequest, error) {
				return models.NewJobRequest{}, errors.New("test error")
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				jid: 15,
			},
			want: models.NewJobRequest{
				Cid: 1,
			},
			wantErr: false,
			mockRepoResponse: func() (models.NewJobRequest, error) {
				return models.NewJobRequest{
					Cid: 1,
				}, nil
			},
		},
		{
			name: "invalid job id",
			args: args{
				ctx: context.Background(),
				jid: 5,
			},
			want:    models.NewJobRequest{},
			wantErr: true,
			mockRepoResponse: func() (models.NewJobRequest, error) {
				return models.NewJobRequest{}, errors.New("invalid job id")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().FetchJobPostingByID(tt.args.ctx, tt.args.jid).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.GetJobPostingByIDService(tt.args.ctx, tt.args.jid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJobPostingByIDService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetJobPostingByIDService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetAllJobPostingsService(t *testing.T) {
	type fields struct {
		UserRepo repository.UserRepo
		auth     auth.Authentication
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name             string
		args             args
		want             []models.Jobs
		wantErr          bool
		mockRepoResponse func() ([]models.Jobs, error)
	}{
		{
			name: "error from db",
			args: args{
				context.Background(),
			},
			want:    nil,
			wantErr: true,
			mockRepoResponse: func() ([]models.Jobs, error) {
				return nil, errors.New("db error")
			},
		},
		{
			name: "success",
			args: args{
				context.Background(),
			},
			want: []models.Jobs{
				{
					Cid: 1,
				},
				{
					Cid: 2,
				},
			},
			wantErr: false,
			mockRepoResponse: func() ([]models.Jobs, error) {
				return []models.Jobs{
					{
						Cid: 1,
					},
					{
						Cid: 2,
					},
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().FetchAllJobPostings(tt.args.ctx).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.GetAllJobPostingsService(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllJobPostingsService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllJobPostingsService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_CreateJobPostingService(t *testing.T) {
	type fields struct {
		UserRepo repository.UserRepo
		auth     auth.Authentication
	}
	type args struct {
		ctx     context.Context
		jobData models.NewJobRequest
		cid     uint64
	}
	tests := []struct {
		name             string
		args             args
		want             models.NewJobResponse
		wantErr          bool
		mockRepoResponse func() (models.NewJobResponse, error)
	}{
		{
			name: "error from db",
			args: args{
				context.Background(),
				models.NewJobRequest{},
				1,
			},
			want:    models.NewJobResponse{},
			wantErr: true,
			mockRepoResponse: func() (models.NewJobResponse, error) {
				return models.NewJobResponse{}, errors.New("db error")
			},
		},
		{
			name: "success",
			args: args{
				context.Background(),
				models.NewJobRequest{
					Cid: 1,
				},
				1,
			},
			want: models.NewJobResponse{
				ID: 1,
			},
			wantErr: false,
			mockRepoResponse: func() (models.NewJobResponse, error) {
				return models.NewJobResponse{
					ID: 1,
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().InsertJobPosting(tt.args.ctx, gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.CreateJobPostingService(tt.args.ctx, tt.args.jobData, tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateJobPostingService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateJobPostingService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ListJobsForCompanyService(t *testing.T) {
	type fields struct {
		UserRepo repository.UserRepo
		auth     auth.Authentication
	}
	type args struct {
		ctx context.Context
		cid uint64
	}
	tests := []struct {
		name             string
		args             args
		want             []models.Jobs
		wantErr          bool
		mockRepoResponse func() ([]models.Jobs, error)
	}{
		{
			name: "error from db",
			args: args{
				context.Background(),
				1,
			},
			want:    nil,
			wantErr: true,
			mockRepoResponse: func() ([]models.Jobs, error) {
				return nil, errors.New("db error")
			},
		},
		{
			name: "success",
			args: args{
				context.Background(),
				1,
			},
			want: []models.Jobs{
				{
					Company: models.Company{
						Name: "TCS",
					},
					Cid: 1,
				},
				{
					Company: models.Company{Name: "TCS"},
					Cid:     1,
				},
			},
			wantErr: false,
			mockRepoResponse: func() ([]models.Jobs, error) {
				return []models.Jobs{
					{
						Company: models.Company{
							Name: "TCS",
						},
						Cid: 1,
					},
					{
						Company: models.Company{Name: "TCS"},
						Cid:     1,
					},
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().FetchJobsForCompany(tt.args.ctx, tt.args.cid).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.ListJobsForCompanyService(tt.args.ctx, tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListJobsForCompanyService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListJobsForCompanyService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ApplicationProcessor(t *testing.T) {
	type args struct {
		ctx       context.Context
		applicant []models.RequestJob
	}
	tests := []struct {
		name    string
		args    args
		want    []models.RequestJob
		wantErr bool
		setup   func(mockRepo *repository.MockUserRepo)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				applicant: []models.RequestJob{
					{
						Name:               "John",
						Jid:                0,
						NoticePeriod:       15,
						Budget:             400,
						LocationsIDs:       []uint{1},
						TechnologyStackIDs: []uint{1, 2},
						WorkModeIDs:        []uint{1, 2},
						Description:        "job sample",
						Experience:         8,
						QualificationIDs:   []uint{},
						ShiftIDs:           []uint{1},
						JobTypeIDs:         []uint{1},
					},
					{
						Name:               "Sam",
						Jid:                1,
						NoticePeriod:       30,
						Budget:             4000,
						LocationsIDs:       []uint{1},
						TechnologyStackIDs: []uint{1, 2},
						WorkModeIDs:        []uint{1, 2},
						Description:        "job sample2",
						Experience:         2,
						QualificationIDs:   []uint{1, 2},
						ShiftIDs:           []uint{1},
						JobTypeIDs:         []uint{1},
					},
				},
			},
			want: []models.RequestJob{
				{
					Name:               "Sam",
					Jid:                1,
					NoticePeriod:       30,
					Budget:             4000,
					LocationsIDs:       []uint{1},
					TechnologyStackIDs: []uint{1, 2},
					WorkModeIDs:        []uint{1, 2},
					Description:        "job sample2",
					Experience:         2,
					QualificationIDs:   []uint{1, 2},
					ShiftIDs:           []uint{1},
					JobTypeIDs:         []uint{1},
				},
			},
			wantErr: false,
			setup: func(mockRepo *repository.MockUserRepo) {
				mockRepo.EXPECT().FetchJobPostingByID(gomock.Any(), uint64(0)).Return(models.Jobs{}, errors.New("test error")).Times(1)
				mockRepo.EXPECT().FetchJobPostingByID(gomock.Any(), uint64(1)).Return(models.Jobs{
					Model:  gorm.Model{ID: 1},
					MinNP:  20,
					MaxNP:  50,
					Budget: 600000,
					Locations: []models.Locations{
						{Model: gorm.Model{ID: 1}},
					},
					TechnologyStacks: []models.TechnologyStacks{
						{Model: gorm.Model{ID: 1}},
					},
					WorkModes: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MinExp: 1,
					MaxExp: 3,
					Qualifications: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shifts: []models.Shifts{
						{Model: gorm.Model{ID: 1}},
					},
					JobTypes: []models.JobTypes{
						{Model: gorm.Model{ID: 1}},
					},
				}, nil).Times(1)
			},
		},

		{
			name: "Validation Check",
			args: args{
				ctx: context.Background(),
				applicant: []models.RequestJob{
					{
						Name:               "A",
						Jid:                1,
						NoticePeriod:       15,
						Budget:             400,
						LocationsIDs:       []uint{1},
						TechnologyStackIDs: []uint{1, 2},
						WorkModeIDs:        []uint{1, 2},
						Experience:         8,
						QualificationIDs:   []uint{},
						ShiftIDs:           []uint{1},
						JobTypeIDs:         []uint{1},
					},
					{
						Name:               "B",
						Jid:                2,
						NoticePeriod:       1,
						Budget:             8000000,
						LocationsIDs:       []uint{1},
						TechnologyStackIDs: []uint{1, 2},
						WorkModeIDs:        []uint{1, 2},
						Experience:         8,
						QualificationIDs:   []uint{},
						ShiftIDs:           []uint{1},
						JobTypeIDs:         []uint{1},
					},
					{
						Name:               "C",
						Jid:                3,
						NoticePeriod:       1,
						Budget:             4000,
						LocationsIDs:       []uint{8, 9},
						TechnologyStackIDs: []uint{1, 2},
						WorkModeIDs:        []uint{1, 2},
						Experience:         8,
						QualificationIDs:   []uint{},
						ShiftIDs:           []uint{1},
						JobTypeIDs:         []uint{1},
					},
					{
						Name:               "D",
						Jid:                4,
						NoticePeriod:       2,
						Budget:             400,
						LocationsIDs:       []uint{1},
						TechnologyStackIDs: []uint{100, 101, 102},
						WorkModeIDs:        []uint{1, 2},
						Experience:         8,
						QualificationIDs:   []uint{},
						ShiftIDs:           []uint{1},
						JobTypeIDs:         []uint{1},
					},
					{
						Name:               "E",
						Jid:                5,
						NoticePeriod:       2,
						Budget:             400,
						LocationsIDs:       []uint{1},
						TechnologyStackIDs: []uint{1, 2},
						WorkModeIDs:        []uint{5},
						Experience:         8,
						QualificationIDs:   []uint{},
						ShiftIDs:           []uint{1},
						JobTypeIDs:         []uint{1},
					},
					{
						Name:               "F",
						Jid:                6,
						NoticePeriod:       2,
						Budget:             400,
						LocationsIDs:       []uint{1},
						TechnologyStackIDs: []uint{1, 2},
						WorkModeIDs:        []uint{1},
						Experience:         8,
						QualificationIDs:   []uint{},
						ShiftIDs:           []uint{1},
						JobTypeIDs:         []uint{1},
					},
					{
						Name:               "G",
						Jid:                7,
						NoticePeriod:       2,
						Budget:             400,
						LocationsIDs:       []uint{1},
						TechnologyStackIDs: []uint{1, 2},
						WorkModeIDs:        []uint{1},
						Experience:         3,
						QualificationIDs:   []uint{5, 9},
						ShiftIDs:           []uint{1},
						JobTypeIDs:         []uint{1},
					},
					{
						Name:               "H",
						Jid:                8,
						NoticePeriod:       2,
						Budget:             400,
						LocationsIDs:       []uint{1},
						TechnologyStackIDs: []uint{1, 2},
						WorkModeIDs:        []uint{1},
						Experience:         3,
						QualificationIDs:   []uint{1},
						ShiftIDs:           []uint{9},
						JobTypeIDs:         []uint{1},
					},
					{
						Name:               "I",
						Jid:                9,
						NoticePeriod:       2,
						Budget:             400,
						LocationsIDs:       []uint{1},
						TechnologyStackIDs: []uint{1, 2},
						WorkModeIDs:        []uint{1},
						Experience:         3,
						QualificationIDs:   []uint{1},
						ShiftIDs:           []uint{1},
						JobTypeIDs:         []uint{10},
					},
				},
			},
			want:    []models.RequestJob{},
			wantErr: false,
			setup: func(mockRepo *repository.MockUserRepo) {
				//mockRepo.EXPECT().GetJobById(gomock.Any(), uint(0)).Return(models.Job{}, errors.New("test error")).Times(1)
				mockRepo.EXPECT().FetchJobPostingByID(gomock.Any(), uint64(1)).Return(models.Jobs{
					Model:  gorm.Model{ID: 1},
					Cid:    1,
					MinNP:  0,
					MaxNP:  2,
					Budget: 600000,
					Locations: []models.Locations{
						{Model: gorm.Model{ID: 1}},
					},
					TechnologyStacks: []models.TechnologyStacks{
						{Model: gorm.Model{ID: 1}}, {Model: gorm.Model{ID: 2}},
					},
					WorkModes: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualifications: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shifts: []models.Shifts{
						{Model: gorm.Model{ID: 1}},
					},
					JobTypes: []models.JobTypes{
						{Model: gorm.Model{ID: 1}},
					},
				}, nil).Times(1)
				mockRepo.EXPECT().FetchJobPostingByID(gomock.Any(), uint64(2)).Return(models.Jobs{
					Model:  gorm.Model{ID: 1},
					Cid:    1,
					MinNP:  0,
					MaxNP:  2,
					Budget: 600000,
					Locations: []models.Locations{
						{Model: gorm.Model{ID: 1}},
					},
					TechnologyStacks: []models.TechnologyStacks{
						{Model: gorm.Model{ID: 1}},
					},
					WorkModes: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualifications: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shifts: []models.Shifts{
						{Model: gorm.Model{ID: 1}},
					},
					JobTypes: []models.JobTypes{
						{Model: gorm.Model{ID: 1}},
					},
				}, nil).Times(1)
				mockRepo.EXPECT().FetchJobPostingByID(gomock.Any(), uint64(3)).Return(models.Jobs{
					Model:  gorm.Model{ID: 1},
					Cid:    1,
					MinNP:  0,
					MaxNP:  2,
					Budget: 600000,
					Locations: []models.Locations{
						{Model: gorm.Model{ID: 1}},
					},
					TechnologyStacks: []models.TechnologyStacks{
						{Model: gorm.Model{ID: 1}},
					},
					WorkModes: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualifications: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shifts: []models.Shifts{
						{Model: gorm.Model{ID: 1}},
					},
					JobTypes: []models.JobTypes{
						{Model: gorm.Model{ID: 1}},
					},
				}, nil).Times(1)
				mockRepo.EXPECT().FetchJobPostingByID(gomock.Any(), uint64(4)).Return(models.Jobs{
					Model:  gorm.Model{ID: 1},
					Cid:    1,
					MinNP:  0,
					MaxNP:  2,
					Budget: 600000,
					Locations: []models.Locations{
						{Model: gorm.Model{ID: 1}},
					},
					TechnologyStacks: []models.TechnologyStacks{
						{Model: gorm.Model{ID: 1}}, {Model: gorm.Model{ID: 2}},
					},
					WorkModes: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualifications: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shifts: []models.Shifts{
						{Model: gorm.Model{ID: 1}},
					},
					JobTypes: []models.JobTypes{
						{Model: gorm.Model{ID: 1}},
					},
				}, nil).Times(1)
				mockRepo.EXPECT().FetchJobPostingByID(gomock.Any(), uint64(5)).Return(models.Jobs{
					Model:  gorm.Model{ID: 1},
					Cid:    1,
					MinNP:  0,
					MaxNP:  2,
					Budget: 600000,
					Locations: []models.Locations{
						{Model: gorm.Model{ID: 1}},
					},
					TechnologyStacks: []models.TechnologyStacks{
						{Model: gorm.Model{ID: 1}}, {Model: gorm.Model{ID: 2}},
					},
					WorkModes: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualifications: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shifts: []models.Shifts{
						{Model: gorm.Model{ID: 1}},
					},
					JobTypes: []models.JobTypes{
						{Model: gorm.Model{ID: 1}},
					},
				}, nil).Times(1)
				mockRepo.EXPECT().FetchJobPostingByID(gomock.Any(), uint64(6)).Return(models.Jobs{
					Model:  gorm.Model{ID: 1},
					Cid:    1,
					MinNP:  0,
					MaxNP:  2,
					Budget: 600000,
					Locations: []models.Locations{
						{Model: gorm.Model{ID: 1}},
					},
					TechnologyStacks: []models.TechnologyStacks{
						{Model: gorm.Model{ID: 1}}, {Model: gorm.Model{ID: 2}},
					},
					WorkModes: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MinExp: 2,
					MaxExp: 3,
					Qualifications: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shifts: []models.Shifts{
						{Model: gorm.Model{ID: 1}},
					},
					JobTypes: []models.JobTypes{
						{Model: gorm.Model{ID: 1}},
					},
				}, nil).Times(1)
				mockRepo.EXPECT().FetchJobPostingByID(gomock.Any(), uint64(7)).Return(models.Jobs{
					Model:  gorm.Model{ID: 1},
					Cid:    1,
					MinNP:  0,
					MaxNP:  2,
					Budget: 600000,
					Locations: []models.Locations{
						{Model: gorm.Model{ID: 1}},
					},
					TechnologyStacks: []models.TechnologyStacks{
						{Model: gorm.Model{ID: 1}}, {Model: gorm.Model{ID: 2}},
					},
					WorkModes: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MinExp: 2,
					MaxExp: 7,
					Qualifications: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shifts: []models.Shifts{
						{Model: gorm.Model{ID: 1}},
					},
					JobTypes: []models.JobTypes{
						{Model: gorm.Model{ID: 1}},
					},
				}, nil).Times(1)
				mockRepo.EXPECT().FetchJobPostingByID(gomock.Any(), uint64(8)).Return(models.Jobs{
					Model:  gorm.Model{ID: 1},
					Cid:    1,
					MinNP:  0,
					MaxNP:  2,
					Budget: 600000,
					Locations: []models.Locations{
						{Model: gorm.Model{ID: 1}},
					},
					TechnologyStacks: []models.TechnologyStacks{
						{Model: gorm.Model{ID: 1}}, {Model: gorm.Model{ID: 2}},
					},
					WorkModes: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MinExp: 2,
					MaxExp: 7,
					Qualifications: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shifts: []models.Shifts{
						{Model: gorm.Model{ID: 1}},
					},
					JobTypes: []models.JobTypes{
						{Model: gorm.Model{ID: 1}},
					},
				}, nil).Times(1)
				mockRepo.EXPECT().FetchJobPostingByID(gomock.Any(), uint64(9)).Return(models.Jobs{
					Model:  gorm.Model{ID: 1},
					Cid:    1,
					MinNP:  0,
					MaxNP:  2,
					Budget: 600000,
					Locations: []models.Locations{
						{Model: gorm.Model{ID: 1}},
					},
					TechnologyStacks: []models.TechnologyStacks{
						{Model: gorm.Model{ID: 1}}, {Model: gorm.Model{ID: 2}},
					},
					WorkModes: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MinExp: 2,
					MaxExp: 7,
					Qualifications: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shifts: []models.Shifts{
						{Model: gorm.Model{ID: 1}},
					},
					JobTypes: []models.JobTypes{
						{Model: gorm.Model{ID: 1}},
					},
				}, nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			tt.setup(mockRepo)
			s := &Service{
				UserRepo: mockRepo,
			}
			got, err := s.ApplicationProcessor(tt.args.ctx, tt.args.applicant)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplicationProcessor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplicationProcessor() got = %v, want %v", got, tt.want)
			}
		})
	}
}
