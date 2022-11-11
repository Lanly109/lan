/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"os"
	"testing"

	"github.com/Lanly109/lan/utils"
	"github.com/matryer/is"
)

func Test_ShareCommand(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name     string
		args     []string
		err      error
		fileName string
		want     string
	}{
		{
			name: "case 1",
			args: []string{"gen", "share"},
			err:  nil,
			want: "sharing.cpp",
		},
		{

			name:     "case 2",
			args:     []string{"gen", "share"},
			err:      nil,
			fileName: "qwq.cpp",
			want:     "qwq.cpp",
		},
	}
	for _, tt := range tests {
		rootCmd.SetArgs(tt.args)

		t.Run(tt.name, func(t *testing.T) {
			if tt.fileName != "" {
				scriptName = tt.fileName
			}

			got := rootCmd.Execute()
			is.Equal(got, tt.err)

			data, err := utils.ReadFile(tt.want)

			is.NoErr(err)

			if got == nil {
				is.Equal(data, sharingCode)
			}
		})

		os.Remove(tt.want)
	}
}
