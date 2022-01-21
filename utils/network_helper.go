package utils

import (
	"fmt"
	"net"
	"time"

	"github.com/merico-dev/lake/logger"
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
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return err
	}
	logger.Info("Connect successfully", map[string]string{
		"Target": target,
		"Remote": conn.RemoteAddr().String(),
		"Local":  conn.LocalAddr().String(),
	})
	return nil
}

func ResolvePort(port string, schema string) (string, error) {
	var defaultPorts = map[string]string{
		"http":  "80",
		"https": "443",
	}
	if port != "" {
		if schema != "" {
			logger.Warn("both port and schema is provided, will using port directly", map[string]string{
				port:   port,
				schema: schema,
			})
		}
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
