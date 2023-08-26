package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TabPostgres struct {
	db *sqlx.DB
}

func NewTabPostgres(db *sqlx.DB) *TabPostgres {
	return &TabPostgres{db: db}
}

func (r *TabPostgres) CreateTabPath(tabs []string) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin tabpath tx: %w", err)
	}

	ids := make([]int, len(tabs))
	var lastTabId int
	for i, tabName := range tabs {
		var tabId int
		err = r.db.QueryRowx(fmt.Sprintf("select id from tab where name='%s';", tabName)).Scan(&tabId)
		if err != sql.ErrNoRows && err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to check for %s tab in db: %w", tabName, err)
		}

		var curTabId int
		if err == sql.ErrNoRows {
			if i == 0 {
				tx.Rollback()
				return fmt.Errorf("No root tab in db")
			}

			err = r.db.QueryRowx(fmt.Sprintf("insert into tab values (default, '%s') returning id;", tabName)).
				Scan(&curTabId)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to insert %s tab: %w", tabName, err)
			}

			err = r.db.QueryRowx(fmt.Sprintf("insert into tabchildren (%d, %d);", lastTabId, curTabId)).Err()
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to insert into tabchildren (%d, %d): %w", lastTabId, curTabId, err)
			}
		} else {
			curTabId = tabId
		}

		lastTabId = curTabId
		ids[i] = curTabId
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("Failed to commit tabpath tx: %w", err)
	}
	return
}

func (r *TabPostgres) GetTabIds(tabs []string) (err error, path []int) {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin GetTabIds tx: %w", err), nil
	}

	path = make([]int, len(tabs))
	for i, tabName := range tabs {
		var tabId int
		err = r.db.QueryRowx(fmt.Sprintf("select id from tab where name='%s';", tabName)).Scan(&tabId)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to get %s tab from db: %w", tabName, err), nil
		}

		path[i] = tabId
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("Failed to commit GetTabIds tx: %w", err)
	}
	return
}

func (r *TabPostgres) AddPosterToQueues(posterId int, path []int) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin poster_tabpath tx: %w", err)
	}

	for _, tabId := range path {
		var postersCount int
		err = r.db.QueryRowx(fmt.Sprintf("select count(*) from tabqueue where tabid=%d;", tabId)).Scan(&postersCount)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to count new post for %d tab: %w", tabId, err)
		}

		err = r.db.QueryRowx(fmt.Sprintf("insert into tabqueue values (%d, %d, %d);", tabId, posterId, postersCount+1)).
			Err()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to insert position into tabqueue values(%d, %d, %d): %w",
				tabId, posterId, postersCount+1, err)
		}
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("Failed to commit tabpath tx: %w", err)
	}
	return
}
