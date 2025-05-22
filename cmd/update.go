package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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
		// Verificar versão atual da CLI
		currentVersion := Version
		fmt.Printf("Versão atual da CLI: %s\n", currentVersion)

		// Obter última versão do GitHub
		latestCliVersion, err := GetLatestGitHubVersion(cliRepo)
		if err != nil {
			return fmt.Errorf("erro ao verificar última versão da CLI: %v", err)
		}

		fmt.Printf("Última versão disponível: %s\n", latestCliVersion)

		// Verificar se já está na versão mais recente
		isLatest := !IsNewerVersion(latestCliVersion, currentVersion)

		if isLatest {
			fmt.Println("Você já está usando a última versão do GIRUS CLI!")
			return nil
		}

		// Confirmar atualização
		fmt.Printf("\nNova versão disponível (%s). Deseja atualizar? (S/n): ", latestCliVersion)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "s" && response != "" {
			fmt.Println("Atualização cancelada.")
			return nil
		}

		// Atualizar CLI
		fmt.Println("\nAtualizando CLI...")
		if err := downloadAndInstall(latestCliVersion); err != nil {
			return fmt.Errorf("erro ao atualizar CLI: %v", err)
		}
		fmt.Printf("CLI atualizada com sucesso para a versão %s!\n", latestCliVersion)

		// Perguntar se deseja recriar o cluster
		fmt.Print("\nDeseja recriar o cluster para garantir compatibilidade com as novas features? (S/n): ")
		fmt.Scanln(&response)
		if strings.ToLower(response) == "s" || response == "" {
			fmt.Println("\nRecriando o cluster...")

			// Executar o comando delete
			deleteCmd := exec.Command("girus", "delete")
			deleteCmd.Stdout = os.Stdout
			deleteCmd.Stderr = os.Stderr
			if err := deleteCmd.Run(); err != nil {
				return fmt.Errorf("erro ao deletar o cluster: %v", err)
			}

			// Executar o comando create
			createCmd := exec.Command("girus", "create")
			createCmd.Stdout = os.Stdout
			createCmd.Stderr = os.Stderr
			if err := createCmd.Run(); err != nil {
				return fmt.Errorf("erro ao criar o cluster: %v", err)
			}

			fmt.Println("\nCluster recriado com sucesso!")
		} else {
			fmt.Println("\nCluster mantido como está. Lembre-se que algumas novas features podem não funcionar corretamente com o cluster atual.")
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
	major1, _ := strconv.Atoi(parts1[0])
	major2, _ := strconv.Atoi(parts2[0])
	if major1 > major2 {
		return true
	}
	if major1 < major2 {
		return false
	}

	// Comparar minor
	minor1, _ := strconv.Atoi(parts1[1])
	minor2, _ := strconv.Atoi(parts2[1])
	if minor1 > minor2 {
		return true
	}
	if minor1 < minor2 {
		return false
	}

	// Comparar patch
	patch1, _ := strconv.Atoi(parts1[2])
	patch2, _ := strconv.Atoi(parts2[2])
	if patch1 > patch2 {
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
