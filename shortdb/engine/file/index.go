package file

import (
	"fmt"
	"strings"

	index "github.com/shortlink-org/shortdb/shortdb/domain/index/v1"
	query "github.com/shortlink-org/shortdb/shortdb/domain/query/v1"
	fileIndex "github.com/shortlink-org/shortdb/shortdb/engine/file/index"
	parser "github.com/shortlink-org/shortdb/shortdb/parser/v1"
)

func (f *File) CreateIndex(in *query.Query) error {
	table := f.database.GetTables()[in.GetTableName()]

	if table.GetIndex() == nil {
		table.Index = make(map[string]*index.Index)
	}

	// check
	for i := range in.GetIndexs() {
		if table.GetIndex()[in.GetIndexs()[i].GetName()] != nil {
			return &CreateExistIndexError{Name: in.GetIndexs()[i].GetName()}
		}
	}

	// create
	for i := range in.GetIndexs() {
		// create index
		table.Index[in.GetIndexs()[i].GetName()] = &index.Index{
			Name:   in.GetIndexs()[i].GetName(),
			Type:   in.GetIndexs()[i].GetType(),
			Fields: in.GetIndexs()[i].GetFields(),
		}

		// get all values
		// TODO: use pattern iterator
		cmd, err := parser.New(fmt.Sprintf("SELECT %s from %s", strings.Join(in.GetIndexs()[i].GetFields(), ","), in.GetTableName()))
		if err != nil {
			return fmt.Errorf("failed to parse in: %w", err)
		}

		rows, err := f.Select(cmd.GetQuery())
		if err != nil { //nolint:staticcheck // ignore
			// NOTE: ignore empty table
		}

		// build index
		tree, err := fileIndex.New(table.GetIndex()[in.GetIndexs()[i].GetName()], rows)
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}

		// save to file
		payload, err := tree.Marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal index: %w", err)
		}

		// save date
		err = f.saveData(in, payload, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *File) saveData(in *query.Query, payload []byte, i int) error {
	openFile, err := f.createFile(fmt.Sprintf("%s_%s_%s.index.json", f.database.GetName(), in.GetTableName(), in.GetIndexs()[i].GetName()))
	if err != nil {
		return err
	}

	defer func() {
		_ = openFile.Close() // #nosec
	}()

	// Write something
	err = f.writeFile(openFile.Name(), payload)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) DropIndex(_ string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// TODO implement me
	panic("implement me")
}
