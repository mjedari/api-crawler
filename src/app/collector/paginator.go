package collector

import (
	"fmt"
)

type PaginationQuery struct {
	Offset int
	Count  int
}

func (q PaginationQuery) toString() string {
	return fmt.Sprintf("&count=%v&offest=%v", q.Count, q.Offset)
}

type PaginationQueryList struct {
	Queries []PaginationQuery
	index   int
}

func NewPaginationQueryList() *PaginationQueryList {
	return &PaginationQueryList{index: 0}
}

func (l *PaginationQueryList) next() *PaginationQuery {
	if len(l.Queries) >= l.index+1 {
		l.index++
		return &l.Queries[l.index]
	}
	return nil
}

func (l *PaginationQueryList) addQuery(query PaginationQuery) {
	l.Queries = append(l.Queries, query)
}

func initiateList(factor, total int) []PaginationQuery {
	iteration := total % factor
	if iteration != 0 {
		iteration = total / factor
	}

	queryList := NewPaginationQueryList()
	query := PaginationQuery{
		Offset: 0,
		Count:  factor,
	}
	queryList.addQuery(query)

	for i := 0; i < iteration; i++ {
		query := PaginationQuery{
			Offset: factor*(i+1) + 1,
			Count:  factor,
		}
		queryList.addQuery(query)
	}
	return queryList.Queries
}
