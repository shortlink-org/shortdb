package main

import (
	"context"
	"fmt"

	session "github.com/shortlink-org/shortdb/shortdb/domain/session/v1"
	"github.com/shortlink-org/shortdb/shortdb/repl"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rootCmd := &cobra.Command{
		Use:   "shortdb",
		Short: "ShortDB it's daabase for experiments",
		Long:  "Implementation simple database like SQLite",
		Run: func(_ *cobra.Command, _ []string) {
			// run new session
			s, err := session.New()
			if err != nil {
				panic(err)
			}

			// run REPL by default
			r, err := repl.New(ctx, s)
			if err != nil {
				panic(err)
			}

			r.Run()
		},
	}

	if err := rootCmd.Execute(); err != nil {
		//nolint:revive,forbidigo // just print error
		fmt.Println(err)

		return
	}

	// Generate docs
	if err := doc.GenMarkdownTree(rootCmd, "./pkg/shortdb/docs"); err != nil {
		//nolint:revive,forbidigo // just print error
		fmt.Println(err)

		return
	}
}
