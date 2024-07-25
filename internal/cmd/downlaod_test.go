package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTfilePath(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
		err      error
	}{
		{"ValidTorrentFile", []string{"file.torrent"}, "file.torrent", nil},
		{"InvalidArgument", []string{}, "", errors.New("Invalid argument")},
		{"NotTorrentFile", []string{"file.txt"}, "", errors.New("It is not torrent file!")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTfilePath(tt.args)
			assert.Equal(t, tt.err, err, "error mismatch")
			assert.Equal(t, tt.expected, got, "output mismatch")
		})
	}
}
