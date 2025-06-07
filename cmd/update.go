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
	"strconv"
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
	Short: common.T("Atualiza o GIRUS CLI para a última versão", "Actualiza el GIRUS CLI a la última versión"),
	Long: common.T(`Verifica e atualiza o GIRUS CLI para a última versão disponível.
Após a atualização, oferece a opção de recriar o cluster para garantir
compatibilidade com as novas features.`,
		`Verifica y actualiza el GIRUS CLI a la última versión disponible.
Después de la actualización, ofrece la opción de recrear el cluster para garantizar
compatibilidad con las nuevas funciones.`),
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
		fmt.Printf("%s: %s\n", bold(common.T("Versão atual da CLI", "Versión actual de la CLI")), magenta(currentVersion))

		// Obter última versão do GitHub
		fmt.Println("\n" + headerColor(common.T("Verificando atualizações...", "Verificando actualizaciones...")))
		latestCliVersion, err := GetLatestGitHubVersion(cliRepo)
		if err != nil {
			return fmt.Errorf("%s %s: %v", red(common.T("ERRO:", "ERROR:")), common.T("erro ao verificar última versão da CLI", "error al verificar la última versión de la CLI"), err)
		}

		fmt.Printf("%s: %s\n", bold(common.T("Última versão disponível", "Última versión disponible")), magenta(latestCliVersion))

		// Verificar se já está na versão mais recente
		isLatest := !IsNewerVersion(latestCliVersion, currentVersion)

		if isLatest {
			fmt.Println("\n" + green(common.T("Você já está usando a última versão do GIRUS CLI!", "¡Ya está usando la última versión del GIRUS CLI!")))
			return nil
		}

		// Confirmar atualização
		fmt.Printf("\n%s (%s). %s ",
			yellow(common.T("Nova versão disponível", "Nueva versión disponible")), magenta(latestCliVersion), common.T("Deseja atualizar? (S/n):", "¿Desea actualizar? (S/n):"))
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "s" && response != "" {
			fmt.Println(yellow(common.T("Atualização cancelada.", "Actualización cancelada.")))
			return nil
		}

		// Atualizar CLI
		fmt.Println("\n" + headerColor(common.T("Atualizando CLI...", "Actualizando CLI...")))
		if err := downloadAndInstall(latestCliVersion); err != nil {
			return fmt.Errorf("%s %s: %v", red(common.T("ERRO:", "ERROR:")), common.T("erro ao atualizar CLI", "error al actualizar la CLI"), err)
		}
		fmt.Printf("%s %s %s %s!\n",
			green(common.T("SUCESSO:", "ÉXITO:")), common.T("CLI atualizada com sucesso para a versão", "CLI actualizada con éxito a la versión"), magenta(latestCliVersion), "")

		// Perguntar se deseja recriar o cluster
		fmt.Print("\n" + yellow(common.T("Deseja recriar o cluster para garantir compatibilidade com as novas features? (S/n): ", "¿Desea recrear el cluster para garantizar compatibilidad con las nuevas funcionalidades? (S/n): ")))
		fmt.Scanln(&response)
		if strings.ToLower(response) == "s" || response == "" {
			fmt.Println("\n" + headerColor(common.T("Recriando o cluster...", "Recreando el cluster...")))

			// Executar o comando delete
			deleteCmd := exec.Command("girus", "delete")
			deleteCmd.Stdout = os.Stdout
			deleteCmd.Stderr = os.Stderr
			if err := deleteCmd.Run(); err != nil {
				return fmt.Errorf("%s %s: %v", red(common.T("ERRO:", "ERROR:")), common.T("erro ao deletar o cluster", "error al eliminar el cluster"), err)
			}

			// Executar o comando create
			createCmd := exec.Command("girus", "create")
			createCmd.Stdout = os.Stdout
			createCmd.Stderr = os.Stderr
			if err := createCmd.Run(); err != nil {
				return fmt.Errorf("%s %s: %v", red(common.T("ERRO:", "ERROR:")), common.T("erro ao criar o cluster", "error al crear el cluster"), err)
			}

			fmt.Println("\n" + green(common.T("Cluster recriado com sucesso!", "¡Cluster recreado con éxito!")))
		} else {
			fmt.Println("\n" + yellow(common.T("Cluster mantido como está. Lembre-se que algumas novas features podem não funcionar corretamente com o cluster atual.", "Cluster mantenido como está. Recuerde que algunas nuevas funcionalidades pueden no funcionar correctamente con el cluster actual.")))
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
	// Remove prefixo 'v' se existir e divide as versões
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Comparar cada parte como inteiro (major, minor, patch)
	for i := 0; i < 3; i++ {
		var p1, p2 int
		if i < len(parts1) {
			p1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			p2, _ = strconv.Atoi(parts2[i])
		}
		if p1 > p2 {
			return true
		}
		if p1 < p2 {
			return false
		}
	}

	// São iguais
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
