//nolint:testpackage // this package is for testing only
package tests

import (
	"testing"

	parser "github.com/shortlink-org/shortdb/shortdb/parser/v1"
	"github.com/stretchr/testify/require"
)

func BenchmarkParser(b *testing.B) {
	b.Run("CREATE TABLE", func(b *testing.B) {
		for range b.N {
			_, err := parser.New("CREATE TABLE users ( id integer, name text );")
			require.NoError(b, err)
		}
	})

	b.Run("SELECT", func(b *testing.B) {
		for range b.N {
			_, err := parser.New("SELECT a, c, d FROM 'b' WHERE a != '1' LIMIT 5")
			require.NoError(b, err)
		}
	})

	b.Run("INSERT INTO", func(b *testing.B) {
		for range b.N {
			_, err := parser.New("INSERT INTO 'a' (b,c,d) VALUES ('1','2','3'),('4','5','6');")
			require.NoError(b, err)
		}
	})
}
