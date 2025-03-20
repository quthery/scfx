package main

import (
	"fmt"
	"scfx/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)


// General stuff for styling the view

func main() {
	initialModel := model.Init()

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

// Sub-views

// The first view, where you're choosing a task
