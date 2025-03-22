package model
import(
	"github.com/fogleman/ease"
  "strconv"
  "fmt"
  "os/exec"
  style "scfx/internal/styling"
  tea "github.com/charmbracelet/bubbletea"
)
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
func chosenView(m *Model) string {
  var msg string
    var projectType string

// Инициализация проекта в зависимости от выбора
switch m.Choice {
case 0:
    projectType = style.KeywordStyle.Render("python")
    cmd := exec.Command("uv", "init", m.ProjectName)
    cmd.Stdout = nil
    cmd.Stderr = nil
    cmd.Run() 
    case 1:
        projectType = "C++"
    case 2:
        projectType = "Golang"
    default:
        projectType = "NodeJS/TS"
    }

    msg = fmt.Sprintf(
        "Creating a %s project named %s\n\n",
        projectType,
        m.ProjectName,
    )
    m.Quitting = true
    return msg + "\n\n"
}

