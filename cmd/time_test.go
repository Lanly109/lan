/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func Test_timeCommand(t *testing.T) {
	is := is.New(t)
	codePath = filepath.Join(testPath, "clean_304")

	tests := []struct {
		name           string
		args           []string
		err            error
		startTimeStr   string
		endTimeStr     string
		abnormalNumber int
		abnormalLog    string
		want           string
	}{
		{
			name: "case 1",
			args: []string{"time"},
			err:  errors.New("Requires args of Code Path"),
		},
		{

			name:           "case 2",
			args:           []string{"time", codePath},
			err:            nil,
			startTimeStr:   "2021-11-17 08:30:00",
			endTimeStr:     "2021-11-17 13:00:00",
			abnormalNumber: 2,
			want:           "abnormal.log",
		},
		{

			name:           "case 3",
			args:           []string{"time", codePath},
			err:            nil,
			startTimeStr:   "2021-11-17 08:30:00",
			endTimeStr:     "2021-11-17 13:00:00",
			abnormalNumber: 2,
			abnormalLog:    "qwq.log",
			want:           "qwq.log",
		},
	}
	for _, tt := range tests {
		rootCmd.SetArgs(tt.args)

		t.Run(tt.name, func(t *testing.T) {
			startTimeStr = tt.startTimeStr
			endTimeStr = tt.endTimeStr

			if tt.abnormalLog != "" {
				abnormalLog = tt.abnormalLog
			}

			got := rootCmd.Execute()
			is.Equal(got, tt.err)

			if got == nil {
				is.Equal(len(abnormalList), tt.abnormalNumber)

				_, err := os.Stat(tt.want)

				is.NoErr(err)

				os.Remove(abnormalLog)
			}
		})
	}
}
