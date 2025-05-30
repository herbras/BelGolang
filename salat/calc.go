package salat

import (
	"math"
	"time"
)

// CalculationMethod represents different methods for calculating prayer times
type CalculationMethod string

const (
	// MWL - Muslim World League
	MWL CalculationMethod = "MWL"
	// ISNA - Islamic Society of North America
	ISNA CalculationMethod = "ISNA"
	// Egypt - Egyptian General Authority of Survey
	Egypt CalculationMethod = "Egypt"
	// Makkah - Umm al-Qura University, Makkah
	Makkah CalculationMethod = "Makkah"
	// Karachi - University of Islamic Sciences, Karachi
	Karachi CalculationMethod = "Karachi"
	// Tehran - Institute of Geophysics, University of Tehran
	Tehran CalculationMethod = "Tehran"
	// Kemenag - Kementerian Agama Republik Indonesia
	Kemenag CalculationMethod = "Kemenag"
	// JAKIM - Jabatan Kemajuan Islam Malaysia
	JAKIM CalculationMethod = "JAKIM"
)

type Location struct {
	Latitude  float64
	Longitude float64
	Method    CalculationMethod
}

type PrayerTimes struct {
	Imsak   time.Time
	Subuh   time.Time
	Dzuhur  time.Time
	Ashar   time.Time
	Maghrib time.Time
	Isya    time.Time
}

type methodParams struct {
	fajrAngle    float64
	ishaAngle    float64
	ishaInterval float64
}

func getMethodParams(method CalculationMethod) methodParams {
	switch method {
	case MWL:
		return methodParams{fajrAngle: 18, ishaAngle: 17, ishaInterval: 0}
	case ISNA:
		return methodParams{fajrAngle: 15, ishaAngle: 15, ishaInterval: 0}
	case Egypt:
		return methodParams{fajrAngle: 19.5, ishaAngle: 17.5, ishaInterval: 0}
	case Makkah:
		return methodParams{fajrAngle: 18.5, ishaAngle: 0, ishaInterval: 90}
	case Karachi:
		return methodParams{fajrAngle: 18, ishaAngle: 18, ishaInterval: 0}
	case Tehran:
		return methodParams{fajrAngle: 17.7, ishaAngle: 14, ishaInterval: 0}
	case Kemenag:
		return methodParams{fajrAngle: 20, ishaAngle: 18, ishaInterval: 0}
	case JAKIM:
		return methodParams{fajrAngle: 20, ishaAngle: 18, ishaInterval: 0}
	default:
		return methodParams{fajrAngle: 18, ishaAngle: 17, ishaInterval: 0}
	}
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func radiansToDegrees(radians float64) float64 {
	return radians * 180.0 / math.Pi
}

func normalizeAngle(angle float64) float64 {
	if angle >= 0 {
		return math.Mod(angle, 360)
	}
	return math.Mod(angle, 360) + 360
}

func normalizeHours(hours float64) float64 {
	if hours < 0 {
		return hours + 24
	}
	if hours >= 24 {
		return hours - 24
	}
	return hours
}

func calculateJulianDate(t time.Time) float64 {
	year, month, day := t.Date()
	if month <= 2 {
		year--
		month += 12
	}

	A := math.Floor(float64(year) / 100.0)
	B := 2 - A + math.Floor(A/4.0)

	jd := math.Floor(365.25*float64(year+4716)) +
		math.Floor(30.6001*float64(month+1)) +
		float64(day) + B - 1524.5

	hour, min, sec := t.Clock()
	jd += (float64(hour) + float64(min)/60.0 + float64(sec)/3600.0) / 24.0

	return jd
}

func calculateSolarPosition(jd float64) (declination, eqOfTime float64) {
	D := jd - 2451545.0

	g := normalizeAngle(357.529 + 0.98560028*D)

	q := normalizeAngle(280.459 + 0.98564736*D)

	L := normalizeAngle(q + 1.915*math.Sin(degreesToRadians(g)) + 0.020*math.Sin(degreesToRadians(2*g)))

	e := 23.439 - 0.00000036*D

	declination = radiansToDegrees(math.Asin(math.Sin(degreesToRadians(e)) * math.Sin(degreesToRadians(L))))

	y := math.Pow(math.Tan(degreesToRadians(e/2.0)), 2)
	eqOfTime = 4.0 * radiansToDegrees(y*math.Sin(2.0*degreesToRadians(q))-
		2.0*math.Sin(degreesToRadians(g))+
		4.0*y*math.Sin(degreesToRadians(g))*math.Cos(2.0*degreesToRadians(q))-
		0.5*y*y*math.Sin(4.0*degreesToRadians(q))-
		1.25*math.Sin(degreesToRadians(2*g)))

	return declination, eqOfTime
}

func calculatePrayerTime(angle, latitude, declination, eqOfTime, timezone, longitude float64, isFajrOrSunrise bool) float64 {
	angleRad := degreesToRadians(angle)
	latRad := degreesToRadians(latitude)
	declRad := degreesToRadians(declination)

	hourAngle := radiansToDegrees(math.Acos(
		(math.Sin(angleRad) - math.Sin(latRad)*math.Sin(declRad)) /
			(math.Cos(latRad) * math.Cos(declRad)),
	))

	noon := 12 + timezone - longitude/15.0 - eqOfTime/60.0

	// Calculate prayer time in hours since midnight
	// For Fajr and Sunrise, we subtract the hour angle
	// For other prayers (Asr, Maghrib, Isha), we add the hour angle
	var time float64
	if isFajrOrSunrise {
		time = noon - hourAngle/15.0
	} else {
		time = noon + hourAngle/15.0
	}

	return normalizeHours(time)
}

// TimesForDate calculates prayer times for a specific date and location
func TimesForDate(date time.Time, loc Location) (PrayerTimes, error) {
	_, offset := date.Zone()
	timezoneOffset := float64(offset) / 3600.0

	jd := calculateJulianDate(date)

	declination, eqOfTime := calculateSolarPosition(jd)

	params := getMethodParams(loc.Method)

	noon := 12 + timezoneOffset - loc.Longitude/15.0 - eqOfTime/60.0

	// Fajr time
	fajrTime := calculatePrayerTime(-params.fajrAngle, loc.Latitude, declination, eqOfTime, timezoneOffset, loc.Longitude, true)

	// Sunrise time (for validation and reference)
	_ = calculatePrayerTime(-0.833, loc.Latitude, declination, eqOfTime, timezoneOffset, loc.Longitude, true)

	// Dzuhur time is same as noon
	dhuhrTime := normalizeHours(noon)

	// Asr time (Shafi'i, shadow factor = 1)
	asrFactor := 1.0
	// Menghitung sudut Ashar dengan memperhitungkan bayangan saat tengah hari
	// Rumus yang benar untuk Ashar: cotg(asrAngle) = asrFactor + tan(abs(latitude-declination))
	asrAngle := radiansToDegrees(math.Atan(1.0 / (asrFactor + math.Tan(degreesToRadians(math.Abs(loc.Latitude-declination))))))
	asrTime := calculatePrayerTime(asrAngle, loc.Latitude, declination, eqOfTime, timezoneOffset, loc.Longitude, false)

	// Maghrib time (same as sunset)
	maghribTime := calculatePrayerTime(-0.833, loc.Latitude, declination, eqOfTime, timezoneOffset, loc.Longitude, false)

	// Isha time
	ishaTime := 0.0
	if params.ishaInterval > 0 {
		// Isha is calculated as minutes after maghrib for some methods
		ishaTime = normalizeHours(maghribTime + params.ishaInterval/60.0)
	} else {
		// Isha is calculated based on sun angle
		ishaTime = calculatePrayerTime(-params.ishaAngle, loc.Latitude, declination, eqOfTime, timezoneOffset, loc.Longitude, false)
	}

	// Imsak time (10 minutes before Fajr)
	imsakTime := normalizeHours(fajrTime - 10.0/60.0)

	// Convert hours to time.Time
	baseDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	convertHoursToTime := func(hours float64) time.Time {
		h := int(hours)
		m := int((hours - float64(h)) * 60.0)
		s := int(((hours-float64(h))*60.0 - float64(m)) * 60.0)
		return baseDate.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second)
	}

	return PrayerTimes{
		Imsak:   convertHoursToTime(imsakTime),
		Subuh:   convertHoursToTime(fajrTime),
		Dzuhur:  convertHoursToTime(dhuhrTime),
		Ashar:   convertHoursToTime(asrTime),
		Maghrib: convertHoursToTime(maghribTime),
		Isya:    convertHoursToTime(ishaTime),
	}, nil
}

// GetCurrentPrayer returns the current active prayer time
func GetCurrentPrayer(t time.Time, times PrayerTimes) (string, bool) {
	// Define prayer times with their start and end times
	prayerPeriods := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
	}{
		{"Imsak", times.Imsak, times.Subuh},
		{"Subuh", times.Subuh, times.Dzuhur},
		{"Dzuhur", times.Dzuhur, times.Ashar},
		{"Ashar", times.Ashar, times.Maghrib},
		{"Maghrib", times.Maghrib, times.Isya},
		{"Isya", times.Isya, times.Imsak.Add(24 * time.Hour)},
	}

	// Check if current time is within any prayer period
	for _, period := range prayerPeriods {
		if (t.Equal(period.startTime) || t.After(period.startTime)) && t.Before(period.endTime) {
			return period.name, true
		}
	}

	return "", false
}

// GetNextPrayer returns the next prayer time after the given time
// GetPrayerEmoji returns an emoji for each prayer time
func GetPrayerEmoji(prayerName string) string {
	switch prayerName {
	case "Imsak":
		return "🌙 " // Bulan - masih malam
	case "Subuh":
		return "🌅 " // Matahari terbit
	case "Dzuhur":
		return "☀️ " // Matahari penuh
	case "Ashar":
		return "🌤️ " // Matahari dengan awan
	case "Maghrib":
		return "🌇 " // Matahari terbenam
	case "Isya":
		return "✨ " // Bintang - malam
	default:
		return ""
	}
}

func GetNextPrayer(t time.Time, times PrayerTimes) (string, time.Time) {
	// Pastikan waktu sholat diurutkan berdasarkan waktu
	prayerTimes := []struct {
		name string
		time time.Time
	}{
		{"Imsak", times.Imsak},
		{"Subuh", times.Subuh},
		{"Dzuhur", times.Dzuhur},
		{"Ashar", times.Ashar},
		{"Maghrib", times.Maghrib},
		{"Isya", times.Isya},
	}

	// Urutkan waktu sholat berdasarkan waktu
	for i := 0; i < len(prayerTimes)-1; i++ {
		for j := i + 1; j < len(prayerTimes); j++ {
			if prayerTimes[i].time.After(prayerTimes[j].time) {
				prayerTimes[i], prayerTimes[j] = prayerTimes[j], prayerTimes[i]
			}
		}
	}

	// Cek waktu sholat berikutnya
	for _, prayer := range prayerTimes {
		if prayer.time.After(t) {
			return prayer.name, prayer.time
		}
	}

	// Jika tidak ada waktu sholat yang tersisa di hari ini, kembalikan waktu sholat pertama untuk hari berikutnya
	// Biasanya ini adalah Imsak atau Subuh
	nextDayTime := t.Add(24 * time.Hour)
	if len(prayerTimes) > 0 {
		// Ambil waktu sholat pertama (setelah diurutkan)
		firstPrayer := prayerTimes[0]
		nextTime := time.Date(
			nextDayTime.Year(), nextDayTime.Month(), nextDayTime.Day(),
			firstPrayer.time.Hour(), firstPrayer.time.Minute(), firstPrayer.time.Second(), 0,
			t.Location(),
		)
		return firstPrayer.name + " (besok)", nextTime
	}

	// Fallback jika tidak ada waktu sholat yang tersedia
	return "Subuh (besok)", times.Subuh.Add(24 * time.Hour)
}
