// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cfg

import (
	"io"
	"io/fs"
	"strings"
)

type mockFile struct {
	name           string
	writeErr       error
	writeStringErr error
	readErr        error
	closeErr       error
}

func (m *mockFile) WriteString(s string) (n int, err error) {
	return 5, m.writeStringErr
}

func (m *mockFile) Name() string {
	return m.name
}

func (m *mockFile) Read(p []byte) (n int, err error) {
	return 5, m.readErr
}

func (m *mockFile) Write(p []byte) (n int, err error) {
	return 5, m.writeErr
}

func (m *mockFile) Close() error {
	return m.closeErr
}

type mockFileSystem struct {
	file                  []byte
	isNotExistOutput      bool
	openedFile            *mockFile
	createdFile           *mockFile
	mkDirErr              error
	openErr               error
	createErr             error
	createUnitTestFileErr error
	createEnvFileErr      error
	readFileErr           error
	writeFileErr          error
}

func (m *mockFileSystem) Mkdir(name string) error {
	return m.mkDirErr
}

func (m *mockFileSystem) Open(name string) (File, error) {
	return m.openedFile, m.openErr
}

func (m *mockFileSystem) Create(name string) (File, error) {
	if strings.Contains(name, configReadFileName) {
		return m.createdFile, m.createErr
	}
	if strings.Contains(name, configReaderUnitTestFileName) {
		return m.createdFile, m.createUnitTestFileErr
	}
	if strings.Contains(name, envFileName) {
		return m.createdFile, m.createEnvFileErr
	}
	return m.createdFile, nil
}

func (m *mockFileSystem) IsNotExist(err error) bool {
	return m.isNotExistOutput
}

func (m *mockFileSystem) ReadFile(name string) ([]byte, error) {
	return m.file, m.readFileErr
}

func (m *mockFileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return m.writeFileErr
}

type mockFormatter struct {
	file      []byte
	sourceErr error
}

func (m *mockFormatter) Source(src []byte) ([]byte, error) {
	return m.file, m.sourceErr
}

type mockTemplateExecutor struct {
	err error
}

func (m *mockTemplateExecutor) Execute(wr io.Writer, data interface{}) error {
	return m.err
}

type mockTemplateProcessor struct {
	te                               templateExecutor
	err                              error
	envFileTemplateParseErr          error
	configReaderUnitTestFileParseErr error
}

func (m *mockTemplateProcessor) Parse(name, text string) (templateExecutor, error) {
	if strings.Contains(name, configReaderMainFileTemplateName) {
		return m.te, m.err
	}
	if strings.Contains(name, configReaderUnitTestFileTemplateName) {
		return m.te, m.configReaderUnitTestFileParseErr
	}
	if strings.Contains(name, envFileTemplateName) {
		return m.te, m.envFileTemplateParseErr
	}
	return m.te, m.err
}

type mockLineReader struct {
	lines []string
	index int
	err   error
}

func (m *mockLineReader) Scan() bool {
	return m.index < len(m.lines)
}

func (m *mockLineReader) Text() string {
	line := m.lines[m.index]
	m.index++
	return line
}

func (m *mockLineReader) Err() error {
	return m.err
}
