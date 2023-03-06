package common

import "strings"

// https://en.wikipedia.org/wiki/List_of_Internet_top-level_domains
var specialDomains map[string]bool

func init() {
	specialDomains = make(map[string]bool)

	for _, name := range [...]string{"com", "org", "net", "int", "edu", "gov", "mil"} {
		specialDomains[name] = true
	}
}

func GetShortDomain(domain string) string {
	split := strings.Split(domain, ".")
	if len(split) <= 2 {
		return domain
	}
	secondLevel := split[len(split)-2]
	if specialDomains[secondLevel] {
		return strings.Join(split[len(split)-3:], ".")
	} else {
		return strings.Join(split[len(split)-2:], ".")
	}
}
