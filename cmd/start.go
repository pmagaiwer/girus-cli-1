package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/badtuxx/girus-cli/internal/common"
	"github.com/badtuxx/girus-cli/internal/k8s"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: common.T("Inicia o ambiente do GIRUS", "Inicia el entorno de GIRUS"),
	Long:  common.T("Inicia o ambiente do GIRUS CLI, reiniciando o deployment do backend e do frontend.", "Inicia el entorno del GIRUS CLI, reiniciando los deployments del backend y del frontend."),
	Run: func(cmd *cobra.Command, args []string) {
		// Define os nomes dos deployments
		frontendDeploymentName := "girus-frontend"
		backendDeploymentName := "girus-backend"
		// Criando um client para interagir com o cluster do Kubernetes
		client, err := k8s.NewKubernetesClient()
		if err != nil {
			fmt.Printf("%s %s: %v\n", color.New().SprintFunc()(common.T("ERRO", "ERROR")), common.T("Erro ao criar cliente Kubernetes", "Error al crear cliente de Kubernetes"), err)
			return
		}

		ctx := context.Background()

		pods, err := client.ListRunningPods(ctx, "girus")
		if err != nil {
			fmt.Printf("%s %s: %v\n", color.New().SprintFunc()(common.T("ERRO", "ERROR")), common.T("Erro ao tentar pegar a lista de pods em execução no namespace do GIRUS", "Error al obtener la lista de pods en ejecución en el namespace de GIRUS"), err)
			fmt.Println(common.T("Leia o erro, se você não conseguir resolvê-lo, recrie o cluster.", "Lea el error; si no puede resolverlo, recree el cluster."))
			return
		}
		// Pega todos os pods do namespace do girus
		var frontendPod string
		var backendPod string
		for _, pod := range pods {
			if strings.Contains(pod, "frontend") {
				frontendPod = pod
			} else if strings.Contains(pod, "backend") {
				backendPod = pod
			}
		}
		// Checa se o frontend já está executando antes de tentar iniciar o deployment
		isFrontendRunning, err := client.IsPodRunning(ctx, "girus", frontendPod)
		if isFrontendRunning {
			fmt.Println(common.T("O pod de frontend já está em execução.", "El pod de frontend ya está en ejecución."))
			fmt.Println(common.T("Tente abrir o browser e navegar até http://localhost:8000.", "Intente abrir el navegador y acceder a http://localhost:8000."))
		}
		if err != nil {
			fmt.Println(common.T("Nenhum pod do frontend encontrado no namespace do GIRUS...", "Ningún pod de frontend encontrado en el namespace de GIRUS..."))
		}

		// Checa se o backend já está executando antes de tentar iniciar o deployment
		isBackendRunning, err := client.IsPodRunning(ctx, "girus", backendPod)
		if isBackendRunning {
			fmt.Println(common.T("O pod de backend já está em execução.", "El pod de backend ya está en ejecución."))
			fmt.Println(common.T("Tente abrir o browser e navegar até http://localhost:8000.", "Intente abrir el navegador y acceder a http://localhost:8000."))
			fmt.Printf("%s %s\n", yellow(common.T("AVISO", "AVISO")), common.T("Cancelando.", "Cancelando."))
			return
		}
		if err != nil {
			fmt.Println(common.T("Nenhum pod do backend encontrado no namespace do GIRUS...", "Ningún pod de backend encontrado en el namespace de GIRUS..."))
		}
		err = startDeployment(client, ctx, backendDeploymentName)
		if err != nil {
			fmt.Printf("%s %s: %v\n", color.New().SprintFunc()(common.T("ERRO", "ERROR")), common.T("Erro ao tentar iniciar o backend", "Error al intentar iniciar el backend"), err)
			return
		}
		err = startDeployment(client, ctx, frontendDeploymentName)
		if err != nil {
			fmt.Printf("%s %s: %v\n", color.New().SprintFunc()(common.T("ERRO", "ERROR")), common.T("Erro ao tentar iniciar o frontend", "Error al intentar iniciar el frontend"), err)
			return
		}
	},
}

func startDeployment(client *k8s.KubernetesClient, ctx context.Context, deploymentName string) error {
	magenta := color.New(color.FgMagenta).SprintFunc()
	err := client.CreateDeployment(ctx, "girus", deploymentName)
	if err != nil {
		fmt.Printf("%s %s %s: %v\n", color.New().SprintFunc()(common.T("ERRO", "ERROR")), common.T("Erro ao tentar iniciar o deploy", "Error al intentar iniciar el deploy"), magenta(deploymentName), err)
		fmt.Println(common.T("Leia o erro, se você não conseguir resolvê-lo, recrie o cluster.", "Lea el error; si no puede resolverlo, recree el cluster."))
		return err
	}

	return nil
}
