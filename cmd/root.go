package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "girus",
	Short: "CLI para administrar a plataforma Girus",
	Long: `Girus CLI é uma ferramenta de linha de comando para administrar
a plataforma Girus de laboratórios interativos baseada em Kubernetes.

Esta ferramenta permite criar, listar e excluir clusters Kubernetes
para execução da plataforma Girus.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute executa o comando raiz
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Adicione os subcomandos aqui
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
} 