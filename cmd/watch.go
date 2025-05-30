package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"jadwalsalat/config"
	"jadwalsalat/salat"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:     "watch",
	Aliases: []string{"w"},
	Short:   "Tampilkan jadwal sholat secara live",
	Long:    `Tampilkan jadwal sholat secara live dengan update setiap menit.`,
	Run: func(cmd *cobra.Command, args []string) {
		notify, _ := cmd.Flags().GetBool("notify")
		watchPrayerTimes(notify)
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().BoolP("notify", "n", false, "Aktifkan notifikasi saat waktu sholat tiba")
}

// watchPrayerTimes displays prayer times in a live updating view
func watchPrayerTimes(notify bool) {
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

	// Setup signal handling for graceful exit
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Setup colors
	headerColor := color.New(color.FgHiCyan, color.Bold)
	activeColor := color.New(color.FgHiYellow, color.Bold)
	nextColor := color.New(color.FgHiGreen)
	normalColor := color.New(color.FgWhite)
	timeColor := color.New(color.FgHiWhite)

	// Create ticker for updates
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Track the current prayer to detect changes
	var lastCurrentPrayer string

	// Main loop
	for {
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
		seconds := int(remaining.Seconds()) % 60
		remainStr := fmt.Sprintf("%dh %02dm %02ds", hours, minutes, seconds)
		if hours == 0 {
			remainStr = fmt.Sprintf("%dm %02ds", minutes, seconds)
		}

		// Clear screen
		fmt.Print("\033[H\033[2J")

		// Print header
		headerColor.Printf("🕌 JADWAL SHOLAT LIVE\n")
		fmt.Printf("📍 %s (%.6f, %.6f) • %s\n", getLocationNameFromConfig(cfg), cfg.Latitude, cfg.Longitude, cfg.Method)
		fmt.Printf("⏰ %s\n\n", now.Format("Monday, 02 January 2006 15:04:05"))

		// Print prayer times
		fmt.Println("Waktu Sholat Hari Ini:")
		fmt.Println("---------------------")

		// Check if the current prayer has changed
		prayerChanged := lastCurrentPrayer != currentName && lastCurrentPrayer != ""
		lastCurrentPrayer = currentName

		// Send notification if enabled and prayer time has changed
		if notify && prayerChanged {
			// Print a box around the notification
			notifyColor := color.New(color.FgHiWhite, color.BgHiRed)
			fmt.Println()
			notifyColor.Println("┌─────────────────────────────────────┐")
			notifyColor.Printf("│ 🔔 WAKTU SHOLAT %s TELAH TIBA! │\n", strings.ToUpper(currentName))
			notifyColor.Println("└─────────────────────────────────────┘")
			fmt.Println()

			// Sound the terminal bell three times
			fmt.Print("\a\a\a")
		}

		// Display prayer times with status
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
			emoji := salat.GetPrayerEmoji(prayer.name)
			timeStr := prayer.time.Format("15:04")

			if prayer.name == currentName {
				activeColor.Printf("%s %s: %s ► AKTIF\n", emoji, prayer.name, timeStr)
			} else if prayer.time.Before(now) {
				normalColor.Printf("%s %s: %s ✓\n", emoji, prayer.name, timeStr)
			} else if prayer.name == nextName {
				nextColor.Printf("%s %s: %s ⏰ %s\n", emoji, prayer.name, timeStr, remainStr)
			} else {
				timeColor.Printf("%s %s: %s\n", emoji, prayer.name, timeStr)
			}
		}

		fmt.Println("\n---------------------")

		// Display next prayer countdown
		nextEmoji := salat.GetPrayerEmoji(strings.Split(nextName, " ")[0]) // Ambil nama sholat tanpa "(besok)"
		nextColor.Printf("⏱️ Sholat berikutnya: %s %s dalam %s\n", nextEmoji, nextName, remainStr)

		// Create ASCII progress bar
		progressBarWidth := 40
		var intervalDuration time.Duration

		// Determine the interval (time between previous prayer and next prayer)
		prayerTimesWithNext := []struct {
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
		for i, prayer := range prayerTimesWithNext {
			if prayer.name == nextName {
				if i == 0 {
					// If it's the first prayer of the day, use midnight as the previous time
					prevPrayer = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
				} else {
					prevPrayer = prayerTimesWithNext[i-1].time
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
		percentColor := color.New(color.FgHiCyan)
		percentColor.Printf("%.1f%% selesai\n", progress*100)

		// Print progress bar
		progressColor := color.New(color.FgHiGreen)
		emptyColor := color.New(color.FgHiBlack)

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

		// Wait for ticker or signal
		select {
		case <-ticker.C:
			// Continue to next iteration
		case <-sigCh:
			fmt.Println("\nExiting watch mode...")
			return
		}
	}
}
