package repl

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"charm.land/lipgloss/v2"

	session "github.com/shortlink-org/shortdb/shortdb/domain/session/v1"
	"github.com/shortlink-org/shortdb/shortdb/engine"
	"github.com/shortlink-org/shortdb/shortdb/engine/file"
	parser "github.com/shortlink-org/shortdb/shortdb/parser/v1"
)

type Repl struct {
	mu sync.Mutex

	engine  engine.Engine
	session *session.Session
}

func New(ctx context.Context, sess *session.Session) (*Repl, error) {
	// set engine
	store, err := engine.New(ctx, "file", file.SetName(sess.GetCurrentDatabase()), file.SetPath("/tmp/shortdb_repl"))
	if err != nil {
		return nil, fmt.Errorf("failed to create engine: %w", err)
	}

	return &Repl{
		session: sess,
		engine:  store,
	}, nil
}

func (r *Repl) Run() {
	err := r.RunTUI()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, errStyle.Render(err.Error()))
	}
}

var (
	errStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	okStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("78"))
	warnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
)

// handleREPLLine processes one submitted line (same semantics as the former go-prompt loop).
func (r *Repl) handleREPLLine(line string) ([]string, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, false
	}

	input := line

	if input[len(input)-1] == ';' || input[0] == '.' {
		input = fmt.Sprintf("%s %s", r.session.GetRaw(), input)
		r.session.Raw = ""
		r.session.Exec = true

		input = strings.TrimSpace(input)
		r.session.History = append(r.session.GetHistory(), input)
	} else {
		r.session.Raw += line + " "
		r.session.Exec = false
	}

	input = strings.TrimSpace(input)

	switch input[0] {
	case '.':
		s := strings.Split(input, " ")

		switch s[0] {
		case ".close":
			err := r.close()
			if err != nil {
				return []string{errStyle.Render(err.Error())}, false
			}

			return []string{warnStyle.Render("Goodbye!")}, true
		case ".open":
			err := r.open(input)
			if err != nil {
				return []string{errStyle.Render(err.Error())}, false
			}

			return []string{okStyle.Render("database switched")}, false
		case ".help":
			return []string{strings.TrimSpace(r.helpString())}, false
		case ".save":
			err := r.save()
			if err != nil {
				return []string{errStyle.Render(err.Error())}, false
			}

			return []string{okStyle.Render("Saved!")}, false
		default:
			return []string{errStyle.Render("incorrect command")}, false
		}
	default:
		if !r.session.GetExec() {
			return nil, false
		}

		p, err := parser.New(input)
		if err != nil {
			return []string{errStyle.Render(err.Error())}, false
		}

		response, err := r.engine.Exec(p.GetQuery())
		if err != nil && err.Error() != "" {
			return []string{errStyle.Render(err.Error())}, false
		}

		if response != nil {
			lines := formatExecTranscript(response)

			out := make([]string, len(lines))
			for i := range lines {
				out[i] = okStyle.Render(lines[i])
			}

			return out, false
		}

		return []string{okStyle.Render("Executed")}, false
	}
}
