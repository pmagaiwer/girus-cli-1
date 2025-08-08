package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/badtuxx/girus-cli/internal/common"
	"github.com/badtuxx/girus-cli/internal/helpers"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var forceDelete bool
var verboseDelete bool

var deleteCmd = &cobra.Command{
	Use:   "delete [subcommand]",
	Short: common.T("Comandos para excluir recursos", "Comandos para eliminar recursos"),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var deleteClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: common.T("Exclui o cluster Girus", "Elimina el cluster Girus"),
	Long:  common.T("Exclui o cluster Girus do sistema, incluindo todos os recursos do Girus.", "Elimina el cluster Girus del sistema, incluyendo todos los recursos de Girus."),
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
			fmt.Fprintf(os.Stderr, "%s %s: %v\n", red("ERRO:"), common.T("Erro ao obter lista de clusters", "Error al obtener la lista de clusters"), err)
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
			fmt.Fprintf(os.Stderr, "%s %s %s %s\n", red("ERRO:"), common.T("Cluster", "Cluster"), magenta("girus"), common.T("não encontrado", "no encontrado"))
			os.Exit(1)
		}

		// Confirmar a exclusão se -f/--force não estiver definido
		if !forceDelete {
			fmt.Printf(common.T("%s Você está prestes a excluir o cluster %s. Esta ação é irreversível.\n",
				"%s Está a punto de eliminar el cluster %s. Esta acción es irreversible.\n"),
				yellow("AVISO:"), magenta(clusterName))
			fmt.Print(common.T("Deseja continuar? [s/N]: ", "¿Desea continuar? [s/N]: "))

			reader := bufio.NewReader(os.Stdin)
			confirmStr, _ := reader.ReadString('\n')
			confirm := strings.TrimSpace(strings.ToLower(confirmStr))

			if confirm != "s" && confirm != "sim" && confirm != "y" && confirm != "yes" {
				fmt.Println(common.T("Operação cancelada pelo usuário.", "Operación cancelada por el usuario."))
				return
			}
		}

		fmt.Println(headerColor(common.T("Excluindo o cluster Girus...", "Eliminando el cluster Girus...")))

		if verboseDelete {
			// Excluir o cluster mostrando o output normal
			deleteCmd := exec.Command("kind", "delete", "cluster", "--name", clusterName)
			deleteCmd.Stdout = os.Stdout
			deleteCmd.Stderr = os.Stderr

			if err := deleteCmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "%s %s: %v\n", red("ERRO:"), common.T("Erro ao excluir o cluster Girus", "Error al eliminar el cluster Girus"), err)
				os.Exit(1)
			}
		} else {
			// Usando barra de progresso (padrão)
			barConfig := helpers.ProgressBarConfig{
				Total:            100,
				Description:      common.T("Excluindo cluster...", "Eliminando cluster..."),
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
				fmt.Fprintf(os.Stderr, "%s %s: %v\n", red("ERRO:"), common.T("Erro ao iniciar o comando", "Error al iniciar el comando"), err)
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
				fmt.Fprintf(os.Stderr, "%s %s: %v\n%s\n", red("ERRO:"), common.T("Erro ao excluir o cluster Girus", "Error al eliminar el cluster Girus"), err, stderr.String())
				os.Exit(1)
			}
		}

		fmt.Println("\n" + green(common.T("SUCESSO:", "ÉXITO:")) + " " + common.T("Cluster", "Cluster") + " " + magenta("Girus") + " " + common.T("excluído com sucesso!", "eliminado con éxito!"))
	},
}

func init() {
	deleteCmd.AddCommand(deleteClusterCmd)

	// Flag para forçar a exclusão sem confirmação
	deleteClusterCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, common.T("Força a exclusão sem confirmação", "Fuerza la eliminación sin confirmación"))

	// Flag para modo detalhado com output completo
	deleteClusterCmd.Flags().BoolVarP(&verboseDelete, "verbose", "v", false, common.T("Modo detalhado com output completo em vez da barra de progresso", "Modo detallado con salida completa en lugar de la barra de progreso"))
}
