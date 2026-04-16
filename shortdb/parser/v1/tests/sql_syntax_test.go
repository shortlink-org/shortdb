//nolint:testpackage // isolated parser checks without godog TestMain
package tests

import (
	"testing"

	query "github.com/shortlink-org/shortdb/shortdb/domain/query/v1"
	parser "github.com/shortlink-org/shortdb/shortdb/parser/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSelectStarNoSemicolon(t *testing.T) {
	t.Parallel()

	p, err := parser.New("SELECT * FROM books LIMIT 10")
	require.NoError(t, err)
	require.NotNil(t, p.GetQuery())
	assert.Equal(t, query.Type_TYPE_SELECT, p.GetQuery().GetType())
}

func TestParseSelectStarLimitTrailingSemicolon(t *testing.T) {
	t.Parallel()

	p, err := parser.New("SELECT * FROM 'books_backup' LIMIT 5;")
	require.NoError(t, err)
	require.NotNil(t, p.GetQuery())
	assert.Equal(t, query.Type_TYPE_SELECT, p.GetQuery().GetType())
}
