package main

import (
	"fmt"
	"scfx/internal/model"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fogleman/ease"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

// General stuff for styling the view

func main() {
	initialModel := model.Init()

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
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
func updateChoices(msg tea.Msg, m model.Model) (tea.Model, tea.Cmd) {
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
func updateChosen(msg tea.Msg, m model.Model) (tea.Model, tea.Cmd) {
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

// Sub-views

// The first view, where you're choosing a task
func choicesView(m model.Model) string {
	c := m.Choice

	tpl := "What to do today?\n\n"
	tpl += "%s\n\n"
	tpl += "Program quits in %s seconds\n\n"
	tpl += style.subtleStyle.Render("j/k, up/down: select") + dotStyle +
		style.subtleStyle.Render("enter: choose") + dotStyle +
		style.subtleStyle.Render("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		checkbox("Python", c == 0),
		checkbox("C++", c == 1),
		checkbox("Golang", c == 2),
		checkbox("NodeJS/TS", c == 3),
	)

	return fmt.Sprintf(tpl, choices, ticksStyle.Render(strconv.Itoa(m.Ticks)))
}

// The second view, after a task has been chosen
func chosenView(m model.Model) string {
	var msg string
	var projectName string = "Hello"
	var projectType string
	switch m.Choice {
	case 0:
		projectType = style.keywordStyle.Render("python")
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
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

// Utils

// Generate a blend of colors.
func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
