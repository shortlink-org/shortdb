package repl

import (
	"strconv"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/paginator"
	tea "charm.land/bubbletea/v2"

	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
)

const (
	obsNavStackMax             = 32
	obsSQLHistMax              = 64
	obsPerPageHeader           = 2
	obsPaginatorInitialPerPage = 10
	obsAltUpKey                = "alt+up"
	obsAltDownKey              = "alt+down"
)

type obsNavFrame struct {
	browseMode obsBrowseMode
	rows       []*page.Row
	inputVal   string
	status     string
	statusErr  bool
	focusTable bool
	page       int
}

func newObsPaginator() paginator.Model {
	p := paginator.New(paginator.WithPerPage(obsPaginatorInitialPerPage))
	p.KeyMap = paginator.KeyMap{
		PrevPage: key.NewBinding(key.WithKeys("pgup")),
		NextPage: key.NewBinding(key.WithKeys("pgdown")),
	}

	return p
}

func obsRowsPerPage(tblHeight int) int {
	return max(3, tblHeight-obsPerPageHeader)
}

func (m *tuiModel) syncObsPaginatorForRows(tblHeight, rowCount int) {
	m.obsPaginator.PerPage = obsRowsPerPage(tblHeight)

	if rowCount < 1 {
		m.obsPaginator.TotalPages = 1
		m.obsPaginator.Page = 0

		return
	}

	totalPages := m.obsPaginator.SetTotalPages(rowCount)
	if m.obsPaginator.Page >= totalPages {
		m.obsPaginator.Page = max(0, totalPages-1)
	}
}

func (m *tuiModel) obsRowRangeSuffix() string {
	rowCount := len(m.obsAllRows)
	if rowCount == 0 || m.obsPaginator.TotalPages <= 1 {
		return ""
	}

	start, end := m.obsPaginator.GetSliceBounds(rowCount)
	lo := start + 1
	hi := max(lo, end)

	return " · showing " + strconv.Itoa(lo) + "–" + strconv.Itoa(hi) + " of " + strconv.Itoa(rowCount)
}

func (m *tuiModel) captureObsNavFrame() obsNavFrame {
	return obsNavFrame{
		browseMode: m.obsBrowseMode,
		rows:       append([]*page.Row(nil), m.obsAllRows...),
		inputVal:   m.obsInput.Value(),
		status:     m.obsStatus,
		statusErr:  m.obsStatusErr,
		focusTable: m.obsFocusTable,
		page:       m.obsPaginator.Page,
	}
}

func (m *tuiModel) pushObsNav() {
	if len(m.obsStack) >= obsNavStackMax {
		m.obsStack = m.obsStack[1:]
	}

	m.obsStack = append(m.obsStack, m.captureObsNavFrame())
}

func (m *tuiModel) popObsNav() bool {
	if len(m.obsStack) == 0 {
		return false
	}

	i := len(m.obsStack) - 1
	fr := m.obsStack[i]
	m.obsStack = m.obsStack[:i]
	m.restoreObsNavFrame(&fr)

	return true
}

func (m *tuiModel) restoreObsNavFrame(frame *obsNavFrame) {
	m.obsBrowseMode = frame.browseMode
	m.obsAllRows = append([]*page.Row(nil), frame.rows...)
	m.obsInput.SetValue(frame.inputVal)
	m.obsStatus = frame.status
	m.obsStatusErr = frame.statusErr
	m.obsFocusTable = frame.focusTable
	m.obsPaginator.Page = frame.page
	m.layout()
	m.applyObsTableFocus()
}

func (m *tuiModel) rebuildObsTableSlice(tblWidth, tblHeight int) {
	if m.obsBrowseMode == obsBrowseCatalog {
		slice := m.obsAllRows
		if len(slice) > 0 {
			start, end := m.obsPaginator.GetSliceBounds(len(slice))
			slice = slice[start:end]
		}

		m.obsTable = newCatalogTableModel(tblWidth, tblHeight, catalogRowsFromPages(slice))

		return
	}

	var pageSlice []*page.Row

	if len(m.obsAllRows) > 0 {
		start, end := m.obsPaginator.GetSliceBounds(len(m.obsAllRows))
		pageSlice = m.obsAllRows[start:end]
	}

	m.obsTable = newDynamicDataTable(tblWidth, tblHeight, pageSlice)
}

func (m *tuiModel) appendObsSQLHistory(sqlLine string) {
	sqlLine = strings.TrimSpace(sqlLine)
	if sqlLine == "" {
		return
	}

	if len(m.obsSQLHistory) > 0 && m.obsSQLHistory[len(m.obsSQLHistory)-1] == sqlLine {
		return
	}

	m.obsSQLHistory = append(m.obsSQLHistory, sqlLine)
	if len(m.obsSQLHistory) > obsSQLHistMax {
		m.obsSQLHistory = m.obsSQLHistory[len(m.obsSQLHistory)-obsSQLHistMax:]
	}
}

func (m *tuiModel) tryObsSQLHistoryKeys(msg tea.KeyPressMsg) bool {
	keyStr := msg.String()
	if keyStr != obsAltUpKey && keyStr != obsAltDownKey {
		return false
	}

	hist := m.obsSQLHistory
	if len(hist) == 0 {
		return false
	}

	if keyStr == obsAltUpKey {
		switch {
		case m.obsHistBrowseIdx < 0:
			m.obsHistBrowseIdx = len(hist) - 1
		case m.obsHistBrowseIdx > 0:
			m.obsHistBrowseIdx--
		default:
			return true
		}
	} else {
		if m.obsHistBrowseIdx < 0 {
			return false
		}

		if m.obsHistBrowseIdx < len(hist)-1 {
			m.obsHistBrowseIdx++
		} else {
			m.obsHistBrowseIdx = -1
			m.obsInput.SetValue("")

			return true
		}
	}

	m.obsInput.SetValue(hist[m.obsHistBrowseIdx])

	return true
}

func (m *tuiModel) clearObsHistBrowseIfNeeded(msg tea.KeyPressMsg) {
	if m.obsHistBrowseIdx < 0 {
		return
	}

	switch msg.String() {
	case obsAltUpKey, obsAltDownKey:
		return
	default:
		m.obsHistBrowseIdx = -1
	}
}

func (m *tuiModel) obsStatusBlock(thm *tuiTheme, innerW int) string {
	var lines []string

	if m.obsStatus != "" {
		line := m.obsStatus
		if !m.obsStatusErr {
			line += m.obsRowRangeSuffix()
		}

		style := thm.subtle.Width(innerW)
		if m.obsStatusErr {
			style = thm.obsErr.Width(innerW)
		}

		lines = append(lines, style.Render(line))
	}

	if m.obsPaginator.TotalPages > 1 && len(m.obsAllRows) > 0 {
		lines = append(lines, thm.muted.Width(innerW).Render(m.obsPaginator.View()+" · pgup/pgdn"))
	}

	return strings.Join(lines, "\n")
}
