package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/badtuxx/girus-cli/internal/repo"
	"github.com/spf13/cobra"
)

var labCmd = &cobra.Command{
	Use:   "lab",
	Short: "Gerencia laboratórios",
	Long:  `Gerencia laboratórios, permitindo listar, instalar e remover laboratórios dos repositórios configurados.`,
}

var labListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos os laboratórios disponíveis",
	Long:  `Lista todos os laboratórios disponíveis em todos os repositórios configurados.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return err
		}

		lm, err := repo.NewLabManager(rm)
		if err != nil {
			return err
		}

		labs, err := lm.ListLabs()
		if err != nil {
			return err
		}

		if len(labs) == 0 {
			fmt.Println("Nenhum laboratório disponível.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NOME\tVERSÃO\tREPOSITÓRIO\tDESCRIÇÃO")
		for name, entries := range labs {
			for _, entry := range entries {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, entry.Version, entry.URL, entry.Description)
			}
		}
		w.Flush()

		return nil
	},
}

var labInstallCmd = &cobra.Command{
	Use:   "install [repositório] [laboratório]",
	Short: "Instala um laboratório",
	Long:  `Instala um laboratório específico de um repositório.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoName := args[0]
		labName := args[1]
		version, _ := cmd.Flags().GetString("version")

		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return err
		}

		lm, err := repo.NewLabManager(rm)
		if err != nil {
			return err
		}

		if err := lm.DownloadLab(repoName, labName, version); err != nil {
			return err
		}

		fmt.Printf("Laboratório '%s' instalado com sucesso.\n", labName)
		return nil
	},
}

var labSearchCmd = &cobra.Command{
	Use:   "search [termo]",
	Short: "Busca laboratórios",
	Long:  `Busca laboratórios por nome ou palavras-chave.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		term := strings.ToLower(args[0])

		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return err
		}

		lm, err := repo.NewLabManager(rm)
		if err != nil {
			return err
		}

		labs, err := lm.ListLabs()
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NOME\tVERSÃO\tREPOSITÓRIO\tDESCRIÇÃO")
		found := false
		for name, entries := range labs {
			for _, entry := range entries {
				// Busca por nome, descrição ou palavras-chave
				if strings.Contains(strings.ToLower(name), term) ||
					strings.Contains(strings.ToLower(entry.Description), term) ||
					containsCaseInsensitive(entry.Keywords, term) {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, entry.Version, entry.URL, entry.Description)
					found = true
				}
			}
		}
		w.Flush()

		if !found {
			fmt.Printf("\nNenhum laboratório encontrado para o termo '%s'\n", args[0])
		}

		return nil
	},
}

func init() {
	labCmd.AddCommand(labListCmd, labInstallCmd, labSearchCmd)

	// Flags para os comandos
	labInstallCmd.Flags().String("version", "", "Versão específica do laboratório")
}

// containsCaseInsensitive verifica se uma string está presente em um slice, ignorando maiúsculas/minúsculas
func containsCaseInsensitive(slice []string, str string) bool {
	str = strings.ToLower(str)
	for _, s := range slice {
		if strings.Contains(strings.ToLower(s), str) {
			return true
		}
	}
	return false
} 