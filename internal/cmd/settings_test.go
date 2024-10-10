package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SetOutDir(t *testing.T) {
	t.Parallel()

	setConfigFile()

	err := os.MkdirAll(viper.GetString("HOME")+"/test-directory/test/", 0o0777)
	if err != nil {
		require.Error(t, err)
	}

	t.Cleanup(func() { os.RemoveAll(viper.GetString("HOME") + "/test-directory/test/") })

	cases := []struct {
		name         string
		args         []string
		Error        error
		expectedPath string
		isErrCase    bool
	}{
		{
			name:         "Success",
			args:         []string{viper.GetString("HOME") + "/test-directory/test/"},
			expectedPath: viper.GetString("HOME") + "/test-directory/test/",
		},
		{
			name:         "Empty path",
			args:         []string{""},
			expectedPath: viper.GetString("out_dir"),
		},
		{
			name:         "Many slash",
			args:         []string{viper.GetString("HOME") + "///test-directory///test///"},
			expectedPath: viper.GetString("HOME") + "/test-directory/test/",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cmd := &cobra.Command{}

			SetOutDir(cmd, tt.args)

			outDir := viper.GetString("out_dir")
			assert.Equal(t, tt.expectedPath, outDir, "Expected output path does not match")
		})
	}
}

func Test_ValidateOutputPath(t *testing.T) {
	t.Parallel()

	viper.AutomaticEnv()

	f, err := os.Create(viper.GetString("HOME") + "/file.txt")
	if err != nil {
		require.Error(t, err)
	}

	t.Cleanup(func() { os.Remove(viper.GetString("HOME") + "/file.txt") })

	if f.Close() != nil {
		require.Error(t, err)
	}

	cases := []struct {
		name   string
		path   string
		expErr error
	}{
		{
			name:   "Success",
			path:   viper.GetString("HOME") + "/",
			expErr: nil,
		},
		{
			name:   "Invalid path",
			path:   viper.GetString("HOME") + "/folder_than_doesn't_exist/",
			expErr: ErrPathNotExist,
		},
		{
			name:   "Path with file",
			path:   viper.GetString("HOME") + "/file.txt",
			expErr: ErrPathNotDir,
		},
		{
			name:   "No end slash",
			path:   viper.GetString("HOME") + "/Видео",
			expErr: nil,
		},
		{
			name:   "Many slash",
			path:   fmt.Sprintf("////%s///Видео/////", viper.GetString("HOME")),
			expErr: ErrPathNotIsHome,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := validateOutputPath(tc.path)
			assert.Equal(t, tc.expErr, err)
		})
	}
}

func TestIsOneArgument(t *testing.T) {
	cases := []struct {
		name        string
		args        []string
		expectedVal bool
		expectedErr error
	}{
		{
			name:        "Success",
			args:        []string{"one"},
			expectedVal: true,
			expectedErr: nil,
		},
		{
			name:        "Empty args",
			args:        []string{},
			expectedVal: false,
			expectedErr: nil,
		},
		{
			name:        "Many Args",
			args:        []string{"one", "two"},
			expectedVal: false,
			expectedErr: ErrOneArgument,
		},
	}

	t.Parallel()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			isOne, err := isOneArgument(tc.args)

			t.Parallel()
			assert.Equal(t, tc.expectedVal, isOne)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
