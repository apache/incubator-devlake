/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
