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

package function

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/jrnd-io/jrv2/pkg/config"
	"math/big"
	"net"
	"text/template"
)

// adding functions to map
func init() {
	AddFuncs(template.FuncMap{
		"http_method":       HTTPMethod,
		"ip":                IP,
		"ipv6":              IPv6,
		"ip_known_protocol": IPKnownProtocol,
		"ip_known_port":     IPKnownPort,
		"mac":               Mac,
		"password":          Password,
		"useragent":         UserAgent,
	})
}

var (
	HTTPMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	Ports       = []string{"80", "81", "443", "22", "631"}
	Protocols   = []string{"TCP", "UDP", "ICMP", "FTP", "HTTP", "SFTP"}
)

// HTTPMethod returns a random http method
func HTTPMethod() string {
	return HTTPMethods[config.Random.IntN(len(HTTPMethods))]
}

func IP(cidr string) string {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "0.0.0.0"
	}

	// Convert the IPNet to a big.Int
	ipInt := big.NewInt(0).SetBytes(ipnet.IP.To4())
	//	maskInt := big.NewInt(0).SetBytes(ipnet.Mask)

	// Calculate the network size
	maskSize, _ := ipnet.Mask.Size()
	networkSize := big.NewInt(0).Sub(big.NewInt(1).Lsh(big.NewInt(1), uint(32-maskSize)), big.NewInt(1)) //nolint:gosec // disable G115
	networkSize.Sub(networkSize, big.NewInt(2))                                                          // Exclude network and broadcast addresses

	// Generate a random offset within the network size
	offset, err := rand.Int(rand.Reader, networkSize)
	if err != nil {
		return "0.0.0.0"
	}

	// Add the offset to the network address
	ipInt.Add(ipInt, offset)
	ipInt.Add(ipInt, big.NewInt(1)) // Skip the network address

	// Convert the big.Int back to an IP address
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, uint32(ipInt.Uint64())) //nolint

	return ip.String()
}

// IPKnownPort returns a random known port number
func IPKnownPort() string {
	return Ports[config.Random.IntN(len(Ports))]
}

// IPKnownProtocol returns a random known protocol
func IPKnownProtocol() string {
	return Protocols[config.Random.IntN(len(Protocols))]
}

// IPv6 returns a random Ipv6 Address
func IPv6() string {
	ip := make(net.IP, net.IPv6len)
	for i := 0; i < net.IPv6len; i++ {
		ip[i] = byte(config.Random.IntN(256))
	}
	ip[0] &= 0xfe // Set the "locally administered" flag
	ip[0] |= 0x02 // Set the "unicast" flag
	return ip.String()
}

// Mac returns a random Mac Address
func Mac() string {
	mac := make(net.HardwareAddr, 8)
	// config.ChaCha8.Read(mac) //nolint
	binary.LittleEndian.PutUint64(mac, config.Random.Uint64())
	mac[0] &= 0xfe // Set the "locally administered" flag
	mac[0] |= 0x02 // Set the "unicast" flag
	return mac[0:6].String()
}

// Password returns a random Password of given length, memorable, and with prefix and suffix
func Password(length int, memorable bool, prefix string, suffix string) string {

	const (
		// Define the set of vowels and consonants that can be used to generate the Password.
		vowels     = "aeiouyAEIOUY"
		consonants = "bcdfghjklmnpqrstvwxzBCDFGHJKLMNPQRSTVWXZ"
	)

	// Generate a memorable Password of the specified length.
	password := make([]byte, length)
	if memorable {
		// Generate a memorable Password.
		for i := range password {
			if i%2 == 0 {
				// Use a vowel.
				char := vowels[config.Random.IntN(len(vowels))]
				password[i] = char
			} else {
				// Use a consonant.
				char := consonants[config.Random.IntN(len(consonants))]
				password[i] = char
			}
		}
	} else {
		// Generate a random Password using the full charset.
		charset := vowels + consonants + "0123456789!@#$%^&*()_+{}:\"<>?,./;'[]\\-=`~"
		for i := range password {
			char := charset[config.Random.IntN(len(charset))]
			password[i] = char
		}
	}

	return prefix + string(password) + suffix
}

// UserAgent returns a random user agent
func UserAgent() string {

	var desktopOperatingSystems = []string{
		"Windows NT 10.0", "Windows NT 6.3", "Macintosh; Intel Mac OS X 10_15_7", "Macintosh; Intel Mac OS X 10_14_5", "X11; Linux x86_64",
	}

	var mobileOperatingSystems = []string{
		"Android 11", "Android 10", "iOS 14_4_2", "iOS 14_0",
	}

	var desktopBrowsers = []string{
		"Chrome", "Safari", "Firefox", "Opera", "Edge",
	}

	var mobileBrowsers = []string{
		"Chrome Mobile", "Safari Mobile", "Firefox Mobile", "Opera Mobile", "Edge Mobile",
	}

	// Generate random desktop user agent
	isDesktop := config.Random.IntN(2) == 0
	var os string
	var browser string
	var version string
	if isDesktop {
		os = desktopOperatingSystems[config.Random.IntN(len(desktopOperatingSystems))]
		browser = desktopBrowsers[config.Random.IntN(len(desktopBrowsers))]
		version = fmt.Sprintf("%d.%d.%d.%d", config.Random.IntN(10), config.Random.IntN(10), config.Random.IntN(10), config.Random.IntN(10))
	} else {
		os = mobileOperatingSystems[config.Random.IntN(len(mobileOperatingSystems))]
		browser = mobileBrowsers[config.Random.IntN(len(mobileBrowsers))]
		switch browser {
		case "Chrome Mobile":
			version = fmt.Sprintf("%d.%d.%d.%d", config.Random.IntN(10), config.Random.IntN(10), config.Random.IntN(10), config.Random.IntN(10))
		case "Safari Mobile":
			version = fmt.Sprintf("%d.%d", config.Random.IntN(14)+1, config.Random.IntN(3)+1)
		case "Firefox Mobile":
			version = fmt.Sprintf("%d.%d", config.Random.IntN(10)+1, config.Random.IntN(10))
		case "Opera Mobile":
			version = fmt.Sprintf("%d.%d.%d.%d", config.Random.IntN(10), config.Random.IntN(10), config.Random.IntN(10), config.Random.IntN(10))
		case "Edge Mobile":
			version = fmt.Sprintf("%d.%d.%d.%d", config.Random.IntN(10)+40, config.Random.IntN(10), config.Random.IntN(10), config.Random.IntN(10))
		}
	}

	userAgent := fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/%d.%d (KHTML, like Gecko) %s/%s Mobile Safari/%d.%d", os, config.Random.IntN(100)+500, config.Random.IntN(100)+1, browser, version, config.Random.IntN(10)+1, config.Random.IntN(10)+1)

	return userAgent

}
