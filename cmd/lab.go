package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/badtuxx/girus-cli/internal/common"
	"github.com/badtuxx/girus-cli/internal/repo"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var labCmd = &cobra.Command{
	Use:   "lab",
	Short: common.T("Gerencia laboratórios", "Gestiona laboratorios"),
	Long:  common.T(`Gerencia laboratórios, permitindo listar, instalar e remover laboratórios dos repositórios configurados.`, `Gestiona laboratorios, permitiendo listar, instalar y eliminar laboratorios de los repositorios configurados.`),
}

var labListCmd = &cobra.Command{
	Use:   "list",
	Short: common.T("Lista todos os laboratórios disponíveis", "Lista todos los laboratorios disponibles"),
	Long:  common.T(`Lista todos os laboratórios disponíveis em todos os repositórios configurados.`, `Lista todos los laboratorios disponibles en todos los repositorios configurados.`),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Criar formatadores de cores
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()
		headerColor := color.New(color.FgCyan, color.Bold).SprintFunc()

		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return fmt.Errorf("%s %v", red(common.T("ERRO:", "ERROR:")), err)
		}

		lm, err := repo.NewLabManager(rm)
		if err != nil {
			return fmt.Errorf("%s %v", red(common.T("ERRO:", "ERROR:")), err)
		}

		labs, err := lm.ListLabs()
		if err != nil {
			return fmt.Errorf("%s %v", red(common.T("ERRO:", "ERROR:")), err)
		}

		fmt.Println(headerColor(common.T("LABORATÓRIOS DISPONÍVEIS", "LABORATORIOS DISPONIBLES")))
		fmt.Println(strings.Repeat("─", 80))

		if len(labs) == 0 {
			fmt.Println(common.T("Nenhum laboratório disponível.", "Ningún laboratorio disponible."))
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
	Short: common.T("Instala um laboratório", "Instala un laboratorio"),
	Long:  common.T(`Instala um laboratório específico de um repositório.`, `Instala un laboratorio específico de un repositorio.`),
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

		fmt.Println(headerColor(common.T("INSTALANDO LABORATÓRIO", "INSTALANDO LABORATORIO")))
		fmt.Println(strings.Repeat("─", 80))
		fmt.Printf(common.T("Instalando laboratório %s do repositório %s...\n", "Instalando el laboratorio %s del repositorio %s...\n"), magenta(labName), magenta(repoName))

		if err := lm.DownloadLab(repoName, labName, version); err != nil {
			return fmt.Errorf("%s %v", red("ERRO:"), err)
		}

		fmt.Printf("%s %s %s %s\n", green(common.T("SUCESSO:", "ÉXITO:")), common.T("Laboratório", "Laboratorio"), magenta(labName), common.T("instalado com sucesso.", "instalado con éxito."))

		// Reinicia o backend
		fmt.Println("\n" + headerColor(common.T("REINICIANDO BACKEND", "REINICIANDO BACKEND")))
		fmt.Println(strings.Repeat("─", 80))
		fmt.Println(common.T("Reiniciando o backend para aplicar as mudanças...", "Reiniciando el backend para aplicar los cambios..."))

		restartCmd := exec.Command("kubectl", "rollout", "restart", "deployment/girus-backend", "-n", "girus")
		if err := restartCmd.Run(); err != nil {
			return fmt.Errorf("%s %s: %v", red("ERRO:"), common.T("Erro ao reiniciar o backend", "Error al reiniciar el backend"), err)
		}

		// Aguarda o reinício completar
		fmt.Println(common.T("Aguardando o reinício do backend completar...", "Esperando a que el backend reinicie por completo..."))
		waitCmd := exec.Command("kubectl", "rollout", "status", "deployment/girus-backend", "-n", "girus", "--timeout=60s")
		if err := waitCmd.Run(); err != nil {
			return fmt.Errorf("%s %s: %v", red("ERRO:"), common.T("Erro ao aguardar reinício do backend", "Error al esperar el reinicio del backend"), err)
		}
		fmt.Printf("%s Backend %s\n", green(common.T("SUCESSO:", "ÉXITO:")), common.T("reiniciado com sucesso.", "reiniciado con éxito."))

		return nil
	},
}

var labSearchCmd = &cobra.Command{
	Use:   "search [termo]",
	Short: common.T("Busca laboratórios por termo", "Busca laboratorios por término"),
	Long:  common.T(`Busca laboratórios por termo, procurando em nomes, descrições e tags`, `Busca laboratorios por término, buscando en nombres, descripciones y etiquetas`),
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
			return fmt.Errorf("%s %s: %v", red(common.T("ERRO:", "ERROR:")), common.T("Erro ao criar gerenciador de repositórios", "Error al crear el gestor de repositorios"), err)
		}

		lm, err := repo.NewLabManager(rm)
		if err != nil {
			return fmt.Errorf("%s %s: %v", red(common.T("ERRO:", "ERROR:")), common.T("Erro ao criar gerenciador de laboratórios", "Error al crear el gestor de laboratorios"), err)
		}

		labs, err := lm.ListLabs()
		if err != nil {
			return fmt.Errorf("%s %s: %v", red(common.T("ERRO:", "ERROR:")), common.T("Erro ao listar laboratórios", "Error al listar laboratorios"), err)
		}

		fmt.Println(headerColor(common.T("BUSCA DE LABORATÓRIOS", "BÚSQUEDA DE LABORATORIOS")))
		fmt.Println(strings.Repeat("─", 80))
		fmt.Printf(common.T("Buscando por: %s\n\n", "Buscando por: %s\n\n"), magenta(term))

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
			fmt.Printf("\n%s %s '%s'\n",
				red(common.T("AVISO:", "AVISO:")), common.T("Nenhum laboratório encontrado para o termo", "Ningún laboratorio encontrado para el término"), magenta(term))
		}

		return nil
	},
}

func init() {
	labCmd.AddCommand(labListCmd, labInstallCmd, labSearchCmd)

	// Flags para os comandos
	labInstallCmd.Flags().String("version", "", common.T("Versão específica do laboratório", "Versión específica del laboratorio"))
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
