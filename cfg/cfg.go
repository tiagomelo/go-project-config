// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cfg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const (
	configReadFileName           = "config.go"
	configReaderUnitTestFileName = "config_test.go"
	envFileName                  = ".env"
)

// For ease of unit testing.
var (
	// fsProvider is a variable of interface type fileSystem. It abstracts
	// file system operations and allows the use of different file system
	// implementations (like mocks for testing).
	fsProvider fileSystem = osFileSystem{}

	// templateProcessorProvider is a variable of interface type templateProcessor.
	// It abstracts template parsing and execution and allows different implementations.
	templateProcessorProvider templateProcessor = textTemplateProcessor{}

	// lr is a function that returns a lineReader interface. It abstracts the creation
	// of a lineReader, which provides functions for scanning lines, retrieving the text
	// of the current line, and obtaining any encountered errors during the scanning process.
	lr = func(r io.Reader) lineReader {
		return bufio.NewScanner(r)
	}
)

// Generator is an interface for generating configuration files.
type Generator interface {
	// GenerateConfigPackage generates '<packagename>/config.go'
	// and '<packagename>/config_test.go' using a provided sample .env file.
	GenerateConfigPackage() ([]string, error)
	// GenerateConfigPackage generates 'config' package with '<packagename>/config.go'
	// and '<packagename>/config_test.go' using a provided .env file.
	GenerateConfigPackageFromEnvFile(envFilePath string) ([]string, error)
}

// generator struct implements the Generator interface.
type generator struct {
	packageName string
}

// NewGenerator creates a new instance of Generator.
func NewGenerator(packageName string) Generator {
	return &generator{
		packageName: packageName,
	}
}

func (g *generator) GenerateConfigPackage() ([]string, error) {
	generatedFiles, err := g.generateConfigReaderFiles()
	if err != nil {
		return nil, err
	}
	return generatedFiles, nil
}

func (g *generator) GenerateConfigPackageFromEnvFile(envFilePath string) ([]string, error) {
	generatedFiles, err := g.generateConfigReaderFilesFromEnvFile(envFilePath)
	if err != nil {
		return nil, err
	}
	return generatedFiles, nil
}

// generateConfigReaderFilesFromEnvFile generates config reader files from .env file.
func (g *generator) generateConfigReaderFilesFromEnvFile(envFilePath string) ([]string, error) {
	var generatedFiles []string
	if err := fsProvider.Mkdir(g.packageName); err != nil && !os.IsExist(err) {
		return nil, errors.Wrapf(err, "creating dir %s", g.packageName)
	}
	mainFilePath, err := g.generateConfigReaderMainFile(envFilePath)
	if err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, mainFilePath)
	unitTestFilePath, err := g.generateConfigReaderUnitTestFile()
	if err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, unitTestFilePath)
	return generatedFiles, nil
}

// generateConfigReaderFiles generates config reader files.
func (g *generator) generateConfigReaderFiles() ([]string, error) {
	var generatedFiles []string
	if err := fsProvider.Mkdir(g.packageName); err != nil && !os.IsExist(err) {
		return nil, errors.Wrapf(err, "creating dir %s", g.packageName)
	}
	mainFilePath, err := g.generateConfigReaderMainFile("")
	if err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, mainFilePath)
	unitTestFilePath, err := g.generateConfigReaderUnitTestFile()
	if err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, unitTestFilePath)
	if err := g.generateEnvFile(); err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, envFileName)
	return generatedFiles, nil
}

// generateEnvFile generates a sample .env file.
func (g *generator) generateEnvFile() error {
	envFile, err := fsProvider.Create(envFileName)
	if err != nil {
		return errors.Wrapf(err, "creating file %s", envFileName)
	}
	defer envFile.Close()
	if err := writeFileFromTemplate(envFileTemplateName,
		envFileTemplate,
		nil,
		envFile); err != nil {
		return err
	}
	return nil
}

// generateConfigReaderMainFile generates config reader main file.
func (g *generator) generateConfigReaderMainFile(envFilePath string) (string, error) {
	configReaderFilePath := fmt.Sprintf("%s/%s", g.packageName, configReadFileName)
	configReaderFile, err := fsProvider.Create(configReaderFilePath)
	if err != nil {
		return "", errors.Wrapf(err, "creating file %s", configReaderFilePath)
	}
	defer configReaderFile.Close()
	templateValues := map[string]string{configReaderPkgPlaceHolder: g.packageName, configStructTemplateName: defaultConfigStructTemplate}
	if envFilePath != "" {
		structFromEnvFile, err := generateConfigStructFromEnvFile(envFilePath)
		if err != nil {
			return "", err
		}
		templateValues[configStructTemplateName] = structFromEnvFile
	}
	if err := writeFileFromTemplate(configReaderMainFileTemplateName,
		configReaderMainFileTemplatePlaceHolder,
		templateValues,
		configReaderFile); err != nil {
		return "", err
	}
	if err := formatGoFile(configReaderFilePath); err != nil {
		return "", err
	}
	return configReaderFilePath, nil
}

// generateConfigStructFromEnvFile generates the 'Config' struct from
// variables defined in the provided .env file.
func generateConfigStructFromEnvFile(envFilePath string) (string, error) {
	envFile, err := fsProvider.Open(envFilePath)
	if err != nil {
		return "", errors.Wrapf(err, "opening env file %s", envFilePath)
	}
	defer envFile.Close()
	structFromEnvFile, err := generateStructFromEnvFile(lr(envFile))
	if err != nil {
		return "", errors.Wrapf(err, "generating struct from env file %s", envFilePath)
	}
	return structFromEnvFile, nil
}

// generateStructFromEnvFile parses the provided .env file and
// generate the 'Config' struct with the correspondent properties.
func generateStructFromEnvFile(lineReader lineReader) (string, error) {
	var sb strings.Builder
	sb.WriteString("// Config holds all configuration needed by this app.\n")
	sb.WriteString("type Config struct {\n")
	sb.WriteString("// TODO: see https://github.com/kelseyhightower/envconfig for all available options\n // for struct tags.\n")
	for lineReader.Scan() {
		line := lineReader.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // skip invalid lines.
		}
		key := parts[0]
		goFieldName := toCamelCase(key)
		sb.WriteString(fmt.Sprintf("\t%s %s `envconfig:\"%s\" required:\"true\"` // TODO: set the correct data type.\n", goFieldName, "interface{}", key))
	}
	sb.WriteString("}\n")
	if err := lineReader.Err(); err != nil {
		return "", errors.Wrap(err, "scanning")
	}
	return sb.String(), nil
}

// generateConfigReaderUnitTestFile generates unit test file.
func (g *generator) generateConfigReaderUnitTestFile() (string, error) {
	configReaderUnitTestFilePath := fmt.Sprintf("%s/%s", g.packageName, configReaderUnitTestFileName)
	configReaderUnitTestFile, err := fsProvider.Create(configReaderUnitTestFilePath)
	if err != nil {
		return "", errors.Wrapf(err, "creating file %s", configReaderUnitTestFilePath)
	}
	defer configReaderUnitTestFile.Close()
	if err := writeFileFromTemplate(configReaderUnitTestFileTemplateName,
		configReaderUnitTestFileTemplate,
		map[string]string{configReaderPkgPlaceHolder: g.packageName},
		configReaderUnitTestFile); err != nil {
		return "", err
	}
	return configReaderUnitTestFilePath, nil
}

// writeFileFromTemplate parses and then executes the given template with
// the given template values.
func writeFileFromTemplate(templateName, templateText string, templateValues map[string]string, file File) error {
	tmplExecutor, err := templateProcessorProvider.Parse(templateName, templateText)
	if err != nil {
		return errors.Wrapf(err, "parsing template %s", templateName)
	}
	err = tmplExecutor.Execute(file, templateValues)
	if err != nil {
		return errors.Wrapf(err, "executing template %s", templateName)
	}
	return nil
}
