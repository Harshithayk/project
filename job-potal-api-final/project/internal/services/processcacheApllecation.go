package services

import (
	"context"
	"encoding/json"

	"project/internal/model"
	"sync"

	"github.com/redis/go-redis/v9"
)

func (s *Service) ProccessApplication(ctx context.Context, applicationData []model.NewUserApplication) ([]model.NewUserApplication, error) {
	var wg = new(sync.WaitGroup)
	ch := make(chan model.NewUserApplication)
	var finalData []model.NewUserApplication

	for _, v := range applicationData {
		wg.Add(1)
		go func(v model.NewUserApplication) {
			defer wg.Done()
			var jobData model.Job
			val, err := s.rdb.GetCache(ctx, uint(v.ID))

			if err != nil {
				dbData, err := s.r.GetTheJobData(uint(v.ID))
				if err != nil {
					return
				}
				err = s.rdb.AddCache(ctx, uint(v.ID), dbData)
				if err != nil {
					return
				}
				jobData = dbData
			} else {
				err = json.Unmarshal([]byte(val), &jobData)
				if err == redis.Nil {
					return
				}
				if err != nil {
					return
				}
			}
			check := compareAndCheck(v, jobData)
			if check {
				ch <- v
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

func compareAndCheck(applicationData model.NewUserApplication, val model.Job) bool {
	if applicationData.Jobs.Experience < val.MinExperience {
		return false
	}
	if applicationData.Jobs.NoticePeriod < val.MinimumNoticePeriod {
		return false
	}
	var count int
	count = compareLocations(applicationData.Jobs.Location, val.Locations)
	if count == 0 {
		return false
	}

	count = compareQualifications(applicationData.Jobs.Qualifications, val.Qualifications)
	if count == 0 {
		return false
	}
	count = compareTechStack(applicationData.Jobs.WorkModeIDs, val.WorkModes)
	if count == 0 {
		return false
	}
	count = compareShifts(applicationData.Jobs.Shift, val.Shifts)
	if count == 0 {
		return false
	}

	return true
}

func compareLocations(locationsID []uint, val []model.Location) int {
	count := 0
	for _, v := range locationsID {
		for _, v1 := range val {
			if v == v1.ID {
				count++
			}
		}
	}
	return count
}

func compareQualifications(qualificationID []uint, val []model.Qualification) int {
	count := 0
	for _, v := range qualificationID {
		for _, v1 := range val {
			if v == v1.ID {
				count++
			}
		}
	}
	return count
}

func compareTechStack(stackID []uint, val []model.WorkMode) int {
	count := 0
	for _, v := range stackID {
		for _, v1 := range val {
			if v == v1.ID {
				count++
			}
		}
	}
	return count
}

func compareShifts(shiftID []uint, val []model.Shift) int {
	count := 0
	for _, v := range shiftID {
		for _, v1 := range val {
			if v == v1.ID {
				count++
			}
		}
	}
	return count
}

// check, v, err := s.compareAndCheck(v)

// if err != nil {
// 	return nil, err
// }
// if check {
// 	finalData = append(finalData, v)
// }
