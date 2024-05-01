// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cfg

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateConfigPackage(t *testing.T) {
	testCases := []struct {
		name           string
		mockClosure    func(mfs *mockFileSystem, mtp *mockTemplateProcessor)
		expectedOutput []string
		expectedError  error
	}{
		{
			name: "happy path",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor) {
				mfs.createdFile = new(mockFile)
				mtp.te = new(mockTemplateExecutor)
			},
			expectedOutput: []string{
				"config/config.go",
				"config/config_test.go",
				".env",
			},
		},
		{
			name: "error when creating config files dir",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor) {
				mfs.mkDirErr = errors.New("mkdir error")
			},
			expectedError: errors.New("creating dir config: mkdir error"),
		},
		{
			name: "error when creating config reader file",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor) {
				mfs.createErr = errors.New("create error")
			},
			expectedError: errors.New("creating file config/config.go: create error"),
		},
		{
			name: "error when writing to config reader file, template parse error",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor) {
				mfs.createdFile = new(mockFile)
				mtp.err = errors.New("parse error")
			},
			expectedError: errors.New("parsing template configReaderMainFile: parse error"),
		},
		{
			name: "error when writing to config reader file, template execute error",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor) {
				mfs.createdFile = new(mockFile)
				te := new(mockTemplateExecutor)
				te.err = errors.New("execute error")
				mtp.te = te
			},
			expectedError: errors.New("executing template configReaderMainFile: execute error"),
		},
		{
			name: "error when creating config reader unit test file",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor) {
				mfs.createdFile = new(mockFile)
				mtp.te = new(mockTemplateExecutor)
				mfs.createUnitTestFileErr = errors.New("create error")
			},
			expectedError: errors.New("creating file config/config_test.go: create error"),
		},
		{
			name: "error when writing to config reader unit file, template parse error",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor) {
				mfs.createdFile = new(mockFile)
				mtp.err = errors.New("parse error")
			},
			expectedError: errors.New("parsing template configReaderMainFile: parse error"),
		},
		{
			name: "error when creating .env file",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor) {
				mfs.createdFile = new(mockFile)
				mtp.te = new(mockTemplateExecutor)
				mfs.createEnvFileErr = errors.New("create error")
			},
			expectedError: errors.New("creating file .env: create error"),
		},
		{
			name: "error when writing .env file, template parse error",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor) {
				mfs.createdFile = new(mockFile)
				mtp.te = new(mockTemplateExecutor)
				mtp.envFileTemplateParseErr = errors.New("template parse error")
			},
			expectedError: errors.New("parsing template envFile: template parse error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mfs := new(mockFileSystem)
			mtp := new(mockTemplateProcessor)
			tc.mockClosure(mfs, mtp)
			fsProvider = mfs
			templateProcessorProvider = mtp
			g := NewGenerator("config")
			output, err := g.GenerateConfigPackage()
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestGenerateConfigPackageFromEnvFile(t *testing.T) {
	testCases := []struct {
		name           string
		mockClosure    func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mlr *mockLineReader, mfr *mockFormatter)
		expectedOutput []string
		expectedError  error
	}{
		{
			name: "happy path",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mlr *mockLineReader, mfr *mockFormatter) {
				mf := new(mockFile)
				mfs.createdFile = mf
				mfs.openedFile = mf
				mtp.te = new(mockTemplateExecutor)
				mlr.lines = []string{"invalid", "var=value"}
			},
			expectedOutput: []string{
				"config/config.go",
				"config/config_test.go",
			},
		},
		{
			name: "error when creating config files dir",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mlr *mockLineReader, mfr *mockFormatter) {
				mfs.mkDirErr = errors.New("mkdir error")
			},
			expectedError: errors.New("creating dir config: mkdir error"),
		},
		{
			name: "error when generating config struct from env file, open env file error",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mlr *mockLineReader, mfr *mockFormatter) {
				mfs.createdFile = new(mockFile)
				mfs.openErr = errors.New("open error")
			},
			expectedError: errors.New("opening env file .env-local: open error"),
		},
		{
			name: "line reader error",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mlr *mockLineReader, mfr *mockFormatter) {
				mf := new(mockFile)
				mfs.createdFile = mf
				mfs.openedFile = mf
				mtp.te = new(mockTemplateExecutor)
				mlr.lines = []string{"invalid", "var=value"}
				mlr.err = errors.New("error")
			},
			expectedError: errors.New("generating struct from env file .env-local: scanning: error"),
		},
		{
			name: "error when creating config reader unit test file",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mlr *mockLineReader, mfr *mockFormatter) {
				mf := new(mockFile)
				mfs.createdFile = mf
				mfs.openedFile = mf
				mtp.te = new(mockTemplateExecutor)
				mlr.lines = []string{"invalid", "var=value"}
				mfs.createUnitTestFileErr = errors.New("create error")
			},
			expectedError: errors.New("creating file config/config_test.go: create error"),
		},
		{
			name: "error when writing to config reader unit file, template parse error",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mlr *mockLineReader, mfr *mockFormatter) {
				mfs.createdFile = new(mockFile)
				mf := new(mockFile)
				mfs.createdFile = mf
				mfs.openedFile = mf
				mtp.te = new(mockTemplateExecutor)
				mtp.configReaderUnitTestFileParseErr = errors.New("parse error")
			},
			expectedError: errors.New("parsing template configReaderUnitTestFile: parse error"),
		},
		{
			name: "error when formating file",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mlr *mockLineReader, mfr *mockFormatter) {
				mf := new(mockFile)
				mfs.createdFile = mf
				mfs.openedFile = mf
				mtp.te = new(mockTemplateExecutor)
				mlr.lines = []string{"invalid", "var=value"}
				mfr.sourceErr = errors.New("source error")
			},
			expectedError: errors.New("formating go file config/config.go: source error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mfs := new(mockFileSystem)
			mtp := new(mockTemplateProcessor)
			mlr := new(mockLineReader)
			mfr := new(mockFormatter)
			tc.mockClosure(mfs, mtp, mlr, mfr)
			fsProvider = mfs
			templateProcessorProvider = mtp
			formatterProvider = mfr
			lr = func(_ io.Reader) lineReader {
				return mlr
			}
			g := NewGenerator("config")
			output, err := g.GenerateConfigPackageFromEnvFile(".env-local")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}
