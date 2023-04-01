/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	. "github.com/Lanly109/lan/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	startTimeStr    string
	endTimeStr      string
	abnormalLog     string
	sourceSizeLimit int64
	startTime       time.Time
	endTime         time.Time
	abnormalList    []Data
)

func visit(path string, info fs.DirEntry, err error) error {

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

		if fileInfo.ModTime().Before(startTime) || fileInfo.ModTime().After(endTime) || fileInfo.Size() > sourceSizeLimit {
			abnor, err := ResolvePath(path[codePathLen+1:])
			if err != nil {
				log.Error(err)
				return nil
			}

			abnor.ModifyTime = fileInfo.ModTime()
			abnor.Size = fileInfo.Size()
			log.Warn(abnor.String())
			abnormalList = append(abnormalList, abnor)
		}
	}
	return nil
}

// validCmd represents the time command
var validCmd = &cobra.Command{
	Use:   "valid <CodePath>",
	Short: "Check modify time and size of code",
	Long:  `Check whether the modify of code is in the competition during and the size of code is in limit.`,
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
		viper.BindPFlag("StartTime", cmd.Flags().Lookup("starttime"))
		viper.BindPFlag("EndTime", cmd.Flags().Lookup("endtime"))
		viper.BindPFlag("AbnormalLog", cmd.Flags().Lookup("abnormallog"))
		viper.BindPFlag("SourceSizeLimit", cmd.Flags().Lookup("sizelimit"))

		startTimeStr = viper.GetString("StartTime")
		endTimeStr = viper.GetString("EndTime")
		abnormalLog = viper.GetString("AbnormalLog")
		sourceSizeLimit = viper.GetInt64("SourceSizeLimit")

		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			loc = time.FixedZone("CST", 8*3600)
		}
		startTime, _ = time.ParseInLocation(TimeTemplate, startTimeStr, loc)
		endTime, _ = time.ParseInLocation(TimeTemplate, endTimeStr, loc)

		codePath = filepath.Clean(codePath)
		codePathLen = len(codePath)

		log.Info("CodePath: ", codePath)
		log.Info("StartTime: ", startTime.Format(TimeTemplate))
		log.Info("EndTime: ", endTime.Format(TimeTemplate))
		log.Info("AbnormalLog: ", abnormalLog)
		log.Info("SourceSizeLimit: ", sourceSizeLimit, " bytes")
	},
	Run: func(cmd *cobra.Command, args []string) {
		abnormalList = []Data{}
		filepath.WalkDir(codePath, visit)

		file, err := os.OpenFile(abnormalLog, os.O_CREATE|os.O_RDWR, 0774)
		if err != nil {
			log.Error("write abnormal error: ", err)
			return
		}
		for _, abnormal := range abnormalList {
			fmt.Fprintln(file, abnormal.Output()) // 向file对应文件中写入数据
		}

		log.Info("Abnormal Log save in ", abnormalLog)
	},
}

func init() {
	rootCmd.AddCommand(validCmd)

	validCmd.Flags().StringVarP(&startTimeStr, "starttime", "s", "0000-00-00 00:00:00", "Start Time of competition")
	validCmd.Flags().StringVarP(&endTimeStr, "endtime", "e", "9999-12-31 23:59:59", "End Time of competition")
	validCmd.Flags().Int64Var(&sourceSizeLimit, "sizelimit", 1024 * 100, "Source Code size limit")
	validCmd.Flags().StringVarP(&abnormalLog, "abnormallog", "a", "abnormal.log", "List of exception modification times")
}
