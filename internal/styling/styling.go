package styling

import (
  "fmt"
  "strconv"
  "github.com/lucasb-eyer/go-colorful"

  "github.com/charmbracelet/lipgloss")

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return
}
var (
	KeywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	SubtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	TicksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	CheckboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	DotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	MainStyle     = lipgloss.NewStyle().MarginLeft(2)

	// Gradient colors we'll use for the progress bar
	Ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
)
