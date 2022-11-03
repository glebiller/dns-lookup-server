package dnslookup

import (
	"regexp"
)

func ValidateIpv4(ipv4 string) bool {
	ipv4Regex := regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
	return ipv4Regex.MatchString(ipv4)
}
