package status

import (
	"fmt"
)

func PrintLoadingBar(received int64, total int64) {
	percent := GetProgressPercent(received, total)
	fmt.Printf("\r[%.2f%%] ", percent)
}
