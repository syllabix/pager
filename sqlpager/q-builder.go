package sqlpager

import "fmt"

type qBuilder struct {
	pqb *PagedQueryBuilder
}

func (q *qBuilder) Join(stmt string) *qJoiner {
	q.pqb.joins = append(q.pqb.joins, fmt.Sprintf("JOIN %s", stmt))
	return &qJoiner{
		pqb: q.pqb,
	}
}

func (q *qBuilder) LeftJoin(stmt string) *qJoiner {
	q.pqb.joins = append(q.pqb.joins, fmt.Sprintf("LEFT JOIN %s", stmt))
	return &qJoiner{
		pqb: q.pqb,
	}
}

func (q *qBuilder) RightJoin(stmt string) *qJoiner {
	q.pqb.joins = append(q.pqb.joins, fmt.Sprintf("RIGHT JOIN %s", stmt))
	return &qJoiner{
		pqb: q.pqb,
	}
}

func (q *qBuilder) Where(stmt string, args ...interface{}) *qBuilder {

	if len(q.pqb.whereStmts) == 0 {
		q.pqb.whereStmts = append(q.pqb.whereStmts, fmt.Sprintf("WHERE (%s)", stmt))
	} else {
		q.pqb.whereStmts = append(q.pqb.whereStmts, fmt.Sprintf("AND (%s)", stmt))
	}
	q.pqb.args = append(q.pqb.args, args...)
	return &qBuilder{
		pqb: q.pqb,
	}
}

func (q *qBuilder) OrderBy(column string) *pageResponseFactory {
	q.pqb.orderby = fmt.Sprintf("ORDER BY %s", column)
	return &pageResponseFactory{
		pqb: q.pqb,
	}
}

type pageResponseSorter struct {
	pqb *PagedQueryBuilder
}

func (sorter *pageResponseSorter) OrderBy(column string) *pageResponseFactory {
	sorter.pqb.orderby = fmt.Sprintf("ORDER BY %s", column)
	return &pageResponseFactory{
		pqb: sorter.pqb,
	}
}

type qJoiner struct {
	pqb *PagedQueryBuilder
}

func (j *qJoiner) On(stmt string) *qBuilder {
	j.pqb.joins = append(j.pqb.joins, fmt.Sprintf("ON %s", stmt))
	return &qBuilder{
		pqb: j.pqb,
	}
}
