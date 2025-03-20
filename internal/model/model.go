package model

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	Choice    int
	Chosen    bool
	Ticks     int
	Frames    int
	Progress  float64
	Loaded    bool
	Quitting  bool
	textInput textinput.Model
	err       error
}

func Init() Model {
	ti := textinput.New()
	ti.Placeholder = "Name of project"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return Model{
		Choice:    0,
		Chosen:    false,
		Ticks:     10,
		Frames:    0,
		Progress:  0,
		Loaded:    false,
		Quitting:  false,
		textInput: ti,
		err:       nil,
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

// Main update function.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	if !m.Chosen {
		return updateChoices(msg, m)
	}
	return updateChosen(msg, m)
}

func (m Model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.Chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return mainStyle.Render("\n" + s + "\n\n")
}
