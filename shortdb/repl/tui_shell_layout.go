package repl

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// workspaceTabPickMsg is sent when the user clicks a workspace tab (View.OnMouse).
type workspaceTabPickMsg shellTab

// shellPlaced returns the full placed shell string and the terminal row index of the tab line (-1 if unknown).
func (m *tuiModel) shellPlaced(observable bool) (string, int) {
	thm := m.theme
	innerW := thm.innerWidth

	header := thm.accent.Width(innerW).Render("ShortDB")
	dividerLine := thm.divider.Render(strings.Repeat("─", max(dividerMinRepeat, innerW-dividerWidthTrim)))
	tabs := thm.subtle.Width(innerW).Render(m.workspaceTabLine())

	var block string

	if observable {
		hint := thm.muted.Width(innerW).Render(
			"Observable · SELECT + enter · tab — focus · enter — run SQL / open row · esc — catalog+CLI · u — back (table) · alt+↑↓ — SQL history · f1 — CLI",
		)
		tbl := catalogTableChrome.Render(m.obsTable.View())
		help := thm.muted.Render("  " + m.obsTable.HelpView())
		input := thm.inputBar.Width(innerW).Render(m.obsInput.View())

		status := m.obsStatusBlock(&thm, innerW)

		footer1 := lipgloss.JoinHorizontal(
			lipgloss.Top,
			thm.footerKey.Render("f1"),
			thm.footerMuted.Render(" CLI · "),
			thm.footerKey.Render("click"),
			thm.footerMuted.Render(" tabs · "),
			thm.footerKey.Render("ctrl+c"),
			thm.footerMuted.Render(" quit · "),
			thm.footerKey.Render("esc"),
			thm.footerMuted.Render(" catalog"),
		)

		if status != "" {
			block = lipgloss.JoinVertical(
				lipgloss.Left,
				header,
				dividerLine,
				"",
				tabs,
				"",
				hint,
				"",
				tbl,
				"",
				help,
				"",
				input,
				"",
				status,
				"",
				footer1,
			)
		} else {
			block = lipgloss.JoinVertical(
				lipgloss.Left,
				header,
				dividerLine,
				"",
				tabs,
				"",
				hint,
				"",
				tbl,
				"",
				help,
				"",
				input,
				"",
				footer1,
			)
		}
	} else {
		body := m.vp.View()
		input := thm.inputBar.Width(innerW).Render(m.ti.View())

		footer1 := lipgloss.JoinHorizontal(
			lipgloss.Top,
			thm.footerKey.Render("ctrl+c"),
			thm.footerMuted.Render(" / esc quit · "),
			thm.footerKey.Render("f2"),
			thm.footerMuted.Render(" observable · "),
			thm.footerKey.Render("f1"),
			thm.footerMuted.Render(" / "),
			thm.footerKey.Render("f2"),
			thm.footerMuted.Render(" workspace tabs · "),
			thm.footerKey.Render(".tables"),
			thm.footerMuted.Render(" · "),
			thm.footerKey.Render("tab"),
			thm.footerMuted.Render(" complete · "),
			thm.footerKey.Render("↑↓"),
			thm.footerMuted.Render(" history"),
		)
		footer2 := lipgloss.JoinHorizontal(
			lipgloss.Top,
			thm.footerKey.Render("pgup/pgdn"),
			thm.footerMuted.Render(" log"),
		)
		footer := lipgloss.JoinVertical(lipgloss.Left, footer1, footer2)

		block = lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			dividerLine,
			"",
			tabs,
			"",
			body,
			"",
			input,
			"",
			footer,
		)
	}

	framed := thm.outer.Render(block)
	placed := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, framed)

	return placed, tabLineYInPlaced(placed)
}

func tabLineYInPlaced(placed string) int {
	lines := strings.Split(strings.TrimSuffix(placed, "\n"), "\n")

	for i, line := range lines {
		plain := ansi.Strip(line)
		if strings.Contains(plain, "CLI") && strings.Contains(plain, "Observable") {
			return i
		}
	}

	return -1
}

func pickWorkspaceTabFromLine(mouseX, mouseY, tabLineY, termWidth int, placed string) (shellTab, bool) {
	if tabLineY < 0 {
		return 0, false
	}

	lines := strings.Split(strings.TrimSuffix(placed, "\n"), "\n")
	if mouseY < 0 || mouseY >= len(lines) {
		return 0, false
	}

	if mouseY != tabLineY {
		return 0, false
	}

	line := lines[mouseY]

	plain := ansi.Strip(line)
	if !strings.Contains(plain, "CLI") || !strings.Contains(plain, "Observable") {
		return 0, false
	}

	lineWidth := lipgloss.Width(line)
	pad := max(0, (termWidth-lineWidth)/tabLineCenterDivisor)

	relX := mouseX - pad
	if relX < 0 || relX >= lineWidth {
		return 0, false
	}

	sep := strings.Index(plain, "│")
	if sep < 0 {
		sep = strings.Index(plain, "|")
	}

	if sep < 0 {
		if relX < lipgloss.Width(plain)/tabLineCenterDivisor {
			return shellTabCLI, true
		}

		return shellTabObservable, true
	}

	leftWidth := lipgloss.Width(plain[:sep])
	if relX < leftWidth {
		return shellTabCLI, true
	}

	return shellTabObservable, true
}

func (m *tuiModel) shellMouseHandler() func(tea.MouseMsg) tea.Cmd {
	return func(msg tea.MouseMsg) tea.Cmd {
		clickMsg, ok := msg.(tea.MouseClickMsg)
		if !ok {
			return nil
		}

		if clickMsg.Button != tea.MouseLeft {
			return nil
		}

		pointer := clickMsg.Mouse()
		obs := m.activeTab == shellTabObservable
		placed, tabY := m.shellPlaced(obs)

		tab, hit := pickWorkspaceTabFromLine(pointer.X, pointer.Y, tabY, m.width, placed)
		if !hit {
			return nil
		}

		if tab == m.activeTab {
			return nil
		}

		picked := tab

		return func() tea.Msg {
			return workspaceTabPickMsg(picked)
		}
	}
}
