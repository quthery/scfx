package model

import (
	"fmt"
  style "scfx/internal/styling"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
      ProjectName string // Добавляем поле для хранения названия проекта
    IsInputtingName bool // Флаг для отслеживания, вводит ли пользователь название проекта
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
    if m.Quitting {
        return m, tea.Quit
    }
    if !m.Chosen {
        return updateChoices(msg, m)
    }

    if m.IsInputtingName {
        // Обрабатываем ввод названия проекта
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "enter":
                // Сохраняем введенное название проекта
                m.ProjectName = m.textInput.Value()
                m.IsInputtingName = false
                m.Loaded = true
                return m, frame()
            }
        }

        // Обновляем текстовый ввод
        var cmd tea.Cmd
        m.textInput, cmd = m.textInput.Update(msg)
        return m, cmd
    }

    return updateChosen(msg, m)
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
func (m Model) View() string {
    var s string
    if m.Quitting {
        return "\n  See you later!\n\n"
    }
    if !m.Chosen {
        s = choicesView(m)
    } else if m.IsInputtingName {
        s = fmt.Sprintf(
            "Enter project name:\n\n%s\n\n",
            m.textInput.View(),
        )
    } else {
        s = chosenView(&m)
    }
    return style.MainStyle.Render("\n" + s + "\n\n")
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
            m.IsInputtingName = true // Переключаемся на ввод названия проекта
            return m, nil
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

func checkbox(label string, checked bool) string {
	if checked {
		return style.CheckboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
