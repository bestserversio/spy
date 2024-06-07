package utils

import (
	"strconv"
	"strings"
)

func ContainsStrings(str string, matches []string) bool {
	for _, m := range matches {
		if strings.Contains(str, m) {
			return true
		}
	}

	return false
}

func ContainsIps(ip string, match_ips []string) (bool, error) {
	var err error
	for _, m_ip := range match_ips {
		// Split by CIDR.
		data := strings.Split(m_ip, "/")

		// If length is 1 or lower, perform simple match.
		if len(data) <= 1 {
			if ip == m_ip {
				return true, err
			}
		} else {
			// Convert CIDR string to integer.
			cidr, err := strconv.Atoi(data[1])

			if err != nil {
				return false, err
			}

			// If we have a /32, perform a simple match.
			if cidr == 32 && ip == m_ip {
				return true, err
			}

			// To Do: Check if IP is in range.
		}
	}
	return false, err
}
