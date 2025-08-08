package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/badtuxx/girus-cli/internal/common"
	"github.com/badtuxx/girus-cli/internal/repo"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	listRepoIndexURL string
)

var listCmd = &cobra.Command{
	Use:   "list [subcommand]",
	Short: common.T("Comandos para listar recursos", "Comandos para listar recursos"),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var listClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: common.T("Lista os clusters Kind disponíveis", "Lista los clusters Kind disponibles"),
	Long:  common.T("Lista todos os clusters Kind disponíveis no sistema, destacando os que executam o Girus.", "Lista todos los clusters Kind disponibles en el sistema, destacando los que ejecutan Girus."),
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(headerColor(common.T("CLUSTERS KIND", "CLUSTERS KIND")))
		fmt.Println(strings.Repeat("─", 80))
		fmt.Println(common.T("Obtendo lista de clusters Kind...", "Obteniendo lista de clusters Kind..."))

		getCmd := exec.Command("kind", "get", "clusters")
		output, err := getCmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s: %v\n", red("ERRO:"), common.T("Erro ao obter clusters Kind", "Error al obtener clusters Kind"), err)
			os.Exit(1)
		}

		clusters := strings.Split(strings.TrimSpace(string(output)), "\n")

		if len(clusters) == 0 || (len(clusters) == 1 && clusters[0] == "") {
			fmt.Println(common.T("Nenhum cluster Kind encontrado.", "Ningún cluster Kind encontrado."))
			return
		}

		fmt.Println("\n" + headerColor(common.T("Clusters Kind disponíveis:", "Clusters Kind disponibles:")))

		for _, cluster := range clusters {
			if cluster == "" {
				continue
			}

			// Verificar se é um cluster Girus verificando o namespace girus
			// Mudar contexto do kubectl para o cluster atual
			contextCmd := exec.Command("kubectl", "config", "use-context", fmt.Sprintf("kind-%s", cluster))
			contextCmd.Run() // Ignoramos erros aqui, pois vamos verificar no próximo comando

			// Verificar se o namespace girus existe
			checkCmd := exec.Command("kubectl", "get", "namespace", "girus", "--no-headers", "--ignore-not-found")
			checkOutput, _ := checkCmd.Output()

			isGirus := strings.Contains(string(checkOutput), "girus")

			if isGirus {
				fmt.Printf("%s Cluster %s (%s)\n", green(common.T("ATIVO", "ACTIVO")), magenta(cluster), "cluster Girus")

				// Verificar o status dos pods no namespace girus
				podsCmd := exec.Command("kubectl", "get", "pods", "-n", "girus", "-o", "custom-columns=NAME:.metadata.name,STATUS:.status.phase,READY:.status.containerStatuses[0].ready", "--no-headers")
				podsOutput, _ := podsCmd.Output()

				if len(podsOutput) > 0 {
					fmt.Println("   " + cyan(common.T("Pods:", "Pods:")))
					podLines := strings.Split(strings.TrimSpace(string(podsOutput)), "\n")
					for _, podLine := range podLines {
						if podLine != "" {
							fmt.Printf("   └─ %s\n", podLine)
						}
					}
				}
			} else {
				fmt.Printf("%s Cluster %s (%s)\n", red(common.T("INATIVO", "INACTIVO")), magenta(cluster), "cluster não-Girus")
			}
		}
	},
}

// Para compatibilidade, mantemos o comando singular, mas ele chamará o plural
var listClusterCmd = &cobra.Command{
	Use:    "cluster",
	Short:  common.T("Lista os clusters Kind disponíveis (alias para 'clusters')", "Lista los clusters Kind disponibles (alias de 'clusters')"),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		listClustersCmd.Run(cmd, args)
	},
}

// Estrutura para processar o JSON dos templates de laboratório
type LabTemplate struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    string `json:"duration"`
}

type LabListResponse struct {
	Templates []LabTemplate `json:"templates"`
}

var listLabsCmd = &cobra.Command{
	Use:   "labs",
	Short: common.T("Lista os laboratórios disponíveis no Girus", "Lista los laboratorios disponibles en Girus"),
	Long:  common.T("Lista todos os laboratórios disponíveis no cluster Girus ativo.", "Lista todos los laboratorios disponibles en el cluster Girus activo."),
	Run: func(cmd *cobra.Command, args []string) {
		// Criar formatadores de cores
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		headerColor := color.New(color.FgCyan, color.Bold).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()

		fmt.Println(headerColor(common.T("LABORATÓRIOS DISPONÍVEIS", "LABORATORIOS DISPONIBLES")))
		fmt.Println(strings.Repeat("─", 80))
		fmt.Println(common.T("Obtendo lista de laboratórios do Girus...", "Obteniendo lista de laboratorios de Girus..."))

		// Verificar se há um cluster Girus ativo
		checkCmd := exec.Command("kubectl", "get", "namespace", "girus", "--no-headers", "--ignore-not-found")
		checkOutput, err := checkCmd.Output()
		if err != nil || !strings.Contains(string(checkOutput), "girus") {
			fmt.Fprintf(os.Stderr, "%s %s\n", red("ERRO:"), common.T("Nenhum cluster Girus ativo encontrado", "Ningún cluster Girus activo encontrado"))
			fmt.Println(common.T("Use 'girus create cluster' para criar um cluster ou 'girus list clusters' para ver os clusters disponíveis.", "Use 'girus create cluster' para crear un cluster o 'girus list clusters' para ver los clusters disponibles."))
			os.Exit(1)
		}

		// Verificar o pod do backend
		backendCmd := exec.Command("kubectl", "get", "pods", "-n", "girus", "-l", "app=girus-backend", "-o", "jsonpath={.items[0].status.phase}")
		backendOutput, err := backendCmd.Output()
		if err != nil || string(backendOutput) != "Running" {
			fmt.Fprintf(os.Stderr, "%s %s\n", red("ERRO:"), common.T("O backend do Girus não está em execução", "El backend de Girus no está en ejecución"))
			fmt.Println(common.T("Verifique o status dos pods com 'kubectl get pods -n girus'", "Verifique el estado de los pods con 'kubectl get pods -n girus'"))
			os.Exit(1)
		}

		// Fazer uma solicitação para a API para obter a lista de laboratórios
		apiCmd := exec.Command("kubectl", "exec", "-n", "girus", "deploy/girus-backend", "--",
			"wget", "-q", "-O-", "http://localhost:8080/api/v1/templates")
		apiOutput, err := apiCmd.Output()

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s: %v\n", red("ERRO:"), common.T("Erro ao obter a lista de laboratórios", "Error al obtener la lista de laboratorios"), err)
			fmt.Println(common.T("Verifique se o serviço do backend está respondendo.", "Verifique si el servicio del backend está respondiendo."))
			os.Exit(1)
		}

		// Processar a resposta JSON
		var response LabListResponse
		if err := json.Unmarshal(apiOutput, &response); err != nil {
			fmt.Fprintf(os.Stderr, "%s %s: %v\n", red("ERRO:"), common.T("Erro ao processar a resposta", "Error al procesar la respuesta"), err)
			fmt.Println(common.T("Resposta da API:", "Respuesta de la API:"))
			fmt.Println(string(apiOutput))
			os.Exit(1)
		}

		// Exibir a lista de laboratórios
		if len(response.Templates) == 0 {
			fmt.Printf("\n%s %s\n", yellow("AVISO:"), common.T("Nenhum laboratório disponível.", "Ningún laboratorio disponible."))
			return
		}

		fmt.Println("\n" + headerColor(common.T("Laboratórios disponíveis:", "Laboratorios disponibles:")))

		for i, lab := range response.Templates {
			fmt.Printf("%d. %s", i+1, bold(lab.Title))
			if lab.Duration != "" {
				fmt.Printf(" (%s)", lab.Duration)
			}
			fmt.Println()
			fmt.Printf("   %s: %s\n", cyan("ID"), magenta(lab.Name))
			if lab.Description != "" {
				fmt.Printf("   %s\n", lab.Description)
			}
			fmt.Println()
		}

		fmt.Println("\n" + common.T("Para criar um laboratório, use:", "Para crear un laboratorio, use:"))
		fmt.Println("  " + magenta("girus create lab <lab-id>"))
	},
}

// Comando para listar laboratórios do repositório remoto
var listRepoLabsCmd = &cobra.Command{
	Use:   "repo-labs",
	Short: common.T("Lista os laboratórios disponíveis no repositório remoto", "Lista los laboratorios disponibles en el repositorio remoto"),
	Long:  common.T("Lista todos os laboratórios disponíveis no repositório remoto do GIRUS.", "Lista todos los laboratorios disponibles en el repositorio remoto de GIRUS."),
	Run: func(cmd *cobra.Command, args []string) {
		// Criar formatadores de cores
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		headerColor := color.New(color.FgCyan, color.Bold).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()

		fmt.Println(headerColor(common.T("LABORATÓRIOS DO REPOSITÓRIO", "LABORATORIOS DEL REPOSITORIO")))
		fmt.Println(strings.Repeat("─", 80))
		fmt.Println(common.T("Buscando laboratórios no repositório remoto...", "Buscando laboratorios en el repositorio remoto..."))

		// Obter o index.yaml
		index, err := repo.GetLabsIndex(listRepoIndexURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %v\n", red("ERRO:"), err)
			os.Exit(1)
		}

		if len(index.Labs) == 0 {
			fmt.Printf("\n%s %s\n", yellow("AVISO:"), common.T("Nenhum laboratório disponível no repositório.", "Ningún laboratorio disponible en el repositorio."))
			return
		}

		fmt.Println("\n" + headerColor(common.T("Laboratórios disponíveis no GIRUS Hub:", "Laboratorios disponibles en GIRUS Hub:")))
		fmt.Println(strings.Repeat("─", 60))

		for i, lab := range index.Labs {
			if i > 0 {
				// Separador entre os laboratórios
				fmt.Println(strings.Repeat("─", 60))
			}

			fmt.Printf("%s: %s\n", cyan("ID"), magenta(lab.ID))
			fmt.Printf("%s: %s\n", cyan("Título"), bold(lab.Title))

			if lab.Description != "" {
				fmt.Printf("%s: %s\n", cyan("Descrição"), lab.Description)
			}

			if lab.Duration != "" {
				fmt.Printf("%s: %s\n", cyan("Duração"), lab.Duration)
			}

			if lab.Version != "" {
				fmt.Printf("%s: %s\n", cyan("Versão"), lab.Version)
			}

			fmt.Printf("%s: %s\n", cyan("Tags"), repo.FormatTags(lab.Tags))
		}

		fmt.Println(strings.Repeat("─", 60))
		fmt.Println("\n" + common.T("Para instalar um laboratório, use:", "Para instalar un laboratorio, use:"))
		fmt.Println("  " + magenta("girus create lab <lab-id>"))
	},
}

func init() {
	listCmd.AddCommand(listClustersCmd)
	listCmd.AddCommand(listClusterCmd) // Para compatibilidade
	listCmd.AddCommand(listLabsCmd)

	// Adicionar o novo comando repo-labs
	listCmd.AddCommand(listRepoLabsCmd)

	// Flags para o comando repo-labs
	listRepoLabsCmd.Flags().StringVarP(&listRepoIndexURL, "url", "u", "", "URL do arquivo index.yaml (opcional)")
}
