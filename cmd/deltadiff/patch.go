package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xrash/deltadiff"
	"io"
	"os"
)

type PatchCommand struct {
	program *Program
}

func (pc *PatchCommand) Run(cmd *cobra.Command, args []string) {

	if len(args) < 2 {
		fmt.Println("command patch requires at least 2 args")
		pc.program.Exit(1)
	}

	if len(args) > 3 {
		fmt.Println("command patch requires at most 3 args")
		pc.program.Exit(1)
	}

	baseReader, err := pc.decideBaseReader(args)
	if err != nil {
		fmt.Println(err)
		pc.program.Exit(1)
	}

	deltaReader, err := pc.decideDeltaReader(args)
	if err != nil {
		fmt.Println(err)
		pc.program.Exit(1)
	}

	resultWriter, err := pc.decideResultWriter(args)
	if err != nil {
		fmt.Println(err)
		pc.program.Exit(1)
	}

	if err := deltadiff.Patch(baseReader, deltaReader, resultWriter); err != nil {
		fmt.Println("Error", err)
		pc.program.Exit(1)
	}

	pc.program.Exit(0)
}

func (p *Program) createPatchCmd() *cobra.Command {

	pc := &PatchCommand{
		program: p,
	}

	cmd := &cobra.Command{
		Use:   "patch <base> <delta> <result>",
		Short: "Apply delta to base",
		Long:  `Apply delta to base`,
		Run:   pc.Run,
	}

	return cmd
}

func (pc *PatchCommand) decideBaseReader(args []string) (io.Reader, error) {
	filename := args[0]
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening base file %s: %v", filename, err)
	}

	return file, nil
}

func (pc *PatchCommand) decideDeltaReader(args []string) (io.Reader, error) {
	filename := args[1]
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening base file %s: %v", filename, err)
	}

	return file, nil
}

func (pc *PatchCommand) decideResultWriter(args []string) (io.Writer, error) {
	if len(args) == 0 || len(args) == 1 || len(args) == 2 {
		return os.Stdout, nil
	}

	filename := args[2]

	if filename == "-" {
		return os.Stdout, nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening result file %s: %v", filename, err)
	}

	return file, nil
}
