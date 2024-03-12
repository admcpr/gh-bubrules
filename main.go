package main

import (
	"fmt"
	"gh-bubrls/models"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var MainModel []tea.Model

func main() {
	MainModel = []tea.Model{models.NewUserModel()}

	p := tea.NewProgram(MainModel[0])

	if _, err := p.Run(); err != nil {
		fmt.Printf("Yeah that didn't work, because: %v", err)
		os.Exit(1)
	}

	// fmt.Println("Oh hi, this is the gh-bubrls extension!")
	// client, err := api.DefaultRESTClient()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// response := struct{ Login string }{}
	// err = client.Get("user", &response)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Printf("running as %s\n", response.Login)
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
