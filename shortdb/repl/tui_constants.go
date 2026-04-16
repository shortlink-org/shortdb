package repl

const (
	minLayoutWidth       = 10
	layoutChromeRows     = 2 // input row + viewport chrome
	minViewportHeight    = 4
	minDataColumnWidth   = 4
	maxDrillTableNameLen = 256
	inputCharLimit       = 8192
	defaultTermWidth     = 80
	defaultTermHeight    = 24
	extraRowsForChrome   = 6 // outer frame, divider, gaps, input bar, footer
	footerRowCount       = 2
	outerPaddingV        = 1
	outerPaddingH        = 2
	dividerMinRepeat     = 4
	dividerWidthTrim     = 2
	tabLineCenterDivisor = 2 // (termWidth - lineWidth) / 2 for centered tab row hit testing
)
