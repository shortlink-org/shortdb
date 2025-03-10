package cursor

import (
	table "github.com/shortlink-org/shortdb/shortdb/domain/table/v1"
)

// Builder helps construct a Cursor.
type Builder struct {
	t          *table.Table
	rowId      int64
	pageId     int32
	endOfTable bool
}

// NewBuilder initializes the builder with the provided table,
// defaulting to a cursor at the beginning.
func NewBuilder(t *table.Table) *Builder {
	return &Builder{
		t:          t,
		rowId:      0,
		pageId:     0,
		endOfTable: false,
	}
}

// AtEnd configures the builder to create a cursor positioned at the end
// of the table. It retrieves the necessary statistics from the table.
func (b *Builder) AtEnd() *Builder {
	stats := b.t.GetStats()

	b.rowId = stats.GetRowsCount()  // assuming GetRowsCount returns an int64
	b.pageId = stats.GetPageCount() // assuming GetPageCount returns an int32

	b.endOfTable = true

	return b
}

// Build constructs the final Cursor using the builder's configuration.
func (b *Builder) Build() *Cursor {
	return &Cursor{
		Table:      b.t,
		RowId:      b.rowId,
		PageId:     b.pageId,
		EndOfTable: b.endOfTable,
	}
}
