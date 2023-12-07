package postgresql

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"task/pkg/api/sorting"
)

type sortOptions struct {
	fields  []string
	orders  []string
	aliases map[string]string
}

func NewSortOptions(so *sorting.Options) *sortOptions {
	return &sortOptions{
		fields:  so.Fields,
		orders:  so.Orders,
		aliases: make(map[string]string),
	}
}

func (so *sortOptions) SpecificAlias(key, value string) *sortOptions {
	so.aliases[key] = value
	return so
}

func (so *sortOptions) Enrich(q squirrel.SelectBuilder, alias string) squirrel.SelectBuilder {
	for i := range so.fields {
		prefix := alias
		if so.aliases[so.fields[i]] != "" {
			prefix = so.aliases[so.fields[i]]
		}

		column := fmt.Sprintf("%s.%s", prefix, so.fields[i])

		q = q.OrderBy(fmt.Sprintf("%s %s", column, so.orders[i]))
	}
	return q
}
