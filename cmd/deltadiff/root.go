package main

import (
	"github.com/spf13/cobra"
)

type RootCommand struct {
	program *Program
}

func (rc *RootCommand) Run(cmd *cobra.Command, args []string) {
	cmd.Help()
	rc.program.Exit(0)
}

func (p *Program) createRootCmd() *cobra.Command {
	rc := &RootCommand{
		program:     p,
	}

	cmd := &cobra.Command{
		Use:   "deltadiff",
		Short: "Calculates delta diff of two files",
		Long:  `Calculates delta diff of two files, base and target, and apply the required changes in order to transform base into target.`,
		Run:   rc.Run,
	}

	return cmd
}
