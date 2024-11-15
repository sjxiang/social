package store

import (
    "testing"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
        expected string
    }{
        {"正常", "2024-11-15 17:18:00", "2024-11-15 17:18:00"},
        {"异常", "invalid", ""},
    }

    for _, test := range tests {
        result := parseTime(test.input)
        if result != test.expected {
            t.Errorf("Expected %s, got %s", test.expected, result)
        }
    }
    
	t.Logf("\nParseTime, 测试完成")
}
