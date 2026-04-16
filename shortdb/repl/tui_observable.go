package repl

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
)

const (
	obsDrillRowLimit          = 200
	obsErrorTableCellMaxChars = 120
)

func newObservableSQLInput() textinput.Model {
	textIn := textinput.New()
	textIn.Prompt = "sql> "
	textIn.Placeholder = "SELECT name, columns FROM 'shortdb_tables'"
	textIn.ShowSuggestions = false
	textIn.CharLimit = inputCharLimit

	hasDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	if hasDark {
		textIn.SetStyles(textinput.DefaultDarkStyles())
	} else {
		textIn.SetStyles(textinput.DefaultLightStyles())
	}

	lightDark := lipgloss.LightDark(hasDark)
	tiStyles := textIn.Styles()
	tiStyles.Cursor.Shape = tea.CursorBar
	tiStyles.Cursor.Color = lightDark(lipgloss.Color("#18181b"), lipgloss.Color("#fafafa"))
	tiStyles.Focused.Prompt = lipgloss.NewStyle().Bold(true).Foreground(lightDark(lipgloss.Color("#5b21b6"), lipgloss.Color("#c4b5fd")))
	tiStyles.Focused.Placeholder = lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("#a1a1aa"), lipgloss.Color("#71717a")))
	tiStyles.Blurred = tiStyles.Focused
	textIn.SetStyles(tiStyles)

	return textIn
}

func (m *tuiModel) workspaceTabLine() string {
	theme := m.theme
	cli := theme.subtle.Render(" CLI ")
	obs := theme.subtle.Render(" Observable ")

	switch m.activeTab {
	case shellTabCLI:
		cli = theme.accent.Render(" CLI ")
	case shellTabObservable:
		obs = theme.accent.Render(" Observable ")
	default:
		// only CLI and Observable tabs exist
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cli, theme.muted.Render(" │ "), obs)
}

func equalColumnWidths(totalWidth, colCount int) []int {
	if colCount <= 0 {
		return nil
	}

	avail := max(colCount*minDataColumnWidth, totalWidth-catalogLayoutPad)
	base := max(minDataColumnWidth, avail/colCount)
	out := make([]int, colCount)

	for i := range out {
		out[i] = base
	}

	for rem := avail - base*colCount; rem > 0; rem-- {
		out[rem%colCount]++
	}

	return out
}

func newDynamicDataTable(width, height int, pages []*page.Row) table.Model {
	if len(pages) == 0 {
		return newDataTableModel(width, height,
			[]table.Column{{Title: " ", Width: max(minDataColumnWidth, width-catalogLayoutPad)}},
			nil)
	}

	keys := sortedFieldKeys(pages)
	if len(keys) == 0 {
		return newDataTableModel(width, height,
			[]table.Column{{Title: " ", Width: max(minDataColumnWidth, width-catalogLayoutPad)}},
			nil)
	}

	trows := make([]table.Row, len(pages))

	for i, pg := range pages {
		line := make([]string, len(keys))

		for j, k := range keys {
			if pg != nil && pg.GetValue() != nil && pg.GetValue()[k] != nil {
				line[j] = string(pg.GetValue()[k])
			}
		}

		trows[i] = line
	}

	widths := equalColumnWidths(max(catalogMinInnerWidth, width), len(keys))
	columns := make([]table.Column, len(keys))

	for i, k := range keys {
		columns[i] = table.Column{Title: k, Width: widths[i]}
	}

	return newDataTableModel(width, height, columns, trows)
}

func (m *tuiModel) observableTableSize() (int, int) {
	tblWidth := max(catalogMinInnerWidth, m.theme.innerWidth-catalogInnerWidthMargin)
	outerV := m.theme.outer.GetVerticalFrameSize()
	innerH := m.height - outerV

	const (
		titleLine  = 1
		divider    = 1
		tabBar     = 1
		hintLine   = 1
		gap        = 1
		inputBlock = 2
		statusLine = 1
		footer     = footerRowCount
		bodyGaps   = 3
	)

	fixed := titleLine + divider + tabBar + hintLine + gap + inputBlock + statusLine + footer + bodyGaps
	tblHeight := max(catalogMinInnerHeight, innerH-fixed)

	return tblWidth, tblHeight
}

func (m *tuiModel) reloadObservableCatalog() {
	rows, err := m.repl.fetchShortdbCatalog()
	tblWidth, tblHeight := m.observableTableSize()

	m.obsAllRows = rows
	m.obsStatusErr = err != nil

	switch {
	case err != nil:
		m.obsStatus = observableSQLMessage(err)
	case len(rows) == 0:
		m.obsStatus = ""
	default:
		m.obsStatus = fmt.Sprintf("%d table(s) in catalog", len(rows))
	}

	m.obsBrowseMode = obsBrowseCatalog
	m.obsPaginator.Page = 0
	m.syncObsPaginatorForRows(tblHeight, len(m.obsAllRows))
	m.rebuildObsTableSlice(tblWidth, tblHeight)
}

func (m *tuiModel) applyObsTableFocus() {
	if m.obsFocusTable {
		m.obsTable.Focus()
		m.obsInput.Blur()
	} else {
		m.obsTable.Blur()
		m.obsInput.Focus()
	}
}

func (m *tuiModel) openObservableTab(appendHistory bool) {
	m.activeTab = shellTabObservable
	m.obsFocusTable = false
	m.obsHistBrowseIdx = -1
	m.reloadObservableCatalog()
	m.applyObsTableFocus()

	if appendHistory {
		m.repl.session.Raw = ""
		m.repl.session.Exec = true
		m.repl.session.History = append(m.repl.session.GetHistory(), ".tables")
	}

	m.ti.Blur()
	m.layout()
}

func (m *tuiModel) switchToCLITab() tea.Cmd {
	m.activeTab = shellTabCLI
	m.obsStack = nil
	m.obsHistBrowseIdx = -1
	m.obsInput.Blur()
	m.obsTable.Blur()
	m.obsFocusTable = false
	m.histIndex = -1
	m.ti.ShowSuggestions = true
	m.layout()
	m.syncViewport()

	return m.ti.Focus()
}

func (m *tuiModel) handleObservableEsc() tea.Cmd {
	if m.obsBrowseMode == obsBrowseCatalog {
		return m.switchToCLITab()
	}

	m.obsStack = nil
	m.obsInput.SetValue("")
	m.reloadObservableCatalog()
	m.obsFocusTable = false
	m.applyObsTableFocus()
	m.layout()

	return m.obsInput.Focus()
}

func (m *tuiModel) runObservableQuery() {
	tblWidth, tblHeight := m.observableTableSize()

	rows, err := m.repl.execSelectRows(m.obsInput.Value())
	if err != nil {
		m.obsHistBrowseIdx = -1
		m.obsStatusErr = true
		m.obsStatus = observableSQLMessage(err)
		m.obsBrowseMode = obsBrowseQuery
		m.obsAllRows = nil
		m.syncObsPaginatorForRows(tblHeight, 0)

		errCell := m.obsStatus
		if len(errCell) > obsErrorTableCellMaxChars {
			errCell = errCell[:obsErrorTableCellMaxChars] + "…"
		}

		m.obsTable = newCatalogTableModel(tblWidth, tblHeight, []table.Row{{"(error)", errCell}})
		m.applyObsTableFocus()
		m.layout()

		return
	}

	m.pushObsNav()
	m.obsPaginator.Page = 0
	m.appendObsSQLHistory(m.obsInput.Value())
	m.obsHistBrowseIdx = -1
	m.obsStatusErr = false
	m.obsAllRows = rows
	m.obsBrowseMode = obsBrowseQuery
	m.syncObsPaginatorForRows(tblHeight, len(rows))
	m.obsStatus = fmt.Sprintf("%d row(s)", len(rows))
	m.rebuildObsTableSlice(tblWidth, tblHeight)
	m.applyObsTableFocus()
	m.layout()
}

func isSafeDrillTableName(name string) bool {
	if name == "" || len(name) > maxDrillTableNameLen {
		return false
	}

	for _, c := range name {
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '_':
		default:
			return false
		}
	}

	if name[0] >= '0' && name[0] <= '9' {
		return false
	}

	ln := strings.ToLower(name)
	if ln == "shortdb_tables" || ln == "shortdbcatalog" {
		return false
	}

	return true
}

func (m *tuiModel) drillObservableSelectedTable() {
	if m.obsBrowseMode != obsBrowseCatalog {
		return
	}

	row := m.obsTable.SelectedRow()
	if len(row) == 0 {
		return
	}

	name := strings.TrimSpace(row[0])
	if name == "" || strings.HasPrefix(name, "(") {
		return
	}

	if !isSafeDrillTableName(name) {
		m.obsStatusErr = false
		m.obsStatus = "cannot open this row as a table"

		return
	}

	escaped := strings.ReplaceAll(name, "'", "''")
	//nolint:unqueryvet // engine expands * to real columns at SELECT execution.
	sql := fmt.Sprintf("SELECT * FROM '%s' LIMIT %d;", escaped, obsDrillRowLimit)
	m.obsInput.SetValue(strings.TrimSuffix(sql, ";"))

	rows, err := m.repl.execSelectRows(sql)

	tblWidth, tblHeight := m.observableTableSize()
	if err != nil {
		m.obsStatusErr = true
		m.obsStatus = observableSQLMessage(err)

		return
	}

	m.pushObsNav()
	m.obsPaginator.Page = 0
	m.appendObsSQLHistory(sql)
	m.obsHistBrowseIdx = -1
	m.obsStatusErr = false
	m.obsAllRows = rows
	m.syncObsPaginatorForRows(tblHeight, len(rows))
	m.obsStatus = fmt.Sprintf("table %q · %d row(s) · esc — catalog · u — back", name, len(rows))
	m.rebuildObsTableSlice(tblWidth, tblHeight)
	m.obsBrowseMode = obsBrowseQuery
	m.obsFocusTable = true
	m.obsTable.Focus()
	m.obsInput.Blur()
	m.layout()
}

//nolint:ireturn // tea.Model contract
func (m *tuiModel) updateObservableKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "f2" {
		return m, nil
	}

	key := msg.String()

	switch key {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		cmd := m.handleObservableEsc()

		return m, cmd
	case "tab":
		m.obsHistBrowseIdx = -1
		m.obsFocusTable = !m.obsFocusTable
		m.applyObsTableFocus()

		return m, nil
	case "u":
		if m.obsFocusTable {
			if m.popObsNav() {
				return m, nil
			}

			return m, nil
		}
	default:
	}

	if m.obsFocusTable && m.obsPaginator.TotalPages > 1 && (key == "pgup" || key == "pgdown") {
		m.obsPaginator, _ = m.obsPaginator.Update(msg)
		m.layout()

		return m, nil
	}

	if key == "enter" {
		if !m.obsFocusTable {
			m.runObservableQuery()

			return m, m.obsInput.Focus()
		}

		if m.obsBrowseMode == obsBrowseCatalog {
			m.drillObservableSelectedTable()

			return m, nil
		}

		var cmd tea.Cmd

		m.obsTable, cmd = m.obsTable.Update(msg)

		return m, cmd
	}

	if m.obsFocusTable {
		var cmd tea.Cmd

		m.obsTable, cmd = m.obsTable.Update(msg)

		return m, cmd
	}

	m.clearObsHistBrowseIfNeeded(msg)

	if m.tryObsSQLHistoryKeys(msg) {
		return m, nil
	}

	var cmd tea.Cmd

	m.obsInput, cmd = m.obsInput.Update(msg)

	return m, cmd
}
