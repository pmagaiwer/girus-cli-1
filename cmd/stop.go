package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/badtuxx/girus-cli/internal/common"
	"github.com/badtuxx/girus-cli/internal/k8s"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Criar formatadores de cores
	red         = color.New(color.FgRed).SprintFunc()
	cyan        = color.New(color.FgCyan).SprintFunc()
	green       = color.New(color.FgGreen).SprintFunc()
	yellow      = color.New(color.FgYellow).SprintFunc()
	magenta     = color.New(color.FgMagenta).SprintFunc()
	headerColor = color.New(color.FgCyan, color.Bold).SprintFunc()
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: common.T("Parar o ambiente do GIRUS", "Detener el entorno de GIRUS"),
	Long:  common.T("Parar o ambiente do GIRUS CLI, removendo todos os recursos criados pelo GIRUS CLI.", "Detener el entorno del GIRUS CLI eliminando todos los recursos creados por el GIRUS CLI."),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(common.T("%s Você está prestes a parar o %s e o %s no cluster %s.\n",
			"%s Está a punto de detener %s y %s en el cluster %s.\n"),
			yellow(common.T("AVISO:", "AVISO:")), magenta("frontend"), magenta("backend"), magenta(clusterName))
		fmt.Print(common.T("Deseja continuar? [s/N]: ", "¿Desea continuar? [s/N]: "))

		reader := bufio.NewReader(os.Stdin)
		confirmStr, _ := reader.ReadString('\n')
		confirm := strings.TrimSpace(strings.ToLower(confirmStr))

		if confirm != "s" && confirm != "sim" && confirm != "y" && confirm != "yes" {
			fmt.Println(common.T("Operação cancelada pelo usuário.", "Operación cancelada por el usuario."))
			return
		}
		// Define os nomes dos deployments
		frontendDeploymentName := "girus-frontend"
		backendDeploymentName := "girus-backend"
		// Criando um client para interagir com o cluster do Kubernetes
		client, err := k8s.NewKubernetesClient()
		if err != nil {
			fmt.Printf("%s %s: %v\n", red(common.T("ERRO:", "ERROR:")), common.T("Erro ao criar cliente Kubernetes", "Error al crear cliente de Kubernetes"), err)
			return
		}

		ctx := context.Background()
		// Pega todos os pods do namespace do girus
		pods, err := client.ListRunningPods(ctx, "girus")
		if err != nil {
			fmt.Printf("%s %s: %v\n", red(common.T("ERRO:", "ERROR:")), common.T("Erro ao tentar pegar a lista de pods", "Error al obtener la lista de pods"), err)
			return
		}

		// Pega o nome dos pods do frontend e do backend
		var frontendPod string
		var backendPod string
		for _, pod := range pods {
			if strings.Contains(pod, "frontend") {
				frontendPod = pod
			} else if strings.Contains(pod, "backend") {
				backendPod = pod
			}
		}

		// Verifica se o backend está em execução, se sim, parar o deploy e remover o serviço
		if isRunning, _ := client.IsPodRunning(ctx, "girus", backendPod); isRunning {
			err := deleteDeployment(client, ctx, backendDeploymentName)
			if err != nil {
				fmt.Printf("falha ao tentar parar o deploy do backend do GIRUS: %v\n", err)
				return
			}
			fmt.Println("✅ Backend parado com sucesso.")

		} else {
			fmt.Println("⚠️ " + common.T("O backend não está em execução.", "El backend no está en ejecución."))
		}

		// Verifica se o frontend está em execução, se sim, parar o deploy e remover o serviço
		if isRunning, _ := client.IsPodRunning(ctx, "girus", frontendPod); isRunning {
			err := deleteDeployment(client, ctx, frontendDeploymentName)
			if err != nil {
				fmt.Printf("%s %s: %v\n", red(common.T("ERRO:", "ERROR:")), common.T("falha ao tentar parar o deploy do frontend do GIRUS", "fallo al intentar detener el deploy del frontend de GIRUS"), err)
				return
			}
			fmt.Println("✅ " + common.T("Frontend parado com sucesso.", "Frontend detenido con éxito."))
		} else {
			fmt.Println("⚠️ " + common.T("O frontend não está em execução..", "El frontend no está en ejecución."))
		}
	},
}

func deleteDeployment(client *k8s.KubernetesClient, ctx context.Context, deploymentName string) error {
	err := client.StopDeployAndWait(ctx, "girus", deploymentName)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "%s %s: %v\n", red(common.T("ERRO:", "ERROR:")), common.T("Erro ao tentar parar o deploy", "Error al intentar detener el deploy"), err)
		if err != nil {
			return err
		}

		fmt.Printf(common.T("%s Você quer forçar a parada do deployment %s?\n",
			"%s ¿Desea forzar la detención del deployment %s?\n"), yellow(common.T("AVISO:", "AVISO:")), magenta(deploymentName))
		fmt.Print(common.T("Deseja continuar? [s/N]: ", "¿Desea continuar? [s/N]: "))

		reader := bufio.NewReader(os.Stdin)
		confirmStr, _ := reader.ReadString('\n')
		confirm := strings.TrimSpace(strings.ToLower(confirmStr))

		if confirm != "s" && confirm != "sim" && confirm != "y" && confirm != "yes" {
			fmt.Println(common.T("Operação cancelada pelo usuário.", "Operación cancelada por el usuario."))
			return err
		}
		return err
	}

	return nil
}
