package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c ComputedValue) FormatMemory() string {
	if c.Percentage == 0 {
		return ByteCountIEC(c.Value)
	} else {
		return fmt.Sprintf("%.2f", c.Percentage*100) + "% (" + ByteCountIEC(c.Value) + ")"
	}
}

func (c ComputedValue) FormatCPU() string {
	if c.Percentage == 0 {
		return Cores(c.Value)
	} else {
		return fmt.Sprintf("%.2f", c.Percentage*100) + "% (" + Cores(c.Value) + ")"
	}
}

func Cores(b float64) string {
	return fmt.Sprintf("%.3f", b) + " Cores"
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

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func ParseSelector(selector string) []KeyValueSelector {
	r := regexp.MustCompile(`((?:([a-z0-9A-Z\/\-\._]+)(=|!=)?([a-z0-9A-Z\-\._]+)?)+)`)
	matchPairs := r.FindAllStringSubmatch(selector, 10)
	selectors := []KeyValueSelector{}

	for _, matchPair := range matchPairs {
		if len(matchPair) == 5 {
			selectors = append(selectors, KeyValueSelector{
				Key:      matchPair[2],
				Operator: matchPair[3],
				Value:    matchPair[4],
			})
		}
	}

	return selectors
}

func MinOf(vars ...float64) float64 {
	min := vars[0]

	for _, i := range vars {
		if min > i || min == 0 {
			min = i
		}
	}

	return min
}

func MaxOf(vars ...float64) float64 {
	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}

	return max
}

func SumOf(vars ...float64) float64 {
	sum := 0.0

	for _, i := range vars {
		sum += i
	}

	return sum
}
