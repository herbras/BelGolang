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

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:     "next",
	Aliases: []string{"n"},
	Short:   "Tampilkan waktu sholat berikutnya",
	Long:    `Tampilkan waktu sholat berikutnya dan hitung mundur.`,
	Run: func(cmd *cobra.Command, args []string) {
		showNextPrayer()
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
}

// showNextPrayer displays the next prayer time and countdown
func showNextPrayer() {
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

	// Print header
	headerColor := color.New(color.FgHiCyan, color.Bold)
	headerColor.Printf("\n🕌 Jadwal Sholat\n")
	fmt.Printf("📍 %s (%.6f, %.6f) • %s\n\n", getLocationNameFromConfig(cfg), cfg.Latitude, cfg.Longitude, cfg.Method)

	// Print current prayer info if active
	if hasActive {
		activeColor := color.New(color.FgHiYellow, color.Bold)
		activeColor.Printf("Waktu sholat saat ini: %s %s\n", salat.GetPrayerEmoji(currentName), currentName)
	}

	// Print next prayer info
	fmt.Println()
	nextColor := color.New(color.FgHiGreen, color.Bold)
	nextEmoji := salat.GetPrayerEmoji(strings.Split(nextName, " ")[0]) // Ambil nama sholat tanpa "(besok)"
	nextColor.Printf("Sholat berikutnya: %s %s\n", nextEmoji, nextName)
	timeColor := color.New(color.FgHiYellow)
	timeColor.Printf("Waktu: %s\n", nextTime.Format("15:04"))
	countdownColor := color.New(color.FgHiCyan)
	countdownColor.Printf("Countdown: %s\n", remainStr)

	// Create ASCII progress bar
	progressBarWidth := 40
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

	// Print enhanced progress bar with percentage
	fmt.Println()
	progressColor := color.New(color.FgHiGreen)
	emptyColor := color.New(color.FgHiBlack)

	// Print percentage first
	percentColor := color.New(color.FgHiCyan)
	percentColor.Printf("%.1f%% selesai\n", progress*100)

	// Print progress bar
	fmt.Print("[")
	progressColor.Print(strings.Repeat("█", filledWidth))
	emptyColor.Print(strings.Repeat("░", emptyWidth))
	fmt.Print("]")
	fmt.Println()

	// Print time markers
	prevTimeStr := prevPrayer.Format("15:04")
	nextTimeStr := nextTime.Format("15:04")
	timeMarkerFmt := fmt.Sprintf("%%-%ds%%s\n", progressBarWidth+1)
	fmt.Printf(timeMarkerFmt, prevTimeStr, nextTimeStr)
	fmt.Println()
}
