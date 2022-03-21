package utils

import (
	"fmt"
	"net"
	"time"
)

func CheckDNS(domain string) error {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return err
	}
	if len(ips) > 0 {
		return nil
	}
	return fmt.Errorf("failed to resolve host: %s", domain)
}

func CheckNetwork(host, port string, timeout time.Duration) error {
	var target = fmt.Sprintf("%s:%s", host, port)
	_, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return err
	}
	return nil
}

func ResolvePort(port string, schema string) (string, error) {
	var defaultPorts = map[string]string{
		"http":  "80",
		"https": "443",
	}
	if port != "" {
		return port, nil
	}
	if schema != "" {
		port, ok := defaultPorts[schema]
		if !ok {
			return "", fmt.Errorf("schema %s not found", schema)
		}
		return port, nil
	}
	return "", fmt.Errorf("you should provide at least one of port or schema")
}
