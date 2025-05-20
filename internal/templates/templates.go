package templates

import (
	"embed"
	"io/fs"
	"path/filepath"
)

//go:embed manifests/*.yaml
var ManifestFS embed.FS

func GetManifest(name string) ([]byte, error) {
	return fs.ReadFile(ManifestFS, "manifests/"+name)
}

// ListManifests retorna uma lista de nomes de todos os arquivos YAML no diretório de manifests
func ListManifests() ([]string, error) {
	entries, err := fs.ReadDir(ManifestFS, "manifests")
	if err != nil {
		return nil, err
	}

	var manifests []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".yaml" {
			manifests = append(manifests, entry.Name())
		}
	}

	return manifests, nil
}
