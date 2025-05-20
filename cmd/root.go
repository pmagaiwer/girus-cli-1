package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "girus",
	Short: "GIRUS - Plataforma de Laboratórios Interativos",
	Long: `GIRUS é uma plataforma open-source de laboratórios interativos que permite a criação,
gerenciamento e execução de ambientes de aprendizado prático para tecnologias como Linux,
Docker, Kubernetes, Terraform e outras ferramentas essenciais para profissionais de DevOps,
SRE, Dev e Platform Engineering.`,
}

// Execute executa o comando raiz
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Adiciona os comandos
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(labCmd)
	rootCmd.AddCommand(repoCmd)

	// Configura flags globais
	rootCmd.PersistentFlags().StringP("config", "c", "", "arquivo de configuração (padrão: $HOME/.girus/config.yaml)")
} 