package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-server-template/internal/conf"
	"log"
	"os"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "base-cmd",
	Short: "base command",
	Long: "base command",
}

var cfgFile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	configPath string
)

func init() {
	cobra.OnInitialize(initConfig)

	//RootCmd.PersistentFlags().StringVar(&configPath, "config_path", "", "config path (default is $PROJECT_HOME/configs/config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("configs")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			cfg := conf.InitDefaultConfig()
			saveToViper(cfg)

			_ = viper.SafeWriteConfig()
		} else {
			// Config file was found but another error was produced
		}
	}

	// Config file found and successfully parsed
	if err := viper.Unmarshal(&conf.Conf); err != nil {
		log.Fatalf("failed to umarshal to conf.Conf")
	}
}

func saveToViper(cfg *conf.Config) {
	viper.Set("database", cfg.Database)
	viper.Set("logger", cfg.Logger)
}