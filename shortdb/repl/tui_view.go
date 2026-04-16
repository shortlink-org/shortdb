package repl

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *tuiModel) View() tea.View {
	if !m.initialized {
		lightDark := lipgloss.LightDark(lipgloss.HasDarkBackground(os.Stdin, os.Stdout))
		s := lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("#71717a"), lipgloss.Color("#a1a1aa"))).Render("…")

		return tea.NewView(s)
	}

	obs := m.activeTab == shellTabObservable
	placed, _ := m.shellPlaced(obs)

	rootView := tea.NewView(placed)
	// CellMotion captures mouse events for the whole terminal, which breaks
	// click-drag text selection in most terminals. Keep mouse only on Observable
	// (tab strip clicks); CLI uses f1/f2 for tabs.
	if obs {
		rootView.MouseMode = tea.MouseModeCellMotion
		rootView.WindowTitle = "ShortDB · observable"
	} else {
		rootView.MouseMode = tea.MouseModeNone
		rootView.WindowTitle = "ShortDB"
	}

	rootView.OnMouse = m.shellMouseHandler()

	return rootView
}
