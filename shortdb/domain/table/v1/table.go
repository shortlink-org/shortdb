package v1

import (
	index "github.com/shortlink-org/shortdb/shortdb/domain/index/v1"
	query "github.com/shortlink-org/shortdb/shortdb/domain/query/v1"
	"github.com/spf13/viper"
)

func New(q *query.Query) *Table {
	return &Table{
		Name:   q.GetTableName(),
		Fields: q.GetTableFields(),
		Stats: &TableStats{
			RowsCount: 0,
			PageCount: -1,
		},
		Option: &Option{
			PageSize: viper.GetInt64("SHORTDB_PAGE_SIZE"),
		},
		Index: map[string]*index.Index{},
	}
}
