// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package function_test

import (
	"net"
	"regexp"
	"testing"

	"github.com/jrnd-io/jrv2/pkg/function"
	"github.com/stretchr/testify/assert"
)

func TestHTTPMethod(t *testing.T) {
	// Define the expected HTTP methods
	expectedMethods := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
		"PATCH":  true,
	}

	// Call the HTTPMethod function multiple times
	for i := 0; i < 100; i++ {
		method := function.HTTPMethod()

		// Check if the returned method is in the list of expected methods
		if !expectedMethods[method] {
			t.Errorf("Unexpected HTTP method: %s", method)
		}
	}
}

func TestIP(t *testing.T) {
	testCases := []struct {
		cidr string
	}{
		{"192.168.1.0/24"},
		{"10.0.0.0/8"},
		{"172.16.0.0/12"},
	}

	for _, tc := range testCases {
		t.Run(tc.cidr, func(t *testing.T) {
			ipStr := function.IP(tc.cidr)
			ip := net.ParseIP(ipStr)
			if ip == nil {
				t.Errorf("Generated IP %s is not a valid IP address", ipStr)
			}

			_, ipnet, err := net.ParseCIDR(tc.cidr)
			if err != nil {
				t.Fatalf("Failed to parse CIDR %s: %v", tc.cidr, err)
			}

			if !ipnet.Contains(ip) {
				t.Errorf("Generated IP %s is not within the CIDR block %s", ipStr, tc.cidr)
			}
		})
	}
}

func TestIPKnownPort(t *testing.T) {
	// Define the expected list of ports
	expectedPorts := function.Ports

	// Call the IPKnownPort function multiple times
	for i := 0; i < 100; i++ {
		port := function.IPKnownPort()

		// Check if the returned port is in the list of expected ports
		assert.Contains(t, expectedPorts, port)
	}
}

func TestIPKnownProtocol(t *testing.T) {
	// Define the expected list of protocols
	expectedProtocols := function.Protocols

	// Call the IPKnownProtocol function multiple times
	for i := 0; i < 100; i++ {
		protocol := function.IPKnownProtocol()

		assert.Contains(t, expectedProtocols, protocol)

	}
}

func TestIPv6(t *testing.T) {
	for i := 0; i < 100; i++ {
		ipStr := function.IPv6()
		ip := net.ParseIP(ipStr)
		if ip == nil || ip.To16() == nil || ip.To4() != nil {
			t.Errorf("Generated IP %s is not a valid IPv6 address", ipStr)
		}

		// Check the "locally administered" and "unicast" flags
		if ip[0]&0x02 == 0 {
			t.Errorf("Generated IP %s does not have the 'unicast' flag set", ipStr)
		}
		if ip[0]&0xfe == 0 {
			t.Errorf("Generated IP %s does not have the 'locally administered' flag set", ipStr)
		}
	}
}

func TestMac(t *testing.T) {
	for i := 0; i < 100; i++ {
		macStr := function.Mac()
		mac, err := net.ParseMAC(macStr)
		if err != nil {
			t.Errorf("Generated MAC %s is not a valid MAC address: %v", macStr, err)
		}

		// Check the "locally administered" and "unicast" flags
		if mac[0]&0x02 == 0 {
			t.Errorf("Generated MAC %s does not have the 'unicast' flag set", macStr)
		}
		if mac[0]&0xfe == 0 {
			t.Errorf("Generated MAC %s does not have the 'locally administered' flag set", macStr)
		}
	}
}

func TestPassword(t *testing.T) {
	tests := []struct {
		length    int
		memorable bool
		prefix    string
		suffix    string
	}{
		{8, true, "pre-", "-suf"},
		{12, false, "", ""},
		{10, true, "start-", "-end"},
		{15, false, "pre-", ""},
	}

	for _, tt := range tests {
		password := function.Password(tt.length, tt.memorable, tt.prefix, tt.suffix)

		// Check the length of the generated password
		expectedLength := tt.length + len(tt.prefix) + len(tt.suffix)
		if len(password) != expectedLength {
			t.Errorf("Expected password length %d, got %d", expectedLength, len(password))
		}

		// Check the prefix and suffix
		if !startsWith(password, tt.prefix) {
			t.Errorf("Expected password to start with %s, got %s", tt.prefix, password)
		}
		if !endsWith(password, tt.suffix) {
			t.Errorf("Expected password to end with %s, got %s", tt.suffix, password)
		}

		// Check the memorable pattern
		if tt.memorable {
			passwordBody := password[len(tt.prefix) : len(password)-len(tt.suffix)]
			for i, char := range passwordBody {
				if i%2 == 0 {
					if !isVowel(char) {
						t.Errorf("Expected vowel at position %d, got %c", i, char)
					}
				} else {
					if !isConsonant(char) {
						t.Errorf("Expected consonant at position %d, got %c", i, char)
					}
				}
			}
		}
	}
}

func TestUserAgent(t *testing.T) {
	// Define a regular expression to match the expected user agent format
	userAgentRegex := regexp.MustCompile(`^Mozilla\/5\.0 \([^)]+\) AppleWebKit\/\d+\.\d+ \(KHTML, like Gecko\) [^\/]+\/[0-9\.]+ Mobile Safari\/\d+\.\d+$`)

	for i := 0; i < 100; i++ {
		userAgent := function.UserAgent()

		// Check if the generated user agent matches the expected format
		if !userAgentRegex.MatchString(userAgent) {
			t.Errorf("Generated user agent does not match the expected format: %s", userAgent)
		}
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func isVowel(c rune) bool {
	vowels := "aeiouyAEIOUY"
	for _, v := range vowels {
		if c == v {
			return true
		}
	}
	return false
}

func isConsonant(c rune) bool {
	consonants := "bcdfghjklmnpqrstvwxzBCDFGHJKLMNPQRSTVWXZ"
	for _, v := range consonants {
		if c == v {
			return true
		}
	}
	return false
}
