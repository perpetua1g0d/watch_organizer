package main

import (
	"fmt"
	"log"

	"github.com/perpetua1g0d/watch_organizer/pkg/backend/repository"
	"github.com/spf13/viper"
)


func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}


func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("err init configs using viper: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.DBConfig{
		Host: "localhost",
		Port: "5436",
		Username: "postgres",
		Password: "postgres",
		DBName: "postgres",
		SSLMode: "disable",
	})
	if err != nil {
		log.Fatalf("failed to init DB: %s", err.Error())
	}

	

	fmt.Println("hi hi hi")
	
}