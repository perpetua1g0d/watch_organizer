package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/perpetua1g0d/watch_organizer/pkg/backend/model"
	"github.com/perpetua1g0d/watch_organizer/pkg/backend/repository"
	"github.com/spf13/viper"
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

// add poster by tab path
func addPoster(db *sqlx.DB) (err error) {
	var tabPath string
	fmt.Println("type tab path like 'tab1/tab2/.../lastTab':")
	if fmt.Scan(&tabPath); len(tabPath) == 0 {
		return fmt.Errorf("%s", "empty tabpath")
	}
	// check tabPath validity

	// tab := root
	tabs := strings.Split(tabPath, "/")
	if len(tabs) == 0 {
		return fmt.Errorf("%s", "path can't be splitted")
	}

	// rootId
	var root model.Tab
	err = db.Get(&root, "select * from tab where name=root;")
	if err != nil {
		return err
	}

	for _, tab := range tabs {
		var tabObj model.Tab
		if err = db.Get(&tabObj, fmt.Sprintf("select * from tab where name=%s;", tab)); err != nil {
			return err
		}

		// tabId := ...
		//
	}

	return nil
}

// get all posters in choosen tab
func getPosters() {

}

func cli() {
	for {

	}
}

func main() {
	var err error
	fmt.Println("hi hi hi")

	if err := initConfig(); err != nil {
		log.Fatalf("err init configs using viper: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.DBConfig{
		Host:     "localhost",
		Port:     "5436",
		Username: "alexey",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("failed to init DB: %s", err.Error())
	}

	// repo := repository.NewRepository(db)
	// service := service.NewService(repo)
	// handler := handler.NewHandler(service)

	// handler.InitRoutes()

	// db.Ping()
	if err = addPoster(db); err != nil {
		log.Fatalf("failed to add poster: %s", err.Error())
	}
}
