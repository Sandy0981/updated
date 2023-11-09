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

func TestService_RegisterUserService(t *testing.T) {
	type fields struct {
		UserRepo repository.UserRepo
		auth     auth.Authentication
	}
	type args struct {
		ctx      context.Context
		userData models.NewUser
	}
	tests := []struct {
		name             string
		args             args
		want             models.User
		wantErr          bool
		mockRepoResponse func() (models.User, error)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "validpassword",
				},
			},
			want: models.User{
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: false,
			mockRepoResponse: func() (models.User, error) {
				return models.User{
					Username: "testuser",
					Email:    "test@example.com",
				}, nil
			},
		},
		{
			name: "error from db",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "validpassword",
				},
			},
			want:    models.User{},
			wantErr: true,
			mockRepoResponse: func() (models.User, error) {
				return models.User{}, errors.New("db error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.RegisterUserService(tt.args.ctx, tt.args.userData)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUserService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterUserService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_UserLoginService(t *testing.T) {
	type fields struct {
		UserRepo repository.UserRepo
		auth     auth.Authentication
	}
	type args struct {
		ctx      context.Context
		userData models.NewUser
	}
	tests := []struct {
		name             string
		args             args
		wantToken        string
		wantErr          bool
		mockRepoResponse func() (models.User, error)
		mockAuthToken    func() (string, error)
	}{
		{
			name: "error verifying user credentials",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Email: "test@example.com",
				},
			},
			wantToken: "",
			wantErr:   true,
			mockRepoResponse: func() (models.User, error) {
				return models.User{}, errors.New("user not found")
			},
		},
		{
			name: "invalid password provided",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Email:    "test@example.com",
					Password: "invalidpassword",
				},
			},
			wantToken: "",
			wantErr:   true,
			mockRepoResponse: func() (models.User, error) {
				return models.User{
					PasswordHash: "$2a$10$3th13U5TSve35mtEP7JHwud5NiLujyZ0BK0cbV9jDb6KvCER8H.zu", // Valid hash for "validpassword"
				}, nil
			},
		},
		{
			name: "failed to generate token",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Email:    "test@example.com",
					Password: "validpassword",
				},
			},
			wantToken: "",
			wantErr:   true,
			mockRepoResponse: func() (models.User, error) {
				return models.User{
					PasswordHash: "$2a$10$xQmztwxwwg2trzNLHpuSq.crH8PojzsVG7Jh4lN96i9tgYrvodV5y", // Valid hash for "validpassword"
				}, nil
			},
			mockAuthToken: func() (string, error) {
				return "", errors.New("error in token")
			},
		},

		{
			name: "success",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Email:    "test@example.com",
					Password: "validpassword",
				},
			},
			wantToken: "testtoken",
			wantErr:   false,
			mockRepoResponse: func() (models.User, error) {
				return models.User{
					PasswordHash: "$2a$10$xQmztwxwwg2trzNLHpuSq.crH8PojzsVG7Jh4lN96i9tgYrvodV5y", // Valid hash for "validpassword"
				}, nil
			},
			mockAuthToken: func() (string, error) {
				return "testtoken", nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			mockAuth := auth.NewMockAuthentication(mc)

			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().VerifyUserCredentials(gomock.Any(), tt.args.userData.Email).Return(tt.mockRepoResponse()).AnyTimes()
			}

			if tt.mockAuthToken != nil {
				mockAuth.EXPECT().GenerateAuthToken(gomock.Any()).Return(tt.mockAuthToken()).AnyTimes()
			}

			s, _ := NewService(mockRepo, mockAuth)
			gotToken, err := s.UserLoginService(tt.args.ctx, tt.args.userData)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserLoginService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("UserLoginService() gotToken = %v, wantToken %v", gotToken, tt.wantToken)
			}
		})
	}
}
