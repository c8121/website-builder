package config

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	ProjectDir string

	SpaHttpRoot string // Root directory to service SPA from. Using embedded FS if empty.

	ListenAddress string //Spa-/Rest-/Ssh-Server in format ip:port
	CertFile      string
	KeyFile       string
}

var (
	AppConfig *Config = nil
)

func NewConfigFromFlags() *Config {

	flagSet := flag.NewFlagSet("website-builder", flag.ExitOnError)

	cmdDataDir := flagSet.String("data", "", "Data directory for storage, meta-data, db...")
	cmdSpaHttpRoot := flagSet.String("spa", "", "HTTP root directory of SPA app")
	cmdListen := flagSet.String("listen", "", "Listen Address (ip:port)")
	cmdCertFile := flagSet.String("cert", "", "Certificate file")
	cmdKeyFile := flagSet.String("key", "", "Key file")

	flagSet.Parse(os.Args[1:])

	config := &Config{}

	useProjectDir := *cmdDataDir
	if useProjectDir == "" {
		if userHome, err := os.UserHomeDir(); err != nil {
			panic("failed to get user home directory")
		} else {
			useProjectDir = filepath.Join(userHome, "website-builder")
		}
	}
	config.ProjectDir = useProjectDir

	if *cmdSpaHttpRoot != "" {
		config.SpaHttpRoot = *cmdSpaHttpRoot
		slog.Info("Init Path", "spa", config.SpaHttpRoot)
	}

	if *cmdListen != "" {
		config.ListenAddress = *cmdListen
		slog.Info("Server", "bind", config.ListenAddress)
	} else {
		config.ListenAddress = "127.0.0.1:9998"
	}

	if *cmdCertFile != "" {
		config.CertFile = *cmdCertFile
		slog.Info("Init TLS", "cert", config.CertFile)
	}

	if *cmdKeyFile != "" {
		config.KeyFile = *cmdKeyFile
		slog.Info("Init TLS", "key", config.KeyFile)
	}

	return config
}
