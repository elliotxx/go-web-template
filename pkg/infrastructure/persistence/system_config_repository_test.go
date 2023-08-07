package persistence

import (
	"context"
	"testing"

	"github.com/elliotxx/go-web-template/pkg/domain/entity"
	"github.com/elliotxx/go-web-template/pkg/domain/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSystemConfigRepository(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		fakeGDB, sqlMock, err := GetMockDB()
		require.NoError(t, err)
		repo := NewSystemConfigRepository(fakeGDB)
		defer CloseDB(t, fakeGDB)
		defer sqlMock.ExpectClose()

		var (
			expectedID, expectedRows uint = 1, 1
			actual                        = entity.SystemConfig{Env: entity.EnvProd}
		)
		sqlMock.ExpectBegin()
		sqlMock.ExpectExec("INSERT").
			WillReturnResult(sqlmock.NewResult(int64(expectedID), int64(expectedRows)))
		sqlMock.ExpectCommit()
		err = repo.Create(context.Background(), &actual)
		require.NoError(t, err)
		require.Equal(t, expectedID, actual.ID)
	})

	t.Run("Delete existed record", func(t *testing.T) {
		fakeGDB, sqlMock, err := GetMockDB()
		require.NoError(t, err)
		repo := NewSystemConfigRepository(fakeGDB)
		defer CloseDB(t, fakeGDB)
		defer sqlMock.ExpectClose()

		var expectedID, expectedRows uint = 1, 1
		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(1))
		sqlMock.ExpectExec("UPDATE").
			WillReturnResult(sqlmock.NewResult(int64(expectedID), int64(expectedRows)))
		sqlMock.ExpectCommit()
		err = repo.Delete(context.Background(), expectedID)
		require.NoError(t, err)
	})

	t.Run("Delete not existing record", func(t *testing.T) {
		fakeGDB, sqlMock, err := GetMockDB()
		require.NoError(t, err)
		repo := NewSystemConfigRepository(fakeGDB)
		defer CloseDB(t, fakeGDB)
		defer sqlMock.ExpectClose()

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		err = repo.Delete(context.Background(), 1)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("Update existed record", func(t *testing.T) {
		fakeGDB, sqlMock, err := GetMockDB()
		require.NoError(t, err)
		repo := NewSystemConfigRepository(fakeGDB)

		var (
			expectedID, expectedRows uint = 1, 1
			actual                        = entity.SystemConfig{
				ID:  1,
				Env: entity.EnvProd,
			}
		)
		sqlMock.ExpectExec("UPDATE").
			WillReturnResult(sqlmock.NewResult(int64(expectedID), int64(expectedRows)))
		err = repo.Update(context.Background(), &actual)
		require.NoError(t, err)
	})

	t.Run("Update not existing record", func(t *testing.T) {
		fakeGDB, sqlMock, err := GetMockDB()
		require.NoError(t, err)
		repo := NewSystemConfigRepository(fakeGDB)
		defer CloseDB(t, fakeGDB)
		defer sqlMock.ExpectClose()

		actual := entity.SystemConfig{Env: entity.EnvProd}
		err = repo.Update(context.Background(), &actual)
		require.ErrorIs(t, err, gorm.ErrMissingWhereClause)
	})

	t.Run("Get", func(t *testing.T) {
		fakeGDB, sqlMock, err := GetMockDB()
		require.NoError(t, err)
		repo := NewSystemConfigRepository(fakeGDB)
		defer CloseDB(t, fakeGDB)
		defer sqlMock.ExpectClose()

		var (
			expectedID  uint = 1
			expectedEnv      = entity.EnvProd
		)
		sqlMock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.NewRows([]string{"id", "env"}).
				AddRow(expectedID, string(expectedEnv)))
		actual, err := repo.Get(context.Background(), expectedID)
		require.NoError(t, err)
		require.Equal(t, expectedID, actual.ID)
		require.Equal(t, expectedEnv, actual.Env)
	})

	t.Run("Find", func(t *testing.T) {
		fakeGDB, sqlMock, err := GetMockDB()
		require.NoError(t, err)
		repo := NewSystemConfigRepository(fakeGDB)
		defer CloseDB(t, fakeGDB)
		defer sqlMock.ExpectClose()

		sqlMock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.NewRows([]string{"id", "env"}).
				AddRow(1, "prod").
				AddRow(2, "dev"))
		actuals, err := repo.Find(context.Background(), repository.Query{
			Offset: 1,
			Limit:  10,
		})
		require.NoError(t, err)
		require.Equal(t, 2, len(actuals))
	})
}
