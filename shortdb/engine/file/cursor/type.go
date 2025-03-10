package cursor

import (
	"sync"

	table "github.com/shortlink-org/shortdb/shortdb/domain/table/v1"
)

// Cursor represents a cursor into a table.
type Cursor struct {
	mu sync.Mutex

	Table      *table.Table
	RowId      int64
	PageId     int32
	EndOfTable bool
}
