package cmd

import (
	"fmt"
	"strings"
	"time"

	"jadwalsalat/config"
	"jadwalsalat/salat"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"s"},
	Short:   "Tampilkan jadwal sholat hari ini",
	Long:    `Tampilkan jadwal sholat hari ini berdasarkan konfigurasi lokasi dan metode perhitungan.`,
	Run: func(cmd *cobra.Command, args []string) {
		compactMode, _ := cmd.Flags().GetBool("compact")
		theme, _ := cmd.Flags().GetString("theme")
		showPrayerTimes(compactMode, theme)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolP("compact", "c", false, "Tampilkan dalam mode compact")
	showCmd.Flags().StringP("theme", "t", "light", "Pilih tema tampilan (light/dark)")
}

// showPrayerTimes displays the prayer times for today
func showPrayerTimes(compactMode bool, theme string) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		fmt.Println("Run 'salat setup' to configure the application.")
		return
	}

	// Parse timezone
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		fmt.Printf("Error parsing timezone: %v\n", err)
		return
	}

	// Get current time in the configured timezone
	now := time.Now().In(loc)

	// Calculate prayer times for today
	location := salat.Location{
		Latitude:  cfg.Latitude,
		Longitude: cfg.Longitude,
		Method:    salat.CalculationMethod(cfg.Method),
	}

	times, err := salat.TimesForDate(now, location)
	if err != nil {
		fmt.Printf("Error calculating prayer times: %v\n", err)
		return
	}

	// Get current and next prayer time
	currentName, _ := salat.GetCurrentPrayer(now, times)
	nextName, nextTime := salat.GetNextPrayer(now, times)
	remaining := nextTime.Sub(now)

	// Format remaining time
	hours := int(remaining.Hours())
	minutes := int(remaining.Minutes()) % 60
	remainStr := fmt.Sprintf("%dh %02dm", hours, minutes)
	if hours == 0 {
		remainStr = fmt.Sprintf("%dm", minutes)
	}

	// Setup colors based on theme
	var headerFgColor, nextFgColor color.Attribute

	if theme == "dark" {
		headerFgColor = color.FgHiCyan
		nextFgColor = color.FgHiYellow
	} else { // light theme (default)
		headerFgColor = color.FgBlue
		nextFgColor = color.FgYellow
	}

	// Create color objects
	headerColor := color.New(headerFgColor, color.Bold)
	nextColor := color.New(nextFgColor)

	// Print header (hanya sekali)
	fmt.Print("\n")
	headerColor.Printf("🕌 Jadwal Sholat - %s\n", now.Format("Monday, 02 January 2006"))
	fmt.Printf("📍 %s (%.6f, %.6f) • %s\n\n", getLocationNameFromConfig(cfg), cfg.Latitude, cfg.Longitude, cfg.Method)

	// Compact mode - single line output
	if compactMode {
		var parts []string
		prayerTimes := []struct {
			name string
			time time.Time
		}{
			{"Subuh", times.Subuh},
			{"Dzuhur", times.Dzuhur},
			{"Ashar", times.Ashar},
			{"Maghrib", times.Maghrib},
			{"Isya", times.Isya},
		}

		for _, prayer := range prayerTimes {
			if prayer.name == currentName {
				parts = append(parts, fmt.Sprintf("%s ►", prayer.name))
			} else if prayer.time.Before(now) {
				parts = append(parts, fmt.Sprintf("%s ✓", prayer.name))
			} else if prayer.name == nextName {
				parts = append(parts, fmt.Sprintf("%s %s", prayer.name, remainStr))
			} else {
				parts = append(parts, prayer.name)
			}
		}

		fmt.Println(strings.Join(parts, " | "))
		fmt.Println()
		return
	}

	// Buat tabel sederhana menggunakan fmt.Printf
	prayerTimes := []struct {
		name string
		time time.Time
	}{
		{"🌙  Imsak", times.Imsak},
		{"🌅  Subuh", times.Subuh},
		{"☀️  Dzuhur", times.Dzuhur},
		{"🌤️  Ashar", times.Ashar},
		{"🌇  Maghrib", times.Maghrib},
		{"✨  Isya", times.Isya},
	}

	// Print header tabel
	fmt.Printf("%-15s %-8s %-10s\n", "WAKTU", "JAM", "STATUS")
	fmt.Println("-------------------------------")

	// Print baris tabel
	for _, prayer := range prayerTimes {
		prayerName := strings.TrimSpace(strings.Split(prayer.name, " ")[1]) // Remove emoji
		timeStr := prayer.time.Format("15:04")
		var status string

		if prayerName == currentName {
			status = "► AKTIF"
		} else if prayer.time.Before(now) {
			status = "✓"
		} else if prayerName == nextName {
			status = remainStr
		} else {
			status = ""
		}

		fmt.Printf("%-15s %-8s %-10s\n", prayer.name, timeStr, status)
	}

	fmt.Println("-------------------------------")

	// Print next prayer info
	fmt.Println()
	nextEmoji := salat.GetPrayerEmoji(strings.Split(nextName, " ")[0])
	nextColor.Printf("⏰ Sholat berikutnya: %s %s dalam %s\n", nextEmoji, nextName, remainStr)
	fmt.Println()
}

// Helper function to get location name (placeholder - could be enhanced with geocoding)
func getLocationNameFromConfig(cfg *config.Config) string {
	// Use saved location name if available
	if cfg.LocationName != "" {
		return cfg.LocationName
	}

	// Fallback to generic name based on coordinates
	if cfg.Latitude >= -7 && cfg.Latitude <= -6 &&
		cfg.Longitude >= 106 && cfg.Longitude <= 107 {
		return "Jakarta"
	}
	return "Lokasi Anda"
}
