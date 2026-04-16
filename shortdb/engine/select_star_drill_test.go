package engine_test

import (
	"context"
	"os"
	"testing"

	"github.com/shortlink-org/shortdb/shortdb/engine"
	"github.com/shortlink-org/shortdb/shortdb/engine/file"
	parser "github.com/shortlink-org/shortdb/shortdb/parser/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mirrors Observable drill: catalog-style quoted table name + SELECT * + LIMIT.
func TestSelectStarQuotedTableObservableDrill(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	dir := t.TempDir()

	store, err := engine.New(ctx, "file", file.SetName("obsDrillDb"), file.SetPath(dir))
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, store.Close())

		require.NoError(t, os.RemoveAll(dir))
	})

	qCreate, err := parser.New("CREATE TABLE items ( sku text, qty integer );")
	require.NoError(t, err)

	_, err = store.Exec(qCreate.GetQuery())
	require.NoError(t, err)

	qIns, err := parser.New("INSERT INTO items ( sku, qty ) VALUES ( 'widget', '3' );")
	require.NoError(t, err)

	require.NoError(t, store.Insert(qIns.GetQuery()))

	drillSQL := "SELECT * FROM 'items' LIMIT 200;" //nolint:unqueryvet // Observable drill; engine expands *.
	pr, err := parser.New(drillSQL)
	require.NoError(t, err, "parser must accept Observable drill SQL: %s", drillSQL)

	assert.Equal(t, "*", pr.GetQuery().GetFields()[0], "first select field should be asterisk token")

	rows, err := store.Select(pr.GetQuery())
	require.NoError(t, err, "Select after CREATE+INSERT should succeed")
	require.Len(t, rows, 1)

	val := rows[0].GetValue()
	assert.Equal(t, "widget", string(val["sku"]))
	assert.Equal(t, "3", string(val["qty"]))
}

func TestSelectStarCatalogShortdbTables(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	dir := t.TempDir()

	store, err := engine.New(ctx, "file", file.SetName("catStarDb"), file.SetPath(dir))
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, store.Close())

		require.NoError(t, os.RemoveAll(dir))
	})

	qCreate, err := parser.New("CREATE TABLE t1 ( id integer );")
	require.NoError(t, err)

	_, err = store.Exec(qCreate.GetQuery())
	require.NoError(t, err)

	pr, err := parser.New("SELECT * FROM 'shortdb_tables';") //nolint:unqueryvet // catalog discovery
	require.NoError(t, err)

	rows, err := store.Select(pr.GetQuery())
	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, "t1", string(rows[0].GetValue()["name"]))
}
