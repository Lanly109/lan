/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"errors"
	"path/filepath"
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/matryer/is"
)

const (
	testPath = "../demo"
)

func Test_getContanstansFromCsv(t *testing.T) {
	verbose = true
	is := is.New(t)
	nameList = filepath.Join(testPath, "namelist.csv")

	tests := []struct {
		name       string
		want       mapset.Set
		wantAbsent mapset.Set
		room       string
	}{
		{
			name:       "case 1",
			want:       mapset.NewSet("GD-00032", "GD-00077", "GD-00081", "GD-00111", "GD-00139", "GD-00192"),
			wantAbsent: mapset.NewSet("GD-00111", "GD-00139"),
			room:       "304",
		},
		{
			name:       "case 2",
			want:       mapset.NewSet("GD-00018", "GD-00062", "GD-00128", "GD-00153", "GD-00291"),
			wantAbsent: mapset.NewSet("GD-00062"),
			room:       "402",
		},
		{
			name:       "case 3",
			want:       mapset.NewSet("GD-00018", "GD-00032", "GD-00062", "GD-00077", "GD-00081", "GD-00111", "GD-00128", "GD-00139", "GD-00153", "GD-00192", "GD-00291"),
			wantAbsent: mapset.NewSet("GD-00062", "GD-00111", "GD-00139"),
			room:       "all",
		},
	}
	for _, tt := range tests {
		room = tt.room
		t.Run(tt.name, func(t *testing.T) {
			got, absent := getContanstansFromCsv()
			is.Equal(got, tt.want)
			is.Equal(absent, tt.wantAbsent)
		})
	}
}

func Test_getContanstansFromCodePath(t *testing.T) {
	is := is.New(t)
	codePath = filepath.Join(testPath, "clean_304")

	tests := []struct {
		name string
		want mapset.Set
	}{
		{
			name: "case 1",
			want: mapset.NewSet("GD-00018", "GD-00032", "GD-00077", "GD-00081", "GD-00129", "GD-00139"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getContanstansFromCodePath()
			is.Equal(got, tt.want)
		})
	}
}

func Test_checkCommand(t *testing.T) {
	is := is.New(t)
	codePath = filepath.Join(testPath, "clean_304")

	tests := []struct {
		name              string
		args              []string
		room              string
		err               error
		absent            mapset.Set
		knownAbsent       mapset.Set
		knownAbsentButNot mapset.Set
		unknownAbsent     mapset.Set
		extra             mapset.Set
	}{
		{
			name:              "case 1",
			args:              []string{"check", codePath},
			room:              "304",
			err:               nil,
			absent:            mapset.NewSet("GD-00111", "GD-00192"),
			knownAbsentButNot: mapset.NewSet("GD-00139"),
			knownAbsent:       mapset.NewSet("GD-00111"),
			unknownAbsent:     mapset.NewSet("GD-00192"),
			extra:             mapset.NewSet("GD-00018", "GD-00129"),
		},
		{
			name:              "case 2",
			args:              []string{"check", codePath},
			room:              "all",
			err:               nil,
			absent:            mapset.NewSet("GD-00062", "GD-00111", "GD-00128", "GD-00153", "GD-00192", "GD-00291"),
			knownAbsentButNot: mapset.NewSet("GD-00139"),
			knownAbsent:       mapset.NewSet("GD-00062", "GD-00111"),
			unknownAbsent:     mapset.NewSet("GD-00128", "GD-00153", "GD-00192", "GD-00291"),
			extra:             mapset.NewSet("GD-00129"),
		},
		{
			name: "case 3",
			args: []string{"check"},
			err:  errors.New("Requires args of Code Path"),
		},
	}
	for _, tt := range tests {
		rootCmd.SetArgs(tt.args)

		room = tt.room

		t.Run(tt.name, func(t *testing.T) {
			got := rootCmd.Execute()
			is.Equal(got, tt.err)

			if got == nil {
				is.Equal(absentContanstants, tt.absent)
				is.Equal(knownAbsentButAcutualNot, tt.knownAbsentButNot)
				is.Equal(knownAbsentContanstants, tt.knownAbsent)
				is.Equal(unknownAbsentContanstants, tt.unknownAbsent)
				is.Equal(extraContanstants, tt.extra)
			}
		})
	}
}
