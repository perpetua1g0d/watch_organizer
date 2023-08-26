package main

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/perpetua1g0d/watch_organizer/internal/handler"
	"github.com/perpetua1g0d/watch_organizer/internal/model"
	"github.com/perpetua1g0d/watch_organizer/internal/repository"
	"github.com/perpetua1g0d/watch_organizer/internal/service"
	"github.com/spf13/viper"
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

// get all posters in choosen tab
func getPosters() {

}

func cli() {
	for {

	}
}

func readPoster() (err error, poster model.Poster) {
	fmt.Println("Enter KP link:")
	if _, err = fmt.Scan(&poster.KpLink); err != nil {
		return
	}

	fmt.Println("Enter rating (float):")
	if _, err = fmt.Scan(&poster.Rating); err != nil {
		return
	}

	fmt.Println("Enter name:")
	if _, err = fmt.Scan(&poster.Name); err != nil {
		return
	}

	fmt.Println("Enter year:")
	if _, err = fmt.Scan(&poster.Year); err != nil {
		return
	}

	var genresStr string
	fmt.Println("Enter genres (divide by comma & space like genre1, genre2, ...):")
	if _, err = fmt.Scan(&genresStr); err != nil {
		return fmt.Errorf("Failed to read genres: %w", err), poster
	}
	if len(genresStr) == 0 {
		return fmt.Errorf("Genres are empty"), poster
	}

	poster.Genres = strings.Split(genresStr, ", ")
	return
}

func readTabPath() (err error, tabPath string) {
	fmt.Println("Type tab path like 'tab1/tab2/.../lastTab':")
	_, err = fmt.Scan(&tabPath)
	return
}

func main() {
	var err error
	fmt.Println("hi hi hi")

	if err = initConfig(); err != nil {
		log.Fatalf("Failed to init configs using viper: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.DBConfig{
		Host:     "localhost",
		Port:     "5436",
		Username: "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("Failed to init DB: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	err, tabPath := readTabPath()
	if err != nil {
		log.Fatalf("Failed to read tabpath: %s", err.Error())
	}

	err, poster := readPoster()
	if err != nil {
		log.Fatalf("Failed to read poster: %s", err.Error())
	}

	handler.CreatePoster(tabPath, poster)
}
