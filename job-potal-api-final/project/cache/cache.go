package cache

import (
	"context"
	"encoding/json"
	"errors"
	"project/internal/model"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=cache.go -destination=cache_mock.go -package=cache
type radisLayer struct {
	rdb *redis.Client
}
type CachingRadis interface {
	AddCache(ctx context.Context, jid uint, jobdata model.Job) error
	GetCache(ctx context.Context, jid uint) (string, error)
	SetEmailCache(ctx context.Context, Email string, otp string) error
	GetEmailCache(ctx context.Context, otp string) (string, error)
}

func NewRadies(rdb *redis.Client) CachingRadis {
	return &radisLayer{
		rdb: rdb,
	}
}

func (r *radisLayer) AddCache(ctx context.Context, jid uint, jobdata model.Job) error {
	jobID := strconv.FormatUint(uint64(jid), 10)
	val, err := json.Marshal(jobdata)
	if err != nil {
		return err
	}
	err = r.rdb.Set(ctx, jobID, val, 1*time.Minute).Err()
	return err
}

func (r *radisLayer) GetCache(ctx context.Context, jid uint) (string, error) {
	jobID := strconv.FormatUint(uint64(jid), 10)
	str, err := r.rdb.Get(ctx, jobID).Result()
	return str, err
}
func (r *radisLayer) SetEmailCache(ctx context.Context, Email string, otp string) error {
	err := r.rdb.Set(ctx, Email, otp, 2*time.Minute)
	if err != nil {
		return errors.New("not set in cache")
	}
	return nil
}
func (r *radisLayer) GetEmailCache(ctx context.Context, otp string) (string, error) {
	otp, err := r.rdb.Get(ctx, otp).Result()
	if err != nil {
		return "", errors.New("not get in cache")
	}
	return otp, nil
}
