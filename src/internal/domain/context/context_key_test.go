package context

import "testing"

func TestContextKeyConstants(t *testing.T) {
	// Test that context key constants have the expected values
	tests := []struct {
		name     string
		key      ContextKey
		expected string
	}{
		{"UserID", UserID, "user_id"},
		{"Logging", Logging, "logging"},
		{"TransactionDB", TransactionDB, "transaction_db"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.key) != tt.expected {
				t.Errorf("ContextKey %s = %v, want %v", tt.name, string(tt.key), tt.expected)
			}
		})
	}
}

func TestContextKeyType(t *testing.T) {
	// Test that ContextKey is a string type
	var key ContextKey = "test"
	if string(key) != "test" {
		t.Error("ContextKey should be convertible to string")
	}
}

func TestContextKeyUniqueness(t *testing.T) {
	// Test that all context keys are unique
	keys := []ContextKey{UserID, Logging, TransactionDB}
	
	seen := make(map[string]bool)
	for _, key := range keys {
		keyStr := string(key)
		if seen[keyStr] {
			t.Errorf("Duplicate context key found: %s", keyStr)
		}
		seen[keyStr] = true
	}
}