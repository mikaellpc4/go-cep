package internalProgressbar

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding  = 2
	maxWidth = 80
)

type ProgressMsg struct {
	Index int
	Value float64
}

type Bar struct {
	Index           int
	ProgressBar     progress.Model
	Message         string
	FinishedMessage string
}

type progressBar struct {
	bars []Bar
	err  error
}

func NewProgressBar(bars []Bar) progressBar {
	return progressBar{
		bars: bars,
	}
}

func (m progressBar) Init() tea.Cmd {
	return nil
}

func (m progressBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		for i := range m.bars {
			m.bars[i].ProgressBar.Width = msg.Width - padding*2 - 4
			if m.bars[i].ProgressBar.Width > maxWidth {
				m.bars[i].ProgressBar.Width = maxWidth
			}
		}
		return m, nil

	case ProgressMsg:
		var cmds []tea.Cmd

		stop := true
		for i := range m.bars {
			if m.bars[i].ProgressBar.Percent() < 1.0 {
				stop = false
			}
		}

		if stop {
			cmds = append(cmds, tea.Quit)
		}

		cmds = append(cmds, m.bars[msg.Index].ProgressBar.SetPercent(msg.Value))

		return m, tea.Batch(cmds...)

	case progress.FrameMsg:
		var cmd tea.Cmd
		for i := range m.bars {
			var progressModel tea.Model
			progressModel, cmd = m.bars[i].ProgressBar.Update(msg)
			m.bars[i].ProgressBar = progressModel.(progress.Model)
		}
		return m, cmd

	default:
		return m, nil
	}
}

func (m progressBar) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	pad := strings.Repeat(" ", padding)
	var progressViews []string
	for _, bar := range m.bars {
		progressViews = append(progressViews, pad+bar.ProgressBar.View())
	}

	return "\n" +
		strings.Join(progressViews, "\n\n") + "\n\n" +
		pad + helpStyle("Press Ctrl+C to stop")
}
