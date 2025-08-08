package main

import (
	"fmt"
	"os"

	"github.com/badtuxx/girus-cli/cmd"
	"github.com/badtuxx/girus-cli/internal/common"
)

func main() {
	cfg := common.LoadConfig()
	common.SetLanguage(cfg.Language)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao executar o comando: %s\n", err)
		os.Exit(1)
	}
}
