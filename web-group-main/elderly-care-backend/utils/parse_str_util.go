package utils

import (
	"strconv"
	"strings"
)

func ParseIntSlice(s string) []uint {
	var result []uint
	for _, v := range strings.Split(s, ",") {
		i, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		result = append(result, uint(i))
	}
	return result
}

func ParseStringSlice(s string) []string {
	return strings.Split(s, ",")
}
