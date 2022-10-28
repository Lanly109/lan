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
	startTimeStr string
	endTimeStr   string
	abnormalLog  string
	startTime    time.Time
	endTime      time.Time
	abnormalList []Data
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

		if fileInfo.ModTime().Before(startTime) || fileInfo.ModTime().After(endTime) {
			abnor, err := ResolvePath(path[codePathLen+1:])
			if err != nil {
				log.Error(err)
				return nil
			}

			abnor.ModifyTime = fileInfo.ModTime()
			log.Warn(abnor.String())
			abnormalList = append(abnormalList, abnor)
		}
	}
	return nil
}

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time <CodePath>",
	Short: "Check modify time of code",
	Long:  `Check whether the modify time of code is in the competition during.`,
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

		startTimeStr = viper.GetString("StartTime")
		endTimeStr = viper.GetString("EndTime")
		abnormalLog = viper.GetString("AbnormalLog")

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
	rootCmd.AddCommand(timeCmd)

	timeCmd.Flags().StringVarP(&startTimeStr, "starttime", "s", "0000-00-00 00:00:00", "Start Time of competition")
	timeCmd.Flags().StringVarP(&endTimeStr, "endtime", "e", "9999-12-31 23:59:59", "End Time of competition")
	timeCmd.Flags().StringVarP(&abnormalLog, "abnormallog", "a", "abnormal.log", "List of exception modification times")
}
