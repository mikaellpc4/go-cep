package progressBar

import (
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

func Create(size int, text string) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(size,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()), //you should install "github.com/k0kubun/go-ansi"
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(text),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]█[reset]",
			SaucerHead:    "[green]▌[reset]",
			SaucerPadding: " ",
			BarStart:      "|",
			BarEnd:        "|",
		}))

	return bar
}
