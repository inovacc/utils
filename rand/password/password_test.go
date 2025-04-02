package password

import (
	"strings"
	"testing"
)

func TestNewPassword(t *testing.T) {
	p := NewPassword(
		WithLength(16),
		WithNumbers(true),
		WithSpecial(true),
		WithLower(true),
		WithUpper(true),
	)

	pass, err := p.Generate()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Generated password: %s", pass)
}

func TestPassword_Generate_SingleCharTypes(t *testing.T) {
	tests := []struct {
		name        string
		options     []Option
		expectedSet string
	}{
		{
			name: "Only numbers",
			options: []Option{
				WithLength(12),
				WithNumbers(true),
				WithSpecial(false),
				WithLower(false),
				WithUpper(false),
			},
			expectedSet: passwordOptions["num"],
		},
		{
			name: "Only upper case",
			options: []Option{
				WithLength(12),
				WithNumbers(false),
				WithSpecial(false),
				WithLower(false),
				WithUpper(true),
			},
			expectedSet: passwordOptions["upperCase"],
		},
		{
			name: "Only lower case",
			options: []Option{
				WithLength(12),
				WithNumbers(false),
				WithSpecial(false),
				WithLower(true),
				WithUpper(false),
			},
			expectedSet: passwordOptions["lowerCase"],
		},
		{
			name: "Only special characters",
			options: []Option{
				WithLength(12),
				WithNumbers(false),
				WithSpecial(true),
				WithLower(false),
				WithUpper(false),
			},
			expectedSet: passwordOptions["specialChar"],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPassword(tt.options...)
			password, err := p.Generate()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(password) != 12 {
				t.Errorf("expected password length 12, got %d", len(password))
			}
			for _, r := range password {
				if !strings.ContainsRune(tt.expectedSet, r) {
					t.Errorf("character %q not in expected set %q", r, tt.expectedSet)
				}
			}
			t.Logf("Generated password: %s", password)
		})
	}
}
