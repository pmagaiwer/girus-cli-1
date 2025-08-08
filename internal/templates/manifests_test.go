package templates_test

import (
	"embed"
	"path/filepath"
	"reflect"
	"testing"

	"io/fs"

	"github.com/badtuxx/girus-cli/internal/templates"
	"gopkg.in/yaml.v3"
)

//go:embed manifests/*.yaml
var manifests embed.FS

func TestYAMLFilesAreValid(t *testing.T) {
	err := fs.WalkDir(manifests, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".yaml" {
			data, err := manifests.ReadFile(path)
			if err != nil {
				t.Errorf("Erro ao ler %s: %v", path, err)
				return nil
			}

			var node yaml.Node
			if err := yaml.Unmarshal(data, &node); err != nil {
				t.Errorf("YAML inválido em %s: %v", path, err)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Erro ao percorrer diretório de manifests: %v", err)
	}
}

func TestListAndGetManifests(t *testing.T) {
	entries, err := fs.ReadDir(manifests, "manifests")
	if err != nil {
		t.Fatalf("erro ao ler diretório embed: %v", err)
	}

	var expected []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".yaml" {
			expected = append(expected, entry.Name())
		}
	}

	names, err := templates.ListManifests()
	if err != nil {
		t.Fatalf("ListManifests retornou erro: %v", err)
	}

	if !reflect.DeepEqual(expected, names) {
		t.Errorf("nomes esperados %v, obtidos %v", expected, names)
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			data, err := templates.GetManifest(name)
			if err != nil {
				t.Fatalf("erro ao obter %s: %v", name, err)
			}

			if len(data) == 0 {
				t.Fatalf("%s retornou dados vazios", name)
			}

			var node yaml.Node
			if err := yaml.Unmarshal(data, &node); err != nil {
				t.Errorf("YAML inválido em %s: %v", name, err)
			}
		})
	}
}
