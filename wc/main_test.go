package main

import (
	"bufio"
	"os"
	"testing"
)

func TestGetStats(t *testing.T) {
	r, _ := os.Open("test.txt")
	reader := bufio.NewReader(r)
	tests := []struct {
		name    string
		reader  *bufio.Reader
		want    FileStats
		wantErr bool
	}{
		{
			name:    "Passing case",
			reader:  reader,
			want:    FileStats{ByteCount: 342190, LineCount: 7145, WordCount: 58164, CharacterCount: 339292},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := GetStats(tt.reader)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetStats() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetStats() succeeded unexpectedly")
			}
			if got.ByteCount != tt.want.ByteCount ||
				got.CharacterCount != tt.want.CharacterCount ||
				got.LineCount != tt.want.LineCount ||
				got.WordCount != tt.want.WordCount {
				t.Errorf("GetStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
