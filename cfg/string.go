// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cfg

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// toCamelCase converts a string to camel case.
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	c := cases.Title(language.Und)
	for i, part := range parts {
		parts[i] = c.String((strings.ToLower(part)))
	}
	return strings.Join(parts, "")
}
