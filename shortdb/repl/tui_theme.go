package repl

import (
	"os"

	"charm.land/lipgloss/v2"
)

// tuiTheme holds Lip Gloss styles derived once from terminal background.
type tuiTheme struct {
	accent      lipgloss.Style
	subtle      lipgloss.Style
	muted       lipgloss.Style
	divider     lipgloss.Style
	outer       lipgloss.Style
	inputBar    lipgloss.Style
	viewport    lipgloss.Style
	footerKey   lipgloss.Style
	footerMuted lipgloss.Style
	obsErr      lipgloss.Style
	innerWidth  int
}

func newTheme(totalWidth int) tuiTheme {
	if totalWidth <= 0 {
		totalWidth = defaultTermWidth
	}

	hasDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark := lipgloss.LightDark(hasDark)

	accentLight := lipgloss.Color("#5b21b6")
	accentDark := lipgloss.Color("#c4b5fd")
	borderLight := lipgloss.Color("#d4d4d8")
	borderDark := lipgloss.Color("#52525b")
	panelLight := lipgloss.Color("#f4f4f5")
	panelDark := lipgloss.Color("#18181b")
	outerLight := lipgloss.Color("#fafafa")
	outerDark := lipgloss.Color("#09090b")
	inputLight := lipgloss.Color("#e4e4e7")
	inputDark := lipgloss.Color("#27272a")
	subtleLight := lipgloss.Color("#52525b")
	subtleDark := lipgloss.Color("#a1a1aa")
	mutedLight := lipgloss.Color("#71717a")
	mutedDark := lipgloss.Color("#71717a")
	errLight := lipgloss.Color("#b91c1c")
	errDark := lipgloss.Color("#fca5a5")

	outer := lipgloss.NewStyle().
		Padding(outerPaddingV, outerPaddingH).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lightDark(accentLight, accentDark)).
		Background(lightDark(outerLight, outerDark))

	inner := max(minLayoutWidth, totalWidth-outer.GetHorizontalFrameSize())

	return tuiTheme{
		innerWidth: inner,
		accent: lipgloss.NewStyle().
			Bold(true).
			Foreground(lightDark(accentLight, accentDark)),
		subtle: lipgloss.NewStyle().
			Foreground(lightDark(subtleLight, subtleDark)),
		muted: lipgloss.NewStyle().
			Foreground(lightDark(mutedLight, mutedDark)),
		divider: lipgloss.NewStyle().
			Foreground(lightDark(borderLight, borderDark)),
		outer: outer,
		inputBar: lipgloss.NewStyle().
			Padding(0, 1).
			Background(lightDark(inputLight, inputDark)).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lightDark(borderLight, borderDark)),
		viewport: lipgloss.NewStyle().
			Padding(0, 1).
			Background(lightDark(panelLight, panelDark)).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lightDark(borderLight, borderDark)),
		footerKey: lipgloss.NewStyle().
			Foreground(lightDark(accentLight, accentDark)),
		footerMuted: lipgloss.NewStyle().
			Foreground(lightDark(mutedLight, mutedDark)),
		obsErr: lipgloss.NewStyle().
			Foreground(lightDark(errLight, errDark)),
	}
}
