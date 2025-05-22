package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/badtuxx/girus-cli/internal/repo"
	"github.com/fatih/color"
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
		// Criar formatadores de cores
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()
		headerColor := color.New(color.FgCyan, color.Bold).SprintFunc()

		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return fmt.Errorf("%s %v", red("ERRO:"), err)
		}

		lm, err := repo.NewLabManager(rm)
		if err != nil {
			return fmt.Errorf("%s %v", red("ERRO:"), err)
		}

		labs, err := lm.ListLabs()
		if err != nil {
			return fmt.Errorf("%s %v", red("ERRO:"), err)
		}

		fmt.Println(headerColor("LABORATÓRIOS DISPONÍVEIS"))
		fmt.Println(strings.Repeat("─", 80))

		if len(labs) == 0 {
			fmt.Println("Nenhum laboratório disponível.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, cyan("NOME")+"\t"+cyan("VERSÃO")+"\t"+cyan("REPOSITÓRIO")+"\t"+cyan("DESCRIÇÃO"))
		for repoName, entries := range labs {
			for _, entry := range entries {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					magenta(entry.ID),
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
		// Criar formatadores de cores
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()
		headerColor := color.New(color.FgCyan, color.Bold).SprintFunc()

		repoName := args[0]
		labName := args[1]
		version, _ := cmd.Flags().GetString("version")

		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return fmt.Errorf("%s %v", red("ERRO:"), err)
		}

		lm, err := repo.NewLabManager(rm)
		if err != nil {
			return fmt.Errorf("%s %v", red("ERRO:"), err)
		}

		fmt.Println(headerColor("INSTALANDO LABORATÓRIO"))
		fmt.Println(strings.Repeat("─", 80))
		fmt.Printf("Instalando laboratório %s do repositório %s...\n", magenta(labName), magenta(repoName))

		if err := lm.DownloadLab(repoName, labName, version); err != nil {
			return fmt.Errorf("%s %v", red("ERRO:"), err)
		}

		fmt.Printf("%s Laboratório %s instalado com sucesso.\n", green("SUCESSO:"), magenta(labName))

		// Reinicia o backend
		fmt.Println("\n" + headerColor("REINICIANDO BACKEND"))
		fmt.Println(strings.Repeat("─", 80))
		fmt.Println("Reiniciando o backend para aplicar as mudanças...")

		restartCmd := exec.Command("kubectl", "rollout", "restart", "deployment/girus-backend", "-n", "girus")
		if err := restartCmd.Run(); err != nil {
			return fmt.Errorf("%s Erro ao reiniciar o backend: %v", red("ERRO:"), err)
		}

		// Aguarda o reinício completar
		fmt.Println("Aguardando o reinício do backend completar...")
		waitCmd := exec.Command("kubectl", "rollout", "status", "deployment/girus-backend", "-n", "girus", "--timeout=60s")
		if err := waitCmd.Run(); err != nil {
			return fmt.Errorf("%s Erro ao aguardar reinício do backend: %v", red("ERRO:"), err)
		}
		fmt.Printf("%s Backend reiniciado com sucesso.\n", green("SUCESSO:"))

		return nil
	},
}

var labSearchCmd = &cobra.Command{
	Use:   "search [termo]",
	Short: "Busca laboratórios por termo",
	Long:  `Busca laboratórios por termo, procurando em nomes, descrições e tags`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Criar formatadores de cores
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()
		headerColor := color.New(color.FgCyan, color.Bold).SprintFunc()

		term := strings.ToLower(args[0])

		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return fmt.Errorf("%s Erro ao criar gerenciador de repositórios: %v", red("ERRO:"), err)
		}

		lm, err := repo.NewLabManager(rm)
		if err != nil {
			return fmt.Errorf("%s Erro ao criar gerenciador de laboratórios: %v", red("ERRO:"), err)
		}

		labs, err := lm.ListLabs()
		if err != nil {
			return fmt.Errorf("%s Erro ao listar laboratórios: %v", red("ERRO:"), err)
		}

		fmt.Println(headerColor("BUSCA DE LABORATÓRIOS"))
		fmt.Println(strings.Repeat("─", 80))
		fmt.Printf("Buscando por: %s\n\n", magenta(term))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, cyan("NOME")+"\t"+cyan("VERSÃO")+"\t"+cyan("REPOSITÓRIO")+"\t"+cyan("DESCRIÇÃO"))

		found := false
		for repoName, entries := range labs {
			for _, entry := range entries {
				// Verifica se o termo está no título, descrição ou tags
				if containsCaseInsensitive(entry.Title, term) ||
					containsCaseInsensitive(entry.Description, term) ||
					containsCaseInsensitive(entry.Tags, term) {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
						magenta(entry.ID),
						entry.Version,
						repoName,
						entry.Description)
					found = true
				}
			}
		}

		w.Flush()

		if !found {
			fmt.Printf("\n%s Nenhum laboratório encontrado para o termo '%s'\n",
				red("AVISO:"), magenta(term))
		}

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
