package utils

import (
	"testing"
)

func TestCheckDNS(t *testing.T) {
	var hostname = "baidu.com"
	var err = CheckDNS(hostname)
	if err != nil {
		t.Error(err)
	}

	var invalidHostname = "baidu.abc"
	err = CheckDNS(invalidHostname)
	t.Log(err)
	if err == nil {
		t.Errorf("Expected %s, Got nil", "failed")
	}
}

func TestResolvePort(t *testing.T) {
	port, err := ResolvePort("80", "https")
	if err != nil {
		t.Error(err)
	}
	if port != "80" {
		t.Errorf("Expected %s, Got %s", "80", port)
	}
	port, err = ResolvePort("", "http")
	if err != nil {
		t.Error(err)
	}
	if port != "80" {
		t.Errorf("Expected %s, Got %s", "80", port)
	}
	_, err = ResolvePort("", "rabbitmq")
	if err == nil {
		t.Errorf("Expected error %s, Got nil", "schema not fount")
	}
	_, err = ResolvePort("", "")
	if err == nil {
		t.Errorf("Expected error %s, Got nil", "you should provide at least one of port or schema")
	}
}
