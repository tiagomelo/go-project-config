// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cfg

import (
	"go/format"

	"github.com/pkg/errors"
)

// formatter is an interface for formatting Go source code.
type formatter interface {
	// Source formats the provided Go source code.
	Source(src []byte) ([]byte, error)
}

// coreFormatter implements the formatter interface using the default Go formatting.
type coreFormatter struct{}

func (c coreFormatter) Source(src []byte) ([]byte, error) {
	return format.Source(src)
}

// For ease of unit testing.
var formatterProvider formatter = coreFormatter{}

// formatGoFile formats the Go source code in the specified file.
func formatGoFile(filePath string) error {
	source, err := fsProvider.ReadFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "reading go file %s", filePath)
	}
	formattedSrc, err := formatterProvider.Source(source)
	if err != nil {
		return errors.Wrapf(err, "formating go file %s", filePath)
	}
	err = fsProvider.WriteFile(filePath, formattedSrc, 0644)
	if err != nil {
		return errors.Wrapf(err, "writing go file %s", filePath)
	}
	return nil
}
