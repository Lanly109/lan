/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"hash"
	"io/fs"
	"path/filepath"

	. "github.com/Lanly109/lan/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	md5File     string
	md5DataList []Data
	md5Totle    hash.Hash
)

func geneMd5(path string, info fs.DirEntry, err error) error {
	if err != nil {
		log.Error(err)
		return err
	}

	if info.IsDir() {
		log.WithField("dir", info.IsDir()).Debug(path)
		return nil
	}

	if fileInfo, err := info.Info(); err != nil {
		log.Error(err)
		return nil

	} else {
		log.WithField("path", path).Debug(fileInfo.ModTime().Format(TimeTemplate))

		file, err := ResolvePath(path[codePathLen+1:])
		if err != nil {
			log.Error(err)
			return nil
		}

		if file.Md5, err = Md5(path); err != nil {
			log.Error(err)
			return nil
		}

		md5Totle.Write([]byte(file.Md5))

		md5DataList = append(md5DataList, file)
	}
	return nil
}

// md5Cmd represents the md5 command
var md5Cmd = &cobra.Command{
	Use:   "md5 <CodePath>",
	Short: "Generate md5 of each file in folds",
	Long:  `A small tools to generate MD5 to each file in folds, and output as a csv file`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			codePath = viper.GetString("CodePath")
			if codePath == "" {
				err := errors.New("Requires args of Code Path")
				log.Error(err)
				return err
			}
			return nil
		}
		codePath = args[0]
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("Md5File", cmd.Flags().Lookup("md5file"))

		md5File = viper.GetString("Md5File")
		md5Totle = md5.New()

		codePath = filepath.Clean(codePath)
		codePathLen = len(codePath)
		log.Info("CodePath: ", codePath)
		log.Info("Md5File: ", md5File)
	},
	Run: func(cmd *cobra.Command, args []string) {
		md5DataList = []Data{}
		filepath.WalkDir(codePath, geneMd5)

		var datas [][]string

		for _, md5Data := range md5DataList {
			if len(md5Data.Name) != 0 {
				datas = append(datas, []string{
					md5Data.Name,
					md5Data.Problem,
					md5Data.Md5,
				})
			}
		}

		WriteCsv(md5File, datas)

		log.Info("Totle Md5: ", hex.EncodeToString(md5Totle.Sum(nil)))
		log.Info("md5 csv save in ", md5File)
	},
}

func init() {
	rootCmd.AddCommand(md5Cmd)

	md5Cmd.Flags().StringVarP(&md5File, "md5file", "", "checker.hash", "md5 CSV file name")
}
