package repository

import (
	"errors"
	"project/internal/model"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=jobDao.go -destination=jobDao_mock.go -package=repository
func (r *Repo) CreateCompany(u model.Company) (model.Company, error) {
	err := r.db.Create(&u).Error
	if err != nil {
		return model.Company{}, err
	}
	return u, nil
}

func (r *Repo) GetAllCompany() ([]model.Company, error) {
	var s []model.Company
	err := r.db.Find(&s).Error
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (r *Repo) GetCompany(id int) (model.Company, error) {
	var m model.Company

	tx := r.db.Where("id = ?", id)
	err := tx.Find(&m).Error
	if err != nil {
		return model.Company{}, err
	}
	return m, nil

}

//	func (r *Repo) CreateJob(j model.Job) (model.Job, error) {
//		err := r.db.Create(&j).Error
//		if err != nil {
//			return model.Job{}, err
//		}
//		return j, nil
//	}
func (r *Repo) PostJob(nj model.Job) (model.Response, error) {

	res := r.db.Create(&nj).Error
	if res != nil {
		log.Info().Err(res).Send()
		return model.Response{}, errors.New("job creation failed")
	}
	return model.Response{ID: uint64(nj.ID)}, nil
}

func (r *Repo) GetJobs(id int) ([]model.Job, error) {
	// var m []model.Job

	// tx := r.db.Where("id = ?", id)
	// err := tx.Find(&m).Error
	// if err != nil {
	// 	return nil, err
	// }
	// return m, nil
	var jobData []model.Job

	// Preload related data using GORM's Preload method
	result := r.db.Preload("Comp").
		Preload("Locations").
		Preload("WorkModes").
		Preload("Qualifications").
		Preload("Shifts").
		Where("id = ?", id).
		Find(&jobData)

	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, result.Error
	}

	return jobData, nil

}

func (r *Repo) GetAllJobs() ([]model.Job, error) {
	var s []model.Job
	err := r.db.Preload("Comp").Preload("Qualifications").Preload("Shifts").Preload("Locations").Preload("JobTypes").Preload("WorkModes").Find(&s).Error
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *Repo) ApplyJob(id uint64) (model.Job, error) {
	var m model.Job
	tx := r.db.Preload("Qualifications").Preload("Locations").Preload("Shifts").Preload("JobTypes").Preload("WorkModes").Where("id = ?", id)
	err := tx.Find(&m).Error
	if err != nil {
		return model.Job{}, err
	}
	return m, nil

}

func (r *Repo) FetchJobData(jid uint64) (model.Job, error) {
	var j model.Job
	result := r.db.Preload("Comp").
		Preload("Locations").
		Preload("Skills").
		Preload("Qualifications").
		Preload("Shifts").
		Where("id = ?", jid).
		Find(&j)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return model.Job{}, result.Error
	}

	return j, nil
}
func (r *Repo) GetTheJobData(jobid uint) (model.Job, error) {
	var jobData model.Job

	// Preload related data using GORM's Preload method
	result := r.db.Preload("Comp").
		Preload("Locations").
		Preload("WorkModes").
		Preload("Qualifications").
		Preload("Shifts").
		Where("id = ?", jobid).
		Find(&jobData)

	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return model.Job{}, result.Error
	}

	return jobData, nil
}
