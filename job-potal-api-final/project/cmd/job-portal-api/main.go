package main

import (
	"fmt"
	"net/http"
	"project/cache"
	"project/config"
	"project/internal/auth"
	"project/internal/database"
	"project/internal/handlers"
	"project/internal/repository"
	"project/internal/services"
	"project/redis"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	err := startApp()
	if err != nil {
		log.Panic().Err(err).Send()
	}
}
func startApp() error {
	cfg := config.GetConfig()
	log.Info().Interface("cfg", cfg).Msg("config")

	log.Info().Msg("started main")

	privatePEM := []byte(cfg.KeysPubPri.Private)
	// if err != nil {
	// 	return fmt.Errorf("cannot find file private.pem %w", err)
	// }
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("cannot convert byte to key %w", err)
	}
	publicPEM := []byte(cfg.KeysPubPri.Public)
	//	publicPEM, err := os.ReadFile(`pubkey.pem`)
	// if err != nil {
	// 	return fmt.Errorf("cannot find file pubkey.pem %w", err)
	// }
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return fmt.Errorf("cannot convert byte to key %w", err)
	}
	a, err := auth.NewAuth(privateKey, publicKey)
	if err != nil {
		return fmt.Errorf("cannot create auth instance %w", err)
	}

	db, err := database.Connection()
	if err != nil {
		return err
	}
	repo, err := repository.NewRepo(db)
	if err != nil {
		return err
	}
	rdb := redis.ConnectRedis()

	redisLayer := cache.NewRadies(rdb)
	se, err := services.NewService(repo, redisLayer)

	if err != nil {
		return err
	}

	api := http.Server{ //server config and settimngs
		Addr:    fmt.Sprintf(":%s", cfg.AppConfig.Port),
		Handler: handlers.Api(a, se),
	}
	api.ListenAndServe()

	return nil

}
