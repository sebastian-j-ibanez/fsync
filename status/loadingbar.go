package status

import (
	"fmt"
	"strings"
)

func (p *Progress) getLoadingBar() string {
	p.UpdateProgressPercent()

	bar := ""
	bar += ""
	bar += "\033[32m"
	bar += strings.Repeat("━", int(p.Percentage/5))
	bar += strings.Repeat(" ", 20-int(p.Percentage/5))
	bar += "\033[0m"
	bar += ""

	return bar
}

func (p *Progress) DisplayProgress() {
	progressBar := "\r%3d%% "
	progressBar += p.getLoadingBar()
	progressBar += " " + p.getTimeElapsed()
	fmt.Printf(progressBar, p.Percentage)
}
