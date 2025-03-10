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
)

func BenchmarkEngine(b *testing.B) {
	ctx, cancel := context.WithCancel(b.Context())

	// set engine
	path := "/tmp/shortdb_test_unit"

	store, err := engine.New(ctx, "file", file.SetName("testDatabase"), file.SetPath(path))
	require.NoError(b, err)

	b.Cleanup(func() {
		cancel()

		err = os.RemoveAll(path)
		require.NoError(b, err)
	})

	b.Run("CREATE TABLE", func(b *testing.B) {
		for range b.N {
			qCreateTable, errParserNew := parser.New("create table users (id integer, name string, active bool)")
			require.NoError(b, errParserNew)

			_, err = store.Exec(qCreateTable.GetQuery())
			require.NoError(b, err)
		}
	})

	b.Run("INSERT INTO USERS", func(b *testing.B) {
		for i := range b.N {
			qInsertUsers, errParserNew := parser.New(fmt.Sprintf("insert into users ('id', 'name', 'active') VALUES ('%d', 'Ivan', 'false')", i))
			require.NoError(b, errParserNew)

			errInsert := store.Insert(qInsertUsers.GetQuery())
			require.NoError(b, errInsert)
		}

		// save data
		err = store.Close()
		require.NoError(b, err)
	})

	b.Run("SELECT USERS", func(b *testing.B) {
		for range b.N {
			qInsertUsers, err := parser.New("select id, name, active from users limit 5")
			require.NoError(b, err)

			resp, err := store.Select(qInsertUsers.GetQuery())
			require.NoError(b, err)
			assert.Len(b, resp, 5)
		}
	})

	b.Run("SELECT USERS WITH WHERE id=99 AND LIMIT 2", func(b *testing.B) {
		for range b.N {
			qSelectUsers, err := parser.New("select id, name, active from users where id='99' limit 2")
			require.NoError(b, err)

			_, err = store.Select(qSelectUsers.GetQuery())
			require.NoError(b, err)
		}
	})

	b.Run("SELECT USERS FULL SCAN", func(b *testing.B) {
		for range b.N {
			qSelectUsers, err := parser.New("select id, name, active from users")
			require.NoError(b, err)

			_, err = store.Select(qSelectUsers.GetQuery())
			require.NoError(b, err)
		}
	})

	b.Run("CREATE INDEX BTREE", func(b *testing.B) {
		for range b.N {
			qCreateIndex, err := parser.New("CREATE INDEX userId ON users USING BTREE (id);")
			require.NoError(b, err)

			err = store.CreateIndex(qCreateIndex.GetQuery())
			require.NoError(b, err)
		}
	})
}
