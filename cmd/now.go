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

// nowCmd represents the now command
var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "Tampilkan waktu sholat saat ini dan countdown",
	Long:  `Tampilkan waktu sholat saat ini dan countdown ke waktu sholat berikutnya dengan tampilan ringkas.`,
	Run: func(cmd *cobra.Command, args []string) {
		showCurrentPrayer()
	},
}

func init() {
	rootCmd.AddCommand(nowCmd)
}

// showCurrentPrayer displays the current prayer time and countdown to next prayer
func showCurrentPrayer() {
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
	currentName, hasActive := salat.GetCurrentPrayer(now, times)
	nextName, nextTime := salat.GetNextPrayer(now, times)
	remaining := nextTime.Sub(now)

	// Format remaining time
	hours := int(remaining.Hours())
	minutes := int(remaining.Minutes()) % 60
	seconds := int(remaining.Seconds()) % 60
	remainStr := fmt.Sprintf("%dh %02dm %02ds", hours, minutes, seconds)
	if hours == 0 {
		remainStr = fmt.Sprintf("%dm %02ds", minutes, seconds)
	}

	// Create color objects
	headerColor := color.New(color.FgHiCyan, color.Bold)
	activeColor := color.New(color.FgHiGreen, color.Bold)
	nextColor := color.New(color.FgHiYellow, color.Bold)
	timeColor := color.New(color.FgHiWhite)

	// Print header with current date and time
	headerColor.Printf("\n🕌 Jadwal Sholat - %s\n", now.Format("15:04:05"))
	fmt.Printf("📍 %s • %s\n\n", getLocationNameFromConfig(cfg), now.Format("Monday, 2 January 2006"))

	// Print current prayer info
	if hasActive {
		currentEmoji := salat.GetPrayerEmoji(currentName)
		activeColor.Printf("Waktu sholat saat ini: %s %s\n", currentEmoji, currentName)

		// Find the time for current prayer
		var currentTime time.Time
		switch currentName {
		case "Subuh":
			currentTime = times.Subuh
		case "Dzuhur":
			currentTime = times.Dzuhur
		case "Ashar":
			currentTime = times.Ashar
		case "Maghrib":
			currentTime = times.Maghrib
		case "Isya":
			currentTime = times.Isya
		}

		// Calculate elapsed time since current prayer started
		elapsed := now.Sub(currentTime)
		elapsedHours := int(elapsed.Hours())
		elapsedMinutes := int(elapsed.Minutes()) % 60
		elapsedSeconds := int(elapsed.Seconds()) % 60
		elapsedStr := fmt.Sprintf("%dh %02dm %02ds", elapsedHours, elapsedMinutes, elapsedSeconds)
		if elapsedHours == 0 {
			elapsedStr = fmt.Sprintf("%dm %02ds", elapsedMinutes, elapsedSeconds)
		}

		timeColor.Printf("Dimulai: %s (%s yang lalu)\n", currentTime.Format("15:04"), elapsedStr)
	} else {
		activeColor.Printf("Tidak ada waktu sholat saat ini\n")
	}

	// Print next prayer info
	fmt.Println()
	nextEmoji := salat.GetPrayerEmoji(strings.Split(nextName, " ")[0]) // Get prayer name without "(besok)"
	nextColor.Printf("Sholat berikutnya: %s %s\n", nextEmoji, nextName)
	timeColor.Printf("Waktu: %s (dalam %s)\n", nextTime.Format("15:04"), remainStr)

	// Create simple progress bar
	progressBarWidth := 30
	var intervalDuration time.Duration

	// Determine the interval (time between previous prayer and next prayer)
	prayerTimes := []struct {
		name string
		time time.Time
	}{
		{"Subuh", times.Subuh},
		{"Dzuhur", times.Dzuhur},
		{"Ashar", times.Ashar},
		{"Maghrib", times.Maghrib},
		{"Isya", times.Isya},
		{"Subuh (besok)", times.Subuh.Add(24 * time.Hour)},
	}

	var prevPrayer time.Time
	for i, prayer := range prayerTimes {
		if prayer.name == nextName {
			if i == 0 {
				// If it's the first prayer of the day, use midnight as the previous time
				prevPrayer = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
			} else {
				prevPrayer = prayerTimes[i-1].time
			}
			break
		}
	}

	intervalDuration = nextTime.Sub(prevPrayer)
	elapsed := now.Sub(prevPrayer)

	var progress float64
	if intervalDuration > 0 {
		progress = float64(elapsed) / float64(intervalDuration)
	} else {
		// Avoid divide by zero if two consecutive prayer times match
		progress = 0
	}

	if progress < 0 {
		progress = 0
	} else if progress > 1 {
		progress = 1
	}

	filledWidth := int(float64(progressBarWidth) * progress)
	emptyWidth := progressBarWidth - filledWidth

	// Print progress bar
	fmt.Println()
	progressColor := color.New(color.FgHiGreen)
	emptyColor := color.New(color.FgHiBlack)

	// Print progress bar with percentage
	fmt.Printf("%3.0f%% ", progress*100)
	progressColor.Print(strings.Repeat("█", filledWidth))
	emptyColor.Print(strings.Repeat("░", emptyWidth))
	fmt.Println()
	fmt.Println()
}
