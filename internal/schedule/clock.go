package schedule

import (
	"os"
	"os/exec"
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
	if out, err := exec.Command("locale", "t_fmt").Output(); err == nil {
		if twelve, ok := tFmtUses12Hour(string(out)); ok {
			return twelve
		}
	}
	for _, key := range []string{"LC_TIME", "LC_ALL", "LANG"} {
		if v := os.Getenv(key); v != "" {
			return langUses12Hour(v)
		}
	}
	return false
}

func tFmtUses12Hour(s string) (twelve bool, ok bool) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "t_fmt=")
	s = strings.Trim(s, `"`)
	if s == "" {
		return false, false
	}
	if strings.Contains(s, "%I") || strings.Contains(s, "%l") || strings.Contains(s, "%p") || strings.Contains(s, "%P") || strings.Contains(s, "%r") {
		return true, true
	}
	if strings.Contains(s, "%H") || strings.Contains(s, "%R") || strings.Contains(s, "%T") || strings.Contains(s, "%k") {
		return false, true
	}
	return false, false
}

func langUses12Hour(lang string) bool {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "" || lang == "c" || lang == "posix" {
		return false
	}
	if i := strings.IndexAny(lang, ".@"); i >= 0 {
		lang = lang[:i]
	}
	switch {
	case strings.HasSuffix(lang, "_us"),
		strings.HasSuffix(lang, "_ph"),
		strings.HasSuffix(lang, "_mx"),
		strings.HasSuffix(lang, "_co"),
		strings.HasSuffix(lang, "_hn"),
		strings.HasSuffix(lang, "_gt"),
		strings.HasSuffix(lang, "_sv"),
		strings.HasSuffix(lang, "_ni"),
		strings.HasSuffix(lang, "_pa"),
		strings.HasSuffix(lang, "_pr"),
		strings.HasSuffix(lang, "_in"),
		strings.HasSuffix(lang, "_pk"),
		strings.HasSuffix(lang, "_eg"),
		strings.HasSuffix(lang, "_sa"),
		strings.HasSuffix(lang, "_au"),
		strings.HasSuffix(lang, "_nz"),
		strings.HasSuffix(lang, "_ie"):
		return true
	case strings.HasPrefix(lang, "en_") && strings.HasSuffix(lang, "_ca"):
		return true
	default:
		return false
	}
}
