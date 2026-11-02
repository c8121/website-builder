package config

import (
	"flag"
	"fmt"
)

var (
	SpaHttpRoot = "" // Root directory to service SPA from. Using embedded FS if empty.

	ListenAddress = "127.0.0.1:9998" //Spa-/Rest-/Ssh-Server in format ip:port
	CertFile      = ""
	KeyFile       = ""
)

// LoadDefault initializes configuration with defaults
func LoadDefault() {

	cmdSpaHttpRoot := flag.String("spa", "", "HTTP root directory of SPA app")
	cmdCertFile := flag.String("cert", "", "Certificate file")
	cmdKeyFile := flag.String("key", "", "Key file")

	flag.Parse()

	if *cmdSpaHttpRoot != "" {
		SpaHttpRoot = *cmdSpaHttpRoot
		fmt.Printf("Using SPA directory: %s\n", SpaHttpRoot)
	}

	if *cmdCertFile != "" {
		CertFile = *cmdCertFile
		fmt.Printf("Certificate file: %s\n", CertFile)
	}

	if *cmdKeyFile != "" {
		KeyFile = *cmdKeyFile
		fmt.Printf("Key file: %s\n", KeyFile)
	}
}
