package main

import (
	"database/sql"
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

func getRootId(db *sqlx.DB) (rootId int, err error) {
	var root model.Tab
	if err = db.Get(&root, "select * from tab where name='root';"); err != nil {
		rootId = root.Id
		// fmt.Println(rootId)
		// fmt.Printf("rootTab: %+v, extracted rootId=%d\n", root, rootId)
	}
	return
}

func createTabPath(db *sqlx.DB, tabs []string) (error, []int) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin tx: %s", err.Error()), nil
	}

	ids := make([]int, len(tabs))
	var lastTabId int
	for i, tabName := range tabs {
		var tabObj model.Tab
		if err = db.Get(&tabObj, fmt.Sprintf("select * from tab where name='%s';", tabName)); err != sql.ErrNoRows && err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to check for %s tab in db: %s", tabName, err.Error()), nil
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
				return fmt.Errorf("Failed to insert %s tab: %s", tabName, err.Error()), nil
			}

			err = db.QueryRowx(fmt.Sprintf("insert into tabchildren (%d, %d);", lastTabId, newTabId)).Err()
			fmt.Println("err inserting into tabchildren: ", err)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to insert into tabchildren (%d, %d): %s", lastTabId, newTabId, err.Error()),
					nil
			}
		} else {
			newTabId = tabObj.Id
		}

		lastTabId = newTabId
		ids[i] = newTabId
	}

	err = tx.Commit()
	return err, ids
}

func createPosterInDB(db *sqlx.DB, poster model.Poster, pathIds []int) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin tx: %s", err.Error())
	}

	var posterId int
	err = db.QueryRowx(fmt.Sprintf("insert into poster values(default, '%s', %f, '%s', %d) returning id;",
		poster.KpLink, poster.Rating, poster.Name, poster.Year)).Scan(&posterId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to insert poster %+v: ", err.Error())
	}

	for _, tabId := range pathIds {
		var postersCount int
		err = db.QueryRowx(fmt.Sprintf("select count(*) from tabqueue where tabid=%d;", tabId)).Scan(&postersCount)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to count new post for %d tab: %s", tabId, err.Error())
		}

		err = db.QueryRowx(fmt.Sprintf("insert into tabqueue values (%d, %d, %d);", tabId, posterId, postersCount+1)).
			Err()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to insert position into tabqueue values(%d, %d, %d): %s",
				tabId, posterId, postersCount+1, err.Error())
		}
	}

	return tx.Commit()
}

func formPoster() (err error, poster model.Poster) {
	fmt.Println("Enter KP link:")
	_, err = fmt.Scan(&poster.KpLink)
	if err != nil {
		return
	}

	fmt.Println("Enter rating (float):")
	_, err = fmt.Scan(&poster.Rating)
	if err != nil {
		return
	}

	fmt.Println("Enter name:")
	_, err = fmt.Scan(&poster.Name)
	if err != nil {
		return
	}

	fmt.Println("Enter year:")
	_, err = fmt.Scan(&poster.Year)

	return
}

func formTabPath() (err error, tabs []string) {
	var tabPath string
	fmt.Println("Type tab path like 'tab1/tab2/.../lastTab':")
	if fmt.Scan(&tabPath); len(tabPath) == 0 {
		return fmt.Errorf("tabpath is empty"), nil
	}
	// check tabPath validity

	if tabs = strings.Split(tabPath, "/"); len(tabs) == 0 {
		return fmt.Errorf("Failed to split tabpath"), nil
	}

	tabs = append([]string{"root"}, tabs...)

	return
}

// add poster by tab path
func addPoster(db *sqlx.DB) (err error) {
	err, tabs := formTabPath()
	fmt.Println("tabs:", tabs)
	if err != nil {
		return err
	}

	err, pathIds := createTabPath(db, tabs)
	fmt.Println("pathIds:", pathIds)
	if err != nil {
		return fmt.Errorf("Failed to create tabpath in db: %s", err.Error())
	}

	err, poster := formPoster()
	if err != nil {
		return fmt.Errorf("Failed to create poster: %s", err.Error())
	}

	err = createPosterInDB(db, poster, pathIds)
	if err != nil {
		return fmt.Errorf("Failed to create poster in db: %s", err.Error())
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

	// db.Ping()
	if err = addPoster(db); err != nil {
		log.Fatalf("Failed to add poster: %s", err.Error())
	}
}
