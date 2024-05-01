package main

import (
	"fmt"
	"os"

	"github.com/tiagomelo/go-project-config/appcfg"
)

func main() {
	cfg, err := appcfg.ReadFromEnvFile(".env-sample")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("cfg: %+v\n", cfg)
}
