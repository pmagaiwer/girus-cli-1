package repo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"sigs.k8s.io/yaml"
)

// Estrutura para o arquivo index.yaml
type IndexFile struct {
	Labs []Lab `yaml:"labs"`
}

// Estrutura para cada laboratório no index.yaml
type Lab struct {
	ID          string   `yaml:"id"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Version     string   `yaml:"version"`
	Duration    string   `yaml:"duration"`
	Tags        []string `yaml:"tags,omitempty"`
	URL         string   `yaml:"url"`
}

// URL padrão do index.yaml
var DefaultIndexURL = "https://raw.githubusercontent.com/badtuxx/girus-labs/main/index.yaml"

// GetIndexURL retorna a URL do index.yaml, considerando variáveis de ambiente
func GetIndexURL() string {
	// Verificar variável de ambiente
	if url := os.Getenv("GIRUS_REPO_URL"); url != "" {
		return url
	}
	return DefaultIndexURL
}

// GetLabsIndex baixa e parseia o index.yaml remoto
func GetLabsIndex(indexURL string) (*IndexFile, error) {
	// Se não for fornecida uma URL, usar a URL padrão
	if indexURL == "" {
		indexURL = GetIndexURL()
	}

	var data []byte
	var err error

	// Verificar se a URL usa o protocolo file://
	if strings.HasPrefix(indexURL, "file://") {
		// Extrair o caminho do arquivo da URL
		filePath := strings.TrimPrefix(indexURL, "file://")
		// Ler o arquivo local
		data, err = os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler o arquivo local %s: %w", filePath, err)
		}
	} else {
		// Configurar cliente HTTP com timeout
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		// Fazer a requisição HTTP
		resp, err := client.Get(indexURL)
		if err != nil {
			return nil, fmt.Errorf("erro ao acessar o repositório remoto: %w", err)
		}
		defer resp.Body.Close()

		// Verificar se a resposta foi bem-sucedida
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("erro HTTP %d ao acessar o repositório remoto", resp.StatusCode)
		}

		// Ler o conteúdo da resposta
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler o conteúdo do repositório: %w", err)
		}
	}

	// Parsear o YAML
	var index IndexFile
	if err := yaml.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("erro ao parsear o arquivo index.yaml: %w", err)
	}

	return &index, nil
}

// FindLabByID busca um laboratório pelo ID no index.yaml
func FindLabByID(id string, indexURL string) (*Lab, error) {
	index, err := GetLabsIndex(indexURL)
	if err != nil {
		return nil, err
	}

	for _, lab := range index.Labs {
		if lab.ID == id {
			return &lab, nil
		}
	}

	return nil, fmt.Errorf("laboratório com ID '%s' não encontrado no repositório", id)
}

// DownloadLabYAML baixa o arquivo lab.yaml para um arquivo temporário
func DownloadLabYAML(url string) (string, error) {
	// Criar arquivo temporário
	tempFile, err := os.CreateTemp("", "girus-lab-*.yaml")
	if err != nil {
		return "", fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	defer tempFile.Close()

	var data []byte

	// Verificar se a URL usa o protocolo file://
	if strings.HasPrefix(url, "file://") {
		// Extrair o caminho do arquivo da URL
		filePath := strings.TrimPrefix(url, "file://")
		// Ler o arquivo local
		data, err = os.ReadFile(filePath)
		if err != nil {
			os.Remove(tempFile.Name())
			return "", fmt.Errorf("erro ao ler o arquivo local %s: %w", filePath, err)
		}

		// Escrever o conteúdo no arquivo temporário
		if _, err := tempFile.Write(data); err != nil {
			os.Remove(tempFile.Name())
			return "", fmt.Errorf("erro ao salvar o arquivo lab.yaml: %w", err)
		}
	} else {
		// Configurar cliente HTTP com timeout
		client := &http.Client{
			Timeout: 20 * time.Second,
		}

		// Fazer a requisição HTTP
		resp, err := client.Get(url)
		if err != nil {
			os.Remove(tempFile.Name()) // Limpar o arquivo temporário em caso de erro
			return "", fmt.Errorf("erro ao baixar o arquivo lab.yaml: %w", err)
		}
		defer resp.Body.Close()

		// Verificar se a resposta foi bem-sucedida
		if resp.StatusCode != http.StatusOK {
			os.Remove(tempFile.Name())
			return "", fmt.Errorf("erro HTTP %d ao baixar o arquivo lab.yaml", resp.StatusCode)
		}

		// Copiar o conteúdo para o arquivo temporário
		_, err = io.Copy(tempFile, resp.Body)
		if err != nil {
			os.Remove(tempFile.Name())
			return "", fmt.Errorf("erro ao salvar o arquivo lab.yaml: %w", err)
		}
	}

	return tempFile.Name(), nil
}

// FormatTags retorna as tags formatadas como string
func FormatTags(tags []string) string {
	if len(tags) == 0 {
		return "Nenhuma"
	}
	return strings.Join(tags, ", ")
}
