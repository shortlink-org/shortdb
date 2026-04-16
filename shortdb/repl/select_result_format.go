package repl

import (
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
)

const (
	transcriptMaxDataRows  = 2000
	transcriptMaxCellRunes = 80
)

func sortedFieldKeys(rows []*page.Row) []string {
	set := map[string]struct{}{}

	for _, r := range rows {
		if r == nil || r.GetValue() == nil {
			continue
		}

		for k := range r.GetValue() {
			set[k] = struct{}{}
		}
	}

	out := make([]string, 0, len(set))

	for k := range set {
		out = append(out, k)
	}

	slices.Sort(out)

	return out
}

func truncateTranscriptCell(val string) string {
	if utf8.RuneCountInString(val) <= transcriptMaxCellRunes {
		return val
	}

	r := []rune(val)

	return string(r[:transcriptMaxCellRunes]) + "…"
}

func formatSelectTranscriptLines(rows []*page.Row) []string {
	if len(rows) == 0 {
		return []string{"(0 rows)"}
	}

	keys := sortedFieldKeys(rows)
	if len(keys) == 0 {
		return []string{"(rows present but no column values)"}
	}

	header := make([]string, len(keys))
	for i, k := range keys {
		header[i] = truncateTranscriptCell(k)
	}

	lines := []string{strings.Join(header, "\t")}

	nShow := len(rows)
	remainder := 0

	if nShow > transcriptMaxDataRows {
		remainder = nShow - transcriptMaxDataRows
		nShow = transcriptMaxDataRows
	}

	for i := range nShow {
		r := rows[i]
		cells := make([]string, len(keys))

		for j, k := range keys {
			val := ""
			if r != nil && r.GetValue() != nil && r.GetValue()[k] != nil {
				val = string(r.GetValue()[k])
			}

			cells[j] = truncateTranscriptCell(val)
		}

		lines = append(lines, strings.Join(cells, "\t"))
	}

	if remainder > 0 {
		lines = append(lines, fmt.Sprintf("… and %d more rows (showing first %d)", remainder, transcriptMaxDataRows))
	}

	return lines
}

func formatExecTranscript(response any) []string {
	rows, ok := response.([]*page.Row)
	if ok {
		return formatSelectTranscriptLines(rows)
	}

	return []string{fmt.Sprint(response)}
}
