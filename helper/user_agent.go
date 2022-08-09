package helper

import (
	"fmt"
	"strings"

	ua "github.com/mileusna/useragent"
)

// return parsed user agent, bool isMobile and bool isApp
func ParseUserAgent(userAgent string) (string, bool, bool) {
	if strings.HasPrefix(userAgent, "BhinnekaApp") {
		return userAgent, true, true
	}
	m := ua.Parse(userAgent)
	browser := fmt.Sprintf("%s %s", m.Name, m.Version)
	if m.OS != "" && m.OSVersion != "" {
		browser = fmt.Sprintf("%s - %s", browser, fmt.Sprintf("%s %s", m.OS, m.OSVersion))
	}

	if m.Device != "" {
		browser += fmt.Sprintf(" [%s]", m.Device)
	}

	return browser, m.Mobile, false
}
