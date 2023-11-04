package services

import (
	"errors"
	"project/internal/model"
	"project/internal/repository"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestService_CompanyCreate(t *testing.T) {
	type args struct {
		nc model.CreateCompany
	}
	tests := []struct {
		name             string
		args             args
		want             model.Company
		mockRepoResponse func() (model.Company, error)
		wantErr          bool
	}{
		{
			name: "check all field is present",
			args: args{
				model.CreateCompany{CompanyName: "tek", Adress: "bangalore", Domain: "java"},
			},
			want:    model.Company{CompanyName: "tek", Adress: "bangalore", Domain: "java"},
			wantErr: false,
			mockRepoResponse: func() (model.Company, error) {
				return model.Company{CompanyName: "tek", Adress: "bangalore", Domain: "java"}, nil
			},
		},
		{
			name: "check feild are not present",
			args: args{
				model.CreateCompany{CompanyName: "", Adress: "bangalore", Domain: "java"},
			},
			want:    model.Company{},
			wantErr: true,
			mockRepoResponse: func() (model.Company, error) {
				return model.Company{}, errors.New("invalid input")
			},
		},
		{
			name: "check feild are not present",
			args: args{
				model.CreateCompany{CompanyName: "", Adress: "bangalore", Domain: "java"},
			},
			want:    model.Company{},
			wantErr: false,
			mockRepoResponse: func() (model.Company, error) {
				return model.Company{}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUsers(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().CreateCompany(gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo)
			got, err := s.CompanyCreate(tt.args.nc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CompanyCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.CompanyCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetAllCompanies(t *testing.T) {
	tests := []struct {
		name             string
		want             []model.Company
		wantErr          bool
		mockRepoResponse func() ([]model.Company, error)
	}{
		{name: "pass case",
			want:    []model.Company{{CompanyName: "tek", Adress: "bangalore", Domain: "python"}, {CompanyName: "tek", Adress: "bangalore", Domain: "python"}},
			wantErr: false,
			mockRepoResponse: func() ([]model.Company, error) {
				return []model.Company{{CompanyName: "tek", Adress: "bangalore", Domain: "python"}, {CompanyName: "tek", Adress: "bangalore", Domain: "python"}},
					nil
			},
		},
		{name: "pass case",
			want:    nil,
			wantErr: true,
			mockRepoResponse: func() ([]model.Company, error) {
				return nil, errors.New("enpty")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUsers(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().GetAllCompany().Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo)
			got, err := s.GetAllCompanies()
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetAllCompanies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetAllCompanies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetCompany(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name             string
		args             args
		want             model.Company
		wantErr          bool
		mockRepoResponse func() (model.Company, error)
	}{
		{name: "gitting the id",
			args:    args{id: 23},
			want:    model.Company{CompanyName: "tek", Adress: "banglore", Domain: "python"},
			wantErr: false,
			mockRepoResponse: func() (model.Company, error) {
				return model.Company{CompanyName: "tek", Adress: "banglore", Domain: "python"}, nil
			},
		},
		{name: "id is not present",
			args:    args{},
			want:    model.Company{},
			wantErr: true,
			mockRepoResponse: func() (model.Company, error) {
				return model.Company{}, errors.New("id is not present")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUsers(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().GetCompany(gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo)
			got, err := s.GetCompany(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetCompany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetCompany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_JobCreate(t *testing.T) {
	type args struct {
		nj model.CreateJob
		id uint64
	}
	tests := []struct {
		name             string
		args             args
		want             model.Job
		wantErr          bool
		mockRepoResponse func() (model.Job, error)
	}{
		{name: "CRAETING THE JOB",
			args:    args{nj: model.CreateJob{JobSalary: "hgh", JobTitle: "gffe"}, id: 56},
			want:    model.Job{JobTitle: "info", JobSalary: "787"},
			wantErr: false,
			mockRepoResponse: func() (model.Job, error) {
				return model.Job{JobTitle: "info", JobSalary: "787"}, nil
			},
		},
		{name: "CRAETING THE JOB",
			args:    args{nj: model.CreateJob{JobSalary: "", JobTitle: "gffe"}, id: 98},
			want:    model.Job{},
			wantErr: true,
			mockRepoResponse: func() (model.Job, error) {
				return model.Job{}, errors.New("FALS")

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUsers(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().CreateJob(gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo)

			got, err := s.JobCreate(tt.args.nj, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.JobCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.JobCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetJobs(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name             string
		args             args
		want             []model.Job
		wantErr          bool
		mockRepoResponse func() ([]model.Job, error)
	}{
		{
			name: "pass",
			args: args{
				id: 2,
			},
			want:    []model.Job{{JobTitle: "java developer", JobSalary: "4536348"}, {JobTitle: "go developer", JobSalary: "45348"}},
			wantErr: false,
			mockRepoResponse: func() ([]model.Job, error) {
				return []model.Job{{JobTitle: "java developer", JobSalary: "4536348"}, {JobTitle: "go developer", JobSalary: "45348"}}, nil
			},
		},
		{
			name: "pass",
			args: args{
				id: 2,
			},
			want:    nil,
			wantErr: true,
			mockRepoResponse: func() ([]model.Job, error) {
				return nil, errors.New("empty slice")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUsers(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().GetJobs(gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo)
			got, err := s.GetJobs(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetJobs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetAllJobs(t *testing.T) {
	tests := []struct {
		name             string
		want             []model.Job
		wantErr          bool
		mockRepoResponse func() ([]model.Job, error)
	}{
		{
			name:    "pass",
			want:    []model.Job{{JobTitle: "java developer", JobSalary: "4536348"}, {JobTitle: "go developer", JobSalary: "45348"}},
			wantErr: false,
			mockRepoResponse: func() ([]model.Job, error) {
				return []model.Job{{JobTitle: "java developer", JobSalary: "4536348"}, {JobTitle: "go developer", JobSalary: "45348"}}, nil
			},
		},
		{
			name:    "fail",
			want:    nil,
			wantErr: true,
			mockRepoResponse: func() ([]model.Job, error) {
				return nil, errors.New("empty job")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUsers(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().GetAllJobs().Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo)
			got, err := s.GetAllJobs()
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetAllJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetAllJobs() = %v, want %v", got, tt.want)
			}
		})
	}
}
