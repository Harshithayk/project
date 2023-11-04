package services

import (
	"errors"
	"fmt"

	"project/internal/model"
	"project/internal/repository"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	r repository.Users
}

func NewService(r repository.Users) (*Service, error) {
	if r == nil {
		return nil, errors.New("db connection not given")
	}

	return &Service{
		r: r,
	}, nil

}

//go:generate mockgen -source=userService.go -destination=userservice_mock.go -package=servicesgo
type Serviceinterface interface {
	CompanyCreate(nc model.CreateCompany) (model.Company, error)
	GetAllCompanies() ([]model.Company, error)
	GetCompany(id int) (model.Company, error)
	JobCreate(nj model.CreateJob, id uint64) (model.Job, error)
	GetJobs(id int) ([]model.Job, error)
	GetAllJobs() ([]model.Job, error)
	UserSignup(nu model.UserSignup) (model.User, error)
	Userlogin(l model.UserLogin) (jwt.RegisteredClaims, error)
}

func (s *Service) UserSignup(nu model.UserSignup) (model.User, error) {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Msg("error occured in hashing password")
		return model.User{}, errors.New("hashing password failed")
	}

	user := model.User{UserName: nu.UserName, Email: nu.Email, PasswordHash: string(hashedPass)}
	cu, err := s.r.CreateUser(user)
	if err != nil {
		log.Error().Err(err).Msg("couldnot create user")
		return model.User{}, errors.New("user creation failed")
	}

	return cu, nil

}
func (s *Service) Userlogin(l model.UserLogin) (jwt.RegisteredClaims, error) {
	fu, err := s.r.FetchUserByEmail(l.Email)
	if err != nil {
		log.Error().Err(err).Msg("couldnot find user")
		return jwt.RegisteredClaims{}, errors.New("user login failed")
	}
	fmt.Println(fu)
	err = bcrypt.CompareHashAndPassword([]byte(fu.PasswordHash), []byte(l.Password))
	if err != nil {
		log.Error().Err(err).Msg("password of user incorrect")
		return jwt.RegisteredClaims{}, errors.New("user login failed")
	}
	c := jwt.RegisteredClaims{
		Issuer:    "service project",
		Subject:   strconv.FormatUint(uint64(fu.ID), 10),
		Audience:  jwt.ClaimStrings{"users"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	fmt.Println(c)

	return c, nil
}
