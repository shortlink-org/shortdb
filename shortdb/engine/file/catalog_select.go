package file

import (
	"fmt"
	"sort"
	"strings"

	field "github.com/shortlink-org/shortdb/shortdb/domain/field/v1"
	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
	query "github.com/shortlink-org/shortdb/shortdb/domain/query/v1"
	table "github.com/shortlink-org/shortdb/shortdb/domain/table/v1"
)

// virtualDiscoveryTable is the unquoted catalog name ([a-zA-Z0-9]+ only) for FROM.
const virtualDiscoveryTable = "shortdbcatalog"

// virtualDiscoveryTableQuoted is the same catalog when FROM uses a quoted identifier (underscores allowed).
const virtualDiscoveryTableQuoted = "shortdb_tables"

func isDiscoveryCatalogTable(name string) bool {
	return strings.EqualFold(name, virtualDiscoveryTable) ||
		strings.EqualFold(name, virtualDiscoveryTableQuoted)
}

func isReservedDiscoveryTableName(name string) bool {
	return isDiscoveryCatalogTable(name)
}

func fieldTypeLabel(t field.Type) string {
	switch t {
	case field.Type_TYPE_INTEGER:
		return "integer"
	case field.Type_TYPE_STRING:
		return "string"
	case field.Type_TYPE_BOOLEAN:
		return "boolean"
	default:
		return "unknown"
	}
}

func formatTableColumns(tbl *table.Table) string {
	if tbl == nil || len(tbl.GetFields()) == 0 {
		return ""
	}

	names := make([]string, 0, len(tbl.GetFields()))
	for n := range tbl.GetFields() {
		names = append(names, n)
	}

	sort.Strings(names)

	parts := make([]string, 0, len(names))
	for _, n := range names {
		parts = append(parts, fmt.Sprintf("%s %s", n, fieldTypeLabel(tbl.GetFields()[n])))
	}

	return strings.Join(parts, ", ")
}

func (f *File) selectDiscoveryCatalog(in *query.Query) ([]*page.Row, error) {
	filterName := ""
	for _, cond := range in.GetConditions() {
		if !cond.GetLValueIsField() {
			return nil, fmt.Errorf("at SELECT: invalid WHERE on catalog %q", virtualDiscoveryTableQuoted)
		}

		if !strings.EqualFold(cond.GetLValue(), "name") {
			return nil, fmt.Errorf("at SELECT: catalog WHERE only supports column %q (got %q)", "name", cond.GetLValue())
		}

		if cond.GetOperator() != query.Operator_OPERATOR_EQ {
			return nil, fmt.Errorf("at SELECT: catalog WHERE only supports = for name")
		}

		// The parser classifies bare RHS tokens as field references; for the discovery catalog,
		// name = books and name = 'books' are both treated as a literal table name filter.
		rhs := cond.GetRValue()
		if rhs == "" {
			return nil, fmt.Errorf("at SELECT: catalog WHERE needs a value for name")
		}

		if filterName != "" && !strings.EqualFold(filterName, rhs) {
			return nil, fmt.Errorf("at SELECT: catalog supports a single name = '...' filter")
		}

		filterName = rhs
	}

	if len(in.GetFields()) == 0 {
		return nil, ErrIncorrectNameFields
	}

	outFields := append([]string(nil), in.GetFields()...)
	if len(outFields) == 1 && outFields[0] == "*" {
		outFields = []string{"name", "columns"}
	}

	valid := map[string]struct{}{"name": {}, "columns": {}}
	for _, fld := range outFields {
		if _, ok := valid[fld]; !ok {
			return nil, &IncorrectNameFieldsError{Field: fld, Table: virtualDiscoveryTableQuoted}
		}
	}

	tables := f.database.GetTables()
	if len(tables) == 0 {
		return []*page.Row{}, nil
	}

	names := make([]string, 0, len(tables))
	for n := range tables {
		if isDiscoveryCatalogTable(n) {
			continue
		}

		names = append(names, n)
	}

	sort.Strings(names)

	response := make([]*page.Row, 0, len(names))

	for _, tblName := range names {
		if filterName != "" && !strings.EqualFold(tblName, filterName) {
			continue
		}

		tbl := tables[tblName]
		row := &page.Row{Value: make(map[string][]byte)}

		for _, fld := range outFields {
			switch fld {
			case "name":
				row.Value["name"] = []byte(tblName)
			case "columns":
				row.Value["columns"] = []byte(formatTableColumns(tbl))
			}
		}

		response = append(response, row)

		if in.IsLimit() {
			in.Limit--
		}

		if !in.IsLimit() {
			break
		}
	}

	return response, nil
}
