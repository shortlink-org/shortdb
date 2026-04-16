package repl

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

//nolint:ireturn,gocognit,varnamelen // tea.Model; branching TUI + observable
func (m *tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case workspaceTabPickMsg:
		if shellTab(msg) == shellTabCLI {
			cmd := m.switchToCLITab()

			return m, cmd
		}

		m.openObservableTab(false)

		return m, m.obsInput.Focus()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.initialized = true
		m.theme = newTheme(m.width)
		m.applyChromeStyles()
		m.layout()
		m.syncViewport()

		return m, nil

	case tea.MouseWheelMsg:
		if m.activeTab == shellTabObservable {
			var cmd tea.Cmd

			m.obsTable, cmd = m.obsTable.Update(msg)

			return m, cmd
		}

		var cmd tea.Cmd

		m.vp, cmd = m.vp.Update(msg)

		return m, cmd

	case tea.KeyPressMsg:
		return m.updateKeyPress(msg)
	}

	return m, nil
}

//nolint:ireturn // tea.Model contract
func (m *tuiModel) updateKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if isSelectCLIKeyMsg(msg) && m.activeTab != shellTabCLI {
		cmd := m.switchToCLITab()

		return m, cmd
	}

	if isSelectObservableKeyMsg(msg) && m.activeTab != shellTabObservable {
		m.openObservableTab(false)

		return m, m.obsInput.Focus()
	}

	if m.activeTab == shellTabObservable {
		return m.updateObservableKeyPress(msg)
	}

	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit

	case "enter":
		line := strings.TrimSpace(m.ti.Value())
		m.ti.Reset()
		m.histIndex = -1
		m.ti.ShowSuggestions = true

		if line == "" {
			return m, m.ti.Focus()
		}

		m.appendLines(m.theme.muted.Render("> " + line))

		if line == ".tables" {
			m.openObservableTab(true)

			return m, m.obsInput.Focus()
		}

		out, quit := m.repl.handleREPLLine(line)
		if len(out) > 0 {
			m.appendLines(out...)
		}

		if quit {
			return m, tea.Quit
		}

		return m, m.ti.Focus()
	}

	if keyScrollsTranscript(msg) {
		var cmd tea.Cmd

		m.vp, cmd = m.vp.Update(msg)

		return m, cmd
	}

	if m.tryHistoryKeys(msg) {
		return m, m.ti.Focus()
	}

	m.clearHistoryBrowseIfNeeded(msg)

	var cmd tea.Cmd

	m.ti, cmd = m.ti.Update(msg)

	return m, cmd
}
