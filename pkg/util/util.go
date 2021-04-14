package util

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func GetOwnerKind(refs []metav1.OwnerReference) string {
	if len(refs) > 0 {
		return refs[0].Kind
	}

	return ""
}

// Contains checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
