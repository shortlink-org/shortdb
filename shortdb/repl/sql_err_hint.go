package repl

import (
	"errors"
	"strings"
)

// observableSQLMessage turns engine/parser errors into short hints for the Observable tab.
func observableSQLMessage(err error) string {
	if err == nil {
		return ""
	}

	switch {
	case errors.Is(err, errObservableEmptySQL):
		return "Enter a SELECT statement. A trailing semicolon is optional."
	case errors.Is(err, errObservableNotSelect):
		return "Only SELECT is allowed in Observable."
	case errors.Is(err, errUnexpectedCatalogResponse):
		return "Unexpected response from the database engine."
	}

	errText := err.Error()

	switch {
	case strings.HasPrefix(errText, "parse:"):
		return "Parse error: " + strings.TrimSpace(strings.TrimPrefix(errText, "parse:")) + " Check syntax; end with ; if unsure."
	case strings.HasPrefix(errText, "exec:"):
		return "Query failed: " + strings.TrimSpace(strings.TrimPrefix(errText, "exec:"))
	case strings.HasPrefix(errText, "catalog:"):
		return "Catalog failed: " + strings.TrimSpace(strings.TrimPrefix(errText, "catalog:"))
	default:
		return errText
	}
}
