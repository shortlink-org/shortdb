package repl

import (
	"errors"
	"fmt"
	"strings"

	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
	query "github.com/shortlink-org/shortdb/shortdb/domain/query/v1"
	parser "github.com/shortlink-org/shortdb/shortdb/parser/v1"
)

var (
	errUnexpectedCatalogResponse = errors.New("unexpected catalog response type")
	errObservableEmptySQL        = errors.New("enter a SQL statement")
	errObservableNotSelect       = errors.New("only SELECT is allowed in the Observable tab")
)

func (r *Repl) fetchShortdbCatalog() ([]*page.Row, error) {
	// Quoted name is required: the lexer treats '_' as terminating an unquoted
	// identifier, so FROM shortdb_tables would parse as table "shortdb".
	rows, err := r.execSelectRows("SELECT name, columns FROM 'shortdb_tables'")
	if err != nil {
		return nil, fmt.Errorf("catalog: %w", err)
	}

	return rows, nil
}

func (r *Repl) execSelectRows(sql string) ([]*page.Row, error) {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil, errObservableEmptySQL
	}

	if !strings.HasSuffix(sql, ";") {
		sql += ";"
	}

	prsr, err := parser.New(sql)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	if prsr.GetQuery().GetType() != query.Type_TYPE_SELECT {
		return nil, errObservableNotSelect
	}

	resp, err := r.engine.Exec(prsr.GetQuery())
	if err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}

	if resp == nil {
		return []*page.Row{}, nil
	}

	rows, ok := resp.([]*page.Row)
	if !ok {
		return nil, fmt.Errorf("%w: %T", errUnexpectedCatalogResponse, resp)
	}

	return rows, nil
}
