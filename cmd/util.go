/*
Copyright © 2024 Sebastian Ibanez <sebas.ibanez219@gmail.com>
*/
package cmd

import (
	"strings"
)

func IsGlobPattern(s string) bool {
	return strings.ContainsAny(s, "*?[]")
}
