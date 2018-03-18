package main

import (
	"regexp"
)

func main() {
	re := regexp.MustCompile(`\[\d+\] ((?P<endpoint>\w+),(?P<user>\w+)(,-*\w*\.*\d*)*)`)
}
