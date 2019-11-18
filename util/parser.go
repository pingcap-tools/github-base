package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseStringSlice convert string to slice
func ParseStringSlice(input string) []string {
	return strings.Split(input, ",")
}

// ParseIntSlice convert string to int slice
func ParseIntSlice(input string) []int {
	var s []int
	err := json.Unmarshal([]byte(fmt.Sprintf("[%s]", input)), &s)
	if err != nil {
		return []int{}
	}
	return s
}

// EncodeStringSlice convert slice to string
func EncodeStringSlice(input []string) string {
	return strings.Join(input, ",")
}
