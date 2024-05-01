// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cfg

import (
	"io"
	"text/template"
)

// templateExecutor interface abstracts the execution of a parsed template.
// It requires an Execute method that writes the executed template to an io.Writer.
type templateExecutor interface {
	Execute(wr io.Writer, data interface{}) error
}

// templateProcessor interface abstracts the parsing of a template.
// It requires a Parse method that takes a template name and text, and returns
// a templateExecutor and an error, if any.
type templateProcessor interface {
	Parse(name, text string) (templateExecutor, error)
}

// textTemplateProcessor struct is an empty struct that implements the
// templateProcessor interface using Go's text/template package.
type textTemplateProcessor struct{}

// Parse implements the templateProcessor interface. It creates a new text
// template with the provided name and text and returns an textTemplateExecutor.
func (textTemplateProcessor) Parse(name, text string) (templateExecutor, error) {
	tmpl, err := template.New(name).Parse(text)
	if err != nil {
		return nil, err
	}
	return textTemplateExecutor{tmpl}, nil
}

// textTemplateExecutor struct holds a reference to a parsed text template.
type textTemplateExecutor struct {
	tmpl *template.Template
}

// Execute implements the templateExecutor interface. It executes the template
// using the provided data and writes the output to the specified io.Writer.
func (r textTemplateExecutor) Execute(wr io.Writer, data interface{}) error {
	return r.tmpl.Execute(wr, data)
}
