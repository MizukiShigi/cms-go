package context

import (
	"context"
	"sync"
	"testing"
)

func TestWithValue_NewContext(t *testing.T) {
	parent := context.Background()
	key := "test_key"
	value := "test_value"
	
	ctx := WithValue(parent, key, value)
	
	// Check that the context has the logging sync.Map
	loggingValue := ctx.Value(Logging)
	if loggingValue == nil {
		t.Error("WithValue() should set Logging in context")
		return
	}
	
	syncMap, ok := loggingValue.(*sync.Map)
	if !ok {
		t.Error("Logging value should be *sync.Map")
		return
	}
	
	// Check that the key-value pair is stored
	stored, exists := syncMap.Load(key)
	if !exists {
		t.Error("Key should be stored in sync.Map")
		return
	}
	
	if stored != value {
		t.Errorf("WithValue() stored value = %v, want %v", stored, value)
	}
}

func TestWithValue_ExistingContext(t *testing.T) {
	// Create context with existing sync.Map
	initialMap := &sync.Map{}
	initialMap.Store("existing_key", "existing_value")
	parent := context.WithValue(context.Background(), Logging, initialMap)
	
	key := "new_key"
	value := "new_value"
	
	ctx := WithValue(parent, key, value)
	
	// Check that the context has the logging sync.Map
	loggingValue := ctx.Value(Logging)
	if loggingValue == nil {
		t.Error("WithValue() should preserve Logging in context")
		return
	}
	
	syncMap, ok := loggingValue.(*sync.Map)
	if !ok {
		t.Error("Logging value should be *sync.Map")
		return
	}
	
	// Check that the new key-value pair is stored
	stored, exists := syncMap.Load(key)
	if !exists {
		t.Error("New key should be stored in sync.Map")
		return
	}
	
	if stored != value {
		t.Errorf("WithValue() stored value = %v, want %v", stored, value)
	}
	
	// Check that existing key-value pair is preserved
	existingStored, existingExists := syncMap.Load("existing_key")
	if !existingExists {
		t.Error("Existing key should be preserved in sync.Map")
		return
	}
	
	if existingStored != "existing_value" {
		t.Errorf("WithValue() preserved value = %v, want %v", existingStored, "existing_value")
	}
}

func TestWithValue_NilParent(t *testing.T) {
	key := "test_key"
	value := "test_value"
	
	ctx := WithValue(nil, key, value)
	
	// Should use background context when parent is nil
	if ctx == nil {
		t.Error("WithValue() should not return nil context")
		return
	}
	
	// Check that the context has the logging sync.Map
	loggingValue := ctx.Value(Logging)
	if loggingValue == nil {
		t.Error("WithValue() should set Logging in context")
		return
	}
	
	syncMap, ok := loggingValue.(*sync.Map)
	if !ok {
		t.Error("Logging value should be *sync.Map")
		return
	}
	
	// Check that the key-value pair is stored
	stored, exists := syncMap.Load(key)
	if !exists {
		t.Error("Key should be stored in sync.Map")
		return
	}
	
	if stored != value {
		t.Errorf("WithValue() stored value = %v, want %v", stored, value)
	}
}

func TestWithValue_MultipleValues(t *testing.T) {
	parent := context.Background()
	
	ctx1 := WithValue(parent, "key1", "value1")
	ctx2 := WithValue(ctx1, "key2", "value2")
	ctx3 := WithValue(ctx2, "key3", "value3")
	
	// Check that all values are stored
	loggingValue := ctx3.Value(Logging)
	syncMap, ok := loggingValue.(*sync.Map)
	if !ok {
		t.Error("Logging value should be *sync.Map")
		return
	}
	
	tests := []struct {
		key   string
		value string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
	}
	
	for _, tt := range tests {
		stored, exists := syncMap.Load(tt.key)
		if !exists {
			t.Errorf("Key %s should be stored in sync.Map", tt.key)
			continue
		}
		
		if stored != tt.value {
			t.Errorf("WithValue() stored value for %s = %v, want %v", tt.key, stored, tt.value)
		}
	}
}

func TestCopySyncMap(t *testing.T) {
	original := &sync.Map{}
	original.Store("key1", "value1")
	original.Store("key2", "value2")
	
	copy := copySyncMap(original)
	
	// Check that copy has the same values
	tests := []struct {
		key   string
		value string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
	}
	
	for _, tt := range tests {
		stored, exists := copy.Load(tt.key)
		if !exists {
			t.Errorf("Key %s should be copied", tt.key)
			continue
		}
		
		if stored != tt.value {
			t.Errorf("copySyncMap() copied value for %s = %v, want %v", tt.key, stored, tt.value)
		}
	}
	
	// Check that modifying copy doesn't affect original
	copy.Store("key3", "value3")
	
	_, exists := original.Load("key3")
	if exists {
		t.Error("Modifying copy should not affect original")
	}
	
	// Check that modifying original doesn't affect copy
	original.Store("key4", "value4")
	
	_, exists = copy.Load("key4")
	if exists {
		t.Error("Modifying original should not affect copy")
	}
}