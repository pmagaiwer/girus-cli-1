package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// LabManager gerencia os laboratórios
type LabManager struct {
	repoManager *RepositoryManager
	cachePath   string
}

// NewLabManager cria uma nova instância do gerenciador de laboratórios
func NewLabManager(repoManager *RepositoryManager) (*LabManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter diretório home: %v", err)
	}

	cachePath := filepath.Join(homeDir, ".girus", "cache")
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de cache: %v", err)
	}

	return &LabManager{
		repoManager: repoManager,
		cachePath:   cachePath,
	}, nil
}

// ListLabs lista todos os laboratórios disponíveis em todos os repositórios
func (lm *LabManager) ListLabs() (map[string][]LabEntry, error) {
	allLabs := make(map[string][]LabEntry)

	for _, repo := range lm.repoManager.ListRepositories() {
		index, err := lm.getIndex(repo)
		if err != nil {
			return nil, fmt.Errorf("erro ao obter índice do repositório %s: %v", repo.Name, err)
		}

		for name, entries := range index.Entries {
			allLabs[name] = append(allLabs[name], entries...)
		}
	}

	return allLabs, nil
}

// GetLab obtém um laboratório específico
func (lm *LabManager) GetLab(repoName, labName, version string) (*LabEntry, error) {
	repo, err := lm.repoManager.GetRepository(repoName)
	if err != nil {
		return nil, err
	}

	index, err := lm.getIndex(repo)
	if err != nil {
		return nil, err
	}

	entries, exists := index.Entries[labName]
	if !exists {
		return nil, fmt.Errorf("laboratório '%s' não encontrado no repositório '%s'", labName, repoName)
	}

	// Se a versão não for especificada, retorna a mais recente
	if version == "" {
		return &entries[0], nil
	}

	// Procura a versão específica
	for _, entry := range entries {
		if entry.Version == version {
			return &entry, nil
		}
	}

	return nil, fmt.Errorf("versão '%s' do laboratório '%s' não encontrada no repositório '%s'", version, labName, repoName)
}

// DownloadLab baixa um laboratório específico
func (lm *LabManager) DownloadLab(repoName, labName, version string) error {
	lab, err := lm.GetLab(repoName, labName, version)
	if err != nil {
		return err
	}

	// Cria o diretório do laboratório
	labPath := filepath.Join(lm.cachePath, repoName, labName, lab.Version)
	if err := os.MkdirAll(labPath, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do laboratório: %v", err)
	}

	// Baixa o arquivo do laboratório
	resp, err := http.Get(lab.URL)
	if err != nil {
		return fmt.Errorf("erro ao baixar laboratório: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao baixar laboratório (status: %d)", resp.StatusCode)
	}

	// Salva o arquivo
	labFile := filepath.Join(labPath, "lab.yaml")
	out, err := os.Create(labFile)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo do laboratório: %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("erro ao salvar laboratório: %v", err)
	}

	return nil
}

// getIndex obtém o índice de um repositório
func (lm *LabManager) getIndex(repo Repository) (*Index, error) {
	// Tenta obter do cache primeiro
	cacheFile := filepath.Join(lm.cachePath, repo.Name, "index.yaml")
	if data, err := os.ReadFile(cacheFile); err == nil {
		var index Index
		if err := json.Unmarshal(data, &index); err == nil {
			return &index, nil
		}
	}

	// Se não estiver em cache, baixa do repositório
	resp, err := http.Get(fmt.Sprintf("%s/index.yaml", repo.URL))
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar índice do repositório: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao acessar índice do repositório (status: %d)", resp.StatusCode)
	}

	var index Index
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return nil, fmt.Errorf("erro ao decodificar índice do repositório: %v", err)
	}

	// Salva no cache
	if err := os.MkdirAll(filepath.Dir(cacheFile), 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de cache: %v", err)
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao codificar índice: %v", err)
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		return nil, fmt.Errorf("erro ao salvar índice em cache: %v", err)
	}

	return &index, nil
} 