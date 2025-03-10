package file

import (
	query "github.com/shortlink-org/shortdb/shortdb/domain/query/v1"
	table "github.com/shortlink-org/shortdb/shortdb/domain/table/v1"
)

func (f *File) CreateTable(item *query.Query) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.database.GetTables() == nil {
		f.database.Tables = make(map[string]*table.Table)
	}

	// check
	if f.database.GetTables()[item.GetTableName()] != nil {
		return ErrExistTable
	}

	f.database.Tables[item.GetTableName()] = table.New(item)

	return nil
}

func (*File) DropTable(_ string) error {
	// TODO implement me
	return nil
}
