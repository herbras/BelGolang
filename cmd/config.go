package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"jadwalsalat/config"
	"jadwalsalat/salat"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Kelola konfigurasi aplikasi",
	Long:  `Kelola konfigurasi aplikasi salat.`,
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Tampilkan konfigurasi saat ini",
	Long:  `Tampilkan konfigurasi aplikasi salat saat ini.`,
	Run: func(cmd *cobra.Command, args []string) {
		showConfig()
	},
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Atur nilai konfigurasi",
	Long: `Atur nilai konfigurasi aplikasi salat.

Contoh penggunaan:
  salat config set timezone Asia/Jakarta
  salat config set location "Jakarta, Indonesia"
  salat config set location "-6.2,106.8"
  salat config set latitude -6.2
  salat config set longitude 106.8
  salat config set method MWL
  salat config set geocoding_api photon`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		setConfig(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}

// showConfig displays the current configuration
func showConfig() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		fmt.Println("Run 'salat setup' to configure the application.")
		return
	}

	// Print configuration
	fmt.Println("Konfigurasi saat ini:")
	fmt.Printf("  timezone: %s\n", cfg.Timezone)
	if cfg.LocationName != "" {
		fmt.Printf("  lokasi: %s\n", cfg.LocationName)
	}
	fmt.Printf("  latitude: %.6f\n", cfg.Latitude)
	fmt.Printf("  longitude: %.6f\n", cfg.Longitude)
	fmt.Printf("  method: %s\n", cfg.Method)
	if cfg.GeocodingAPI != "" {
		fmt.Printf("  geocoding_api: %s\n", cfg.GeocodingAPI)
	}
}

// setConfig sets a configuration value
func setConfig(key, value string) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		fmt.Println("Run 'salat setup' to configure the application.")
		return
	}

	// Update configuration based on key
	switch strings.ToLower(key) {
	case "timezone":
		cfg.Timezone = value
		fmt.Printf("Timezone diatur ke: %s\n", value)

	case "location", "lokasi":
		// Handle location setting with geocoding
		apiType := cfg.GeocodingAPI
		if apiType == "" {
			apiType = "nominatim"
		}

		// Check if value is coordinates
		parts := strings.Split(value, ",")
		var isCoordinates bool
		if len(parts) == 2 {
			// Check if both parts can be parsed as float
			_, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			_, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			isCoordinates = (err1 == nil && err2 == nil)
		}

		if isCoordinates {
			lat, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			lon, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				cfg.Latitude = lat
				cfg.Longitude = lon

				// Try to get location name
				if name, err := config.ReverseGeocode(apiType, lat, lon); err == nil {
					cfg.LocationName = name
					fmt.Printf("Lokasi diatur ke: %s (%.6f, %.6f)\n", name, lat, lon)
				} else {
					cfg.LocationName = fmt.Sprintf("Lokasi (%.6f, %.6f)", lat, lon)
					fmt.Printf("Koordinat diatur ke: %.6f, %.6f\n", lat, lon)
				}
			} else {
				fmt.Printf("Error: koordinat tidak valid: %v, %v\n", err1, err2)
				return
			}
		} else {
			// Forward geocode
			fmt.Println("🔍 Mencari koordinat lokasi...")
			lat, lon, locationName, err := config.ForwardGeocode(apiType, value)
			if err != nil {
				fmt.Printf("Error: tidak dapat menemukan lokasi: %v\n", err)
				return
			}
			cfg.Latitude = lat
			cfg.Longitude = lon
			cfg.LocationName = locationName
			fmt.Printf("Lokasi diatur ke: %s (%.6f, %.6f)\n", locationName, lat, lon)
		}

	case "latitude":
		lat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Printf("Error: nilai latitude tidak valid: %v\n", err)
			return
		}
		cfg.Latitude = lat
		fmt.Printf("Latitude diatur ke: %.6f\n", lat)

		// Update location name if we have both coordinates
		if cfg.Longitude != 0 {
			apiType := cfg.GeocodingAPI
			if apiType == "" {
				apiType = "nominatim"
			}
			if name, err := config.ReverseGeocode(apiType, lat, cfg.Longitude); err == nil {
				cfg.LocationName = name
				fmt.Printf("Lokasi diperbarui: %s\n", name)
			}
		}

	case "longitude":
		lon, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Printf("Error: nilai longitude tidak valid: %v\n", err)
			return
		}
		cfg.Longitude = lon
		fmt.Printf("Longitude diatur ke: %.6f\n", lon)

		// Update location name if we have both coordinates
		if cfg.Latitude != 0 {
			apiType := cfg.GeocodingAPI
			if apiType == "" {
				apiType = "nominatim"
			}
			if name, err := config.ReverseGeocode(apiType, cfg.Latitude, lon); err == nil {
				cfg.LocationName = name
				fmt.Printf("Lokasi diperbarui: %s\n", name)
			}
		}

	case "method":
		// Validate method
		validMethods := map[string]bool{
			string(salat.MWL):     true,
			string(salat.ISNA):    true,
			string(salat.Egypt):   true,
			string(salat.Makkah):  true,
			string(salat.Karachi): true,
			string(salat.Tehran):  true,
			string(salat.Kemenag): true,
			string(salat.JAKIM):   true,
		}

		if !validMethods[value] {
			fmt.Printf("Error: metode tidak valid. Pilih salah satu dari: MWL, ISNA, Egypt, Makkah, Karachi, Tehran, Kemenag, JAKIM\n")
			return
		}

		cfg.Method = value
		fmt.Printf("Metode perhitungan diatur ke: %s\n", value)

	case "geocoding_api":
		if value != "nominatim" && value != "photon" {
			fmt.Printf("Error: API tidak valid. Pilih salah satu dari: nominatim, photon\n")
			return
		}
		cfg.GeocodingAPI = value
		fmt.Printf("Geocoding API diatur ke: %s\n", value)

	default:
		fmt.Printf("Error: kunci konfigurasi tidak valid. Pilih salah satu dari: timezone, location, latitude, longitude, method, geocoding_api\n")
		return
	}

	// Save configuration
	err = config.SaveConfig(cfg)
	if err != nil {
		fmt.Printf("Error menyimpan konfigurasi: %v\n", err)
		return
	}

	fmt.Println("Konfigurasi berhasil disimpan!")
}
