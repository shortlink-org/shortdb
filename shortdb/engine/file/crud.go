package file

import (
	"fmt"

	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
	query "github.com/shortlink-org/shortdb/shortdb/domain/query/v1"
	"github.com/shortlink-org/shortdb/shortdb/engine/file/cursor"
)

func (f *File) Select(in *query.Query) ([]*page.Row, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// check table
	table := f.database.GetTables()[in.GetTableName()]
	if table == nil {
		return nil, &NotExistTableError{
			Table: in.GetTableName(),
			Type:  "SELECT",
		}
	}

	if len(in.GetFields()) == 0 {
		return nil, ErrIncorrectNameFields
	}

	// response
	response := make([]*page.Row, 0)

	currentRow := cursor.NewBuilder(table).Build()
	for !currentRow.EndOfTable {
		// load data
		if table.GetPages()[currentRow.PageId] == nil {
			pagePath := fmt.Sprintf("%s/%s_%s_%d.page", f.path, f.database.GetName(), table.GetName(), currentRow.PageId)

			payload, errLoadPage := f.loadPage(pagePath)
			if errLoadPage != nil {
				return nil, errLoadPage
			}

			if table.GetPages() == nil {
				table.Pages = make(map[int32]*page.Page, 0)
			}

			table.Pages[currentRow.PageId] = payload
		}

		// get value
		record, errGetValue := currentRow.Value()
		if errGetValue != nil {
			return nil, fmt.Errorf("get value error: %w", errGetValue)
		}

		for _, field := range in.GetFields() {
			if record.GetValue()[field] == nil {
				return nil, &IncorrectNameFieldsError{
					Field: field,
					Table: in.GetTableName(),
				}
			}
		}

		if in.IsFilter(record, table.GetFields()) {
			response = append(response, record)

			if in.IsLimit() {
				in.Limit--
			}
		}

		if !in.IsLimit() {
			break
		}

		currentRow.Advance()
	}

	return response, nil
}

func (*File) Update(_ *query.Query) error {
	// TODO implement me
	return nil
}

func (f *File) Insert(in *query.Query) error {
	err := f.insertToTable(in)
	if err != nil {
		return err
	}

	err = f.insertToIndex(in)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) insertToTable(in *query.Query) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// check the table's existence
	table := f.database.GetTables()[in.GetTableName()]
	if table == nil {
		return &NotExistTableError{
			Table: in.GetTableName(),
			Type:  "INSERT",
		}
	}

	// check if a new page needs to be created
	_, err := f.addPage(in.GetTableName())
	if err != nil {
		return ErrCreatePage
	}

	if table.GetStats().GetPageCount() > -1 && table.GetPages()[table.GetStats().GetPageCount()] == nil {
		// load page
		pagePath := fmt.Sprintf("%s/%s_%s_%d.page", f.path, f.database.GetName(), table.GetName(), table.GetStats().GetPageCount())

		payload, errLoadPage := f.loadPage(pagePath)
		if errLoadPage != nil {
			return errLoadPage
		}

		if table.GetPages() == nil {
			table.Pages = make(map[int32]*page.Page, 0)
		}

		table.Pages[table.GetStats().GetPageCount()] = payload
	}

	// insert to last page
	currentRow := cursor.NewBuilder(table).AtEnd().Build()

	row, err := currentRow.Value()
	if err != nil {
		return ErrCreateCursor
	}

	// check values and create row record
	record := page.Row{
		Value: make(map[string][]byte),
	}

	for index, field := range in.GetFields() {
		if table.GetFields()[field].String() == "" {
			return &IncorrectTypeFieldsError{
				Field: field,
				Table: in.GetTableName(),
			}
		}

		record.Value[field] = []byte(in.GetInserts()[0].GetItems()[index])
	}

	row.Value = record.GetValue()

	// update stats
	table.Stats.RowsCount += 1

	// iterator to next value
	currentRow.Advance()

	return nil
}

func (*File) insertToIndex(_ *query.Query) error {
	// TODO implement me
	return nil
}

func (*File) Delete(_ *query.Query) error {
	// TODO implement me
	return nil
}
