package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"jadwalsalat/config"
	"jadwalsalat/salat"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup [lokasi]",
	Short: "Setup konfigurasi aplikasi",
	Long: `Setup interaktif untuk konfigurasi aplikasi salat.

Lokasi bisa berupa:
  - Koordinat: "-6.2,106.8" 
  - Alamat: "Monas Jakarta" atau "Jakarta, Indonesia"

Jika tidak ada parameter lokasi, akan menggunakan mode interaktif.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return setupWithLocation(cmd, args[0])
		}
		return setupInteractive()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringP("api", "a", "nominatim", "API geocoding yang digunakan (nominatim/photon)")
}

// setupWithLocation performs setup with a location argument
func setupWithLocation(cmd *cobra.Command, locInput string) error {
	// Get API type from flag
	apiType, err := cmd.Flags().GetString("api")
	if err != nil {
		apiType = "nominatim" // fallback to default
	}

	fmt.Printf("🌍 Mengatur lokasi: %s\n", locInput)

	var lat, lon float64
	var locationName string

	// 1. Cek apakah input matching "lat,lon"
	parts := strings.Split(locInput, ",")
	var isCoordinates bool
	if len(parts) == 2 {
		// Check if both parts can be parsed as float
		_, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		_, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		isCoordinates = (err1 == nil && err2 == nil)
	}

	if isCoordinates {
		lat, err = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err == nil {
			lon, err = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		}
		if err != nil {
			return fmt.Errorf("koordinat tidak valid: %v", err)
		}

		fmt.Printf("📍 Koordinat: %.6f, %.6f\n", lat, lon)

		// Dapatkan nama lokasi dari reverse geocoding
		fmt.Println("🔍 Mencari nama lokasi...")
		locationName, err = config.ReverseGeocode(apiType, lat, lon)
		if err != nil {
			fmt.Printf("Warning: Tidak dapat menentukan nama lokasi: %v\n", err)
			locationName = fmt.Sprintf("Lokasi (%.6f, %.6f)", lat, lon)
		} else {
			fmt.Printf("📍 Lokasi ditemukan: %s\n", locationName)
		}
	} else {
		// 2. Forward‐geocode jika bukan koordinat
		fmt.Println("🔍 Mencari koordinat lokasi...")
		lat, lon, locationName, err = config.ForwardGeocode(apiType, locInput)
		if err != nil {
			return fmt.Errorf("error geocoding: %v", err)
		}
		fmt.Printf("📍 Lokasi: %s\n", locationName)
		fmt.Printf("📍 Koordinat: %.6f, %.6f\n", lat, lon)
	}

	// Setup timezone dan method dengan default values
	loc, err := config.DetectTimezone()
	if err != nil {
		return fmt.Errorf("error mendeteksi timezone: %v", err)
	}
	timezone := loc.String()

	// Ask untuk override timezone jika perlu
	overrideTz := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Timezone terdeteksi: %s. Ubah?", timezone),
		Default: false,
	}
	if err := survey.AskOne(prompt, &overrideTz); err != nil {
		return fmt.Errorf("error reading timezone confirmation: %v", err)
	}

	if overrideTz {
		tzPrompt := &survey.Input{
			Message: "Masukkan timezone (contoh: Asia/Jakarta):",
			Default: "Asia/Jakarta",
		}
		if err := survey.AskOne(tzPrompt, &timezone); err != nil {
			return fmt.Errorf("error reading timezone input: %v", err)
		}
	}

	// Ask untuk method
	methods := []string{
		string(salat.MWL),
		string(salat.ISNA),
		string(salat.Egypt),
		string(salat.Makkah),
		string(salat.Karachi),
		string(salat.Tehran),
		string(salat.Kemenag),
		string(salat.JAKIM),
	}

	method := string(salat.MWL)
	methodPrompt := &survey.Select{
		Message: "Pilih metode perhitungan waktu sholat:",
		Options: methods,
		Default: method,
	}
	if err := survey.AskOne(methodPrompt, &method); err != nil {
		return fmt.Errorf("error reading method selection: %v", err)
	}

	// Simpan konfigurasi
	cfg := &config.Config{
		Timezone:     timezone,
		Latitude:     lat,
		Longitude:    lon,
		Method:       method,
		LocationName: locationName,
		GeocodingAPI: apiType,
	}

	err = config.SaveConfig(cfg)
	if err != nil {
		return fmt.Errorf("error menyimpan konfigurasi: %v", err)
	}

	fmt.Println("✅ Konfigurasi berhasil disimpan!")
	fmt.Println("Jalankan 'salat show' untuk melihat jadwal sholat hari ini.")
	return nil
}

// setupInteractive runs the interactive setup process
func setupInteractive() error {
	fmt.Println("Selamat datang di setup salat CLI!")

	// Detect timezone
	loc, err := config.DetectTimezone()
	if err != nil {
		fmt.Printf("Error mendeteksi timezone: %v\n", err)
		return err
	}

	fmt.Printf("Timezone terdeteksi: %s\n", loc.String())

	// Ask if user wants to override timezone
	overrideTz := false
	prompt := &survey.Confirm{
		Message: "Apakah Anda ingin mengubah timezone?",
		Default: false,
	}
	if err := survey.AskOne(prompt, &overrideTz); err != nil {
		return fmt.Errorf("error reading timezone confirmation: %v", err)
	}

	timezone := loc.String()
	if overrideTz {
		tzPrompt := &survey.Input{
			Message: "Masukkan timezone (contoh: Asia/Jakarta):",
			Default: "Asia/Jakarta",
		}
		if err := survey.AskOne(tzPrompt, &timezone); err != nil {
			return fmt.Errorf("error reading timezone input: %v", err)
		}
	}

	// Ask for latitude and longitude
	latStr := "-6.2"
	lonStr := "106.8"

	latPrompt := &survey.Input{
		Message: "Masukkan latitude lokasi Anda:",
		Default: latStr,
	}
	if err := survey.AskOne(latPrompt, &latStr); err != nil {
		return fmt.Errorf("error reading latitude input: %v", err)
	}

	lonPrompt := &survey.Input{
		Message: "Masukkan longitude lokasi Anda:",
		Default: lonStr,
	}
	if err := survey.AskOne(lonPrompt, &lonStr); err != nil {
		return fmt.Errorf("error reading longitude input: %v", err)
	}

	// Parse latitude and longitude
	lat, err := strconv.ParseFloat(strings.TrimSpace(latStr), 64)
	if err != nil {
		fmt.Printf("Error parsing latitude: %v\n", err)
		return err
	}

	lon, err := strconv.ParseFloat(strings.TrimSpace(lonStr), 64)
	if err != nil {
		fmt.Printf("Error parsing longitude: %v\n", err)
		return err
	}

	// Ask for calculation method
	methods := []string{
		string(salat.MWL),
		string(salat.ISNA),
		string(salat.Egypt),
		string(salat.Makkah),
		string(salat.Karachi),
		string(salat.Tehran),
		string(salat.Kemenag),
		string(salat.JAKIM),
	}

	method := string(salat.MWL)
	methodPrompt := &survey.Select{
		Message: "Pilih metode perhitungan waktu sholat:",
		Options: methods,
		Default: method,
	}
	if err := survey.AskOne(methodPrompt, &method); err != nil {
		return fmt.Errorf("error reading method selection: %v", err)
	}

	// Try to get location name via reverse geocoding
	locationName := ""
	apiType := "nominatim"
	fmt.Println("🔍 Mencari nama lokasi...")
	if name, err := config.ReverseGeocode(apiType, lat, lon); err == nil {
		locationName = name
		fmt.Printf("📍 Lokasi ditemukan: %s\n", locationName)
	} else {
		fmt.Printf("Warning: Tidak dapat menentukan nama lokasi: %v\n", err)
		locationName = fmt.Sprintf("Lokasi (%.6f, %.6f)", lat, lon)
	}

	// Save configuration
	cfg := &config.Config{
		Timezone:     timezone,
		Latitude:     lat,
		Longitude:    lon,
		Method:       method,
		LocationName: locationName,
		GeocodingAPI: apiType,
	}

	err = config.SaveConfig(cfg)
	if err != nil {
		fmt.Printf("Error menyimpan konfigurasi: %v\n", err)
		return err
	}

	fmt.Println("Konfigurasi berhasil disimpan!")
	fmt.Println("Jalankan 'salat show' untuk melihat jadwal sholat hari ini.")
	return nil
}
