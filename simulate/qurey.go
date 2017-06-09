package simulate

import (
	"time"
	"fmt"
)

type queryBuilder struct {
	p_SELECT_FIELD []string
	p_FROM string
	p_WHERE []string
	p_ORDER_BY string
	p_GROUP_BY string
	p_LIMIT int64
}

func newQuery() (*queryBuilder) {
	c := &queryBuilder{
		[]string{},
		"",
		[]string{},
		"",
		"",
		-1,
	}
	return c
}
func (query *queryBuilder)Select(field... string) *queryBuilder {
	query.p_SELECT_FIELD = append(field)
	return query
}
func (query *queryBuilder)SelectCount(field string) *queryBuilder {
	query.p_SELECT_FIELD = append([]string{"COUNT("+field+")"})
	return query
}

func (query *queryBuilder)From(from string) *queryBuilder {
	query.p_FROM = from
	return query
}
//WHERE
func (query *queryBuilder)TAG(tag,value string) *queryBuilder {
	q := fmt.Sprintf("%s = '%s'", tag,value)
	query.p_WHERE = append(query.p_WHERE ,q)
	return query
}

func (query *queryBuilder)TIME(start,end time.Time) *queryBuilder {
	layout := "2006-01-02 15:04:05"
	start_str := start.Format(layout)
	end_str := end.Format(layout)
	q := fmt.Sprintf("time <= '%s' AND time >= '%s'",end_str,start_str)
	query.p_WHERE = append(query.p_WHERE ,q)
	return query
}
//GROUP BY
func (query *queryBuilder)GroupByTime(filed string) *queryBuilder {
	query.p_GROUP_BY = "time("+filed+")"
	return query
}
//GROUP BY time(12m)

//ORDER BY time ASC
func (query *queryBuilder)ASC(filed string) *queryBuilder {
	query.p_ORDER_BY = filed+" ASC"
	return query
}

func (query *queryBuilder)DESC(filed string) *queryBuilder {
	query.p_ORDER_BY = filed+" DESC"
	return query
}
func (query *queryBuilder)Limit(limit int64) *queryBuilder {
	query.p_LIMIT = limit
	return query
}

//linear
func (query *queryBuilder)Build() string{
	final_query := ""
	//SELECT
	switch len(query.p_SELECT_FIELD) {
	case 0:
		final_query += "SELECT *"
	case 1:
		final_query += "SELECT "+ query.p_SELECT_FIELD[0]
	default:
		final_query += "SELECT "
		for i:=0;i< len(query.p_SELECT_FIELD);i++  {
			final_query += query.p_SELECT_FIELD[i]
			if i+1 != len(query.p_SELECT_FIELD) {
				final_query += ","
			}
		}
	}
	//FROM
	final_query += " FROM " + query.p_FROM +" "

	//WHERE
	switch len(query.p_WHERE) {
	case 0:

	case 1:
		final_query += "WHERE "+ query.p_WHERE[0] +" "
	default:
		final_query += "WHERE "
		for i:=0;i< len(query.p_WHERE);i++  {
			final_query += query.p_WHERE[i]
			if i+1 != len(query.p_WHERE) {
				final_query += " AND "
			}else{
				final_query += " "
			}
		}
	}

	if query.p_GROUP_BY != "" {
		final_query += "GROUP BY " + query.p_GROUP_BY + " "
	}
	final_query += "Fill(linear) "
	//ORDER BY
	if query.p_ORDER_BY != "" {
		final_query += "ORDER BY " + query.p_ORDER_BY + " "
	}

	//LIMIT
	if query.p_LIMIT != -1 {
		q := fmt.Sprintf("LIMIT %d", query.p_LIMIT)
		final_query += q
	}
	return final_query
}
