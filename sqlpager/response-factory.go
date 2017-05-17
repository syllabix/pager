package sqlpager

import (
	"errors"
	"fmt"
	"strings"

	"github.com/syllabix/pager"
)

type pageResponseFactory struct {
	pqb *PagedQueryBuilder
}

func (f *pageResponseFactory) ToPagedResponse(result interface{}) (*pager.PagedResult, error) {

	//declare paging and count statements, that depending on driver, will be set accordingly
	var pagingStmt string
	var countStmt string

	switch f.pqb.db.DriverName() {
	case "mssql", "sqlmock":

		pagingStmt = "OFFSET ? ROWS FETCH NEXT ? ROWS ONLY"
		countStmt = "SELECT Count(*)"

	default:
		return nil, errors.New("SQL Driver not currently supported")
	}

	//begin TRANSACTION
	tx, err := f.pqb.db.Beginx()
	if err != nil {
		return nil, err
	}

	//depending on whether a full query or "fluent migrator" syntax was used, assemble base query
	var query string
	if !f.pqb.useFullQuery {
		joins := strings.TrimSpace(strings.Join(f.pqb.joins, " "))
		where := strings.TrimSpace(strings.Join(f.pqb.whereStmts, " "))
		query = strings.TrimSpace(fmt.Sprintf("%s %s %s", f.pqb.from, joins, where))
	} else {
		query = strings.TrimSpace(f.pqb.fullQuery)
	}

	//get the count
	countQuery := fmt.Sprintf("%s %s", countStmt, query)
	var count int
	err = tx.Get(&count, countQuery, f.pqb.args...)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//assemble fully query
	order := strings.TrimSpace(f.pqb.orderby)
	query = fmt.Sprintf("%s %s %s %s", f.pqb.selectStmt, query, order, pagingStmt)
	args := append(f.pqb.args, f.pqb.result.Offset())
	args = append(args, f.pqb.result.PageSize)
	err = tx.Select(result, query, args...)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	response := f.pqb.result.Result(count)
	response.Results = result
	return &response, nil
}
