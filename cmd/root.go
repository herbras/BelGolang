// Package cmd provides command-line interface functionality for the salat application
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"jadwalsalat/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "salat",
	Short: "Aplikasi CLI untuk jadwal sholat",
	Long: `Salat adalah aplikasi command line untuk menampilkan jadwal sholat
berdasarkan lokasi dan metode perhitungan yang dikonfigurasi.`,
	// Default command: show prayer times when no subcommand is provided
	Run: func(cmd *cobra.Command, args []string) {
		// Check if config exists
		cfg, err := config.LoadConfig()
		if err != nil || cfg.Latitude == 0 || cfg.Longitude == 0 {
			// If no config or incomplete config, run setup
			fmt.Println("Konfigurasi belum lengkap. Menjalankan setup...")
			setupCmd.Run(cmd, args)
			return
		}

		// Otherwise show prayer times
		showPrayerTimes(false, "light")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/salat/config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".salat" (without extension).
		configDir := filepath.Join(home, ".config", "salat")

		// Ensure config directory exists
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			err = os.MkdirAll(configDir, 0755)
			cobra.CheckErr(err)
		}

		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// Config file found and successfully parsed
	} else {
		// Config file not found or error reading it
		// This is normal on first run, so we don't need to show an error
	}
}
