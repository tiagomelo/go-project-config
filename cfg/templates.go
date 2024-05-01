// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cfg

const (
	configReaderPkgPlaceHolder  = "ConfigReaderPkgName"
	configStructTemplateName    = "ConfigStruct"
	defaultConfigStructTemplate = `// Config holds all configuration needed by this app.
type Config struct {
	SampleEnvVar string ` + "`envconfig:\"SAMPLE_ENV_VAR\" required:\"true\"`" + `
}`

	configReaderUnitTestFileTemplateName    = "configReaderUnitTestFile"
	configReaderMainFileTemplateName        = "configReaderMainFile"
	configReaderMainFileTemplatePlaceHolder = `package {{ .ConfigReaderPkgName }}

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

{{ .ConfigStruct }}

// For ease of unit testing.
var (
	godotenvLoad     = godotenv.Load
	envconfigProcess = envconfig.Process
)

// Read reads configuration from environment variables.
// It assumes that an '.env' file is present at current path.
func Read() (*Config, error) {
	if err := godotenvLoad(); err != nil {
		return nil, errors.Wrap(err, "loading env vars from .env file")
	}
	config := new(Config)
	if err := envconfigProcess("", config); err != nil {
		return nil, errors.Wrap(err, "processing env vars")
	}
	return config, nil
}

// ReadFromEnvFile reads configuration from the specified environment file.
func ReadFromEnvFile(envFilePath string) (*Config, error) {
	if err := godotenvLoad(envFilePath); err != nil {
		return nil, errors.Wrapf(err, "loading env vars from %s", envFilePath)
	}
	config := new(Config)
	if err := envconfigProcess("", config); err != nil {
		return nil, errors.Wrap(err, "processing env vars")
	}
	return config, nil
}
`

	configReaderUnitTestFileTemplate = `package {{ .ConfigReaderPkgName }}

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		name                   string
		mockedGodotenvLoad     func(filenames ...string) (err error)
		mockedEnvconfigProcess func(prefix string, spec interface{}) error
		expectedError          error
	}{
		{
			name: "happy path",
			mockedGodotenvLoad: func(filenames ...string) (err error) {
				return nil
			},
			mockedEnvconfigProcess: func(prefix string, spec interface{}) error {
				return nil
			},
		},
		{
			name: "error loading env vars",
			mockedGodotenvLoad: func(filenames ...string) (err error) {
				return errors.New("random error")
			},
			expectedError: errors.New("loading env vars from .env file: random error"),
		},
		{
			name: "error processing env vars",
			mockedGodotenvLoad: func(filenames ...string) (err error) {
				return nil
			},
			mockedEnvconfigProcess: func(prefix string, spec interface{}) error {
				return errors.New("random error")
			},
			expectedError: errors.New("processing env vars: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			godotenvLoad = tc.mockedGodotenvLoad
			envconfigProcess = tc.mockedEnvconfigProcess
			config, err := Read()
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Nil(t, config)
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error, got nil")
				}
				require.NotNil(t, config)
			}
		})
	}
}

func TestReadFromEnvFile(t *testing.T) {
	testCases := []struct {
		name                   string
		mockedGodotenvLoad     func(filenames ...string) (err error)
		mockedEnvconfigProcess func(prefix string, spec interface{}) error
		expectedError          error
	}{
		{
			name: "happy path",
			mockedGodotenvLoad: func(filenames ...string) (err error) {
				return nil
			},
			mockedEnvconfigProcess: func(prefix string, spec interface{}) error {
				return nil
			},
		},
		{
			name: "error loading env vars",
			mockedGodotenvLoad: func(filenames ...string) (err error) {
				return errors.New("random error")
			},
			expectedError: errors.New("loading env vars from path/to/.env: random error"),
		},
		{
			name: "error processing env vars",
			mockedGodotenvLoad: func(filenames ...string) (err error) {
				return nil
			},
			mockedEnvconfigProcess: func(prefix string, spec interface{}) error {
				return errors.New("random error")
			},
			expectedError: errors.New("processing env vars: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			godotenvLoad = tc.mockedGodotenvLoad
			envconfigProcess = tc.mockedEnvconfigProcess
			config, err := ReadFromEnvFile("path/to/.env")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Nil(t, config)
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error, got nil")
				}
				require.NotNil(t, config)
			}
		})
	}
}
`
	envFileTemplateName = "envFile"
	envFileTemplate     = `SAMPLE_ENV_VAR=some value`
)
