package utils

import "strconv"

func ParseUint(s string) (uint, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	return uint(i), err
}
