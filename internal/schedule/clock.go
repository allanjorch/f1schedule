package schedule

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	hour12Once sync.Once
	hour12     bool
	hour12Test *bool
)

func use12Hour() bool {
	if hour12Test != nil {
		return *hour12Test
	}
	hour12Once.Do(func() {
		hour12 = detect12HourClock()
	})
	return hour12
}

func clockLayout() string {
	if use12Hour() {
		return "3:04 PM"
	}
	return "15:04"
}

func FormatCompact(t time.Time) string {
	return t.Format("Mon 02 Jan  " + clockLayout())
}

func FormatCacheTime(t time.Time) string {
	return t.Format("Mon 2 Jan, " + clockLayout())
}

func FormatDateTime(t time.Time) string {
	return t.Format("Mon 2 Jan 2006, " + clockLayout() + " MST")
}

func detect12HourClock() bool {
	if twelve, ok := zoneUses12Hour(localIANAZone()); ok {
		return twelve
	}
	return false
}

func localIANAZone() string {
	if tz := strings.TrimSpace(os.Getenv("TZ")); tz != "" {
		tz = strings.TrimPrefix(tz, ":")
		if tz != "" && tz != "localtime" {
			return tz
		}
	}
	target, err := filepath.EvalSymlinks("/etc/localtime")
	if err != nil {
		return time.Local.String()
	}
	const marker = "/zoneinfo/"
	slash := filepath.ToSlash(target)
	if i := strings.LastIndex(slash, marker); i >= 0 {
		return slash[i+len(marker):]
	}
	return time.Local.String()
}

func zoneUses12Hour(zone string) (twelve bool, ok bool) {
	zone = strings.TrimPrefix(strings.TrimSpace(zone), ":")
	if zone == "" || strings.EqualFold(zone, "Local") {
		return false, false
	}
	if isUTCZone(zone) {
		return false, true
	}
	countries, found := countriesForZone(zone)
	if !found || len(countries) == 0 {
		return false, false
	}
	return countryUses12Hour(countries[0]), true
}

func isUTCZone(zone string) bool {
	z := strings.ToUpper(zone)
	return z == "UTC" || z == "GMT" || z == "UCT" || strings.HasPrefix(z, "ETC/")
}

func countryUses12Hour(cc string) bool {
	switch strings.ToUpper(cc) {
	case "US", "PH", "MX", "CO", "HN", "GT", "SV", "NI", "PA", "PR",
		"IN", "PK", "EG", "SA", "AU", "NZ", "IE", "CA":
		return true
	default:
		return false
	}
}

var (
	zoneCountriesOnce sync.Once
	zoneCountries     map[string][]string
)

func countriesForZone(zone string) ([]string, bool) {
	zoneCountriesOnce.Do(loadZoneCountries)
	c, ok := zoneCountries[zone]
	return c, ok
}

func loadZoneCountries() {
	zoneCountries = map[string][]string{}
	for _, path := range []string{
		"/usr/share/zoneinfo/zone.tab",
		"/usr/share/zoneinfo/zone1970.tab",
	} {
		loadZoneTab(path)
	}
}

func loadZoneTab(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, "	")
		if len(fields) < 3 {
			continue
		}
		name := fields[2]
		if _, exists := zoneCountries[name]; exists {
			continue
		}
		zoneCountries[name] = strings.Split(fields[0], ",")
	}
}
