package manifest

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/artemlive/gh-crossplane/internal/domain"
	"gopkg.in/yaml.v3"
)

type ManifestLoader struct {
	dir    string // Directory to load YAML files from
	groups []GroupFile
}

// GroupFile represents a loaded RepositoriesGroup + source path.
type GroupFile struct {
	Path     string
	Manifest domain.RepositoriesGroup
}

func (g GroupFile) Title() string {
	return g.Manifest.Metadata.Name
}

func (g GroupFile) Description() string {
	return fmt.Sprintf("File: %s", filepath.Base(g.Path))
}

func (g GroupFile) FilterValue() string {
	return g.Manifest.Metadata.Name
}

// NewManifestLoader creates a new ManifestLoader instance.
func NewManifestLoader(path string) *ManifestLoader {
	manifestLoader := &ManifestLoader{
		dir: path,
	}
	if err := manifestLoader.LoadGroupsFromFS(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading groups from %s: %v\n", path, err)
	}
	return manifestLoader
}

// LoadGroupsFromDir scans all YAMLs with kind=RepositoriesGroup
func (m *ManifestLoader) LoadGroupsFromFS() error {
	err := filepath.Walk(m.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var g domain.RepositoriesGroup
		if err := yaml.Unmarshal(content, &g); err != nil {
			return fmt.Errorf("unmarshal %s: %w", path, err)
		}

		if g.Kind != "RepositoriesGroup" {
			return nil // ignore unrelated YAMLs
		}

		m.groups = append(m.groups, GroupFile{
			Path:     path,
			Manifest: g,
		})

		return nil
	})

	return err
}

// Groups returns loaded RepositoriesGroup manifests.
func (m *ManifestLoader) Groups() []GroupFile {
	if m.groups == nil {
		if err := m.LoadGroupsFromFS(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading groups: %v\n", err)
			return nil
		}
	}
	return m.groups
}

// GetGroup returns a RepositoriesGroup manifest by its name.
func (m *ManifestLoader) GetGroup(name string) *GroupFile {
	for _, group := range m.groups {
		if group.Manifest.Metadata.Name == name {
			return &group
		}
	}
	return nil
}

func (m *ManifestLoader) SaveGroupFile(gf *GroupFile) error {
	if gf == nil {
		return fmt.Errorf("group file is nil")
	}
	f, err := os.Create(gf.Path)
	if err != nil {
		return fmt.Errorf("create file %s: %w", gf.Path, err)
	}
	defer f.Close() //nolint:errcheck

	enc := yaml.NewEncoder(f)
	enc.SetIndent(1)
	defer enc.Close() //nolint:errcheck
	return enc.Encode(gf.Manifest)
}
