package repository

import (
	"context"
	"errors"
	"project/internal/model"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) (*Repo, error) {
	if db == nil {
		return nil, errors.New("db connection not given")
	}

	return &Repo{db: db}, nil

}

//go:generate mockgen -source=userDao.go -destination=usergoDao_mock.go -package=repository

type Users interface {
	CreateUser(model.User) (model.User, error)
	FetchUserByEmail(string) (model.User, error)

	CreateCompany(model.Company) (model.Company, error)
	GetAllCompany() ([]model.Company, error)
	GetCompany(id int) (model.Company, error)

	FetchJobData(jid uint64) (model.Job, error)
	GetJobs(id int) ([]model.Job, error)
	PostJob(nj model.Job) (model.Response, error)
	GetAllJobs() ([]model.Job, error)

	GetTheJobData(jobid uint) (model.Job, error)

	EmailForgotPassword(ctx context.Context, s string) (model.User, error)
	DobForgotPassword(s string) (model.User, error)
	ChangePassword(email string, user model.User) (model.User, error)
	//Emailfround(s string) (string, error)
}

func (r *Repo) CreateUser(u model.User) (model.User, error) {
	err := r.db.Create(&u).Error
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (r *Repo) FetchUserByEmail(s string) (model.User, error) {
	var u model.User
	tx := r.db.Where("email=?", s).First(&u)
	if tx.Error != nil {
		return model.User{}, nil
	}
	return u, nil

}
func (r *Repo) EmailForgotPassword(ctx context.Context, s string) (model.User, error) {
	var u model.User
	tx := r.db.Where("email=?", s).First(&u)
	if tx.Error != nil {
		return model.User{}, nil
	}
	return u, nil

}
func (r *Repo) DobForgotPassword(s string) (model.User, error) {
	var u model.User
	tx := r.db.Where("Dob=?", s).First(&u)
	if tx.Error != nil {
		return model.User{}, nil
	}
	return u, nil

}
func (r *Repo) ChangePassword(email string, user model.User) (model.User, error) {
	var v model.User
	err := r.db.Where("email=?", email).First(&v).Error
	if err != nil {
		return model.User{}, errors.New("error taking the data")
	}
	v.PasswordHash = user.PasswordHash
	err = r.db.Save(&v).Error
	if err != nil {
		return model.User{}, errors.New("error in saving the password in db")
	}
	return v, nil
}

// func (r *Repo) Emailfround(s string) (string, error) {
// 	var u string
// 	tx := r.db.Where("email=?", s).First(&u)
// 	if tx.Error != nil {
// 		return "", nil
// 	}
// 	return u, nil

// }
