package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/badtuxx/girus-cli/internal/common"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	cliRepo   = "badtuxx/girus-cli"
	githubAPI = "https://api.github.com/repos"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Atualiza o GIRUS CLI para a última versão",
	Long: `Verifica e atualiza o GIRUS CLI para a última versão disponível.
Após a atualização, oferece a opção de recriar o cluster para garantir
compatibilidade com as novas features.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Criar formatadores de cores
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc() // Para informações importantes
		yellow := color.New(color.FgYellow).SprintFunc()

		// Criar formatador para títulos
		headerColor := color.New(color.FgCyan, color.Bold).SprintFunc()

		// Exibir cabeçalho
		fmt.Println(strings.Repeat("─", 80))
		fmt.Println(headerColor("GIRUS UPDATE"))
		fmt.Println(strings.Repeat("─", 80))

		// Verificar versão atual da CLI
		currentVersion := common.Version
		fmt.Printf("%s: %s\n", bold("Versão atual da CLI"), magenta(currentVersion))

		// Obter última versão do GitHub
		fmt.Println("\n" + headerColor("Verificando atualizações..."))
		latestCliVersion, err := GetLatestGitHubVersion(cliRepo)
		if err != nil {
			return fmt.Errorf("%s erro ao verificar última versão da CLI: %v", red("ERRO:"), err)
		}

		fmt.Printf("%s: %s\n", bold("Última versão disponível"), magenta(latestCliVersion))

		// Verificar se já está na versão mais recente
		isLatest := !IsNewerVersion(latestCliVersion, currentVersion)

		if isLatest {
			fmt.Println("\n" + green("Você já está usando a última versão do GIRUS CLI!"))
			return nil
		}

		// Confirmar atualização
		fmt.Printf("\n%s (%s). Deseja atualizar? (S/n): ",
			yellow("Nova versão disponível"), magenta(latestCliVersion))
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "s" && response != "" {
			fmt.Println(yellow("Atualização cancelada."))
			return nil
		}

		// Atualizar CLI
		fmt.Println("\n" + headerColor("Atualizando CLI..."))
		if err := downloadAndInstall(latestCliVersion); err != nil {
			return fmt.Errorf("%s erro ao atualizar CLI: %v", red("ERRO:"), err)
		}
		fmt.Printf("%s CLI atualizada com sucesso para a versão %s!\n",
			green("SUCESSO:"), magenta(latestCliVersion))

		// Perguntar se deseja recriar o cluster
		fmt.Print("\n" + yellow("Deseja recriar o cluster para garantir compatibilidade com as novas features? (S/n): "))
		fmt.Scanln(&response)
		if strings.ToLower(response) == "s" || response == "" {
			fmt.Println("\n" + headerColor("Recriando o cluster..."))

			// Executar o comando delete
			deleteCmd := exec.Command("girus", "delete")
			deleteCmd.Stdout = os.Stdout
			deleteCmd.Stderr = os.Stderr
			if err := deleteCmd.Run(); err != nil {
				return fmt.Errorf("%s erro ao deletar o cluster: %v", red("ERRO:"), err)
			}

			// Executar o comando create
			createCmd := exec.Command("girus", "create")
			createCmd.Stdout = os.Stdout
			createCmd.Stderr = os.Stderr
			if err := createCmd.Run(); err != nil {
				return fmt.Errorf("%s erro ao criar o cluster: %v", red("ERRO:"), err)
			}

			fmt.Println("\n" + green("Cluster recriado com sucesso!"))
		} else {
			fmt.Println("\n" + yellow("Cluster mantido como está. Lembre-se que algumas novas features podem não funcionar corretamente com o cluster atual."))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

// GetLatestGitHubVersion obtém a última versão do GitHub
func GetLatestGitHubVersion(repo string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/releases/latest", githubAPI, repo))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}

	return strings.TrimPrefix(release.TagName, "v"), nil
}

// IsNewerVersion compara duas versões semânticas
// retorna TRUE se v1 é MAIS NOVA que v2
func IsNewerVersion(v1, v2 string) bool {
	// Padronizar para ter 3 componentes (major.minor.patch)
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Garantir que temos pelo menos 3 partes
	for len(parts1) < 3 {
		parts1 = append(parts1, "0")
	}
	for len(parts2) < 3 {
		parts2 = append(parts2, "0")
	}

	// Comparar major
	if parts1[0] > parts2[0] {
		return true
	}
	if parts1[0] < parts2[0] {
		return false
	}

	// Comparar minor
	if parts1[1] > parts2[1] {
		return true
	}
	if parts1[1] < parts2[1] {
		return false
	}

	// Comparar patch
	if parts1[2] > parts2[2] {
		return true
	}

	// São iguais ou v1 é menor
	return false
}

// downloadAndInstall baixa e instala a nova versão da CLI
func downloadAndInstall(version string) error {
	// Determinar sistema operacional e arquitetura
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Construir URL do binário
	binaryName := fmt.Sprintf("girus-cli-%s-%s", osName, arch)
	if osName == "windows" {
		binaryName += ".exe"
	}
	url := fmt.Sprintf("https://github.com/%s/releases/download/v%s/%s", cliRepo, version, binaryName)

	// Criar diretório temporário
	tempDir, err := os.MkdirTemp("", "girus-update")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// Baixar binário
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao baixar binário: %s", resp.Status)
	}

	// Salvar binário
	binaryPath := filepath.Join(tempDir, "girus")
	if osName == "windows" {
		binaryPath += ".exe"
	}
	out, err := os.Create(binaryPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	// Tornar binário executável
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return err
	}

	// Obter caminho do binário atual
	currentBinary, err := os.Executable()
	if err != nil {
		return err
	}

	// Se estiver instalado globalmente, usar sudo
	if strings.HasPrefix(currentBinary, "/usr/local/bin/") || strings.HasPrefix(currentBinary, "/usr/bin/") {
		cmd := exec.Command("sudo", "mv", binaryPath, currentBinary)
		if err := cmd.Run(); err != nil {
			return err
		}
	} else {
		// Caso contrário, apenas mover para o mesmo diretório
		if err := os.Rename(binaryPath, currentBinary); err != nil {
			return err
		}
	}

	return nil
}
