package model

import (
	"fmt"
  style "scfx/internal/styling"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fogleman/ease"
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
	return style.MainStyle.Render("\n" + s + "\n\n")
}

type (
	tickMsg  struct{}
	frameMsg struct{}
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice++
			if m.Choice > 3 {
				m.Choice = 3
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			return m, frame()
		}

	case tickMsg:
		if m.Ticks == 0 {
			m.Quitting = true
			return m, tea.Quit
		}
		m.Ticks--
		return m, tick()
	}

	return m, nil
}

// Update loop for the second view after a choice has been made
func updateChosen(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case frameMsg:
		if !m.Loaded {
			m.Frames++
			m.Progress = ease.OutBounce(float64(m.Frames) / float64(100))
			if m.Progress >= 1 {
				m.Progress = 1
				m.Loaded = true
				m.Ticks = 3
				return m, tick()
			}
			return m, frame()
		}

	case tickMsg:
		if m.Loaded {
			if m.Ticks == 0 {
				m.Quitting = true
				return m, tea.Quit
			}
			m.Ticks--
			return m, tick()
		}
	}

	return m, nil
}

func choicesView(m Model) string {
	c := m.Choice

	tpl := "What to do today?\n\n"
	tpl += "%s\n\n"
	tpl += "Program quits in %s seconds\n\n"
	tpl += style.SubtleStyle.Render("j/k, up/down: select") + style.DotStyle +
		style.SubtleStyle.Render("enter: choose") + style.DotStyle +
		style.SubtleStyle.Render("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		checkbox("Python", c == 0),
		checkbox("C++", c == 1),
		checkbox("Golang", c == 2),
		checkbox("NodeJS/TS", c == 3),
	)

	return fmt.Sprintf(tpl, choices, style.TicksStyle.Render(strconv.Itoa(m.Ticks)))
}

// The second view, after a task has been chosen
func chosenView(m Model) string {
	var msg string
	var projectName string = "Hello"
	var projectType string
	switch m.Choice {
	case 0:
		projectType = style.KeywordStyle.Render("python")
	case 1:
		projectType = "C++"
	case 2:
		projectType = "Golang"
	default:
		projectType = "NodeJS/TS"

	}

	msg = fmt.Sprintf(
		"Creating a %s project in %s\n\n",
		"Enter project name:",
		projectType,
		projectName,
	)
	return msg + "\n\n"
}

func checkbox(label string, checked bool) string {
	if checked {
		return style.CheckboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
