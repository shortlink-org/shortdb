package engine_test

import (
	"context"
	"testing"

	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
	"github.com/shortlink-org/shortdb/shortdb/engine"
	"github.com/shortlink-org/shortdb/shortdb/engine/file"
	parser "github.com/shortlink-org/shortdb/shortdb/parser/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoveryShortdbTables(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	store, err := engine.New(ctx, "file", file.SetName("discoverydb"), file.SetPath(t.TempDir()))
	require.NoError(t, err)

	t.Cleanup(func() {
		cancel()
		_ = store.Close()
	})

	qBooks, err := parser.New("create table books (id integer, title string)")
	require.NoError(t, err)

	_, err = store.Exec(qBooks.GetQuery())
	require.NoError(t, err)

	qMeta, err := parser.New("SELECT name, columns FROM shortdbcatalog")
	require.NoError(t, err)

	resp, err := store.Exec(qMeta.GetQuery())
	require.NoError(t, err)

	rows, ok := resp.([]*page.Row)
	require.True(t, ok)
	require.Len(t, rows, 1)
	assert.Equal(t, "books", string(rows[0].GetValue()["name"]))
	assert.Contains(t, string(rows[0].GetValue()["columns"]), "id integer")
	assert.Contains(t, string(rows[0].GetValue()["columns"]), "title string")

	t.Run("reserved table name rejected", func(t *testing.T) {
		qBad, err := parser.New("create table shortdbcatalog (x integer)")
		require.NoError(t, err)

		_, err = store.Exec(qBad.GetQuery())
		require.Error(t, err)
	})

	t.Run("WHERE name filter on catalog", func(t *testing.T) {
		qWhere, err := parser.New("SELECT name, columns FROM shortdbcatalog WHERE name = 'books'")
		require.NoError(t, err)

		out, err := store.Exec(qWhere.GetQuery())
		require.NoError(t, err)

		rows, ok := out.([]*page.Row)
		require.True(t, ok)
		require.Len(t, rows, 1)
		assert.Equal(t, "books", string(rows[0].GetValue()["name"]))
	})

	t.Run("WHERE unsupported column on catalog", func(t *testing.T) {
		qBad, err := parser.New("SELECT name FROM shortdbcatalog WHERE columns = 'x'")
		require.NoError(t, err)

		_, err = store.Exec(qBad.GetQuery())
		require.Error(t, err)
	})

	t.Run("quoted catalog alias", func(t *testing.T) {
		q, err := parser.New("SELECT name FROM 'shortdb_tables'")
		require.NoError(t, err)

		out, err := store.Exec(q.GetQuery())
		require.NoError(t, err)

		r, ok := out.([]*page.Row)
		require.True(t, ok)
		require.Len(t, r, 1)
		assert.Equal(t, "books", string(r[0].GetValue()["name"]))
	})

	t.Run("REPL catalog query string", func(t *testing.T) {
		q, err := parser.New("SELECT name, columns FROM 'shortdb_tables';")
		require.NoError(t, err)

		out, err := store.Exec(q.GetQuery())
		require.NoError(t, err)

		r, ok := out.([]*page.Row)
		require.True(t, ok)
		require.Len(t, r, 1)
		assert.Equal(t, "books", string(r[0].GetValue()["name"]))
	})
}
