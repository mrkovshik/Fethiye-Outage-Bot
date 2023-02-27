package util

import (
	"reflect"
	"testing"
)

func Test_Normalize(t *testing.T) {

	tests := []struct {
		name    string
		s       string
		want    []string
		wantErr bool
	}{
		{"normal", "fethiye tuzla", []string{"fethiye", "tuzla"}, false},
		{"caps", "FETHIYE TUZLA", []string{"fethiye", "tuzla"}, false},
		{"accents", "fethiye yeşilüzümlü", []string{"fethiye", "yesiluzumlu"}, false},
		{"scpecial symbols", "fethiye/ \\tuzla!", []string{"fethiye", "tuzla"}, false},
		{"a lot of scpecial symbols", "!#@$%^/\\&#^#!a", []string{"a"}, false},
		{"two worded mahalesi", "Milas Firuzpaşa-Gazipaşa", []string{"milas", "firuzpasa", "gazipasa"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Normalize(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Normalize error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tNormalize = %v, want %v", got, tt.want)
			}
		})
	}
}
