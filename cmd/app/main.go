package main

import (
	"flag"
	"os"

	"github.com/artemlive/gh-crossplane/internal/app"
	tea "github.com/charmbracelet/bubbletea/v2"
)

func main() {
	groupsDir := flag.String("groups-dir", "flux/resources/github/management/repositories", "Path to the directory with RepositoriesGroup YAMLs")
	flag.Parse()
	p := tea.NewProgram(app.NewAppModel(*groupsDir), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
