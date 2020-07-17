package dirscanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirScannerWithEntries(t *testing.T) {
	scanner := NewScanner()
	includes := []string{`.*\.crt`}
	excludes := []string{`.*\\node_modules`}

	if err := scanner.AddNeedle("crt", includes, excludes); err != nil {
		t.Errorf("ERROR crt regex")
	}

	includes = []string{`.*\.csproj`}
	excludes = []string{`.*\\bin`}

	if err := scanner.AddNeedle("nuget", includes, excludes); err != nil {
		t.Errorf("ERROR csproj regex ")
	}

	if results, err := scanner.Scan("test_data"); err != nil {
		t.Errorf("ERROR")
	} else if len(*results) != 2 {
		t.Errorf("Invalid length returned")
	}
}

func TestDirScannerWithNoEntries(t *testing.T) {
	scanner := NewScanner()

	_, err := scanner.Scan("test_data")

	assert.Error(t, err)
}

func TestDirScannerWithInvalidSrcPath(t *testing.T) {
	scanner := NewScanner()
	includes := []string{`.*\.crt`}
	excludes := []string{`.*\\node_modules`}

	if err := scanner.AddNeedle("crt", includes, excludes); err != nil {
		t.Errorf("ERROR crt regex")
	}

	_, err := scanner.Scan("asdfasdfasdffsd")
	assert.Error(t, err)
}
