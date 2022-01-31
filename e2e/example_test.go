package e2e

import "testing"

// This test should only run once main_test is complete and ready

func TestSum(t *testing.T) {
	total := 5
	if total != 10 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	}
}
