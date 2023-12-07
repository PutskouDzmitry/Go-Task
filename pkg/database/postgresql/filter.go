package postgresql

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"task/pkg/api/filter"
	"unicode/utf8"
)

type Filterable interface {
	Enrich(q sq.SelectBuilder, alias string) sq.SelectBuilder
}

type filterOptions struct {
	fields  []FilterField
	limit   uint64
	offset  uint64
	query   string
	aliases map[string]string
	names   map[string]string
	exprs   map[string]func(interface{}) sq.Sqlizer
}

type FilterField struct {
	Column   string
	Operator string
	Value    interface{}
}

func NewFilterOptions(o *filter.Options) *filterOptions {
	fields := make([]FilterField, 0)
	for _, f := range o.Fields {
		field := FilterField{
			Column:   f.Name,
			Operator: f.Operator,
			Value:    f.Value,
		}
		fields = append(fields, field)
	}

	return &filterOptions{
		limit:   o.Limit,
		offset:  o.Offset,
		fields:  fields,
		aliases: make(map[string]string),
		names:   make(map[string]string),
		exprs:   make(map[string]func(interface{}) sq.Sqlizer),
		query:   o.Query,
	}
}

func (o *filterOptions) SpecificAlias(key, value string) *filterOptions {
	o.aliases[key] = value
	return o
}

func (o *filterOptions) SpecificName(key, value string) *filterOptions {
	o.names[key] = value
	return o
}

func (o *filterOptions) SpecificExpr(key string, value func(interface{}) sq.Sqlizer) *filterOptions {
	o.exprs[key] = value
	return o
}

func (o *filterOptions) Enrich(q sq.SelectBuilder, alias string, searchIn ...string) sq.SelectBuilder {
	and := sq.And{}

	for _, field := range o.fields {
		column := field.Column
		value := field.Value
		if o.exprs[column] != nil {
			and = append(and, o.exprs[column](value))
			continue
		}

		columnAlias := ""
		if o.aliases[column] != "-" {
			if o.aliases[column] != "" {
				columnAlias = o.aliases[column]
			} else {
				columnAlias = alias
			}
		}
		if o.names[column] != "" {
			column = o.names[column]
			if column == "-" {
				continue
			}
		}

		if columnAlias != "" {
			column = fmt.Sprintf("%s.%s", columnAlias, column)
		}

		var e sq.Sqlizer
		switch field.Operator {
		case filter.OperatorEq:
			e = sq.Eq{column: value}
		case filter.OperatorNotEq:
			e = sq.NotEq{column: value}
		case filter.OperatorGt:
			e = sq.Gt{column: value}
		case filter.OperatorGte:
			e = sq.GtOrEq{column: value}
		case filter.OperatorLt:
			e = sq.Lt{column: value}
		case filter.OperatorLte:
			e = sq.LtOrEq{column: value}
		default:
			panic("bad operator")
		}
		and = append(and, e)
	}
	q = q.Where(and)

	if o.limit > 0 {
		q = q.Limit(o.limit).Offset(o.offset)
	}

	if utf8.RuneCountInString(o.query) >= 3 {

		for i := range searchIn {
			if o.aliases[searchIn[i]] == "*" {
				continue
			}
			columnAlias := ""
			if o.aliases[searchIn[i]] != "-" {
				if o.aliases[searchIn[i]] != "" {
					columnAlias = o.aliases[searchIn[i]]
				} else {
					columnAlias = alias
				}
			}
			searchIn[i] = fmt.Sprintf("%s.%s", columnAlias, searchIn[i])
		}

		for _, row := range searchIn {
			q = q.Where(sq.Like{row: "%" + o.query + "%"})
		}
	}

	return q
}
