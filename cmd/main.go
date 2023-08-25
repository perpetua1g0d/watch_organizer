package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/perpetua1g0d/watch_organizer/internal/model"
	"github.com/perpetua1g0d/watch_organizer/internal/repository"
	"github.com/spf13/viper"
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func createTabPath(db *sqlx.DB, tabs []string) (error, []int) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin tabpath tx: %w", err), nil
	}

	ids := make([]int, len(tabs))
	var lastTabId int
	for i, tabName := range tabs {
		var tabObj model.Tab
		if err = db.Get(&tabObj, fmt.Sprintf("select * from tab where name='%s';", tabName)); err != sql.ErrNoRows && err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to check for %s tab in db: %w", tabName, err), nil
		}

		var newTabId int
		if err == sql.ErrNoRows {
			if i == 0 {
				tx.Rollback()
				return fmt.Errorf("No root tab in db"), nil
			}

			err = db.QueryRowx(fmt.Sprintf("insert into tab values (default, '%s') returning id;", tabName)).
				Scan(&newTabId)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to insert %s tab: %w", tabName, err), nil
			}

			err = db.QueryRowx(fmt.Sprintf("insert into tabchildren (%d, %d);", lastTabId, newTabId)).Err()
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to insert into tabchildren (%d, %d): %w", lastTabId, newTabId, err), nil
			}
		} else {
			newTabId = tabObj.Id
		}

		lastTabId = newTabId
		ids[i] = newTabId
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit tabpath tx: %w", err), ids
	}
	return err, ids
}

func createPosterInDB(db *sqlx.DB, poster model.Poster, pathIds []int) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin poster tx: %w", err)
	}

	var posterId int
	err = db.QueryRowx(fmt.Sprintf("insert into poster values(default, '%s', %f, '%s', %d) returning id;",
		poster.KpLink, poster.Rating, poster.Name, poster.Year)).Scan(&posterId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to insert poster %+v: %w", poster, err)
	}

	for _, tabId := range pathIds {
		var postersCount int
		err = db.QueryRowx(fmt.Sprintf("select count(*) from tabqueue where tabid=%d;", tabId)).Scan(&postersCount)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to count new post for %d tab: %w", tabId, err)
		}

		err = db.QueryRowx(fmt.Sprintf("insert into tabqueue values (%d, %d, %d);", tabId, posterId, postersCount+1)).
			Err()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to insert position into tabqueue values(%d, %d, %d): %w",
				tabId, posterId, postersCount+1, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit poster tx: %w", err)
	}
	return
}

func createGenresInDB(db *sqlx.DB, posterId int, genres []string) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin genres tx: %w", err)
	}

	for _, genre := range genres {
		err = db.QueryRowx(fmt.Sprintf("insert into postergenre values (%d, %s);", posterId, genre)).Err()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to insert genre %s: %w", genre, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit genres tx: %w", err)
	}
	return
}

func readPoster() (err error, poster model.Poster, genres []string) {
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
		return fmt.Errorf("Failed to read genres: %w", err), poster, nil
	}
	if len(genresStr) == 0 {
		return fmt.Errorf("Genres are empty"), poster, nil
	}

	genres = strings.Split(genresStr, ", ")
	return
}

func formTabPath(tabPath string) (err error, tabs []string) {
	if len(tabPath) == 0 {
		return fmt.Errorf("tabpath is empty"), nil
	}
	if tabs = strings.Split(tabPath, "/"); len(tabs) == 0 {
		return fmt.Errorf("Failed to split tabpath"), nil
	}

	tabs = append([]string{"root"}, tabs...)
	return
}

// add poster by tab path
func addPoster(db *sqlx.DB) (err error) {
	var tabPath string
	fmt.Println("Type tab path like 'tab1/tab2/.../lastTab':")
	fmt.Scan(&tabPath)
	err, tabs := formTabPath(tabPath)
	fmt.Println("tabs:", tabs)
	if err != nil {
		return err
	}

	err, pathIds := createTabPath(db, tabs)
	// fmt.Println("pathIds:", pathIds)
	if err != nil {
		return fmt.Errorf("Failed to create tabpath in db: %w", err)
	}

	err, poster, genres := readPoster()
	if err != nil {
		return fmt.Errorf("Failed to create poster: %w", err)
	}

	err = createPosterInDB(db, poster, pathIds)
	if err != nil {
		return fmt.Errorf("Failed to create poster in db: %w", err)
	}

	err = createGenresInDB(db, poster.Id, genres)
	if err != nil {
		return fmt.Errorf("Failed to create genres in db: %w", err)
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
		log.Fatalf("Err init configs using viper: %s", err.Error())
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

	// repo := repository.NewRepository(db)
	// service := service.NewService(repo)
	// handler := handler.NewHandler(service)

	// handler.InitRoutes()

	if err = addPoster(db); err != nil {
		log.Fatalf("Failed to add poster: %s", err.Error())
	}
}
