package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/perpetua1g0d/watch_organizer/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestPosterPostgres_CreateGenresInDB(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create new sqlmock db: %s", err.Error())
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	testPosters := []model.Poster{
		{
			Id:     1,
			KpLink: "https://test0.com",
			Rating: 8.1,
			Name:   "test0 and some spaces",
			Year:   1913,
			Genres: []string{"genre1", "genre 2", "genre3"},
		},
		{
			Id:     2,
			KpLink: "https://test1.ru",
			Rating: 1.1,
			Name:   "test1: no genres",
			Year:   1966,
			Genres: []string{},
		},
		{
			Id:     3,
			KpLink: "https://test2.",
			Rating: 5.1,
			Name:   "test2: nil genres",
			Year:   1966,
		},
		{},
	}
	tests := []struct {
		name    string
		runMock func()
		wantErr bool
	}{
		{
			name: "pos: simple",
			runMock: func() {
				mock.ExpectBegin()

				idRows := sqlmock.NewRows([]string{"id"})
				for _, genre := range testPosters[0].Genres {
					mock.ExpectQuery("insert into postergenre").WithArgs(testPosters[0].Id, genre).
						WillReturnRows(idRows)
				}

				mock.ExpectCommit()
			},
		},
		{
			name: "pos: no genres",
			runMock: func() {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
		},
		{
			name: "pos: nil genres",
			runMock: func() {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
		},
		{
			name: "pos: nil poster",
			runMock: func() {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.runMock()

			err := createGenresInDB(sqlxDB, testPosters[i])

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPosterPostgres_Create(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create new sqlmock db: %s", err.Error())
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(sqlxDB)

	testPosters := []model.Poster{
		{
			Id:     1,
			KpLink: "https://test0.com",
			Rating: 8.1,
			Name:   "test0 and some spaces",
			Year:   1913,
			Genres: []string{"genre1", "genre 2", "genre3"},
		},
	}
	tests := []struct {
		name    string
		runMock func()
		wantErr bool
	}{
		{
			name: "pos: simple",
			runMock: func() {
				mock.ExpectBegin()
				resId := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("insert into poster").
					WithArgs(testPosters[0].KpLink, testPosters[0].Rating, testPosters[0].Name, testPosters[0].Year).
					WillReturnRows(resId)

				mock.ExpectBegin()
				for _, genre := range testPosters[0].Genres {
					mock.ExpectQuery("insert into postergenre").WithArgs(testPosters[0].Id, genre).
						WillReturnRows(resId)
				}
				mock.ExpectCommit()

				mock.ExpectCommit()
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.runMock()

			err := repo.Create(testPosters[i])

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
