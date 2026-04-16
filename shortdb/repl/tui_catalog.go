package repl

import (
	"unicode/utf8"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"

	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
)

const (
	catalogLayoutPad        = 6
	catalogMinInnerWidth    = 20
	catalogMinInnerHeight   = 6
	catalogMinNameCol       = 8
	catalogMinColumnsCol    = 16
	catalogInnerWidthMargin = 2
	catalogVPHeightBoost    = 4
	catalogMinViewportBody  = 8
)

var catalogTableChrome = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func catalogRowsFromPages(rows []*page.Row) []table.Row {
	out := make([]table.Row, 0, len(rows))
	for _, row := range rows {
		name := ""
		cols := ""

		if row != nil && row.GetValue() != nil {
			if b := row.GetValue()["name"]; len(b) > 0 {
				name = string(b)
			}

			if b := row.GetValue()["columns"]; len(b) > 0 {
				cols = string(b)
			}
		}

		out = append(out, table.Row{name, cols})
	}

	return out
}

func catalogColumnWidths(width int, rows []table.Row) (int, int) {
	avail := width - catalogLayoutPad
	avail = max(avail, catalogMinNameCol+catalogMinColumnsCol)

	nameW, colsW := catalogMinNameCol, catalogMinColumnsCol

	for _, r := range rows {
		if len(r) > 0 {
			nameW = max(nameW, utf8.RuneCountInString(r[0])+1)
		}

		if len(r) > 1 {
			colsW = max(colsW, utf8.RuneCountInString(r[1])+1)
		}
	}

	if nameW+colsW > avail {
		extra := nameW + colsW - avail
		takeName := min(extra, max(0, nameW-catalogMinNameCol))
		nameW -= takeName
		extra -= takeName
		colsW = max(catalogMinColumnsCol, colsW-extra)
	}

	return nameW, colsW
}

func newDataTableModel(width, height int, columns []table.Column, data []table.Row) table.Model {
	if len(columns) == 0 {
		columns = []table.Column{{Title: " ", Width: max(minDataColumnWidth, width-catalogLayoutPad)}}
	}

	if len(data) == 0 {
		row := make(table.Row, len(columns))
		for i := range row {
			row[i] = " "
		}

		row[0] = "(no rows)"

		data = []table.Row{row}
	}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(data),
		table.WithFocused(true),
		table.WithHeight(max(catalogMinInnerHeight, max(1, height))),
		table.WithWidth(max(catalogMinInnerWidth, max(1, width))),
	)

	tblStyles := table.DefaultStyles()
	tblStyles.Header = tblStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	tblStyles.Selected = tblStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	tbl.SetStyles(tblStyles)

	return tbl
}

func newCatalogTableModel(width, height int, data []table.Row) table.Model {
	if len(data) == 0 {
		data = []table.Row{{"(no tables yet)", ""}}
	}

	nameW, colsW := catalogColumnWidths(max(catalogMinInnerWidth, width), data)

	columns := []table.Column{
		{Title: "name", Width: nameW},
		{Title: "columns", Width: colsW},
	}

	return newDataTableModel(width, height, columns, data)
}
