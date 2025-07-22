package field

import (
	"github.com/artemlive/gh-crossplane/internal/domain"
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/artemlive/gh-crossplane/internal/util"
	tea "github.com/charmbracelet/bubbletea"
)

type RepoComponent struct {
	repo    *domain.Repository
	index   int
	focused bool
}

func (r *RepoComponent) Repo() *domain.Repository {
	return r.repo
}

func NewRepoComponent(repo domain.Repository, index int) *RepoComponent {
	return &RepoComponent{
		repo:  &repo,
		index: index,
	}
}

func (r *RepoComponent) View() string {
	prefix := "  "
	if r.focused {
		prefix = "→ "
	}
	name := util.IfEmpty(r.repo.Name, "<unnamed>")
	return prefix + name
}

func (r *RepoComponent) Focus() tea.Cmd {
	r.focused = true
	return nil
}

func (r *RepoComponent) Blur() { r.focused = false }
func (r *RepoComponent) Update(msg tea.Msg, mode ui.FocusMode) (FieldComponent, tea.Cmd) {
	return r, nil
}

func (r *RepoComponent) IsFocused() bool {
	return r.focused
}

func (r *RepoComponent) Init() tea.Cmd {
	return nil
}
