package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/badtuxx/girus-cli/internal/helpers"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var forceDelete bool
var verboseDelete bool

var deleteCmd = &cobra.Command{
	Use:   "delete [subcommand]",
	Short: "Comandos para excluir recursos",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var deleteClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Exclui o cluster Girus",
	Long:  "Exclui o cluster Girus do sistema, incluindo todos os recursos do Girus.",
	Run: func(cmd *cobra.Command, args []string) {
		// Criar formatadores de cores
		red := color.New(color.FgRed).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()
		headerColor := color.New(color.FgCyan, color.Bold).SprintFunc()

		clusterName := "girus"

		// Verificar se o cluster existe
		checkCmd := exec.Command("kind", "get", "clusters")
		output, err := checkCmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s Erro ao obter lista de clusters: %v\n", red("ERRO:"), err)
			os.Exit(1)
		}

		clusters := strings.Split(strings.TrimSpace(string(output)), "\n")
		clusterExists := false
		for _, cluster := range clusters {
			if cluster == clusterName {
				clusterExists = true
				break
			}
		}

		if !clusterExists {
			fmt.Fprintf(os.Stderr, "%s Cluster %s não encontrado\n", red("ERRO:"), magenta("girus"))
			os.Exit(1)
		}

		// Confirmar a exclusão se -f/--force não estiver definido
		if !forceDelete {
			fmt.Printf("%s Você está prestes a excluir o cluster %s. Esta ação é irreversível.\n",
				yellow("AVISO:"), magenta(clusterName))
			fmt.Print("Deseja continuar? [s/N]: ")

			reader := bufio.NewReader(os.Stdin)
			confirmStr, _ := reader.ReadString('\n')
			confirm := strings.TrimSpace(strings.ToLower(confirmStr))

			if confirm != "s" && confirm != "sim" && confirm != "y" && confirm != "yes" {
				fmt.Println("Operação cancelada pelo usuário.")
				return
			}
		}

		fmt.Println(headerColor("Excluindo o cluster Girus..."))

		if verboseDelete {
			// Excluir o cluster mostrando o output normal
			deleteCmd := exec.Command("kind", "delete", "cluster", "--name", clusterName)
			deleteCmd.Stdout = os.Stdout
			deleteCmd.Stderr = os.Stderr

			if err := deleteCmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "%s Erro ao excluir o cluster Girus: %v\n", red("ERRO:"), err)
				os.Exit(1)
			}
		} else {
			// Usando barra de progresso (padrão)
			barConfig := helpers.ProgressBarConfig{
				Total:            100,
				Description:      "Excluindo cluster...",
				Width:            80,
				Throttle:         65,
				SpinnerType:      14,
				RenderBlankState: true,
				ShowBytes:        false,
				SetPredictTime:   false,
			}
			bar := helpers.CreateProgressBar(barConfig)

			// Executar comando sem mostrar saída
			deleteCmd := exec.Command("kind", "delete", "cluster", "--name", clusterName)
			var stderr bytes.Buffer
			deleteCmd.Stderr = &stderr

			// Iniciar o comando
			err := deleteCmd.Start()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s Erro ao iniciar o comando: %v\n", red("ERRO:"), err)
				os.Exit(1)
			}

			// Atualizar a barra de progresso enquanto o comando está em execução
			done := make(chan struct{})
			go func() {
				for {
					select {
					case <-done:
						return
					default:
						bar.Add(1)
						time.Sleep(150 * time.Millisecond)
					}
				}
			}()

			// Aguardar o final do comando
			err = deleteCmd.Wait()
			close(done)
			bar.Finish()

			if err != nil {
				fmt.Fprintf(os.Stderr, "%s Erro ao excluir o cluster Girus: %v\n%s\n", red("ERRO:"), err, stderr.String())
				os.Exit(1)
			}
		}

		fmt.Println("\n" + green("SUCESSO:") + " Cluster " + magenta("Girus") + " excluído com sucesso!")
	},
}

func init() {
	deleteCmd.AddCommand(deleteClusterCmd)

	// Flag para forçar a exclusão sem confirmação
	deleteClusterCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Força a exclusão sem confirmação")

	// Flag para modo detalhado com output completo
	deleteClusterCmd.Flags().BoolVarP(&verboseDelete, "verbose", "v", false, "Modo detalhado com output completo em vez da barra de progresso")
}
