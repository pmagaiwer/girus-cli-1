package main

import (
	"fmt"
	"os"

	"github.com/badtuxx/girus-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao executar o comando: %s\n", err)
		os.Exit(1)
	}
} 