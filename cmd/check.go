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
	nameList                  string
	codePath                  string
	codePathLen               int
	room                      string
	absentContanstants        mapset.Set
	knownAbsentContanstants   mapset.Set
	unknownAbsentContanstants mapset.Set
	extraContanstants         mapset.Set
	knownAbsentButAcutualNot  mapset.Set
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

		expectContestants, knownAbsentContestants := getContanstansFromCsv()
		realContestants := getContanstansFromCodePath()

		log.Debug("Expected Contestant: ", utils.SetToOrderStringSlice(expectContestants))
		log.Debug("Known Absent Contestant: ", utils.SetToOrderStringSlice(knownAbsentContestants))
		log.Debug("Real Contestant: ", utils.SetToOrderStringSlice(realContestants))

		absentContanstants = expectContestants.Difference(realContestants)
		knownAbsentButAcutualNot = knownAbsentContestants.Difference(absentContanstants)
		if knownAbsentButAcutualNot.Cardinality() != 0 {
			log.Error("Known Absent But actual NOT absent: ", utils.SetToOrderStringSlice(knownAbsentButAcutualNot))
		}

		knownAbsentContanstants = absentContanstants.Intersect(knownAbsentContestants)
		unknownAbsentContanstants = absentContanstants.Difference(knownAbsentContestants)
		extraContanstants = realContestants.Difference(expectContestants)

		log.Info("Absent Contestants: ", utils.SetToOrderStringSlice(absentContanstants))
		log.Info("Known Absent Contestants: ", utils.SetToOrderStringSlice(knownAbsentContanstants))
		log.Warn("Unknown Absent Contestants: ", utils.SetToOrderStringSlice(unknownAbsentContanstants))
		log.Warn("Contestants should not in this room: ", utils.SetToOrderStringSlice(extraContanstants))

		log.Infof("Required attendance:\t%d people", expectContestants.Cardinality())
		log.Infof("Actual attendance:\t%d people", realContestants.Cardinality()-extraContanstants.Cardinality())
	},
}

func getContanstansFromCsv() (mapset.Set, mapset.Set) {
	datas := utils.ReadCsv(nameList)

	log.Debug("NameList: ", datas)

	var filteData, absentData []interface{}
	for _, data := range datas {
		if data[1] == room || room == "all" {
			filteData = append(filteData, data[0])
			if len(data) > 2 && data[2] == "0" {
				absentData = append(absentData, data[0])
			}
		}
	}

	return mapset.NewSetFromSlice(filteData), mapset.NewSetFromSlice(absentData)
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
