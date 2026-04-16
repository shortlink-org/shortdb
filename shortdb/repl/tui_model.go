package repl

import (
	"charm.land/bubbles/v2/paginator"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"

	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
)

type shellTab byte

const (
	shellTabCLI shellTab = iota
	shellTabObservable
)

type obsBrowseMode byte

const (
	obsBrowseCatalog obsBrowseMode = iota
	obsBrowseQuery
)

// tuiModel is the Bubble Tea state for the ShortDB shell.
type tuiModel struct {
	repl *Repl

	vp          viewport.Model
	ti          textinput.Model
	obsTable    table.Model
	obsInput    textinput.Model
	transcript  []string
	width       int
	height      int
	initialized bool
	theme       tuiTheme

	// histIndex is an index into session command history while browsing with ↑/↓ (-1 = off).
	histIndex int

	activeTab     shellTab
	obsBrowseMode obsBrowseMode
	obsFocusTable bool
	obsStatus     string
	obsStatusErr  bool

	obsAllRows    []*page.Row
	obsPaginator  paginator.Model
	obsStack      []obsNavFrame
	obsSQLHistory []string
	// obsHistBrowseIdx is an index into obsSQLHistory while browsing with alt+↑/↓ (-1 = off).
	obsHistBrowseIdx int
}
