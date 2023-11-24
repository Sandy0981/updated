package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"job-portal-api/internal/config"
	"job-portal-api/internal/redis"
	"net/http"
	"os"
	"os/signal"
	"time"

	"job-portal-api/internal/auth"
	"job-portal-api/internal/database"
	"job-portal-api/internal/handler"
	"job-portal-api/internal/repository"
	"job-portal-api/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	err := StartApp()
	if err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Msg("Hello this is our app")

}

func StartApp() error {
	// =========================================================================
	// initializing the authentication support
	cfg := config.GetConfig()
	log.Info().Msg("main started : initializing the authentication support")

	//reading the private key file
	decodedPVKBytes, err := base64.StdEncoding.DecodeString(cfg.AuthConfig.PrivateKey)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return err
	}

	decodedPKBytes, err := base64.StdEncoding.DecodeString(cfg.AuthConfig.PublicKey)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(decodedPVKBytes))
	if err != nil {
		return fmt.Errorf("error in parsing auth private key : %w", err) // %w is used for error wraping
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(decodedPKBytes))
	if err != nil {
		return fmt.Errorf("error in parsing auth public key : %w", err) // %w is used for error wraping
	}
	a, err := auth.NewAuth(privateKey, publicKey)
	if err != nil {
		return fmt.Errorf("error in constructing auth %w", err)
	}

	// =========================================================================
	// start the database

	log.Info().Msg("main started : initializing the data")

	db, err := database.ConnectToDatabase(cfg)
	if err != nil {
		return fmt.Errorf("error in opening the database connection : %w", err)
	}

	pg, err := db.DB()
	if err != nil {
		return fmt.Errorf("error in getting the database instance")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = pg.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("database is not connected: %w", err)
	}

	// =========================================================================
	//initializing redis

	rdb := redis.NewRedisClient(cfg)
	// initialize the repository layer
	repo, err := repository.NewRepository(db)
	if err != nil {
		return err
	}

	svc, err := service.NewService(repo, a, rdb)
	if err != nil {
		return err
	}

	// initializing the http server
	api := http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.AppConfig.Host, cfg.AppConfig.Port),
		ReadTimeout:  8000 * time.Second,
		WriteTimeout: 800 * time.Second,
		IdleTimeout:  800 * time.Second,
		Handler:      handler.SetupApi(a, svc),
	}

	// channel to store any errors while setting up the service
	serverErrors := make(chan error, 1)

	go func() {
		log.Info().Str("Port", api.Addr).Msg("main started : api is listening")
		serverErrors <- api.ListenAndServe()
	}()

	//shutdown channel intercepts ctrl+c signals
	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error : %w", err)

	case sig := <-shutdown:
		log.Info().Msgf("main: Start shutdown %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			err := api.Close()
			return fmt.Errorf("could not stop server gracefully : %w", err)
		}
	}
	return nil

}
