package status

import (
	"fmt"
	"strings"
)

func (p *Progress) getLoadingBar() string {
	p.UpdateProgressPercent()

	bar := ""
	bar += "["
	bar += strings.Repeat("-", int(p.Percentage/5))
	bar += strings.Repeat(" ", 20-int(p.Percentage/5))
	bar += "]"

	return bar
}

func (p *Progress) DisplayProgress() {
	progressBar := "\r[%3d%%] "
	progressBar += p.getLoadingBar()
	progressBar += " " + p.getTimeElapsed()
	fmt.Printf(progressBar, p.Percentage)
}
