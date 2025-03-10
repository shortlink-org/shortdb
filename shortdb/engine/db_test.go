package engine_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/shortlink-org/shortdb/shortdb/engine"
	"github.com/shortlink-org/shortdb/shortdb/engine/file"
	parser "github.com/shortlink-org/shortdb/shortdb/parser/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)

	os.Exit(m.Run())
}

//nolint:revive,maintidx // ignore style for tests
func TestDatabase(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	// set engine
	path := "/tmp/shortdb_test_unit"

	store, err := engine.New(ctx, "file", file.SetName("testDatabase"), file.SetPath(path))
	require.NoError(t, err)

	t.Cleanup(func() {
		cancel()

		err = os.RemoveAll(path)
		require.NoError(t, err)

		_ = store.Close()
	})

	t.Run("CREATE TABLE", func(t *testing.T) {
		// create table
		qCreateTable, errParser := parser.New("create table users (id integer, name string, active bool)")
		require.NoError(t, errParser)

		_, errExec := store.Exec(qCreateTable.GetQuery())
		require.NoError(t, errExec)

		// save data
		errClose := store.Close()
		require.NoError(t, errClose)
	})

	t.Run("INSERT INTO USERS SINGLE", func(t *testing.T) {
		qInsertUsers, errParser := parser.New("insert into users ('id', 'name', 'active') VALUES ('1', 'Ivan', 'false')")
		require.NoError(t, errParser)

		errParser = store.Insert(qInsertUsers.GetQuery())
		require.NoError(t, errParser)

		errParser = store.Insert(qInsertUsers.GetQuery())
		require.NoError(t, errParser)

		errParser = store.Insert(qInsertUsers.GetQuery())
		require.NoError(t, errParser)

		// save data
		errClose := store.Close()
		require.NoError(t, errClose)
	})

	t.Run("INSERT INTO USERS", func(t *testing.T) {
		for i := range 1000 {
			qInsertUsers, errParserNew := parser.New(fmt.Sprintf("insert into users ('id', 'name', 'active') VALUES ('%d', 'Ivan', 'false')", i))
			require.NoError(t, errParserNew)

			errInsert := store.Insert(qInsertUsers.GetQuery())
			require.NoError(t, errInsert)
		}

		// save data
		err = store.Close()
		require.NoError(t, err)
	})

	t.Run("INSERT INTO USERS +173", func(t *testing.T) {
		for i := range 173 {
			qInsertUsers, errParserNew := parser.New(fmt.Sprintf("insert into users ('id', 'name', 'active') VALUES ('%d', 'Ivan', 'false')", i))
			require.NoError(t, errParserNew)

			errInsert := store.Insert(qInsertUsers.GetQuery())
			require.NoError(t, errInsert)
		}

		// save data
		err = store.Close()
		require.NoError(t, err)
	})

	t.Run("INSERT INTO USERS +207", func(t *testing.T) {
		for i := range 207 {
			qInsertUsers, errParserNew := parser.New(fmt.Sprintf("insert into users ('id', 'name', 'active') VALUES ('%d', 'Ivan', 'false')", i))
			require.NoError(t, errParserNew)

			errInsert := store.Insert(qInsertUsers.GetQuery())
			require.NoError(t, errInsert)
		}

		// save data
		err = store.Close()
		require.NoError(t, err)
	})

	t.Run("SELECT USERS WITH LIMIT 300", func(t *testing.T) {
		qSelectUsers, err := parser.New("select id, name, active from users limit 300")
		require.NoError(t, err)

		resp, err := store.Select(qSelectUsers.GetQuery())
		require.NoError(t, err)
		assert.Len(t, resp, 300)
	})

	t.Run("SELECT USERS WITH WHERE id=99 AND LIMIT 2", func(t *testing.T) {
		qSelectUsers, err := parser.New("select id, name, active from users where id='99' limit 2")
		require.NoError(t, err)

		resp, err := store.Select(qSelectUsers.GetQuery())
		require.NoError(t, err)
		assert.Len(t, resp, 2)
	})

	t.Run("SELECT USERS FULL SCAN", func(t *testing.T) {
		qSelectUsers, err := parser.New("select id, name, active from users")
		require.NoError(t, err)

		resp, err := store.Select(qSelectUsers.GetQuery())
		require.NoError(t, err)
		assert.Len(t, resp, 1383)
	})

	t.Run("CREATE INDEX BINARY", func(t *testing.T) {
		qCreateIndex, err := parser.New("CREATE INDEX userId ON users USING BINARY (id);")
		require.NoError(t, err)

		err = store.CreateIndex(qCreateIndex.GetQuery())
		require.NoError(t, err)
	})
}
