package repl

import (
	"os"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func newTUIModel(repl *Repl) *tuiModel {
	textIn := textinput.New()
	textIn.Prompt = "> "
	textIn.Placeholder = "SQL or .help · tab complete"
	textIn.ShowSuggestions = true
	textIn.CharLimit = inputCharLimit
	textIn.SetSuggestions(SuggestionTexts())

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
	tiStyles.Focused.Suggestion = lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("#71717a"), lipgloss.Color("#a1a1aa")))
	tiStyles.Blurred = tiStyles.Focused
	textIn.SetStyles(tiStyles)

	scroll := viewport.New()
	scroll.SoftWrap = true
	scroll.MouseWheelEnabled = false

	shellModel := &tuiModel{
		repl:             repl,
		ti:               textIn,
		vp:               scroll,
		obsTable:         newCatalogTableModel(defaultTermWidth, defaultTermHeight, nil),
		obsInput:         newObservableSQLInput(),
		obsPaginator:     newObsPaginator(),
		histIndex:        -1,
		obsHistBrowseIdx: -1,
		width:            defaultTermWidth,
		height:           defaultTermHeight,
		theme:            newTheme(defaultTermWidth),
		activeTab:        shellTabCLI,
	}
	shellModel.applyChromeStyles()
	shellModel.layout()

	return shellModel
}

func (m *tuiModel) Init() tea.Cmd {
	return m.ti.Focus()
}
