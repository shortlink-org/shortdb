package file

import (
	"context"
	"fmt"
	"os"

	"github.com/sasha-s/go-deadlock"
	database "github.com/shortlink-org/shortdb/shortdb/domain/database/v1"
	query "github.com/shortlink-org/shortdb/shortdb/domain/query/v1"
	table "github.com/shortlink-org/shortdb/shortdb/domain/table/v1"
	"github.com/shortlink-org/shortdb/shortdb/engine/options"
	"github.com/shortlink-org/shortdb/shortdb/io_uring"
	"github.com/sourcegraph/conc"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
)

type File struct {
	mu deadlock.RWMutex

	database *database.DataBase
	path     string
}

func New(ctx context.Context, opts ...options.Option) (*File, error) {
	const SHORTDB_PAGE_SIZE = 100

	viper.AutomaticEnv()
	viper.SetDefault("SHORTDB_DEFAULT_DATABASE", "public")   // ShortDB default database
	viper.SetDefault("SHORTDB_PAGE_SIZE", SHORTDB_PAGE_SIZE) // ShortDB default page of size

	var err error

	file := &File{
		database: &database.DataBase{
			Name:   viper.GetString("SHORTDB_DEFAULT_DATABASE"),
			Tables: make(map[string]*table.Table),
		},
	}

	for _, opt := range opts {
		if errApplyOptions := opt(file); errApplyOptions != nil {
			panic(errApplyOptions)
		}
	}

	// if not set a path, set temp directory
	if file.path == "" {
		file.path = os.TempDir()
	}

	// init db
	err = file.init(ctx)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *File) Exec(in *query.Query) (any, error) {
	switch in.GetType() {
	case query.Type_TYPE_UNSPECIFIED:
		return nil, ErrIncorrectType
	case query.Type_TYPE_SELECT:
		return f.Select(in)
	case query.Type_TYPE_UPDATE:
		return nil, f.Update(in)
	case query.Type_TYPE_INSERT:
		return nil, f.Insert(in)
	case query.Type_TYPE_DELETE:
		return nil, f.Delete(in)
	case query.Type_TYPE_CREATE_TABLE:
		return nil, f.CreateTable(in)
	case query.Type_TYPE_DROP_TABLE:
		return nil, f.DropTable(in.GetTableName())
	case query.Type_TYPE_CREATE_INDEX:
		return nil, f.CreateIndex(in)
	case query.Type_TYPE_DELETE_INDEX:
		return nil, f.DropIndex(in.GetTableName())
	}

	//nolint:nilnil // it's correct
	return nil, nil
}

func (f *File) init(ctx context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// create directory if not exist
	err := os.MkdirAll(f.path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// create a file if not exist
	fileOpenFile, err := f.createFile(f.database.GetName() + ".db")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		_ = fileOpenFile.Close() // #nosec
	}()

	// init io_uring
	err = io_uring.Init()
	if err != nil {
		return fmt.Errorf("failed to init io_uring: %w", err)
	}
	defer io_uring.Cleanup()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case errIOUring := <-io_uring.Err():
				//nolint:revive,forbidigo // just print error
				fmt.Println(errIOUring)
			}
		}
	}()

	var (
		wg      conc.WaitGroup
		payload []byte
	)

	// Read a file.
	err = io_uring.ReadFile(fileOpenFile.Name(), func(buf []byte) {
		wg.Go(func() {
			payload = buf
		})
	})
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	io_uring.Poll()
	wg.Wait()

	if len(payload) != 0 {
		err = proto.Unmarshal(payload, f.database)
		if err != nil {
			return fmt.Errorf("failed to unmarshal database: %w", err)
		}
	}

	return nil
}

func (f *File) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// create a database if not exist
	databaseFile, err := f.createFile(f.database.GetName() + ".db")
	if err != nil {
		return err
	}

	defer func() {
		_ = databaseFile.Close() // #nosec
	}()

	// init io_uring
	err = io_uring.Init()
	if err != nil {
		return fmt.Errorf("failed to init io_uring: %w", err)
	}
	defer io_uring.Cleanup()

	var wg conc.WaitGroup

	// save last page
	for tableName := range f.database.GetTables() {
		err = f.savePage(tableName, f.database.GetTables()[tableName].GetStats().GetPageCount())
		if err != nil {
			return err
		}

		// clear cache
		err = f.clearPages(tableName)
		if err != nil {
			return err
		}
	}

	payload, err := proto.Marshal(f.database)
	if err != nil {
		return fmt.Errorf("failed to marshal database: %w", err)
	}

	// save database
	err = io_uring.WriteFile(databaseFile.Name(), payload, 0o644, func(n int) { //nolint:mnd,revive // #nosec
		wg.Go(func() {})
		// handle n
	})
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Call Poll to let the kernel know to read the entries.
	io_uring.Poll()
	// Wait till all callbacks are done.
	wg.Wait()

	return nil
}

func (f *File) createFile(name string) (*os.File, error) {
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", f.path, name), os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return file, nil
}

func (*File) writeFile(name string, payload []byte) error {
	err := os.WriteFile(name, payload, 0o600) //nolint:mnd,revive // #nosec
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
