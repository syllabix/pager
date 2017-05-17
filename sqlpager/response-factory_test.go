package sqlpager

import (
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"stash.cqlcorp.net/go/pager"
)

type Item struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func mockDb() *sqlx.DB {
	mockDB, mock, _ := sqlmock.New()
	count := sqlmock.NewRows([]string{"Count"}).AddRow(2)
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, "Widget", "A thing to do stuff with").AddRow(2, "Apple", "Nutritious Snack")

	mock.ExpectBegin()
	mock.ExpectQuery(`^SELECT Count\(\*\) FROM items`).WithArgs().WillReturnRows(count)
	mock.ExpectQuery(`^SELECT (.+) FROM items ORDER BY name OFFSET \? ROWS FETCH NEXT \? ROWS ONLY$`).WithArgs(0, 15).WillReturnRows(rows)
	mock.ExpectCommit()

	return sqlx.NewDb(mockDB, "sqlmock")
}

func TestQueryBuilder(t *testing.T) {
	assert := assert.New(t)

	db := mockDb()
	defer db.Close()
	pagedQueryBuilder := New(db, &pager.PageDetails{})

	items := make([]Item, 0, pagedQueryBuilder.result.PageSize)
	paged, err := pagedQueryBuilder.Select("*").From("items").OrderBy("name").ToPagedResponse(&items)

	assert.Nil(err)
	assert.Equal(paged.TotalItemCount, 2)
	assert.Equal(paged.PageSize, 15)
	assert.Equal(paged.PageNumber, 1)
	assert.Equal(paged.TotalPageCount, 1)

	itemResponse, ok := paged.Results.(*[]Item)
	assert.True(ok)
	assert.Equal("Apple", (*itemResponse)[1].Name)
	assert.Equal("A thing to do stuff with", (*itemResponse)[0].Description)
}
