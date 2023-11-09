package service

import (
	"context"
	"errors"
	"go.uber.org/mock/gomock"
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
		jobData models.Jobs
		cid     uint64
	}
	tests := []struct {
		name             string
		args             args
		want             models.Jobs
		wantErr          bool
		mockRepoResponse func() (models.Jobs, error)
	}{
		{
			name: "error from db",
			args: args{
				context.Background(),
				models.Jobs{},
				1,
			},
			want:    models.Jobs{},
			wantErr: true,
			mockRepoResponse: func() (models.Jobs, error) {
				return models.Jobs{}, errors.New("db error")
			},
		},
		{
			name: "success",
			args: args{
				context.Background(),
				models.Jobs{
					Company: models.Company{Name: "TCS"},
					Cid:     1,
				},
				1,
			},
			want: models.Jobs{
				Company: models.Company{Name: "TCS"},
				Cid:     1,
			},
			wantErr: false,
			mockRepoResponse: func() (models.Jobs, error) {
				return models.Jobs{
					Company: models.Company{Name: "TCS"},
					Cid:     1,
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
