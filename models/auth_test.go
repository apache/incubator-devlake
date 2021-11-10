package models

import (
	"testing"
)

func TestGetAuthToken(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"test for GetAuthToken", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAuthToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
