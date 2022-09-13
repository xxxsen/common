package es

import "fmt"

const (
	DefaultVersion = "v1"
)

func Index(table, ver string) (string, string) {
	index := fmt.Sprintf("index_%s", table)
	alias := index + "_" + DefaultVersion
	return index, alias
}
