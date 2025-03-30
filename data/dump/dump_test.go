package dump

import (
	"os"
	"testing"
)

func TestDump(t *testing.T) {
	data := struct {
		Name    string
		Age     int
		Balance float64
	}{
		Name:    "John Doe",
		Age:     30,
		Balance: 1000.50,
	}

	if err := Dump("test.pgo", data); err != nil {
		t.Errorf("Failed to dump data: %v", err)
		return
	}

	var loadedData struct {
		Name    string
		Age     int
		Balance float64
	}

	if err := Load("test.pgo", &loadedData); err != nil {
		t.Errorf("Failed to load data: %v", err)
		return
	}

	if loadedData.Name != data.Name || loadedData.Age != data.Age {
		t.Errorf("Loaded data does not match original data: got %v, want %v", loadedData, data)
		return
	}

	if err := os.Remove("test.pgo"); err != nil {
		t.Errorf("Failed to remove test file: %v", err)
		return
	}
}
