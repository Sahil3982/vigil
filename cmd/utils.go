package cmd

import "strings"

func barFor(value, max float64) string {
	perc := value / max
	width := 10
	filled := int(perc * float64(width))
	empty := width - filled

	bar := strings.Repeat("■", filled) + strings.Repeat("□", empty)
	return "[" + bar + "]"
}