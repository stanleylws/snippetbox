package main

import (
	"testing"
	"time"

	"snippetbox.stanley.net/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm time.Time
		want string
	}{
		{
			name: "UTC",
			tm: time.Date(2024, 11, 27, 10, 15, 0, 0, time.UTC),
			want: "27 Nov 2024 at 10:15",
		},
		{
			name: "Empty",
			tm: time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm: time.Date(2024, 11, 27, 10, 15, 0, 0, time.FixedZone("UTC+8", 8*60*60)),
			want: "27 Nov 2024 at 02:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}
}