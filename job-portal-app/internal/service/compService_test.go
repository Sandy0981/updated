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

func TestService_ListCompaniesService(t *testing.T) {
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
		want             []models.Company
		wantErr          bool
		mockRepoResponse func() ([]models.Company, error)
	}{
		{
			name: "error from db",
			args: args{
				context.Background(),
			},
			want:    nil,
			wantErr: true,
			mockRepoResponse: func() ([]models.Company, error) {
				return nil, errors.New("db error")
			},
		},
		{
			name: "success",
			args: args{
				context.Background(),
			},
			want: []models.Company{
				{
					Name:     "TCS",
					Location: "India",
				},
				{
					Name:     "TekSystem",
					Location: "USA",
				},
			},
			wantErr: false,
			mockRepoResponse: func() ([]models.Company, error) {
				return []models.Company{
					{
						Name:     "TCS",
						Location: "India",
					},
					{
						Name:     "TekSystem",
						Location: "USA",
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
				mockRepo.EXPECT().FetchAllCompanies(tt.args.ctx).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.ListCompaniesService(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListCompaniesService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListCompaniesService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetCompanyService(t *testing.T) {
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
		want             models.Company
		wantErr          bool
		mockRepoResponse func() (models.Company, error)
	}{
		{
			name: "error from db",
			args: args{
				context.Background(),
				1,
			},
			want:    models.Company{},
			wantErr: true,
			mockRepoResponse: func() (models.Company, error) {
				return models.Company{}, errors.New("db error")
			},
		},
		{
			name: "success",
			args: args{
				context.Background(),
				1,
			},
			want:    models.Company{Name: "TCS", Location: "INDIA"},
			wantErr: false,
			mockRepoResponse: func() (models.Company, error) {
				return models.Company{Name: "TCS", Location: "INDIA"}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().FetchCompanyByID(tt.args.ctx, tt.args.cid).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.GetCompanyService(tt.args.ctx, tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCompanyService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCompanyService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_CreateCompanyService(t *testing.T) {
	type fields struct {
		UserRepo repository.UserRepo
		auth     auth.Authentication
	}
	type args struct {
		ctx         context.Context
		companyData models.Company
	}
	tests := []struct {
		name             string
		args             args
		want             models.Company
		wantErr          bool
		mockRepoResponse func() (models.Company, error)
	}{
		{
			name: "error from db",
			args: args{
				context.Background(),
				models.Company{Name: "TestCompany", Location: "TestLocation"},
			},
			want:    models.Company{},
			wantErr: true,
			mockRepoResponse: func() (models.Company, error) {
				return models.Company{}, errors.New("db error")
			},
		},
		{
			name: "success",
			args: args{
				context.Background(),
				models.Company{Name: "TestCompany", Location: "TestLocation"},
			},
			want:    models.Company{Name: "TestCompany", Location: "TestLocation"},
			wantErr: false,
			mockRepoResponse: func() (models.Company, error) {
				return models.Company{Name: "TestCompany", Location: "TestLocation"}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().InsertCompany(tt.args.ctx, tt.args.companyData).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.CreateCompanyService(tt.args.ctx, tt.args.companyData)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCompanyService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateCompanyService() got = %v, want %v", got, tt.want)
			}
		})
	}
}
