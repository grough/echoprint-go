package echo

import (
	"os"
	"testing"
)

func TestNewRenderer(t *testing.T) {
	testFile := "testdata/mono.wav"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test WAV file %s not found, skipping", testFile)
	}

	renderer, err := NewRenderer(testFile, 120, 8, 1)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}
	if renderer.WavFile == nil {
		t.Error("Expected WavFile to be non-nil")
	}
}
