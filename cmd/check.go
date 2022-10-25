/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/Lanly109/lan/utils"

	mapset "github.com/deckarep/golang-set"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	nameList           string
	codePath           string
	codePathLen        int
	room               string
	absentContanstants mapset.Set
	extraContanstants  mapset.Set
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check <CodePath>",
	Short: "Check namelist in a room",
	Long:  `Check whether the candidate's file name is missing or extra according to the room namelist.`,
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
		viper.BindPFlag("NameList", cmd.Flags().Lookup("namelist"))
		viper.BindPFlag("Room", cmd.Flags().Lookup("room"))

		nameList = viper.GetString("NameList")
		room = viper.GetString("Room")
		codePath = filepath.Clean(codePath)

		log.Info("NameList: ", nameList)
		log.Info("CodePath: ", codePath)
		log.Info("Room: ", room)
	},
	Run: func(cmd *cobra.Command, args []string) {

		expectContestants := getContanstansFromCsv()
		realContestants := getContanstansFromCodePath()

		log.Debug("Expected Contestant: ", expectContestants.String())
		log.Debug("Real Contestant: ", realContestants.String())

		absentContanstants = expectContestants.Difference(realContestants)
		extraContanstants = realContestants.Difference(expectContestants)

		log.Warn("Contestant without submissions: ", absentContanstants)
		log.Warn("Contestant should not in this room: ", extraContanstants)
	},
}

func getContanstansFromCsv() mapset.Set {
	datas := utils.ReadCsv(nameList)

	log.Debug("NameList: ", datas)

	var filteData []interface{}
	for _, data := range datas {
		if data[1] == room || room == "all" {
			filteData = append(filteData, data[0])
		}
	}

	return mapset.NewSetFromSlice(filteData)
}

func getContanstansFromCodePath() mapset.Set {
	var filteData []interface{}

	log.Debug("Reading code path: ", codePath)
	dirs, err := ioutil.ReadDir(codePath)

	if err != nil {
		log.Error(err)
		cobra.CheckErr(err)
	}

	for _, dir := range dirs {
		log.WithFields(log.Fields{
			"Name": dir.Name(),
			"Dir":  dir.IsDir(),
		}).Debug("")
		if dir.IsDir() {
			filteData = append(filteData, dir.Name())
		}
	}

	return mapset.NewSetFromSlice(filteData)
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVarP(&nameList, "namelist", "l", "./namelist.csv", "path to namelist")
	checkCmd.Flags().StringVarP(&room, "room", "r", "all", "checking room, all default")
}
