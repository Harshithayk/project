package services

import (
	"errors"
	"project/internal/model"

	"sync"

	//"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (s *Service) CompanyCreate(nc model.CreateCompany) (model.Company, error) {
	company := model.Company{CompanyName: nc.CompanyName, Adress: nc.Adress, Domain: nc.Domain}
	cu, err := s.r.CreateCompany(company)
	if err != nil {
		log.Error().Err(err).Msg("couldnot create user")
		return model.Company{}, errors.New("user creation failed")
	}

	return cu, nil
}

func (s *Service) GetAllCompanies() ([]model.Company, error) {

	AllCompanies, err := s.r.GetAllCompany()
	if err != nil {
		return nil, err
	}
	return AllCompanies, nil

}

func (s *Service) GetCompany(id int) (model.Company, error) {

	AllCompanies, err := s.r.GetCompany(id)
	if err != nil {
		return model.Company{}, err
	}
	return AllCompanies, nil

}

func (s *Service) GetJobs(id int) ([]model.Job, error) {
	AllCompanies, err := s.r.GetJobs(id)
	if err != nil {
		return nil, errors.New("job retreval failed")
	}
	return AllCompanies, nil
}

func (s *Service) GetAllJobs() ([]model.Job, error) {

	AllJobs, err := s.r.GetAllJobs()
	if err != nil {
		return nil, err
	}
	return AllJobs, nil

}

func (s *Service) JobCreate(newJob model.NewJobRequest, id uint64) (model.Response, error) {

	app := model.Job{
		CompanyId:           id,
		JobTitle:            newJob.JobTitle,
		Salary:              newJob.Salary,
		MinimumNoticePeriod: newJob.MinimumNoticePeriod,
		MaximumNoticePeriod: newJob.MaximumNoticePeriod,
		Budget:              newJob.Budget,
		JobDescription:      newJob.JobDescription,
		MinExperience:       newJob.MinExperience,
	}
	for _, v := range newJob.QualificationIDs {
		tempData := model.Qualification{
			Model: gorm.Model{
				ID: v,
			},
		}
		app.Qualifications = append(app.Qualifications, tempData)
	}
	for _, v := range newJob.LocationIDs {
		tempData := model.Location{
			Model: gorm.Model{
				ID: v,
			},
		}
		app.Locations = append(app.Locations, tempData)
	}
	for _, v := range newJob.SkillIDs {
		tempData := model.Skill{
			Model: gorm.Model{
				ID: v,
			},
		}
		app.Skills = append(app.Skills, tempData)
	}
	for _, v := range newJob.WorkModeIDs {
		tempData := model.WorkMode{
			Model: gorm.Model{
				ID: v,
			},
		}
		app.WorkModes = append(app.WorkModes, tempData)
	}
	for _, v := range newJob.ShiftIDs {
		tempData := model.Shift{
			Model: gorm.Model{
				ID: v,
			},
		}
		app.Shifts = append(app.Shifts, tempData)
	}
	for _, v := range newJob.JobTypeIDs {
		tempData := model.JobType{
			Model: gorm.Model{
				ID: v,
			},
		}
		app.JobTypes = append(app.JobTypes, tempData)
	}
	jobData, err := s.r.PostJob(app)
	if err != nil {
		return model.Response{}, err
	}
	return jobData, nil
}

// func (s *Service) ProcessJobApplications(appData []model.NewUserApplication) ([]model.NewUserApplication, error) {
// 	wg := new(sync.WaitGroup)
// 	var finalApplications []model.NewUserApplication
// 	for _, v := range appData {
// 		wg.Add(1)
// 		go func(v model.NewUserApplication) {
// 			defer wg.Done()
// 			check, v, err := s.Compare(v)
// 			if err != nil {
// 				return
// 			}
// 			if check {
// 				finalApplications = append(finalApplications, v)
// 			}

//			}(v)
//		}
//		wg.Wait()
//		return finalApplications, nil
//	}
func (s *Service) ProcessJobApplications(appData []model.NewUserApplication) ([]model.NewUserApplication, error) {
	var wg = new(sync.WaitGroup)
	ch := make(chan model.NewUserApplication)
	var finalData []model.NewUserApplication

	for _, v := range appData {
		wg.Add(1)
		go func(v model.NewUserApplication) {
			defer wg.Done()

			val, err := s.r.FetchJobData(uint64(v.ID))
			if err != nil {
				return
			}
			check, value, err := s.Compare(v, val)
			if err != nil {
				return
			}
			if check {
				ch <- value
			}
		}(v)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for v := range ch {
		finalData = append(finalData, v)
	}

	return finalData, nil
}

func (s *Service) Compare(appData model.NewUserApplication, jobData model.Job) (bool, model.NewUserApplication, error) {
	matchingConditions := 0
	totalConditions := 8
	if appData.Jobs.Experience >= jobData.MinExperience {
		matchingConditions++
	}

	if appData.Jobs.NoticePeriod >= jobData.MinimumNoticePeriod {
		matchingConditions++
	}

	for _, v := range appData.Jobs.WorkModeIDs {
		for _, v1 := range jobData.WorkModes {
			if v == v1.ID {
				matchingConditions++
				break
			}
		}
	}

	for _, v := range appData.Jobs.JobTypeIDs {
		for _, v1 := range jobData.JobTypes {
			if v == v1.ID {
				matchingConditions++
				break
			}
		}
	}

	for _, v := range appData.Jobs.Location {
		for _, v1 := range jobData.Locations {
			if v == v1.ID {
				matchingConditions++
				break
			}
		}
	}

	for _, v := range appData.Jobs.Qualifications {
		for _, v1 := range jobData.Qualifications {
			if v == v1.ID {
				matchingConditions++
				break
			}
		}
	}

	for _, v := range appData.Jobs.Skills {
		for _, v1 := range jobData.Skills {
			if v == v1.ID {
				matchingConditions++
				break
			}
		}
	}

	for _, v := range appData.Jobs.Shift {
		for _, v1 := range jobData.Shifts {
			if v == v1.ID {
				matchingConditions++
				break
			}
		}
	}

	if matchingConditions*2 >= totalConditions {
		return true, appData, nil
	}

	return false, model.NewUserApplication{}, nil
}

// func (s *Service) Compare(appData model.NewUserApplication) (bool, model.NewUserApplication, error) {
// 	jobData, err := s.r.FetchJobData(appData.ID)
// 	if err != nil {
// 		return false, model.NewUserApplication{}, err
// 	}

// 	matchingConditions := 0
// 	totalConditions := 8

// 	if appData.Jobs.Experience >= jobData.MinExperience {
// 		matchingConditions++
// 	}

// 	if appData.Jobs.NoticePeriod >= jobData.MinimumNoticePeriod {
// 		matchingConditions++
// 	}

// 	for _, v := range appData.Jobs.WorkModeIDs {
// 		for _, v1 := range jobData.WorkModes {
// 			if v == v1.ID {
// 				matchingConditions++
// 				break
// 			}
// 		}
// 	}

// 	for _, v := range appData.Jobs.JobTypeIDs {
// 		for _, v1 := range jobData.JobTypes {
// 			if v == v1.ID {
// 				matchingConditions++
// 				break
// 			}
// 		}
// 	}

// 	for _, v := range appData.Jobs.Location {
// 		for _, v1 := range jobData.Locations {
// 			if v == v1.ID {
// 				matchingConditions++
// 				break
// 			}
// 		}
// 	}

// 	for _, v := range appData.Jobs.Qualifications {
// 		for _, v1 := range jobData.Qualifications {
// 			if v == v1.ID {
// 				matchingConditions++
// 				break
// 			}
// 		}
// 	}

// 	for _, v := range appData.Jobs.Skills {
// 		for _, v1 := range jobData.Skills {
// 			if v == v1.ID {
// 				matchingConditions++
// 				break
// 			}
// 		}
// 	}

// 	for _, v := range appData.Jobs.Shift {
// 		for _, v1 := range jobData.Shifts {
// 			if v == v1.ID {
// 				matchingConditions++
// 				break
// 			}
// 		}
// 	}

// 	if matchingConditions*2 >= totalConditions {
// 		return true, appData, nil
// 	}

// 	return false, model.NewUserApplication{}, nil
// }
