package conflictless

import (
	"fmt"
	"log"
	"strings"
)

func SectionLink(baseURL, sectionName string, useVPrefixInLinks bool) string {
	if baseURL == "" {
		return ""
	}

	prefix := ""
	if useVPrefixInLinks {
		prefix = "v"
	}

	if strings.Contains(baseURL, "github.com") {
		return fmt.Sprintf("[%s]: %s/releases/tag/%s%s", sectionName, baseURL, prefix, sectionName)
	}

	if strings.Contains(baseURL, "gitlab.") {
		return fmt.Sprintf("[%s]: %s/-/releases/%s%s", sectionName, baseURL, prefix, sectionName)
	}

	log.Print("Unknown repository host, skipping section link generation")

	return ""
}
