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
		clusterName := "girus"

		// Verificar se o cluster existe
		checkCmd := exec.Command("kind", "get", "clusters")
		output, err := checkCmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao obter lista de clusters: %v\n", err)
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
			fmt.Fprintf(os.Stderr, "Erro: cluster 'girus' não encontrado\n")
			os.Exit(1)
		}

		// Confirmar a exclusão se -f/--force não estiver definido
		if !forceDelete {
			fmt.Printf("Você está prestes a excluir o cluster Girus. Esta ação é irreversível.\n")
			fmt.Print("Deseja continuar? (s/N): ")

			reader := bufio.NewReader(os.Stdin)
			confirmStr, _ := reader.ReadString('\n')
			confirm := strings.TrimSpace(strings.ToLower(confirmStr))

			if confirm != "s" && confirm != "sim" && confirm != "y" && confirm != "yes" {
				fmt.Println("Operação cancelada pelo usuário.")
				return
			}
		}

		fmt.Println("Excluindo o cluster Girus...")

		if verboseDelete {
			// Excluir o cluster mostrando o output normal
			deleteCmd := exec.Command("kind", "delete", "cluster", "--name", clusterName)
			deleteCmd.Stdout = os.Stdout
			deleteCmd.Stderr = os.Stderr

			if err := deleteCmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao excluir o cluster Girus: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Usando barra de progresso (padrão)
			barConfig := helpers.ProgressBarConfig{
				Total:            100,
				Description:      "Excluindo cluster...",
				Width:            50,
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
				fmt.Fprintf(os.Stderr, "Erro ao iniciar o comando: %v\n", err)
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
				fmt.Fprintf(os.Stderr, "Erro ao excluir o cluster Girus: %v\n%s\n", err, stderr.String())
				os.Exit(1)
			}
		}

		fmt.Println("Cluster Girus excluído com sucesso!")
	},
}

func init() {
	deleteCmd.AddCommand(deleteClusterCmd)

	// Flag para forçar a exclusão sem confirmação
	deleteClusterCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Força a exclusão sem confirmação")

	// Flag para modo detalhado com output completo
	deleteClusterCmd.Flags().BoolVarP(&verboseDelete, "verbose", "v", false, "Modo detalhado com output completo em vez da barra de progresso")
}
