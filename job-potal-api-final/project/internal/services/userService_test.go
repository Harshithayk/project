package services

import (
	"errors"
	"project/cache"
	"project/internal/model"
	"project/internal/repository"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
)

func TestService_UserSignup(t *testing.T) {
	type args struct {
		nu model.UserSignup
	}
	tests := []struct {
		name             string
		args             args
		want             model.User
		wantErr          bool
		mockRepoResponse func() (model.User, error)
	}{{
		name:    "success case for signup",
		args:    args{model.UserSignup{UserName: "harshi", Email: "harshi@gmail.com", Password: "1223"}},
		want:    model.User{UserName: "harshi", Email: "harshi@gmail.com"},
		wantErr: false,
		mockRepoResponse: func() (model.User, error) {
			return model.User{UserName: "harshi", Email: "harshi@gmail.com"}, nil
		},
	},
		{
			name:    "success case for signup 1",
			args:    args{model.UserSignup{UserName: "harshi", Email: "harshi@gmail.com", Password: "12211111111111111365555555558521638745723458273456597460978888888888888888888888888888888888888888888888888888888888888888888888888888888888888999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999993"}},
			want:    model.User{},
			wantErr: true,
			mockRepoResponse: func() (model.User, error) {
				return model.User{}, errors.New("invalid user input for signup")
			},
		},
		{name: "error case for signup",
			args:    args{model.UserSignup{UserName: "", Email: "harshi@gmail.com", Password: "1223"}},
			want:    model.User{},
			wantErr: true,
			mockRepoResponse: func() (model.User, error) {
				return model.User{}, errors.New("invalid user input for signup")
			}},
		{name: "error case for converting password to hashpassword in signup",
			args:    args{model.UserSignup{UserName: "harshi", Email: "harshi@gmail.com"}},
			want:    model.User{},
			wantErr: true,
			mockRepoResponse: func() (model.User, error) {
				return model.User{}, errors.New("invalid user input for signup")
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUsers(mc)
			mockRepo1 := cache.NewMockCachingRadis(mc)

			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().CreateUser(gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, mockRepo1)
			got, err := s.UserSignup(tt.args.nu)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UserSignup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.UserSignup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Userlogin(t *testing.T) {
	type args struct {
		l model.UserLogin
	}
	tests := []struct {
		name             string
		args             args
		want             jwt.RegisteredClaims
		wantErr          bool
		mockRepoResponse func() (model.User, error)
	}{
		{name: "fail case for login",
			args:    args{model.UserLogin{Email: "harshimail.com", Password: ""}},
			want:    jwt.RegisteredClaims{},
			wantErr: true,
			mockRepoResponse: func() (model.User, error) {
				return model.User{Email: "harshimail.com", PasswordHash: "$2a$10$votXUqKwkXe6l5.2aVKSU.08QEPzZYuXy47OP7JuHebrZSppBlYSW"}, nil
			},
		},
		{name: "success case for login",
			args:    args{model.UserLogin{Email: "harshimail.com", Password: ""}},
			want:    jwt.RegisteredClaims{},
			wantErr: true,
			mockRepoResponse: func() (model.User, error) {
				return model.User{}, errors.New("ERROR")
			},
		},
		{name: "success case",
			args:    args{model.UserLogin{Email: "harshimail.com", Password: "hfhhfhfh"}},
			want:    jwt.RegisteredClaims{Issuer: "service project", Subject: "0", Audience: jwt.ClaimStrings{"users"}, ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), IssuedAt: jwt.NewNumericDate(time.Now())},
			wantErr: false,
			mockRepoResponse: func() (model.User, error) {
				return model.User{Email: "harshimail.com", PasswordHash: "$2a$10$votXUqKwkXe6l5.2aVKSU.08QEPzZYuXy47OP7JuHebrZSppBlYSW"}, nil
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUsers(mc)
			mockRepo1 := cache.NewMockCachingRadis(mc)

			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().FetchUserByEmail(gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, mockRepo1)
			got, err := s.Userlogin(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Userlogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Userlogin() = %v, want %v", got, tt.want)
			}
		})
	}
}
