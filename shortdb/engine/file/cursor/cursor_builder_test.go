package cursor_test

import (
	"testing"

	tablepb "github.com/shortlink-org/shortdb/shortdb/domain/table/v1"
	cursor2 "github.com/shortlink-org/shortdb/shortdb/engine/file/cursor"
	"github.com/stretchr/testify/assert"
)

// TestBuilder_Build_Default ensures that the builder produces a Cursor
// with default values (i.e. beginning of the table).
func TestBuilder_Build_Default(t *testing.T) {
	// Arrange: create a real table with predefined stats.
	tbl := &tablepb.Table{
		Stats: &tablepb.TableStats{
			//nolint:revive // it's ok for tests
			RowsCount: 42,
			//nolint:revive // it's ok for tests
			PageCount: 5,
		},
	}

	// Act: build a cursor using the default builder.
	builder := cursor2.NewBuilder(tbl)
	cursor := builder.Build()

	// Assert: verify the cursor is positioned at the beginning.
	assert.Equal(t, tbl, cursor.Table, "Expected table to be the same")
	assert.Equal(t, int64(0), cursor.RowId, "Expected default row id to be 0")
	assert.Equal(t, int32(0), cursor.PageId, "Expected default page id to be 0")
	assert.False(t, cursor.EndOfTable, "Expected EndOfTable to be false")
}

// TestBuilder_Build_AtEnd ensures that the builder produces a Cursor
// with values set from the table stats when AtEnd() is invoked.
func TestBuilder_Build_AtEnd(t *testing.T) {
	// Arrange: create a real table with predefined stats.
	tbl := &tablepb.Table{
		Stats: &tablepb.TableStats{
			RowsCount: 100,
			PageCount: 10,
		},
	}

	// Act: build a cursor positioned at the end of the table.
	builder := cursor2.NewBuilder(tbl).AtEnd()
	cursor := builder.Build()

	// Assert: verify the cursor is positioned at the end.
	assert.Equal(t, tbl, cursor.Table, "Expected table to be the same")
	assert.Equal(t, int64(100), cursor.RowId, "Expected row id to match TableStats.RowsCount")
	assert.Equal(t, int32(10), cursor.PageId, "Expected page id to match TableStats.PageCount")
	assert.True(t, cursor.EndOfTable, "Expected EndOfTable to be true")
}
