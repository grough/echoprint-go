package echo

import (
	"os"
	"testing"
)

func TestNewRenderer(t *testing.T) {
	inputFile := "testdata/mono.wav"
	outputFile := "ignore/test.wav"
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		t.Skipf("Test WAV file %s not found, skipping", inputFile)
	}

	renderer, err := NewRenderer(inputFile, outputFile, 120, 8, 1)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}
	if renderer.InputDecoder == nil {
		t.Error("Expected WavFile to be non-nil")
	}
}
