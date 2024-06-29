package download

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

func File(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "go-cep-download-*.tmp")
	if err != nil {
		return err
	}
	defer func() {
		if cerr := tmpFile.Close(); cerr != nil && err == nil {
			err = cerr
		}
		if err != nil {
			os.Remove(tmpFile.Name())
		}
	}()

	text := fmt.Sprintf("[cyan][1/3][reset] Download cep data to %s", tmpFile.Name())

	bar := progressbar.NewOptions(int(resp.ContentLength),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()), //you should install "github.com/k0kubun/go-ansi"
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(text),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	_, err = io.Copy(io.MultiWriter(tmpFile, bar), resp.Body)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filePath); err == nil {
		oldPath := filePath + ".old"

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Remove(oldPath); err != nil {
				return err
			}
			fmt.Printf("\nDeleted existing .old file: %s", oldPath)
		}

		if err := os.Rename(filePath, oldPath); err != nil {
			return err
		}
		fmt.Printf("\nMoved existing file %s to %s", filePath, oldPath)
	}

  fmt.Printf("\nMoved temp file %s to %s", tmpFile.Name(), filePath)
	if err := os.Rename(tmpFile.Name(), filePath); err != nil {
		return err
	}

	return nil
}
