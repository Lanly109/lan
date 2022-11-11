/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func Test_cleanCommand(t *testing.T) {
	is := is.New(t)
	codePath = filepath.Join(testPath, "304")
	srcPath = filepath.Join(testPath, "raw_304")

	tests := []struct {
		name     string
		args     []string
		problems []string
		err      error
		copyNum  int
	}{
		{
			name: "case 1",
			args: []string{"clean"},
			err:  errors.New("Requires args of Code Path"),
		},
		{

			name: "case 2",
			args: []string{"clean", codePath},
			err:  errors.New("Requires args of Source Path"),
		},
		{
			name:     "case 3",
			args:     []string{"clean", codePath, srcPath},
			problems: []string{"expr", "live", "number", "power"},
			copyNum:  18,
			err:      nil,
		},
		{
			name:     "case 4",
			args:     []string{"clean", codePath, srcPath},
			problems: []string{"expr"},
			copyNum:  4,
			err:      nil,
		},
		{
			name:     "case 5",
			args:     []string{"clean", codePath, srcPath},
			problems: []string{"live", "power"},
			copyNum:  10,
			err:      nil,
		},
		{
			name:     "case 6",
			args:     []string{"clean", codePath, srcPath},
			problems: []string{"number"},
			copyNum:  4,
			err:      nil,
		},
	}
	for _, tt := range tests {
		rootCmd.SetArgs(tt.args)

		problemList = tt.problems

		t.Run(tt.name, func(t *testing.T) {
			got := rootCmd.Execute()
			is.Equal(got, tt.err)

			if got == nil {
				is.Equal(copyNum, tt.copyNum)
			}
		})
	}
}
