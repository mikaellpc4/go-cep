package download

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dustin/go-humanize"
)

type WriteCounter struct {
  InitialBytes bool
	TotalWritten uint64
	DataTotal    uint64
}

func (wc *WriteCounter) Write(data []byte) (int, error) {
	bytesWriten := len(data)
	wc.TotalWritten += uint64(bytesWriten)
	wc.PrintProgress(&wc.InitialBytes)

	return bytesWriten, nil
}

func (wc WriteCounter) PrintProgress(initial *bool) {
	percent := float64(wc.TotalWritten) / float64(wc.DataTotal) * 100
	written := humanize.Bytes(wc.TotalWritten)
	total := humanize.Bytes(wc.DataTotal)

	if !*initial {
		fmt.Printf("\033[1A\033[K")
	} else {
    *initial = false
  }

	fmt.Printf("\rDownloding %s/%s, %.2f%% complete\n", written, total, percent)
}

func File(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	progressCounter := &WriteCounter{DataTotal: uint64(resp.ContentLength), InitialBytes: true}

	_, err = io.Copy(io.MultiWriter(out, progressCounter), resp.Body)
	if err != nil {
		return err
	}

	return nil
}
