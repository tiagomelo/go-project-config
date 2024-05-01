// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/tiagomelo/go-project-config/cfg"
)

type options struct {
	ConfigPackageName string `short:"p" long:"packageName" description:"package name" required:"true"`
	EnvFile           string `short:"e" long:"envFile" description:"env file" default:""`
}

func run(opts *options) ([]string, error) {
	generator := cfg.NewGenerator(opts.ConfigPackageName)
	if opts.EnvFile != "" {
		return generator.GenerateConfigPackageFromEnvFile(opts.EnvFile)
	}
	return generator.GenerateConfigPackage()
}

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				fmt.Println(err)
				os.Exit(0)
			}
			fmt.Println(err)
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
	generatedFiles, err := run(&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, f := range generatedFiles {
		fmt.Println("created:", f)
	}
}
