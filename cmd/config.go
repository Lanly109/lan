/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"github.com/Lanly109/lan/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var fileName string

const config string = `# use in all command except config
CodePath = "304"

# use in check command
# "all" for all in csv
Room = "304"
NameList = "namelist.csv"

# use in time command
StartTime = "2021-11-17 08:30:00"
EndTime = "2021-11-17 13:00:00"

# use in clean command
SourcePath = "raw_304"
Problems = [ "expr", "live", "number", "power" ]
Extensions = [ ".cpp", ".c", ".pas" ]

# use in moss command
ReviewProblem = "expr"
ReviewUserID = ""
ReviewLanguage = "cc"
ReviewComment = "expr"
ReviewMaxLimit = 10
ReviewExperimental = false
ReviewNumberResult = 250

# ===Don't Edit it if you do NOT know what you are doing=== #

AbnormalLog = "error.log"
Md5File = "checker.hash.csv"
Debug = false`

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate Config file(default config.toml)",
	Long: `In order to reduce duplicate parameter input, creating configuration files is a good way. 
This instruction will create a file in sample toml format.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("filename: %s", fileName)
		if utils.FileExist(fileName) {
			log.Errorf("The Config File %s Exists!", fileName)
			return
		}
		err := utils.WriteFile(fileName, config)
		if err != nil {
			log.Error(err)
		} else {
			log.Infof("Successfully generated configuration file [%s]", fileName)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&fileName, "name", "n", "config.toml", "The config name(default config.toml)")
}
