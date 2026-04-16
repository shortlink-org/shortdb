package repl

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// RunTUI starts the Bubble Tea REPL interface.
func (r *Repl) RunTUI() error {
	err := r.init()
	if err != nil {
		return fmt.Errorf("repl init: %w", err)
	}

	m := newTUIModel(r)
	m.appendLines(strings.TrimSpace(r.helpString()))

	// When stdin is not a TTY (e.g. some IDE runners), bubbletea otherwise opens
	// /dev/tty for input and can fail with ENXIO. Explicit stdin/stdout skips that path.
	p := tea.NewProgram(m, tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))

	_, err = p.Run()
	if err != nil {
		return fmt.Errorf("repl: %w", err)
	}

	return nil
}
