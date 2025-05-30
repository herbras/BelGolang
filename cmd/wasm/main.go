package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"jadwalsalat/salat"
)

// SalatAPI represents the main API for WASM
type SalatAPI struct{}

// ProcessPrayerTime processes prayer time calculation and returns JSON
func (api *SalatAPI) ProcessPrayerTime(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return map[string]interface{}{
			"error": "Missing required arguments: latitude, longitude",
		}
	}

	lat := args[0].Float()
	lng := args[1].Float()

	// Optional method parameter (default to Kemenag for Indonesia)
	method := salat.Kemenag
	if len(args) >= 3 && args[2].String() != "" {
		methodStr := args[2].String()
		switch methodStr {
		case "MWL":
			method = salat.MWL
		case "ISNA":
			method = salat.ISNA
		case "Egypt":
			method = salat.Egypt
		case "Makkah":
			method = salat.Makkah
		case "Karachi":
			method = salat.Karachi
		case "Tehran":
			method = salat.Tehran
		case "JAKIM":
			method = salat.JAKIM
		default:
			method = salat.Kemenag
		}
	}

	// Create location
	location := salat.Location{
		Latitude:  lat,
		Longitude: lng,
		Method:    method,
	}

	// Calculate for today
	now := time.Now()
	times, err := salat.TimesForDate(now, location)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to calculate prayer times: %v", err),
		}
	}

	// Get current prayer
	currentPrayer, _ := salat.GetCurrentPrayer(now, times)
	nextPrayerName, nextPrayerTime := salat.GetNextPrayer(now, times)

	// Ensure all values are basic types
	currentPrayerStr := currentPrayer
	if currentPrayerStr == "" {
		currentPrayerStr = "Unknown"
	}

	nextPrayerStr := nextPrayerName
	if nextPrayerStr == "" {
		nextPrayerStr = "Unknown"
	}

	// Build result with only basic types
	result := map[string]interface{}{
		"location": map[string]interface{}{
			"latitude":  lat,
			"longitude": lng,
		},
		"method": string(method),
		"date":   now.Format("2006-01-02"),
		"prayers": map[string]interface{}{
			"imsak":   times.Imsak.Format("15:04"),
			"subuh":   times.Subuh.Format("15:04"),
			"dzuhur":  times.Dzuhur.Format("15:04"),
			"ashar":   times.Ashar.Format("15:04"),
			"maghrib": times.Maghrib.Format("15:04"),
			"isya":    times.Isya.Format("15:04"),
		},
		"current": map[string]interface{}{
			"prayer": currentPrayerStr,
			"emoji":  salat.GetPrayerEmoji(currentPrayerStr),
		},
		"next": map[string]interface{}{
			"prayer": nextPrayerStr,
			"time":   nextPrayerTime.Format("15:04"),
			"emoji":  salat.GetPrayerEmoji(nextPrayerStr),
		},
		"timestamp": now.Format(time.RFC3339),
	}

	return result
}

// GetVersion returns the current version
func (api *SalatAPI) GetVersion(this js.Value, args []js.Value) interface{} {
	return map[string]string{
		"version": "1.6.1",
		"build":   "wasm",
		"runtime": "browser",
		"methods": "MWL,ISNA,Egypt,Makkah,Karachi,Tehran,Kemenag,JAKIM",
	}
}

// ProcessCommand processes CLI-like commands
func (api *SalatAPI) ProcessCommand(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{
			"error": "No command provided",
		}
	}

	command := args[0].String()

	switch command {
	case "help":
		return map[string]interface{}{
			"commands": []string{
				"prayer <lat> <lng> [method] - Calculate prayer times",
				"version - Get version info",
				"methods - List available calculation methods",
				"help - Show this help",
			},
		}
	case "version":
		return api.GetVersion(this, args[1:])
	case "methods":
		return map[string]interface{}{
			"methods": []string{
				"MWL - Muslim World League",
				"ISNA - Islamic Society of North America",
				"Egypt - Egyptian General Authority of Survey",
				"Makkah - Umm al-Qura University, Makkah",
				"Karachi - University of Islamic Sciences, Karachi",
				"Tehran - Institute of Geophysics, University of Tehran",
				"Kemenag - Kementerian Agama Republik Indonesia (default)",
				"JAKIM - Jabatan Kemajuan Islam Malaysia",
			},
		}
	case "prayer":
		if len(args) < 3 {
			return map[string]interface{}{
				"error": "Prayer command requires: lat, lng, [method]",
			}
		}
		// Convert string args to proper types
		lat, err := strconv.ParseFloat(args[1].String(), 64)
		if err != nil {
			return map[string]interface{}{
				"error": "Invalid latitude: " + args[1].String(),
			}
		}
		lng, err := strconv.ParseFloat(args[2].String(), 64)
		if err != nil {
			return map[string]interface{}{
				"error": "Invalid longitude: " + args[2].String(),
			}
		}

		// Create new args with proper types
		newArgs := []js.Value{
			js.ValueOf(lat),
			js.ValueOf(lng),
		}
		if len(args) >= 4 {
			newArgs = append(newArgs, args[3])
		}

		return api.ProcessPrayerTime(this, newArgs)
	default:
		return map[string]interface{}{
			"error": fmt.Sprintf("Unknown command: %s", command),
		}
	}
}

func main() {
	// Keep the program running
	c := make(chan struct{}, 0)

	// Create API instance
	api := &SalatAPI{}

	// Register global functions
	js.Global().Set("salat", js.ValueOf(map[string]interface{}{
		"prayerTime":     js.FuncOf(api.ProcessPrayerTime),
		"version":        js.FuncOf(api.GetVersion),
		"processCommand": js.FuncOf(api.ProcessCommand),
	}))

	// Register console-style API for terminal libraries
	js.Global().Set("salatConsole", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return "Usage: salatConsole('command arg1 arg2')"
		}

		input := args[0].String()
		// Simple command parsing
		parts := strings.Fields(input)
		if len(parts) == 0 {
			return api.ProcessCommand(this, []js.Value{js.ValueOf("help")})
		}

		jsArgs := make([]js.Value, len(parts))
		for i, part := range parts {
			jsArgs[i] = js.ValueOf(part)
		}

		result := api.ProcessCommand(this, jsArgs)
		jsonBytes, _ := json.Marshal(result)
		return string(jsonBytes)
	}))

	fmt.Println("🕌 Salat WASM API ready!")
	fmt.Println("Available functions:")
	fmt.Println("- salat.prayerTime(lat, lng, [method])")
	fmt.Println("- salat.version()")
	fmt.Println("- salat.processCommand(command, ...args)")
	fmt.Println("- salatConsole('command args')")

	<-c
}
