package sqlpager

import (
	"fmt"
	"log"
	"regexp"

	"github.com/jmoiron/sqlx"
	"github.com/syllabix/pager"

)

//New is a factory constructor for a PagedQueryBuilder
func New(db *sqlx.DB, pageDetails *pager.PageDetails, defaults ...int) *PagedQueryBuilder {
	pagedResult := new(pager.PagedResult)

	//check if details were passed in
	if pageDetails.PageSize == 0 {
		//fall back to defaults
		if len(defaults) > 1 {
			panic("Page Size Default")
		}

		if len(defaults) == 1 {
			pageDetails.PageSize = defaults[0]
		} else {
			//fallback to hard default value of 15
			pageDetails.PageSize = 15
		}
	}

	pagedResult.PageDetails = *pageDetails

	return &PagedQueryBuilder{
		db:     db,
		result: pagedResult,
	}
}

//PagedQueryBuilder is used to construct queries and call a type that implements the PagedResponse interface
type PagedQueryBuilder struct {
	db           *sqlx.DB
	result       *pager.PagedResult
	useFullQuery bool
	fullQuery    string
	selectStmt   string
	whereStmts   []string
	args         []interface{}
	from         string
	joins        []string
	orderby      string
}

//Select method builds the select statement (returns unexported builder struct with intention to only be used in the context of a chained built query)
func (r *PagedQueryBuilder) Select(stmt string) *queryBuilderTableSelector {
	r.selectStmt = fmt.Sprintf("SELECT %s", stmt)
	return &queryBuilderTableSelector{
		pqb: r,
	}
}

var selectRegExp = regexp.MustCompile(`(SELECT.+)\sFROM`)
var queryStmtRegExp = regexp.MustCompile(`(FROM.+$)`)

func (r *PagedQueryBuilder) Query(stmt string, args ...interface{}) *pageResponseSorter {
	r.useFullQuery = true
	groups := selectRegExp.FindStringSubmatch(stmt)
	qGroups := queryStmtRegExp.FindStringSubmatch(stmt)
	if len(groups) < 2 || len(qGroups) < 2 {
		log.Panicln("Query looks to be malformed...")
	}
	r.selectStmt = groups[1]
	r.fullQuery = qGroups[1]
	r.args = args
	return &pageResponseSorter{
		pqb: r,
	}
}

type queryBuilderTableSelector struct {
	pqb *PagedQueryBuilder
}

func (q *queryBuilderTableSelector) From(table string) *qBuilder {
	q.pqb.from = fmt.Sprintf("FROM %s", table)
	q.pqb.whereStmts = make([]string, 0, 20)
	q.pqb.joins = make([]string, 0, 20)
	return &qBuilder{
		pqb: q.pqb,
	}
}
