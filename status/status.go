package status

func GetProgressPercent(received int64, total int64) float64 {
	return ((float64(received) / float64(total)) * 100)
}
