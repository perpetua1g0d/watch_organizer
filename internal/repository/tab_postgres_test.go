package repository

import (
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTabPostgres_createTab(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock db: %s", err.Error())
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	input := []struct {
		prevTabId int
		tabName   string
	}{
		{
			prevTabId: 0,
			tabName:   "tab1",
		},
	}
	tests := []struct {
		name    string
		runMock func()
		wantVal int
		isNeg   bool
	}{
		{
			name: "pos simple",
			runMock: func() {
				mock.ExpectQuery("insert into tab").
					WithArgs(input[0].tabName).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("insert into tabchildren").
					WithArgs(input[0].prevTabId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantVal: 1,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.runMock()

			err, got := createTab(sqlxDB, input[i].prevTabId, input[i].tabName)
			if tt.isNeg {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantVal, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTabPostgres_CreateTabPath(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create sqlmock DB: %s", err.Error())
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repo := NewRepository(sqlxDB)

	input := []struct {
		tabs []string
	}{
		{
			tabs: []string{"root", "tab1", "tab2"},
		},
		{
			tabs: []string{"root", "t"},
		},
		{
			tabs: []string{"root"},
		},
		{
			tabs: []string{"not root", "tab1"},
		},
		{
			tabs: []string{"root", "tab1"},
		},
	}
	tests := []struct {
		name    string
		runMock func()
		isNeg   bool
	}{
		{
			name: "pos simple",
			runMock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("root").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("tab1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("insert into tab").
					WithArgs("tab1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
				mock.ExpectQuery("insert into tabchildren").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("tab2").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("insert into tab").
					WithArgs("tab2").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))
				mock.ExpectQuery("insert into tabchildren").
					WithArgs(2, 3).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectCommit()
			},
		},
		{
			name: "pos one tab",
			runMock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("root").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("t").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("insert into tab").
					WithArgs("t").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
				mock.ExpectQuery("insert into tabchildren").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectCommit()
			},
		},
		{
			name: "pos only root",
			runMock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("root").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
		},
		{
			name: "neg first tab is not root",
			runMock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("not root").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectRollback()
			},
			isNeg: true,
		},
		{
			name: "neg root is absent in db",
			runMock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("root").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectRollback()
			},
			isNeg: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.runMock()

			err := repo.CreateTabPath(input[i].tabs)
			if tt.isNeg {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTabPostgres_GetTabIds(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create sqlmock DB: %s", err.Error())
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repo := NewRepository(sqlxDB)

	input := []struct {
		tabs []string
	}{
		{
			tabs: []string{"root", "tab1", "tab2"},
		},
		{
			tabs: []string{"root", "tab1", "tab2"},
		},
	}

	tests := []struct {
		name    string
		runMock func()
		isNeg   bool
		wantVal []int
	}{
		{
			name: "pos simple",
			runMock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("root").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("tab1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("tab2").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

				mock.ExpectCommit()
			},
			wantVal: []int{1, 2, 3},
		},
		{
			name: "neg last tab is absent in db",
			runMock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("root").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("tab1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

				mock.ExpectQuery("select (.+) from tab where (.+)").
					WithArgs("tab2").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectRollback()
			},
			isNeg: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.runMock()

			err, path := repo.GetTabIds(input[i].tabs)
			if tt.isNeg {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantVal, path)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTabPostgres_AddPosterToQueues(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create sqlmock DB: %s", err.Error())
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(sqlxDB)

	input := []struct {
		posterId int
		path     []int
	}{
		{
			posterId: 1,
			path:     []int{1, 2, 3},
		},
		{
			posterId: 5,
			path: []int{1, 7},
		},
	}

	tests := []struct {
		name    string
		runMock func()
		isNeg   bool
	}{
		{
			name: "pos simple",
			runMock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("select (.+) from tabqueue where (.+)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
				mock.ExpectQuery("insert into tabqueue").
					WithArgs(1, 1, 3).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectQuery("select (.+) from tabqueue where (.+)").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectQuery("insert into tabqueue").
					WithArgs(2, 1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectQuery("select (.+) from tabqueue where (.+)").
					WithArgs(3).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectQuery("insert into tabqueue").
					WithArgs(3, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectCommit()
			},
		},
		{
			name: "neg tabqueue doesn't contain last tab",
			runMock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("select (.+) from tabqueue where (.+)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
				mock.ExpectQuery("insert into tabqueue").
					WithArgs(1, 5, 4).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectQuery("select (.+) from tabqueue where (.+)").
					WithArgs(7).
					WillReturnRows(sqlmock.NewRows([]string{"count"}))
				
				mock.ExpectRollback()
			},
			isNeg: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.runMock()

			err := repo.AddPosterToQueues(input[i].posterId, input[i].path)
			if tt.isNeg {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
