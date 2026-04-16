package repl

import "strings"

func (m *tuiModel) applyChromeStyles() {
	m.vp.Style = m.theme.viewport
}

func (m *tuiModel) layout() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	const tabBarRows = 1

	if m.activeTab == shellTabObservable {
		tblWidth, tblHeight := m.observableTableSize()
		m.syncObsPaginatorForRows(tblHeight, len(m.obsAllRows))
		m.rebuildObsTableSlice(tblWidth, tblHeight)
		inner := m.theme.innerWidth
		obsInner := inner - m.theme.inputBar.GetHorizontalFrameSize()
		m.obsInput.SetWidth(max(minLayoutWidth, obsInner))

		return
	}

	const (
		titleLines  = 1
		dividerLine = 1
		gap         = 1
	)

	vpH := m.height - titleLines - tabBarRows - dividerLine - footerRowCount - gap - layoutChromeRows - extraRowsForChrome
	vpH = max(minViewportHeight, vpH)

	inner := m.theme.innerWidth
	vpInner := inner - m.vp.Style.GetHorizontalFrameSize()
	tiInner := inner - m.theme.inputBar.GetHorizontalFrameSize()

	m.vp.SetWidth(max(minLayoutWidth, vpInner))
	m.vp.SetHeight(vpH)
	m.ti.SetWidth(max(minLayoutWidth, tiInner))
}

func (m *tuiModel) syncViewport() {
	m.vp.SetContent(strings.Join(m.transcript, "\n"))
	m.vp.GotoBottom()
}

func (m *tuiModel) appendLines(lines ...string) {
	m.transcript = append(m.transcript, lines...)
	m.syncViewport()
}
