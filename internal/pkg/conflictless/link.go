package conflictless

import (
	"fmt"
	"log"
	"strings"
)

func SectionLink(baseURL string, sectionName string) string {
	if baseURL == "" {
		return ""
	}

	if strings.Contains(baseURL, "github.com") {
		return fmt.Sprintf("[%s]: %s/releases/tag/v%s", sectionName, baseURL, sectionName)
	}

	if strings.Contains(baseURL, "gitlab.") {
		return fmt.Sprintf("[%s]: %s/-/releases/v%s", sectionName, baseURL, sectionName)
	}

	log.Print("Unknown repository host, skipping section link generation")

	return ""
}
