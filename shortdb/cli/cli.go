package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	session "github.com/shortlink-org/shortdb/shortdb/domain/session/v1"
	"github.com/shortlink-org/shortdb/shortdb/repl"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	var gendocsOut string

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

	gendocsCmd := &cobra.Command{
		Use:   "gendocs",
		Short: "Generate Cobra command reference as Markdown files",
		RunE: func(_ *cobra.Command, _ []string) error {
			return doc.GenMarkdownTree(rootCmd, gendocsOut)
		},
	}
	gendocsCmd.Flags().StringVarP(&gendocsOut, "output", "o", "shortdb/docs", "directory to write Markdown files into")

	rootCmd.AddCommand(gendocsCmd)

	err := rootCmd.Execute()

	cancel()

	if err != nil {
		//nolint:revive,forbidigo // CLI error path
		fmt.Println(err)
		os.Exit(1)
	}
}
