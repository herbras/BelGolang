package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Timezone     string  `mapstructure:"timezone"`
	Latitude     float64 `mapstructure:"latitude"`
	Longitude    float64 `mapstructure:"longitude"`
	Method       string  `mapstructure:"method"`
	LocationName string  `mapstructure:"location_name"`
	GeocodingAPI string  `mapstructure:"geocoding_api"`
}

// GetConfigDir returns the directory where config is stored
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, ".config", "salat")
	return configDir, nil
}

// LoadConfig reads the configuration from disk
func LoadConfig() (*Config, error) {
	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %v", err)
	}

	return &config, nil
}

// SaveConfig writes the configuration to disk
func SaveConfig(config *Config) error {
	viper.Set("timezone", config.Timezone)
	viper.Set("latitude", config.Latitude)
	viper.Set("longitude", config.Longitude)
	viper.Set("method", config.Method)
	viper.Set("location_name", config.LocationName)
	viper.Set("geocoding_api", config.GeocodingAPI)

	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.yaml")
	return viper.WriteConfigAs(configPath)
}

// DetectTimezone attempts to detect the system timezone
func DetectTimezone() (*time.Location, error) {
	// Get local timezone
	return time.LoadLocation("")
}

// ForwardGeocode converts an address to coordinates using the specified API
func ForwardGeocode(apiType, query string) (float64, float64, string, error) {
	var apiURL string
	switch apiType {
	case "photon":
		apiURL = fmt.Sprintf("https://photon.komoot.io/api/?q=%s&limit=1", url.QueryEscape(query))
	default: // nominatim
		apiURL = fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1",
			url.QueryEscape(query))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return 0, 0, "", err
	}

	if apiType != "photon" {
		req.Header.Set("User-Agent", "jadwalsalat/1.0")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, "", err
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return 0, 0, "", fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	if apiType == "photon" {
		var r struct {
			Features []struct {
				Geometry struct {
					Coordinates []float64 `json:"coordinates"`
				} `json:"geometry"`
				Properties struct {
					Name        string `json:"name"`
					City        string `json:"city"`
					State       string `json:"state"`
					Country     string `json:"country"`
					DisplayName string `json:"display_name"`
				} `json:"properties"`
			} `json:"features"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return 0, 0, "", err
		}
		if len(r.Features) > 0 {
			// Photon: [lon, lat]
			lat := r.Features[0].Geometry.Coordinates[1]
			lon := r.Features[0].Geometry.Coordinates[0]
			props := r.Features[0].Properties

			// Build location name from available properties
			locationName := props.Name
			if props.City != "" {
				locationName = props.City
			}
			if props.State != "" && locationName != "" {
				locationName += ", " + props.State
			}
			if props.Country != "" && locationName != "" {
				locationName += ", " + props.Country
			}
			if locationName == "" {
				locationName = props.DisplayName
			}

			return lat, lon, locationName, nil
		}
	} else {
		var r []struct {
			Lat         string `json:"lat"`
			Lon         string `json:"lon"`
			DisplayName string `json:"display_name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return 0, 0, "", err
		}
		if len(r) > 0 {
			lat, err := strconv.ParseFloat(r[0].Lat, 64)
			if err != nil {
				return 0, 0, "", fmt.Errorf("invalid latitude from API: %v", err)
			}
			lon, err := strconv.ParseFloat(r[0].Lon, 64)
			if err != nil {
				return 0, 0, "", fmt.Errorf("invalid longitude from API: %v", err)
			}
			return lat, lon, r[0].DisplayName, nil
		}
	}

	return 0, 0, "", fmt.Errorf("no results found for %q", query)
}

// ReverseGeocode converts coordinates to a location name
func ReverseGeocode(apiType string, lat, lon float64) (string, error) {
	var apiURL string
	switch apiType {
	case "photon":
		apiURL = fmt.Sprintf("https://photon.komoot.io/reverse?lat=%f&lon=%f", lat, lon)
	default: // nominatim
		apiURL = fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json",
			lat, lon)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	if apiType != "photon" {
		req.Header.Set("User-Agent", "jadwalsalat/1.0")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	if apiType == "photon" {
		var r struct {
			Features []struct {
				Properties struct {
					Name        string `json:"name"`
					City        string `json:"city"`
					State       string `json:"state"`
					Country     string `json:"country"`
					DisplayName string `json:"display_name"`
				} `json:"properties"`
			} `json:"features"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return "", err
		}
		if len(r.Features) > 0 {
			props := r.Features[0].Properties

			// Build location name from available properties
			locationName := props.Name
			if props.City != "" {
				locationName = props.City
			}
			if props.State != "" && locationName != "" {
				locationName += ", " + props.State
			}
			if props.Country != "" && locationName != "" {
				locationName += ", " + props.Country
			}
			if locationName == "" {
				locationName = props.DisplayName
			}

			return locationName, nil
		}
	} else {
		var r struct {
			DisplayName string `json:"display_name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return "", err
		}
		if r.DisplayName != "" {
			return r.DisplayName, nil
		}
	}

	return "", fmt.Errorf("no location found for coordinates %f, %f", lat, lon)
}
