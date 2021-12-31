package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xrash/deltadiff"
	"io"
	"os"
)

type SignatureCommand struct {
	program *Program

	options struct {
		hasher    string
		blockSize uint32
	}
}

func (sc *SignatureCommand) Run(cmd *cobra.Command, args []string) {

	if len(args) > 2 {
		fmt.Println("command signature requires at most two args")
		sc.program.Exit(1)
	}

	baseReader, baseSize, err := sc.decideBaseReader(args)
	if err != nil {
		fmt.Println(err)
		sc.program.Exit(1)
	}

	signatureWriter, err := sc.decideSignatureWriter(args)
	if err != nil {
		fmt.Println(err)
		sc.program.Exit(1)
	}

	config := &deltadiff.SignatureConfig{
		Hasher:    sc.options.hasher,
		BlockSize: int(sc.options.blockSize),
		BaseSize:  baseSize,
	}

	if err := deltadiff.Signature(baseReader, signatureWriter, config); err != nil {
		fmt.Println("Error", err)
		sc.program.Exit(1)
	}

	sc.program.Exit(0)
}

func (p *Program) createSignatureCmd() *cobra.Command {

	sc := &SignatureCommand{
		program: p,
	}

	cmd := &cobra.Command{
		Use:   "signature <base> <signature>",
		Short: "Create the signature of base",
		Long:  `Create the signature of base`,
		Run:   sc.Run,
	}

	cmd.Flags().StringVarP(
		&sc.options.hasher,
		"hasher",
		"",
		"polyroll",
		"Hasher to be used, can be md5, crc32 or polyroll",
	)

	cmd.Flags().Uint32VarP(
		&sc.options.blockSize,
		"block-size",
		"",
		1024,
		"Size of the blocks used in the rolling hash algorithm",
	)

	return cmd
}

func (sc *SignatureCommand) decideBaseReader(args []string) (io.Reader, int, error) {
	if len(args) == 0 {
		return os.Stdin, -1, nil
	}

	if args[0] == "-" {
		return os.Stdin, -1, nil
	}

	filename := args[0]

	file, err := os.Open(filename)
	if err != nil {
		return nil, -1, fmt.Errorf("Error opening base file %s: %v", filename, err)
	}

	info, err := file.Stat()
	if err != nil {
		return nil, -1, fmt.Errorf("Error running Stat() on file %s: %v", filename, err)
	}

	return file, int(info.Size()), nil
}

func (sc *SignatureCommand) decideSignatureWriter(args []string) (io.Writer, error) {
	if len(args) == 0 || len(args) == 1 {
		return os.Stdout, nil
	}

	filename := args[1]

	if filename == "-" {
		return os.Stdout, nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening signature file %s: %v", filename, err)
	}

	return file, nil
}
