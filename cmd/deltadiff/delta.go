package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xrash/deltadiff"
	"io"
	"os"
)

type DeltaCommand struct {
	program *Program

	options struct {
		debug     bool
		debugFile string
	}
}

func (dc *DeltaCommand) Run(cmd *cobra.Command, args []string) {

	if len(args) < 2 {
		fmt.Println("command delta requires at least 2 args")
		dc.program.Exit(1)
	}

	if len(args) > 3 {
		fmt.Println("command delta requires at most 3 args")
		dc.program.Exit(1)
	}

	signatureReader, err := dc.decideSignatureReader(args)
	if err != nil {
		fmt.Println(err)
		dc.program.Exit(1)
	}

	targetReader, err := dc.decideTargetReader(args)
	if err != nil {
		fmt.Println(err)
		dc.program.Exit(1)
	}

	deltaWriter, err := dc.decideDeltaWriter(args)
	if err != nil {
		fmt.Println(err)
		dc.program.Exit(1)
	}

	debugFile, err := dc.decideDebugFile(dc.options.debug, dc.options.debugFile)
	if err != nil {
		fmt.Println(err)
		dc.program.Exit(1)
	}

	c := &deltadiff.DeltaConfig{
		Debug:       dc.options.debug,
		DebugWriter: debugFile,
	}

	if err := deltadiff.Delta(signatureReader, targetReader, deltaWriter, c); err != nil {
		fmt.Println("Error", err)
		dc.program.Exit(1)
	}

	dc.program.Exit(0)
}

func (p *Program) createDeltaCmd() *cobra.Command {

	dc := &DeltaCommand{
		program: p,
	}

	cmd := &cobra.Command{
		Use:   "delta <signature> <target> <delta>",
		Short: "Produce delta of signature and target",
		Long:  `Produce delta of signature and target`,
		Run:   dc.Run,
	}

	cmd.Flags().BoolVarP(
		&dc.options.debug,
		"debug",
		"",
		false,
		"If enabled, displays debug information, also check --debug-file",
	)

	cmd.Flags().StringVarP(
		&dc.options.debugFile,
		"debug-file",
		"",
		"/dev/stderr",
		"File to write debug information to",
	)

	return cmd
}

func (dc *DeltaCommand) decideDebugFile(debug bool, debugFile string) (io.Writer, error) {
	if !debug {
		return nil, nil
	}

	if debugFile == "-" {
		return os.Stderr, nil
	}

	file, err := os.Open(debugFile)
	if err != nil {
		return nil, fmt.Errorf("Error opening debug file %s: %v", debugFile, err)
	}

	return file, nil
}

func (dc *DeltaCommand) decideSignatureReader(args []string) (io.Reader, error) {
	filename := args[0]
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening base file %s: %v", filename, err)
	}

	return file, nil
}

func (dc *DeltaCommand) decideTargetReader(args []string) (io.Reader, error) {
	filename := args[1]
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening base file %s: %v", filename, err)
	}

	return file, nil
}

func (dc *DeltaCommand) decideDeltaWriter(args []string) (io.Writer, error) {
	if len(args) == 0 || len(args) == 1 || len(args) == 2 {
		return os.Stdout, nil
	}

	filename := args[2]

	if filename == "-" {
		return os.Stdout, nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening delta file %s: %v", filename, err)
	}

	return file, nil
}
