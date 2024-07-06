package services

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/GoCEP/api/cep/repository"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// "github.com/GoCEP/internal/download"
	"github.com/GoCEP/internal/insertData"
)

type CepService struct {
	repos []repository.CepRepositary
}

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func NewCepService(cepRepository []repository.CepRepositary) *CepService {
	return &CepService{
		repos: cepRepository,
	}
}

type tickMsg time.Time

type model struct {
	progress progress.Model
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		cmd := m.progress.SetPercent(0.9)
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press ctrl+c to cancel the update")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (cepService *CepService) UpdateData(ctx context.Context) error {
	m := model{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}

	return nil

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	dataLocation := dir + os.Getenv("CEP_ZIP_LOCATION")

	dataDir := path.Dir(dataLocation)
	_, err = os.Stat(dataDir)
	if err != nil {
		err = os.MkdirAll(dataDir, 0755)

		if err != nil {
			return fmt.Errorf("failed to mkdir, %s | error: %s", dataDir, err)
		}
	}

	// err = download.File(os.Getenv("CEP_DATA_URL"), dataLocation)
	// if err != nil {
	// 	return fmt.Errorf("failed to download CEP data: %w", err)
	// }

	unprocessedFilesChan := make(chan []string)
	filesJSON := make(chan [][]byte)
	doneZipChan := make(chan bool)
	doneChan := make(chan bool)
	var wg sync.WaitGroup

	wg.Add(1)
	go insertData.UnzipCeps(dataLocation, unprocessedFilesChan, doneZipChan, &wg)

	wg.Add(1)
	go insertData.CleanJSON(unprocessedFilesChan, filesJSON, doneChan, doneZipChan, &wg)

	for _, repo := range cepService.repos {
		wg.Add(1)
		go insertData.InsertToDB(repo, filesJSON, doneChan, &wg)
	}

	wg.Wait()
	close(filesJSON)

	fmt.Println("finished CEP Update")

	return nil
}

func (cepService *CepService) Reset(ctx context.Context) error {
	return nil
}
