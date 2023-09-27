package main

import (
	"frps-panel/pkg/server"
	"frps-panel/pkg/server/controller"
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const version = "1.0.0"

var (
	showVersion bool
	configFile  string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of frpc-panel")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "./frpc-panel.toml", "config file of frpc-panel")
}

var rootCmd = &cobra.Command{
	Use:   "frpc-panel",
	Short: "frpc-panel is the server plugin of frp to support multiple users.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			log.Println(version)
			return nil
		}
		executable, err := os.Executable()
		if err != nil {
			log.Printf("error get program path: %v", err)
			return err
		}
		rootDir := filepath.Dir(executable)

		configDir := filepath.Dir(configFile)
		tokensFile := filepath.Join(configDir, "frps-tokens.toml")

		config, tls, err := parseConfigFile(configFile, tokensFile)
		if err != nil {
			log.Printf("fail to start frpc-panel : %v", err)
			return err
		}

		s, err := server.New(
			rootDir,
			config,
			tls,
		)
		if err != nil {
			return err
		}
		err = s.Run()
		if err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func parseConfigFile(configFile, tokensFile string) (controller.HandleController, server.TLS, error) {
	var common controller.Common
	_, err := toml.DecodeFile(configFile, &common)
	if err != nil {
		log.Fatalf("decode config file %v error: %v", configFile, err)
	}

	common.Common.DashboardTls = strings.HasPrefix(strings.ToLower(common.Common.DashboardAddr), "https://")

	tls := server.TLS{
		Enable:   common.Common.TlsMode,
		Protocol: "HTTP",
		Cert:     common.Common.TlsCertFile,
		Key:      common.Common.TlsKeyFile,
	}

	if tls.Enable {
		tls.Protocol = "HTTPS"

		if strings.TrimSpace(tls.Cert) == "" || strings.TrimSpace(tls.Key) == "" {
			tls.Enable = false
			tls.Protocol = "HTTP"
			log.Printf("fail to enable tls: tls cert or key not exist, use http as default.")
		}
	}

	return controller.HandleController{
		CommonInfo: common.Common,
		Version:    version,
		ConfigFile: configFile,
		TokensFile: tokensFile,
	}, tls, nil
}
