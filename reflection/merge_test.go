package reflection

import (
	"reflect"
	"testing"
	"time"
)

type sampleConfig struct {
	Name       string
	Enabled    bool
	MaxSize    int
	Timeout    time.Duration
	Nested     *string
	SliceField []string
}

func TestMergeZeroFields(t *testing.T) {
	defaultVal := "default"
	base := &sampleConfig{
		Name:       "Existing",
		Enabled:    true,
		MaxSize:    100,
		Timeout:    10 * time.Second,
		Nested:     &defaultVal,
		SliceField: []string{"keep"},
	}

	override := &sampleConfig{
		Name:       "",
		Enabled:    false,
		MaxSize:    0,
		Timeout:    0,
		Nested:     nil,
		SliceField: nil,
	}

	MergeZeroFields(override, base)

	if !reflect.DeepEqual(override, base) {
		t.Errorf("Expected merged config to match base, got: %+v", override)
	}
}
