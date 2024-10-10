package cmd

import (
	"fmt"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_getUniqueName(t *testing.T) {
	t.Parallel()

	formatTime := func() string {
		return time.Now().Format("2006-01-02_15-04-05")
	}

	tests := []struct {
		name    string
		newPath string
		want    string
	}{
		{
			name:    "Success",
			newPath: "/test_dir/file.txt",
			want:    fmt.Sprintf("%sfile_(%s).txt", viper.GetString("out_dir"), formatTime()),
		},
		{
			name:    "No file",
			newPath: "/test_dir/",
			want: fmt.Sprintf("%stest_dir_(%s)", viper.GetString("out_dir"),
				formatTime()),
		},
		{
			name:    "Empty input",
			newPath: "",
			want: fmt.Sprintf("%s_(%s).", viper.GetString("out_dir"),
				formatTime()),
		},
		{
			name:    "Only file",
			newPath: "file.txt",
			want: fmt.Sprintf("%sfile_(%s).txt", viper.GetString("out_dir"),
				formatTime()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := getUniqueName(tt.newPath); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
