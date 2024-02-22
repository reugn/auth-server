package main

import (
	"log/slog"
	"os"

	"github.com/reugn/auth-server/internal/auth"
	"github.com/reugn/auth-server/internal/config"
	"github.com/reugn/auth-server/internal/http"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	version = "0.4.0"
)

func run() int {
	rootCmd := &cobra.Command{
		Short:   "Authentication and authorization service",
		Version: version,
	}

	var configFilePath string
	rootCmd.Flags().StringVarP(&configFilePath, "config", "c", "config.yaml", "configuration file path")

	rootCmd.RunE = func(_ *cobra.Command, _ []string) error {
		// read configuration file
		config, err := readConfiguration(configFilePath)
		if err != nil {
			return err
		}
		// load ssl keys
		keys, err := auth.NewKeys(config.Secret)
		if err != nil {
			return err
		}
		// set default logger
		slogHandler, err := config.Logger.SlogHandler()
		if err != nil {
			return err
		}
		slog.SetDefault(slog.New(slogHandler))
		// start http server
		server, err := http.NewServer(version, keys, config)
		if err != nil {
			return err
		}
		slog.Info("Starting service", "config", config)
		return server.Start()
	}

	err := rootCmd.Execute()
	if err != nil {
		return 1
	}
	return 0
}

func readConfiguration(path string) (*config.Service, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := config.NewServiceDefault()
	err = yaml.Unmarshal(data, config)
	return config, err
}

func main() {
	// start the service
	os.Exit(run())
}
