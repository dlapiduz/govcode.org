package main

import (
	"strconv"
)

func ForceStoInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
