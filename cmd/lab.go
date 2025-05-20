package cmd

import (
	"fmt"
	"os"
	"os/exec"
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
		for repoName, entries := range labs {
			for _, entry := range entries {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					entry.ID,
					entry.Version,
					repoName,
					entry.Description)
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

		// Reinicia o backend
		fmt.Println("Reiniciando o backend...")
		restartCmd := exec.Command("kubectl", "rollout", "restart", "deployment/girus-backend", "-n", "girus")
		if err := restartCmd.Run(); err != nil {
			return fmt.Errorf("erro ao reiniciar o backend: %v", err)
		}

		// Aguarda o reinício completar
		fmt.Println("Aguardando o reinício do backend completar...")
		waitCmd := exec.Command("kubectl", "rollout", "status", "deployment/girus-backend", "-n", "girus", "--timeout=60s")
		if err := waitCmd.Run(); err != nil {
			return fmt.Errorf("erro ao aguardar reinício do backend: %v", err)
		}
		fmt.Println("Backend reiniciado com sucesso.")

		return nil
	},
}

var labSearchCmd = &cobra.Command{
	Use:   "search [termo]",
	Short: "Busca laboratórios por termo",
	Long:  `Busca laboratórios por termo, procurando em nomes, descrições e tags`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		term := strings.ToLower(args[0])

		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return fmt.Errorf("erro ao criar gerenciador de repositórios: %v", err)
		}

		lm, err := repo.NewLabManager(rm)
		if err != nil {
			return fmt.Errorf("erro ao criar gerenciador de laboratórios: %v", err)
		}

		labs, err := lm.ListLabs()
		if err != nil {
			return fmt.Errorf("erro ao listar laboratórios: %v", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NOME\tVERSÃO\tREPOSITÓRIO\tDESCRIÇÃO")

		found := false
		for repoName, entries := range labs {
			for _, entry := range entries {
				// Verifica se o termo está no título, descrição ou tags
				if containsCaseInsensitive(entry.Title, term) ||
					containsCaseInsensitive(entry.Description, term) ||
					containsCaseInsensitive(entry.Tags, term) {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
						entry.ID,
						entry.Version,
						repoName,
						entry.Description)
					found = true
				}
			}
		}

		if !found {
			fmt.Printf("\nNenhum laboratório encontrado para o termo '%s'\n", term)
		}

		w.Flush()
		return nil
	},
}

func init() {
	labCmd.AddCommand(labListCmd, labInstallCmd, labSearchCmd)

	// Flags para os comandos
	labInstallCmd.Flags().String("version", "", "Versão específica do laboratório")
}

// containsCaseInsensitive verifica se uma string está contida em outra, ignorando maiúsculas/minúsculas
func containsCaseInsensitive(s interface{}, term string) bool {
	switch v := s.(type) {
	case string:
		return strings.Contains(strings.ToLower(v), term)
	case []string:
		for _, str := range v {
			if strings.Contains(strings.ToLower(str), term) {
				return true
			}
		}
	}
	return false
}
