package services

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/GoCEP/api/cep/repository"
	"github.com/GoCEP/internal/insertData"
	"github.com/GoCEP/internal/internalProgressbar"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type CepService struct {
	repos []repository.CepRepositary
}

func NewCepService(repos []repository.CepRepositary) CepService {
	return CepService{
		repos: repos,
	}
}

func (cepService *CepService) UpdateData(ctx context.Context) error {
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

	/*
		err = download.File(os.Getenv("CEP_DATA_URL"), dataLocation)
		if err != nil {
			return fmt.Errorf("failed to download CEP data: %w", err)
		}
	*/

	filesJSON := make(chan [][]byte)
	doneZipChan := make(chan bool)
	doneChan := make(chan bool)
	var wg sync.WaitGroup

	unzipBar := internalProgressbar.Bar{
    Index: 0,
		ProgressBar:     progress.New(),
		Message:         "Extraindo arquivos",
		FinishedMessage: "Arquivos extraindos com sucesso",
	}

	insertBar := internalProgressbar.Bar{
    Index: 1,
		ProgressBar:     progress.New(),
		Message:         "Inserindo arquivos no banco",
		FinishedMessage: "Arquivos extraidos com sucesso",
	}

	bars := []internalProgressbar.Bar{unzipBar, insertBar}

	progressBars := internalProgressbar.NewProgressBar(bars)
	p := tea.NewProgram(progressBars)

	setUnzipProgress := func(percentage float64) {
		p.Send(internalProgressbar.ProgressMsg{Index: 0, Value: percentage})
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		insertData.UnzipCeps(dataLocation, filesJSON, doneZipChan, setUnzipProgress)
	}()

	setInsertOnRepoProgress := func(percentage float64, index int) {
		p.Send(internalProgressbar.ProgressMsg{Index: index, Value: percentage})
	}

	for i, repo := range cepService.repos {
		wg.Add(1)
		setProgress := func(percentage float64) {
			setInsertOnRepoProgress(percentage, i+1)
		}
		go func(repo repository.CepRepositary) {
			defer wg.Done()
			insertData.InsertToDB(repo, filesJSON, doneChan, setProgress)
		}(repo)
	}

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running progress bars: %w", err)
	}

	defer wg.Wait()
	defer close(filesJSON)

	return nil
}

func (cepService *CepService) Reset(ctx context.Context) error {
	return nil
}
