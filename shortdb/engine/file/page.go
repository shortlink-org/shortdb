package file

import (
	"fmt"

	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
	"google.golang.org/protobuf/proto"
)

func (f *File) getPage(nameTable string, p int32) error { //nolint:unused // ignore
	table := f.database.GetTables()[nameTable]

	// read page
	pageFile, err := f.loadPage(f.pageName(nameTable))
	if err != nil {
		return err
	}

	table.Pages[p] = pageFile

	return nil
}

func (f *File) addPage(nameTable string) (int32, error) {
	table := f.database.GetTables()[nameTable]

	if table.GetStats().GetRowsCount()%table.GetOption().GetPageSize() == 0 { //nolint:nestif // ignore
		if table.GetPages() == nil {
			table.Pages = make(map[int32]*page.Page, 0)
		}

		table.Stats.PageCount += 1
		table.Pages[table.GetStats().GetPageCount()] = &page.Page{Rows: []*page.Row{}}

		// create a page file
		newPageFile, err := f.createFile(f.pageName(nameTable))
		if err != nil {
			return table.GetStats().GetPageCount(), err
		}

		err = newPageFile.Close()
		if err != nil {
			return table.GetStats().GetPageCount(), err
		}

		// if this not first page, save current date
		if table.GetStats().GetPageCount() > 0 && table.GetPages()[table.GetStats().GetPageCount()-1] != nil {
			// save data after clear memory page
			err = f.savePage(nameTable, table.GetStats().GetPageCount()-1)
			if err != nil {
				return table.GetStats().GetPageCount(), err
			}

			// clear old page
			err = f.clearPage(nameTable, table.GetStats().GetPageCount()-1)
			if err != nil {
				return table.GetStats().GetPageCount(), err
			}
		}
	}

	return table.GetStats().GetPageCount(), nil
}

func (f *File) savePage(nameTable string, pageCount int32) error {
	table := f.database.GetTables()[nameTable]

	if pageCount == -1 {
		return nil
	}

	// save date
	openFile, err := f.createFile(fmt.Sprintf("%s_%s_%d.page", f.database.GetName(), nameTable, pageCount))
	if err != nil {
		return err
	}

	defer func() {
		_ = openFile.Close() // #nosec
	}()

	payload, err := proto.Marshal(table.GetPages()[pageCount])
	if err != nil {
		return fmt.Errorf("failed to marshal page: %w", err)
	}

	// Write something
	err = f.writeFile(openFile.Name(), payload)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) clearPage(nameTable string, pageCount int32) error { //nolint:unparam // ignore param
	f.database.Tables[nameTable].Pages[pageCount] = nil

	return nil
}

//nolint:unparam // ignore param
func (f *File) clearPages(nameTable string) error {
	f.database.Tables[nameTable].Pages = nil

	return nil
}

func (f *File) pageName(nameTable string) string {
	return fmt.Sprintf("%s_%s_%d.page", f.database.GetName(), nameTable, f.database.GetTables()[nameTable].GetStats().GetPageCount())
}
