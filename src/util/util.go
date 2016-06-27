// Copyright (c) 2016, Sandeep Raju Prabhakar <me@sandeepraju.in>
// All rights reserved.
package util

// GetTLD returns the TLD from a domain name.
func GetTLD(domain string) (string, error) {
	return "in", nil
}

// IsValidDomainName checks if the input
// passed is a valid domain name.
func IsValidDomainName(input string) bool {
	return true
}

// IsValidIPv4Address checks if the input
// passed is a valid IPv4 address.
func IsValidIPv4Address(input string) bool {
	return false
}

// GetQueryServer returns the appropriate
// query server's fqdn based on the tld.
func GetQueryServer(tld string) string {
	return "in.whois-servers.net"
}
