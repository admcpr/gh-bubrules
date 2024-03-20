package models

import (
	"fmt"
	"gh-bubrls/structs"
	"gh-bubrls/style"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

const (
	padding  = 2
	maxWidth = 80
)

type orgQueryMsg structs.OrganizationQuery
type repoQueryMsg structs.RepositoryQuery

type OrgModel struct {
	progress  progress.Model
	login     string
	repoCount int
	repos     []structs.Repository
	list      list.Model
}

func NewOrgModel(login string) OrgModel {
	orgModel := OrgModel{
		list:      list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		login:     login,
		repoCount: 0,
		repos:     []structs.Repository{},
		progress:  progress.New(progress.WithDefaultGradient()),
	}
	return orgModel
}

func (m OrgModel) Init() tea.Cmd {
	return getRepoList(m.login)

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

	case orgQueryMsg:
		repos := msg.Organization.Repositories.Edges
		cmds := []tea.Cmd{m.progress.SetPercent(0.1)}
		m.repoCount = len(msg.Organization.Repositories.Edges)
		for _, repo := range repos {
			cmds = append(cmds, getRepoDetails(m.login, repo.Node.Name))
		}
		return m, tea.Batch(cmds...)

	case repoQueryMsg:
		m.repos = append(m.repos, msg.Repository)
		cmd := m.progress.IncrPercent(0.9 / float64(m.repoCount))

		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m OrgModel) View() string {
	if m.progress.Percent() < 1.0 {
		return m.ProgressView()
	}

	for idx, repo := range m.repos {
		m.list.InsertItem(idx, structs.NewListItem(repo.Name, repo.Url))
	}

	return style.App.Render(m.list.View())
}

func (m OrgModel) ProgressView() string {
	pad := strings.Repeat(" ", padding)
	progress := "\n" + pad + m.progress.View() + "\n\n" + pad + "Getting repositories ... "
	if m.repoCount < 1 {
		return progress
	}
	return progress + fmt.Sprintf("%d of %d", len(m.repos), m.repoCount)
}

func getRepoDetails(owner string, name string) tea.Cmd {
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

func getRepoList(login string) tea.Cmd {
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
