package services

import (
	"context"
	"math/rand"
	//"crypto/rand"
	//"math/big"

	"errors"
	"fmt"

	"net/smtp"

	"project/cache"
	"project/internal/model"
	"project/internal/repository"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	r   repository.Users
	rdb cache.CachingRadis
}

var otpCode string

//go:generate mockgen -source=userService.go -destination=userservice_mock.go -package=servicesgo
type Serviceinterface interface {
	CompanyCreate(nc model.CreateCompany) (model.Company, error)
	GetAllCompanies() ([]model.Company, error)
	GetCompany(id int) (model.Company, error)
	//JobCreate(nj model.CreateJob, id uint64) (model.Job, error)
	JobCreate(newJob model.NewJobRequest, id uint64) (model.Response, error)
	GetJobs(id int) ([]model.Job, error)
	GetAllJobs() ([]model.Job, error)
	ProcessJobApplications(appData []model.NewUserApplication) ([]model.NewUserApplication, error)

	UserSignup(nu model.UserSignup) (model.User, error)
	Userlogin(l model.UserLogin) (jwt.RegisteredClaims, error)

	// JobProcess(idc uint64, idj uint64, jr []model.JobRequest) ([]model.Job, error)
	ProccessApplication(ctx context.Context, applicationData []model.NewUserApplication) ([]model.NewUserApplication, error)

	Forgetpassword(ctx context.Context, forotp model.ForgotPass) (bool, string, error)
	Newpassword(ctx context.Context, newpass model.Newpassword) (string, error)
}

func NewService(r repository.Users, rdb cache.CachingRadis) (*Service, error) {
	if r == nil {
		return nil, errors.New("db connection not given")
	}

	return &Service{
		r:   r,
		rdb: rdb,
	}, nil

}

func (s *Service) UserSignup(nu model.UserSignup) (model.User, error) {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Msg("error occured in hashing password")
		return model.User{}, errors.New("hashing password failed")
	}

	user := model.User{UserName: nu.UserName, Email: nu.Email, PasswordHash: string(hashedPass), Dob: nu.Dob}
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

	return c, nil
}

func (s *Service) Forgetpassword(ctx context.Context, f model.ForgotPass) (bool, string, error) {
	e, err := s.r.EmailForgotPassword(ctx, f.Email)

	if err != nil {
		log.Error().Err(err).Msg("could not find email")
		return false, "", errors.New("cannot get the email")
	}
	// _, err = s.r.DobForgotPassword(f.Dob)
	// if err != nil {
	// 	log.Error().Err(err).Msg("could not find dob")
	// 	return errors.New("cannot get the dob")
	// }

	otpCode = generateRandomOTP(5)

	from := "ykharshitha@gmail.com"
	password := "pkem svxo xgbh tyhg"
	// Recipient's email address
	to := e
	// SMTP server details
	smtpServer := "smtp.gmail.com"
	smtpPort := 587

	message := []byte(fmt.Sprintf("Subject: Test Email\n\notp:" + otpCode))
	// calling the set method to set email,otp
	err = s.rdb.SetEmailCache(ctx, e.Email, otpCode)
	// if err != nil {
	// 	fmt.Println("Error storing the email,otp:", err)
	// 	return false, "", errors.New("error in cache")
	// }
	// Authentication information
	auth := smtp.PlainAuth("", from, password, smtpServer)
	// SMTP connection
	smtpAddr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)
	err = smtp.SendMail(smtpAddr, auth, from, []string{to.Email}, message)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return false, "", errors.New("cannot send  email")
	}
	fmt.Println("Email sent successfully!")

	return true, otpCode, nil
}

// func generateOTP() (string, error) {
// 	// Define the length of the OTP code
// 	otpLength := 6

// 	// Create a slice to store the random digits
// 	otp := make([]byte, otpLength)

// 	// Generate random digits using crypto/rand
// 	for i := 0; i < otpLength; i++ {
// 		randomDigit, err := rand.Int(rand.Reader, big.NewInt(10))
// 		if err != nil {
// 			return "", err
// 		}
// 		otp[i] = byte(randomDigit.Int64())
// 	}

// 	// Convert the byte slice to a string
// 	otpCode := fmt.Sprintf("%06d", otp)

//		return otpCode, nil
//	}
func generateRandomOTP(length int) string {
	rand.Seed(time.Now().UnixNano())

	// Define the characters allowed in the OTP
	otpChars := "0123456789abcdefghijklmnopqrstABCDEFGHIHIJKLMNOPQRSTUVWXYZ"

	// Generate the OTP
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = otpChars[rand.Intn(len(otpChars))]
	}

	return string(otp)
}

func (s *Service) Newpassword(ctx context.Context, np model.Newpassword) (string, error) {
	// otpData, err := s.rdb.GetEmailCache(ctx, np.Email)
	// if err != nil {
	// 	return "", errors.New("error in geting cache opt data in cache")
	// }

	if np.Otp == otpCode {
		if np.Newpassword == np.ConfirmPassword {
			//newpassuser, err := s.r.FetchUserByEmail(np.Email)
			// if err != nil {
			// 	return "", errors.New("error in fetching data by email in server")
			// }
			hashedPass, err := bcrypt.GenerateFromPassword([]byte(np.ConfirmPassword), bcrypt.DefaultCost)
			if err != nil {
				return "", errors.New("error in hassing password for new password")
			}
			user := model.User{PasswordHash: string(hashedPass)}

			_, err = s.r.ChangePassword(np.Email, user)
			if err != nil {
				log.Error().Err(err).Msg("couldnot update user")
				return "", errors.New("user updating failed")
			}

		} else {
			return "", errors.New("password is mismatch")
		}
	} else {
		return "", errors.New("otp is not valid")
	}
	return "success fully set the new password", nil
}
