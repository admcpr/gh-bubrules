package models

import (
	"gh-bubrls/structs"
	"gh-bubrls/style"

	"github.com/charmbracelet/bubbles/list"
)

func NewOrgListModel(organisations []structs.Organisation, width, height int, user structs.User) list.Model {
	items := make([]list.Item, len(organisations))
	for i, org := range organisations {
		items[i] = structs.NewListItem(org.Login, org.Url)
	}

	list := list.New(items, style.DefaultDelegate, width, height-2)

	list.Title = "User: " + user.Name
	list.SetStatusBarItemName("Organisation", "Organisations")
	list.Styles.Title = style.Title
	list.SetShowTitle(true)

	return list
}
