package repl

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/c-bata/go-prompt"
	"github.com/pterm/pterm"
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

func (r *Repl) Run() { //nolint:gocyclo,gocognit // ignore
	// load history
	if err := r.init(); err != nil {
		pterm.FgRed.Println(err)
	}

	// Show help snippet
	r.help()

	for {
		input := prompt.Input("> ", completer,
			prompt.OptionTitle("shortdb"),
			prompt.OptionHistory(r.session.GetHistory()),
			prompt.OptionPrefixTextColor(prompt.Yellow),
			prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSuggestionBGColor(prompt.DarkGray),
		)

		if input == "" {
			continue
		}

		// if this next line
		if input[len(input)-1] == ';' || input[0] == '.' {
			input = fmt.Sprintf("%s %s", r.session.GetRaw(), input)
			r.session.Raw = ""
			r.session.Exec = true

			// set in history
			input = strings.TrimSpace(input)
			r.session.History = append(r.session.GetHistory(), input)
		} else {
			r.session.Raw += input + " "
			r.session.Exec = false
		}

		input = strings.TrimSpace(input)

		switch input[0] {
		case '.': // if this command
			s := strings.Split(input, " ")

			switch s[0] {
			case ".close":
				if err := r.close(); err != nil {
					pterm.FgRed.Println(err)
				}

				pterm.FgYellow.Println("Good buy!")

				return
			case ".open":
				if err := r.open(input); err != nil {
					pterm.FgRed.Println(err)
				}
			case ".help":
				r.help()
			case ".save":
				if err := r.save(); err != nil {
					pterm.FgRed.Println(err)
					continue
				}

				pterm.FgGreen.Println("Saved!")
			default:
				pterm.FgRed.Println("incorrect command")
			}
		default: // if this not command then this SQL-expression
			// if this multiline then skip
			if !r.session.GetExec() {
				continue
			}

			p, err := parser.New(input)
			if err != nil {
				pterm.FgRed.Println(err)
				continue
			}

			// exec query
			response, err := r.engine.Exec(p.GetQuery())
			if err != nil && err.Error() != "" {
				pterm.FgRed.Println(err)
				continue
			}

			if response != nil {
				pterm.FgGreen.Println(response)
			} else {
				pterm.FgGreen.Println(`Executed`)
			}
		}
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}

	return prompt.FilterHasPrefix(suggestions, w, true)
}
