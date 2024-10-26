package status

import (
	"fmt"
	"time"
)

type Progress struct {
	TimeElapsed    int64
	Percentage     int64
	BytesReceived  int64
	TotalFileBytes int64
}

// Get progress percentage
func (p *Progress) UpdateProgressPercent() {
	p.Percentage = int64((float64(p.BytesReceived) / float64(p.TotalFileBytes)) * 100)
}

func (p *Progress) getTimeElapsed() string {
	timeElapsed := time.Now().Unix() - p.TimeElapsed
	minutes := timeElapsed / 60
	seconds := timeElapsed % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
