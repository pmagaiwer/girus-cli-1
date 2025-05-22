package repo

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

	// Procura o laboratório pelo ID
	for _, lab := range index.Labs {
		if lab.ID == labName {
			// Se a versão não for especificada, retorna o laboratório encontrado
			if version == "" {
				return &lab, nil
			}

			// Procura a versão específica
			if lab.Version == version {
				return &lab, nil
			}
		}
	}

	return nil, fmt.Errorf("laboratório '%s' não encontrado no repositório '%s'", labName, repoName)
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
	fmt.Printf("Verificando cache em: %s\n", cacheFile)

	// Verifica se o arquivo de cache existe e não está expirado
	if info, err := os.Stat(cacheFile); err == nil {
		// Verifica se o arquivo tem menos de 7 dias
		if time.Since(info.ModTime()) < 7*24*time.Hour {
			fmt.Println("Usando índice do cache")
			data, err := os.ReadFile(cacheFile)
			if err == nil {
				var index Index
				if err := yaml.Unmarshal(data, &index); err == nil {
					return &index, nil
				}
			}
		} else {
			fmt.Println("Cache expirado, baixando novo índice")
		}
	}

	// Se não estiver em cache ou estiver expirado, baixa do repositório
	indexURL := fmt.Sprintf("%s/index.yaml", repo.URL)
	fmt.Printf("Buscando índice em: %s\n", indexURL)
	resp, err := http.Get(indexURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar índice do repositório: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao acessar índice do repositório (status: %d)", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler conteúdo do repositório: %v", err)
	}

	var index Index
	if err := yaml.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("erro ao decodificar índice do repositório: %v", err)
	}

	// Salva no cache
	if err := os.MkdirAll(filepath.Dir(cacheFile), 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de cache: %v", err)
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		return nil, fmt.Errorf("erro ao salvar índice em cache: %v", err)
	}

	return &index, nil
}
