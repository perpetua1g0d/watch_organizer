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

	type args struct {
		poster model.Poster
	}
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
		name       string
		funcToTest func()
		input      args
		wantErr    bool
	}{
		{
			name: "pos: simple",
			funcToTest: func() {
				mock.ExpectBegin()

				idRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				for _, genre := range testPosters[0].Genres {
					mock.ExpectQuery("insert into postergenre").WithArgs(testPosters[0].Id, genre).
						WillReturnRows(idRows)
				}

				mock.ExpectCommit()
			},
			input: args{
				poster: testPosters[0],
			},
		},
		{
			name: "pos: no genres",
			funcToTest: func() {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			input: args{
				poster: testPosters[1],
			},
		},
		{
			name: "pos: nil genres",
			funcToTest: func() {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			input: args{
				poster: testPosters[2],
			},
		},
		{
			name: "pos: nil poster",
			funcToTest: func() {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			input: args{
				poster: testPosters[3],
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.funcToTest()

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
