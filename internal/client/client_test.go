package client

import (
	"testing"
)

func TestParametersForQuery(t *testing.T) {
	// Test creating a ParametersForQuery instance
	param := QueryParameters{
		Key:   "query",
		Value: "Stranger Things",
	}

	if param.Key != "query" {
		t.Errorf("Expected Key to be 'query', got '%s'", param.Key)
	}

	if param.Value != "Stranger Things" {
		t.Errorf("Expected Value to be 'query', got '%s'", param.Value)
	}
}
