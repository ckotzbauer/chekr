package util

import "fmt"

func (c ComputedValue) FormatMemory() string {
	return fmt.Sprintf("%.2f", c.Percentage*100) + "% (" + ByteCountIEC(c.Value) + ")"
}

func (c ComputedValue) FormatCPU() string {
	return fmt.Sprintf("%.2f", c.Percentage*100) + "% (" + fmt.Sprintf("%.3f", c.Value) + " Cores)"
}

func ByteCountIEC(b float64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%b B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
