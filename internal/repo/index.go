package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"gopkg.in/yaml.v3"
)

// Repository representa um repositório de laboratórios
type Repository struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

// Index representa o arquivo de índice de um repositório
type Index struct {
	APIVersion string                 `yaml:"apiVersion"`
	Generated  string                 `yaml:"generated"`
	Entries    map[string][]LabEntry  `yaml:"entries"`
}

// LabEntry representa um laboratório no índice
type LabEntry struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Description string   `yaml:"description"`
	Keywords    []string `yaml:"keywords"`
	Maintainers []string `yaml:"maintainers"`
	URL         string   `yaml:"url"`
	Created     string   `yaml:"created"`
	Digest      string   `yaml:"digest"`
}

// RepositoryManager gerencia os repositórios de laboratórios
type RepositoryManager struct {
	configPath string
	repos      map[string]Repository
}

// NewRepositoryManager cria uma nova instância do gerenciador de repositórios
func NewRepositoryManager() (*RepositoryManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter diretório home: %v", err)
	}

	configPath := filepath.Join(homeDir, ".girus", "repositories.json")
	rm := &RepositoryManager{
		configPath: configPath,
		repos:      make(map[string]Repository),
	}

	// Carrega repositórios existentes
	if err := rm.loadRepositories(); err != nil {
		return nil, err
	}

	return rm, nil
}

// AddRepository adiciona um novo repositório
func (rm *RepositoryManager) AddRepository(name, url, description string) error {
	// Verifica se o repositório já existe
	if _, exists := rm.repos[name]; exists {
		return fmt.Errorf("repositório '%s' já existe", name)
	}

	// Valida o repositório
	if err := rm.validateRepository(url); err != nil {
		return fmt.Errorf("repositório inválido: %v", err)
	}

	// Adiciona o repositório
	rm.repos[name] = Repository{
		Name:        name,
		URL:         url,
		Description: description,
		Version:     "v1",
	}

	// Salva as alterações
	return rm.saveRepositories()
}

// RemoveRepository remove um repositório
func (rm *RepositoryManager) RemoveRepository(name string) error {
	if _, exists := rm.repos[name]; !exists {
		return fmt.Errorf("repositório '%s' não encontrado", name)
	}

	delete(rm.repos, name)
	return rm.saveRepositories()
}

// ListRepositories lista todos os repositórios
func (rm *RepositoryManager) ListRepositories() []Repository {
	repos := make([]Repository, 0, len(rm.repos))
	for _, repo := range rm.repos {
		repos = append(repos, repo)
	}
	return repos
}

// GetRepository obtém um repositório específico
func (rm *RepositoryManager) GetRepository(name string) (Repository, error) {
	repo, exists := rm.repos[name]
	if !exists {
		return Repository{}, fmt.Errorf("repositório '%s' não encontrado", name)
	}
	return repo, nil
}

// UpdateRepository atualiza um repositório existente
func (rm *RepositoryManager) UpdateRepository(name, url, description string) error {
	if _, exists := rm.repos[name]; !exists {
		return fmt.Errorf("repositório '%s' não encontrado", name)
	}

	// Valida o repositório
	if err := rm.validateRepository(url); err != nil {
		return fmt.Errorf("repositório inválido: %v", err)
	}

	rm.repos[name] = Repository{
		Name:        name,
		URL:         url,
		Description: description,
		Version:     "v1",
	}

	return rm.saveRepositories()
}

// loadRepositories carrega os repositórios do arquivo de configuração
func (rm *RepositoryManager) loadRepositories() error {
	// Cria o diretório se não existir
	dir := filepath.Dir(rm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de configuração: %v", err)
	}

	// Verifica se o arquivo existe
	if _, err := os.Stat(rm.configPath); os.IsNotExist(err) {
		return nil
	}

	// Lê o arquivo
	data, err := os.ReadFile(rm.configPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de configuração: %v", err)
	}

	// Decodifica o JSON
	if err := json.Unmarshal(data, &rm.repos); err != nil {
		return fmt.Errorf("erro ao decodificar arquivo de configuração: %v", err)
	}

	return nil
}

// saveRepositories salva os repositórios no arquivo de configuração
func (rm *RepositoryManager) saveRepositories() error {
	// Codifica para JSON
	data, err := json.MarshalIndent(rm.repos, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao codificar configuração: %v", err)
	}

	// Salva no arquivo
	if err := os.WriteFile(rm.configPath, data, 0644); err != nil {
		return fmt.Errorf("erro ao salvar arquivo de configuração: %v", err)
	}

	return nil
}

// validateRepository valida se um repositório é acessível e válido
func (rm *RepositoryManager) validateRepository(url string) error {
	_, err := fetchAndParseIndex(url)
	if err != nil {
		return fmt.Errorf("falha ao validar repositório: %v", err)
	}
	return nil
}

// fetchAndParseIndex baixa e parseia o arquivo index.yaml de um repositório
func fetchAndParseIndex(url string) (*Index, error) {
	// Se a URL usa o protocolo file://
	if strings.HasPrefix(url, "file://") {
		filePath := strings.TrimPrefix(url, "file://")
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler arquivo local: %v", err)
		}

		var index Index
		if err := yaml.Unmarshal(data, &index); err != nil {
			return nil, fmt.Errorf("erro ao decodificar arquivo local: %v", err)
		}
		return &index, nil
	}

	// Para URLs HTTP/HTTPS
	resp, err := http.Get(fmt.Sprintf("%s/index.yaml", url))
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar repositório: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao acessar repositório (status: %d)", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler conteúdo do repositório: %v", err)
	}

	var index Index
	if err := yaml.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("erro ao decodificar índice do repositório: %v", err)
	}

	return &index, nil
}