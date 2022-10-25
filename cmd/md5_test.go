/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Lanly109/lan/utils"
	"github.com/matryer/is"
)

func Test_md5Command(t *testing.T) {
	is := is.New(t)
	codePath = filepath.Join(testPath, "clean_304")

	tests := []struct {
		name     string
		args     []string
		err      error
		md5      string
		md5File  string
		wantFile string
		want     string
	}{
		{
			name: "case 1",
			args: []string{"md5"},
			err:  errors.New("Requires args of Code Path"),
		},
		{
			name:     "case 2",
			args:     []string{"md5", codePath},
			err:      nil,
			wantFile: "checker.hash",
			md5:      "adb176103b113909144795dd33dd22a6",
			want: `GD-00032,expr,e68166389de01200c6d970fc3be10a4b
GD-00032,live,9aac92da434f99f644e7740db50eca25
GD-00032,number,51e45352120bc7067179481c5c56a350
GD-00032,power,29e1307417b6144ca247fd974d7bd95d
GD-00077,expr,9020723a27fbcc726f10f6ff2306688c
GD-00077,live,5c279a6a8af486c3e5e4929a42f1f175
GD-00077,number,28472492102c2795772612699e2a0a0f
GD-00077,power,dc6b9790f84d698f6e257ea382a38b93
GD-00081,expr,3476863bccc3adf7db7aa63cf06bc6bd
GD-00081,live,edeb6ec1f9fdf2c7741511c86e85b69e
GD-00081,power,aa930e7b907ad924ce34bb800bcf17f7
GD-00111,expr,2eed8e9c2f2938d139dad0fbd00706b1
GD-00111,live,912390dee0225056e11e812035666817
GD-00111,number,c784af6c7c10fadeb29f496ce6197efb
GD-00111,power,f1848fca1415245270e83096cea88bdf
GD-00139,expr,b5f9dc19c1e093d35c31ab0002cab3b5
GD-00139,live,1c0a55027ce93cb3d4db976b740e50b8
GD-00139,number,2e30638f33f913145e4e129eac238163
GD-00139,power,6a733844e618c38f366f9a149a6c68a0
GD-00142,expr,6a8e5c7a4c0eb6197d2de77a8f663d89
GD-00142,live,abcbb58bf383db2954ca669288cbed0a
GD-00142,number,ae2cd346028e80ea22de88d711e47ee7
GD-00142,power,3828e142aaa229ab94a711d510cbe6a2
GD-00192,expr,06c81b60f732b5e676f5094c9574d6eb
GD-00192,live,13a0063b0e624ca25f6aaec813ef9bb6
GD-00192,number,c5b8cdc8deabd20422a3762c1ee855d2
GD-00192,power,3f643c3f3785ce045e070c1c6faac8eb
`,
		},
		{
			name:     "case 3",
			args:     []string{"md5", codePath},
			err:      nil,
			md5File:  "qwq.hash",
			wantFile: "qwq.hash",
			md5:      "adb176103b113909144795dd33dd22a6",
			want: `GD-00032,expr,e68166389de01200c6d970fc3be10a4b
GD-00032,live,9aac92da434f99f644e7740db50eca25
GD-00032,number,51e45352120bc7067179481c5c56a350
GD-00032,power,29e1307417b6144ca247fd974d7bd95d
GD-00077,expr,9020723a27fbcc726f10f6ff2306688c
GD-00077,live,5c279a6a8af486c3e5e4929a42f1f175
GD-00077,number,28472492102c2795772612699e2a0a0f
GD-00077,power,dc6b9790f84d698f6e257ea382a38b93
GD-00081,expr,3476863bccc3adf7db7aa63cf06bc6bd
GD-00081,live,edeb6ec1f9fdf2c7741511c86e85b69e
GD-00081,power,aa930e7b907ad924ce34bb800bcf17f7
GD-00111,expr,2eed8e9c2f2938d139dad0fbd00706b1
GD-00111,live,912390dee0225056e11e812035666817
GD-00111,number,c784af6c7c10fadeb29f496ce6197efb
GD-00111,power,f1848fca1415245270e83096cea88bdf
GD-00139,expr,b5f9dc19c1e093d35c31ab0002cab3b5
GD-00139,live,1c0a55027ce93cb3d4db976b740e50b8
GD-00139,number,2e30638f33f913145e4e129eac238163
GD-00139,power,6a733844e618c38f366f9a149a6c68a0
GD-00142,expr,6a8e5c7a4c0eb6197d2de77a8f663d89
GD-00142,live,abcbb58bf383db2954ca669288cbed0a
GD-00142,number,ae2cd346028e80ea22de88d711e47ee7
GD-00142,power,3828e142aaa229ab94a711d510cbe6a2
GD-00192,expr,06c81b60f732b5e676f5094c9574d6eb
GD-00192,live,13a0063b0e624ca25f6aaec813ef9bb6
GD-00192,number,c5b8cdc8deabd20422a3762c1ee855d2
GD-00192,power,3f643c3f3785ce045e070c1c6faac8eb
`,
		},
	}
	for _, tt := range tests {
		rootCmd.SetArgs(tt.args)

		t.Run(tt.name, func(t *testing.T) {
			if tt.md5File != "" {
				md5File = tt.md5File
			}

			got := rootCmd.Execute()
			is.Equal(got, tt.err)

			if got == nil {
				is.Equal(hex.EncodeToString(md5Totle.Sum(nil)), tt.md5)

				_, err := os.Stat(tt.wantFile)
				is.NoErr(err)

				data, _ := utils.ReadFile(tt.wantFile)

				is.Equal(data, tt.want)

				os.Remove(md5File)
			}
		})
	}

}
