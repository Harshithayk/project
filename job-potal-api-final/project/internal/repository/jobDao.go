package repository

import "project/internal/model"

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

func (r *Repo) CreateJob(j model.Job) (model.Job, error) {
	err := r.db.Create(&j).Error
	if err != nil {
		return model.Job{}, err
	}
	return j, nil
}

func (r *Repo) GetJobs(id int) ([]model.Job, error) {
	var m []model.Job

	tx := r.db.Where("uid = ?", id)
	err := tx.Find(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil

}

func (r *Repo) GetAllJobs() ([]model.Job, error) {
	var s []model.Job
	err := r.db.Find(&s).Error
	if err != nil {
		return nil, err
	}

	return s, nil
}
