package repl

import tea "charm.land/bubbletea/v2"

// tryHistoryKeys handles ↑/↓ for session history only when inline completions are inactive.
func (m *tuiModel) tryHistoryKeys(msg tea.KeyPressMsg) bool {
	key := msg.String()
	if key != "up" && key != "down" {
		return false
	}

	if len(m.ti.MatchedSuggestions()) > 0 {
		return false
	}

	if key == "up" {
		return m.historyUp()
	}

	return m.historyDown()
}

func (m *tuiModel) clearHistoryBrowseIfNeeded(msg tea.KeyPressMsg) {
	if m.histIndex < 0 {
		return
	}

	switch msg.String() {
	case "up", "down":
		return
	default:
		m.histIndex = -1
		m.ti.ShowSuggestions = true
	}
}

func (m *tuiModel) historyUp() bool {
	hist := m.repl.session.GetHistory()
	if len(hist) == 0 {
		return false
	}

	switch {
	case m.histIndex < 0:
		m.histIndex = len(hist) - 1
	case m.histIndex > 0:
		m.histIndex--
	default:
		return true
	}

	m.ti.ShowSuggestions = false
	m.ti.SetValue(hist[m.histIndex])

	return true
}

func (m *tuiModel) historyDown() bool {
	hist := m.repl.session.GetHistory()
	if len(hist) == 0 || m.histIndex < 0 {
		return false
	}

	if m.histIndex < len(hist)-1 {
		m.histIndex++
		m.ti.ShowSuggestions = false
		m.ti.SetValue(hist[m.histIndex])

		return true
	}

	m.histIndex = -1
	m.ti.ShowSuggestions = true
	m.ti.SetValue("")

	return true
}
