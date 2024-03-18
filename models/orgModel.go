package models

import (
	"fmt"
	"gh-bubrls/structs"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type tickMsg time.Time
type orgQueryMsg structs.OrganizationQuery
type repoQueryMsg structs.RepositoryQuery

type OrgModel struct {
	progress  progress.Model
	login     string
	repoCount int
	repos     []structs.Repository
}

func NewOrgModel(login string) OrgModel {
	orgModel := OrgModel{}
	orgModel.login = login
	orgModel.repoCount = 0
	orgModel.repos = []structs.Repository{}
	orgModel.progress = progress.New(progress.WithDefaultGradient())
	return orgModel
}

func (m OrgModel) Init() tea.Cmd {
	return tea.Batch(tickCmd(), getRepos(m.login))

}

func (m OrgModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.01)
		return m, tea.Batch(tickCmd(), cmd)

	case orgQueryMsg:
		repos := msg.Organization.Repositories.Edges
		cmds := []tea.Cmd{m.progress.SetPercent(0.1)}
		m.repoCount = len(msg.Organization.Repositories.Edges)
		for _, repo := range repos {
			cmds = append(cmds, getRepo(m.login, repo.Node.Name))
		}
		return m, tea.Batch(cmds...)

	case repoQueryMsg:
		m.repos = append(m.repos, msg.Repository)
		cmd := m.progress.IncrPercent(0.9 / float64(m.repoCount))

		return m, cmd

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m OrgModel) View() string {
	pad := strings.Repeat(" ", padding)
	progress := "\n" + pad + m.progress.View() + "\n\n" + pad + "Getting repositories ... "
	if m.repoCount < 1 {
		return progress
	}
	return progress + fmt.Sprintf("%d of %d", len(m.repos), m.repoCount)
}

func getRepo(owner string, name string) tea.Cmd {
	return func() tea.Msg {
		client, err := api.DefaultGraphQLClient()
		if err != nil {
			log.Fatal(err)
		}
		repoQuery := structs.RepositoryQuery{}

		variables := map[string]interface{}{
			"owner": graphql.String(owner),
			"name":  graphql.String(name),
		}
		err = client.Query("Repository", &repoQuery, variables)
		if err != nil {
			log.Fatal(err)
		}
		return repoQueryMsg(repoQuery)
	}
}

func getRepos(login string) tea.Cmd {
	return func() tea.Msg {
		client, err := api.DefaultGraphQLClient()
		if err != nil {
			log.Fatal(err)
		}
		organizationQuery := structs.OrganizationQuery{}

		variables := map[string]interface{}{
			"login": graphql.String(login),
			"first": graphql.Int(100),
		}
		err = client.Query("OrganizationRepositories", &organizationQuery, variables)
		if err != nil {
			log.Fatal(err)
		}
		return orgQueryMsg(organizationQuery)
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
