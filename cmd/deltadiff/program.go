package main

import (
	"fmt"
	"os"
)

type Program struct {
}

func NewProgram() *Program {
	return &Program{}
}

func (p *Program) Run() {
	rootCmd := p.createRootCmd()
	signatureCmd := p.createSignatureCmd()
	deltaCmd := p.createDeltaCmd()
	patchCmd := p.createPatchCmd()

	rootCmd.AddCommand(signatureCmd)
	rootCmd.AddCommand(deltaCmd)
	rootCmd.AddCommand(patchCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		p.Exit(1)
	}

	p.Exit(0)
}

func (p *Program) Exit(code int) {
	os.Exit(code)
}
